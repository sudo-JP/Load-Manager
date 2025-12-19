package grpc

import (
	"log"

	"github.com/sudo-JP/Load-Manager/load-manager/api/proto/order"
	"github.com/sudo-JP/Load-Manager/load-manager/api/proto/product"
	"github.com/sudo-JP/Load-Manager/load-manager/api/proto/user"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type BackendClient struct {
	conn 		*grpc.ClientConn
	Users 		user.UserServiceClient
	Products 	product.ProductServiceClient
	Orders 		order.OrderServiceClient
}

// Constructor 
func NewBackendClient(address string) (*BackendClient, error) {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err 
	}

	return &BackendClient{
		conn: conn, 
		Users: user.NewUserServiceClient(conn), 
		Products: product.NewProductServiceClient(conn), 
		Orders: order.NewOrderServiceClient(conn),
	}, nil
}

// For defer 
func (bc *BackendClient) Close() {
	err := bc.conn.Close()

	if err != nil {
		log.Println(err)
	}
}

