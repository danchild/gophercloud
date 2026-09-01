package auth

import (
	"fmt"
	"maps"

	"github.com/gophercloud/gophercloud/v2"
)

type V3MultifactorOpts struct {
	AuthMethods []AuthV3
	Scope       *Scope
}

func (opts V3MultifactorOpts) ToAuthBody() ([]AuthType, AuthData, error) {
	var result AuthData
	var authTypes []AuthType
	for _, authMethod := range opts.AuthMethods {
		var authResult AuthData
		var err error

		switch authMethod.(type) {
		case V3PasswordOpts:
			authTypes = append(authTypes, V3Password)
		case V3TOTPOpts:
			authTypes = append(authTypes, V3TOTP)
		case V3ApplicationCredentialOpts:
			authTypes = append(authTypes, V3ApplicationCredential)
		case V3TokenOpts:
			authTypes = append(authTypes, V3Token)
		default:
			return authTypes, nil, gophercloud.ErrUnsupportedAuthType{AuthType: fmt.Sprintf("%T", authMethod)}
		}

		_, authResult, err = authMethod.ToAuthBody()
		if err != nil {
			return authTypes, nil, err
		}

		maps.Copy(result, authResult)
	}

	return authTypes, result, nil
}

func (opts V3MultifactorOpts) ToAuthHeaders() (map[string]any, error) {
	return nil, nil
}

func (opts V3MultifactorOpts) CanReauth() bool {
	return false
}
