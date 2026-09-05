package auth

import (
	"HailowSellerService/internal/domain"
	"context"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func (u *AuthUseCase) generateTokenPair(ctx context.Context, seller *domain.Seller) (*domain.TokenPair, error) {
	now := time.Now()

	accessTokenStr, err := u.generateAccessToken(now, seller)

	if err != nil {
		return nil, err
	}

	refreshTokenStr, err := u.generateRefreshToken(now, seller)

	if err != nil {
		return nil, err
	}

	session := &domain.RefreshSession{
		SellerID:     uuid.MustParse(seller.ID),
		RefreshToken: refreshTokenStr,
		CreatedAt:    now,
		ExpiresAt:    now.Add(7 * 24 * time.Hour),
	}

	if err := u.repo.CreateSession(ctx, session); err != nil {
		return nil, err
	}

	return &domain.TokenPair{
		AccessToken:  accessTokenStr,
		RefreshToken: refreshTokenStr,
	}, nil
}

func (u *AuthUseCase) generateAccessToken(now time.Time, seller *domain.Seller) (string, error) {
	claims := jwtClaims{
		ID:               uuid.MustParse(seller.ID),
		Email:            seller.Email,
		OrganizationForm: string(seller.OrganizationForm),
		TokenType:        "access",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}
	accessTokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessTokenStr, err := accessTokenObj.SignedString(u.getAccessSigningKey())
	return accessTokenStr, err
}

func (u *AuthUseCase) generateRefreshToken(now time.Time, seller *domain.Seller) (string, error) {
	claims := jwtClaims{
		ID:               uuid.MustParse(seller.ID),
		Email:            seller.Email,
		OrganizationForm: string(seller.OrganizationForm),
		TokenType:        "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(7 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}
	refreshTokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	refreshTokenStr, err := refreshTokenObj.SignedString(u.getRefreshSigningKey())
	return refreshTokenStr, err
}

func (u *AuthUseCase) getAccessSigningKey() []byte {
	if key := os.Getenv("JWT_ACCESS_KEY_SELLER_SERVICE"); key != "" {
		return []byte(key)
	}

	return []byte("secret")
}

func (u *AuthUseCase) getRefreshSigningKey() []byte {
	if key := os.Getenv("JWT_REFRESH_KEY_SELLER_SERVICE"); key != "" {
		return []byte(key)
	}

	return []byte("secret")
}

func (u *AuthUseCase) ParseAccessToken(tokenStr string) (*jwtClaims, error) {
	var claims jwtClaims
	token, err := jwt.ParseWithClaims(tokenStr, &claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, domain.ErrUnauthorized
		}
		return u.getAccessSigningKey(), nil
	})

	if err != nil || !token.Valid {
		return nil, domain.ErrUnauthorized
	}

	if claims.TokenType != "" && claims.TokenType != "access" {
		return nil, domain.ErrUnauthorized
	}

	return &claims, nil
}
