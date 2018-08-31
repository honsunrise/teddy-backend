package components

import (
	"github.com/zhsyourai/teddy-backend/uaa/repositories"
	"math/rand"
)

type UidGenerator interface {
	NexID() (string, error)
}

func NewUidGenerator(repo repositories.AccountRepository) (UidGenerator, error) {
	return &uidGenerator{
		repo: repo,
	}, nil
}

type uidGenerator struct {
	repo repositories.AccountRepository
}

func (t *uidGenerator) NexID() (string, error) {
	s := ""
	for i := 0; i < 10; i++ {
		s += (string)(rand.Intn(10) + 48)
	}
	return s, nil
}
