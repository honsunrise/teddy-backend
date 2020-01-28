package repositories

import (
	"context"
	"teddy-backend/internal/models"
	"teddy-backend/internal/proto/uaa"
	"upper.io/db.v3"
)

type AccountRepository interface {
	InsertAccount(account *models.Account) error
	FindOne(principal string) (*models.Account, error)
	FindAccountByUsername(username string) (*models.Account, error)
	FindAccountByEmail(email string) (*models.Account, error)
	FindAccountByPhone(phone string) (*models.Account, error)
	FindAll(page, size uint, sorts []*uaa.Sort) ([]*models.Account, error)
	DeleteOne(uid string) error
	UpdateOne(uid string, account map[string]interface{}) error
}

func NewAccountRepository(sess db.Database) (AccountRepository, error) {
	return &accountRepository{
		ctx:        context.Background(),
		sess:       sess,
		collection: sess.Collection("account"),
	}, nil
}

type accountRepository struct {
	ctx        context.Context
	sess       db.Database
	collection db.Collection
}

func (repo *accountRepository) InsertAccount(account *models.Account) error {
	err := repo.collection.InsertReturning(account)
	if err != nil {
		return err
	}
	return nil
}

func (repo *accountRepository) FindOne(principal string) (*models.Account, error) {
	var account models.Account
	result := repo.collection.Find(db.Or(
		db.Cond{"id": principal},
		db.Cond{"username": principal},
		db.Cond{"email": principal},
		db.Cond{"phone": principal},
	))
	err := result.One(&account)
	if err != nil {
		return nil, err
	}
	return &account, nil
}

func (repo *accountRepository) FindAccountByUsername(username string) (*models.Account, error) {
	var account models.Account
	result := repo.collection.Find(
		db.Cond{"username": username},
	)
	err := result.One(&account)
	if err != nil {
		return nil, err
	}
	return &account, nil
}

func (repo *accountRepository) FindAccountByEmail(email string) (*models.Account, error) {
	var account models.Account
	result := repo.collection.Find(
		db.Cond{"email": email},
	)
	err := result.One(&account)
	if err != nil {
		return nil, err
	}
	return &account, nil
}

func (repo *accountRepository) FindAccountByPhone(phone string) (*models.Account, error) {
	var account models.Account
	result := repo.collection.Find(
		db.Cond{"phone": phone},
	)
	err := result.One(&account)
	if err != nil {
		return nil, err
	}
	return &account, nil
}

func (repo *accountRepository) FindAll(page, size uint, sorts []*uaa.Sort) ([]*models.Account, error) {
	var itemsSorts = make([]interface{}, 0, len(sorts))
	if len(sorts) != 0 {
		for _, sort := range sorts {
			if sort.Asc {
				itemsSorts = append(itemsSorts, sort.Name)
			} else {
				itemsSorts = append(itemsSorts, "-"+sort.Name)
			}
		}
	}

	accounts := make([]*models.Account, 0, size)
	err := repo.collection.Find().OrderBy(itemsSorts...).Paginate(size).Page(page).All(accounts)
	if err != nil {
		return nil, err
	}
	return accounts, nil
}

func (repo *accountRepository) DeleteOne(uid string) error {
	err := repo.collection.Find(db.Cond{"id": uid}).Delete()
	if err != nil {
		return nil
	}
	return err
}

func (repo *accountRepository) UpdateOne(uid string, fields map[string]interface{}) error {
	err := repo.collection.Find(db.Cond{"id": uid}).Update(fields)
	if err != nil {
		return err
	}
	return nil
}
