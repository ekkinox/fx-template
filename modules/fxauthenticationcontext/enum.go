package fxauthenticationcontext

import (
	"encoding/json"
	"fmt"
)

type EntityType int

const (
	UnknownEntity EntityType = iota
	GuestEntity
	UserEntity
	AdminEntity
	MachineEntity
)

func (e EntityType) String() string {
	switch e {
	case GuestEntity:
		return "guest"
	case UserEntity:
		return "user"
	case AdminEntity:
		return "admin"
	case MachineEntity:
		return "machine"
	default:
		return "unknown"
	}
}

func (e *EntityType) UnmarshalJSON(data []byte) error {
	var str string

	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}

	switch str {
	case "guest":
		*e = GuestEntity
	case "user":
		*e = UserEntity
	case "admin":
		*e = AdminEntity
	case "machine":
		*e = MachineEntity
	default:
		return fmt.Errorf("invalid entity type: %s", str)
	}

	return nil
}

type AccountType int

const (
	UnknownAccount AccountType = iota
	BrandAccount
	RetailerAccount
)

func (a AccountType) String() string {
	switch a {
	case BrandAccount:
		return "brand"
	case RetailerAccount:
		return "retailer"
	default:
		return "unknown"
	}
}

func (a *AccountType) UnmarshalJSON(data []byte) error {
	var str string

	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}

	switch str {
	case "brand":
		*a = BrandAccount
	case "retailer":
		*a = RetailerAccount
	default:
		return fmt.Errorf("invalid account type: %s", str)
	}

	return nil
}

type IdentityProviderType int

const (
	UnknownIdentityProvider IdentityProviderType = iota
	GuestIdentityProvider
	UserIdentityProvider
	AdminIdentityProvider
	MachineIdentityProvider
	ImpersonationIdentityProvider
)

func (i IdentityProviderType) String() string {
	switch i {
	case GuestIdentityProvider:
		return "aks_guest"
	case UserIdentityProvider:
		return "aks_user"
	case AdminIdentityProvider:
		return "aks_admin"
	case MachineIdentityProvider:
		return "aks_machine"
	case ImpersonationIdentityProvider:
		return "aks_imp"
	default:
		return "unknown"
	}
}

func (i *IdentityProviderType) UnmarshalJSON(data []byte) error {
	var str string

	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}

	switch str {
	case "aks_guest":
		*i = GuestIdentityProvider
	case "aks_user":
		*i = UserIdentityProvider
	case "aks_admin":
		*i = AdminIdentityProvider
	case "aks_machine":
		*i = MachineIdentityProvider
	case "aks_imp":
		*i = ImpersonationIdentityProvider
	default:
		return fmt.Errorf("invalid identity provider type: %s", str)
	}

	return nil
}
