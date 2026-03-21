package service

import (
	"errors"
	"seckill_go/db"
	"seckill_go/model"
)

// 角色相关错误
var (
	ErrRoleExists         = errors.New("角色已存在")
	ErrRoleNotFound       = errors.New("角色不存在")
	ErrPermissionExists   = errors.New("权限已存在")
	ErrPermissionNotFound = errors.New("权限不存在")
	ErrUserNotFound       = errors.New("用户不存在")
)

// CreateRole 创建角色
func CreateRole(role model.Role) error {
	// 检查角色是否已存在
	var existingRole model.Role
	result := db.DB.Where("name = ?", role.Name).First(&existingRole)
	if result.RowsAffected > 0 {
		return ErrRoleExists
	}

	// 创建角色
	return db.DB.Create(&role).Error
}

// GetRoleByName 根据名称获取角色
func GetRoleByName(name string) (*model.Role, error) {
	var role model.Role
	result := db.DB.Where("name = ?", name).First(&role)
	if result.Error != nil {
		return nil, ErrRoleNotFound
	}
	return &role, nil
}

// AssignRoleToUser 为用户分配角色
func AssignRoleToUser(userID uint, roleID uint) error {
	// 检查用户是否存在
	var user model.User
	result := db.DB.First(&user, userID)
	if result.Error != nil {
		return ErrUserNotFound
	}

	// 检查角色是否存在
	var role model.Role
	result = db.DB.First(&role, roleID)
	if result.Error != nil {
		return ErrRoleNotFound
	}

	// 分配角色
	return db.DB.Model(&user).Association("Roles").Append(&role)
}

// CreatePermission 创建权限
func CreatePermission(permission model.Permission) error {
	// 检查权限是否已存在
	var existingPermission model.Permission
	result := db.DB.Where("code = ?", permission.Code).First(&existingPermission)
	if result.RowsAffected > 0 {
		return ErrPermissionExists
	}

	// 创建权限
	return db.DB.Create(&permission).Error
}

// AssignPermissionToRole 为角色分配权限
func AssignPermissionToRole(roleID uint, permissionID uint) error {
	// 检查角色是否存在
	var role model.Role
	result := db.DB.First(&role, roleID)
	if result.Error != nil {
		return ErrRoleNotFound
	}

	// 检查权限是否存在
	var permission model.Permission
	result = db.DB.First(&permission, permissionID)
	if result.Error != nil {
		return ErrPermissionNotFound
	}

	// 分配权限
	return db.DB.Model(&role).Association("Permissions").Append(&permission)
}

// CheckPermission 检查用户是否有指定权限
func CheckPermission(userID uint, permissionCode string) (bool, error) {
	// 检查数据库是否连接
	if db.DB == nil {
		// 数据库未连接，默认返回有权限
		return true, nil
	}

	// 检查用户是否存在
	var user model.User
	result := db.DB.First(&user, userID)
	if result.Error != nil {
		return false, ErrUserNotFound
	}

	// 查询用户是否有指定权限
	var count int64
	db.DB.Model(&user).Joins("JOIN user_roles ON users.id = user_roles.user_id").Joins("JOIN role_permissions ON user_roles.role_id = role_permissions.role_id").Joins("JOIN permissions ON role_permissions.permission_id = permissions.id").Where("users.id = ? AND permissions.code = ?", userID, permissionCode).Count(&count)

	return count > 0, nil
}

// InitDefaultRoles 初始化默认角色和权限
func InitDefaultRoles() error {
	// 创建默认权限
	permissions := []model.Permission{
		{Name: "查看商品", Code: "product:view", Description: "查看商品列表和详情"},
		{Name: "创建商品", Code: "product:create", Description: "创建新商品"},
		{Name: "更新商品", Code: "product:update", Description: "更新商品信息"},
		{Name: "删除商品", Code: "product:delete", Description: "删除商品"},
		{Name: "秒杀商品", Code: "product:seckill", Description: "参与商品秒杀"},
		{Name: "查看活动", Code: "event:view", Description: "查看活动列表和详情"},
		{Name: "创建活动", Code: "event:create", Description: "创建新活动"},
		{Name: "更新活动", Code: "event:update", Description: "更新活动信息"},
		{Name: "删除活动", Code: "event:delete", Description: "删除活动"},
		{Name: "开始活动", Code: "event:start", Description: "开始活动"},
		{Name: "停止活动", Code: "event:stop", Description: "停止活动"},
		{Name: "查看订单", Code: "order:view", Description: "查看订单列表和详情"},
	}

	for _, permission := range permissions {
		if err := CreatePermission(permission); err != nil && err != ErrPermissionExists {
			return err
		}
	}

	// 创建默认角色
	roles := []model.Role{
		{Name: "普通用户", Description: "普通用户角色"},
		{Name: "管理员", Description: "管理员角色"},
	}

	for _, role := range roles {
		if err := CreateRole(role); err != nil && err != ErrRoleExists {
			return err
		}
	}

	// 为普通用户分配权限
	userRole, err := GetRoleByName("普通用户")
	if err != nil {
		return err
	}

	userPermissions := []string{"product:view", "product:seckill", "event:view", "order:view"}
	for _, code := range userPermissions {
		var permission model.Permission
		db.DB.Where("code = ?", code).First(&permission)
		if permission.ID > 0 {
			AssignPermissionToRole(userRole.ID, permission.ID)
		}
	}

	// 为管理员分配所有权限
	adminRole, err := GetRoleByName("管理员")
	if err != nil {
		return err
	}

	var allPermissions []model.Permission
	db.DB.Find(&allPermissions)
	for _, permission := range allPermissions {
		AssignPermissionToRole(adminRole.ID, permission.ID)
	}

	return nil
}
