package auth

type AuthType = string

const (
	V2Password              AuthType = "v2password"
	V2Token                 AuthType = "v2token"
	V3Password              AuthType = "v3password"
	V3TOTP                  AuthType = "v3totp"
	V3ApplicationCredential AuthType = "v3applicationcredential"
	V3Token                 AuthType = "v3token"
	V3MultiFactor           AuthType = "v3multifactor"
)

type AuthData map[string]map[string]any

type AuthResult struct {
	// TODO (danchild)
}

type AuthV2 interface {
	ToAuthBody() ([]AuthType, AuthData, error)
	CanReauth() bool
}

type AuthV3 interface {
	ToAuthBody() ([]AuthType, AuthData, error)
	ToAuthHeaders() (map[string]any, error)
	CanReauth() bool
}

type AuthOptionsBuilder interface {
	Authenticate() error
}

type AuthOptionsV2 struct {
	AuthURL string
	Auth    AuthV2
}

func (ao AuthOptionsV2) Authenticate() error {
	// TODO (danchild)
	return nil
}

type AuthOptionsV3 struct {
	AuthURL     string
	AuthType    AuthType
	AuthMethods []string
	Auth        AuthV3
}

func (ao AuthOptionsV3) Authenticate() error {
	// TODO (danchild)
	return nil
}
