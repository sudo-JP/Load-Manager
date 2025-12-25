package service

import (
	"context"
	"errors"
	"fmt"
	"runtime"
	"strings"
	"sync"

	pb "github.com/sudo-JP/Load-Manager/backend/api/proto/user"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/sudo-JP/Load-Manager/backend/internal/hash"
	"github.com/sudo-JP/Load-Manager/backend/internal/model"
	"github.com/sudo-JP/Load-Manager/backend/internal/repository"
	"github.com/sudo-JP/Load-Manager/backend/internal/salt"
)

type User struct {
	repo repository.UserRepositoryInterface
}

// verifyUser checks the password before updates
func verifyUser(ctx context.Context, us *User, unverified *model.User) error {
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

func protoToUser(user *pb.User) model.User {
	return model.User{
		UserId:   int(user.UserId),
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
	}
}

func protoToUsers(users []*pb.User) []model.User {
	result := make([]model.User, len(users))
	for i, user := range users {
		result[i] = protoToUser(user)
	}
	return result
}

func (us *User) ProtoCreateUsers(ctx context.Context, req *pb.CreateUsersRequest) (*emptypb.Empty, error) {
	users := protoToUsers(req.Users)
	if len(users) == 0 {
		return &emptypb.Empty{}, nil
	}

	if err := us.CreateUsers(ctx, users); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

// CreateUser creates a single user with password hashing
func (us *User) CreateUser(ctx context.Context, user model.User) (*model.User, error) {
	saltPW := salt.Salt()
	hashed := hash.SHA256(user.Password + saltPW)
	user.Password = saltPW + ":" + hashed

	created, err := us.repo.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}
	return created, nil
}

// CreateUsers hashes passwords concurrently and inserts into repo
func (us *User) CreateUsers(ctx context.Context, users []model.User) error {

	threadsNum := runtime.NumCPU()
	jobs := make(chan model.User, threadsNum*2)
	results := make(chan model.User, len(users))
	var wg sync.WaitGroup

	for range threadsNum {
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

	if err := us.repo.CreateUsers(ctx, hashedUsers); err != nil {
		return err
	}
	return nil
}

func (us *User) ProtoGetUsers(ctx context.Context, req *pb.GetUsersRequest) (*pb.GetUsersResponse, error) {
	var users []model.User
	var u model.User
	var err error

	if req.Email == "" {
		users, err = us.ListUsers(ctx)
		if err != nil {
			return nil, err
		}
	} else {
		u, err = us.GetUser(ctx, req.Email)
		users = []model.User{u}
	}

	if err != nil {
		return nil, err
	}

	pbUsers := make([]*pb.User, len(users))
	for i, user := range users {
		pbUsers[i] = &pb.User{
			UserId: int64(user.UserId),
			Name:   user.Name,
			Email:  user.Email,
		}
	}

	return &pb.GetUsersResponse{
		Users: pbUsers,
	}, nil
}

// GetUser fetches a single user
func (us *User) GetUser(ctx context.Context, email string) (model.User, error) {
	u, err := us.repo.GetByEmail(ctx, email)
	if err != nil {
		return model.User{}, err
	}
	return *u, nil
}

// GetUserById fetches a single user by id
func (us *User) GetUserById(ctx context.Context, id int) (model.User, error) {
	u, err := us.repo.GetById(ctx, id)
	if err != nil {
		return model.User{}, err
	}
	if u == nil {
		return model.User{}, fmt.Errorf("user not found")
	}
	return *u, nil
}

// ListUsers fetches all users
func (us *User) ListUsers(ctx context.Context) ([]model.User, error) {
	return us.repo.ListAll(ctx)
}

func (us *User) ProtoUpdateUsers(ctx context.Context, req *pb.UpdateUsersRequest) (*emptypb.Empty, error) {
	users := protoToUsers(req.Users)
	err := us.UpdateUsers(ctx, users)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, err
}

// UpdateUsers updates usernames and/or passwords in bulk
func (us *User) UpdateUsers(ctx context.Context, updates []model.User) error {
	for _, u := range updates {
		if u.Password != "" {
			if err := verifyUser(ctx, us, &u); err != nil {
				return err
			}
			// New password
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

func protoToEmail(users []*pb.User) []string {
	result := make([]string, len(users))

	for i, user := range users {
		result[i] = user.Email
	}
	return result
}

func (us *User) ProtoDeleteUsers(ctx context.Context, req *pb.DeleteUsersRequest) (*emptypb.Empty, error) {
	emails := protoToEmail(req.Users)
	err := us.DeleteUsers(ctx, emails)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

// DeleteUsers deletes multiple users by email
func (us *User) DeleteUsers(ctx context.Context, emails []string) error {
	return us.repo.DeleteUsers(ctx, emails)
}

// UpdateUser updates a single user (full update)
func (us *User) UpdateUser(ctx context.Context, user model.User) error {
	if err := verifyUser(ctx, us, &user); err != nil {
		return err
	}

	// Hash new password
	saltPW := salt.Salt()
	hashed := hash.SHA256(user.Password + saltPW)
	user.Password = saltPW + ":" + hashed

	return us.repo.UpdateUser(ctx, user)
}

// DeleteUser deletes a single user by email
func (us *User) DeleteUser(ctx context.Context, email string) error {
	return us.repo.DeleteUser(ctx, email)
}

// Constructor
func NewUserService(repo repository.UserRepositoryInterface) UserServiceInterface {
	return &User{repo: repo}
}
