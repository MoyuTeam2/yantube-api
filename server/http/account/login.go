package account

import (
	"api/db"
	"api/server/http/common"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type UserLoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Login(c *gin.Context) {
	var req UserLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error().Err(err).Msg("login request bind json failed")
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(http.StatusBadRequest, "bad request"))
		return
	}

	user, err := db.Get().GetUserByUsername(req.Username)
	if err != nil {
		if errors.Is(err, db.ErrUserNotFound) {
			c.JSON(http.StatusBadRequest, common.NewErrorResponse(http.StatusBadRequest, "user not exists"))
			return
		}
		log.Error().Err(err).Msg("get user by username failed")
		c.JSON(http.StatusInternalServerError, common.NewErrorResponse(http.StatusInternalServerError, "server get user failed"))
		return
	}

	if valid, err := user.ValidatePassword(req.Password); err != nil {
		log.Error().Err(err).Msg("validate password failed")
		c.JSON(http.StatusInternalServerError, common.NewErrorResponse(http.StatusInternalServerError, "server can't validate password"))
		return
	} else if !valid {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(http.StatusBadRequest, "password incorrect"))
		return
	}

	c.JSON(http.StatusOK, common.NewSuccessResponse(nil))
}
