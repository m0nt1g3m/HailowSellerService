package auth

import (
	"HailowSellerService/internal/domain"
	redis_repository "HailowSellerService/internal/infrastructure/redis/repository"
	"HailowSellerService/internal/repository"
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthUseCase struct {
	repo       *redis_repository.SessionRepository
	sellerRepo *repository.SellerRepository
	accessTTL  time.Duration
	refreshTTL time.Duration
}

type IAuthUseCase interface {
	SignUp(ctx context.Context, input *domain.SellerInfo) (*domain.Seller, error)
	SignIn(ctx context.Context, input *domain.SellerInfo) (*domain.AuthInfo, error)
}

type jwtClaims struct {
	ID               uuid.UUID `json:"id"`
	Email            string    `json:"email"`
	OrganizationForm string    `json:"organization_form"`
	StoreName        string    `json:"store_name"`
	TokenType        string    `json:"token_type,omitempty"`
	jwt.RegisteredClaims
}

func NewAuthUseCase(sellerRepo *repository.SellerRepository, repo *redis_repository.SessionRepository) *AuthUseCase {
	return &AuthUseCase{
		repo:       repo,
		sellerRepo: sellerRepo,
	}
}

func (u *AuthUseCase) SignUp(ctx context.Context, input *domain.SellerInfo) (*domain.Seller, error) {
	if input.Email == "" || input.Password == "" || input.PSRN == "" || input.KPP == "" || input.TIN == "" {
		return nil, domain.ErrInvalidCredentials
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	input.Password = string(hashedPassword)

	seller, err := u.sellerRepo.CreateSeller(ctx, input)
	if err != nil {
		return nil, err
	}

	return seller, nil
}

func (u *AuthUseCase) SignIn(ctx context.Context, input *domain.SellerInfo) (*domain.AuthInfo, error) {
	if input.Email == "" || input.Password == "" {
		return nil, domain.ErrInvalidCredentials
	}

	seller, err := u.sellerRepo.GetSellerByEmail(ctx, input.Email)

	if err != nil {
		return nil, domain.ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(seller.PasswordHash), []byte(input.Password)); err != nil {
		return nil, domain.ErrInvalidCredentials
	}

	tokens, err := u.generateTokenPair(ctx, seller)

	if err != nil {
		return nil, err
	}

	return &domain.AuthInfo{
		SellerID:  seller.ID,
		TokenPair: tokens,
	}, nil
}
