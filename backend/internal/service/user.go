package service

import (
	"context"
	"errors"
	"runtime"
	"strings"
	"sync"

	"github.com/sudo-JP/Load-Manager/backend/internal/hash"
	"github.com/sudo-JP/Load-Manager/backend/internal/model"
	"github.com/sudo-JP/Load-Manager/backend/internal/repository"
	"github.com/sudo-JP/Load-Manager/backend/internal/salt"
)

type UserService struct {
	repo repository.UserRepositoryInterface
}

// verifyUser checks the password before updates
func verifyUser(ctx context.Context, us *UserService, unverified *model.User) error {
	stored, err := us.repo.GetByEmail(ctx, unverified.Email)
	if err != nil {
		return err
	}

	idx := strings.Index(stored.Password, ":")
	if idx < 0 {
		return errors.New("invalid stored password format")
	}
	storedSalt := stored.Password[:idx]
	storedHash := stored.Password[idx+1:]
	checkHash := hash.SHA256(unverified.Password + storedSalt)
	if storedHash != checkHash {
		return errors.New("unauthorized")
	}
	return nil
}

// CreateUsers hashes passwords concurrently and inserts into repo
func (us *UserService) CreateUsers(ctx context.Context, users []model.User) error {
	if len(users) == 0 {
		return nil
	}

	threadsNum := runtime.NumCPU()
	jobs := make(chan model.User, threadsNum*2)
	results := make(chan model.User, len(users))
	var wg sync.WaitGroup

	for i := 0; i < threadsNum; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for u := range jobs {
				saltPW := salt.Salt()
				hashed := hash.SHA256(u.Password + saltPW)
				u.Password = saltPW + ":" + hashed
				results <- u
			}
		}()
	}

	for _, u := range users {
		jobs <- u
	}
	close(jobs)

	go func() {
		wg.Wait()
		close(results)
	}()

	var hashedUsers []model.User
	for u := range results {
		hashedUsers = append(hashedUsers, u)
	}

	return us.repo.CreateUsers(ctx, hashedUsers)
}

// UpdateUsers updates usernames and/or passwords in bulk
func (us *UserService) UpdateUsers(ctx context.Context, updates []model.User) error {
	for _, u := range updates {
		if u.Password != "" {
			if err := verifyUser(ctx, us, &u); err != nil {
				return err
			}
			saltPW := salt.Salt()
			hashed := hash.SHA256(u.Password + saltPW)
			u.Password = saltPW + ":" + hashed
			if err := us.repo.UpdatePassword(ctx, u.Email, u.Password); err != nil {
				return err
			}
		}
		if u.Name != "" {
			if err := verifyUser(ctx, us, &u); err != nil {
				return err
			}
			if err := us.repo.UpdateUsername(ctx, u.Email, u.Name); err != nil {
				return err
			}
		}
	}
	return nil
}

// DeleteUsers deletes multiple users by email
func (us *UserService) DeleteUsers(ctx context.Context, emails []string) error {
	return us.repo.DeleteUsers(ctx, emails)
}

// GetUser fetches a single user
func (us *UserService) GetUser(ctx context.Context, email string) (model.User, error) {
	u, err := us.repo.GetByEmail(ctx, email)
	if err != nil {
		return model.User{}, err
	}
	return *u, nil
}

// ListUsers fetches all users
func (us *UserService) ListUsers(ctx context.Context) ([]model.User, error) {
	return us.repo.ListAll(ctx)
}

// Constructor
func NewUserService(repo repository.UserRepositoryInterface) UserServiceInterface {
	return &UserService{repo: repo}
}
