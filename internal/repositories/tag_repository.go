package repositories

import (
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"github.com/mongodb/mongo-go-driver/mongo"
	"teddy-backend/internal/models"
	"teddy-backend/internal/proto/content"
	"time"
)

type TagRepository interface {
	Insert(ctx mongo.SessionContext, tag *models.Tag) error
	FindByTypeAndTag(ctx mongo.SessionContext, tp string, tag string) (*models.Tag, error)
	FindOne(ctx mongo.SessionContext, id primitive.ObjectID) (*models.Tag, error)
	FindAll(ctx mongo.SessionContext, tp string,
		page, size uint64, sorts []*content.Sort) ([]*models.Tag, uint64, error)
	IncUsage(ctx mongo.SessionContext, id primitive.ObjectID, inc int64) error
	UpdateLastUse(ctx mongo.SessionContext, id primitive.ObjectID, lastUse time.Time) error
	DeleteOne(ctx mongo.SessionContext, id primitive.ObjectID) error
	DeleteMany(ctx mongo.SessionContext, ids []primitive.ObjectID) error
}

func NewTagRepository(client *mongo.Client) (TagRepository, error) {
	return &tagRepository{
		collections: client.Database("teddy").Collection("tag"),
	}, nil
}

type tagRepository struct {
	collections *mongo.Collection
}

func (repo *tagRepository) Insert(ctx mongo.SessionContext, tag *models.Tag) error {
	tag.ID = primitive.NewObjectID()
	_, err := repo.collections.InsertOne(ctx, tag)
	if err != nil {
		return err
	}
	return nil
}

func (repo *tagRepository) FindByTypeAndTag(ctx mongo.SessionContext, tp string, tag string) (*models.Tag, error) {
	var tagEntry models.Tag
	filter := bson.D{{"type", tp}, {"tag", tag}}
	result := repo.collections.FindOne(ctx, filter)
	err := result.Decode(&tagEntry)
	if err != nil {
		return nil, err
	}
	return &tagEntry, nil
}

func (repo *tagRepository) FindOne(ctx mongo.SessionContext, id primitive.ObjectID) (*models.Tag, error) {
	var tagEntry models.Tag
	filter := bson.D{{"_id", id}}
	result := repo.collections.FindOne(ctx, filter)
	err := result.Decode(&tagEntry)
	if err != nil {
		return nil, err
	}
	return &tagEntry, nil
}

func (repo *tagRepository) FindAll(ctx mongo.SessionContext, tp string,
	page, size uint64, sorts []*content.Sort) ([]*models.Tag, uint64, error) {
	pipeline := mongo.Pipeline{}
	countPipeline := mongo.Pipeline{}

	if tp != "" {
		pipeline = append(pipeline, bson.D{{"$match", bson.D{{"type", tp}}}})
		countPipeline = append(countPipeline, bson.D{{"$count", "count"}})
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
	items := make([]*models.Tag, 0, size)
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var item models.Tag
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
	return items, totalCount, nil
}

func (repo *tagRepository) IncUsage(ctx mongo.SessionContext, id primitive.ObjectID, inc int64) error {
	filter := bson.D{{"_id", id}}
	update := bson.D{{"$inc", bson.D{{Key: "usage", Value: inc}}}}
	ur, err := repo.collections.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	} else if ur.ModifiedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}

func (repo *tagRepository) UpdateLastUse(ctx mongo.SessionContext, id primitive.ObjectID, lastUse time.Time) error {
	filter := bson.D{{"_id", id}}
	update := bson.D{{"$set", bson.D{{Key: "lastUseTime", Value: lastUse}}}}
	ur, err := repo.collections.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	} else if ur.ModifiedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}

func (repo *tagRepository) DeleteOne(ctx mongo.SessionContext, id primitive.ObjectID) error {
	filter := bson.D{{"_id", id}}
	_, err := repo.collections.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	return nil
}

func (repo *tagRepository) DeleteMany(ctx mongo.SessionContext, ids []primitive.ObjectID) error {
	filter := bson.D{{"_id", bson.D{{"$in", ids}}}}
	_, err := repo.collections.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	return nil
}
