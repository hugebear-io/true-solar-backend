package middleware

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/hugebear-io/true-solar-backend/pkg/config"
	"github.com/hugebear-io/true-solar-backend/pkg/deliver"
)

func Authentication() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// validate header
		requestToken := ctx.Request.Header.Get("Authorization")
		if !strings.HasPrefix(requestToken, "Bearer") {
			deliver.ResponseUnauthorized(ctx)
			return
		}
		requestToken = strings.TrimPrefix(requestToken, "Bearer")
		token := strings.TrimSpace(requestToken)
		if len(token) == 0 {
			deliver.ResponseUnauthorized(ctx)
			return
		}

		// extract jwt token
		tokenClaims := jwt.MapClaims{}
		jwtToken, err := jwt.ParseWithClaims(token, tokenClaims,
			func(token *jwt.Token) (interface{}, error) {
				return []byte(config.Config.API.SecretKey), nil
			})
		if err != nil {
			deliver.ResponseUnauthorized(ctx)
			return
		}
		claims, exits := jwtToken.Claims.(jwt.MapClaims)
		if !exits || !jwtToken.Valid {
			deliver.ResponseUnauthorized(ctx)
			return
		}

		// validate expired time
		if _, exits := claims["expired_time"]; !exits {
			deliver.ResponseUnauthorized(ctx)
			return
		}
		expiredTimeString := claims["expired_time"].(string)
		if len(expiredTimeString) == 0 {
			deliver.ResponseUnauthorized(ctx)
			return
		}
		expiredTime, err := time.Parse(time.RFC3339Nano, expiredTimeString)
		if err != nil {
			deliver.ResponseUnauthorized(ctx)
			return
		}
		if expiredTime.UTC().Before(time.Now().UTC()) {
			deliver.ResponseUnauthorized(ctx)
			return
		}

		// validate userID
		userID, exits := claims["user_id"]
		if !exits {
			deliver.ResponseUnauthorized(ctx)
			return
		}
		userID, err = strconv.Atoi(fmt.Sprintf("%.f", userID))
		if err != nil {
			deliver.ResponseUnauthorized(ctx)
			return
		}
		ctx.Set("user_id", userID)
	}
}
