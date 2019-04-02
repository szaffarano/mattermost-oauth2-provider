package oauth

import (
	"log"
	"net/http"

	"github.com/pkg/errors"

	"github.com/gin-gonic/gin"
	"github.com/go-session/session"
	"github.com/szaffarano/oauth-server/domain"
	"gopkg.in/oauth2.v3"
	oauterrors "gopkg.in/oauth2.v3/errors"
	"gopkg.in/oauth2.v3/server"
)

const (
	returnURI      = "ReturnUri"
	loggedInUserID = "LoggedInUserID"
	tokenFormKey   = "token"
	signFormKey    = "sign"
)

// LoginHandler returns a handler that manages authentication against a signed token
func LoginHandler(srv *server.Server, cfg domain.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		var token string
		var exists bool

		store := getStore(c)

		if token, exists = c.GetPostForm(tokenFormKey); !exists {
			c.AbortWithError(
				http.StatusBadRequest, errors.New("Token not found"))
		}

		if _, exists = c.GetPostForm(signFormKey); !exists {
			c.AbortWithError(
				http.StatusBadRequest, errors.New("Sign not found"))
		}

		// TODO token&sign should be base64 coded
		// TODO validate token with the signature

		store.Set(loggedInUserID, token)
		store.Save()

		c.Header("Location", cfg.MattermostURL)
		c.Status(http.StatusFound)
	}
}

// AuthorizeHandler returns a handler
func AuthorizeHandler(srv *server.Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		err := srv.HandleAuthorizeRequest(c.Writer, c.Request)
		if err != nil {
			c.AbortWithError(
				http.StatusBadRequest,
				errors.Wrap(err, "Error authorizing request"))
		}
	}
}

// UserAuthorizationHandler returns a handler
func UserAuthorizationHandler(cfg domain.Config) func(http.ResponseWriter, *http.Request) (string, error) {
	return func(w http.ResponseWriter, r *http.Request) (string, error) {
		var userID string
		var store session.Store
		var err error

		if store, err = session.Start(nil, w, r); err == nil {
			if uid, exists := store.Get(loggedInUserID); exists {
				userID = uid.(string)
				store.Delete(loggedInUserID)
				store.Save()
			} else {
				w.Header().Set("Location", cfg.SUAURL)
				w.WriteHeader(http.StatusFound)
			}
		}
		return userID, err
	}
}

// RenderHTMLHandler render a static resource
// func RenderHTMLHandler(res string) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		c.HTML(http.StatusOK, res, nil)
// 	}
// }

// TokenHandler returns a handler
func TokenHandler(srv *server.Server) func(*gin.Context) {
	return func(c *gin.Context) {
		err := srv.HandleTokenRequest(c.Writer, c.Request)
		if err != nil {
			c.AbortWithError(
				http.StatusInternalServerError,
				errors.Wrap(err, "Error handling token request"))
		}
	}
}

// GetUserHandler returns a handler
func GetUserHandler(srv *server.Server) func(*gin.Context) {
	return func(c *gin.Context) {
		var ticket domain.Ticket
		var err error
		var bearer oauth2.TokenInfo

		bearer, err = srv.ValidationBearerToken(c.Request)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, errors.Wrap(err, ""))
		}

		if ticket, err = domain.ParseTicket(bearer.GetUserID()); err != nil {
			c.AbortWithError(
				http.StatusInternalServerError,
				errors.New("Error parsing bearer token"))
		}

		c.JSON(http.StatusOK, ticket.ToUser())
	}
}

// InternalErrorHandler handles internal errors
func InternalErrorHandler(err error) *oauterrors.Response {
	var re oauterrors.Response

	log.Println("Internal Error:", err.Error())

	re.Description = err.Error()
	re.Error = err
	re.ErrorCode = 1
	re.StatusCode = http.StatusInternalServerError

	return &re
}

// ResponseErrorHandler handles response errors
func ResponseErrorHandler(re *oauterrors.Response) {
	log.Println("Response Error:", re.Error.Error())
}
