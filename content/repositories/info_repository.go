package repositories

import (
	"context"
	"errors"
	"fmt"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/zhsyourai/teddy-backend/common/models"
	"github.com/zhsyourai/teddy-backend/common/types"
)

var ErrTitleExisted = errors.New("info title has been existed")
var ErrUpdateInfo = errors.New("info update error")

type InfoRepository interface {
	InsertInfo(info *models.Info) error
	FindAll(page uint32, size uint32, sorts []types.Sort) ([]models.Info, error)
	FindByTags(tags []string, page uint32, size uint32, sorts []types.Sort) ([]models.Info, error)
	FindByUser(uid string, page uint32, size uint32, sorts []types.Sort) ([]models.Info, error)
	FindByTitle(title string, page uint32, size uint32, sorts []types.Sort) ([]models.Info, error)
	FindByTitleAndUser(title string, uid string, page uint32, size uint32, sorts []types.Sort) ([]models.Info, error)
	DeleteInfo(id string) error
	UpdateInfo(id string, fields map[string]interface{}) error
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

func (repo *infoRepository) InsertInfo(info *models.Info) error {
	filter := bson.D{{"title", info.Title}}
	result := repo.collections.FindOne(repo.ctx, filter)
	if result.Decode(nil) == mongo.ErrNoDocuments {
		_, err := repo.collections.InsertOne(repo.ctx, info)
		if err != nil {
			return err
		}
	} else {
		return ErrTitleExisted
	}
	return nil
}

func (repo *infoRepository) internalFindInfo(uid string, title string, tags []string, page uint32,
	size uint32, sorts []types.Sort) ([]models.Info, error) {
	var cur mongo.Cursor
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
		{"$match", dynFilter},
		{"$skip", int64(size * page)},
		{"$limit", int64(size)},
		{"$sort", itemsSorts},
	}
	repo.collections.Aggregate(repo.ctx, pipeline)
	items := make([]models.Info, 0, 50)
	defer cur.Close(repo.ctx)
	for cur.Next(repo.ctx) {
		var item models.Info
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

func (repo *infoRepository) FindAll(page uint32, size uint32, sorts []types.Sort) ([]models.Info, error) {
	return repo.internalFindInfo("", "", nil, page, size, sorts)
}

func (repo *infoRepository) FindByTags(tags []string, page uint32, size uint32, sorts []types.Sort) ([]models.Info, error) {
	return repo.internalFindInfo("", "", tags, page, size, sorts)
}

func (repo *infoRepository) FindByUser(uid string, page uint32, size uint32, sorts []types.Sort) ([]models.Info, error) {
	return repo.internalFindInfo(uid, "", nil, page, size, sorts)
}

func (repo *infoRepository) FindByTitle(title string, page uint32, size uint32, sorts []types.Sort) ([]models.Info, error) {
	return repo.internalFindInfo("", title, nil, page, size, sorts)
}

func (repo *infoRepository) FindByTitleAndUser(title string, uid string, page uint32, size uint32, sorts []types.Sort) ([]models.Info, error) {
	return repo.internalFindInfo(uid, title, nil, page, size, sorts)
}

func (repo *infoRepository) DeleteInfo(id string) error {
	filter := bson.D{{"_id", id}}
	_, err := repo.collections.DeleteOne(repo.ctx, filter)
	if err != nil {
		return err
	}
	return nil
}

func (repo *infoRepository) UpdateInfo(id string, fields map[string]interface{}) error {
	filter := bson.D{{"_id", id}}
	var bsonFields = make(bson.D, len(fields))
	for k, v := range fields {
		bsonFields = append(bsonFields, bson.E{Key: fmt.Sprintf("items.%s", k), Value: v})
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
