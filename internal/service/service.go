package service

import (
	errs "golang-arch/internal/errors"
	"golang-arch/internal/model"
)

type PartRepository interface {
	GetByID(id int64) (model.Part, error)
	Withdraw(id int64, quantity int) error
	GetAll() []model.Part
	Create(part model.Part) model.Part
}

type partService struct {
	repo PartRepository
}

func NewPartService(repo PartRepository) *partService {
	return &partService{repo: repo}
}

func (s *partService) GetAllParts() []model.Part { return s.repo.GetAll() }

func (s *partService) CreatePart(name, partType string, quantity int, weight float64) (model.Part, error) {
	part := model.Part{
		Name:     name,
		Type:     partType,
		Quantity: quantity,
		Weight:   weight,
	}

	part = s.repo.Create(part)
	return part, nil
}

func (s *partService) Withdraw(id int64, quantity int) error {
	part, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	if part.Quantity < quantity {
		return errs.ErrNotEnoughParts
	}
	return s.repo.Withdraw(id, quantity)
}

func (s *partService) GetByID(id int64) (model.Part, error) {
	return s.repo.GetByID(id)
}
