package fxauthenticationcontext

type Account struct {
	Uuid        string      `json:"id"`
	AccountType AccountType `json:"type"`
}

type AuthenticationContext struct {
	Uuid                 string               `json:"sub"`
	ClientId             string               `json:"cid,omitempty"`
	IdentityProviderType IdentityProviderType `json:"idp"`
	Aks                  struct {
		EntityType           EntityType             `json:"entity"`
		Account              *Account               `json:"account,omitempty"`
		ImpersonationContext *AuthenticationContext `json:"imp,omitempty"`
	} `json:"aks"`
}

func (c *AuthenticationContext) EntityType() EntityType {
	return c.Aks.EntityType
}

func (c *AuthenticationContext) Account() *Account {
	return c.Aks.Account
}

func (c *AuthenticationContext) ImpersonationContext() *AuthenticationContext {
	return c.Aks.ImpersonationContext
}

func (c *AuthenticationContext) IsGuestEntity() bool {
	return c.Aks.EntityType == GuestEntity
}

func (c *AuthenticationContext) IsUserEntity() bool {
	return c.Aks.EntityType == UserEntity
}

func (c *AuthenticationContext) IsAdminEntity() bool {
	return c.Aks.EntityType == AdminEntity
}

func (c *AuthenticationContext) IsMachineEntity() bool {
	return c.Aks.EntityType == MachineEntity
}

func (c *AuthenticationContext) IsBrandAccount() bool {
	if c.Aks.Account == nil {
		return false
	}

	return c.Aks.Account.AccountType == BrandAccount
}

func (c *AuthenticationContext) IsRetailerAccount() bool {
	if c.Aks.Account == nil {
		return false
	}

	return c.Aks.Account.AccountType == RetailerAccount
}

func (c *AuthenticationContext) IsFromGuestIdentityProvider() bool {
	return c.IdentityProviderType == GuestIdentityProvider
}

func (c *AuthenticationContext) IsFromUserIdentityProvider() bool {
	return c.IdentityProviderType == UserIdentityProvider
}

func (c *AuthenticationContext) IsFromAdminIdentityProvider() bool {
	return c.IdentityProviderType == AdminIdentityProvider
}

func (c *AuthenticationContext) IsFromMachineIdentityProvider() bool {
	return c.IdentityProviderType == MachineIdentityProvider
}

func (c *AuthenticationContext) IsFromImpersonationIdentityProvider() bool {
	return c.IdentityProviderType == ImpersonationIdentityProvider
}

func (c *AuthenticationContext) IsImpersonation() bool {
	return c.IsFromImpersonationIdentityProvider() && c.Aks.ImpersonationContext != nil
}
