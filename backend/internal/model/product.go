package model

import "time"

type Product struct {
    ProductId int       
    Name      string   
    Version   string    
    CreatedAt time.Time 
}
