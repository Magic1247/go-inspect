package routes

import (
	"go-inspect/controllers"
	"go-inspect/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	// 公开路由
	r.SetTrustedProxies([]string{"127.0.0.1"})
	public := r.Group("/api")
	{
		public.POST("/register", controllers.Register)
		public.POST("/login", controllers.Login)
	}

	// 需要认证的路由
	protected := r.Group("/api")
	protected.Use(middleware.JWTAuth())
	{
		// 用户管理路由
		protected.GET("/user", controllers.GetUserInfo)
		protected.PUT("/user", controllers.UpdateUserInfo)
		protected.POST("/user/changePassword", controllers.ChangePassword)

		// 巡检点位管理路由
		inspectionPoints := protected.Group("/inspectionPoints")
		{
			inspectionPoints.POST("/", controllers.CreateInspectionPoint)
			inspectionPoints.GET("/", controllers.ListInspectionPoints)
			inspectionPoints.GET("/:id", controllers.GetInspectionPoint)
			inspectionPoints.PUT("/:id", controllers.UpdateInspectionPoint)
			inspectionPoints.DELETE("/:id", controllers.DeleteInspectionPoint)
		}

		// 巡检路线管理路由
		inspectionRoutes := protected.Group("/inspectionRoutes")
		{
			inspectionRoutes.POST("/", controllers.CreateInspectionRoute)
			inspectionRoutes.GET("/", controllers.ListInspectionRoutes)
			inspectionRoutes.GET("/:id", controllers.GetInspectionRoute)
			inspectionRoutes.PUT("/:id", controllers.UpdateInspectionRoute)
			inspectionRoutes.DELETE("/:id", controllers.DeleteInspectionRoute)
			inspectionRoutes.POST("/:id/points", controllers.AddPointToRoute)
			inspectionRoutes.DELETE("/:id/points/:pointId", controllers.RemovePointFromRoute)
		}

		// 巡检项路由
		inspectionItems := protected.Group("/inspectionItems")
		{
			inspectionItems.POST("/", controllers.CreateInspectionItem)
			inspectionItems.GET("/", controllers.ListInspectionItems)
			inspectionItems.GET("/:id", controllers.GetInspectionItem)
			inspectionItems.PUT("/:id", controllers.UpdateInspectionItem)
			inspectionItems.DELETE("/:id", controllers.DeleteInspectionItem)
			inspectionItems.POST("/:id/points", controllers.AddItemToPoint)
			inspectionItems.DELETE("/:id/points/:pointId", controllers.RemoveItemFromPoint)
		}

		// 项目管理路由
		projects := protected.Group("/projects")
		{
			projects.POST("/", controllers.CreateProject)
			projects.GET("/", controllers.ListProjects)
			projects.GET("/:id", controllers.GetProject)
			projects.PUT("/:id", controllers.UpdateProject)
			projects.DELETE("/:id", controllers.DeleteProject)
			projects.GET("/:id/tree", controllers.GetProjectTree)
		}

		// 巡检计划管理路由
		inspectionPlans := protected.Group("/inspectionPlans")
		{
			inspectionPlans.POST("/", controllers.CreateInspectionPlan)
			inspectionPlans.GET("/", controllers.ListInspectionPlans)
			inspectionPlans.GET("/:id", controllers.GetInspectionPlan)
			inspectionPlans.PUT("/:id", controllers.UpdateInspectionPlan)
			inspectionPlans.DELETE("/:id", controllers.DeleteInspectionPlan)
			inspectionPlans.POST("/:id/trigger", controllers.TriggerInspectionPlan)
		}

		// 巡检工单管理路由
		inspectionOrders := protected.Group("/inspectionOrders")
		{
			inspectionOrders.GET("/", controllers.ListInspectionOrders)
			inspectionOrders.GET("/:id", controllers.GetInspectionOrder)
			inspectionOrders.POST("/:id/assign", controllers.AssignInspectionOrder)
			inspectionOrders.POST("/:id/start", controllers.StartInspectionOrder)
			inspectionOrders.POST("/:id/complete", controllers.CompleteInspectionOrder)
			inspectionOrders.POST("/:id/points/:pointId/check", controllers.CheckInspectionPoint)
		}
	}
}
