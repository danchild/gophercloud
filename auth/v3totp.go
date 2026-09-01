package auth

import "github.com/gophercloud/gophercloud/v2"

type V3TOTPOpts struct {
	Username       string
	UserID         string
	Passcode       string
	UserDomainID   string
	UserDomainName string
	Scope          *Scope
}

func (opts V3TOTPOpts) ToAuthBody() ([]AuthType, AuthData, error) {
	type domainReq struct {
		ID   *string `json:"id,omitempty"`
		Name *string `json:"name,omitempty"`
	}

	type userReq struct {
		ID       *string    `json:"id,omitempty"`
		Name     *string    `json:"name,omitempty"`
		Passcode string     `json:"passcode"`
		Domain   *domainReq `json:"domain,omitempty"`
	}

	type totpReq struct {
		User userReq `json:"user"`
	}

	authTypes := []AuthType{V3TOTP}
	req := totpReq{
		User: userReq{
			Passcode: opts.Passcode,
		},
	}

	if opts.Passcode == "" {
		// A passcode must be specified.
		return authTypes, nil, gophercloud.ErrMissingPasscode{}
	}

	// Exactly one of Username and UserID must be specified
	if opts.Username == "" && opts.UserID == "" {
		return authTypes, nil, gophercloud.ErrUsernameOrUserID{}
	} else if opts.Username != "" && opts.UserID != "" {
		return authTypes, nil, gophercloud.ErrUsernameOrUserID{}
	}

	if opts.Username != "" {
		// Exactly one of DomainID or DomainName must be specified
		if opts.UserDomainID == "" && opts.UserDomainName == "" {
			return authTypes, nil, gophercloud.ErrDomainIDOrDomainName{}
		} else if opts.UserDomainID != "" && opts.UserDomainName != "" {
			return authTypes, nil, gophercloud.ErrDomainIDOrDomainName{}
		}

		var domain *domainReq

		if opts.UserDomainID != "" {
			domain = &domainReq{ID: &opts.UserDomainID}
		} else { // opts.DomainName != ""
			domain = &domainReq{Name: &opts.UserDomainName}
		}

		req.User.Name = &opts.Username
		req.User.Domain = domain
	} else { // opts.UserID != ""
		// None of DomainID or DomainName may be specified
		if opts.UserDomainID != "" {
			return authTypes, nil, gophercloud.ErrDomainIDWithUserID{}
		} else if opts.UserDomainName != "" {
			return authTypes, nil, gophercloud.ErrDomainNameWithUserID{}
		}
		req.User.ID = &opts.UserID
	}

	b, err := gophercloud.BuildRequestBody(req, "")
	if err != nil {
		return authTypes, nil, err
	}

	result := AuthData{
		"totp": b,
	}

	return authTypes, result, nil
}

func (opts V3TOTPOpts) ToAuthHeaders() (map[string]any, error) {
	return nil, nil
}

func (opts V3TOTPOpts) CanReauth() bool {
	return false
}
