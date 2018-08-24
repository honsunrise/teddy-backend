package repositories

import (
	"context"
	"errors"
	"fmt"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/findopt"
	"github.com/zhsyourai/teddy-backend/common/models"
	"reflect"
)

type InBoxRepository interface {
	InsertInBoxItem(uid string, item *models.InBoxItem) error
	FindInBoxItems(uid string, itemType models.InBoxType, page uint32, size uint32) ([]models.InBoxItem, error)
	FindInBoxItem(uid string, id string) (models.InBoxItem, error)
	DeleteAllInBoxItem(uid string) error
	DeleteInBoxItems(uid string, ids []string) error
	UpdateInBoxItem(uid string, id string, fields map[string]interface{}) error
	UpdateInBoxItems(uid string, ids []string, fields map[string]interface{}) error
}

func NewInBoxRepository(client *mongo.Client) InBoxRepository {
	return &inboxRepository{
		ctx:         context.Background(),
		client:      client,
		collections: client.Database("InBox").Collection("InBox"),
	}
}

type inboxRepository struct {
	ctx         context.Context
	client      *mongo.Client
	collections *mongo.Collection
}

func (repo *inboxRepository) InsertInBoxItem(uid string, item *models.InBoxItem) error {
	filter := bson.NewDocument(bson.EC.String("uid", uid))
	update := bson.NewDocument(
		bson.EC.SubDocumentFromElements("$addToSet", bson.EC.Interface("items", item)),
		bson.EC.SubDocumentFromElements("$inc", bson.EC.Int64("unread_count", 1)))
	result := repo.collections.FindOneAndUpdate(repo.ctx, filter, update)
	if result.Decode(nil) == mongo.ErrNoDocuments {
		inbox := models.InBox{
			Uid: uid,
			Items: []models.InBoxItem{
				*item,
			},
		}
		_, err := repo.collections.InsertOne(repo.ctx, &inbox)
		if err != nil {
			return err
		}
	}
	return nil
}

func (repo *inboxRepository) FindInBoxItems(uid string, itemType models.InBoxType, page uint32,
	size uint32) (items []models.InBoxItem, err error) {
	filter := bson.NewDocument(
		bson.EC.String("uid", uid),
		bson.EC.SubDocumentFromElements("items",
			bson.EC.ArrayFromElements("$slice", bson.VC.Int64(int64(page*size)), bson.VC.Int64(int64(size))),
		),
	)
	var inbox models.InBox
	err = repo.collections.FindOne(repo.ctx, filter).Decode(&inbox)
	if err != nil {
		return
	}
	items = append(items, inbox.Items...)
	return
}

func (repo *inboxRepository) FindInBoxItem(uid string, id string) (item models.InBoxItem, err error) {
	filter := bson.NewDocument(
		bson.EC.String("uid", uid),
		bson.EC.SubDocumentFromElements("items",
			bson.EC.SubDocumentFromElements("$elemMatch", bson.EC.String("id", id)),
		),
	)
	var inbox models.InBox
	err = repo.collections.FindOne(repo.ctx, filter, findopt.Projection(bson.NewDocument(
		bson.EC.Int32("items.$", 1),
		bson.EC.Int32("uid", 1),
	))).Decode(&inbox)
	if err != nil {
		return
	}
	item = inbox.Items[0]
	return
}

func (repo *inboxRepository) DeleteAllInBoxItem(uid string) error {
	filter := bson.NewDocument(bson.EC.String("uid", uid))
	update := bson.NewDocument(
		bson.EC.SubDocumentFromElements("$set", bson.EC.Int64("unread_count", 0)),
		bson.EC.SubDocumentFromElements("$unset", bson.EC.String("items", "")),
	)
	_, err := repo.collections.UpdateOne(repo.ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}

func (repo *inboxRepository) DeleteInBoxItems(uid string, ids []string) error {
	filter := bson.NewDocument(bson.EC.String("uid", uid))
	var bsonIds []*bson.Value
	for _, id := range ids {
		bsonIds = append(bsonIds, bson.VC.String(id))
	}
	update := bson.NewDocument(
		bson.EC.SubDocumentFromElements("$pull",
			bson.EC.SubDocumentFromElements("items",
				bson.EC.SubDocumentFromElements("id", bson.EC.ArrayFromElements("$in", bsonIds...)),
			),
		),
	)
	_, err := repo.collections.UpdateOne(repo.ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}

func (repo *inboxRepository) UpdateInBoxItem(uid string, id string, fields map[string]interface{}) error {
	var inbox models.InBox
	filter := bson.NewDocument(
		bson.EC.String("uid", uid),
		bson.EC.SubDocumentFromElements("items",
			bson.EC.SubDocumentFromElements("$elemMatch", bson.EC.String("id", id)),
		),
	)
	err := repo.collections.FindOne(repo.ctx, filter).Decode(&inbox)
	if err != nil {
		return err
	}

	s := reflect.ValueOf(&inbox.Items[0]).Elem()
	for k, v := range fields {
		field := s.FieldByName(k)
		if field.IsValid() {
			field.Set(reflect.ValueOf(v))
		} else {
			return errors.New(fmt.Sprintf("field %s not exist", k))
		}
	}

	update := bson.NewDocument(
		bson.EC.SubDocumentFromElements("$set",
			bson.EC.Interface("items.$", inbox.Items[0])),
	)
	_, err = repo.collections.UpdateOne(repo.ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}

func (repo *inboxRepository) UpdateInBoxItems(uid string, ids []string, fields map[string]interface{}) error {
	var inbox models.InBox
	var bsonIds []*bson.Value
	for _, id := range ids {
		bsonIds = append(bsonIds, bson.VC.String(id))
	}
	filter := bson.NewDocument(
		bson.EC.String("uid", uid),
		bson.EC.SubDocumentFromElements("items",
			bson.EC.SubDocumentFromElements("$elemMatch",
				bson.EC.SubDocumentFromElements("id", bson.EC.ArrayFromElements("$in", bsonIds...)),
			),
		),
	)
	err := repo.collections.FindOne(repo.ctx, filter).Decode(&inbox)
	if err != nil {
		return err
	}

	s := reflect.ValueOf(&inbox.Items[0]).Elem()
	for k, v := range fields {
		field := s.FieldByName(k)
		if field.IsValid() {
			field.Set(reflect.ValueOf(v))
		} else {
			return errors.New(fmt.Sprintf("field %s not exist", k))
		}
	}

	update := bson.NewDocument(
		bson.EC.SubDocumentFromElements("$set",
			bson.EC.Interface("items.$", inbox.Items[0])),
	)
	_, err = repo.collections.UpdateOne(repo.ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}
