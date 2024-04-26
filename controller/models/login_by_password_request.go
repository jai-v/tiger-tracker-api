package models

type LoginByPasswordRequest struct {
	Username       string `json:"userName"`
	Password       string `json:"password"`
	RememberMe     bool   `json:"rememberMe"`
	LoginChallenge string `json:"loginChallenge"`
}
