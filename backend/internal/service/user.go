package service

import (
	"context"
	"errors"
	"strings"
	"github.com/sudo-JP/Load-Manager/backend/internal/repository"
	"github.com/sudo-JP/Load-Manager/backend/internal/model"
	"github.com/sudo-JP/Load-Manager/backend/internal/salt"
	"github.com/sudo-JP/Load-Manager/backend/internal/hash"
)

type UserService struct {
	repo repository.UserRepositoryInterface
}


func (us *UserService) Create(ctx context.Context, u *model.User) (bool, error) {
	saltPW := salt.Salt()	
	hashed := hash.SHA256(u.Password + saltPW)
	u.Password = saltPW + ":" + hashed
	_, err := us.repo.Create(ctx, u)
	if err != nil {
		return false, err 
	}
	return true, nil
}

func (us *UserService) Delete(ctx context.Context, email string) (bool, error) {
	_, err := us.repo.DeleteUser(ctx, email) 
	if err != nil {
		return false, err
	}
	return true, nil
}

func verifyUser(ctx context.Context, us *UserService, unverified *model.User) (bool, error) {
	verified, err := us.repo.GetByEmail(ctx, unverified.Email)
	if err != nil {
		return false, err
	}

	idx := strings.Index(verified.Password, ":")
	if idx < 0 {
    	return false, errors.New("invalid stored password format")
	}
	verifiedSalt := verified.Password[:idx]
	verifiedHash := verified.Password[idx + 1:]
	checkUserHashed := hash.SHA256(unverified.Password + verifiedSalt) 
	if verifiedHash != checkUserHashed {
		return false, errors.New("unauthorized to change username")
	}
	return true, nil
}

func (us *UserService) UpdateUsername(ctx context.Context, u *model.User) (bool, error) {
	verified, err := verifyUser(ctx, us, u)
	if err != nil || !verified {
		return false, err
	}

	boolean, err := us.repo.UpdateUsername(ctx, u.Email, u.Name)
	if err != nil || !boolean {
		return false, err
	}

	return true, nil
}

func (us *UserService) UpdatePassword(ctx context.Context, u *model.User) (bool, error) {
	verified, err := verifyUser(ctx, us, u)
	if err != nil || !verified {
		return false, err
	}

	salty := salt.Salt()
	hashed := hash.SHA256(u.Password + salty)
	u.Password = salty + ":" + hashed

	boolean, err := us.repo.UpdatePassword(ctx, u.Email, u.Password)
	if err != nil || !boolean {
		return false, err
	}

	return true, nil
}

func (us *UserService) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	result, err := us.repo.GetByEmail(ctx, email) 
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (us *UserService) GetAll(ctx context.Context, email string) ([]model.User, error) {
	result, err := us.repo.ListAll(ctx)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func NewUserService(repo repository.UserRepositoryInterface) UserServiceInterface {
    return &UserService{repo: repo}
}

