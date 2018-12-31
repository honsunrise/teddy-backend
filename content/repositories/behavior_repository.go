package repositories

import (
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/objectid"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/options"
	"github.com/mongodb/mongo-go-driver/x/bsonx"
	"github.com/zhsyourai/teddy-backend/common/proto/content"
	"github.com/zhsyourai/teddy-backend/content/models"
	"time"
)

type BehaviorRepository interface {
	Insert(ctx mongo.SessionContext, uid string, thumb *models.BehaviorInfoItem) error
	FindInfoByUser(ctx mongo.SessionContext, uid string,
		page, size uint64, sorts []*content.Sort) ([]*models.BehaviorInfoItem, error)
	FindUserByInfo(ctx mongo.SessionContext, infoID objectid.ObjectID,
		page, size uint64, sorts []*content.Sort) ([]*models.BehaviorUserItem, error)
	IsExists(ctx mongo.SessionContext, uid string, infoID objectid.ObjectID) (bool, error)
	CountByInfo(ctx mongo.SessionContext, infoID objectid.ObjectID) (uint64, error)
	CountByUser(ctx mongo.SessionContext, uid string) (uint64, error)
	Delete(ctx mongo.SessionContext, uid string, infoID objectid.ObjectID) error
}

func NewThumbUpRepository(client *mongo.Client) (BehaviorRepository, error) {
	return &behaviorRepository{
		collections: client.Database("teddy").Collection("thumb_up"),
	}, nil
}

func NewThumbDownRepository(client *mongo.Client) (BehaviorRepository, error) {
	return &behaviorRepository{
		collections: client.Database("teddy").Collection("thumb_down"),
	}, nil
}

func NewFavoriteRepository(client *mongo.Client) (BehaviorRepository, error) {
	return &behaviorRepository{
		collections: client.Database("teddy").Collection("favorite"),
	}, nil
}

type behaviorRepository struct {
	collections *mongo.Collection
}

func (repo *behaviorRepository) Insert(ctx mongo.SessionContext, uid string, thumb *models.BehaviorInfoItem) error {
	result := repo.collections.FindOne(ctx, bson.D{{"uid", uid}})
	if err := result.Decode(nil); err == mongo.ErrNoDocuments {
		now := time.Now()
		f := models.Behavior{
			Id:        objectid.New(),
			UID:       uid,
			LastTime:  now,
			FirstTime: now,
			Count:     1,
			Items: []*models.BehaviorInfoItem{
				thumb,
			},
		}
		_, err := repo.collections.InsertOne(ctx, f)
		if err != nil {
			return err
		}
	} else if err == nil {
		filter := bson.D{{"uid", uid}, {"items", bson.D{{"$ne", thumb}}}}
		update := bson.D{
			{"$addToSet", bson.D{{"items", bson.D{{"$each", bson.A{thumb}}}}}},
			{"$inc", bson.D{{"count", 1}}},
			{"$currentDate", bson.D{{"lastTime", bson.D{{"$type", "timestamp"}}}}},
		}
		ur, err := repo.collections.UpdateOne(ctx, filter, update)
		if err != nil {
			return err
		} else if ur.ModifiedCount == 0 {
			return mongo.ErrNoDocuments
		}
		return nil
	}
	return nil
}

func (repo *behaviorRepository) FindInfoByUser(ctx mongo.SessionContext, uid string,
	page, size uint64, sorts []*content.Sort) ([]*models.BehaviorInfoItem, error) {
	var cur mongo.Cursor
	pipeline := mongo.Pipeline{
		bson.D{{"$unwind", "$items"}},
		bson.D{{"$match", bson.D{{"uid", uid}}}},
		bson.D{{"$skip", int64(size * page)}},
		bson.D{{"$limit", int64(size)}},
		bson.D{{"$project", bson.D{
			{"_id", 0},
			{"infoId", "$items.infoId"},
			{"time", "$items.time"},
		}}},
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

	cur, err := repo.collections.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	items := make([]*models.BehaviorInfoItem, 0, size)
	for cur.Next(ctx) {
		var item models.BehaviorInfoItem
		err := cur.Decode(&item)
		if err != nil {
			return nil, err
		}
		items = append(items, &item)
	}
	err = cur.Err()
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (repo *behaviorRepository) FindUserByInfo(ctx mongo.SessionContext, infoID objectid.ObjectID,
	page, size uint64, sorts []*content.Sort) ([]*models.BehaviorUserItem, error) {
	var cur mongo.Cursor

	pipeline := mongo.Pipeline{
		bson.D{{"$match", bson.D{{"items.infoId", infoID}}}},
		bson.D{{"$unwind", "$items"}},
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

	pipeline = append(pipeline, bson.D{{"$project", bson.D{
		{"_id", 0},
		{"uid", 1},
	}}})

	cur, err := repo.collections.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	items := make([]*models.BehaviorUserItem, 0, size)
	for cur.Next(ctx) {
		var item models.BehaviorUserItem
		err := cur.Decode(&item)
		if err != nil {
			return nil, err
		}
		items = append(items, &item)
	}
	err = cur.Err()
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (repo *behaviorRepository) IsExists(ctx mongo.SessionContext, uid string, infoID objectid.ObjectID) (bool, error) {
	filter := bson.D{{"uid", uid}, {"items.infoId", infoID}}
	result := repo.collections.FindOne(ctx, filter)
	err := result.Decode(nil)
	if err == mongo.ErrNoDocuments {
		return false, nil
	} else if err == nil {
		return true, nil
	}
	return false, err
}

func (repo *behaviorRepository) CountByInfo(ctx mongo.SessionContext, infoID objectid.ObjectID) (uint64, error) {
	pipeline := mongo.Pipeline{
		bson.D{{"$unwind", "$items"}},
		bson.D{{"$match", bson.D{{"items.infoId", infoID}}}},
		bson.D{{"$group", bson.D{
			{"_id", bsonx.Null()},
			{"count", bson.D{{"$sum", 1}}},
		}}},
		bson.D{{"$project", bson.D{{"_id", 0}}}},
	}
	cur, err := repo.collections.Aggregate(ctx, pipeline)
	if err != nil {
		return 0, err
	}
	defer cur.Close(ctx)
	if cur.Next(ctx) {
		item := make(map[string]interface{})
		err := cur.Decode(&item)
		if err != nil {
			return 0, err
		}
		return uint64(item["count"].(int64)), nil
	} else {
		return 0, cur.Err()
	}
}

func (repo *behaviorRepository) CountByUser(ctx mongo.SessionContext, uid string) (uint64, error) {
	behavior := models.Behavior{}
	filter := bson.D{{"uid", uid}}
	result := repo.collections.FindOne(ctx, filter,
		options.FindOne().SetProjection(bson.D{{"count", true}}))
	if err := result.Decode(&behavior); err != nil {
		return 0, err
	} else {
		return behavior.Count, nil
	}
}

func (repo *behaviorRepository) Delete(ctx mongo.SessionContext, uid string, infoID objectid.ObjectID) error {
	filter := bson.D{{"uid", uid}, {"items.infoId", infoID}}
	result := repo.collections.FindOne(ctx, filter)
	if err := result.Decode(nil); err != nil {
		return err
	} else {
		update := bson.D{
			{"$pull", bson.D{{"items", bson.D{{"infoId", infoID}}}}},
			{"$inc", bson.D{{"count", -1}}},
		}
		ur, err := repo.collections.UpdateOne(ctx, filter, update)
		if err != nil {
			return err
		} else if ur.ModifiedCount == 0 {
			return mongo.ErrNoDocuments
		}
		return nil
	}
}
