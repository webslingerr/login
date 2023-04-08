package postgresql

import (
	"context"
	"database/sql"
	"fmt"

	"app/models"
	"app/pkg/helper"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

type userRepo struct {
	db *pgxpool.Pool
}

func NewUserRepo(db *pgxpool.Pool) *userRepo {
	return &userRepo{
		db: db,
	}
}

func (u *userRepo) Create(ctx context.Context, req *models.CreateUser) (string, error) {
	id := uuid.New().String()

	query := `
		INSERT INTO users(
			id,
			first_name,
			last_name,
			login,
			password, 
			phone_number, 
			updated_at
		) VALUES($1, $2, $3, $4, $5, $6, NOW())
	`

	_, err := u.db.Exec(ctx, query,
		id,
		helper.NewNullString(req.FirstName),
		helper.NewNullString(req.LastName),
		helper.NewNullString(req.Login),
		helper.NewNullString(req.Password),
		helper.NewNullString(req.PhoneNumber),
	)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (u *userRepo) GetById(ctx context.Context, req *models.UserPrimaryKey) (*models.User, error) {
	var (
		id           sql.NullString
		first_name   sql.NullString
		last_name    sql.NullString
		login        sql.NullString
		password     sql.NullString
		phone_number sql.NullString
		created_at   sql.NullString
		updated_at   sql.NullString
	)

	query := `
		SELECT 
			id,
			first_name,
			last_name,
			login,
			password, 
			phone_number,
			created_at,
			updated_at
		FROM users
		WHERE id = $1
	`

	err := u.db.QueryRow(ctx, query, req.Id).Scan(
		&id,
		&first_name,
		&last_name,
		&login,
		&password,
		&phone_number,
		&created_at,
		&updated_at,
	)
	if err != nil {
		return nil, err
	}

	return &models.User{
		Id:          id.String,
		FirstName:   first_name.String,
		LastName:    last_name.String,
		Login:       login.String,
		Password:    password.String,
		PhoneNumber: phone_number.String,
		CreatedAt:   created_at.String,
		UpdatedAt:   updated_at.String,
	}, nil
}

func (u *userRepo) GetByLoginPassword(ctx context.Context, req *models.Login) (*models.User, error) {
	var (
		id           sql.NullString
		first_name   sql.NullString
		last_name    sql.NullString
		login        sql.NullString
		password     sql.NullString
		phone_number sql.NullString
		created_at   sql.NullString
		updated_at   sql.NullString
	)

	query := `
		SELECT 
			id,
			first_name,
			last_name,
			login,
			password, 
			phone_number,
			created_at,
			updated_at
		FROM users
		WHERE login = $1 AND password = $2
	`

	err := u.db.QueryRow(ctx, query, req.Login, req.Password).Scan(
		&id,
		&first_name,
		&last_name,
		&login,
		&password,
		&phone_number,
		&created_at,
		&updated_at,
	)
	if err != nil {
		return nil, err
	}

	return &models.User{
		Id:          id.String,
		FirstName:   first_name.String,
		LastName:    last_name.String,
		Login:       login.String,
		Password:    password.String,
		PhoneNumber: phone_number.String,
		CreatedAt:   created_at.String,
		UpdatedAt:   updated_at.String,
	}, nil
}

func (u *userRepo) GetList(ctx context.Context, req *models.GetListUserRequest) (*models.GetListUserResponse, error) {
	resp := models.GetListUserResponse{}

	var (
		query  string
		filter = " WHERE TRUE "
		offset string
		limit  string
	)

	query = `
		SELECT
			id,
			first_name,
			last_name,
			login,
			password, 
			phone_number,
			created_at,
			updated_at
		FROM users
	`

	if len(req.Search) > 0 {
		filter += " AND (first_name || ' ' last_name) ILIKE '%' || '" + req.Search + "' || '%' "
	}
	if req.Offset > 0 {
		offset = fmt.Sprintf(" OFFSET %d", req.Offset)
	}
	if req.Limit > 0 {
		limit = fmt.Sprintf(" LIMIT %d", req.Limit)
	}

	query += filter + offset + limit

	rows, err := u.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var id, first_name, last_name, login, password, phone_number, created_at, updated_at sql.NullString

		err = rows.Scan(
			&id,
			&first_name,
			&last_name,
			&login,
			&password,
			&phone_number,
			&created_at,
			&updated_at,
		)
		if err != nil {
			return nil, err
		}

		resp.Users = append(resp.Users, &models.User{
			Id:          id.String,
			FirstName:   first_name.String,
			LastName:    last_name.String,
			Login:       login.String,
			Password:    password.String,
			PhoneNumber: phone_number.String,
			CreatedAt:   created_at.String,
			UpdatedAt:   updated_at.String,
		})
	}
	resp.Count = len(resp.Users)

	return &resp, nil
}

func (u *userRepo) Update(ctx context.Context, req *models.UpdateUser) (int64, error) {
	query := `
		UPDATE users
		SET 
			first_name = :first_name,
			last_name = :last_name,
			login = :login,
			password = :password,
			phone_number = :phone_number
			updated_at = NOW()
		WHERE id = :id
	`

	params := map[string]interface{}{
		"id":           req.Id,
		"first_name":   req.FirstName,
		"last_name":    req.LastName,
		"login":        req.Login,
		"password":     req.Password,
		"phone_number": req.PhoneNumber,
	}

	query, args := helper.ReplaceQueryParams(query, params)

	rowsAffected, err := u.db.Exec(ctx, query, args...)
	if err != nil {
		return 0, nil
	}
	return rowsAffected.RowsAffected(), nil
}

func (u *userRepo) Delete(ctx context.Context, req *models.UserPrimaryKey) (int64, error) {
	rowsAffected, err := u.db.Exec(ctx, `DELETE FROM users WHERE id = $1`, req.Id)
	if err != nil {
		return 0, err
	}

	return rowsAffected.RowsAffected(), nil
}
