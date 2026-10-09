package controller

import (
	"agent/cmd/hazmatAgent/model"
	"agent/util"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RegisterReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Register(c *gin.Context) {

	/*
		#### 加密等待移动到前端实现
	*/

	var req RegisterReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "参数错误"})
		return
	}

	hashPwd, err := util.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "加密失败"})
		return
	}

	user := model.User{
		Username: req.Username,
		Password: hashPwd,
	}

	err = util.GetMySQLClient().GormDB.Create(&user).Error
	if err != nil {
		slog.Error("注册失败", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"msg": "注册失败，请更换用户名"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"msg": "注册成功", "user_id": user.ID, "username": user.Username})
}

func Login(c *gin.Context) {
	var req LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "参数错误"})
		return
	}

	var user model.User
	err := util.GetMySQLClient().GormDB.Where("username = ?", req.Username).First(&user).Error
	if err != nil {
		slog.Error("登录失败", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"msg": "登录失败，请检查用户名或密码"})
		return
	}

	if !util.CheckPassword(user.Password, req.Password) {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "登录失败，请检查用户名或密码"})
		return
	}

	token, err := util.GenToken(user.ID)
	if err != nil {
		slog.Error("生成token失败", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "生成token失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token, "user": user})
}


