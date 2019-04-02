package oauth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-session/session"
	"github.com/pkg/errors"
)

func getStore(c *gin.Context) session.Store {
	store, err := session.Start(nil, c.Writer, c.Request)
	if err != nil {
		c.AbortWithError(
			http.StatusInternalServerError,
			errors.Wrap(err, "Error getting session"))
	}

	return store
}
