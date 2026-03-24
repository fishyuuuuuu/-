package service

import (
	"errors"
	"seckill_go/db"
	"seckill_go/model"
	"seckill_go/utils"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

// ErrUserExists 自定义错误
var (
	ErrUserExists = errors.New("用户已存在")     // 用户已存在错误
	ErrUserLogin  = errors.New("用户帐号或密码错误") // 用户帐号或密码错误
)

func Register(user model.User) error {
	// 检查用户是否已存在
	existingID, err := db.CheckUserPhone(user.Phone)
	if err != nil {
		utils.Logger.Error("检查用户手机号失败", zap.String("phone", user.Phone), zap.Error(err))
		return err
	}

	if existingID > 0 {
		utils.Logger.Warn("用户已存在", zap.String("phone", user.Phone))
		return ErrUserExists
	}

	// 加密密码
	hashed, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		utils.Logger.Error("加密密码失败", zap.Error(err))
		return err
	}
	user.Password = string(hashed)

	// 创建用户到数据库
	userID, err := db.CreateUser(&user)
	if err != nil {
		if err.Error() == "索引冲突" {
			utils.Logger.Warn("用户已存在（索引冲突）", zap.String("phone", user.Phone))
			return ErrUserExists
		}
		utils.Logger.Error("创建用户失败", zap.Error(err))
		return err
	}

	utils.Logger.Info("用户注册成功",
		zap.String("username", user.Username),
		zap.String("phone", user.Phone),
		zap.Uint("id", userID))
	return nil
}

func Login(user model.User) (uint, error) {
	// 检查用户密码
	isValid, err := db.CheckUserPassword(user.Phone, user.Password)
	if err != nil {
		if err.Error() == "用户不存在" || err.Error() == "密码错误" {
			utils.Logger.Warn("用户不存在或密码错误", zap.String("phone", user.Phone))
			return 0, ErrUserLogin
		}
		utils.Logger.Error("检查用户密码失败", zap.String("phone", user.Phone), zap.Error(err))
		return 0, err
	}

	if !isValid {
		utils.Logger.Warn("密码错误", zap.String("phone", user.Phone))
		return 0, ErrUserLogin
	}

	// 获取用户ID
	userID, err := db.CheckUserPhone(user.Phone)
	if err != nil {
		utils.Logger.Error("获取用户ID失败", zap.String("phone", user.Phone), zap.Error(err))
		return 0, err
	}

	utils.Logger.Info("密码验证成功",
		zap.String("phone", user.Phone))
	return userID, nil
}

// GetAllUsers 获取所有用户列表
func GetAllUsers(page, pageSize int) ([]model.User, int64, error) {
	return db.GetAllUsers(page, pageSize)
}

// GetUsersByCondition 根据条件查询用户
func GetUsersByCondition(condition map[string]interface{}) ([]model.User, error) {
	return db.GetUsersByCondition(condition)
}

// GetUserByID 根据ID获取用户
func GetUserByID(id uint) (*model.User, error) {
	return db.GetUserByID(id)
}

// UpdateUserInfo 更新用户信息
func UpdateUserInfo(user *model.User) error {
	return db.UpdateUser(user)
}

// DeleteUserByID 删除用户
func DeleteUserByID(id uint) error {
	return db.DeleteUser(id)
}

// InitDefaultAdminUser 初始化默认管理员用户
func InitDefaultAdminUser() error {
	// 检查管理员用户是否已存在
	adminPhone := "admin"
	existingID, err := db.CheckUserPhone(adminPhone)
	if err == nil && existingID > 0 {
		utils.Logger.Info("管理员用户已存在，跳过创建")
		return nil
	}

	// 创建管理员用户
	adminUser := model.User{
		Username: "admin",
		Phone:    "admin",
		Password: "admin123",
	}

	// 加密密码
	hashed, err := bcrypt.GenerateFromPassword([]byte(adminUser.Password), bcrypt.DefaultCost)
	if err != nil {
		utils.Logger.Error("加密管理员密码失败", zap.Error(err))
		return err
	}
	adminUser.Password = string(hashed)

	// 创建用户到数据库
	userID, err := db.CreateUser(&adminUser)
	if err != nil {
		utils.Logger.Error("创建管理员用户失败", zap.Error(err))
		return err
	}

	utils.Logger.Info("管理员用户创建成功",
		zap.String("username", adminUser.Username),
		zap.Uint("id", userID))

	// 为管理员用户分配管理员角色
	adminRole, err := GetRoleByName("管理员")
	if err != nil {
		utils.Logger.Warn("获取管理员角色失败，跳过角色分配", zap.Error(err))
		return nil
	}

	if err := AssignRoleToUser(userID, adminRole.ID); err != nil {
		utils.Logger.Warn("为管理员用户分配角色失败", zap.Error(err))
		return nil
	}

	utils.Logger.Info("管理员角色分配成功",
		zap.String("username", adminUser.Username),
		zap.Uint("user_id", userID),
		zap.Uint("role_id", adminRole.ID))

	return nil
}
