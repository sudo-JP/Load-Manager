package grpc

import (
	pb "github.com/sudo-JP/Load-Manager/load-manager/api/proto/user"
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (bc *BackendClient) CreateUsers(ctx context.Context, req *pb.CreateUsersRequest) (*emptypb.Empty, error) {
	return bc.Users.CreateUsers(ctx, req)
}

func (bc *BackendClient) GetUsers(ctx context.Context,
	req *pb.GetUsersRequest) (*pb.GetUsersResponse, error) {
	return bc.Users.GetUsers(ctx, req)
}

func (bc *BackendClient) UpdateUsers(ctx context.Context,
	req *pb.UpdateUsersRequest) (*emptypb.Empty, error) {
	return bc.Users.UpdateUsers(ctx, req)
}

func (bc *BackendClient) DeleteUsers(ctx context.Context,
	req *pb.DeleteUsersRequest) (*emptypb.Empty, error) {
	return bc.Users.DeleteUsers(ctx, req)
}
