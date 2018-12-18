package repositories

import (
	"context"
	"errors"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/zhsyourai/teddy-backend/uaa/models"
)

var ErrUpdateAccount = errors.New("uaa update error")

type AccountRepository interface {
	InsertAccount(account *models.Account) error
	FindAccountByUsername(username string) (*models.Account, error)
	FindAll() ([]*models.Account, error)
	DeleteAccountByUsername(username string) error
	UpdateAccountByUsername(username string, account map[string]interface{}) error
}

func NewAccountRepository(client *mongo.Client) (AccountRepository, error) {
	return &accountRepository{
		ctx:         context.Background(),
		client:      client,
		collections: client.Database("teddy").Collection("account"),
	}, nil
}

type accountRepository struct {
	ctx         context.Context
	client      *mongo.Client
	collections *mongo.Collection
}

func (repo *accountRepository) InsertAccount(account *models.Account) error {
	_, err := repo.collections.InsertOne(repo.ctx, account)
	if err != nil {
		return err
	}
	return nil
}

func (repo *accountRepository) FindAccountByUsername(username string) (*models.Account, error) {
	var account models.Account
	filter := bson.D{{"username", username}}
	err := repo.collections.FindOne(repo.ctx, filter).Decode(&account)
	if err != nil {
		return nil, err
	}
	return &account, nil
}

func (repo *accountRepository) FindAll() ([]*models.Account, error) {
	accounts := make([]*models.Account, 0, 100)
	var cur mongo.Cursor
	cur, err := repo.collections.Find(repo.ctx, nil)
	if err != nil {
		return nil, err
	}
	defer cur.Close(repo.ctx)
	for cur.Next(repo.ctx) {
		var elem models.Account
		err = cur.Decode(&elem)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, &elem)
	}
	err = cur.Err()
	return accounts, nil
}

func (repo *accountRepository) DeleteAccountByUsername(username string) (err error) {
	filter := bson.D{{"username", username}}
	_, err = repo.collections.DeleteOne(repo.ctx, filter)
	if err != nil {
		return
	}
	return
}

func (repo *accountRepository) UpdateAccountByUsername(username string,
	fields map[string]interface{}) error {
	filter := bson.D{{"username", username}}
	var bsonFields []bson.E
	for k, v := range fields {
		bsonFields = append(bsonFields, bson.E{Key: k, Value: v})
	}
	update := bson.D{{"$set", bson.A{bsonFields}}}
	ur, err := repo.collections.UpdateOne(repo.ctx, filter, update)
	if err != nil {
		return err
	} else if ur.ModifiedCount == 0 {
		return ErrUpdateAccount
	}
	return nil
}
