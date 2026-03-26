package internal

type partService struct {
	repo *partRepository
}

func NewPartService(repo *partRepository) *partService {
	return &partService{repo: repo}
}

func (s *partService) GetAllParts() []Part {
	return s.repo.GetAll()
}

func (s *partService) CreatePart(name, partType string, quantity int, weight float64) (Part, error) {
	part := Part{
		Name:     name,
		Type:     partType,
		Quantity: quantity,
		Weight:   weight,
	}

	created := s.repo.Create(part)
	return created, nil
}

func (s *partService) Withdraw(id int64, quantity int) error {
	return s.repo.Withdraw(id, quantity)
}

func (s *partService) GetByID(id int64) (Part, error) {
	return s.repo.GetByID(id)
}
