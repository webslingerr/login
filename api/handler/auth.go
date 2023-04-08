package handler

import (
	"app/models"
	"app/pkg/helper"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Register User godoc
// @ID register_user
// @Router /register [POST]
// @Summary Register User
// @Description Register User
// @Tags Auth
// @Accept json
// @Produce json
// @Param User body models.Register true "RegisterUserRequest"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server Error"
func (h *Handler) RegisterUser(c *gin.Context) {
	var registerUser models.Register

	err := c.ShouldBindJSON(&registerUser)
	if err != nil {
		h.handlerResponse(c, "register user", http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.storages.User().Create(context.Background(), &models.CreateUser{
		FirstName:   registerUser.FirstName,
		LastName:    registerUser.LastName,
		Login:       registerUser.Login,
		Password:    registerUser.Password,
		PhoneNumber: registerUser.PhoneNumber,
	})
	if err != nil {
		h.handlerResponse(c, "storage register/create user", http.StatusInternalServerError, err.Error())
		return
	}

	user, err := h.storages.User().GetById(context.Background(), &models.UserPrimaryKey{
		Id: id,
	})
	if err != nil {
		h.handlerResponse(c, "storage get by id user inside register", http.StatusInternalServerError, err.Error())
		return
	}

	h.handlerResponse(c, "create user", http.StatusCreated, user)
}

// Login User godoc
// @ID login_user
// @Router /login [POST]
// @Summary Login User
// @Description Login User
// @Tags Auth
// @Accept json
// @Produce json
// @Param User body models.Login true "LoginUserRequest"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server Error"
func (h *Handler) Login(c *gin.Context) {
	var loginUser models.Login

	err := c.ShouldBindJSON(&loginUser)
	if err != nil {
		h.handlerResponse(c, "login user", http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.storages.User().GetByLoginPassword(context.Background(), &loginUser)
	if err != nil && err.Error() == "no rows in result set" {
		h.handlerResponse(c, "invalid password or login", http.StatusBadRequest, err)
		return
	}
	if err != nil {
		fmt.Println(err)
		h.handlerResponse(c, "get by user login and password", http.StatusInternalServerError, err)
		return
	}

	data := map[string]interface{}{
		"id":           user.Id,
		"first_name":   user.FirstName,
		"last_name":    user.LastName,
		"phone_number": user.PhoneNumber,
	}

	token, err := helper.GenerateJWT(data, time.Hour*24, h.cfg.SecretKey)
	if err != nil {
		h.handlerResponse(c, "error while generating token", http.StatusInternalServerError, err)
		return
	}

	h.handlerResponse(c, "login user", http.StatusOK, models.LoginResponse{
		AccessToken: token,
	})
}
