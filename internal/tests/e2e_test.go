package tests

import (
	"testing"

	"github.com/stretchr/testify/require"

	errs "golang-arch/internal/errors"
	"golang-arch/internal/model"
	"golang-arch/internal/service"
)

// mockPartRepo — ручной мок для тестирования
type mockPartRepo struct {
	part     model.Part
	parts    []model.Part
	err      error
	created  model.Part
	lastPart model.Part // последняя переданная деталь в Create
}

func (m *mockPartRepo) GetByID(id int64) (model.Part, error) {
	return m.part, m.err
}

func (m *mockPartRepo) Withdraw(id int64, quantity int) error {
	return nil
}

func (m *mockPartRepo) GetAll() []model.Part {
	return m.parts
}

func (m *mockPartRepo) Create(part model.Part) model.Part {
	m.lastPart = part
	if m.created.ID != 0 {
		return m.created
	}
	part.ID = 1
	return part
}

func TestWithdraw_NotEnoughParts(t *testing.T) {
	repo := &mockPartRepo{part: model.Part{ID: 1, Quantity: 5}}
	svc := service.NewPartService(repo)

	err := svc.Withdraw(1, 10)

	require.ErrorIs(t, err, errs.ErrNotEnoughParts)
}

func TestWithdraw_Success(t *testing.T) {
	repo := &mockPartRepo{part: model.Part{ID: 1, Quantity: 10}}
	svc := service.NewPartService(repo)

	err := svc.Withdraw(1, 5)

	require.NoError(t, err)
}

func TestWithdraw_NotFound(t *testing.T) {
	repo := &mockPartRepo{err: errs.ErrNotFound}
	svc := service.NewPartService(repo)

	err := svc.Withdraw(1, 5)

	require.ErrorIs(t, err, errs.ErrNotFound)
}

func TestGetAllParts_ReturnsEmptyList(t *testing.T) {
	repo := &mockPartRepo{parts: []model.Part{}}
	svc := service.NewPartService(repo)

	parts := svc.GetAllParts()

	require.Empty(t, parts)
}

func TestGetAllParts_ReturnsParts(t *testing.T) {
	expected := []model.Part{
		{ID: 1, Name: "Двигатель", Type: "engine", Quantity: 10, Weight: 100},
		{ID: 2, Name: "Корпус", Type: "hull", Quantity: 5, Weight: 500},
	}
	repo := &mockPartRepo{parts: expected}
	svc := service.NewPartService(repo)

	parts := svc.GetAllParts()

	require.Len(t, parts, 2)
	require.Equal(t, expected, parts)
}

func TestCreatePart_Success(t *testing.T) {
	repo := &mockPartRepo{}
	svc := service.NewPartService(repo)

	part, err := svc.CreatePart("Новый двигатель", "engine", 5, 150.5)

	require.NoError(t, err)
	require.Equal(t, int64(1), part.ID)
	require.Equal(t, "Новый двигатель", part.Name)
	require.Equal(t, "engine", part.Type)
	require.Equal(t, 5, part.Quantity)
	require.Equal(t, 150.5, part.Weight)
}

func TestCreatePart_PassesCorrectDataToRepo(t *testing.T) {
	repo := &mockPartRepo{
		created: model.Part{ID: 42, Name: "Титановая обшивка", Type: "hull", Quantity: 3, Weight: 50},
	}
	svc := service.NewPartService(repo)

	part, err := svc.CreatePart("Титановая обшивка", "hull", 3, 50)

	require.NoError(t, err)
	require.Equal(t, int64(42), part.ID)
	require.Equal(t, "Титановая обшивка", repo.lastPart.Name)
	require.Equal(t, "hull", repo.lastPart.Type)
	require.Equal(t, 3, repo.lastPart.Quantity)
	require.Equal(t, 50.0, repo.lastPart.Weight)
}
