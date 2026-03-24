package db

import (
	"errors"
	"seckill_go/model"
	"seckill_go/utils"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// CreateUser 新增用户到数据库
// 参数：user - 包含用户信息的结构体（需包含加密后的密码）
// 返回：创建成功的用户ID和可能的错误
func CreateUser(user *model.User) (uint, error) {
	if user == nil {
		return 0, errors.New("用户信息不能为空")
	}

	// 执行插入操作
	result := DB.Create(user)
	if result.Error != nil {
		// 捕获唯一索引冲突（如邮箱已存在）
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return 0, errors.New("索引冲突")
		}
		return 0, result.Error
	}

	// 返回自增ID
	return user.ID, nil
}

// GetUserByID 根据ID查询用户
// 参数：id - 用户ID
// 返回：用户信息结构体和可能的错误
func GetUserByID(id uint) (*model.User, error) {
	if id == 0 {
		return nil, errors.New("用户ID不能为空")
	}

	var user model.User
	result := DB.Where("id = ?", id).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, result.Error
	}

	return &user, nil
}

// CheckUserPhone 检查用户手机号是否被注册
// 参数：phone - 要检查的手机号
// 返回：若查询到则返回主键id，若没查询到则返回0
func CheckUserPhone(phone string) (uint, error) {
	if phone == "" {
		return 0, errors.New("手机号不能为空")
	}

	var userId uint
	// 只查询id字段，提高效率
	result := DB.Model(&model.User{}).Where("phone = ?", phone).Pluck("id", &userId)

	if result.Error != nil {
		return 0, result.Error
	}

	// 如果没有查询到记录，GORM不会返回错误，此时userId会保持默认值0
	return userId, nil
}

// CheckUserPassword 检查用户密码是否正确
// 参数：phone - 要检查的手机号
// 返回：是否存在的布尔值和可能的错误
func CheckUserPassword(phone string, password string) (bool, error) {
	if phone == "" {
		return false, errors.New("手机号不能为空")
	}

	var existingUser model.User
	err := DB.Where("phone = ?", phone).First(&existingUser).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		utils.Logger.Warn("用户不存在",
			zap.String("phone", phone))
		return false, errors.New("用户不存在")
	} else if err != nil {
		utils.Logger.Error("用户检查过程中出错",
			zap.Error(err),
			zap.String("phone", phone))
		return false, err
	}
	// 获取用户密码并校验
	if err := bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(password)); err != nil {
		utils.Logger.Warn("密码错误",
			zap.String("phone", phone))
		return false, errors.New("密码错误")
	}

	return true, nil
}

// GetAllUsers 获取所有用户列表
// 参数：page - 页码（从1开始），pageSize - 每页大小
// 返回：用户列表和总记录数
func GetAllUsers(page, pageSize int) ([]model.User, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	var users []model.User
	var total int64

	// 获取总记录数
	err := DB.Model(&model.User{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	err = DB.Offset(offset).Limit(pageSize).Find(&users).Error
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// GetUsersByCondition 根据条件查询用户
// 参数：condition - 查询条件（如用户名、手机号等）
// 返回：符合条件的用户列表
func GetUsersByCondition(condition map[string]interface{}) ([]model.User, error) {
	var users []model.User

	query := DB.Model(&model.User{})

	// 构建查询条件
	if username, ok := condition["username"]; ok && username != "" {
		query = query.Where("username LIKE ?", "%"+username.(string)+"%")
	}

	if phone, ok := condition["phone"]; ok && phone != "" {
		query = query.Where("phone LIKE ?", "%"+phone.(string)+"%")
	}

	// 执行查询
	err := query.Find(&users).Error
	if err != nil {
		return nil, err
	}

	return users, nil
}

// DeleteUser 删除用户（软删除）
// 参数：id - 用户ID
// 返回：删除操作的错误
func DeleteUser(id uint) error {
	if id == 0 {
		return errors.New("用户ID不能为空")
	}

	// 执行软删除
	result := DB.Delete(&model.User{}, id)
	return result.Error
}

// UpdateUser 更新用户信息
// 参数：user - 包含更新信息的用户结构体（需包含ID）
// 返回：更新操作的错误
func UpdateUser(user *model.User) error {
	if user == nil || user.ID == 0 {
		return errors.New("用户ID不能为空")
	}

	// 只更新非零值字段（避免覆盖未修改的字段）
	result := DB.Model(&model.User{}).Where("id = ?", user.ID).Updates(user)
	return result.Error
}
