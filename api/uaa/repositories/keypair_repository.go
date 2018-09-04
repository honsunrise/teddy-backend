package repositories

import (
	"context"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/zhsyourai/teddy-backend/common/models"
	"time"
)

type KeyValuePairRepository interface {
	InsertKeyValuePair(kv *models.KeyValuePair) error
	FindKeyValuePairByKey(key string) (models.KeyValuePair, error)
	FindKeyValuePairByKeyAndValueAndExpire(key string, value string, time time.Time) (models.KeyValuePair, error)
	DeleteKeyValuePairByKey(key string) error
	DeleteKeyValuePairLT(time time.Time) error
}

func NewKeyValuePairRepository(client *mongo.Client) (KeyValuePairRepository, error) {
	return &keyValuePairRepository{
		ctx:         context.Background(),
		client:      client,
		collections: client.Database("Teddy").Collection("Account"),
	}, nil
}

type keyValuePairRepository struct {
	ctx         context.Context
	client      *mongo.Client
	collections *mongo.Collection
}

func (repo *keyValuePairRepository) FindKeyValuePairByKeyAndValueAndExpire(key string,
	value string, time time.Time) (models.KeyValuePair, error) {
	var kvp models.KeyValuePair
	filter := bson.NewDocument(
		bson.EC.String("key", key),
		bson.EC.String("value", value),
		bson.EC.SubDocumentFromElements("expire_time",
			bson.EC.Int64("$lt", time.UnixNano()),
		),
	)
	err := repo.collections.FindOne(repo.ctx, filter).Decode(&kvp)
	if err != nil {
		return models.KeyValuePair{}, err
	}
	return kvp, nil
}

func (repo *keyValuePairRepository) DeleteKeyValuePairLT(time time.Time) error {
	filter := bson.NewDocument(
		bson.EC.SubDocumentFromElements("expire_time",
			bson.EC.Int64("$lt", time.UnixNano()),
		),
	)
	_, err := repo.collections.DeleteMany(repo.ctx, filter)
	if err != nil {
		return err
	}
	return nil
}

func (repo *keyValuePairRepository) InsertKeyValuePair(kv *models.KeyValuePair) error {
	_, err := repo.collections.InsertOne(repo.ctx, kv)
	if err != nil {
		return err
	}
	return nil
}

func (repo *keyValuePairRepository) FindKeyValuePairByKey(key string) (models.KeyValuePair, error) {
	var kvp models.KeyValuePair
	filter := bson.NewDocument(bson.EC.String("key", key))
	err := repo.collections.FindOne(repo.ctx, filter).Decode(&kvp)
	if err != nil {
		return models.KeyValuePair{}, err
	}
	return kvp, nil
}

func (repo *keyValuePairRepository) DeleteKeyValuePairByKey(key string) error {
	filter := bson.NewDocument(bson.EC.String("key", key))
	_, err := repo.collections.DeleteOne(repo.ctx, filter)
	if err != nil {
		return err
	}
	return nil
}
