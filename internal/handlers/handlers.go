package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/aAmer0neee/rest-mindspace-blog/internal/auth"
	"github.com/aAmer0neee/rest-mindspace-blog/internal/domain"
	"github.com/aAmer0neee/rest-mindspace-blog/internal/service"
)

func handlePostArticle(ctx *gin.Context, s *service.Service) {

	article := domain.ArticleEntity{}

	if err := ctx.ShouldBindBodyWithJSON(&article); err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}
	username, _ := ctx.Get("username")

	article.Author = username.(string)
	if err := s.NewPost(article); err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}

	ctx.Status(http.StatusOK)
}

func handleGetFeed(ctx *gin.Context, s *service.Service) {
	limit := 3
	page, err := strconv.Atoi(ctx.Query("page"))
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	feed, total, err := s.GetFeed(page, limit)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"total": total,
		"posts": feed,
	})
}

func handleGetArticle(ctx *gin.Context, s *service.Service) {

	id, err := strconv.Atoi(ctx.Query("id"))
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	article, err := s.GetArticle(id)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	println(article.Preview)

	ctx.HTML(http.StatusOK, "article.html", gin.H{
		"Title":   article.Title,
		"Preview": article.Preview,
		"Date":    article.CreatedAt,
		"Content": article.Content,
		"Author":  article.Author,
	})
}

func handleRegister(ctx *gin.Context, s *service.Service) {
	user := domain.UserEntity{}
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}

	cookie, _ := ctx.Request.Cookie("id")
	if cookie != nil {
		user.ID = cookie.Value
	}

	if err := s.RegisterUser(user); err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}
	ctx.Status(http.StatusOK)

}

func handleAuth(ctx *gin.Context, s *service.Service, j *auth.JWTService) {
	user := domain.UserEntity{}
	if err := ctx.ShouldBindBodyWithJSON(&user); err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}

	id, err := s.AuthUser(user)

	if err != nil {
		ctx.Status(http.StatusUnauthorized)
		return
	}

	token, err := j.GenerateJWT(user.Username, id)
	if err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}

	ctx.SetCookie("JWT", token, 0, "/", "localhost", false, false)
	ctx.SetCookie("id", id, 0, "/", "localhost", false, false)

	next := ctx.Query("next")
	if next == "" {
		next = "/admin"
	}

	ctx.Redirect(http.StatusFound, next)
}

func authError(ctx *gin.Context) {
	ctx.Redirect(http.StatusFound, "/auth")
	ctx.Abort()
}
