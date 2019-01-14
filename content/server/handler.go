package server

import (
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/mongodb/mongo-go-driver/bson/objectid"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/rs/xid"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"teddy-backend/common/proto/content"
	"teddy-backend/content/models"
	"teddy-backend/content/repositories"
	"time"
)

func NewContentServer(client *mongo.Client) (content.ContentServer, error) {
	// New Repository
	infoRepo, err := repositories.NewInfoRepository(client)
	if err != nil {
		return nil, err
	}
	favoriteRepo, err := repositories.NewFavoriteRepository(client)
	if err != nil {
		return nil, err
	}
	segmentRepo, err := repositories.NewSegmentRepository(client)
	if err != nil {
		return nil, err
	}
	tagRepo, err := repositories.NewTagRepository(client)
	if err != nil {
		return nil, err
	}
	thumbUpRepo, err := repositories.NewThumbUpRepository(client)
	if err != nil {
		return nil, err
	}
	thumbDownRepo, err := repositories.NewThumbDownRepository(client)
	if err != nil {
		return nil, err
	}
	valueRepo, err := repositories.NewValueRepository(client)
	if err != nil {
		return nil, err
	}
	instance := &contentHandler{
		client:        client,
		infoRepo:      infoRepo,
		segRepo:       segmentRepo,
		tagRepo:       tagRepo,
		valueRepo:     valueRepo,
		favoriteRepo:  favoriteRepo,
		thumbDownRepo: thumbDownRepo,
		thumbUpRepo:   thumbUpRepo,
	}
	return instance, nil
}

type contentHandler struct {
	client        *mongo.Client
	infoRepo      repositories.InfoRepository
	segRepo       repositories.SegmentRepository
	tagRepo       repositories.TagRepository
	valueRepo     repositories.ValueRepository
	favoriteRepo  repositories.BehaviorRepository
	thumbUpRepo   repositories.BehaviorRepository
	thumbDownRepo repositories.BehaviorRepository
}

