package fxauthentication

type Entity interface {
	IdentityProvider() IdentityProviderKind
	Entity() EntityKind
	Uuid() string
}

type UserAccount struct {
	kind AccountKind
	uuid string
}

func (a *UserAccount) Kind() AccountKind {
	return a.kind
}

func (a *UserAccount) Uuid() string {
	return a.uuid
}

func NewUserAccount(kind AccountKind, uuid string) *UserAccount {
	return &UserAccount{
		kind: kind,
		uuid: uuid,
	}
}

type User struct {
	identityProvider IdentityProviderKind
	entity           EntityKind
	account          *UserAccount
	uuid             string
}

func NewUser(identityProvider IdentityProviderKind, entity EntityKind, account *UserAccount, uuid string) *User {
	return &User{
		identityProvider: identityProvider,
		entity:           entity,
		account:          account,
		uuid:             uuid,
	}
}
