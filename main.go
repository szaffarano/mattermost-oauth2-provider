package main

import (
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/szaffarano/oauth-server/domain"
	"github.com/szaffarano/oauth-server/oauth"
	"gopkg.in/oauth2.v3/generates"
	"gopkg.in/oauth2.v3/manage"
	"gopkg.in/oauth2.v3/models"
	"gopkg.in/oauth2.v3/server"
	"gopkg.in/oauth2.v3/store"
)

var (
	config = domain.Config{}
)

func init() {
	config.Client.ID = "222222"
	config.Client.Secret = "1234567890123456789012345678901234567890"
	config.Client.Domain = "http://zatopek.zaffarano.com.ar"

	config.JWTKey = []byte("1234567890")
	config.AuthorizeCodeExpiration = time.Minute * 2
	config.AccessTokenExpiration = time.Hour * 2
	config.RefreshTokenExpiration = time.Hour * 24 * 3
	config.GenerateRefresh = true

	config.MattermostURL = "http://zatopek.zaffarano.com.ar/oauth/gitlab/login"
	config.SUAURL = "http://localhost:8080"
}

func main() {
	var router = setupRouter()
	var srv = setupOauthServer()

	router.POST("/login", oauth.LoginHandler(srv, config))

	o2 := router.Group("/oauth")
	{
		o2.GET("/authorize", oauth.AuthorizeHandler(srv))
		o2.POST("/token", oauth.TokenHandler(srv))
	}

	api := router.Group("/api/v4")
	{
		api.GET("/user", oauth.GetUserHandler(srv))
	}

	log.Fatal(http.ListenAndServe(":9096", router))
}

func setupRouter() *gin.Engine {
	r := gin.Default()

	r.Use(cors.Default())

	//r.LoadHTMLGlob("static/*.html")

	return r
}

func setupOauthServer() *server.Server {
	srv := server.NewServer(server.NewConfig(), setupOauthManager())

	srv.ClientInfoHandler = server.ClientFormHandler

	srv.SetUserAuthorizationHandler(oauth.UserAuthorizationHandler(config))
	srv.SetInternalErrorHandler(oauth.InternalErrorHandler)
	srv.SetResponseErrorHandler(oauth.ResponseErrorHandler)
	srv.SetAllowGetAccessRequest(true)

	return srv
}

func setupOauthManager() *manage.Manager {
	manager := manage.NewDefaultManager()

	accessGenerate := generates.NewJWTAccessGenerate(
		config.JWTKey, jwt.SigningMethodHS512)

	cfg := &manage.Config{
		AccessTokenExp:    config.AccessTokenExpiration,
		RefreshTokenExp:   config.RefreshTokenExpiration,
		IsGenerateRefresh: config.GenerateRefresh,
	}

	clientStore := store.NewClientStore()
	clientStore.Set(config.Client.ID, &models.Client{
		ID:     config.Client.ID,
		Secret: config.Client.Secret,
		Domain: config.Client.Domain,
	})

	manager.SetAuthorizeCodeExp(config.AuthorizeCodeExpiration)
	manager.MapAccessGenerate(accessGenerate)
	manager.SetAuthorizeCodeTokenCfg(cfg)
	manager.MustTokenStorage(store.NewMemoryTokenStore())

	manager.MapClientStorage(clientStore)

	return manager
}
