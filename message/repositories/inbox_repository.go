package repositories

import (
	"context"
	"fmt"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/zhsyourai/teddy-backend/common/models"
	"github.com/zhsyourai/teddy-backend/common/types"
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
		collections: client.Database("Teddy").Collection("InBox"),
	}, nil
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
	)
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
	var cur mongo.Cursor
	var itemsTypeFilter *bson.Element = nil
	if itemType != models.ALL {
		itemsTypeFilter = bson.EC.Int64("items.type", int64(itemType))
	}
	var itemsIdsFilter *bson.Element = nil
	if len(ids) != 0 {
		var bsonIds = make([]*bson.Value, 0, len(ids))
		for _, id := range ids {
			bsonIds = append(bsonIds, bson.VC.String(id))
		}
		itemsIdsFilter = bson.EC.SubDocumentFromElements("items.id",
			bson.EC.ArrayFromElements("$in", bsonIds...),
		)
	}
	var dynFilter = make([]*bson.Element, 0, 2)
	if itemsTypeFilter != nil {
		dynFilter = append(dynFilter, itemsTypeFilter)
	}
	if itemsIdsFilter != nil {
		dynFilter = append(dynFilter, itemsIdsFilter)
	}
	var itemsSorts []*bson.Element
	if len(sorts) != 0 {
		itemsSorts = make([]*bson.Element, 0, len(sorts))
		for _, sort := range sorts {
			if sort.Order == types.ASC {
				itemsSorts = append(itemsSorts, bson.EC.Int32(sort.Name, 1))
			} else {
				itemsSorts = append(itemsSorts, bson.EC.Int32(sort.Name, -1))
			}
		}
	}
	pipeline := bson.NewDocument(
		bson.EC.SubDocumentFromElements("$match", bson.EC.String("uid", uid)),
		bson.EC.String("$unwind", "$items"),
		bson.EC.SubDocumentFromElements("$match", dynFilter...),
		bson.EC.String("$count", "count"),
		bson.EC.Int64("$skip", int64(size*page)),
		bson.EC.Int64("$limit", int64(size)),
		bson.EC.SubDocumentFromElements("$sort", itemsSorts...),
	)
	repo.collections.Aggregate(repo.ctx, pipeline)
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
	err := cur.Err()
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
	pipeline := bson.NewDocument(
		bson.EC.SubDocumentFromElements("$match", bson.EC.String("uid", uid)),
		bson.EC.String("$unwind", "$items"),
		bson.EC.SubDocumentFromElements("$match", bson.EC.Boolean("items.unread", true)),
		bson.EC.String("$count", "count"),
	)
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
	filter := bson.NewDocument(bson.EC.String("uid", uid))
	update := bson.NewDocument(
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
	var bsonIds = make([]*bson.Value, 0, len(ids))
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

func (repo *inboxRepository) UpdateInBoxItems(uid string, ids []string, fields map[string]interface{}) error {
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

	var bsonFields []*bson.Element
	for k, v := range fields {
		bsonFields = append(bsonFields, bson.EC.Interface(fmt.Sprintf("items.%s", k), v))
	}
	update := bson.NewDocument(
		bson.EC.SubDocumentFromElements("$set", bsonFields...),
	)
	_, err := repo.collections.UpdateOne(repo.ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}
