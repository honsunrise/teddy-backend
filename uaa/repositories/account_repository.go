package repositories

import (
	"context"
	"errors"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
	"teddy-backend/common/proto/uaa"
	"teddy-backend/uaa/models"
)

var ErrUpdateAccount = errors.New("uaa update error")

type AccountRepository interface {
	InsertAccount(account *models.Account) error
	FindOne(principal string) (*models.Account, error)
	FindAll(page, size uint32, sorts []*uaa.Sort) ([]*models.Account, error)
	DeleteOne(uid string) error
	UpdateOne(uid string, account map[string]interface{}) error
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

func (repo *accountRepository) FindOne(principal string) (*models.Account, error) {
	var account models.Account
	filter := bson.D{{"$or", bson.A{
		bson.D{{"_id", bson.D{
			{"$exists", true},
			{"$ne", ""},
			{"$eq", principal}}}},
		bson.D{{"username", bson.D{
			{"$exists", true},
			{"$ne", ""},
			{"$eq", principal}}}},
		bson.D{{"email", bson.D{
			{"$exists", true},
			{"$ne", ""},
			{"$eq", principal}}}},
		bson.D{{"phone", bson.D{
			{"$exists", true},
			{"$ne", ""},
			{"$eq", principal}}}},
	}}}
	err := repo.collections.FindOne(repo.ctx, filter).Decode(&account)
	if err != nil {
		return nil, err
	}
	return &account, nil
}

func (repo *accountRepository) FindAccountByEmail(email string) (*models.Account, error) {
	var account models.Account
	filter := bson.D{{"email", email}}
	err := repo.collections.FindOne(repo.ctx, filter).Decode(&account)
	if err != nil {
		return nil, err
	}
	return &account, nil
}

func (repo *accountRepository) FindAccountByPhone(phone string) (*models.Account, error) {
	var account models.Account
	filter := bson.D{{"phone", phone}}
	err := repo.collections.FindOne(repo.ctx, filter).Decode(&account)
	if err != nil {
		return nil, err
	}
	return &account, nil
}

func (repo *accountRepository) FindAll(page, size uint32, sorts []*uaa.Sort) ([]*models.Account, error) {
	pipeline := mongo.Pipeline{
		bson.D{{"$skip", int64(size * page)}},
		bson.D{{"$limit", int64(size)}},
	}

	var itemsSorts = make(bson.D, 0, len(sorts))
	if len(sorts) != 0 {
		for _, sort := range sorts {
			if sort.Asc {
				itemsSorts = append(itemsSorts, bson.E{Key: sort.Name, Value: 1})
			} else {
				itemsSorts = append(itemsSorts, bson.E{Key: sort.Name, Value: -1})
			}
		}
		pipeline = append(pipeline, bson.D{{"$sort", itemsSorts}})
	}

	cur, err := repo.collections.Aggregate(repo.ctx, pipeline)
	if err != nil {
		return nil, err
	}
	accounts := make([]*models.Account, 0, size)
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

func (repo *accountRepository) DeleteOne(uid string) (err error) {
	filter := bson.D{{"_id", uid}}
	_, err = repo.collections.DeleteOne(repo.ctx, filter)
	if err != nil {
		return
	}
	return
}

func (repo *accountRepository) UpdateOne(uid string,
	fields map[string]interface{}) error {
	filter := bson.D{{"_id", uid}}
	var bsonFields = make(bson.D, 0, len(fields))
	for k, v := range fields {
		bsonFields = append(bsonFields, bson.E{Key: k, Value: v})
	}
	update := bson.D{{"$set", bsonFields}}
	ur, err := repo.collections.UpdateOne(repo.ctx, filter, update)
	if err != nil {
		return err
	} else if ur.ModifiedCount == 0 {
		return ErrUpdateAccount
	}
	return nil
}
