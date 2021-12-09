package service

import (
	"context"
)

type ProdService struct {
}

func (ps *ProdService) GetProdStock(ctx context.Context, in *ProductRequest) (*ProductResponse, error) {
	start := in.ProdId + 100
	return &ProductResponse{ProdStock: start}, nil
}
