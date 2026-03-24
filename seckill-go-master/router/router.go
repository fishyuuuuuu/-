package router

import (
	"seckill_go/controller"
	"seckill_go/utils"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// 初始化路由
func Init() *gin.Engine {
	r := gin.Default()
	// 指定信任的代理服务器
	if err := r.SetTrustedProxies([]string{"127.0.0.1"}); err != nil {
		return nil
	}
	// 全局中间件（跨域）
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://localhost:5175", "http://localhost:5176", "http://localhost:5177"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))
	// 性能监控中间件
	r.Use(utils.MonitoringMiddleware())
	// 安全中间件（限流、黑名单检查）
	r.Use(utils.SecurityMiddleware())
	// 静态文件服务
	r.Static("/temp", "./temp")
	r.Static("/uploads", "./uploads")
	// 路由分组：API 前缀
	api := r.Group("/api")
	{
		// 验证码路由
		api.GET("/captcha", utils.CaptchaHandler)

		// 安全管理路由（需要管理员权限）
		security := api.Group("/security")
		{
			// 黑名单管理
			security.GET("/blacklist/:type/:id", controller.AuthMiddleware(), controller.RBACMiddleware("security:read"), controller.GetBlacklistInfoHandler)
			security.POST("/blacklist", controller.AuthMiddleware(), controller.RBACMiddleware("security:write"), controller.AddToBlacklistHandler)
			security.DELETE("/blacklist/:type/:id", controller.AuthMiddleware(), controller.RBACMiddleware("security:write"), controller.RemoveFromBlacklistHandler)
			// 安全统计
			security.GET("/stats", controller.AuthMiddleware(), controller.RBACMiddleware("security:read"), controller.GetSecurityStatsHandler)
			// 限流配置
			security.GET("/rate-limit/config", controller.AuthMiddleware(), controller.RBACMiddleware("security:read"), controller.GetRateLimitConfigHandler)
			security.PUT("/rate-limit/config", controller.AuthMiddleware(), controller.RBACMiddleware("security:write"), controller.UpdateRateLimitConfigHandler)
		}

		// 用户相关路由
		user := api.Group("/user")
		{
			user.POST("/register", controller.RegisterHandler) // 注册接口
			user.POST("/login", controller.LoginHandler)       // 登录接口
			// 用户管理接口，暂时移除认证方便测试
			user.GET("/list", controller.GetUserListHandler)
			user.GET("/get/:id", controller.GetUserByIDHandler)
			user.PUT("/update/:id", controller.UpdateUserHandler)
			user.DELETE("/delete/:id", controller.DeleteUserHandler)
		}
		product := api.Group("/product")
		{
			product.GET("/list", controller.GetProductListHandler)
			// 秒杀接口使用增强的安全中间件
		product.POST("/seckill", utils.BehaviorAnalysisMiddleware(), utils.CaptchaSecurityMiddleware(), utils.RateLimitMiddleware(), controller.SeckillHandler)
			product.POST("/reset-stock", controller.ResetProductStockHandler) // 重置库存（测试用）
			product.POST("/clear-test-data", controller.ClearTestDataHandler) // 清理测试数据（测试用）
			// 商品管理接口，需要认证和权限
			product.POST("/create", controller.AuthMiddleware(), controller.RBACMiddleware("product:create"), controller.CreateProductHandler)
			product.PUT("/update/:id", controller.AuthMiddleware(), controller.RBACMiddleware("product:update"), controller.UpdateProductHandler)
			product.DELETE("/delete/:id", controller.AuthMiddleware(), controller.RBACMiddleware("product:delete"), controller.DeleteProductHandler)
			product.GET("/get/:id", controller.GetProductByIDHandler)
		}

		// 文件上传路由
		upload := api.Group("/upload")
		{
			upload.POST("/image", controller.AuthMiddleware(), controller.UploadImageHandler)
			upload.POST("/images", controller.AuthMiddleware(), controller.UploadMultipleImagesHandler)
			upload.DELETE("/image/:filename", controller.AuthMiddleware(), controller.DeleteImageHandler)
		}

		// 活动管理路由
		event := api.Group("/event")
		{
			event.GET("/list", controller.GetEventListHandler)
			event.GET("/get/:id", controller.GetEventByIDHandler)
			// 活动管理接口，需要认证和权限
			event.POST("/create", controller.AuthMiddleware(), controller.RBACMiddleware("event:create"), controller.CreateEventHandler)
			event.PUT("/update/:id", controller.AuthMiddleware(), controller.RBACMiddleware("event:update"), controller.UpdateEventHandler)
			event.DELETE("/delete/:id", controller.AuthMiddleware(), controller.RBACMiddleware("event:delete"), controller.DeleteEventHandler)
			event.POST("/start/:id", controller.AuthMiddleware(), controller.RBACMiddleware("event:start"), controller.StartEventHandler)
			event.POST("/stop/:id", controller.AuthMiddleware(), controller.RBACMiddleware("event:stop"), controller.StopEventHandler)
		}

		// 订单管理路由
		order := api.Group("/order")
		{
			// 用户订单接口，需要认证
			order.GET("/list", controller.AuthMiddleware(), controller.GetOrderListHandler)
			order.GET("/get/:id", controller.AuthMiddleware(), controller.GetOrderByIDHandler)
			order.POST("/cancel/:id", controller.AuthMiddleware(), controller.CancelOrderHandler)
			// 管理员订单接口，需要管理员权限
			order.GET("/admin/list", controller.AuthMiddleware(), controller.RBACMiddleware("order:read"), controller.GetAllOrdersHandler)
			order.GET("/admin/get/:id", controller.AuthMiddleware(), controller.RBACMiddleware("order:read"), controller.GetOrderByIDAdminHandler)
			order.PUT("/admin/status/:id", controller.AuthMiddleware(), controller.RBACMiddleware("order:update"), controller.UpdateOrderStatusHandler)
		}
	}
	return r
}
