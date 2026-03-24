package controller

import (
	"net/http"
	"seckill_go/service"
	"seckill_go/utils"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// GetProductListHandler 获取商品列表请求
func GetProductListHandler(c *gin.Context) {
	products, err := service.GetProductList(c)
	if err != nil {
		utils.Logger.Error("获取商品列表失败", zap.Error(err))
		c.JSON(500, gin.H{
			"code":    500,
			"message": "无法获取商品列表",
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    products,
	})
}

// SeckillHandler 秒杀请求
func SeckillHandler(c *gin.Context) {
	// 从请求参数或生成模拟用户ID（方便测试）
	type SeckillRequest struct {
		ProductId  uint   `json:"productId"`
		CaptchaId  string `json:"captchaId"`
		CaptchaStr string `json:"captchaStr"`
		UserId     uint   `json:"userId"` // 可选，方便测试
	}
	var req SeckillRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	userID := req.UserId
	if contextUserID, exists := c.Get("user_id"); exists {
		if uid, ok := contextUserID.(uint); ok && uid > 0 {
			userID = uid
		}
	}
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未登录"})
		return
	}

	// 校验验证码，防止机器人刷单
	if req.CaptchaId == "" || req.CaptchaStr == "" || !utils.VerifyCaptcha(req.CaptchaId, req.CaptchaStr) {
		// 记录秒杀失败审计日志
		utils.AuditSeckill(strconv.Itoa(int(userID)), strconv.Itoa(int(req.ProductId)), false, c.ClientIP())
		c.JSON(http.StatusBadRequest, gin.H{"error": "验证码错误"})
		return
	}

	// 限流检查已在中间件完成，这里不再重复检查

	utils.Logger.Info("处理秒杀请求，商品id:" + strconv.Itoa(int(req.ProductId)) + "，用户id:" + strconv.Itoa(int(userID)))
	err := service.ReduceProductStock(c, req.ProductId, strconv.Itoa(int(userID)))
	if err != nil {
		// 记录秒杀失败审计日志
		utils.AuditSeckill(strconv.Itoa(int(userID)), strconv.Itoa(int(req.ProductId)), false, c.ClientIP())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "秒杀失败"})
		return
	}
	// 记录秒杀成功审计日志
	utils.AuditSeckill(strconv.Itoa(int(userID)), strconv.Itoa(int(req.ProductId)), true, c.ClientIP())
	c.JSON(http.StatusOK, gin.H{"message": "秒杀成功"})
}

// CreateProductHandler 创建商品请求
func CreateProductHandler(c *gin.Context) {
	var product map[string]interface{}
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	// 验证必要字段
	if _, ok := product["name"]; !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "商品名称不能为空"})
		return
	}
	if _, ok := product["price"]; !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "商品价格不能为空"})
		return
	}
	if _, ok := product["stock"]; !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "商品库存不能为空"})
		return
	}

	err := service.CreateProduct(c, product)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建商品失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "商品创建成功"})
}

// UpdateProductHandler 更新商品请求
func UpdateProductHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的商品ID"})
		return
	}

	var product map[string]interface{}
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	err = service.UpdateProduct(c, id, product)
	if err != nil {
		if err == service.ErrProductNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "商品不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新商品失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "商品更新成功"})
}

// DeleteProductHandler 删除商品请求
func DeleteProductHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的商品ID"})
		return
	}

	err = service.DeleteProduct(c, id)
	if err != nil {
		if err == service.ErrProductNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "商品不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除商品失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "商品删除成功"})
}

// GetProductByIDHandler 根据ID获取商品请求
func GetProductByIDHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的商品ID"})
		return
	}

	product, err := service.GetProductByID(c, id)
	if err != nil {
		if err == service.ErrProductNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "商品不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取商品失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    product,
	})
}

// ResetProductStockHandler 重置商品库存（用于性能测试）
func ResetProductStockHandler(c *gin.Context) {
	type ResetStockRequest struct {
		ProductId uint `json:"productId" binding:"required"`
		Stock     int  `json:"stock" binding:"required"`
	}
	var req ResetStockRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	err := service.ResetProductStock(c, req.ProductId, req.Stock)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "库存重置成功",
		"productId": req.ProductId,
		"stock":     req.Stock,
	})
}

// ClearTestDataHandler 清理测试数据（用于重置测试环境）
func ClearTestDataHandler(c *gin.Context) {
	utils.Logger.Info("清理测试数据")

	// 清空订单数据
	service.ClearMockOrders()

	// 清理Redis中的限流键
	utils.ClearRateLimitCache()

	c.JSON(http.StatusOK, gin.H{
		"message": "测试数据已清理",
	})
}
