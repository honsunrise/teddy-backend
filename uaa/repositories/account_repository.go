package repositories

import (
	"context"
	"errors"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/zhsyourai/teddy-backend/common/models"
)

var ErrUpdateAccount = errors.New("uaa update error")

type AccountRepository interface {
	InsertAccount(account *models.Account) error
	FindAccountByUsername(username string) (models.Account, error)
	FindAll() ([]models.Account, error)
	DeleteAccountByUsername(username string) error
	UpdateAccountByUsername(username string, account map[string]interface{}) error
}

func NewAccountRepository(client *mongo.Client) (AccountRepository, error) {
	return &accountRepository{
		ctx:         context.Background(),
		client:      client,
		collections: client.Database("Teddy").Collection("Account"),
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

func (repo *accountRepository) FindAccountByUsername(username string) (account models.Account, err error) {
	filter := bson.NewDocument(bson.EC.String("username", username))
	err = repo.collections.FindOne(repo.ctx, filter).Decode(&account)
	if err != nil {
		return
	}
	return
}

func (repo *accountRepository) FindAll() (accounts []models.Account, err error) {
	var cur mongo.Cursor
	cur, err = repo.collections.Find(repo.ctx, nil)
	if err != nil {
		return
	}
	defer cur.Close(repo.ctx)
	for cur.Next(repo.ctx) {
		var elem models.Account
		err = cur.Decode(&elem)
		if err != nil {
			return
		}
		accounts = append(accounts, elem)
	}
	err = cur.Err()
	return
}

func (repo *accountRepository) DeleteAccountByUsername(username string) (err error) {
	filter := bson.NewDocument(bson.EC.String("username", username))
	_, err = repo.collections.DeleteOne(repo.ctx, filter)
	if err != nil {
		return
	}
	return
}

func (repo *accountRepository) UpdateAccountByUsername(username string,
	fields map[string]interface{}) error {
	filter := bson.NewDocument(
		bson.EC.String("username", username),
	)

	var bsonFields []*bson.Element
	for k, v := range fields {
		bsonFields = append(bsonFields, bson.EC.Interface(k, v))
	}
	update := bson.NewDocument(
		bson.EC.SubDocumentFromElements("$set", bsonFields...),
	)
	ur, err := repo.collections.UpdateOne(repo.ctx, filter, update)
	if err != nil {
		return err
	} else if ur.ModifiedCount == 0 {
		return ErrUpdateAccount
	}
	return nil
}
