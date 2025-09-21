package controllers

import (
	"app_api/middleware"
	"app_api/models"
	"app_api/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// LoginUser 用户登录接口
func LoginUser(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ParamError(c)
		return
	}

	// 邮箱格式校验
	if !utils.ValidateEmailFormat(req.Email) {
		utils.Error(c, "邮箱格式不合法")
		return
	}

	// 密码长度校验
	if len(req.Password) < 6 {
		utils.Error(c, "密码必须大于6位")
		return
	}

	// 通过邮箱获取用户
	user, err := models.GetUserByEmail(req.Email)
	if err != nil {
		utils.Error(c, "用户不存在")
		return
	}

	// 验证密码
	if utils.HashPassword(req.Password, user.Salt) != user.Password {
		utils.Error(c, "密码错误")
		return
	}

	// 生成JWT令牌
	token, err := middleware.GenerateJWT(user.UID)
	if err != nil {
		utils.Error(c, "登录失败")
		return
	}

	// 更新最后登录信息
	user.LastLoginAt = time.Now()
	user.LastLoginIP = c.ClientIP()
	user.Token = token
	if err := models.UpdateUser(&user); err == nil {
		utils.Success(c, "登录成功", gin.H{
			"uid":      user.UID,
			"email":    user.Email,
			"token":    token,
			"avatar":   user.Avatar,
			"username": user.Username,
			"phone":    user.Phone,
			"role":     user.Role,
		})
	} else {
		utils.Error(c, "登录失败")
		return
	}
}

// CreateUser 创建用户接口
func CreateUser(c *gin.Context) {
	var request models.User
	// 绑定JSON数据到结构体
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.ParamError(c)
		return
	}

	// 邮箱格式校验
	if !utils.ValidateEmailFormat(request.Email) {
		utils.Error(c, "邮箱格式不合法")
		return
	}

	// 密码长度校验
	if len(request.Password) < 6 {
		utils.Error(c, "密码必须大于6位")
		return
	}

	// 生成盐值
	salt, err := utils.GenerateRandomString(32)
	if err != nil {
		utils.Error(c, "系统错误")
		return
	}

	token, err := middleware.GenerateJWT(request.UID)
	if err != nil {
		utils.Error(c, "系统错误")
		return
	}

	// 生成UID
	request.UID = uuid.New().String()
	request.CreatedAt = time.Now()
	request.UpdatedAt = time.Now()
	request.LastLoginAt = time.Now()
	request.Salt = salt
	request.LastLoginIP = c.ClientIP()
	request.Token = token
	request.Password = utils.HashPassword(request.Password, request.Salt) // 修复变量名称错误
	// 保存用户到数据库
	if err := models.CreateUser(&request); err != nil {
		utils.Error(c, err.Error())
		return
	}

	// 返回成功响应
	utils.Success(c, "注册成功", gin.H{
		"email":     request.Email,
		"username":  request.Username,
		"phone":     request.Phone,
		"uid":       request.UID,
		"createdAt": request.CreatedAt.Format(time.RFC3339), // 格式化时间
		"avatar":    request.Avatar,                         // 返回头像地址
		"role":      request.Role,
		"token":     token, // 返回用户角色
	})
}

// GetUserInfo 查询用户信息接口
func GetUserInfo(c *gin.Context) {
	// 从上下文中获取UID
	uid, exists := c.Get("uid")
	if !exists {
		utils.Error(c, "未授权访问")
		return
	}

	// 查询用户信息
	user, err := models.GetUserByUID(uid.(string))
	if err != nil {
		utils.Error(c, "用户不存在")
		return
	}

	// 返回用户信息（过滤敏感字段）
	utils.Success(c, "查询成功", gin.H{
		"uid":       user.UID,
		"username":  user.Username,
		"email":     user.Email,
		"avatar":    user.Avatar,
		"phone":     user.Phone,
		"role":      user.Role,
		"lastLogin": user.LastLoginAt.Format(time.RFC3339),
		"createdAt": user.CreatedAt.Format(time.RFC3339),
	})
}
