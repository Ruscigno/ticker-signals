package server

import (
	"github.com/Ruscigno/ticker-signals/internal/config"
	"github.com/gin-gonic/gin"
)

func registerRoutes(router *gin.Engine, conf *config.AppConfig) {
	// Enables automatic redirection if the current route cannot be matched but a
	// handler for the path with (without) the trailing slash exists.
	router.RedirectTrailingSlash = true

	// JSON-REST API Version 1
	// v1 := router.Group(conf.BaseUri(config.ApiUri))
	// {
	// // Config options.
	// api.GetConfig(v1)
	// api.GetConfigOptions(v1)
	// api.SaveConfigOptions(v1)

	// // User profile and settings.
	// api.GetSettings(v1)
	// api.SaveSettings(v1)
	// api.ChangePassword(v1)
	// api.CreateSession(v1)
	// api.DeleteSession(v1)

	// // External account management.
	// api.SearchAccounts(v1)
	// api.GetAccount(v1)
	// api.GetAccountFolders(v1)
	// api.ShareWithAccount(v1)
	// api.CreateAccount(v1)
	// api.DeleteAccount(v1)
	// api.UpdateAccount(v1)
	// }

	// // Default HTML page for client-side rendering and routing via VueJS.
	// router.NoRoute(func(c *gin.Context) {
	// 	signUp := gin.H{"message": config.MsgSponsor, "url": config.SignUpURL}
	// 	values := gin.H{"signUp": signUp, "config": conf.PublicConfig()}
	// 	c.HTML(http.StatusOK, conf.TemplateName(), values)
	// })
}
