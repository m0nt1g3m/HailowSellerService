package profile

import (
	"HailowSellerService/internal/domain"
	"HailowSellerService/internal/repository"
	"context"
	"fmt"

	"github.com/google/uuid"
)

type ProfileUseCase struct {
	sellerRepo *repository.SellerRepository
}

type IProfileUseCase interface {
	GetProfile(ctx context.Context, id string) (*domain.Seller, error)
}

func NewProfileUseCase(sellerRepo *repository.SellerRepository) *ProfileUseCase {
	return &ProfileUseCase{
		sellerRepo: sellerRepo,
	}
}

func (u *ProfileUseCase) GetProfile(ctx context.Context, id string) (*domain.Seller, error) {
	if id == "" {
		return nil, domain.ErrInvalidCredentials
	}

	sellerID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", domain.ErrInvalidSellerID.Message, err)
	}

	seller, err := u.sellerRepo.GetSellerByID(ctx, sellerID)
	if err != nil {
		return nil, err
	}

	return seller, nil
}

// func (u *ProfileUseCase) (ctx context.Context, id string) (*domain.Seller, error) {
// 	claims, err := auth.ParseAccessToken(id)

// 	if err != nil {
// 		return nil, domain.ErrUnauthorized
// 	}

// }
