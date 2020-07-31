package app

import (
	"github.com/gin-gonic/gin"
	"github.com/studingprojects/bookstore_oauth-api/src/client/cassandra"
	"github.com/studingprojects/bookstore_oauth-api/src/http"
	"github.com/studingprojects/bookstore_oauth-api/src/repository/db"
	"github.com/studingprojects/bookstore_oauth-api/src/repository/rest"
	"github.com/studingprojects/bookstore_oauth-api/src/service/access_token"
)

var (
	router = gin.Default()
)

func StartApplication() {
	session := cassandra.GetSession()
	defer session.Close()

	atService := access_token.NewService(db.NewRepository(), rest.NewRepository())
	atHandler := http.NewAccessTokenHandler(atService)

	router.GET("/oauth/access_token/:access_token_id", atHandler.GetById)
	router.POST("/oauth/login", atHandler.Login)
	router.POST("/oauth/token", atHandler.Create)
	router.Run(":8089")
}
