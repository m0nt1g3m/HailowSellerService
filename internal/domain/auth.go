package domain

import (
	"database/sql/driver"
	"fmt"
	"time"

	"github.com/google/uuid"
)

const (
	OrganizationUnspecified OrganizationForm = "UNSPECIFIED"
	OrganizationLLC         OrganizationForm = "LLC" // Limited Liability Company (ООО)
	OrganizationJSC         OrganizationForm = "JSC" // Joint Stock Company (АО)
	OrganizationOOO         OrganizationForm = "OOO" // Open Joint-Stock Company (ПАО)
	OrganizationIP          OrganizationForm = "IP"  // Individual Entrepreneur (ИП)
)

func (o *OrganizationForm) Scan(value any) error {
	if value == nil {
		*o = OrganizationUnspecified
		return nil
	}

	switch v := value.(type) {
	case string:
		*o = OrganizationForm(v)
	case []byte:
		*o = OrganizationForm(v)
	default:
		return fmt.Errorf("Cannot scan %T into domain.OrganizationForm", value)
	}
	return nil
}

func (o OrganizationForm) Value() (driver.Value, error) {
	return string(o), nil
}

type TokenPair struct {
	AccessToken  string
	RefreshToken string
}

type RefreshSession struct {
	ID           string
	SellerID     uuid.UUID
	DeviceID     string
	RefreshToken string
	UserAgent    string
	IP           string
	CreatedAt    time.Time
	ExpiresAt    time.Time
}

type AuthInfo struct {
	SellerID  string
	TokenPair *TokenPair
}

func (s *RefreshSession) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}
