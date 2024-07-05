package service

import "shop-server/internal/repo"

type ProductServcie interface {
}

type productService struct {
	repository repo.ProductRepo
}

func NewProductService(repo repo.ProductRepo) ProductServcie {
	return &productService{repository: repo}
}
