package auth

import (
	"github.com/Nistagram-Organization/nistagram-auth/src/dto/registration_request"
	"github.com/Nistagram-Organization/nistagram-auth/src/services/auth"
	"github.com/Nistagram-Organization/nistagram-shared/src/utils/rest_error"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AuthController interface {
	Register(ctx *gin.Context)
}

type authController struct {
	authService auth.AuthService
}

func NewAuthController(authService auth.AuthService) AuthController {
	return &authController{
		authService,
	}
}

func (c *authController) Register(ctx *gin.Context) {
	var registrationRequest registration_request.RegistrationRequest

	if err := ctx.ShouldBindJSON(&registrationRequest); err != nil {
		restErr := rest_error.NewBadRequestError("invalid json body")
		ctx.JSON(restErr.Status(), restErr)
		return
	}

	if err := c.authService.Register(registrationRequest); err != nil {
		ctx.JSON(err.Status(), err)
		return
	}

	ctx.Status(http.StatusCreated)
}
