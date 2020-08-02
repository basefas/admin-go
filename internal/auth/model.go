package auth

type CreatePolicy struct {
	UID          string `json:"uid" binding:"required"`
	PolicyUrl    string `json:"policy_url" binding:"required"`
	PolicyMethod string `json:"policy_method" binding:"required"`
}
