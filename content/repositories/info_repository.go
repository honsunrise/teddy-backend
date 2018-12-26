package repositories

import (
	"fmt"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/objectid"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/zhsyourai/teddy-backend/common/proto/content"
	"github.com/zhsyourai/teddy-backend/content/models"
)

type InfoRepository interface {
	Insert(ctx mongo.SessionContext, info *models.Info) error
	IncWatchCount(ctx mongo.SessionContext, id objectid.ObjectID, count int64) error
	FindOne(ctx mongo.SessionContext, id objectid.ObjectID) (*models.Info, error)
	FindAll(ctx mongo.SessionContext, uid string, tags []*models.TypeAndTag,
		page uint32, size uint32, sorts []*content.Sort) ([]*models.Info, uint64, error)
	Delete(ctx mongo.SessionContext, id objectid.ObjectID) error
	Update(ctx mongo.SessionContext, id objectid.ObjectID, fields map[string]interface{}) error
}

func NewInfoRepository(client *mongo.Client) (InfoRepository, error) {
	return &infoRepository{
		collections: client.Database("teddy").Collection("info"),
	}, nil
}

type infoRepository struct {
	collections *mongo.Collection
}

func (repo *infoRepository) Insert(ctx mongo.SessionContext, info *models.Info) error {
	info.ID = objectid.New()
	_, err := repo.collections.InsertOne(ctx, info)
	if err != nil {
		return err
	}
	return nil
}

func (repo *infoRepository) IncWatchCount(ctx mongo.SessionContext, id objectid.ObjectID, count int64) error {
	filter := bson.D{{"_id", id}}
	update := bson.D{{"$inc", bson.D{{Key: "watchCount", Value: count}}}}
	ur, err := repo.collections.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	} else if ur.ModifiedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}

func (repo *infoRepository) FindOne(ctx mongo.SessionContext, id objectid.ObjectID) (*models.Info, error) {
	var info models.Info
	filter := bson.D{{"_id", id}}
	err := repo.collections.FindOne(ctx, filter).Decode(&info)
	if err != nil {
		return nil, err
	}
	return &info, nil
}

func (repo *infoRepository) internalFindInfo(ctx mongo.SessionContext, uid string, tags []*models.TypeAndTag,
	page uint32, size uint32, sorts []*content.Sort) ([]*models.Info, uint64, error) {
	var cur mongo.Cursor

	pipeline := mongo.Pipeline{}
	countPipeline := mongo.Pipeline{}

	var dynFilter = make(bson.D, 0, 2)
	if uid != "" {
		dynFilter = append(dynFilter, bson.E{Key: "uid", Value: uid})
	}
	if len(tags) != 0 {
		dynFilter = append(dynFilter, bson.E{Key: "tags", Value: bson.D{{"$all", tags}}})
	}
	if len(dynFilter) != 0 {
		pipeline = append(pipeline, bson.D{{"$match", dynFilter}})
		countPipeline = append(countPipeline, bson.D{{"$match", dynFilter}})
	}

	countPipeline = append(countPipeline, bson.D{{"$count", "count"}})
	cur, err := repo.collections.Aggregate(ctx, countPipeline)
	if err != nil {
		return nil, 0, err
	}
	elem := struct {
		Count uint64 `bson:"count"`
	}{}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		err := cur.Decode(&elem)
		if err != nil {
			return nil, 0, err
		}
		break
	}
	err = cur.Err()
	if err != nil {
		return nil, 0, err
	}
	totalCount := elem.Count

	pipeline = append(pipeline, bson.D{{"$skip", int64(size * page)}})
	pipeline = append(pipeline, bson.D{{"$limit", int64(size)}})
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
	cur, err = repo.collections.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, 0, err
	}
	defer cur.Close(ctx)
	items := make([]*models.Info, 0, size)
	for cur.Next(ctx) {
		var item models.Info
		err := cur.Decode(&item)
		if err != nil {
			return nil, 0, err
		}
		items = append(items, &item)
	}
	err = cur.Err()
	if err != nil {
		return nil, 0, err
	}
	return items, uint64(totalCount), nil
}

func (repo *infoRepository) FindAll(ctx mongo.SessionContext, uid string, tags []*models.TypeAndTag,
	page uint32, size uint32, sorts []*content.Sort) ([]*models.Info, uint64, error) {
	return repo.internalFindInfo(ctx, uid, tags, page, size, sorts)
}

func (repo *infoRepository) Delete(ctx mongo.SessionContext, id objectid.ObjectID) error {
	filter := bson.D{{"_id", id}}
	dr, err := repo.collections.DeleteOne(ctx, filter)
	if err != nil {
		return err
	} else if dr.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}

func (repo *infoRepository) Update(ctx mongo.SessionContext, id objectid.ObjectID, fields map[string]interface{}) error {
	filter := bson.D{{"_id", id}}
	var bsonFields = make(bson.D, 0, len(fields))
	for k, v := range fields {
		bsonFields = append(bsonFields, bson.E{Key: fmt.Sprintf("%s", k), Value: v})
	}
	update := bson.D{{"$set", bsonFields}}
	ur, err := repo.collections.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	} else if ur.ModifiedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}
