package manager

import (
	"sync"

	"shop-server/infra"
	"shop-server/internal/service"
)

type ServiceManager interface {
}

type serviceManager struct {
	infra infra.Infra
	repo  RepoManager
}

func NewServiceManager(infra infra.Infra) ServiceManager {
	return &serviceManager{
		infra: infra,
		repo:  NewRepoManager(infra),
	}
}

var (
	productServiceOnce sync.Once
	productService     service.ProductServcie
)

func (sm *serviceManager) ProductService() service.ProductServcie {
	productServiceOnce.Do(func() {
		productService = service.NewProductService(sm.repo.ProductRepo())
	})
	return productService
}
