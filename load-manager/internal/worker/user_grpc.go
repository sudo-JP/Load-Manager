package worker

import (
	"context"
	"encoding/json"
	"time"

	"github.com/sudo-JP/Load-Manager/load-manager/internal/queue"
	"github.com/sudo-JP/Load-Manager/load-manager/internal/registry"

	pb "github.com/sudo-JP/Load-Manager/load-manager/api/proto/user"

	"log"
)


func (w *Worker) GetUsers(node *registry.BackendNode, 
	jobs []*queue.Job) {
	type GetUserDTO struct {
		Email string `json:"email"`
	}

	for _, job := range jobs {
		var dto GetUserDTO 
		if err := json.Unmarshal(job.Payload, &dto); err != nil {
			log.Printf("Failed to Unmarshal user: %v", err)
			continue
		}
		req := &pb.GetUsersRequest{
			Email: dto.Email,
		}

		go func() {
		    ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
		    defer cancel()
			
			client, err := w.getClient(node)
			if err != nil {
				log.Println(err)
				return 
			}
			resp, err := client.Users.GetUsers(ctx, req)
			if err != nil {
				log.Printf("gRPC GetUsers failed: %v", err)
				return 
			}
			log.Printf("Retrived %d users", len(resp.Users))
		}()

	}


}

func (w *Worker) CreateUsers(node *registry.BackendNode,
	jobs []*queue.Job) {
    // Unmarshal all payloads
    type CreateUserDTO struct {
    	Name 	 string `json:"name"`
    	Email 	 string `json:"email"`
    	Password string `json:"password"`
    }

    var users []*pb.User 
    for _, job := range jobs {
		var dto CreateUserDTO 
		if err := json.Unmarshal(job.Payload, &dto); err != nil {
			log.Printf("Failed to Unmarshal user: %v", err)
			continue
		}

		// To proto
		users = append(users, &pb.User{
			Name: 		dto.Name,
			Email: 		dto.Email,
			Password: 	dto.Password,
		})
    }

    if len(users) == 0 {
    	return 
    }

    // make grpc req
    req := &pb.CreateUsersRequest{
    	Users: users,
    }

    // Send grpc 
    ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
    defer cancel()

	client, err := w.getClient(node)
	if err != nil {
		log.Println(err)
		return 
	}

	_, err = client.Users.CreateUsers(ctx, req)

	if err != nil {
		log.Printf("gRPC CreateUsers failed for node %s:%d: %v",
			node.Host, node.Port, err)
		return 
	}
	log.Printf("Created users on mode %s:%d", node.Host, node.Port)
}

func (w *Worker) DeleteUsers(node *registry.BackendNode,
	jobs []*queue.Job) {
    // Unmarshal all payloads
    type DeleteUserDTO struct {
    	Name 	 string `json:"name"`
    	Email 	 string `json:"email"`
    	Password string `json:"password"`
    }

    var users []*pb.User 
    for _, job := range jobs {
		var dto DeleteUserDTO 
		if err := json.Unmarshal(job.Payload, &dto); err != nil {
			log.Printf("Failed to Unmarshal user: %v", err)
			continue
		}

		// To proto
		users = append(users, &pb.User{
			Name: 		dto.Name,
			Email: 		dto.Email,
			Password: 	dto.Password,
		})
    }

    if len(users) == 0 {
    	return 
    }

    // make grpc req
    req := &pb.DeleteUsersRequest{
    	Users: users,
    }

    // Send grpc 
    ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
    defer cancel()

	client, err := w.getClient(node)
	if err != nil {
		log.Println(err)
		return 
	}

	_, err = client.Users.DeleteUsers(ctx, req)

	if err != nil {
		log.Printf("gRPC CreateUsers failed for node %s:%d: %v",
			node.Host, node.Port, err)
		return 
	}
	log.Printf("Deleted users on mode %s:%d", node.Host, node.Port)
}

func (w *Worker) UpdateUsers(node *registry.BackendNode,
	jobs []*queue.Job) {
    // Unmarshal all payloads
    type UpdateUserDTO struct {
    	Name 	 string `json:"name"`
    	Email 	 string `json:"email"`
    	Password string `json:"password"`
    }

    var users []*pb.User 
    for _, job := range jobs {
		var dto UpdateUserDTO 
		if err := json.Unmarshal(job.Payload, &dto); err != nil {
			log.Printf("Failed to Unmarshal user: %v", err)
			continue
		}

		// To proto
		users = append(users, &pb.User{
			Name: 		dto.Name,
			Email: 		dto.Email,
			Password: 	dto.Password,
		})
    }

    if len(users) == 0 {
    	return 
    }

    // make grpc req
    req := &pb.UpdateUsersRequest{
    	Users: users,
    }

    // Send grpc 
    ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
    defer cancel()

	client, err := w.getClient(node)
	if err != nil {
		log.Println(err)
		return 
	}

	_, err = client.Users.UpdateUsers(ctx, req)

	if err != nil {
		log.Printf("gRPC CreateUsers failed for node %s:%d: %v",
			node.Host, node.Port, err)
		return 
	}
	log.Printf("Updated users on mode %s:%d", node.Host, node.Port)
}

