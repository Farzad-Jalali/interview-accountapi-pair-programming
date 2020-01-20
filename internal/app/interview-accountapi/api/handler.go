package api

import (
	"context"
	"github.com/form3tech/go-security/security"
	"github.com/form3tech-oss/interview-accountapi-pair-programming/internal/app/interview-accountapi/api/errors"
	models "github.com/form3tech-oss/interview-accountapi-pair-programming/internal/app/interview-accountapi/api/externalmodels"
	"github.com/form3tech-oss/interview-accountapi-pair-programming/internal/app/interview-accountapi/log"
	"github.com/gin-gonic/gin"
	pkgerr "github.com/pkg/errors"
	"net/http"
)

func WithUserContext(handler func(ctx *context.Context, c *gin.Context) error) func(ctx *gin.Context) {
	return func(c *gin.Context) {
		ctx := context.Background()
		if err := handler(&ctx, c); err != nil {
			switch e := pkgerr.Cause(err).(type) {
			case *security.AuthError:
				log.Infof("%v", e)
				c.JSON(http.StatusForbidden, models.APIError{ErrorMessage: "forbidden"})
				return
			case *errors.AccessDeniedError:
				log.Infof("%v", e)
				c.JSON(http.StatusForbidden, models.APIError{ErrorMessage: e.Error()})
				return
			case *errors.NotFoundError:
				log.Infof("%v", e)
				c.JSON(http.StatusNotFound, models.APIError{ErrorMessage: e.Error()})
				return
			case *errors.NotAcceptableError:
				log.Infof("%v", e)
				c.JSON(http.StatusNotAcceptable, models.APIError{ErrorMessage: e.Error()})
				return
			case *errors.DuplicateError:
				log.Infof("%v", e)
				c.JSON(http.StatusConflict, models.APIError{ErrorMessage: e.Error()})
				return
			case *errors.ConflictError:
				log.Infof("%v", e)
				c.JSON(http.StatusConflict, models.APIError{ErrorMessage: e.Error()})
				return
			case *errors.IllegalArgumentError:
				log.Infof("%v", e)
				c.JSON(http.StatusBadRequest, models.APIError{ErrorMessage: e.Error()})
				return
			default:
				log.Errorf("server error:, %v", err)
				c.JSON(http.StatusInternalServerError, models.APIError{ErrorMessage: "server error"})
			}
		}

	}
}
