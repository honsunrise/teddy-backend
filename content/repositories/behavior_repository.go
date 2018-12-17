package repositories

import (
	"context"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/objectid"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/zhsyourai/teddy-backend/common/proto"
	"github.com/zhsyourai/teddy-backend/content/models"
	"time"
)

type BehaviorRepository interface {
	Insert(uid objectid.ObjectID, thumb *models.BehaviorInfoItem) error
	FindInfoByUser(uid objectid.ObjectID, page uint32, size uint32, sorts []*proto.Sort) ([]models.BehaviorInfoItem, error)
	FindUserByInfo(infoID objectid.ObjectID, page uint32, size uint32, sorts []*proto.Sort) ([]models.BehaviorUserItem, error)
	IsExist(uid, infoID objectid.ObjectID) (bool, error)
	CountByInfo(infoID objectid.ObjectID) (uint64, error)
	CountByUser(uid objectid.ObjectID) (uint64, error)
	Delete(uid, infoID objectid.ObjectID) error
}

func NewThumbUpRepository(client *mongo.Client) (BehaviorRepository, error) {
	return &behaviorRepository{
		ctx:         context.Background(),
		client:      client,
		collections: client.Database("teddy").Collection("thumb_up"),
	}, nil
}

func NewThumbDownRepository(client *mongo.Client) (BehaviorRepository, error) {
	return &behaviorRepository{
		ctx:         context.Background(),
		client:      client,
		collections: client.Database("teddy").Collection("thumb_down"),
	}, nil
}

func NewFavoriteRepository(client *mongo.Client) (BehaviorRepository, error) {
	return &behaviorRepository{
		ctx:         context.Background(),
		client:      client,
		collections: client.Database("teddy").Collection("favorite"),
	}, nil
}

type behaviorRepository struct {
	ctx         context.Context
	client      *mongo.Client
	collections *mongo.Collection
}

func (repo *behaviorRepository) Insert(uid objectid.ObjectID, thumb *models.BehaviorInfoItem) error {
	filter := bson.D{{"uid", uid}}
	result := repo.collections.FindOne(repo.ctx, filter)
	if result.Decode(nil) == mongo.ErrNoDocuments {
		now := time.Now()
		f := models.Behavior{
			Id:        objectid.New(),
			UID:       uid,
			LastTime:  now,
			FirstTime: now,
			Count:     0,
			Items:     []models.BehaviorInfoItem{},
		}
		_, err := repo.collections.InsertOne(repo.ctx, f)
		if err != nil {
			return err
		}
	} else {
		update := bson.D{
			{"$addToSet", bson.D{{"items", bson.D{{"$each", bson.A{thumb}}}}}},
			{"$inc", bson.D{{"count", 1}}},
			{"$currentDate", bson.D{{"lastTime", bson.D{{"$type", "timestamp"}}}}},
		}
		ur, err := repo.collections.UpdateOne(repo.ctx, filter, update)
		if err != nil {
			return err
		} else if ur.ModifiedCount == 0 {
			return ErrUpdateInfo
		}
		return nil
	}
	return nil
}

