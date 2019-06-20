package admin_http_handler

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/x-insane/ngrokex/log"
	"github.com/x-insane/ngrokex/server/db"
)

func HttpSuccess(c *gin.Context) {
	ApiResult(c, 0, "success")
}

func ApiResult(c *gin.Context, code int64, msg string) {
	c.JSON(200, gin.H{
		"error": code,
		"msg": msg,
	})
}

// returns true if error
func CommonError(c *gin.Context, err error, scope string) bool {
	if err != nil {
		_ = log.Error("[%s] http error: %+v", scope, err)
		c.JSON(200, gin.H{
			"error": 500,
			"msg": fmt.Sprintf("http error: %+v", err),
		})
		return true
	}
	return false
}

func CheckLogin(c *gin.Context) int64 {
	session := sessions.Default(c)
	var userId int64 = 0
	if v := session.Get("user"); v != nil {
		userId = v.(int64)
	}
	if userId == 0 {
		ApiResult(c, 403, "请先登录")
	}
	return userId
}

func CheckAdminRole(c *gin.Context) (*db.User, error) {
	session := sessions.Default(c)
	var userId int64 = 0
	if v := session.Get("user"); v != nil {
		userId = v.(int64)
	}
	if userId == 0 {
		ApiResult(c, 403, "请先登录")
		return nil, fmt.Errorf("not login")
	}
	user, err := db.GetUserById(userId)
	if err != nil {
		CommonError(c, err, "check_admin_role")
		return user, err
	}
	if user.Role != "admin" {
		ApiResult(c, 403, "权限不足")
		return user, fmt.Errorf("not admin")
	}
	return user, nil
}
