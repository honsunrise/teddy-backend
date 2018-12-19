package repositories

import (
	"context"
	"fmt"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/objectid"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/zhsyourai/teddy-backend/common/proto"
	"github.com/zhsyourai/teddy-backend/content/models"
)

type InfoRepository interface {
	Insert(info *models.Info) error
	IncWatchCount(id objectid.ObjectID, count int64) error
	FindOne(id objectid.ObjectID) (*models.Info, error)
	FindAll(page uint32, size uint32, sorts []*proto.Sort) ([]*models.Info, error)
	FindByTags(tags []string, page uint32, size uint32, sorts []*proto.Sort) ([]*models.Info, error)
	FindByUser(uid string, page uint32, size uint32, sorts []*proto.Sort) ([]*models.Info, error)
	FindByTitle(title string, page uint32, size uint32, sorts []*proto.Sort) ([]*models.Info, error)
	FindByTitleAndUser(title string, uid string, page uint32, size uint32, sorts []*proto.Sort) ([]*models.Info, error)
	Delete(id objectid.ObjectID) error
	Update(id objectid.ObjectID, fields map[string]interface{}) error
}

func NewInfoRepository(client *mongo.Client) (InfoRepository, error) {
	return &infoRepository{
		ctx:         context.Background(),
		client:      client,
		collections: client.Database("teddy").Collection("info"),
	}, nil
}

type infoRepository struct {
	ctx         context.Context
	client      *mongo.Client
	collections *mongo.Collection
}

func (repo *infoRepository) Insert(info *models.Info) error {
	filter := bson.D{{"title", info.Title}}
	result := repo.collections.FindOne(repo.ctx, filter)
	if result.Decode(nil) == mongo.ErrNoDocuments {
		info.Id = objectid.New()
		_, err := repo.collections.InsertOne(repo.ctx, info)
		if err != nil {
			return err
		}
	} else {
		return ErrTitleExisted
	}
	return nil
}

func (repo *infoRepository) IncWatchCount(id objectid.ObjectID, count int64) error {
	filter := bson.D{{"_id", id}}
	update := bson.D{{"$inc", bson.D{{Key: "watchCount", Value: count}}}}
	ur, err := repo.collections.UpdateOne(repo.ctx, filter, update)
	if err != nil {
		return err
	} else if ur.ModifiedCount == 0 {
		return ErrUpdateTag
	}
	return nil
}

func (repo *infoRepository) FindOne(id objectid.ObjectID) (*models.Info, error) {
	var info models.Info
	filter := bson.D{{"_id", id}}
	err := repo.collections.FindOne(repo.ctx, filter).Decode(&info)
	if err != nil {
		return nil, err
	}
	return &info, nil
}

func (repo *infoRepository) internalFindInfo(uid string, title string, tags []string, page uint32,
	size uint32, sorts []*proto.Sort) ([]*models.Info, error) {
	var cur mongo.Cursor

	pipeline := mongo.Pipeline{}

	var dynFilter = make(bson.D, 0, 3)
	if uid != "" {
		dynFilter = append(dynFilter, bson.E{Key: "uid", Value: uid})
	}
	if len(tags) != 0 {
		dynFilter = append(dynFilter, bson.E{Key: "tags", Value: bson.D{{"$in", bson.A{tags}}}})
	}
	if title != "" {
		dynFilter = append(dynFilter, bson.E{Key: "$text", Value: bson.D{{"$search", title}}})
	}
	if len(dynFilter) != 0 {
		pipeline = append(pipeline, bson.D{{"$match", dynFilter}})
	}

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

	cur, err := repo.collections.Aggregate(repo.ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cur.Close(repo.ctx)
	items := make([]*models.Info, 0, size)
	for cur.Next(repo.ctx) {
		var item models.Info
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

func (repo *infoRepository) FindAll(page uint32, size uint32, sorts []*proto.Sort) ([]*models.Info, error) {
	return repo.internalFindInfo("", "", nil, page, size, sorts)
}

func (repo *infoRepository) FindByTags(tags []string, page uint32, size uint32, sorts []*proto.Sort) ([]*models.Info, error) {
	return repo.internalFindInfo("", "", tags, page, size, sorts)
}

func (repo *infoRepository) FindByUser(uid string, page uint32, size uint32, sorts []*proto.Sort) ([]*models.Info, error) {
	return repo.internalFindInfo(uid, "", nil, page, size, sorts)
}

func (repo *infoRepository) FindByTitle(title string, page uint32, size uint32, sorts []*proto.Sort) ([]*models.Info, error) {
	return repo.internalFindInfo("", title, nil, page, size, sorts)
}

func (repo *infoRepository) FindByTitleAndUser(title string, uid string, page uint32, size uint32, sorts []*proto.Sort) ([]*models.Info, error) {
	return repo.internalFindInfo(uid, title, nil, page, size, sorts)
}

func (repo *infoRepository) Delete(id objectid.ObjectID) error {
	filter := bson.D{{"_id", id}}
	dr, err := repo.collections.DeleteOne(repo.ctx, filter)
	if err != nil {
		return err
	} else if dr.DeletedCount == 0 {
		return ErrDeleteInfo
	}
	return nil
}

func (repo *infoRepository) Update(id objectid.ObjectID, fields map[string]interface{}) error {
	filter := bson.D{{"_id", id}}
	var bsonFields = make(bson.D, 0, len(fields))
	for k, v := range fields {
		bsonFields = append(bsonFields, bson.E{Key: fmt.Sprintf("%s", k), Value: v})
	}
	update := bson.D{{"$set", bsonFields}}
	ur, err := repo.collections.UpdateOne(repo.ctx, filter, update)
	if err != nil {
		return err
	} else if ur.ModifiedCount == 0 {
		return ErrUpdateInfo
	}
	return nil
}
