package model

import "time"

type Order struct {
    OrderId     int       // corresponds to order_id
    UserId      int       // foreign key to users
    ProductId   int 
    Quantity    int
    CreatedAt   time.Time
}
