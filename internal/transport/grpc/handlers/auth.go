package handlers

import (
	pb "HailowSellerService/HailowProto/build/go/SellerService/v1"
	"HailowSellerService/internal/domain"
	"HailowSellerService/internal/transport/grpc/response/errorcode"
	"HailowSellerService/internal/usecase/auth"
	"HailowSellerService/internal/usecase/utils"
	"HailowSellerService/pkg/logging"
	"context"
)

type AuthHandler struct {
	usecase auth.IAuthUseCase
	logger  *logging.Logger
	pb.UnimplementedSellerServiceServer
}

func NewAuthHandler(usecase auth.IAuthUseCase, logger *logging.Logger) *AuthHandler {
	return &AuthHandler{usecase: usecase, logger: logger}
}

func (h *AuthHandler) SignUp(ctx context.Context, req *pb.SignUpRequest) (*pb.SignUpResponse, error) {
	input := domain.SellerInfo{
		StoreName:        req.GetStoreName(),
		StoreDescription: req.GetStoreDescription(),
		TIN:              req.GetTin(),
		PSRN:             req.GetPsrn(),
		KPP:              req.GetKpp(),
		OrganizationForm: domain.OrganizationLLC,
		Email:            req.GetEmail(),
		City:             req.GetCity(),
		Street:           req.GetStreet(),
		Building:         req.GetBuilding(),
		Password:         req.GetPassword(),
	}

	seller, err := h.usecase.SignUp(ctx, &input)
	if err != nil {
		h.logger.Errorf("SignUp error: %v", err)
		return nil, errorcode.ToStatus(err)
	}

	return &pb.SignUpResponse{Seller: utils.MapSeller(seller)}, nil
}

func (h *AuthHandler) SignIn(ctx context.Context, req *pb.SignInRequest) (*pb.SignInResponse, error) {
	input := domain.SellerInfo{
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
	}

	authInfo, err := h.usecase.SignIn(ctx, &input)

	if err != nil {
		h.logger.Errorf("SignIn error: %v", err)
		return nil, errorcode.ToStatus(err)
	}

	return &pb.SignInResponse{Id: authInfo.SellerID, Tokens: utils.MapTokenPair(authInfo.TokenPair)}, nil
}
