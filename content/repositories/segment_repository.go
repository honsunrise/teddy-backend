package repositories

import (
	"fmt"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/objectid"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/zhsyourai/teddy-backend/common/proto/content"
	"github.com/zhsyourai/teddy-backend/content/models"
)

type SegmentRepository interface {
	Insert(ctx mongo.SessionContext, segment *models.Segment) error
	FindByInfoIDAndNoAndTitleAndLabels(ctx mongo.SessionContext,
		infoID objectid.ObjectID, no uint64, title string, labels []string) (*models.Segment, error)
	FindOne(ctx mongo.SessionContext, infoID objectid.ObjectID, id objectid.ObjectID) (*models.Segment, error)
	FindAll(ctx mongo.SessionContext, infoID objectid.ObjectID,
		labels []string, page, size uint32, sorts []*content.Sort) ([]*models.Segment, uint64, error)

	DeleteByInfoID(ctx mongo.SessionContext, infoID objectid.ObjectID) error
	DeleteOne(ctx mongo.SessionContext, infoID objectid.ObjectID, id objectid.ObjectID) error
	DeleteMany(ctx mongo.SessionContext, ids []objectid.ObjectID) error

	Update(ctx mongo.SessionContext, id objectid.ObjectID, fields map[string]interface{}) error
}

func NewSegmentRepository(client *mongo.Client) (SegmentRepository, error) {
	return &segmentRepository{
		collections: client.Database("teddy").Collection("segment"),
	}, nil
}

type segmentRepository struct {
	collections *mongo.Collection
}

func (repo *segmentRepository) FindByInfoIDAndNoAndTitleAndLabels(ctx mongo.SessionContext,
	infoID objectid.ObjectID, no uint64, title string, labels []string) (*models.Segment, error) {
	var segment models.Segment
	filter := bson.D{
		{"infoID", infoID},
		{"title", title},
		{"labels", bson.D{{"$all", labels}}},
		{"no", no},
	}
	result := repo.collections.FindOne(ctx, filter)
	err := result.Decode(&segment)
	if err != nil {
		return nil, err
	}
	return &segment, nil
}

func (repo *segmentRepository) Insert(ctx mongo.SessionContext, segment *models.Segment) error {
	segment.ID = objectid.New()
	_, err := repo.collections.InsertOne(ctx, segment)
	if err != nil {
		return err
	}
	return nil
}

func (repo *segmentRepository) FindOne(ctx mongo.SessionContext, infoID objectid.ObjectID, id objectid.ObjectID) (*models.Segment, error) {
	var segment models.Segment
	filter := bson.D{{"_id", id}, {"infoID", infoID}}
	result := repo.collections.FindOne(ctx, filter)
	err := result.Decode(&segment)
	if err != nil {
		return nil, err
	}
	return &segment, nil
}

func (repo *segmentRepository) FindAll(ctx mongo.SessionContext,
	infoID objectid.ObjectID, labels []string, page, size uint32, sorts []*content.Sort) ([]*models.Segment, uint64, error) {
	pipeline := mongo.Pipeline{}
	countPipeline := mongo.Pipeline{}

	var dynFilter = make(bson.D, 0, 1)
	dynFilter = append(dynFilter, bson.E{Key: "infoID", Value: infoID})

	if len(labels) != 0 {
		dynFilter = append(dynFilter,
			bson.E{Key: "labels", Value: bson.D{{"$all", labels}}})
	}
	pipeline = append(pipeline, bson.D{{"$match", dynFilter}})
	countPipeline = append(countPipeline,
		bson.D{{"$match", dynFilter}},
		bson.D{{"$count", "count"}})

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

	pipeline = append(pipeline, bson.D{{"$skip", int64(size * page)}})
	pipeline = append(pipeline, bson.D{{"$limit", int64(size)}})

	cur, err = repo.collections.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, 0, err
	}
	items := make([]*models.Segment, 0, size)
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var item models.Segment
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

func (repo *segmentRepository) DeleteByInfoID(ctx mongo.SessionContext, infoID objectid.ObjectID) error {
	filter := bson.D{{"infoID", infoID}}
	dr, err := repo.collections.DeleteOne(ctx, filter)
	if err != nil {
		return err
	} else if dr.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}

func (repo *segmentRepository) DeleteOne(ctx mongo.SessionContext, infoID objectid.ObjectID, id objectid.ObjectID) error {
	filter := bson.D{{"_id", id}, {"infoID", infoID}}
	dr, err := repo.collections.DeleteOne(ctx, filter)
	if err != nil {
		return err
	} else if dr.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}

func (repo *segmentRepository) DeleteMany(ctx mongo.SessionContext, ids []objectid.ObjectID) error {
	filter := bson.D{{"_id", bson.D{{"$in", ids}}}}
	dr, err := repo.collections.DeleteOne(ctx, filter)
	if err != nil {
		return err
	} else if dr.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}

func (repo *segmentRepository) Update(ctx mongo.SessionContext, id objectid.ObjectID, fields map[string]interface{}) error {
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
