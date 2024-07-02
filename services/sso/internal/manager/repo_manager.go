package manager

import (
	"sync"

	"shop-server/infra"
)

type RepoManager interface {
}

type repoManager struct {
	infra infra.Infra
}

func NewRepoManager(infra infra.Infra) RepoManager {
	return &repoManager{infra: infra}
}

var (
	orderRepoOnce sync.Once
	orderRepo     repo.OrderRepo
)

func (rm *repoManager) OrderRepo() repo.OrderRepo {
	orderRepoOnce.Do(func() {
		orderRepo = repo.NewOrderRepo(rm.infra.GormDB())
	})
	return orderRepo
}
