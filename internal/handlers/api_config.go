package handlers

import (
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/aAmer0neee/rest-mindspace-blog/internal/auth"
	"github.com/aAmer0neee/rest-mindspace-blog/internal/service"
)

func Configurehandlers(r *gin.Engine, s *service.Service, j *auth.JWTService) {

	// Static для обработки фронтэнда
	{
		r.SetFuncMap(template.FuncMap{
			"safeHTML": func(s string) template.HTML {
				return template.HTML(s)
			},
		})

		r.LoadHTMLGlob("templates/*.html")
		r.Static("/static", "./static")
	}
	//глобальные middleware
	{
		r.Use(AssignId())

		r.Use(RequestLimiter())

	}
	// отдельные группы
	{
		admin := r.Group("/admin")
		admin.Use(JWTmiddleware(j))

		admin.GET("/", func(ctx *gin.Context) {
			ctx.HTML(http.StatusOK, "post_form.html", nil)
		})
		admin.POST("/", func(ctx *gin.Context) {
			handlePostArticle(ctx, s)
		})
	}

	{
		auth := r.Group("auth")

		auth.GET("/", func(ctx *gin.Context) {
			ctx.HTML(http.StatusOK, "auth.html", nil)
		})
		auth.POST("/", func(ctx *gin.Context) {
			handleAuth(ctx, s, j)
		})
		auth.GET("/register", func(ctx *gin.Context) {
			ctx.HTML(http.StatusOK, "register.html", nil)
		})
		auth.POST("/register", func(ctx *gin.Context) {
			handleRegister(ctx, s)
		})
	}

	{
		feed := r.Group("/")

		feed.GET("/", func(ctx *gin.Context) {
			ctx.HTML(http.StatusOK, "feed.html", nil)

		})

		feed.GET("/data", func(ctx *gin.Context) {
			handleGetFeed(ctx, s)
		})

		feed.GET("/article", func(ctx *gin.Context) {
			handleGetArticle(ctx, s)
		})
	}

}
