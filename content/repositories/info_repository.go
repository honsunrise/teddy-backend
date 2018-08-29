package repositories

import (
	"context"
	"errors"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/zhsyourai/teddy-backend/common/models"
	"github.com/zhsyourai/teddy-backend/common/types"
)

var ErrTitleExisted = errors.New("info title has been existed")

type InfoRepository interface {
	InsertInfo(info *models.Info) error
	FindAll(page uint32, size uint32, sorts []types.Sort) ([]models.Info, error)
	FindByTags(tags []string, page uint32, size uint32, sorts []types.Sort) ([]models.Info, error)
	FindByUser(uid string, page uint32, size uint32, sorts []types.Sort) ([]models.Info, error)
	FindByTitle(title string, page uint32, size uint32, sorts []types.Sort) ([]models.Info, error)
	FindBYTitleAndUser(title string, uid string, page uint32, size uint32, sorts []types.Sort) ([]models.Info, error)
	DeleteInfo(id string) error
	UpdateInfo(id string, info *models.Info) error
}

func NewInfoRepository(client *mongo.Client) (InfoRepository, error) {
	return &infoRepository{
		ctx:         context.Background(),
		client:      client,
		collections: client.Database("Teddy").Collection("Content"),
	}, nil
}

type infoRepository struct {
	ctx         context.Context
	client      *mongo.Client
	collections *mongo.Collection
}

func (repo *infoRepository) InsertInfo(info *models.Info) error {
	filter := bson.NewDocument(bson.EC.String("title", info.Title))
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
	var uidFilter *bson.Element = nil
	if uid != "" {
		uidFilter = bson.EC.String("uid", uid)
	}
	var tagsFilter *bson.Element = nil
	if len(tags) != 0 {
		var bsonTags = make([]*bson.Value, 0, len(tags))
		for _, id := range tags {
			bsonTags = append(bsonTags, bson.VC.String(id))
		}
		tagsFilter = bson.EC.SubDocumentFromElements("tags",
			bson.EC.ArrayFromElements("$in", bsonTags...),
		)
	}
	var titleFilter *bson.Element = nil
	if title != "" {
		titleFilter = bson.EC.SubDocumentFromElements("$text",
			bson.EC.String("$search", title))
	}

	var dynFilter = make([]*bson.Element, 0, 3)
	if uidFilter != nil {
		dynFilter = append(dynFilter, uidFilter)
	}
	if tagsFilter != nil {
		dynFilter = append(dynFilter, tagsFilter)
	}
	if titleFilter != nil {
		dynFilter = append(dynFilter, titleFilter)
	}
	var itemsSorts []*bson.Element
	if len(sorts) != 0 {
		itemsSorts = make([]*bson.Element, 0, len(sorts))
		for _, sort := range sorts {
			if sort.Order == types.ASC {
				itemsSorts = append(itemsSorts, bson.EC.Int32(sort.Name, 1))
			} else {
				itemsSorts = append(itemsSorts, bson.EC.Int32(sort.Name, -1))
			}
		}
	}
	pipeline := bson.NewDocument(
		bson.EC.SubDocumentFromElements("$match", dynFilter...),
		bson.EC.Int64("$skip", int64(size*page)),
		bson.EC.Int64("$limit", int64(size)),
		bson.EC.SubDocumentFromElements("$sort", itemsSorts...),
	)
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
	panic("implement me")
}

func (repo *infoRepository) UpdateInfo(id string, info *models.Info) error {
	panic("implement me")
}
