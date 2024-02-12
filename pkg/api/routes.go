package api

// registerRoutes defines all routes in our api
func (app *App) registerRoutes() {
	app.router.POST("/signup", app.createNewUser )
	app.router.POST("/signin", app.signIn)
	app.router.POST("/refresh-token", app.refreshToken)
	orgGroup := app.router.Group("/organization")
	{
		orgGroup.POST("",app.createNewOrganization)
		orgGroup.GET("",app.getAllOrganizations)
		orgGroup.GET("/:organization_id", app.getOrganizationByID)
		orgGroup.PUT("/:organization_id",app.updateOrganization)
		orgGroup.DELETE("/:organization_id", app.deleteOrganization)
		orgGroup.POST("/:organization_id/invite",app.inviteUser)

	}
}