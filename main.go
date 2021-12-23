package main

import (
	"fmt"
	"io/ioutil"
	"path"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

// HTTP
func cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("server", "CluckCDN 1.0")
		c.Writer.Header().Set("Content-Type", "text/html;charset=utf-8")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

func main() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(cors())
	router.Use(gin.Recovery())
	store := cookie.NewStore([]byte("clucknetwork"))
	router.Use(sessions.Sessions("cluck_token", store))

	// Index
	router.GET("/", func(c *gin.Context) {
		c.Redirect(302, "/page/index")
	})
	// 404
	router.NoRoute(func(context *gin.Context) {
		context.Writer.Header().Set("Content-Type", "application/json;charset=utf-8")
		context.String(404, `{"status":"404"}`)
	})
	// Status
	router.GET("/status", func(context *gin.Context) {
		context.Writer.Header().Set("Content-Type", "application/json;charset=utf-8")
		context.String(200, `{"status":"ok"}`)
	})
	// Static
	router.GET("/static/:filename", func(context *gin.Context) {
		context.Writer.Header().Set("Content-Type", "text/plain;charset=utf-8")
		filename := context.Param("filename")
		filetext, err := ioutil.ReadFile("static/" + filename)
		if err != nil {
			context.String(404, `{"status":"404"}`)
		}
		context.String(200, string(filetext))
	})
	// Script (CSS/JS)
	router.GET("/script/:filename", func(context *gin.Context) {
		filename := context.Param("filename")
		fileExt := path.Ext(filename)
		var filetext []byte
		if fileExt == ".js" || fileExt == ".css" {
			if fileExt == ".js" {
				context.Writer.Header().Set("Content-Type", "text/javascript;charset=utf-8")
			} else if fileExt == ".css" {
				context.Writer.Header().Set("Content-Type", "text/css;charset=utf-8")
			}
			var err error
			filetext, err = ioutil.ReadFile("script/" + filename)
			if err != nil {
				context.String(404, `{"status":"404"}`)
			}
		} else {
			context.String(403, `{"status":"403"}`)
		}
		context.String(200, string(filetext))
	})
	// Html
	router.GET("/page/:name", func(c *gin.Context) {
		session := sessions.Default(c)
		name := c.Param("name")
		if session.Get("login") == "ok" || name == "login" {
			file, err := ioutil.ReadFile("html/" + name + ".html")
			if err != nil {
				c.Writer.Header().Set("Content-Type", "application/json;charset=utf-8")
				c.String(404, `{"status":"404"}`)
				return
			}
			c.String(200, string(file))
		} else {
			c.Redirect(302, "/page/login")
		}
	})

	router.GET("/install.sh", getInstall) // Install Shell
	router.GET("/config.yaml", getConfig) // Get Config
	router.GET("/vhost.conf", getVhost)   // Get Vhost Config
	router.GET("tls.zip", getTls)         // Get TLS

	// API
	router.POST("/api/login", loginApi)                    // Login
	router.POST("/api/postNodeApi", postNodeApi)           // POST 節點 API
	router.POST("/api/nodeList", nodeList)                 // 節點列表
	router.POST("/api/updateNodeStatus", updateNodeStatus) // 節點狀態
	router.POST("/api/setAllStatus", setAllStatus)         // 設定所有節點狀態
	router.POST("/api/getHttpToken", getHttpToken)         // Get Token

	router.POST("/api/websiteList", websiteList) // 網站列表

	router.POST("/api/delNode", delNode)       // 刪節點
	router.POST("/api/addNode", addNode)       // 加節點
	router.POST("/api/addWebSite", addWebSite) // 加網站
	router.POST("/api/delWebSite", delWebSite) // 刪網站

	// Run
	fmt.Println("CluckCDN Ctrl Starting ...")
	router.Run(":80")
}
