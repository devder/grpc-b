package api

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/devder/grpc-b/token"
	"github.com/gin-gonic/gin"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "bearer"
	authorizationPayloadKey = "authorization_payload"
)

func authMiddleware(tokenMaker token.Maker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		key := ctx.GetHeader(authorizationHeaderKey)
		if len(key) == 0 {
			err := fmt.Errorf("authorization header is not provided")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		fields := strings.Fields(key)
		if len(fields) < 2 {
			err := fmt.Errorf("authorization header is not provided")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		authType, token := fields[0], fields[1]

		authorizationType := strings.ToLower(authType)
		if authorizationType != authorizationTypeBearer {
			err := fmt.Errorf("authorization type %s is not supported", authorizationType)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		payload, err := tokenMaker.VerifyToken(token)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		ctx.Set(authorizationPayloadKey, payload)
		ctx.Next()

	}
}
