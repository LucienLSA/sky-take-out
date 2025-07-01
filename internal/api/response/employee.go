package response

type EmployeeLogin struct {
	Id           uint64 `json:"id"`
	Name         string `json:"name"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	UserName     string `json:"userName"`
}
