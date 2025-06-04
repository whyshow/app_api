package models

import (
	"app_api/database"
	"fmt"
	"time"
)

// User 用户模型
type User struct {
	UID         string     `gorm:"primaryKey;type:varchar(36);uniqueIndex;<-:create"`         // 仅允许创建时写入
	Email       string     `gorm:"size:255;not null;<-:create" json:"email"`                  // 同UID保护规则  邮箱字段
	Phone       string     `gorm:"default:null;size:20" json:"phone"`                         // 手机号字段
	Username    string     `gorm:"default:null;size:100" json:"username"`                     // 用户名字段
	Password    string     `gorm:"size:255;not null;comment:密码哈希" json:"password"`            // 密码字段
	Token       string     `gorm:"size:255" json:"token"`                                     // token字段
	Salt        string     `gorm:"size:32;comment:密码盐" json:"salt"`                           // 盐值字段
	Role        string     `gorm:"index;size:50;default:'user'" json:"role"`                  // 角色字段
	Gender      int        `gorm:"type:tinyint;default:0;comment:0-未知 1-男 2-女" json:"gender"` // 性别字段
	Country     string     `gorm:"default:null;size:100" json:"country"`                      // 国家字段
	LastLoginAt time.Time  `gorm:"default:null;index" json:"lastLoginAt"`                     // 最后登录时间字段
	LastLoginIP string     `gorm:"default:null;size:45" json:"lastLoginIp"`                   // 最后登录IP字段
	LoginDevice string     `gorm:"default:null;size:100" json:"loginDevice"`                  // 登录设备字段
	Avatar      string     `gorm:"default:null;type:text" json:"avatar"`                      // 头像字段
	Birthday    *time.Time `gorm:"type:datetime" json:"birthday,omitempty"`                   // 生日字段
	Address     string     `gorm:"default:null;size:255" json:"address"`                      // 地址字段
	CreatedAt   time.Time  `gorm:"default:null;index" json:"createdAt"`                       // 创建时间字段
	UpdatedAt   time.Time  `json:"default:null;updatedAt"`                                    // 更新时间字段
}

// LoginRequest 登录请求体
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// GetUserByEmail 根据邮箱获取用户
func GetUserByEmail(email string) (User, error) {
	var user User
	err := database.GetDB().Where("email = ?", email).First(&user).Error
	return user, err
}

// GetUserByUID 根据UID获取用户
func GetUserByUID(uid string) (User, error) {
	var user User
	err := database.GetDB().Where("uid = ?", uid).First(&user).Error
	return user, err
}

// CreateUser 创建用户时增加校验
func CreateUser(user *User) error {
	// 检查邮箱是否已存在
	if exists, _ := CheckEmailExists(user.Email); exists {
		return fmt.Errorf("邮箱已被注册")
	}

	return database.GetDB().Create(user).Error
}

// CheckEmailExists 检查邮箱是否存在
func CheckEmailExists(email string) (bool, error) {
	var count int64
	err := database.GetDB().Model(&User{}).Where("email = ?", email).Count(&count).Error
	return count > 0, err
}

func UpdateUser(user *User) error {
	// 创建更新白名单，排除不允许修改的字段
	updates := map[string]interface{}{
		"token":         user.Token,
		"username":      user.Username,
		"password":      user.Password,
		"phone":         user.Phone,
		"role":          user.Role,
		"salt":          user.Salt,
		"gender":        user.Gender,
		"country":       user.Country,
		"avatar":        user.Avatar,
		"birthday":      user.Birthday,
		"address":       user.Address,
		"last_login_at": user.LastLoginAt,
		"last_login_ip": user.LastLoginIP,
		"login_device":  user.LoginDevice,
		"updated_at":    time.Now(),
	}

	return database.GetDB().Model(user).
		Select("*").          // 重置选择所有字段
		Omit("uid", "email"). // 明确排除保护字段
		Updates(updates).Error
}