func (repo *behaviorRepository) FindInfoByUser(uid objectid.ObjectID,
	page uint32, size uint32, sorts []*proto.Sort) ([]models.BehaviorInfoItem, error) {
	var cur mongo.Cursor
	var itemsSorts = make(bson.D, 0, len(sorts))
	if len(sorts) != 0 {
		for _, sort := range sorts {
			if sort.Asc {
				itemsSorts = append(itemsSorts, bson.E{Key: sort.Name, Value: 1})
			} else {
				itemsSorts = append(itemsSorts, bson.E{Key: sort.Name, Value: -1})
			}
		}
	}
	pipeline := bson.D{
		{"$unwind", "$items"},
		{"$match", bson.D{{"uid", uid}}},
		{"$skip", int64(size * page)},
		{"$limit", int64(size)},
		{"$sort", itemsSorts},
		{"$project", bson.D{
			{"_id", 0},
			{"infoId", "$items.infoId"},
			{"time", "$items.time"},
		}},
	}
	cur, err := repo.collections.Aggregate(repo.ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cur.Close(repo.ctx)
	items := make([]models.BehaviorInfoItem, 0, size)
	for cur.Next(repo.ctx) {
		var item models.BehaviorInfoItem
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

func (repo *behaviorRepository) FindUserByInfo(infoID objectid.ObjectID,
	page uint32, size uint32, sorts []*proto.Sort) ([]models.BehaviorUserItem, error) {
	var cur mongo.Cursor
	var itemsSorts = make(bson.D, 0, len(sorts))
	if len(sorts) != 0 {
		for _, sort := range sorts {
			if sort.Asc {
				itemsSorts = append(itemsSorts, bson.E{Key: sort.Name, Value: 1})
			} else {
				itemsSorts = append(itemsSorts, bson.E{Key: sort.Name, Value: -1})
			}
		}
	}
	pipeline := bson.D{
		{"$match", bson.D{{"items.infoId", infoID}}},
		{"$unwind", "$items"},
		{"$skip", int64(size * page)},
		{"$limit", int64(size)},
		{"$sort", itemsSorts},
		{"$project", bson.D{
			{"_id", 0},
			{"uid", 1},
		}},
	}
	cur, err := repo.collections.Aggregate(repo.ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cur.Close(repo.ctx)
	items := make([]models.BehaviorUserItem, 0, size)
	for cur.Next(repo.ctx) {
		var item models.BehaviorUserItem
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

func (repo *behaviorRepository) IsExist(uid, infoID objectid.ObjectID) (bool, error) {
	filter := bson.D{{"uid", uid}, {"items.infoId", infoID}}
	result := repo.collections.FindOne(repo.ctx, filter)
	err := result.Decode(nil)
	if err == mongo.ErrNoDocuments {
		return false, nil
	} else if err == nil {
		return true, nil
	}
	return false, err
}

func (repo *behaviorRepository) CountByInfo(infoID objectid.ObjectID) (uint64, error) {
	pipeline := bson.D{
		{"$unwind", "$items"},
		{"$match", bson.D{{"items.infoId", infoID}}},
		{"$group", bson.D{{"_id", nil}, {"count", bson.D{{"$sum", 1}}}}},
		{"$project", bson.D{{"_id", 0}}},
	}
	cur, err := repo.collections.Aggregate(repo.ctx, pipeline)
	if err != nil {
		return 0, err
	}
	defer cur.Close(repo.ctx)
	if cur.Next(repo.ctx) {
		item := make(map[interface{}]interface{})
		err := cur.Decode(&item)
		if err != nil {
			return 0, err
		}
		return uint64(item["count"].(int64)), nil
	} else {
		err = cur.Err()
		if err != nil {
			return 0, err
		}
		return 0, ErrThumbCount
	}
}

func (repo *behaviorRepository) CountByUser(uid objectid.ObjectID) (uint64, error) {
	pipeline := bson.D{
		{"$match", bson.D{{"uid", uid}}},
		{"$unwind", "$items"},
		{"$group", bson.D{{"_id", nil}, {"count", bson.D{{"$sum", 1}}}}},
		{"$project", bson.D{{"_id", 0}}},
	}
	cur, err := repo.collections.Aggregate(repo.ctx, pipeline)
	if err != nil {
		return 0, err
	}
	defer cur.Close(repo.ctx)
	if cur.Next(repo.ctx) {
		item := make(map[interface{}]interface{})
		err := cur.Decode(&item)
		if err != nil {
			return 0, err
		}
		return uint64(item["count"].(int64)), nil
	} else {
		err = cur.Err()
		if err != nil {
			return 0, err
		}
		return 0, ErrThumbCount
	}
}

func (repo *behaviorRepository) Delete(uid, infoID objectid.ObjectID) error {
	filter := bson.D{{"uid", uid}, {"items.infoId", infoID}}
	result := repo.collections.FindOne(repo.ctx, filter)
	if result.Decode(nil) == mongo.ErrNoDocuments {
		return ErrThumbNotExisted
	} else {
		update := bson.D{
			{"$pull", bson.D{{"items", bson.D{{"infoId", infoID}}}}},
			{"$inc", bson.D{{"count", -1}}},
		}
		ur, err := repo.collections.UpdateOne(repo.ctx, filter, update)
		if err != nil {
			return err
		} else if ur.ModifiedCount == 0 {
			return ErrUpdateThumb
		}
		return nil
	}
}
