package cmd

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/szaffarano/oauth-server/domain"
	"github.com/szaffarano/oauth-server/oauth"
	"gopkg.in/oauth2.v3/generates"
	"gopkg.in/oauth2.v3/manage"
	"gopkg.in/oauth2.v3/models"
	"gopkg.in/oauth2.v3/server"
	"gopkg.in/oauth2.v3/store"
	"gopkg.in/yaml.v2"
)

var (
	cfgFile            string
	printConfigExample bool

	config domain.Config

	rootCmd = &cobra.Command{
		Use:   "oauth-server-provider",
		Short: "Oauth server provider",
		Long:  `Oauth server provider`,
		Run: func(cmd *cobra.Command, args []string) {
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

			log.Fatal(http.ListenAndServe(
				fmt.Sprintf(":%d", config.App.Port), router))
		},
	}
)

// Execute is the CLI entry point
func Execute() {
	if rootCmd.Execute() != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initViper)

	// init global flags
	rootCmd.
		PersistentFlags().
		BoolVarP(
			&config.App.Verbose,
			"verbose",
			"v",
			false,
			"Print verbose info")

	rootCmd.
		PersistentFlags().
		StringVarP(
			&cfgFile,
			"config",
			"c",
			"",
			`Configure cli through a yaml file`)

	rootCmd.
		PersistentFlags().
		BoolVarP(
			&printConfigExample,
			"printConfigExample",
			"p",
			false,
			`Print configuration template`)
}

func initViper() {
	if printConfigExample {
		s, _ := yaml.Marshal(config)
		fmt.Println(string(s))
		os.Exit(0)
	}

	viperConfig := viper.New()

	if cfgFile != "" {
		viperConfig.SetConfigFile(cfgFile)
	} else {
		viperConfig.AddConfigPath(".")
		viperConfig.SetConfigName("oauth2-mm")
	}

	viperConfig.SetConfigType("yaml")
	viperConfig.AutomaticEnv()
	viperConfig.SetEnvPrefix("OAUTH2_MM")

	if err := viperConfig.ReadInConfig(); err != nil {
		log.Fatal("Error obteniendo configuración: ", err)
	}

	if err := viperConfig.Unmarshal(&config); err != nil {
		log.Fatal("Error leyendo configuración: ", err)
	}

	if err := config.Validate(); err != nil {
		log.Fatal("Archivo de configuración inválido: ", err)
	}
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
		[]byte(config.Oauth.JWTKey), jwt.SigningMethodHS512)

	cfg := &manage.Config{
		AccessTokenExp:    config.Oauth.AccessTokenExpiration,
		RefreshTokenExp:   config.Oauth.RefreshTokenExpiration,
		IsGenerateRefresh: config.Oauth.GenerateRefresh,
	}

	clientStore := store.NewClientStore()
	clientStore.Set(config.Client.ID, &models.Client{
		ID:     config.Client.ID,
		Secret: config.Client.Secret,
		Domain: config.Client.Domain,
	})

	manager.SetAuthorizeCodeExp(config.Oauth.AuthorizeCodeExpiration)
	manager.MapAccessGenerate(accessGenerate)
	manager.SetAuthorizeCodeTokenCfg(cfg)
	manager.MustTokenStorage(store.NewMemoryTokenStore())

	manager.MapClientStorage(clientStore)

	return manager
}
