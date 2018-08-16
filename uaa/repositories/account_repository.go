package repositories

import (
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/zhsyourai/teddy-backend/common/models"
)

type AccountRepository interface {
	InsertAccount(account *models.Account) error
	FindAccountByUsername(username string) (models.Account, error)
	FindAll() ([]models.Account, error)
	DeleteAccountByUsername(username string) (models.Account, error)
	UpdateAccountByUsername(username string, account map[string]interface{}) (models.Account, error)
}

func NewAccountRepository(client *mongo.Client) AccountRepository {
	return &accountMemoryRepository{client: client}
}

type accountMemoryRepository struct {
	client *mongo.Client
}

func (repo *accountMemoryRepository) InsertAccount(account *models.Account) error {
	panic("implement me")
}

func (repo *accountMemoryRepository) FindAccountByUsername(username string) (models.Account, error) {
	panic("implement me")
}

func (repo *accountMemoryRepository) FindAll() ([]models.Account, error) {
	panic("implement me")
}

func (repo *accountMemoryRepository) DeleteAccountByUsername(username string) (models.Account, error) {
	panic("implement me")
}

func (repo *accountMemoryRepository) UpdateAccountByUsername(username string, account map[string]interface{}) (models.Account, error) {
	panic("implement me")
}
