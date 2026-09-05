package server

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"google.golang.org/grpc"

	pb "HailowSellerService/HailowProto/build/go/SellerService/v1"
	"HailowSellerService/internal/domain"
	redis_repository "HailowSellerService/internal/infrastructure/redis/repository"
	"HailowSellerService/internal/repository"
	"HailowSellerService/internal/transport/grpc/handlers"
	"HailowSellerService/internal/transport/grpc/interceptors"
	"HailowSellerService/internal/usecase/auth"
	"HailowSellerService/internal/usecase/profile"
	"HailowSellerService/pkg/database"
	"HailowSellerService/pkg/logging"
	cache "HailowSellerService/pkg/redis"
	s3_storage "HailowSellerService/pkg/s3"
)

type Server struct {
	grpcServer     *grpc.Server
	debug          bool
	host           string
	port           int
	listener       net.Listener
	postgresClient *database.Postgres
	redisClient    *cache.RedisClient
	s3Client       *s3_storage.S3Client
	logger         *logging.Logger
}

func New(debug bool, host string, port int) (*Server, error) {
	runtime.GOMAXPROCS(runtime.NumCPU())

	var logger *logging.Logger
	var err error

	if debug {
		logger, err = logging.New("development")
	} else {
		logger, err = logging.New("production")
	}

	if err != nil {
		return nil, fmt.Errorf("Failed to initialize logger: %v", err)
	}

	logger.Infof("Initializing gRPC server on %s:%d", host, port)

	var dbURL string
	var redisHost string
	var redisPort string

	if debug {
		dbURL = os.Getenv("DEV_DATABASE_URL")
		redisHost = os.Getenv("DEV_REDIS_HOST")
		redisPort = os.Getenv("DEV_REDIS_PORT")
	} else {
		dbURL = os.Getenv("DATABASE_URL")
		redisHost = os.Getenv("REDIS_HOST")
		redisPort = os.Getenv("REDIS_PORT")
	}

	if redisHost == "" || redisPort == "" {
		return nil, fmt.Errorf("Redis host or port is not set")
	}

	postgresClient, err := database.NewPostgresClient(dbURL, logger)

	if err != nil {
		return nil, fmt.Errorf("Failed to connect PostgreSQL database: %v", err)
	}

	sellerRepo := repository.NewSellerRepository(postgresClient)

	redisClient, err := cache.New(redisHost, redisPort, logger)
	if err != nil {
		return nil, fmt.Errorf("Failed to connect Redis client: %v", err)
	}

	sessionRepo := redis_repository.NewSessionRepository(redisClient)

	var s3Client *s3_storage.S3Client
	if accessKeyID := os.Getenv("S3_ACCESS_KEY_ID"); accessKeyID != "" || os.Getenv("S3_SECRET_ACCESS_KEY") != "" || os.Getenv("S3_BUCKET") != "" {
		cfg := s3_storage.Config{
			AccessKeyID:     os.Getenv("S3_ACCESS_KEY_ID"),
			SecretAccessKey: os.Getenv("S3_SECRET_ACCESS_KEY"),
			Region:          os.Getenv("S3_REGION"),
			Bucket:          os.Getenv("S3_BUCKET"),
			Endpoint:        os.Getenv("S3_ENDPOINT"),
		}
		s3Client, err = s3_storage.NewS3Client(context.Background(), cfg, logger)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize S3 client: %w", err)
		}
	}

	authUseCase := auth.NewAuthUseCase(sellerRepo, sessionRepo)
	authHandler := handlers.NewAuthHandler(authUseCase, logger)

	profileUseCase := profile.NewProfileUseCase(sellerRepo)
	profileHandler := handlers.NewProfileHandler(profileUseCase, logger)

	grpcServer := grpc.NewServer(grpc.ChainUnaryInterceptor(
		interceptors.RecoveryInterceptor(logger),
		interceptors.RefreshTokenInterceptor(logger),
		interceptors.AuthInterceptor(logger),
		interceptors.LoggingInterceptor(logger),
	))

	pb.RegisterSellerServiceServer(grpcServer, authHandler)
	pb.RegisterSellerProfileServiceServer(grpcServer, profileHandler)

	return &Server{
		grpcServer:     grpcServer,
		debug:          debug,
		host:           host,
		port:           port,
		logger:         logger,
		postgresClient: postgresClient,
		redisClient:    redisClient,
		s3Client:       s3Client,
	}, nil

}

func (s *Server) Run() error {
	listenAddr := fmt.Sprintf("%s:%d", s.host, s.port)
	listener, err := net.Listen("tcp", listenAddr)
	if err != nil {
		return fmt.Errorf("%s: %v", domain.ErrListenAddress.Message, err)
	}

	s.listener = listener

	s.logger.Infof("gRPC server listening on %s", listener.Addr().String())

	go func() {
		if err := s.grpcServer.Serve(listener); err != nil && err != grpc.ErrServerStopped {
			s.logger.Errorf("%s: %v", domain.ErrGRPCServer.Message, err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	s.logger.Info("Shutdowm signal received")

	s.grpcServer.GracefulStop()

	if s.postgresClient.Pool != nil {
		s.postgresClient.Close()
	}

	if s.redisClient != nil {
		s.redisClient.CloseClient()
	}

	s.logger.Info("gRPC server stopped")

	return nil
}
