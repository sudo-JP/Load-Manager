# BACKEND ARCHITECTURE

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
│   └── backend
│       └── main.go
├── go.mod
├── go.sum
├── internal
│   ├── config
│   │   └── config.go
│   ├── database
│   │   └── connection.go
│   ├── hash
│   │   └── sha256.go
│   ├── migrations
│   │   └── 0001_User.sql
│   ├── model
│   │   ├── order.go
│   │   ├── product.go
│   │   └── user.go
│   ├── repository
│   │   ├── order.go
│   │   ├── order_interface.go
│   │   ├── product.go
│   │   ├── product_interface.go
│   │   ├── user.go
│   │   └── user_interface.go
│   ├── routes
│   │   ├── order.go
│   │   ├── product.go
│   │   └── user.go
│   ├── salt
│   │   └── random.go
│   ├── server
│   │   ├── order.go
│   │   ├── product.go
│   │   └── user.go
│   └── service
│       ├── order.go
│       ├── order_interface.go
│       ├── product.go
│       ├── product_interface.go
│       ├── user.go
│       └── user_interface.go
└── README.md
```
