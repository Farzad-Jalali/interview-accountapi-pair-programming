package convert

import (
	"github.com/form3tech/go-form3-web/web"
	"github.com/form3tech-oss/interview-accountapi-pair-programming/internal/swagger-client/interview-accountapi-pair-programming/models"
)

func FromLinks(links web.Links) *models.Links {
	return &models.Links{
		Self:  links.Self,
		First: links.First,
		Next:  links.Next,
		Prev:  links.Prev,
		Last:  links.Last,
	}
}
