package auth

import "github.com/gophercloud/gophercloud/v2"

type V3TokenOpts struct {
	Token string
	Scope *Scope
}

func (opts V3TokenOpts) ToAuthBody() ([]AuthType, AuthData, error) {
	type tokenReq struct {
		Token string `json:"token"`
	}

	authTypes := []AuthType{V3Token}
	req := tokenReq{
		Token: opts.Token,
	}

	b, err := gophercloud.BuildRequestBody(req, "")
	if err != nil {
		return authTypes, nil, err
	}

	result := AuthData{
		"token": b,
	}

	return authTypes, result, nil
}

func (opts V3TokenOpts) ToAuthHeaders() (map[string]any, error) {
	return nil, nil
}

func (opts V3TokenOpts) CanReauth() bool {
	return false
}
