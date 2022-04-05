package service

import (
	"github.com/kirktriplefive/test"
	"github.com/kirktriplefive/test/pkg/repository"
)

type ItemService struct {
	repo repository.Item
}

func NewItemService(repo repository.Item) *ItemService {
	return &ItemService{repo: repo}
}

func (s *ItemService) Create(item test.Item) (string, error) {
	return s.repo.Create(item)	
}

func (s *ItemService) GetAll() ([]test.Item, error) {
	return s.repo.GetAll()
}
