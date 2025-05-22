package types

type GithubUser struct {
	ID        int    `json:"id"`
	Username  string `json:"login"`
	Email     string `json:"email"`
	AvatarUrl string `json:"avatar_url"`
}
