package views

import "github.com/PotatoesFall/ots/domain"

type ClaimSecretResponse struct {
	ID      int    `json:"id"`
	Message string `json:"message"`
	Content string `json:"content"`
}

func NewClaimSecretResponse(s domain.Secret) ClaimSecretResponse {
	return ClaimSecretResponse{
		ID:      s.ID,
		Message: s.Message,
		Content: s.Content,
	}
}
