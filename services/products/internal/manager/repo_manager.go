package manager

import (
	"sync"

	"shop-server/infra"
	"shop-server/internal/repo"
)

type RepoManager interface {
	ProductRepo() repo.ProductRepo
}

type repoManager struct {
	infra infra.Infra
}

func NewRepoManager(infra infra.Infra) RepoManager {
	return &repoManager{infra: infra}
}

var (
	productRepoOnce sync.Once
	productRepo     repo.ProductRepo
)

func (rm *repoManager) ProductRepo() repo.ProductRepo {
	productRepoOnce.Do(func() {
		productRepo = repo.NewProductRepo(rm.infra.GormDB())
	})

	return productRepo
}
