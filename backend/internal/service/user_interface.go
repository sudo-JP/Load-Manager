package service

import (
	"context"

	"github.com/sudo-JP/Load-Manager/backend/internal/model"

	pb "github.com/sudo-JP/Load-Manager/backend/api/proto/user"
	"google.golang.org/protobuf/types/known/emptypb"
)

type UserServiceInterface interface {
	// Internal - Singleton
	CreateUser(ctx context.Context, user model.User) (*model.User, error)
	GetUser(ctx context.Context, email string) (model.User, error)
	GetUserById(ctx context.Context, id int) (model.User, error)
	UpdateUser(ctx context.Context, user model.User) error
	DeleteUser(ctx context.Context, email string) error
	ListUsers(ctx context.Context) ([]model.User, error)

	// Internal - Batch
	CreateUsers(ctx context.Context, users []model.User) error
	UpdateUsers(ctx context.Context, updates []model.User) error
	DeleteUsers(ctx context.Context, emails []string) error

	// Proto
	ProtoCreateUsers(ctx context.Context, req *pb.CreateUsersRequest) (*emptypb.Empty, error)
	ProtoGetUsers(ctx context.Context, req *pb.GetUsersRequest) (*pb.GetUsersResponse, error)
	ProtoUpdateUsers(ctx context.Context, req *pb.UpdateUsersRequest) (*emptypb.Empty, error)
	ProtoDeleteUsers(ctx context.Context, req *pb.DeleteUsersRequest) (*emptypb.Empty, error)
}
