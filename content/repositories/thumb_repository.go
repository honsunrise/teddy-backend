package repositories

import (
	"context"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/objectid"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/zhsyourai/teddy-backend/common/types"
	"github.com/zhsyourai/teddy-backend/content/models"
	"time"
)

type ThumbRepository interface {
	Insert(userID objectid.ObjectID, thumb *models.ThumbItem) error
	FindInfoByUser(userID objectid.ObjectID, page uint32, size uint32, sorts []types.Sort) ([]models.ThumbItem, error)
	FindUserByInfo(infoID objectid.ObjectID, page uint32, size uint32, sorts []types.Sort) ([]objectid.ObjectID, error)
	IsExist(userID, infoID objectid.ObjectID) (bool, error)
	CountByInfo(infoID objectid.ObjectID) (uint64, error)
	CountByUser(userID objectid.ObjectID) (uint64, error)
	Delete(userID, infoID objectid.ObjectID) error
}

func NewThumbUpRepository(client *mongo.Client) (ThumbRepository, error) {
	return &thumbRepository{
		ctx:         context.Background(),
		client:      client,
		collections: client.Database("teddy").Collection("thumb_up"),
	}, nil
}

func NewThumbDownRepository(client *mongo.Client) (ThumbRepository, error) {
	return &thumbRepository{
		ctx:         context.Background(),
		client:      client,
		collections: client.Database("teddy").Collection("thumb_down"),
	}, nil
}

type thumbRepository struct {
	ctx         context.Context
	client      *mongo.Client
	collections *mongo.Collection
}

func (repo *thumbRepository) Insert(userID objectid.ObjectID, thumb *models.ThumbItem) error {
	filter := bson.D{{"userId", userID}}
	result := repo.collections.FindOne(repo.ctx, filter)
	if result.Decode(nil) == mongo.ErrNoDocuments {
		now := time.Now()
		f := models.Thumb{
			Id:        objectid.New(),
			UserId:    userID,
			LastTime:  now,
			FirstTime: now,
			Count:     0,
			Items:     []models.ThumbItem{},
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

func (repo *thumbRepository) FindInfoByUser(userID objectid.ObjectID,
	page uint32, size uint32, sorts []types.Sort) ([]models.ThumbItem, error) {
	var cur mongo.Cursor
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
		{"$unwind", "$items"},
		{"$match", bson.D{{"userId", userID}}},
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
	items := make([]models.ThumbItem, 0, size)
	for cur.Next(repo.ctx) {
		var item models.ThumbItem
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

func (repo *thumbRepository) FindUserByInfo(infoID objectid.ObjectID,
	page uint32, size uint32, sorts []types.Sort) ([]objectid.ObjectID, error) {
	var cur mongo.Cursor
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
		{"$match", bson.D{{"items.infoId", infoID}}},
		{"$unwind", "$items"},
		{"$skip", int64(size * page)},
		{"$limit", int64(size)},
		{"$sort", itemsSorts},
		{"$project", bson.D{
			{"_id", 0},
			{"userId", 1},
		}},
	}
	cur, err := repo.collections.Aggregate(repo.ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cur.Close(repo.ctx)
	items := make([]objectid.ObjectID, 0, size)
	for cur.Next(repo.ctx) {
		item := make(map[interface{}]interface{})
		err := cur.Decode(&item)
		if err != nil {
			return nil, err
		}
		items = append(items, item["userId"].(objectid.ObjectID))
	}
	err = cur.Err()
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (repo *thumbRepository) IsExist(userID, infoID objectid.ObjectID) (bool, error) {
	filter := bson.D{{"userId", userID}, {"items.infoId", infoID}}
	result := repo.collections.FindOne(repo.ctx, filter)
	err := result.Decode(nil)
	if err == mongo.ErrNoDocuments {
		return false, nil
	} else if err == nil {
		return true, nil
	}
	return false, err
}

func (repo *thumbRepository) CountByInfo(infoID objectid.ObjectID) (uint64, error) {
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

func (repo *thumbRepository) CountByUser(userID objectid.ObjectID) (uint64, error) {
	pipeline := bson.D{
		{"$match", bson.D{{"userId", userID}}},
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

func (repo *thumbRepository) Delete(userID, infoID objectid.ObjectID) error {
	filter := bson.D{{"userId", userID}, {"items.infoId", infoID}}
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
