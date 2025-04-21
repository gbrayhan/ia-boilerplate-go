package middlewares

import (
	"errors"
	"github.com/gin-gonic/gin"
	"ia-boilerplate/src/db"
	"net/http"
)

type MessagesResponse struct {
	Message string `json:"message"`
}

func Handler(c *gin.Context) {
	c.Next()
	errs := c.Errors

	if len(errs) > 0 {
		var err *db.AppError
		ok := errors.As(errs[0].Err, &err)
		if ok {
			resp := MessagesResponse{Message: err.Error()}
			switch err.Type {
			case db.NotFound:
				c.JSON(http.StatusNotFound, resp)
				return
			case db.ValidationError:
				c.JSON(http.StatusBadRequest, resp)
				return
			case db.ResourceAlreadyExists:
				c.JSON(http.StatusConflict, resp)
				return
			case db.NotAuthenticated:
				c.JSON(http.StatusUnauthorized, resp)
				return
			case db.NotAuthorized:
				c.JSON(http.StatusForbidden, resp)
				return
			case db.RepositoryError:
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
