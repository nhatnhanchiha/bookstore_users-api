package app

import (
	"github.com/nhatnhanchiha/bookstore_users-api/controllers/ping"
	"github.com/nhatnhanchiha/bookstore_users-api/controllers/user"
)

func mapUrls() {
	router.GET("/ping", ping.Ping)

	//region For Users
	router.GET("/users/:user_id", user.Get)
	router.POST("/users", user.Create)
	router.PUT("/users/:user_id", user.Update)
	router.PATCH("/users/:user_id", user.Update)
	router.DELETE("/users/:user_id", user.Delete)
	router.GET("/internal/users/search", user.Search)
	router.POST("/users/login", user.Login)
	//endregion

}
