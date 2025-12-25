package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/sudo-JP/Load-Manager/backend/internal/database"
	"github.com/sudo-JP/Load-Manager/backend/internal/model"
)

type User struct {
	db *database.Database
}

func (r *User) CreateUser(ctx context.Context, user model.User)  (*model.User, error) {
	u := model.User {
		Name: user.Name, 
		Email: user.Email, 
		Password: user.Password,
	}

	err := r.db.Pool.QueryRow(
		ctx,
		"INSERT INTO users (name, email, password) VALUES ($1, $2, $3) RETURNING user_id;",
		user.Name, user.Email, user.Password,
	).Scan(&u.UserId)

	if err != nil {
		return nil, err
	}
	return &u, nil
}

// Bulk create users
func (r *User) CreateUsers(ctx context.Context, users []model.User) error {
	if len(users) == 0 {
		return nil
	}

	rows := make([][]any, 0, len(users))
	for _, u := range users {
		rows = append(rows, []any{
			u.Name,
			u.Email,
			u.Password,
		})
	}

	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	_, err = tx.CopyFrom(
		ctx,
		pgx.Identifier{"users"},
		[]string{"name", "email", "password"},
		pgx.CopyFromRows(rows),
	)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

// Get user by email
func (r *User) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	var u model.User
	err := r.db.Pool.QueryRow(
		ctx,
		"SELECT user_id, name, password, email FROM users WHERE email=$1",
		email,
	).Scan(&u.UserId, &u.Name, &u.Password, &u.Email)

	if err != nil {
		return nil, errors.New("unable to get user by email")
	}
	return &u, nil
}

// GetById fetches a user by id
func (r *User) GetById(ctx context.Context, userId int) (*model.User, error) {
	var u model.User
	err := r.db.Pool.QueryRow(ctx,
		"SELECT user_id, name, email, password FROM users WHERE user_id=$1",
		userId,
	).Scan(&u.UserId, &u.Name, &u.Email, &u.Password)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

// List all users
func (r *User) ListAll(ctx context.Context) ([]model.User, error) {
	rows, err := r.db.Pool.Query(ctx, "SELECT user_id, name, email, password FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]model.User, 0)
	for rows.Next() {
		var u model.User
		if err := rows.Scan(&u.UserId, &u.Name, &u.Email, &u.Password); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

// Update user password
func (r *User) UpdatePassword(ctx context.Context, email string, password string) error {
	res, err := r.db.Pool.Exec(
		ctx,
		"UPDATE users SET password=$1 WHERE email=$2",
		password, email,
	)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return errors.New("user not found, unable to update password")
	}
	return nil
}

// Update user name
func (r *User) UpdateUsername(ctx context.Context, email string, name string) error {
	res, err := r.db.Pool.Exec(
		ctx,
		"UPDATE users SET name=$1 WHERE email=$2",
		name, email,
	)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return errors.New("user not found, unable to update name")
	}
	return nil
}


// Bulk delete users by email
func (r *User) DeleteUsers(ctx context.Context, emails []string) error {
	_, err := r.db.Pool.Exec(
		ctx,
		"DELETE FROM users WHERE email = ANY($1);",
		emails,
	)

	return err
}

// Constructor
func NewUserRepository(db *database.Database) UserRepositoryInterface {
	return &User{db: db}
}
