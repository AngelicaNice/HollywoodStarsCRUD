package rest

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

//func Logger() gin.HandlerFunc {
//	return func(c *gin.Context) {
//		log.WithFields(log.Fields{
//			"method": c.Request.Method,
//			"uri":    c.Request.URL,
//		}).Info()
//	}
//}

func authMiddleware(h *Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := getTokenFromRequest(c.Request)
		if err != nil {
			log.WithField("authMiddleware", "getting token").Error(err)
			c.Writer.WriteHeader(http.StatusUnauthorized)

			return
		}

		userId, err := h.usersService.ParseToken(c.Request.Context(), token)
		if err != nil {
			log.WithField("authMiddleware", "parsing token").Error(err)
			c.Writer.WriteHeader(http.StatusUnauthorized)

			return
		}
		//ctx := context.WithValue(r.Context(), ctxUserID, userId)
		//r = r.WithContext(ctx)
		ctx := context.WithValue(c.Request.Context(), "id", userId)
		c.Request = c.Request.WithContext(ctx)
	}
}

func getTokenFromRequest(r *http.Request) (string, error) {
	header := r.Header.Get("Authorization")

	if header == "" {
		return "", errors.New("empty authorization header")
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		return "", errors.New("invalid authorization header")
	}

	if len(headerParts[1]) == 0 {
		return "", errors.New("empty token")
	}

	return headerParts[1], nil
}
