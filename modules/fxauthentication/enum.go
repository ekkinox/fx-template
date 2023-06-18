package fxauthentication

type EntityKind int

const (
	UnknownEntityKind EntityKind = iota
	GuestEntityKind
	UserEntityKind
	AdminEntityKind
	MachineEntityKind
)

func (e EntityKind) String() string {
	switch e {
	case GuestEntityKind:
		return "guest"
	case UserEntityKind:
		return "user"
	case AdminEntityKind:
		return "admin"
	case MachineEntityKind:
		return "machine"
	default:
		return "unknown"
	}
}

type AccountKind int

const (
	UnknownAccountKind AccountKind = iota
	BrandAccountKind
	RetailerAccountKind
)

func (a AccountKind) String() string {
	switch a {
	case BrandAccountKind:
		return "brand"
	case RetailerAccountKind:
		return "retailer"
	default:
		return "unknown"
	}
}

type IdentityProviderKind int

const (
	UnknownIdentityProviderKind IdentityProviderKind = iota
	GuestIdentityProviderKind
	UserIdentityProviderKind
	AdminIdentityProviderKind
	MachineIdentityProviderKind
	ImpersonationIdentityProviderKind
)

func (i IdentityProviderKind) String() string {
	switch i {
	case GuestIdentityProviderKind:
		return "aks_guest"
	case UserIdentityProviderKind:
		return "aks_user"
	case AdminIdentityProviderKind:
		return "aks_admin"
	case MachineIdentityProviderKind:
		return "aks_machine"
	case ImpersonationIdentityProviderKind:
		return "aks_imp"
	default:
		return "unknown"
	}
}
