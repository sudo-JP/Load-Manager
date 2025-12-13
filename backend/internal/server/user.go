package server

import (
	"context"

	pb "github.com/sudo-JP/Load-Manager/backend/api/proto/user"
	"github.com/sudo-JP/Load-Manager/backend/internal/service"
	"google.golang.org/protobuf/types/known/emptypb"
)

type UserServer struct {
	pb.UnimplementedUserServiceServer
	svc service.UserServiceInterface
}

func (s *UserServer) CreateUsers(ctx context.Context, 
	req *pb.CreateUsersRequest) (*emptypb.Empty, error) {
	return s.svc.ProtoCreateUsers(ctx, req)
}

func (s *UserServer) GetUsers(ctx context.Context,
	req *pb.GetUsersRequest) (*pb.GetUsersResponse, error) {
	return s.svc.ProtoGetUsers(ctx, req)
}

func (s *UserServer) UpdateUsers(ctx context.Context,
	req *pb.UpdateUsersRequest) (*emptypb.Empty, error) {
	return s.svc.ProtoUpdateUsers(ctx, req)
}

func NewUserServer(svc service.UserServiceInterface) *UserServer {
	return &UserServer{svc: svc}
}
