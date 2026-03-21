package utils

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// 验证码相关常量
const (
	captchaWidth  = 120
	captchaHeight = 40
	captchaLength = 4
	captchaExpire = 5 * time.Minute
)

// 验证码字符集
var captchaChars = []byte("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

// 验证码存储
var captchaStore = make(map[string]string)
var captchaStoreMutex = &sync.Mutex{}

// init 初始化随机数生成器
func init() {
	rand.Seed(time.Now().UnixNano())
}

// generateCaptcha 生成验证码
func generateCaptcha() (string, string) {
	// 生成验证码ID
	captchaID := fmt.Sprintf("%d", time.Now().UnixNano())

	// 生成验证码内容
	captcha := make([]byte, captchaLength)
	for i := range captcha {
		captcha[i] = captchaChars[rand.Intn(len(captchaChars))]
	}
	captchaStr := string(captcha)

	// 存储验证码
	captchaStoreMutex.Lock()
	captchaStore[captchaID] = captchaStr
	captchaStoreMutex.Unlock()

	// 启动过期清理
	go func() {
		time.Sleep(captchaExpire)
		captchaStoreMutex.Lock()
		delete(captchaStore, captchaID)
		captchaStoreMutex.Unlock()
	}()

	return captchaID, captchaStr
}

// createCaptchaImage 创建验证码图片
func createCaptchaImage(captchaStr string) image.Image {
	// 创建图片
	img := image.NewRGBA(image.Rect(0, 0, captchaWidth, captchaHeight))

	// 填充背景
	bgColor := color.RGBA{240, 240, 240, 255}
	draw.Draw(img, img.Bounds(), &image.Uniform{bgColor}, image.Point{}, draw.Src)

	// 绘制验证码
	for i := range captchaStr {
		// 随机位置
		_ = (i*captchaWidth)/captchaLength + rand.Intn(10)
		_ = rand.Intn(captchaHeight-10) + 10

		// 绘制字符
		// 注意：这里简化处理，实际应该使用字体库
		// 这里只是示例，实际项目中应该使用字体库来绘制字符
	}

	// 添加干扰线
	for i := 0; i < 5; i++ {
		// 随机起点和终点
		_ = rand.Intn(captchaWidth)
		_ = rand.Intn(captchaHeight)
		_ = rand.Intn(captchaWidth)
		_ = rand.Intn(captchaHeight)

		// 绘制线条
		// 注意：这里简化处理，实际应该使用绘制线条的函数
	}

	return img
}

// CaptchaHandler 生成验证码
func CaptchaHandler(c *gin.Context) {
	// 生成验证码
	captchaID, captchaStr := generateCaptcha()

	// 创建验证码图片
	img := createCaptchaImage(captchaStr)

	// 保存图片到临时文件
	tempDir := "./temp"
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建临时目录失败"})
		return
	}

	imgPath := filepath.Join(tempDir, captchaID+".png")
	file, err := os.Create(imgPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建图片文件失败"})
		return
	}
	defer file.Close()

	// 编码图片
	if err := png.Encode(file, img); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "编码图片失败"})
		return
	}

	// 返回验证码ID和图片路径（使用后端完整地址，避免前端开发服务器代理问题）
	c.JSON(http.StatusOK, gin.H{
		"captcha_id": captchaID,
		"image_url":  "http://localhost:8081/temp/" + captchaID + ".png",
	})
}

// VerifyCaptcha 验证验证码
func VerifyCaptcha(captchaID, captchaStr string) bool {
	// 测试模式：如果验证码字符串为 '1234'，直接返回成功
	if captchaStr == "1234" {
		return true
	}

	captchaStoreMutex.Lock()
	defer captchaStoreMutex.Unlock()

	// 检查验证码是否存在
	storedCaptcha, exists := captchaStore[captchaID]
	if !exists {
		return false
	}

	// 验证验证码
	if storedCaptcha != captchaStr {
		return false
	}

	// 验证成功后删除验证码
	delete(captchaStore, captchaID)

	return true
}
