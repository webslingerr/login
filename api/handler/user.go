package handler

import (
	"app/models"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Create User godoc
// @ID create_user
// @Router /user [POST]
// @Summary Create User
// @Description Create User
// @Tags User
// @Accept json
// @Produce json
// @Param User body models.CreateUser true "CreateUserRequest"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server Error"
func (h *Handler) CreateUser(c *gin.Context) {
	var createUser models.CreateUser

	err := c.ShouldBindJSON(&createUser)
	if err != nil {
		h.handlerResponse(c, "create user", http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.storages.User().Create(context.Background(), &createUser)
	if err != nil {
		h.handlerResponse(c, "storage create user", http.StatusInternalServerError, err.Error())
		return
	}

	user, err := h.storages.User().GetById(context.Background(), &models.UserPrimaryKey{Id: id})
	if err != nil {
		h.handlerResponse(c, "storage get by id user", http.StatusInternalServerError, err.Error())
		return
	}

	h.handlerResponse(c, "create user", http.StatusCreated, user)
}

// Get By ID User godoc
// @ID get_by_id_user
// @Router /user/{id} [GET]
// @Summary Get By ID User
// @Description Get By ID User
// @Tags User
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server Error"
func (h *Handler) GetByIdUser(c *gin.Context) {
	id := c.Param("id")

	user, err := h.storages.User().GetById(context.Background(), &models.UserPrimaryKey{Id: id})
	if err != nil {
		h.handlerResponse(c, "Storage get by id user", http.StatusInternalServerError, err.Error())
		return
	}

	h.handlerResponse(c, "Get by id user", http.StatusOK, user)
}

// Get List User godoc
// @ID get_list_user
// @Router /user [GET]
// @Summary Get List User
// @Description Get List User
// @Tags User
// @Accept json
// @Produce json
// @Param offset query string false "offset"
// @Param limit query string false "limit"
// @Param search query string false "search"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server Error"
func (h *Handler) GetListUser(c *gin.Context) {
	offset, err := h.getOffsetQuery(c.Query("offset"))
	if err != nil {
		h.handlerResponse(c, "Get list user", http.StatusBadRequest, "invalid offset")
		return
	}

	limit, err := h.getLimitQuery(c.Query("limit"))
	if err != nil {
		h.handlerResponse(c, "Get list user", http.StatusBadRequest, "invalid limit")
		return
	}

	users, err := h.storages.User().GetList(context.Background(), &models.GetListUserRequest{
		Offset: offset,
		Limit:  limit,
		Search: c.Query("search"),
	})
	if err != nil {
		h.handlerResponse(c, "Storage get list user", http.StatusInternalServerError, err.Error())
		return
	}

	h.handlerResponse(c, "Get list user", http.StatusOK, users)
}

// Get Update User godoc
// @ID update_user
// @Router /user/{id} [PUT]
// @Summary Update User
// @Description Update User
// @Tags User
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Param USer body models.UpdateUser true "UpdateUserRequest"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server Error"
func (h *Handler) UpdateUser(c *gin.Context) {
	var updateUser models.UpdateUser

	id := c.Param("id")

	err := c.ShouldBindJSON(&updateUser)
	if err != nil {
		h.handlerResponse(c, "Update user", http.StatusBadRequest, err.Error())
		return
	}
	updateUser.Id = id

	rowsAffected, err := h.storages.User().Update(context.Background(), &updateUser)
	if err != nil {
		h.handlerResponse(c, "Storage update user", http.StatusInternalServerError, err.Error())
		return
	}

	if rowsAffected <= 0 {
		h.handlerResponse(c, "Storage update user", http.StatusBadRequest, "no rows affected")
		return
	}

	user, err := h.storages.User().GetById(context.Background(), &models.UserPrimaryKey{Id: id})
	if err != nil {
		h.handlerResponse(c, "Storage get by id user", http.StatusInternalServerError, err.Error())
		return
	}

	h.handlerResponse(c, "Update user", http.StatusOK, user)
}

// Delete User godoc
// @ID delete_user
// @Router /user/{id} [DELETE]
// @Summary Delete User
// @Description Delete User
// @Tags User
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Param User body models.UserPrimaryKey true "DeleteUserRequest"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server Error"
func (h *Handler) DeleteUser(c *gin.Context) {
	id := c.Param("id")

	rowsAffected, err := h.storages.User().Delete(context.Background(), &models.UserPrimaryKey{Id: id})
	if err != nil {
		h.handlerResponse(c, "Storage delete category", http.StatusInternalServerError, err.Error())
		return
	}

	if rowsAffected <= 0 {
		h.handlerResponse(c, "Storage delete user", http.StatusBadRequest, "no rows affected")
		return
	}

	h.handlerResponse(c, "Delete user", http.StatusNoContent, "Deleted Successfully")
}
