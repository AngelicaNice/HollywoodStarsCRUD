package rest

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/AngelicaNice/HollywoodStarsCRUD/internal/domain"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// Auth godoc
//
//	@Summary		SignUp
//	@Description	Registration in system
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Param			input body domain.SignUpInput true "user's info"
//	@Success		201	{integer} integer 1
//	@Failure		400,404,500 {integer} integer 0
//	@Router			/auth/sign-up [post]
func (h *Handler) SignUp(c *gin.Context) {
	var user domain.SignUpInput

	decoder := json.NewDecoder(c.Request.Body)
	if err := decoder.Decode(&user); err != nil {
		log.WithFields(log.Fields{
			"handler": "SignUp",
			"issue":   "failed unmarshalling request body",
		}).Error(err)
		c.Writer.WriteHeader(http.StatusBadRequest)

		return
	}

	if err := user.Validate(); err != nil {
		log.WithFields(log.Fields{
			"handler": "SignUp",
			"issue":   "wrong params",
		}).Error(err)
		c.Writer.WriteHeader(http.StatusBadRequest)

		return
	}

	if _, err := h.usersService.Create(context.TODO(), user); err != nil {
		log.WithFields(log.Fields{
			"handler": "SignUp",
			"issue":   "internal error",
		}).Error(err)
		c.Writer.WriteHeader(http.StatusInternalServerError)

		return
	}

	c.Writer.WriteHeader(http.StatusCreated)
}

// Auth godoc
//
//	@Summary		SignIn
//	@Description	Login in system
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Param			input body domain.SignInInput true "user's info"
//	@Success 		200 {string} string "token"
//	@Failure 		400,404 {object} error
//	@Failure 		500 {object} 	 error
//	@Failure 		default {object} error
//	@Router			/auth/sign-in [post]
func (h *Handler) SignIn(c *gin.Context) {
	var user domain.SignInInput

	decoder := json.NewDecoder(c.Request.Body)
	if err := decoder.Decode(&user); err != nil {
		log.WithFields(log.Fields{
			"handler": "SignIn",
			"issue":   "failed unmarshalling request body",
		}).Error(err)
		c.Writer.WriteHeader(http.StatusBadRequest)

		return
	}

	if err := user.Validate(); err != nil {
		log.WithFields(log.Fields{
			"handler": "SignIn",
			"issue":   "wrong params",
		}).Error(err)
		c.Writer.WriteHeader(http.StatusBadRequest)

		return
	}

	accessToken, refreshToken, err := h.usersService.GetToken(context.TODO(), user)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			handleNotFoundError(c.Writer, err)

			return
		}

		log.WithFields(log.Fields{
			"handler": "SignIn",
			"issue":   "internal error",
		}).Error(err)
		c.Writer.WriteHeader(http.StatusInternalServerError)

		return
	}

	encoder := json.NewEncoder(c.Writer)
	encoder.Encode(map[string]string{
		"token": accessToken,
	})

	c.Writer.Header().Add("Set-Cookie", fmt.Sprintf("refresh-token=%s; HttpOnly", refreshToken))
	c.Writer.Header().Add("Content-Type", "application/json")
}

// Auth godoc
//
//	@Summary		Refresh
//	@Description	Refresh token
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Success 		200 {string} string "token"
//	@Failure		400,404,500  {object} error
//	@Router			/auth/refresh [get]
func (h *Handler) Refresh(c *gin.Context) {
	c.Writer.Header().Add("Content-Type", "application/json")

	cookie, err := c.Cookie("refresh-token")
	if err != nil {
		log.WithFields(log.Fields{
			"handler": "Refresh",
			"issue":   "bad request",
		}).Error(err)
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}

	accessToken, refreshToken, err := h.usersService.RefreshToken(c, cookie)
	if err != nil {
		log.WithFields(log.Fields{
			"handler": "Refresh",
			"issue":   "internal error",
		}).Error(err)
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	encoder := json.NewEncoder(c.Writer)
	encoder.Encode(map[string]string{
		"token": accessToken,
	})

	c.Header("Set-Cookie", fmt.Sprintf("refresh-token='%s'; HttpOnly", refreshToken))
}

func handleNotFoundError(w gin.ResponseWriter, err error) {
	encoder := json.NewEncoder(w)
	encoder.Encode(map[string]string{
		"error": err.Error(),
	})

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
}
