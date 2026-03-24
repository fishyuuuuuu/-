package controller

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"seckill_go/utils"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// 允许的图片格式
var allowedImageFormats = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".gif":  true,
	".webp": true,
}

// 最大文件大小：10MB
const maxFileSize = 10 * 1024 * 1024

// UploadImageHandler 上传单张图片
func UploadImageHandler(c *gin.Context) {
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "请选择要上传的图片",
		})
		return
	}

	// 验证文件大小
	if file.Size > maxFileSize {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "文件大小不能超过10MB",
		})
		return
	}

	// 验证文件格式
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !allowedImageFormats[ext] {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "不支持的图片格式，仅支持 jpg、jpeg、png、gif、webp",
		})
		return
	}

	// 创建上传目录
	uploadDir := "./uploads/images"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		utils.Logger.Error("创建上传目录失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "上传失败",
		})
		return
	}

	// 生成唯一文件名
	filename := fmt.Sprintf("%s_%d%s", uuid.New().String(), time.Now().Unix(), ext)
	filepath := filepath.Join(uploadDir, filename)

	// 保存文件
	if err := c.SaveUploadedFile(file, filepath); err != nil {
		utils.Logger.Error("保存文件失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "上传失败",
		})
		return
	}

	// 返回文件URL
	fileURL := fmt.Sprintf("/uploads/images/%s", filename)
	utils.Logger.Info("图片上传成功", zap.String("filename", filename), zap.String("url", fileURL))

	c.JSON(http.StatusOK, gin.H{
		"code":     0,
		"message":  "上传成功",
		"data":     fileURL,
		"filename": file.Filename,
	})
}

// UploadMultipleImagesHandler 上传多张图片
func UploadMultipleImagesHandler(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "请选择要上传的图片",
		})
		return
	}

	files := form.File["images"]
	if len(files) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "请选择要上传的图片",
		})
		return
	}

	// 创建上传目录
	uploadDir := "./uploads/images"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		utils.Logger.Error("创建上传目录失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "上传失败",
		})
		return
	}

	var uploadedFiles []map[string]string
	var failedFiles []string

	for _, file := range files {
		// 验证文件大小
		if file.Size > maxFileSize {
			failedFiles = append(failedFiles, file.Filename+" (文件过大)")
			continue
		}

		// 验证文件格式
		ext := strings.ToLower(filepath.Ext(file.Filename))
		if !allowedImageFormats[ext] {
			failedFiles = append(failedFiles, file.Filename+" (格式不支持)")
			continue
		}

		// 生成唯一文件名
		filename := fmt.Sprintf("%s_%d%s", uuid.New().String(), time.Now().Unix(), ext)
		filePath := filepath.Join(uploadDir, filename)

		// 保存文件
		if err := c.SaveUploadedFile(file, filePath); err != nil {
			utils.Logger.Error("保存文件失败", zap.String("filename", file.Filename), zap.Error(err))
			failedFiles = append(failedFiles, file.Filename+" (保存失败)")
			continue
		}

		// 返回文件URL
		fileURL := fmt.Sprintf("/uploads/images/%s", filename)
		uploadedFiles = append(uploadedFiles, map[string]string{
			"url":      fileURL,
			"filename": file.Filename,
		})
	}

	utils.Logger.Info("批量图片上传完成",
		zap.Int("success", len(uploadedFiles)),
		zap.Int("failed", len(failedFiles)))

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": fmt.Sprintf("成功上传 %d 张图片", len(uploadedFiles)),
		"data":    uploadedFiles,
		"failed":  failedFiles,
	})
}

// DeleteImageHandler 删除图片
func DeleteImageHandler(c *gin.Context) {
	filename := c.Param("filename")
	if filename == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "文件名不能为空",
		})
		return
	}

	// 防止路径遍历攻击
	if strings.Contains(filename, "..") || strings.Contains(filename, "/") {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的文件名",
		})
		return
	}

	filepath := fmt.Sprintf("./uploads/images/%s", filename)
	if err := os.Remove(filepath); err != nil {
		if os.IsNotExist(err) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "文件不存在",
			})
			return
		}
		utils.Logger.Error("删除文件失败", zap.String("filename", filename), zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "删除失败",
		})
		return
	}

	utils.Logger.Info("图片删除成功", zap.String("filename", filename))
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "删除成功",
	})
}