func (h *contentHandler) GetValue(ctx context.Context, req *content.ValueOneReq) (*content.ValueResp, error) {
	var resp content.ValueResp
	if err := validateValueOneReq(req); err != nil {
		return nil, err
	}

	infoID, err := objectid.FromHex(req.InfoID)
	if err != nil {
		return nil, err
	}

	segID, err := objectid.FromHex(req.SegID)
	if err != nil {
		return nil, err
	}

	err = h.client.UseSession(ctx, func(sessionContext mongo.SessionContext) error {
		var err error
		value, err := h.valueRepo.FindOne(sessionContext, infoID, segID, req.ValID)

		if err != nil {
			return err
		}

		resp.Id = value.ID
		resp.Value = value.Value
		var pbTime *timestamp.Timestamp
		pbTime, err = ptypes.TimestampProto(value.Time)
		if err != nil {
			return err
		}
		resp.Time = pbTime
		return nil
	})

	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (h *contentHandler) GetValues(ctx context.Context, req *content.GetValuesReq) (*content.ValuesResp, error) {
	var resp content.ValuesResp
	if err := validateGetValuesReq(req); err != nil {
		return nil, err
	}

	infoID, err := objectid.FromHex(req.InfoID)
	if err != nil {
		return nil, err
	}

	segID, err := objectid.FromHex(req.SegID)
	if err != nil {
		return nil, err
	}

	err = h.client.UseSession(ctx, func(sessionContext mongo.SessionContext) error {
		var err error
		values, totalCount, err := h.valueRepo.FindAll(sessionContext, infoID, segID, req.Page, req.Size, req.Sorts)

		if err != nil {
			return err
		}

		result := make([]*content.ValueResp, 0, len(values))
		for _, value := range values {
			pbValue := &content.ValueResp{}
			copyFromValueToPBValue(value, pbValue)
			result = append(result, pbValue)
		}
		resp.Items = result
		resp.TotalCount = totalCount
		return nil
	})

	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (h *contentHandler) InsertValue(ctx context.Context, req *content.InsertValueReq) (*content.InsertValueResp, error) {
	var resp content.InsertValueResp
	if err := validateInsertValueReq(req); err != nil {
		return nil, err
	}

	infoID, err := objectid.FromHex(req.InfoID)
	if err != nil {
		return nil, err
	}

	segID, err := objectid.FromHex(req.SegID)
	if err != nil {
		return nil, err
	}

	err = h.client.UseSession(ctx, func(sessionContext mongo.SessionContext) error {
		var err error
		err = sessionContext.StartTransaction()
		if err != nil {
			return err
		}

		defer func() {
			if err == nil {
				err = sessionContext.CommitTransaction(sessionContext)
			} else {
				err = sessionContext.AbortTransaction(context.Background())
			}
		}()

		_, err = h.segRepo.FindOne(sessionContext, infoID, segID)
		if err != nil {
			log.Errorf("info can't find error %v", err)
			return err
		}

		reqTime, err := ptypes.Timestamp(req.Time)
		if err != nil {
			return err
		}

		value := models.Value{
			ID:    xid.New().String(),
			Time:  reqTime,
			Value: req.Value,
		}

		err = h.valueRepo.Insert(sessionContext, infoID, segID, &value)
		if err != nil {
			return err
		}
		resp.ValueID = value.ID
		return nil
	})

	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (h *contentHandler) EditValue(ctx context.Context, req *content.EditValueReq) (*empty.Empty, error) {
	var resp empty.Empty
	if err := validateEditValueReq(req); err != nil {
		return nil, err
	}

	infoID, err := objectid.FromHex(req.InfoID)
	if err != nil {
		return nil, err
	}

	segID, err := objectid.FromHex(req.SegID)
	if err != nil {
		return nil, err
	}

	err = h.client.UseSession(ctx, func(sessionContext mongo.SessionContext) error {
		var err error
		err = sessionContext.StartTransaction()
		if err != nil {
			return err
		}

		defer func() {
			if err == nil {
				err = sessionContext.CommitTransaction(sessionContext)
			} else {
				err = sessionContext.AbortTransaction(context.Background())
			}
		}()

		segment, err := h.segRepo.FindOne(sessionContext, infoID, segID)
		if err != nil {
			log.Errorf("find segment error %v", err)
			return err
		}

		updateMap := map[string]interface{}{}
		if req.Time != nil {
			var reqTime time.Time
			reqTime, err = ptypes.Timestamp(req.Time)
			if err != nil {
				return err
			}
			updateMap["time"] = reqTime
		}

		if req.Value != "" {
			updateMap["value"] = req.Value
		}

		err = h.valueRepo.Update(sessionContext, infoID, segment.ID, req.ValID, updateMap)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (h *contentHandler) DeleteValue(ctx context.Context, req *content.ValueOneReq) (*empty.Empty, error) {
	var resp empty.Empty
	if err := validateValueOneReq(req); err != nil {
		return nil, err
	}

	infoID, err := objectid.FromHex(req.InfoID)
	if err != nil {
		return nil, ErrInternal
	}

	segID, err := objectid.FromHex(req.SegID)
	if err != nil {
		return nil, ErrInternal
	}

	err = h.client.UseSession(ctx, func(sessionContext mongo.SessionContext) error {
		var err error
		err = h.valueRepo.DeleteOne(sessionContext, infoID, segID, req.ValID)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (h *contentHandler) GetSegments(ctx context.Context, req *content.GetSegmentsReq) (*content.SegmentsResp, error) {
	var resp content.SegmentsResp
	if err := validateGetSegmentsReq(req); err != nil {
		return nil, err
	}

	infoID, err := objectid.FromHex(req.InfoID)
	if err != nil {
		return nil, err
	}

	err = h.client.UseSession(ctx, func(sessionContext mongo.SessionContext) error {
		var err error
		segments, totalCount, err := h.segRepo.FindAll(sessionContext, infoID, req.Labels, req.Page, req.Size, req.Sorts)

		if err != nil {
			return err
		}

		result := make([]*content.SegmentResp, 0, len(segments))
		for _, segment := range segments {
			pbSegment := &content.SegmentResp{}
			copyFromSegmentToPBSegment(segment, pbSegment)
			result = append(result, pbSegment)
		}
		resp.Items = result
		resp.TotalCount = totalCount
		return nil
	})

	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (h *contentHandler) GetSegment(ctx context.Context, req *content.SegmentOneReq) (*content.SegmentResp, error) {
	var resp content.SegmentResp
	if err := validateSegmentOneReq(req); err != nil {
		return nil, err
	}

	infoID, err := objectid.FromHex(req.InfoID)
	if err != nil {
		return nil, err
	}

	segID, err := objectid.FromHex(req.SegID)
	if err != nil {
		return nil, err
	}

	err = h.client.UseSession(ctx, func(sessionContext mongo.SessionContext) error {
		var err error
		segment, err := h.segRepo.FindOne(sessionContext, infoID, segID)
		if err == mongo.ErrNoDocuments {
			return ErrSegmentNotExists
		} else if err != nil {
			return ErrInternal
		}

		copyFromSegmentToPBSegment(segment, &resp)
		return nil
	})

	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (h *contentHandler) PublishSegment(ctx context.Context, req *content.PublishSegmentReq) (*content.PublishSegmentResp, error) {
	var resp content.PublishSegmentResp
	if err := validatePublishSegmentReq(req); err != nil {
		return nil, err
	}

	infoID, err := objectid.FromHex(req.InfoID)
	if err != nil {
		return nil, err
	}

	err = h.client.UseSession(ctx, func(sessionContext mongo.SessionContext) error {
		var err error
		err = sessionContext.StartTransaction()
		if err != nil {
			return err
		}

		defer func() {
			if err == nil {
				err = sessionContext.CommitTransaction(sessionContext)
			} else {
				err = sessionContext.AbortTransaction(context.Background())
			}
		}()

		info, err := h.infoRepo.FindOne(sessionContext, infoID)
		if err != nil {
			log.Errorf("info can't find error %v", err)
			return err
		}

		_, err = h.segRepo.FindByInfoIDAndNoAndTitleAndLabels(sessionContext, info.ID, req.No, req.Title, req.Labels)
		if err != mongo.ErrNoDocuments {
			log.Errorf("check segment error %v", err)
			return ErrInternal
		} else if err == nil {
			return ErrSegmentExists
		}

		segment := models.Segment{
			ID:         objectid.New(),
			InfoID:     infoID,
			No:         req.No,
			Title:      req.Title,
			Labels:     req.Labels,
			WatchCount: 0,
			Count:      0,
			Values:     []models.Value{},
		}

		err = h.segRepo.Insert(sessionContext, &segment)
		if err != nil {
			return err
		}
		resp.SegID = segment.ID.Hex()
		return nil
	})

	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (h *contentHandler) EditSegment(ctx context.Context, req *content.EditSegmentReq) (*empty.Empty, error) {
	var resp empty.Empty
	if err := validateEditSegmentReq(req); err != nil {
		return nil, err
	}

	segID, err := objectid.FromHex(req.SegID)
	if err != nil {
		return nil, err
	}

	infoID, err := objectid.FromHex(req.InfoID)
	if err != nil {
		return nil, err
	}

	err = h.client.UseSession(ctx, func(sessionContext mongo.SessionContext) error {
		var err error
		err = sessionContext.StartTransaction()
		if err != nil {
			return err
		}

		defer func() {
			if err == nil {
				err = sessionContext.CommitTransaction(sessionContext)
			} else {
				err = sessionContext.AbortTransaction(context.Background())
			}
		}()

		segment, err := h.segRepo.FindOne(sessionContext, infoID, segID)
		if err != nil {
			log.Errorf("find segment error %v", err)
			return err
		}

		if req.No >= 0 {
			segment.No = req.No
		}

		if req.Title != "" {
			segment.Title = req.Title
		}

		if len(req.Labels) != 0 {
			segment.Labels = req.Labels
		}

		_, err = h.segRepo.FindByInfoIDAndNoAndTitleAndLabels(sessionContext, infoID, req.No, req.Title, req.Labels)
		if err != mongo.ErrNoDocuments {
			log.Errorf("check segment error %v", err)
			return ErrInternal
		} else if err == nil {
			return ErrSegmentExists
		}

		err = h.segRepo.Update(sessionContext, segment.ID, map[string]interface{}{
			"no":     segment.No,
			"title":  segment.Title,
			"labels": segment.Labels,
		})
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (h *contentHandler) DeleteSegment(ctx context.Context, req *content.SegmentOneReq) (*empty.Empty, error) {
	var resp empty.Empty
	if err := validateSegmentOneReq(req); err != nil {
		return nil, err
	}

	infoID, err := objectid.FromHex(req.InfoID)
	if err != nil {
		return nil, ErrInternal
	}

	segID, err := objectid.FromHex(req.SegID)
	if err != nil {
		return nil, ErrInternal
	}

	err = h.client.UseSession(ctx, func(sessionContext mongo.SessionContext) error {
		var err error
		err = h.segRepo.DeleteOne(sessionContext, infoID, segID)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (h *contentHandler) GetTags(ctx context.Context, req *content.GetTagsReq) (*content.TagsResp, error) {
	var resp content.TagsResp
	if err := validateGetTagsReq(req); err != nil {
		return nil, err
	}

	err := h.client.UseSession(ctx, func(sessionContext mongo.SessionContext) error {
		var err error
		tags, totalCount, err := h.tagRepo.FindAll(sessionContext, req.Type, req.Page, req.Size, req.Sorts)

		if err != nil {
			return err
		}

		result := make([]*content.TagResp, 0, len(tags))
		for _, tag := range tags {
			pbTag := &content.TagResp{}
			copyFromTagToPBTag(tag, pbTag)
			result = append(result, pbTag)
		}
		resp.Items = result
		resp.TotalCount = totalCount
		return nil
	})

	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (h *contentHandler) GetTag(ctx context.Context, req *content.GetTagReq) (*content.TagResp, error) {
	var resp content.TagResp
	if err := validateGetTagReq(req); err != nil {
		return nil, err
	}

	tagID, err := objectid.FromHex(req.Id)
	if err != nil {
		return nil, err
	}

	err = h.client.UseSession(ctx, func(sessionContext mongo.SessionContext) error {
		var err error
		tag, err := h.tagRepo.FindOne(sessionContext, tagID)

		if err != nil {
			return err
		}

		resp.Tag = tag.Tag
		resp.Type = tag.Type
		return nil
	})

	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (h *contentHandler) PublishInfo(ctx context.Context, req *content.PublishInfoReq) (*content.PublishInfoResp, error) {
	var resp content.PublishInfoResp
	if err := validatePublishInfoReq(req); err != nil {
		return nil, err
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()

	err := h.client.UseSession(timeoutCtx, func(sessionContext mongo.SessionContext) error {
		var err error
		now := time.Now()
		err = sessionContext.StartTransaction()
		if err != nil {
			return err
		}

		defer func() {
			if err == nil {
				err = sessionContext.CommitTransaction(sessionContext)
			} else {
				err = sessionContext.AbortTransaction(context.Background())
			}
		}()

		tags := make([]*models.TypeAndTag, 0, len(req.Tags))
		for _, v := range req.Tags {
			result, err := h.tagRepo.FindByTypeAndTag(sessionContext, v.Type, v.Tag)
			if err != nil {
				log.Errorf("tag can't find error %v", err)
				return err
			}
			tags = append(tags, &models.TypeAndTag{
				Type: result.Type,
				Tag:  result.Tag,
			})
			err = h.tagRepo.IncUsage(sessionContext, result.ID, 1)
			if err != nil {
				log.Errorf("inc tag usage error %v", err)
				return err
			}
		}

		if exists, err := h.infoRepo.ExistsByTitleAndAuthorAndCountry(sessionContext,
			req.Title, req.Author, req.Country); err != nil {
			return err
		} else if exists {
			return ErrInfoExists
		}

		contentTime, err := ptypes.Timestamp(req.ContentTime)
		if err != nil {
			return err
		}

		info := models.Info{
			ID:               objectid.New(),
			UID:              req.Uid,
			Title:            req.Title,
			Author:           req.Author,
			Country:          req.Country,
			Summary:          req.Summary,
			CoverResources:   req.CoverResources,
			PublishTime:      now,
			LastReviewTime:   time.Now(),
			Valid:            true,
			WatchCount:       0,
			Tags:             tags,
			LatestModifyTime: now,
			CanReview:        req.CanReview,
			Archived:         false,
			ContentTime:      contentTime,
			LatestSegmentID:  objectid.NilObjectID,
			SegmentCount:     0,
		}

		err = h.infoRepo.Insert(sessionContext, &info)
		if err != nil {
			return err
		}

		resp.InfoID = info.ID.Hex()
		return nil
	})

	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func difference(a, b []objectid.ObjectID) ([]objectid.ObjectID, []objectid.ObjectID) {
	ma := map[objectid.ObjectID]bool{}
	for _, x := range a {
		ma[x] = true
	}
	mb := map[objectid.ObjectID]bool{}
	for _, x := range b {
		mb[x] = true
	}
	ab := make([]objectid.ObjectID, 0, 20)
	for _, x := range a {
		if _, ok := mb[x]; !ok {
			ab = append(ab, x)
		}
	}
	ba := make([]objectid.ObjectID, 0, 20)
	for _, x := range a {
		if _, ok := mb[x]; !ok {
			ab = append(ab, x)
		}
	}
	return ab, ba
}

func (h *contentHandler) EditInfo(ctx context.Context, req *content.EditInfoReq) (*empty.Empty, error) {
	var resp empty.Empty
	if err := validateEditInfoReq(req); err != nil {
		return nil, err
	}

	infoID, err := objectid.FromHex(req.InfoID)
	if err != nil {
		return nil, err
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	err = h.client.UseSession(timeoutCtx, func(sessionContext mongo.SessionContext) error {
		var err error
		now := time.Now()
		err = sessionContext.StartTransaction()
		if err != nil {
			return err
		}

		defer func() {
			if err == nil {
				err = sessionContext.CommitTransaction(sessionContext)
			} else {
				err = sessionContext.AbortTransaction(context.Background())
			}
		}()

		curInfo, err := h.infoRepo.FindOne(sessionContext, infoID)
		if err != nil {
			return err
		}

		curTagIDs := make([]objectid.ObjectID, 0, len(curInfo.Tags))
		for _, v := range curInfo.Tags {
			result, err := h.tagRepo.FindByTypeAndTag(sessionContext, v.Type, v.Tag)
			if err != nil {
				log.Errorf("tag can't find error %v", err)
				return err
			}
			curTagIDs = append(curTagIDs, result.ID)
		}

		tagIDs := make([]objectid.ObjectID, 0, len(req.Tags))
		tags := make([]*models.TypeAndTag, 0, len(req.Tags))
		for _, v := range req.Tags {
			result, err := h.tagRepo.FindByTypeAndTag(sessionContext, v.Type, v.Tag)
			if err != nil {
				log.Errorf("tag can't find error %v", err)
				return err
			}
			tagIDs = append(tagIDs, result.ID)
			tags = append(tags, &models.TypeAndTag{
				Tag:  result.Tag,
				Type: result.Type,
			})
		}

		sub, inc := difference(curTagIDs, tagIDs)

		for _, v := range inc {
			err = h.tagRepo.IncUsage(sessionContext, v, 1)
			if err != nil {
				log.Errorf("inc tag usage error %v", err)
				return err
			}
		}

		for _, v := range sub {
			err = h.tagRepo.IncUsage(sessionContext, v, -1)
			if err != nil {
				log.Errorf("sub tag usage error %v", err)
				return err
			}
		}

		err = h.infoRepo.Update(sessionContext, infoID, map[string]interface{}{
			"uid":            req.Uid,
			"title":          req.Title,
			"author":         req.Author,
			"summary":        req.Summary,
			"country":        req.Country,
			"coverResources": req.CoverResources,
			"canReview":      req.CanReview,
			"tags":           tags,
			"lastModifyTime": now,
		})
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (h *contentHandler) fillInfo(sessionContext mongo.SessionContext, uid string, info *models.Info) (*content.InfoResp, error) {
	isThumbUp, err := h.thumbUpRepo.IsExists(sessionContext, uid, info.ID)
	if err != nil {
		return nil, err
	}

	isThumbDown, err := h.thumbDownRepo.IsExists(sessionContext, uid, info.ID)
	if err != nil {
		return nil, err
	}

	isFavorite, err := h.favoriteRepo.IsExists(sessionContext, uid, info.ID)
	if err != nil {
		return nil, err
	}

	thumbUps, err := h.thumbUpRepo.CountByInfo(sessionContext, info.ID)
	if err != nil {
		return nil, err
	}

	thumbDowns, err := h.thumbDownRepo.CountByInfo(sessionContext, info.ID)
	if err != nil {
		return nil, err
	}

	favorites, err := h.favoriteRepo.CountByInfo(sessionContext, info.ID)
	if err != nil {
		return nil, err
	}

	tmpList, err := h.thumbUpRepo.FindUserByInfo(sessionContext, info.ID, 0, 10, []*content.Sort{
		{
			Name: "time",
			Asc:  false,
		},
	})
	if err != nil {
		return nil, err
	}
	thumbUpList := make([]string, 0, len(tmpList))
	for _, v := range tmpList {
		thumbUpList = append(thumbUpList, v.UID)
	}

	tmpList, err = h.thumbDownRepo.FindUserByInfo(sessionContext, info.ID, 0, 10, []*content.Sort{
		{
			Name: "time",
			Asc:  false,
		},
	})
	if err != nil {
		return nil, err
	}
	thumbDownList := make([]string, 0, len(tmpList))
	for _, v := range tmpList {
		thumbDownList = append(thumbDownList, v.UID)
	}

	tmpList, err = h.favoriteRepo.FindUserByInfo(sessionContext, info.ID, 0, 10, []*content.Sort{
		{
			Name: "time",
			Asc:  false,
		},
	})
	if err != nil {
		return nil, err
	}
	favoriteList := make([]string, 0, len(tmpList))
	for _, v := range tmpList {
		favoriteList = append(favoriteList, v.UID)
	}

	tagList := make([]*content.TagAndType, 0, len(info.Tags))
	for _, v := range info.Tags {
		tagList = append(tagList, &content.TagAndType{
			Tag:  v.Tag,
			Type: v.Type,
		})
	}

	resp := &content.InfoResp{
		InfoID:  info.ID.Hex(),
		Uid:     info.UID,
		Title:   info.Title,
		Author:  info.Author,
		Summary: info.Summary,
		Country: info.Country,
		ContentTime: &timestamp.Timestamp{
			Seconds: info.ContentTime.Unix(),
			Nanos:   int32(info.ContentTime.Nanosecond()),
		},
		CoverResources: info.CoverResources,
		PublishTime: &timestamp.Timestamp{
			Seconds: info.PublishTime.Unix(),
			Nanos:   int32(info.PublishTime.Nanosecond()),
		},
		LastReviewTime: &timestamp.Timestamp{
			Seconds: info.LastReviewTime.Unix(),
			Nanos:   int32(info.LastReviewTime.Nanosecond()),
		},
		Valid:         info.Valid,
		WatchCount:    info.WatchCount,
		Tags:          tagList,
		ThumbUps:      thumbUps,
		IsThumbUp:     isThumbUp,
		ThumbUpList:   thumbUpList,
		ThumbDowns:    thumbDowns,
		IsThumbDown:   isThumbDown,
		ThumbDownList: thumbDownList,
		Favorites:     favorites,
		IsFavorite:    isFavorite,
		FavoriteList:  favoriteList,
		LastModifyTime: &timestamp.Timestamp{
			Seconds: info.LatestModifyTime.Unix(),
			Nanos:   int32(info.LatestModifyTime.Nanosecond()),
		},
		CanReview:       info.CanReview,
		Archived:        info.Archived,
		LatestSegmentID: info.LatestSegmentID.Hex(),
		SegmentCount:    info.SegmentCount,
	}
	return resp, nil
}

func (h *contentHandler) GetInfo(ctx context.Context, req *content.GetInfoReq) (*content.InfoResp, error) {
	var resp *content.InfoResp
	var err error
	if err = validateGetInfoReq(req); err != nil {
		return nil, err
	}

	infoID, err := objectid.FromHex(req.InfoID)
	if err != nil {
		return nil, err
	}

	err = h.client.UseSession(ctx, func(sessionContext mongo.SessionContext) error {
		info, err := h.infoRepo.FindOne(sessionContext, infoID)
		if err != nil {
			return err
		}

		resp, err = h.fillInfo(sessionContext, req.Uid, info)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (h *contentHandler) GetInfos(ctx context.Context, req *content.GetInfosReq) (*content.InfosResp, error) {
	var resp content.InfosResp
	if err := validateGetInfosReq(req); err != nil {
		return nil, err
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	err := h.client.UseSession(timeoutCtx, func(sessionContext mongo.SessionContext) error {
		var err error
		tags := make([]*models.TypeAndTag, 0, len(req.Tags))
		for _, v := range req.Tags {
			tags = append(tags, &models.TypeAndTag{
				Type: v.Type,
				Tag:  v.Tag,
			})
		}

		startTime, err := ptypes.Timestamp(req.StartTime)
		if err != nil {
			return err
		}

		endTime, err := ptypes.Timestamp(req.EndTime)
		if err != nil {
			return err
		}

		var infos []*models.Info
		var totalCount uint64
		if req.Title == "" {
			infos, totalCount, err = h.infoRepo.FindAll(sessionContext, "", req.Country,
				&startTime, &endTime, tags, req.Page, req.Size, req.Sorts)
		} else {
			//TODO: search use elasticsearch
		}
		if err != nil {
			return err
		}

		pInfos := make([]*content.InfoResp, 0, len(infos))
		for _, info := range infos {
			pInfo, err := h.fillInfo(sessionContext, req.Uid, info)
			if err != nil {
				return err
			}
			pInfos = append(pInfos, pInfo)
		}
		resp.Items = pInfos
		resp.TotalCount = totalCount
		return nil
	})

	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (h *contentHandler) DeleteInfo(ctx context.Context, req *content.InfoOneReq) (*empty.Empty, error) {
	var resp empty.Empty
	if err := validateInfoOneReq(req); err != nil {
		return nil, err
	}

	err := h.client.UseSession(ctx, func(sessionContext mongo.SessionContext) error {
		var err error
		infoID, err := objectid.FromHex(req.InfoID)
		if err != nil {
			return err
		}

		err = h.infoRepo.Delete(sessionContext, infoID)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (h *contentHandler) WatchInfo(ctx context.Context, req *content.InfoOneReq) (*empty.Empty, error) {
	var resp empty.Empty
	if err := validateInfoOneReq(req); err != nil {
		return nil, err
	}

	err := h.client.UseSession(ctx, func(sessionContext mongo.SessionContext) error {
		var err error
		infoID, err := objectid.FromHex(req.InfoID)
		if err != nil {
			return err
		}

		err = h.infoRepo.IncWatchCount(sessionContext, infoID, 1)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (h *contentHandler) _behaviorInsert(ctx context.Context,
	checkRepo, repo repositories.BehaviorRepository, req *content.InfoIDWithUIDReq) (*empty.Empty, error) {
	var resp empty.Empty

	if err := validateInfoIDWithUIDReq(req); err != nil {
		return nil, err
	}

	infoID, err := objectid.FromHex(req.InfoID)
	if err != nil {
		return nil, err
	}

	err = h.client.UseSession(ctx, func(sessionContext mongo.SessionContext) error {
		var err error
		if checkRepo != nil {
			isExist, err := checkRepo.IsExists(sessionContext, req.Uid, infoID)
			if err != nil {
				return err
			}

			if isExist {
				return ErrBehaviorExists
			}
		}

		isExist, err := repo.IsExists(sessionContext, req.Uid, infoID)
		if err != nil {
			return err
		}

		if isExist {
			return ErrBehaviorExists
		}

		err = repo.Insert(sessionContext, req.Uid, &models.BehaviorInfoItem{
			InfoId: infoID,
			Time:   time.Now(),
		})
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func (h *contentHandler) _behaviorFindInfoByUser(ctx context.Context,
	repo repositories.BehaviorRepository, req *content.UIDPageReq) (*content.InfoIDsResp, error) {
	var resp content.InfoIDsResp
	if err := validateUIDPageReq(req); err != nil {
		return nil, err
	}

	err := h.client.UseSession(ctx, func(sessionContext mongo.SessionContext) error {
		items, err := repo.FindInfoByUser(sessionContext, req.Uid, req.Page, req.Size, req.Sorts)

		if err != nil {
			return err
		}

		results := make([]*content.InfoIDWithTime, 0, len(items))
		for _, item := range items {
			results = append(results, &content.InfoIDWithTime{
				InfoId: item.InfoId.Hex(),
				Time: &timestamp.Timestamp{
					Seconds: item.Time.Unix(),
					Nanos:   int32(item.Time.Nanosecond()),
				},
			})
		}
		resp.Items = results
		return nil
	})

	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func (h *contentHandler) _behaviorFindUserByInfo(ctx context.Context,
	repo repositories.BehaviorRepository, req *content.InfoIDPageReq) (*content.UserIDsResp, error) {
	var resp content.UserIDsResp
	if err := validateInfoIDPageReq(req); err != nil {
		return nil, err
	}

	infoID, err := objectid.FromHex(req.InfoID)
	if err != nil {
		return nil, err
	}

	err = h.client.UseSession(ctx, func(sessionContext mongo.SessionContext) error {
		items, err := repo.FindUserByInfo(sessionContext, infoID, req.Page, req.Size, req.Sorts)

		if err != nil {
			return err
		}

		results := make([]*content.UIDWithTime, 0, len(items))
		for _, item := range items {
			results = append(results, &content.UIDWithTime{
				Uid: item.UID,
				Time: &timestamp.Timestamp{
					Seconds: item.Time.Unix(),
					Nanos:   int32(item.Time.Nanosecond()),
				},
			})
		}
		resp.Items = results
		return nil
	})

	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func (h *contentHandler) _behaviorDelete(ctx context.Context,
	repo repositories.BehaviorRepository, req *content.InfoIDWithUIDReq) (*empty.Empty, error) {
	var resp empty.Empty
	if err := validateInfoIDWithUIDReq(req); err != nil {
		return nil, err
	}

	infoID, err := objectid.FromHex(req.InfoID)
	if err != nil {
		return nil, err
	}

	err = h.client.UseSession(ctx, func(sessionContext mongo.SessionContext) error {
		err = repo.Delete(sessionContext, req.Uid, infoID)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (h *contentHandler) ThumbUp(ctx context.Context, req *content.InfoIDWithUIDReq) (*empty.Empty, error) {
	return h._behaviorInsert(ctx, h.thumbDownRepo, h.thumbUpRepo, req)
}

func (h *contentHandler) ThumbDown(ctx context.Context, req *content.InfoIDWithUIDReq) (*empty.Empty, error) {
	return h._behaviorInsert(ctx, h.thumbUpRepo, h.thumbDownRepo, req)
}

func (h *contentHandler) GetUserThumbUp(ctx context.Context, req *content.UIDPageReq) (*content.InfoIDsResp, error) {
	return h._behaviorFindInfoByUser(ctx, h.thumbUpRepo, req)
}

func (h *contentHandler) GetUserThumbDown(ctx context.Context, req *content.UIDPageReq) (*content.InfoIDsResp, error) {
	return h._behaviorFindInfoByUser(ctx, h.thumbDownRepo, req)
}

func (h *contentHandler) GetInfoThumbUp(ctx context.Context, req *content.InfoIDPageReq) (*content.UserIDsResp, error) {
	return h._behaviorFindUserByInfo(ctx, h.thumbUpRepo, req)
}

func (h *contentHandler) GetInfoThumbDown(ctx context.Context, req *content.InfoIDPageReq) (*content.UserIDsResp, error) {
	return h._behaviorFindUserByInfo(ctx, h.thumbDownRepo, req)
}

func (h *contentHandler) DeleteThumbUp(ctx context.Context, req *content.InfoIDWithUIDReq) (*empty.Empty, error) {
	return h._behaviorDelete(ctx, h.thumbUpRepo, req)
}

func (h *contentHandler) DeleteThumbDown(ctx context.Context, req *content.InfoIDWithUIDReq) (*empty.Empty, error) {
	return h._behaviorDelete(ctx, h.thumbDownRepo, req)
}

func (h *contentHandler) Favorite(ctx context.Context, req *content.InfoIDWithUIDReq) (*empty.Empty, error) {
	return h._behaviorInsert(ctx, nil, h.favoriteRepo, req)
}

func (h *contentHandler) GetUserFavorite(ctx context.Context, req *content.UIDPageReq) (*content.InfoIDsResp, error) {
	return h._behaviorFindInfoByUser(ctx, h.favoriteRepo, req)
}

func (h *contentHandler) GetInfoFavorite(ctx context.Context, req *content.InfoIDPageReq) (*content.UserIDsResp, error) {
	return h._behaviorFindUserByInfo(ctx, h.favoriteRepo, req)
}

func (h *contentHandler) DeleteFavorite(ctx context.Context, req *content.InfoIDWithUIDReq) (*empty.Empty, error) {
	return h._behaviorDelete(ctx, h.favoriteRepo, req)
}
