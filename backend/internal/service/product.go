package service

import (
	"context"
	"errors"
	"runtime"
	"sync"
	"github.com/sudo-JP/Load-Manager/backend/internal/model"
	"github.com/sudo-JP/Load-Manager/backend/internal/repository"
)

type ProductServer struct {
	repo repository.ProductRepositoryInterface
}
