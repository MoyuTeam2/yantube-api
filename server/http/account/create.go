package account

import (
	"api/config"
	"api/db"
	"api/models"
	"api/server/http/common"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type CreateUserRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Create(c *gin.Context) {
	if !config.Config.User.AllowRegister {
		c.JSON(http.StatusForbidden, common.NewErrorResponse(http.StatusForbidden, "register is not allowed"))
		return
	}

	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error().Err(err).Msg("create user request bind json failed")
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(http.StatusBadRequest, "bad request"))
		return
	}

	user, err := models.NewUserWithPassword(req.Username, req.Password)
	if err != nil {
		log.Error().Err(err).Msg("create user failed")
		c.JSON(http.StatusInternalServerError, common.NewErrorResponse(http.StatusInternalServerError, "server create user failed"))
		return
	}

	if err := db.Get().CreateUser(user); err != nil {
		log.Error().Err(err).Msg("save user failed")
		c.JSON(http.StatusInternalServerError, common.NewErrorResponse(http.StatusInternalServerError, "server save user failed"))
		return
	}

	c.JSON(http.StatusOK, common.NewSuccessResponse(nil))
}
