package handlers

import (
	pb "HailowSellerService/HailowProto/build/go/SellerService/v1"
	"HailowSellerService/internal/transport/grpc/response/errorcode"
	"HailowSellerService/internal/usecase/profile"
	"HailowSellerService/internal/usecase/utils"
	"HailowSellerService/pkg/logging"
	"context"
)

type ProfileHandler struct {
	usecase profile.IProfileUseCase
	logger  *logging.Logger
	pb.UnimplementedSellerProfileServiceServer
}

func NewProfileHandler(usecase *profile.ProfileUseCase, logger *logging.Logger) *ProfileHandler {
	return &ProfileHandler{usecase: usecase, logger: logger}
}

func (h *ProfileHandler) GetProfile(ctx context.Context, req *pb.GetProfileRequest) (*pb.GetProfileResponse, error) {
	seller, err := h.usecase.GetProfile(ctx, req.GetId())

	if err != nil {
		h.logger.Errorf("GetProfile error: %v", err)
		return nil, errorcode.ToStatus(err)
	}

	return &pb.GetProfileResponse{Seller: utils.MapSeller(seller)}, nil
}
