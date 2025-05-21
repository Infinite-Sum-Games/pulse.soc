package types

type GithubUser struct {
	ID        string `json:"id"`
	Username  string `json:"login"`
	Email     string `json:"email"`
	AvatarUrl string `json:"avatar_url"`
}
