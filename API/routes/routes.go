package routes

import (
	"github.com/gin-gonic/gin"
	//Importation des controleurs
	"llio-api/controllers"
	"llio-api/middleware"
	"llio-api/models/enums"
)

// Importation du middleware
func RegisterRoutes(r *gin.Engine) {
	/*------------------- Users -------------------*/
	userGroup := r.Group("/user", middleware.RoleValidationMiddleware(enums.Employee))
	{
		userGroup.GET("/me", controllers.GetUserInfo)
		userGroup.DELETE("/:id", middleware.RoleValidationMiddleware(enums.Administrator), controllers.DeleteUserById)
		userGroup.PATCH("/:id/role", middleware.RoleValidationMiddleware(enums.Administrator), controllers.UpdateUserRole)
		userGroup.POST("/calculate", controllers.CalculateTimeBank)
	}

	usersGroup := r.Group("/users", middleware.RoleValidationMiddleware(enums.ProjectManager))
	{

		usersGroup.GET("", controllers.GetAllUsers)
	}

	/*------------------- Activities -------------------*/
	activityGroup := r.Group("/activity", middleware.RoleValidationMiddleware((enums.Employee)))
	{
		activityGroup.POST("", controllers.CreateActivity)
		activityGroup.GET("/:id", controllers.GetActivityById)
		activityGroup.PUT("", controllers.UpdateActivity)
		activityGroup.DELETE("/:id", controllers.DeleteActivity)

	}

	activitiesGroup := r.Group("/activities", middleware.RoleValidationMiddleware(enums.Employee))
	{
		usersActivitiesGroup := activitiesGroup.Group("/me")
		{
			usersActivitiesGroup.GET("", controllers.GetActivitiesFromRange)
			usersActivitiesGroup.GET("/detailed", controllers.GetDetailedActivitiesFromRange)
		}
	}

	/*------------------- Categories -------------------*/
	categoryGroup := r.Group("/category", middleware.RoleValidationMiddleware(enums.Employee))
	{
		categoryGroup.POST("", controllers.CreateCategory)
		categoryGroup.GET("/:id", controllers.GetCategoryById)
		categoryGroup.PUT("", controllers.UpdateCategory)
	}

	categoriesGroup := r.Group("/categories", middleware.RoleValidationMiddleware(enums.Employee))
	{
		categoriesGroup.GET("", controllers.GetCategories)
	}

	/*------------------- Projets -------------------*/
	projectGroup := r.Group("/project")
	{
		projectGroup.POST("", middleware.RoleValidationMiddleware(enums.ProjectManager), controllers.CreatedProject)
		projectGroup.GET("/:id", middleware.RoleValidationMiddleware(enums.Employee), controllers.GetProjectById)
		projectGroup.PUT("", middleware.RoleValidationMiddleware(enums.ProjectManager), controllers.UpdateProject)
		projectGroup.GET("/:id/categories", middleware.RoleValidationMiddleware(enums.Employee), controllers.GetCategoriesByProjectId)
	}

	projectsGroup := r.Group("/projects", middleware.RoleValidationMiddleware(enums.Employee))
	{
		projectsGroup.GET("", controllers.GetProjects)
		projectsGroup.GET("/detailed", controllers.GetDetailedProjects)
		projectsGroup.GET("/me/detailed", controllers.GetDetailedProjectsByUser)
	}

}

// Routes pour l'authentification
func AuthRoutes(r *gin.Engine) {
	r.GET("/auth/:provider/callback", controllers.GetAuthCallback)
	r.POST("/logout", controllers.Logout)
	r.GET("/auth/:provider", controllers.Auth)
}
