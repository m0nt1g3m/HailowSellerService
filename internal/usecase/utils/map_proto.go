package utils

import (
	pb "HailowSellerService/HailowProto/build/go/SellerService/v1"
	"HailowSellerService/internal/domain"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func MapSeller(seller *domain.Seller) *pb.Seller {
	if seller == nil {
		return nil
	}

	return &pb.Seller{
		Id:               seller.ID,
		LogoUrl:          seller.LogoURL,
		StoreName:        seller.StoreName,
		StoreDescription: seller.StoreDescription,
		Tin:              seller.TIN,
		Psrn:             seller.PSRN,
		Kpp:              seller.KPP,
		OrganizationForm: MapOrganizationForm(seller.OrganizationForm),
		Email:            seller.Email,
		City:             seller.City,
		Street:           seller.Street,
		Building:         seller.Building,
		IsActive:         seller.IsActive,
		CreatedAt:        timestamppb.New(seller.CreatedAt),
		UpdatedAt:        timestamppb.New(seller.UpdatedAt),
	}
}

func MapOrganizationForm(organization_form domain.OrganizationForm) pb.OrganizationForm {
	switch organization_form {
	case domain.OrganizationLLC:
		return pb.OrganizationForm_LLC
	case domain.OrganizationJSC:
		return pb.OrganizationForm_JSC
	case domain.OrganizationOOO:
		return pb.OrganizationForm_OOO
	case domain.OrganizationIP:
		return pb.OrganizationForm_IP
	default:
		return pb.OrganizationForm_UNSPECIFIED
	}
}

func MapTokenPair(tokens *domain.TokenPair) *pb.TokenPair {
	if tokens == nil {
		return nil
	}

	return &pb.TokenPair{AccessToken: tokens.AccessToken, RefreshToken: tokens.RefreshToken}
}
