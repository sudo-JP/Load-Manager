# LOAD MANAGER ARCHITECTURE 


## Tree Structure
```bash
.
├── api
│   └── proto
│       ├── order
│       │   ├── order_grpc.pb.go
│       │   ├── order.pb.go
│       │   └── order.proto
│       ├── product
│       │   ├── product_grpc.pb.go
│       │   ├── product.pb.go
│       │   └── product.proto
│       └── user
│           ├── user_grpc.pb.go
│           ├── user.pb.go
│           └── user.proto
├── cmd
│   └── load-manager
│       └── main.go
├── go.mod
├── go.sum
├── internal
│   ├── batcher
│   │   └── batcher.go
│   ├── grpc
│   │   ├── client.go
│   │   ├── order.go
│   │   ├── product.go
│   │   └── user.go
│   ├── model
│   │   ├── order.go
│   │   ├── product.go
│   │   └── user.go
│   ├── queue
│   │   ├── algorithms
│   │   │   ├── fcfs.go
│   │   │   ├── fcfs_test.go
│   │   │   ├── ljf.go
│   │   │   ├── queue_test.go
│   │   │   ├── random.go
│   │   │   └── sjf.go
│   │   ├── job.go
│   │   └── queue.go
│   ├── registry
│   │   └── registry.go
│   ├── routes
│   │   ├── order.go
│   │   ├── product.go
│   │   └── user.go
│   ├── selector
│   │   ├── rand.go
│   │   ├── rr.go
│   │   └── selector.go
│   └── worker
│       ├── order_grpc.go
│       ├── product_grpc.go
│       ├── user_grpc.go
│       └── worker.go
└── README.md
```

