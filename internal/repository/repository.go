package repository

import (
	"encoding/csv"
	"fmt"
	errs "golang-arch/internal/errors"
	"golang-arch/internal/model"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"sync"
)

type partRepository struct {
	mu      sync.Mutex
	storage map[int64]model.Part
	nextID  int64
}

func NewPartRepository() *partRepository {
	return &partRepository{
		storage: make(map[int64]model.Part),
		nextID:  1,
	}
}

func (r *partRepository) LoadFromCSV(path string) error {
	cleanPath := filepath.Clean(path)
	file, err := os.Open(cleanPath)
	if err != nil {
		return fmt.Errorf("не удалось открыть файл: %w", err)
	}
	defer func() {
		err = file.Close()
		if err != nil {
			log.Printf("Ошибка закрытия файла: %v", err)
		}
	}()

	reader := csv.NewReader(file)

	// Пропускаем заголовок
	_, err = reader.Read()
	if err != nil {
		return fmt.Errorf("не удалось прочитать заголовок: %w", err)
	}

	var (
		record   []string
		id       int64
		quantity int
		weight   float64
	)

	for {
		record, err = reader.Read()
		if err != nil {
			break
		}

		id, err = strconv.ParseInt(record[0], 10, 64)
		if err != nil {
			log.Printf("Ошибка парсинга ID '%s': %v", record[0], err)
			continue
		}
		quantity, err = strconv.Atoi(record[3])
		if err != nil {
			log.Printf("Ошибка парсинга количества '%s': %v", record[3], err)
			continue
		}
		weight, err = strconv.ParseFloat(record[4], 64)
		if err != nil {
			log.Printf("Ошибка парсинга веса '%s': %v", record[4], err)
			continue
		}

		r.storage[id] = model.Part{
			ID:       id,
			Name:     record[1],
			Type:     record[2],
			Quantity: quantity,
			Weight:   weight,
		}

		if id >= r.nextID {
			r.nextID = id + 1
		}
	}

	return nil
}

func (r *partRepository) GetAll() []model.Part {
	r.mu.Lock()
	defer r.mu.Unlock()

	parts := make([]model.Part, 0, len(r.storage))
	for _, p := range r.storage {
		parts = append(parts, p)
	}

	return parts
}

func (r *partRepository) Create(part model.Part) model.Part {
	r.mu.Lock()
	defer r.mu.Unlock()

	part.ID = r.nextID
	r.storage[r.nextID] = part
	r.nextID++
	return part
}

func (r *partRepository) Withdraw(id int64, quantity int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	part, exists := r.storage[id]
	if !exists {
		return errs.ErrNotFound
	}

	part.Quantity -= quantity
	r.storage[id] = part
	return nil
}

func (r *partRepository) GetByID(id int64) (model.Part, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	part, exists := r.storage[id]
	if !exists {
		return model.Part{}, errs.ErrNotFound
	}

	return part, nil
}
