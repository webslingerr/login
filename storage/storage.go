package storage

import (
	"app/models"
	"context"
)

type StorageI interface {
	CloseDb()
	User() UserRepoI
}

type UserRepoI interface {
	Create(context.Context, *models.CreateUser) (string, error)
	GetById(context.Context, *models.UserPrimaryKey) (*models.User, error)
	GetByLoginPassword(context.Context, *models.Login) (*models.User, error)
	GetList(context.Context, *models.GetListUserRequest) (*models.GetListUserResponse, error)
	Update(context.Context, *models.UpdateUser) (int64, error)
	Delete(context.Context, *models.UserPrimaryKey) (int64, error)
}