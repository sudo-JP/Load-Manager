# Backend Package 

## Run 
To run, ensure `.env` is in place, then 
```bash
go run cmd/backend/main.go --host GRPC_HOST --port GRPC_PORT
```
Note that the `GRPC_PORT` is used to create HTTP Port as well. The code specified the `GRPC_PORT` to be 50000 or higher, then `GRPC_PORT` - 50000 is the port for HTTP Server starting from 9000.
