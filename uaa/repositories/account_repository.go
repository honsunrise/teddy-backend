package repositories

import (
	"context"
	"errors"
	"fmt"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/zhsyourai/teddy-backend/common/models"
	"reflect"
)

type AccountRepository interface {
	InsertAccount(account *models.Account) error
	FindAccountByUsername(username string) (models.Account, error)
	FindAll() ([]models.Account, error)
	DeleteAccountByUsername(username string) error
	UpdateAccountByUsername(username string, account map[string]interface{}) (models.Account, error)
}

func NewAccountRepository(client *mongo.Client) AccountRepository {
	return &accountMemoryRepository{
		ctx:         context.Background(),
		client:      client,
		collections: client.Database("Account").Collection("Account"),
	}
}

type accountMemoryRepository struct {
	ctx         context.Context
	client      *mongo.Client
	collections *mongo.Collection
}

func (repo *accountMemoryRepository) InsertAccount(account *models.Account) error {
	_, err := repo.collections.InsertOne(repo.ctx, account)
	if err != nil {
		return err
	}
	return nil
}

func (repo *accountMemoryRepository) FindAccountByUsername(username string) (account models.Account, err error) {
	filter := bson.NewDocument(bson.EC.String("username", username))
	err = repo.collections.FindOne(repo.ctx, filter).Decode(&account)
	if err != nil {
		return
	}
	return
}

func (repo *accountMemoryRepository) FindAll() (accounts []models.Account, err error) {
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

func (repo *accountMemoryRepository) DeleteAccountByUsername(username string) (err error) {
	filter := bson.NewDocument(bson.EC.String("username", username))
	_, err = repo.collections.DeleteOne(repo.ctx, filter)
	if err != nil {
		return
	}
	return
}

func (repo *accountMemoryRepository) UpdateAccountByUsername(username string,
	fields map[string]interface{}) (account models.Account, err error) {
	filter := bson.NewDocument(bson.EC.String("username", username))
	err = repo.collections.FindOne(repo.ctx, filter).Decode(&account)
	if err != nil {
		return
	}

	s := reflect.ValueOf(&account).Elem()
	for k, v := range fields {
		field := s.FieldByName(k)
		if field.IsValid() {
			field.Set(reflect.ValueOf(v))
		} else {
			err = errors.New(fmt.Sprintf("field %s not exist", k))
			return
		}
	}

	_, err = repo.collections.UpdateOne(repo.ctx, filter, account)
	if err != nil {
		return
	}
	return
}
