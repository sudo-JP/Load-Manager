package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/sudo-JP/Load-Manager/backend/internal/database"
	"github.com/sudo-JP/Load-Manager/backend/internal/model"
)

type UserRepository struct {
    db *database.Database
}


func (r *UserRepository) Create(ctx context.Context, u *model.User) (bool, error) {
	err := r.db.Pool.QueryRow(
		ctx, 
		"INSERT INTO users (name, email, password) VALUES ($1, $2, $3) RETURNING user_id;",
		u.Name, u.Email, u.Password,
	).Scan(&u.UserId)

	if err != nil {
		// Check for PostgreSQL-specific error
		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code == "23505" {
				return false, errors.New("DUPLICATE EMAIL")
			}
		}
		return false, err
	}

	return true, nil
}


func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	var name string 
	var password string 
	var user_id int 

	err := r.db.Pool.QueryRow(
		ctx,
		"SELECT name, password, user_id FROM users WHERE email = $1;",
		email,
	).Scan(&name, &password, &user_id)

	if err != nil {
		return nil, errors.New("UNABLE TO GET USER BY EMAIL;") 
	}


	u := &model.User{
    	Name:     name,
    	Email:    email,
    	Password: password,
    	UserId:   user_id,
	}
	return u, nil
}

func (r *UserRepository) ListAll(ctx context.Context) ([]model.User, error) {
	result, err := r.db.Pool.Query(
		ctx, 
		"SELECT user_id, name, email, password from users;", 
	)

	if err != nil {
		return nil, errors.New("UNABLE TO GET ALL USERS")
	}

	defer result.Close()

	users := []model.User{}

	for result.Next() {
		var u model.User
		err := result.Scan(&u.UserId, &u.Name, &u.Email, &u.Password)
		if err != nil {
			return nil, errors.New("UNABLE TO PARSE USERS")
		}
		users = append(users, u)
	}

	return users, nil 
}

func (r *UserRepository) DeleteUser(ctx context.Context, email string) (bool, error) {
	result, err := r.db.Pool.Exec(
		ctx, 
		"DELETE FROM users WHERE email = $1;", 
		email, 
	)
	
	if err != nil {
		return false, err
	}
	if result.RowsAffected() == 0 {
		return false, errors.New("USER DOESN'T EXIST, CAN'T BE DELETED")	
	}

	return true, nil 
}

func (r *UserRepository) UpdatePassword(ctx context.Context, email string, password string) (bool, error) {
	result, err := r.db.Pool.Exec(
		ctx, 
		"UPDATE users SET password = $1 WHERE email = $2;", 
		password, email,
	)

	if err != nil {
		return false, err
	}

	if result.RowsAffected() == 0 {
		return false, errors.New("CAN'T FIND USER WITH EMAIL, UNABLE TO UPDATE PASSWORD")
	}

	return true, nil 
}

func (r *UserRepository) UpdateUsername(ctx context.Context, email string, name string) (bool, error) {
	result, err := r.db.Pool.Exec(
		ctx, 
		"UPDATE users SET name = $1 WHERE email = $2;", 
		name, email,
	)

	if err != nil {
		return false, err
	}

	if result.RowsAffected() == 0 {
		return false, errors.New("CAN'T FIND USER WITH EMAIL, UNABLE TO UPDATE name")
	}

	return true, nil 
}

func (r *UserRepository) CreateUsers(ctx context.Context, users []model.User) (bool, error) {

	rows := [][]any{}

	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
    	return false, err
	}
	// defer rollback, but ignore the error if commit succeeds
	defer func() {
    	_ = tx.Rollback(ctx)
	}()

	_, err = tx.CopyFrom(
    	ctx,
    	pgx.Identifier{"users"},
    	[]string{"name", "email", "password"},
    	pgx.CopyFromRows(rows),
	)
	if err != nil {
    	return false, err
	}

	err = tx.Commit(ctx)
	if err != nil {
    	return false, err
	}

	return true, nil
}

func NewUserRepository(db *database.Database) UserRepositoryInterface {
	return &UserRepository{ db: db }	
}
