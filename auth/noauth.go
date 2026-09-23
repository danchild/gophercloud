package auth

type NoAuthOpts struct {
	Username    string `json:"username,omitempty"`
	ProjectName string `json:"project_name,omitempty"`
}

func (opts NoAuthOpts) ToAuthBody() (map[string]map[string]any, error) {
	return map[string]map[string]any{}, nil
}

func (opts NoAuthOpts) CanReauth() bool {
	return false
}
