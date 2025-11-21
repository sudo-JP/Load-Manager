package repository

import "github.com/sudo-JP/Load-Manager/backend/internal/model"

type OrderInterface interface {
    Create(order *model.Order) (bool, error)           
    GetById(orderId int) (*model.Order, error)
    GetByUser(userId int) ([]model.Order, error)
    Update(order model.Order) (bool, error)         
    Delete(orderId int) (bool, error)
    ListAll() ([]model.Order, error)               
}
