package admin_http_handler

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/x-insane/ngrokex/server/db"
)

func Login(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.BindJSON(&req); CommonError(c, err, "login") {
		return
	}

	// login
	user, err := db.LoginGetUser(req.Username, req.Password)
	if err != nil {
		ApiResult(c, 403, "账号或密码错误")
		return
	}

	// save to session
	session := sessions.Default(c)
	session.Set("user", user.UserId)
	_ = session.Save()

	HttpSuccess(c)
}

func Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Delete("user")
	_ = session.Save()
	HttpSuccess(c)
}

func Register(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.BindJSON(&req); CommonError(c, err, "login") {
		return
	}

	user := db.User{
		Username: req.Username,
		Password: req.Password,
	}
	err := db.CreateUser(&user)
	if CommonError(c, err, "register") {
		return
	}

	// save to session
	session := sessions.Default(c)
	session.Set("user", user.UserId)
	_ = session.Save()

	c.JSON(200, gin.H{
		"error": 0,
		"user": user,
	})
}

func GetAllPortAuth(c *gin.Context) {
	userId := CheckLogin(c)
	if userId == 0 {
		return
	}
	ports, err := db.GetAllPortAuth(userId)
	if CommonError(c, err, "get_all_ports_auth") {
		return
	}
	c.JSON(200, gin.H{
		"error": 0,
		"ports": ports,
	})
}

func GetAllSubDomainAuth(c *gin.Context) {
	userId := CheckLogin(c)
	if userId == 0 {
		return
	}
	subDomains, err := db.GetAllSubDomainAuth(userId)
	if CommonError(c, err, "get_all_sub_domain_auth") {
		return
	}
	c.JSON(200, gin.H{
		"error": 0,
		"sub_domains": subDomains,
	})
}

func GetUser(c *gin.Context) {
	userId := CheckLogin(c)
	if userId == 0 {
		return
	}
	user, err := db.GetUserById(userId)
	if CommonError(c, err, "get_user") {
		return
	}
	user.Password = ""
	c.JSON(200, gin.H{
		"error": 0,
		"user": user,
	})
}

func ListUsers(c *gin.Context) {
	_, err := CheckAdminRole(c)
	if err != nil {
		return
	}
	users, err := db.ListUsers()
	if CommonError(c, err, "list_users") {
		return
	}
	c.JSON(200, gin.H{
		"error": 0,
		"users": users,
	})
}

func AuthPortToUser(c *gin.Context) {
	_, err := CheckAdminRole(c)
	if err != nil {
		return
	}

	var req struct {
		UserId int64 `json:"user_id"`
		Port   int64 `json:"port"`
	}
	if err := c.BindJSON(&req); CommonError(c, err, "auth_port_to_user") {
		return
	}

	err = db.AuthPortToUser(req.Port, req.UserId)
	if CommonError(c, err, "auth_port_to_user") {
		return
	}
	HttpSuccess(c)
}

func AuthSubDomainToUser(c *gin.Context) {
	_, err := CheckAdminRole(c)
	if err != nil {
		return
	}

	var req struct {
		UserId    int64  `json:"user_id"`
		SubDomain string `json:"sub_domain"`
	}
	if err := c.BindJSON(&req); CommonError(c, err, "auth_sub_domain_to_user") {
		return
	}

	if req.SubDomain == "" {
		ApiResult(c, 400, "请填写子域名")
		return
	}

	err = db.AuthSubDomainToUser(req.SubDomain, req.UserId)
	if CommonError(c, err, "auth_sub_domain_to_user") {
		return
	}
	HttpSuccess(c)
}
