package handler

import (
	"context"
	"fmt"
	"net/http"

	coreServer "github.com/red-gold/telar-core/server"
	"github.com/red-gold/telar-web/src/controllers"
	cf "github.com/red-gold/telar-web/src/controllers/users/auth/config"
	"github.com/red-gold/telar-web/src/controllers/users/auth/handlers"
	// handlers "github.com/red-gold/telar-web/src/controllers/users/auth/handlers"
)

// Cache state
var server *coreServer.ServerRouter
var db interface{}

func init() {

	cf.InitConfig()
}

// Handler function
func Handle(w http.ResponseWriter, r *http.Request) {

	ctx := context.Background()

	// Start
	if db == nil {
		var startErr error
		db, startErr = controllers.Start(ctx)
		if startErr != nil {
			fmt.Printf("Error startup: %s", startErr.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(startErr.Error()))
		}
	}

	// Server Routing
	if server == nil {

		// Routing
		server = coreServer.NewServerRouter()

		// Admin
		server.POSTWR("/check/admin", handlers.CheckAdminHandler(db), coreServer.RouteProtectionHMAC)
		server.POSTWR("/signup/admin", handlers.AdminSignupHandle(db), coreServer.RouteProtectionHMAC)
		server.POSTWR("/login/admin", handlers.LoginAdminHandler(db), coreServer.RouteProtectionHMAC)

		// Signup
		server.POSTWR("/signup/verify", handlers.VerifySignupHandle(db), coreServer.RouteProtectionPublic)
		server.POSTWR("/signup", handlers.SignupTokenHandle(db), coreServer.RouteProtectionPublic)
		server.GET("/signup", handlers.SignupPageHandler, coreServer.RouteProtectionPublic)

		// Password
		server.GETWR("/password/reset/:verifyId", handlers.ResetPasswordPageHandler(db), coreServer.RouteProtectionPublic)
		server.POSTWR("/password/reset/:verifyId", handlers.ResetPasswordFormHandler(db), coreServer.RouteProtectionPublic)
		server.GET("/password/forget", handlers.ForgetPasswordPageHandler, coreServer.RouteProtectionPublic)
		server.POSTWR("/password/forget", handlers.ForgetPasswordFormHandler(db), coreServer.RouteProtectionPublic)

		// Login
		server.GET("/login", handlers.LoginPageHandler, coreServer.RouteProtectionPublic)
		server.POSTWR("/login", handlers.LoginTelarHandler(db), coreServer.RouteProtectionPublic)
		server.POSTWR("/login/telar", handlers.LoginTelarHandler(db), coreServer.RouteProtectionPublic)
		server.GETWR("/login/github", handlers.LoginGithubHandler, coreServer.RouteProtectionPublic)
		server.GETWR("/login/google", handlers.LoginGoogleHandler, coreServer.RouteProtectionPublic)
		server.GETWR("/oauth2/authorized", handlers.OAuth2Handler(db), coreServer.RouteProtectionPublic)
		server.GETWR("/oauth2/google", handlers.OauthGoogleCallback, coreServer.RouteProtectionPublic)

		// Profile
		server.PUTWR("/profile", handlers.UpdateProfileHandle(db), coreServer.RouteProtectionCookie)
	}
	server.ServeHTTP(w, r)
}
