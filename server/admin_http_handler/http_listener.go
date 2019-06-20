package admin_http_handler

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/elazarl/go-bindata-assetfs"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"

	"github.com/x-insane/ngrokex/server/assets"
	"github.com/x-insane/ngrokex/server/config"
)

type binaryFileSystem struct {
	fs http.FileSystem
}

func (b *binaryFileSystem) Open(name string) (http.File, error) {
	return b.fs.Open(name)
}

func (b *binaryFileSystem) Exists(prefix string, filepath string) bool {
	if p := strings.TrimPrefix(filepath, prefix); len(p) < len(filepath) {
		if _, err := b.fs.Open(p); err != nil {
			return false
		}
		return true
	}
	return false
}

func BinaryFileSystem(root string) *binaryFileSystem {
	fs := &assetfs.AssetFS{
		Asset: assets.Asset,
		AssetDir: assets.AssetDir,
		AssetInfo: assets.AssetInfo,
		Prefix: root,
	}
	return &binaryFileSystem{
		fs,
	}
}

func AdminHttpListener(conf config.AdminConf) {
	r := gin.Default()
	r.Use(static.Serve("/", BinaryFileSystem("assets/server/static")))

	store := cookie.NewStore([]byte("Drrmpm5zB8XOBGFgngdKbRxJwiNbHfNFyfIhNC2YmkQhNCekuXSZS13fpvTqoJ4B"))
	r.Use(sessions.Sessions("default", store))
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.Group("api").
		POST("login", Login).
		POST("logout", Logout).
		POST("register", Register).
		POST("user/info", GetUser).
		POST("ports/list", GetAllPortAuth).
		POST("sub_domains/list", GetAllSubDomainAuth).
		POST("users/list", ListUsers).
		POST("ports/auth", AuthPortToUser).
		POST("sub_domains/auth", AuthSubDomainToUser)

	if conf.HttpPort == 0 {
		conf.HttpPort = 18080
	}
	err := r.Run(fmt.Sprintf(":%d", conf.HttpPort))
	if err != nil {
		panic(err)
	}
}
