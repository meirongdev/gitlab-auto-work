package config

type WorkConfig struct {
	Token        string
	Image        string
	BaseUrl      string
	BranchPrefix string
	UserEmail    string
	UserName     string
	CommitMsg    string
	Repositories []Repository
}

type Repository struct {
	Name     string
	Url      string
	Versions []string
}
