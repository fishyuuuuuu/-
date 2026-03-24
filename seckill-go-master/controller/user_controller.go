package controller

import (
	"errors"
	"net/http"
	"seckill_go/model"
	"seckill_go/service"
	"seckill_go/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// LoginRequest 登录请求结构
type LoginRequest struct {
	Username    string `json:"username" binding:"required_without=Phone,omitempty,min=2,max=50"`
	Phone       string `json:"phone" binding:"required_without=Username,omitempty,min=5,max=20"`
	Password    string `json:"password" binding:"required,min=6,max=72"`
	CaptchaID   string `json:"captcha_id" binding:"required"`
	CaptchaCode string `json:"captcha_code" binding:"required,min=4,max=8"`
}

// RegisterRequest 注册请求结构
type RegisterRequest struct {
	Username    string `json:"username" binding:"required,min=2,max=50"`
	Phone       string `json:"phone" binding:"required,min=5,max=20"`
	Password    string `json:"password" binding:"required,min=6,max=72"`
	CaptchaID   string `json:"captcha_id" binding:"required"`
	CaptchaCode string `json:"captcha_code" binding:"required,min=4,max=8"`
}

// RegisterHandler 处理用户注册请求
func RegisterHandler(c *gin.Context) {
	utils.Logger.Info("处理用户注册请求")

	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Logger.Warn("注册请求数据无效", zap.String("error", err.Error()))
		utils.AuditRegister(req.Username, false, c.ClientIP())
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的注册数据：" + err.Error(),
		})
		return
	}

	// 验证验证码
	if !utils.VerifyCaptcha(req.CaptchaID, req.CaptchaCode) {
		utils.Logger.Warn("验证码错误",
			zap.String("captcha_id", req.CaptchaID))
		utils.AuditRegister(req.Username, false, c.ClientIP())
		c.JSON(http.StatusBadRequest, gin.H{"error": "验证码错误或已过期"})
		return
	}

	user := model.User{
		Username: req.Username,
		Phone:    req.Phone,
		Password: req.Password,
	}

	utils.Logger.Info("用户注册数据验证成功",
		zap.String("username", user.Username),
		zap.String("phone", user.Phone),
	)
	if err := service.Register(user); err != nil {
		utils.Logger.Error("用户注册失败", zap.Error(err))
		utils.AuditRegister(user.Username, false, c.ClientIP())
		switch {
		case errors.Is(err, service.ErrUserExists):
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器内部错误"})
		}
		return
	}
	utils.AuditRegister(user.Username, true, c.ClientIP())
	c.JSON(http.StatusOK, gin.H{"message": "注册成功"})
}

// LoginHandler 处理用户登录请求
func LoginHandler(c *gin.Context) {
	utils.Logger.Info("处理用户登录请求")

	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Logger.Warn("登录请求数据无效", zap.String("error", err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的登录数据：" + err.Error(),
		})
		return
	}

	// 验证验证码
	if !utils.VerifyCaptcha(req.CaptchaID, req.CaptchaCode) {
		utils.Logger.Warn("验证码错误",
			zap.String("captcha_id", req.CaptchaID))
		c.JSON(http.StatusBadRequest, gin.H{"error": "验证码错误或已过期"})
		return
	}

	// 构建用户对象
	user := model.User{
		Username: req.Username,
		Phone:    req.Phone,
		Password: req.Password,
	}

	utils.Logger.Info("用户登录数据验证成功",
		zap.String("username", user.Username),
		zap.String("phone", user.Phone),
	)
	id, err := service.Login(user)
	if err != nil {
		utils.Logger.Error("用户登录失败", zap.Error(err))
		utils.AuditLogin(user.Username, false, c.ClientIP())
		switch {
		case errors.Is(err, service.ErrUserLogin):
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器内部错误"})
			return
		}
		return
	}
	utils.AuditLogin(user.Username, true, c.ClientIP())
	token, err := utils.GenerateToken(id)
	if err != nil {
		utils.Logger.Error("生成令牌失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器内部错误"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "登录成功",
		"token":   token,
		"user_id": id,
	})
}

// GetUserListHandler 获取用户列表
func GetUserListHandler(c *gin.Context) {
	utils.Logger.Info("处理用户列表请求")

	// 获取查询参数
	page := c.DefaultQuery("page", "1")
	pageSize := c.DefaultQuery("pageSize", "10")
	username := c.Query("username")
	phone := c.Query("phone")

	// 转换为整数
	pageInt := utils.StringToInt(page)
	pageSizeInt := utils.StringToInt(pageSize)
	if pageInt <= 0 {
		pageInt = 1
	}
	if pageSizeInt <= 0 || pageSizeInt > 100 {
		pageSizeInt = 10
	}

	// 如果有条件查询，使用条件查询
	if username != "" || phone != "" {
		condition := make(map[string]interface{})
		if username != "" {
			condition["username"] = username
		}
		if phone != "" {
			condition["phone"] = phone
		}

		users, err := service.GetUsersByCondition(condition)
		if err != nil {
			utils.Logger.Error("获取用户列表失败", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "获取用户列表失败"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"users":    users,
			"total":    len(users),
			"page":     pageInt,
			"pageSize": pageSizeInt,
		})
		return
	}

	// 否则使用分页查询
	users, total, err := service.GetAllUsers(pageInt, pageSizeInt)
	if err != nil {
		utils.Logger.Error("获取用户列表失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取用户列表失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"users":      users,
		"total":      total,
		"page":       pageInt,
		"pageSize":   pageSizeInt,
		"totalPages": (int(total) + pageSizeInt - 1) / pageSizeInt,
	})
}

// GetUserByIDHandler 根据ID获取用户详情
func GetUserByIDHandler(c *gin.Context) {
	utils.Logger.Info("处理用户详情请求")

	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户ID不能为空"})
		return
	}

	idUint := utils.StringToUint(id)
	if idUint == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
		return
	}

	user, err := service.GetUserByID(idUint)
	if err != nil {
		if err.Error() == "用户不存在" {
			c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
			return
		}
		utils.Logger.Error("获取用户详情失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取用户详情失败"})
		return
	}

	// 隐藏敏感信息
	user.Password = ""

	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}

// UpdateUserHandler 更新用户信息
func UpdateUserHandler(c *gin.Context) {
	utils.Logger.Info("处理更新用户请求")

	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户ID不能为空"})
		return
	}

	idUint := utils.StringToUint(id)
	if idUint == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
		return
	}

	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
		return
	}

	user.ID = idUint

	if err := service.UpdateUserInfo(&user); err != nil {
		utils.Logger.Error("更新用户信息失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新用户信息失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "用户信息更新成功",
	})
}

// DeleteUserHandler 删除用户
func DeleteUserHandler(c *gin.Context) {
	utils.Logger.Info("处理删除用户请求")

	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户ID不能为空"})
		return
	}

	idUint := utils.StringToUint(id)
	if idUint == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
		return
	}

	if err := service.DeleteUserByID(idUint); err != nil {
		utils.Logger.Error("删除用户失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除用户失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "用户删除成功",
	})
}

// AuthMiddleware JWT认证中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头获取令牌
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "未提供认证令牌"})
			c.Abort()
			return
		}

		// 检查令牌格式
		if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的认证令牌格式"})
			c.Abort()
			return
		}

		// 提取令牌
		tokenString := authHeader[7:]

		// 验证令牌
		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		// 将用户ID存储到上下文中
		c.Set("user_id", claims.UserID)
		c.Set("userID", claims.UserID)
		c.Next()
	}
}
