package repository

import (
	"HailowSellerService/internal/domain"
	"HailowSellerService/pkg/database"
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type SellerRepository struct {
	postgresClient *database.Postgres
}

func NewSellerRepository(postgresClient *database.Postgres) *SellerRepository {
	return &SellerRepository{
		postgresClient: postgresClient,
	}
}

func (r *SellerRepository) CreateSeller(ctx context.Context, input *domain.SellerInfo) (*domain.Seller, error) {
	query := `
		INSERT INTO sellers_schema.sellers (store_name, store_description, tin, kpp, psrn, organization_form, email, city, street, building, password_hash)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING *`
	rows, err := r.postgresClient.Pool.Query(ctx, query, input.StoreName, input.StoreDescription, input.TIN, input.KPP, input.PSRN, input.OrganizationForm, input.Email, input.City, input.Street, input.Building, input.Password)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	seller, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[domain.Seller])

	if err != nil {
		return nil, err
	}

	return &seller, nil
}

func (r *SellerRepository) GetSellerByID(ctx context.Context, id uuid.UUID) (*domain.Seller, error) {
	query := `
		SELECT *
		FROM sellers_schema.sellers
		WHERE id = $1
		LIMIT 1
	`

	rows, err := r.postgresClient.Pool.Query(ctx, query, id)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	seller, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[domain.Seller])

	if err != nil {
		return nil, err
	}

	return &seller, nil
}

func (r *SellerRepository) GetSellerByEmail(ctx context.Context, email string) (*domain.Seller, error) {
	query := `
		SELECT *
		FROM sellers_schema.sellers
		WHERE email = $1
		LIMIT 1
	`

	rows, err := r.postgresClient.Pool.Query(ctx, query, email)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	seller, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[domain.Seller])

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrSellerNotFound
		}
		return nil, err
	}

	return &seller, nil
}
