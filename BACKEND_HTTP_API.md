# Backend HTTP API Documentation

The backend now supports both gRPC (for load manager) and HTTP (for direct testing and comparison).

## Server Ports

- **gRPC Server**: Port specified via `--port` flag (e.g., 50051)
- **HTTP Server**: Port 9000 (fixed)

## Starting the Backend

```bash
cd backend
go run cmd/backend/main.go --host localhost --port 50051
```

This will start both:
- gRPC server on localhost:50051
- HTTP server on localhost:9000

## HTTP Endpoints

### Users

#### Create User
```bash
POST http://localhost:9000/users
Content-Type: application/json

{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "password123"
}

Response: 201 Created
{
  "user_id": 1,
  "name": "John Doe",
  "email": "john@example.com",
  "password": "salt:hashed_password"
}
```

#### Get User(s)
```bash
# Get specific user
GET http://localhost:9000/users?email=john@example.com

Response: 200 OK
{
  "user_id": 1,
  "name": "John Doe",
  "email": "john@example.com"
}

# List all users
GET http://localhost:9000/users

Response: 200 OK
[
  {
    "user_id": 1,
    "name": "John Doe",
    "email": "john@example.com"
  },
  ...
]
```

#### Update User
```bash
PUT http://localhost:9000/users
Content-Type: application/json

{
  "email": "john@example.com",
  "name": "John Smith",
  "password": "newpassword123"
}

Response: 200 OK
{
  "message": "user updated"
}
```

#### Delete User
```bash
DELETE http://localhost:9000/users?email=john@example.com

Response: 200 OK
{
  "message": "user deleted"
}
```

### Products

#### Create Product
```bash
POST http://localhost:9000/products
Content-Type: application/json

{
  "name": "Widget",
  "version": "1.0.0"
}

Response: 201 Created
{
  "product_id": 1,
  "name": "Widget",
  "version": "1.0.0",
  "created_at": "2024-01-01T00:00:00Z"
}
```

#### Get Product(s)
```bash
# Get specific product
GET http://localhost:9000/products?product_id=1

Response: 200 OK
{
  "product_id": 1,
  "name": "Widget",
  "version": "1.0.0",
  "created_at": "2024-01-01T00:00:00Z"
}

# List all products
GET http://localhost:9000/products

Response: 200 OK
[
  {
    "product_id": 1,
    "name": "Widget",
    "version": "1.0.0",
    "created_at": "2024-01-01T00:00:00Z"
  },
  ...
]
```

#### Update Product
```bash
PUT http://localhost:9000/products
Content-Type: application/json

{
  "product_id": 1,
  "name": "Widget Pro",
  "version": "2.0.0"
}

Response: 200 OK
{
  "message": "product updated"
}
```

#### Delete Product
```bash
DELETE http://localhost:9000/products?product_id=1

Response: 200 OK
{
  "message": "product deleted"
}
```

### Orders

#### Create Order
```bash
POST http://localhost:9000/orders
Content-Type: application/json

{
  "user_id": 1,
  "product_id": 1,
  "quantity": 5
}

Response: 201 Created
{
  "order_id": 1,
  "user_id": 1,
  "product_id": 1,
  "quantity": 5,
  "created_at": "2024-01-01T00:00:00Z"
}
```

#### Get Order(s)
```bash
# Get specific order for user
GET http://localhost:9000/orders?user_id=1&order_id=1

Response: 200 OK
{
  "order_id": 1,
  "user_id": 1,
  "product_id": 1,
  "quantity": 5,
  "created_at": "2024-01-01T00:00:00Z"
}

# List all orders for user
GET http://localhost:9000/orders?user_id=1

Response: 200 OK
[
  {
    "order_id": 1,
    "user_id": 1,
    "product_id": 1,
    "quantity": 5,
    "created_at": "2024-01-01T00:00:00Z"
  },
  ...
]
```

#### Update Order
```bash
PUT http://localhost:9000/orders
Content-Type: application/json

{
  "order_id": 1,
  "user_id": 1,
  "product_id": 1,
  "quantity": 10
}

Response: 200 OK
{
  "message": "order updated"
}
```

#### Delete Order
```bash
DELETE http://localhost:9000/orders?order_id=1

Response: 200 OK
{
  "message": "order deleted"
}
```

## Python Testing Example

```python
import requests
import time

# Backend HTTP endpoint (direct access)
BACKEND_URL = "http://localhost:9000"

# Load Manager endpoint (with batching/queueing)
LOAD_MANAGER_URL = "http://localhost:8000"

def test_create_user_backend():
    """Test direct backend access"""
    start = time.time()
    response = requests.post(f"{BACKEND_URL}/users", json={
        "name": "Test User",
        "email": "test@example.com",
        "password": "password123"
    })
    latency = time.time() - start
    print(f"Backend latency: {latency*1000:.2f}ms")
    return response.json()

def test_create_user_load_manager():
    """Test through load manager (batched)"""
    start = time.time()
    response = requests.post(f"{LOAD_MANAGER_URL}/users", json={
        "name": "Test User",
        "email": "test@example.com",
        "password": "password123"
    })
    latency = time.time() - start
    print(f"Load Manager latency: {latency*1000:.2f}ms")
    return response.status_code  # Should be 202 Accepted

# Compare backend vs load manager
user = test_create_user_backend()  # Fast, synchronous
status = test_create_user_load_manager()  # Returns 202, processes in batch
```

## Key Differences

### Backend HTTP (Port 9000)
- **Synchronous**: Returns data immediately
- **No batching**: Each request processed individually
- **Returns data**: Create operations return the created object with ID
- **Use for**: Baseline comparison, measuring load manager overhead

### Load Manager HTTP (Port 8000)
- **Asynchronous writes**: Returns 202 Accepted immediately
- **Batching**: Accumulates requests, sends in batches to backend
- **No return data**: Create/update/delete return 202 status only
- **Use for**: Testing load balancing strategies, measuring throughput

## Next Steps

1. Start backend: `go run cmd/backend/main.go --host localhost --port 50051`
2. Start load manager: See load-manager README
3. Write Python tests to compare:
   - Direct backend access (baseline)
   - Load manager with different strategies (M, PR, PO, PRO)
   - Measure p50, p95, p99 latencies
   - Generate comparison graphs
