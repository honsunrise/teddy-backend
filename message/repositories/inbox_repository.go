package repositories

import (
	"context"
	"fmt"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/zhsyourai/teddy-backend/common/types"
	"github.com/zhsyourai/teddy-backend/message/models"
)

type InBoxRepository interface {
	InsertInBoxItem(uid string, item *models.InBoxItem) error
	FindInBoxItems(uid string, itemType models.InBoxType, page uint32, size uint32, sorts []types.Sort) ([]models.InBoxItem, error)
	FindInBoxItem(uid string, id string) (models.InBoxItem, error)
	FindInBoxUnreadCount(uid string) (int64, error)
	DeleteAllInBoxItem(uid string) error
	DeleteInBoxItems(uid string, ids []string) error
	UpdateInBoxItems(uid string, ids []string, fields map[string]interface{}) error
}

func NewInBoxRepository(client *mongo.Client) (InBoxRepository, error) {
	return &inboxRepository{
		ctx:         context.Background(),
		client:      client,
		collections: client.Database("teddy").Collection("inbox"),
	}, nil
}

type inboxRepository struct {
	ctx         context.Context
	client      *mongo.Client
	collections *mongo.Collection
}

func (repo *inboxRepository) InsertInBoxItem(uid string, item *models.InBoxItem) error {
	filter := bson.D{{"uid", uid}}
	update := bson.D{{"$addToSet", bson.D{{"items", item}}}}
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

func (repo *inboxRepository) internalFindInBoxItems(uid string, itemType models.InBoxType, ids []string, page uint32,
	size uint32, sorts []types.Sort) ([]models.InBoxItem, error) {
	var dynFilter = make(bson.D, 0, 2)
	if itemType != models.ALL {
		dynFilter = append(dynFilter, bson.E{Key: "items.type", Value: int64(itemType)})
	}
	if len(ids) != 0 {
		dynFilter = append(dynFilter, bson.E{Key: "items.id", Value: bson.D{{"$in", bson.A{ids}}}})
	}

	var itemsSorts = make(bson.D, 0, len(sorts))
	if len(sorts) != 0 {
		for _, sort := range sorts {
			if sort.Order == types.ASC {
				itemsSorts = append(itemsSorts, bson.E{Key: sort.Name, Value: 1})
			} else {
				itemsSorts = append(itemsSorts, bson.E{Key: sort.Name, Value: -1})
			}
		}
	}
	pipeline := bson.D{
		{"$match", bson.D{{"uid", uid}}},
		{"$unwind", "$items"},
		{"$match", dynFilter},
		{"$skip", int64(size * page)},
		{"$limit", int64(size)},
		{"$sort", itemsSorts},
	}

	cur, err := repo.collections.Aggregate(repo.ctx, pipeline)
	if err != nil {
		return nil, err
	}
	items := make([]models.InBoxItem, 0, 50)
	defer cur.Close(repo.ctx)
	for cur.Next(repo.ctx) {
		var item models.InBoxItem
		err := cur.Decode(&item)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	err = cur.Err()
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (repo *inboxRepository) FindInBoxItems(uid string, itemType models.InBoxType, page uint32,
	size uint32, sorts []types.Sort) ([]models.InBoxItem, error) {
	return repo.internalFindInBoxItems(uid, itemType, nil, page, size, sorts)
}

func (repo *inboxRepository) FindInBoxItem(uid string, id string) (models.InBoxItem, error) {
	result, err := repo.internalFindInBoxItems(uid, models.ALL, []string{id}, 0, 1, nil)
	if err != nil {
		return models.InBoxItem{}, nil
	}
	return result[0], nil
}

func (repo *inboxRepository) FindInBoxUnreadCount(uid string) (int64, error) {
	var cur mongo.Cursor
	pipeline := bson.D{
		{"$match", bson.D{{"uid", uid}}},
		{"$unwind", "$items"},
		{"$match", bson.D{{"items.unread", true}}},
		{"$count", "count"},
	}

	repo.collections.Aggregate(repo.ctx, pipeline)
	elem := make(map[string]interface{})
	defer cur.Close(repo.ctx)
	for cur.Next(repo.ctx) {
		err := cur.Decode(&elem)
		if err != nil {
			return 0, err
		}
		break
	}
	err := cur.Err()
	if err != nil {
		return 0, err
	}
	return elem["count"].(int64), nil
}

func (repo *inboxRepository) DeleteAllInBoxItem(uid string) error {
	filter := bson.D{{"uid", uid}}
	update := bson.D{{"$unset", bson.D{{"items", ""}}}}
	_, err := repo.collections.UpdateOne(repo.ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}

func (repo *inboxRepository) DeleteInBoxItems(uid string, ids []string) error {
	filter := bson.D{{"uid", uid}}
	update := bson.D{{
		"$pull", bson.D{{
			"items", bson.D{{
				"id", bson.D{{
					"$in", bson.A{ids},
				}},
			}},
		}},
	}}
	_, err := repo.collections.UpdateOne(repo.ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}

func (repo *inboxRepository) UpdateInBoxItems(uid string, ids []string, fields map[string]interface{}) error {
	filter := bson.D{
		{"uid", uid},
		{"items", bson.D{
			{"$elemMatch", bson.D{{"id", bson.A{ids}}}},
		}},
	}

	var bsonFields = make(bson.D, 0, len(fields))
	for k, v := range fields {
		bsonFields = append(bsonFields, bson.E{Key: fmt.Sprintf("items.%s", k), Value: v})
	}

	update := bson.D{{"$set", bsonFields}}
	ur, err := repo.collections.UpdateOne(repo.ctx, filter, update)
	if err != nil {
		return err
	} else if ur.ModifiedCount == 0 {
		return err
	}
	return nil
}
