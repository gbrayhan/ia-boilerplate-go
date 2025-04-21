package middlewares

import (
	"errors"
	"github.com/gin-gonic/gin"
	"ia-boilerplate/src/repository"
	"net/http"
)

type MessagesResponse struct {
	Message string `json:"message"`
}

func Handler(c *gin.Context) {
	c.Next()
	errs := c.Errors

	if len(errs) > 0 {
		var err *repository.AppError
		ok := errors.As(errs[0].Err, &err)
		if ok {
			resp := MessagesResponse{Message: err.Error()}
			switch err.Type {
			case repository.NotFound:
				c.JSON(http.StatusNotFound, resp)
				return
			case repository.ValidationError:
				c.JSON(http.StatusBadRequest, resp)
				return
			case repository.ResourceAlreadyExists:
				c.JSON(http.StatusConflict, resp)
				return
			case repository.NotAuthenticated:
				c.JSON(http.StatusUnauthorized, resp)
				return
			case repository.NotAuthorized:
				c.JSON(http.StatusForbidden, resp)
				return
			case repository.RepositoryError:
				c.JSON(http.StatusInternalServerError, MessagesResponse{Message: "We are working to improve the flow of this request."})
				return
			default:
				c.JSON(http.StatusInternalServerError, MessagesResponse{Message: "We are working to improve the flow of this request."})
				return
			}
		}
		return
	}
}
