package repositories

import (
	"context"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/zhsyourai/teddy-backend/common/types"
	"github.com/zhsyourai/teddy-backend/content/models"
	"time"
)

type TagRepository interface {
	Insert(tag *models.Tag) error
	FindOne(tag string) (models.Tag, error)
	FindAll(page uint32, size uint32, sorts []types.Sort) ([]models.Tag, error)
	IncUsage(tag string, inc uint64) error
	UpdateLastUse(tag string, lastUse time.Time) error
	DeleteOne(tag string) error
	DeleteAll(tags []string) error
}

func NewTagRepository(client *mongo.Client) (TagRepository, error) {
	return &tagRepository{
		ctx:         context.Background(),
		client:      client,
		collections: client.Database("teddy").Collection("tag"),
	}, nil
}

type tagRepository struct {
	ctx         context.Context
	client      *mongo.Client
	collections *mongo.Collection
}

func (repo *tagRepository) Insert(tag *models.Tag) error {
	filter := bson.D{{"_id", tag.Tag}}
	result := repo.collections.FindOne(repo.ctx, filter)
	if result.Decode(nil) == mongo.ErrNoDocuments {
		_, err := repo.collections.InsertOne(repo.ctx, tag)
		if err != nil {
			return err
		}
	} else {
		return ErrTagExisted
	}
	return nil
}

func (repo *tagRepository) FindOne(tag string) (models.Tag, error) {
	var tagEntry models.Tag
	filter := bson.D{{"_id", tag}}
	result := repo.collections.FindOne(repo.ctx, filter)
	err := result.Decode(&tagEntry)
	if err == nil {
		return models.Tag{}, err
	}
	return tagEntry, nil
}

func (repo *tagRepository) FindAll(page uint32, size uint32, sorts []types.Sort) ([]models.Tag, error) {
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
		{"$skip", int64(size * page)},
		{"$limit", int64(size)},
		{"$sort", itemsSorts},
	}
	repo.collections.Aggregate(repo.ctx, pipeline)
	items := make([]models.Tag, 0, size)
	defer cur.Close(repo.ctx)
	for cur.Next(repo.ctx) {
		var item models.Tag
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

func (repo *tagRepository) IncUsage(tag string, inc uint64) error {
	filter := bson.D{{"_id", tag}}
	update := bson.D{{"$inc", bson.D{{Key: "usage", Value: inc}}}}
	ur, err := repo.collections.UpdateOne(repo.ctx, filter, update)
	if err != nil {
		return err
	} else if ur.ModifiedCount == 0 {
		return ErrUpdateTag
	}
	return nil
}

func (repo *tagRepository) UpdateLastUse(tag string, lastUse time.Time) error {
	filter := bson.D{{"_id", tag}}
	update := bson.D{{"$set", bson.D{{Key: "lastUseTime", Value: lastUse}}}}
	ur, err := repo.collections.UpdateOne(repo.ctx, filter, update)
	if err != nil {
		return err
	} else if ur.ModifiedCount == 0 {
		return ErrUpdateTag
	}
	return nil
}

func (repo *tagRepository) DeleteOne(tag string) error {
	filter := bson.D{{"_id", tag}}
	_, err := repo.collections.DeleteOne(repo.ctx, filter)
	if err != nil {
		return err
	}
	return nil
}

func (repo *tagRepository) DeleteAll(tags []string) error {
	filter := bson.D{{"_id", bson.D{{"$in", bson.A{tags}}}}}
	_, err := repo.collections.DeleteOne(repo.ctx, filter)
	if err != nil {
		return err
	}
	return nil
}
