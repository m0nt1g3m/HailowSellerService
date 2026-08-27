package domain

import (
	"fmt"
)

type Error struct {
	Code    string
	Message string
	Cause   error
}

func (e *Error) Error() string {
	if e == nil {
		return ""
	}
	if e.Cause == nil {
		return e.Message
	}
	return fmt.Sprintf("%s: %v", e.Message, e.Cause)
}

func (e *Error) Unwrap() error {
	return e.Cause
}

func (e *Error) Is(target error) bool {
	if t, ok := target.(*Error); ok {
		return e.Code == t.Code
	}
	return false
}

var (
	ErrGRPCServer            = &Error{Code: "grpc_server_error", Message: "gRPC server error"}
	ErrListenAddress         = &Error{Code: "failed_listen_address", Message: "Failed listen address"}
	ErrInitServer            = &Error{Code: "failed_init_server", Message: "Failed to initialize server"}
	ErrRunServer             = &Error{Code: "failed_run_server", Message: "Failed to run server"}
	ErrRedisClose            = &Error{Code: "failed_closing_redis", Message: "Failed closing Redis"}
	ErrSellerAlreadyExists   = &Error{Code: "seller_already_exists", Message: "Seller already exists"}
	ErrSellerNotFound        = &Error{Code: "seller_not_found", Message: "Seller not found"}
	ErrInvalidSellerID       = &Error{Code: "invalid_seller_id", Message: "Invalid seller id"}
	ErrInvalidCredentials    = &Error{Code: "invalid_credentials", Message: "Invalid credentials"}
	ErrRefreshTokenNotFound  = &Error{Code: "refresh_token_not_found", Message: "refresh_token not found"}
	ErrSessionNotFound       = &Error{Code: "session_not_found", Message: "Session not found"}
	ErrUnauthorized          = &Error{Code: "unauthorized", Message: "Unauthorized"}
	ErrLogoEmpty             = &Error{Code: "logo_image_empty", Message: "Logo image is empty"}
	ErrLogoUploadFailed      = &Error{Code: "logo_upload_failed", Message: "Failed to upload logo"}
	ErrLogoUploadUnavailable = &Error{Code: "logo_upload_unavailable", Message: "Logo upload is unavailable"}
	ErrInvalidRequest        = &Error{Code: "invalid_auth_request", Message: "Invalid auth request"}
)
