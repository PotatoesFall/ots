package router

import "github.com/PotatoesFall/ots/domain"

type NewSecretRequest struct {
	Message string `form:"message"`
	Content string `form:"content"`
}

func (nsr NewSecretRequest) Secret() domain.NewSecret {
	return domain.NewSecret{
		Message: nsr.Message,
		Content: nsr.Content,
	}
}
