package auth

import "github.com/gophercloud/gophercloud/v2"

type V3PasswordOpts struct {
	Username       string
	UserID         string
	Password       string
	UserDomainID   string
	UserDomainName string
	Scope          *Scope
}

func (opts V3PasswordOpts) ToAuthBody() ([]AuthType, AuthData, error) {
	type domainReq struct {
		ID   *string `json:"id,omitempty"`
		Name *string `json:"name,omitempty"`
	}

	type userReq struct {
		ID       *string    `json:"id,omitempty"`
		Name     *string    `json:"name,omitempty"`
		Password string     `json:"password"`
		Domain   *domainReq `json:"domain,omitempty"`
	}

	type passwordReq struct {
		User userReq `json:"user"`
	}

	authTypes := []AuthType{V3Password}
	req := passwordReq{
		User: userReq{
			Password: opts.Password,
		},
	}

	if opts.Password == "" {
		// A password must be specified.
		return authTypes, nil, gophercloud.ErrMissingPassword{}
	}

	// Exactly one of Username and UserID must be specified
	if opts.Username == "" && opts.UserID == "" {
		return authTypes, nil, gophercloud.ErrUsernameOrUserID{}
	} else if opts.Username != "" && opts.UserID != "" {
		return authTypes, nil, gophercloud.ErrUsernameOrUserID{}
	}

	if opts.Username != "" {
		// Exactly one of UserDomainID or UserDomainName must be specified
		if opts.UserDomainID == "" && opts.UserDomainName == "" {
			return authTypes, nil, gophercloud.ErrDomainIDOrDomainName{}
		} else if opts.UserDomainID != "" && opts.UserDomainName != "" {
			return authTypes, nil, gophercloud.ErrDomainIDOrDomainName{}
		}

		var domain *domainReq

		if opts.UserDomainID != "" {
			domain = &domainReq{ID: &opts.UserDomainID}
		} else { // opts.UserDomainName != ""
			domain = &domainReq{Name: &opts.UserDomainName}
		}

		req.User.Name = &opts.Username
		req.User.Domain = domain
	} else { // opts.UserID != ""
		// None of UserDomainID or UserDomainName may be specified
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
		"password": b,
	}

	return authTypes, result, nil
}

func (opts V3PasswordOpts) ToAuthHeaders() (map[string]any, error) {
	return nil, nil
}

func (opts V3PasswordOpts) CanReauth() bool {
	return true
}
