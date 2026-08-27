package domain

import "time"

type OrganizationForm string

type SellerInfo struct {
	ID               string
	StoreName        string
	StoreDescription string
	LogoURL          *string
	TIN              string
	KPP              string
	PSRN             string
	OrganizationForm OrganizationForm
	Email            string
	City             string
	Street           string
	Building         string
	Password         string
}

type Seller struct {
	ID               string
	StoreName        string
	StoreDescription string
	LogoURL          *string
	TIN              string
	PSRN             string
	KPP              string
	OrganizationForm OrganizationForm
	Email            string
	City             string
	Street           string
	Building         string
	IsActive         bool
	PasswordHash     string
	CreatedAt        time.Time
	UpdatedAt        time.Time
}
