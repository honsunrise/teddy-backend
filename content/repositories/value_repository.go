package repositories

import (
	"fmt"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/objectid"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/options"
	"github.com/zhsyourai/teddy-backend/common/proto/content"
	"github.com/zhsyourai/teddy-backend/content/models"
)

type ValueRepository interface {
	Insert(ctx mongo.SessionContext, infoID objectid.ObjectID,
		segID objectid.ObjectID, value *models.Value) error

	FindOne(ctx mongo.SessionContext, infoID objectid.ObjectID,
		segID objectid.ObjectID, id string) (*models.Value, error)
	FindAll(ctx mongo.SessionContext, infoID objectid.ObjectID, segID objectid.ObjectID,
		page, size uint32, sorts []*content.Sort) ([]*models.Value, uint64, error)

	DeleteAll(ctx mongo.SessionContext, infoID objectid.ObjectID, segID objectid.ObjectID) error
	DeleteOne(ctx mongo.SessionContext, infoID objectid.ObjectID, segID objectid.ObjectID, id string) error
	DeleteMany(ctx mongo.SessionContext, infoID objectid.ObjectID, segID objectid.ObjectID, ids []string) error

	Update(ctx mongo.SessionContext, infoID objectid.ObjectID, segID objectid.ObjectID,
		id string, fields map[string]interface{}) error
}

func NewValueRepository(client *mongo.Client) (ValueRepository, error) {
	return &valueRepository{
		collections: client.Database("teddy").Collection("segment"),
	}, nil
}

type valueRepository struct {
	collections *mongo.Collection
}

func (repo *valueRepository) Insert(ctx mongo.SessionContext, infoID objectid.ObjectID,
	segID objectid.ObjectID, value *models.Value) error {
	result := repo.collections.FindOne(ctx, bson.D{{"_id", segID}, {"infoID", infoID}})
	if err := result.Decode(nil); err != nil {
		return err
	}

	filter := bson.D{
		{"_id", segID},
		{"infoID", infoID},
		{"values.time", bson.D{{"$ne", value.Time}}},
		{"values.value", bson.D{{"$ne", value.Value}}},
	}
	update := bson.D{
		{"$addToSet", bson.D{{"values", value}}},
		{"$inc", bson.D{{"count", 1}}},
	}
	ur, err := repo.collections.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	} else if ur.ModifiedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}

func (repo *valueRepository) FindOne(ctx mongo.SessionContext, infoID objectid.ObjectID,
	segID objectid.ObjectID, id string) (*models.Value, error) {
	pipeline := mongo.Pipeline{
		bson.D{{"$unwind", "$values"}},
		bson.D{{"$match", bson.D{
			{"_id", segID},
			{"infoID", infoID},
			{"values.id", id},
		}}},
		bson.D{{"$project", bson.D{
			{"id", "$values.id"},
			{"time", "$values.time"},
			{"value", "$values.value"},
		}}},
	}
	cur, err := repo.collections.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	defer cur.Close(ctx)
	if cur.Next(ctx) {
		item := models.Value{}
		err := cur.Decode(&item)
		if err != nil {
			return nil, err
		}
		return &item, nil
	} else {
		return nil, cur.Err()
	}
}

func (repo *valueRepository) FindAll(ctx mongo.SessionContext, infoID objectid.ObjectID, segID objectid.ObjectID,
	page, size uint32, sorts []*content.Sort) ([]*models.Value, uint64, error) {

	seg := models.Segment{}
	filter := bson.D{
		{"_id", segID},
		{"infoID", infoID},
	}
	result := repo.collections.FindOne(ctx, filter,
		options.FindOne().SetProjection(bson.D{{"count", true}}))
	if err := result.Decode(&seg); err != nil {
		return nil, 0, err
	}
	totalCount := seg.Count

	pipeline := mongo.Pipeline{
		bson.D{{"$unwind", "$values"}},
		bson.D{{"$match", bson.D{
			{"_id", segID},
			{"infoID", infoID},
		}}},
		bson.D{{"$skip", int64(size * page)}},
		bson.D{{"$limit", int64(size)}},
		bson.D{{"$project", bson.D{
			{"id", "$values.id"},
			{"time", "$values.time"},
			{"value", "$values.value"},
		}}},
	}

	cur, err := repo.collections.Aggregate(ctx, pipeline)
	defer cur.Close(ctx)
	items := make([]*models.Value, 0, size)
	for cur.Next(ctx) {
		var item models.Value
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

func (repo *valueRepository) DeleteAll(ctx mongo.SessionContext, infoID objectid.ObjectID, segID objectid.ObjectID) error {
	filter := bson.D{
		{"_id", segID},
		{"infoID", infoID},
	}
	result := repo.collections.FindOne(ctx, filter)
	if err := result.Decode(nil); err != nil {
		return err
	} else {
		update := bson.D{
			{"$pull", bson.D{{"values", bson.D{{"$slice", 0}}}}},
			{"$set", bson.D{{"count", 0}}},
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

func (repo *valueRepository) DeleteOne(ctx mongo.SessionContext, infoID objectid.ObjectID, segID objectid.ObjectID, id string) error {
	filter := bson.D{
		{"_id", segID},
		{"infoID", infoID},
	}
	result := repo.collections.FindOne(ctx, filter)
	if err := result.Decode(nil); err != nil {
		return err
	} else {
		update := bson.D{
			{"$pull", bson.D{{"values", bson.D{{"id", id}}}}},
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

func (repo *valueRepository) DeleteMany(ctx mongo.SessionContext, infoID objectid.ObjectID, segID objectid.ObjectID, ids []string) error {
	filter := bson.D{
		{"_id", segID},
		{"infoID", infoID},
	}
	result := repo.collections.FindOne(ctx, filter)
	if err := result.Decode(nil); err != nil {
		return err
	} else {
		update := bson.D{
			{"$pullAll", bson.D{{"values", bson.D{{"id", ids}}}}},
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

func (repo *valueRepository) Update(ctx mongo.SessionContext, infoID objectid.ObjectID, segID objectid.ObjectID,
	id string, fields map[string]interface{}) error {
	filter := bson.D{
		{"_id", segID},
		{"infoID", infoID},
	}
	result := repo.collections.FindOne(ctx, filter)
	if err := result.Decode(nil); err != nil {
		return err
	} else {
		var bsonFields = make(bson.D, 0, len(fields))
		for k, v := range fields {
			bsonFields = append(bsonFields, bson.E{Key: fmt.Sprintf("values.$[elem].%s", k), Value: v})
		}
		update := bson.D{{"$set", bsonFields}}
		ur, err := repo.collections.UpdateOne(ctx, filter, update, options.Update().
			SetArrayFilters(options.ArrayFilters{
				Filters: []interface{}{
					bson.D{{"elem.id", id}},
				},
			}).
			SetUpsert(true),
		)
		if err != nil {
			return err
		} else if ur.ModifiedCount == 0 {
			return mongo.ErrNoDocuments
		}
		return nil
	}
}
