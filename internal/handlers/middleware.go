package handlers

import (
	"log"
	"net/http"
	_ "sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ulule/limiter/v3"
	middleware "github.com/ulule/limiter/v3/drivers/middleware/gin"
	"github.com/ulule/limiter/v3/drivers/store/memory"

	"github.com/aAmer0neee/rest-mindspace-blog/internal/auth"
	"github.com/aAmer0neee/rest-mindspace-blog/internal/service"
)

func JWTmiddleware(j *auth.JWTService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		cookie, err := ctx.Request.Cookie("JWT")
		if err != nil {
			authError(ctx)
			return
		}

		token, err := j.ValidateJWT(cookie.Value)
		if err != nil || !token.Valid {
			authError(ctx)
			return
		}

		exp := j.GetExpTime(token)
		if time.Now().After(time.Unix(int64(exp), 0)) {
			authError(ctx)
			return
		}

		username := j.GetUsername(token)
		ctx.Set("username", username)

		ctx.Next()
	}
}

func AssignId() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		cookie, err := ctx.Request.Cookie("id")
		if err != nil || cookie.Value == "" {
			id := service.GenerateUUID().String()
			ctx.SetCookie("id", id, int(time.Hour)*1, "/", "localhost", false, false)
		}

	}
}

func RequestLimiter() gin.HandlerFunc {
	rate, err := limiter.NewRateFromFormatted("100-S")
	if err != nil {
		log.Fatal("error setup limiter middleware", err.Error())
	}

	store := memory.NewStore()

	return middleware.NewMiddleware(limiter.New(store, rate),
		middleware.WithLimitReachedHandler(func(ctx *gin.Context) {
			ctx.Status(http.StatusTooManyRequests)
			ctx.Abort()
		}))
}

/*
var (
	maxRequests = 100
	totalRequests int
	mu sync.Mutex
	resetTime time.Time
)

 func RequestLimiter() gin.HandlerFunc{
	return func(ctx *gin.Context) {
		mu.Lock()
		defer mu.Unlock()

		if time.Since(resetTime) >= time.Second {
			totalRequests = 0
			resetTime = time.Now()
		}

		if totalRequests >= maxRequests {
			ctx.Status(http.StatusTooManyRequests)
			ctx.Abort()
			return
		}
		totalRequests++
		ctx.Next()
	}
}
*/
