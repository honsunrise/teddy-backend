// Code generated by protoc-gen-go. DO NOT EDIT.
// source: proto/content.proto

package proto

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import _ "github.com/golang/protobuf/ptypes/empty"
import timestamp "github.com/golang/protobuf/ptypes/timestamp"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type Tag struct {
	Tag                  string               `protobuf:"bytes,1,opt,name=tag,proto3" json:"tag,omitempty"`
	Hot                  uint64               `protobuf:"varint,2,opt,name=hot,proto3" json:"hot,omitempty"`
	FirstShow            *timestamp.Timestamp `protobuf:"bytes,3,opt,name=firstShow,proto3" json:"firstShow,omitempty"`
	LastUse              *timestamp.Timestamp `protobuf:"bytes,4,opt,name=lastUse,proto3" json:"lastUse,omitempty"`
	LastUseBy            string               `protobuf:"bytes,5,opt,name=lastUseBy,proto3" json:"lastUseBy,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *Tag) Reset()         { *m = Tag{} }
func (m *Tag) String() string { return proto.CompactTextString(m) }
func (*Tag) ProtoMessage()    {}
func (*Tag) Descriptor() ([]byte, []int) {
	return fileDescriptor_content_85356c321e0b0cd2, []int{0}
}
func (m *Tag) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Tag.Unmarshal(m, b)
}
func (m *Tag) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Tag.Marshal(b, m, deterministic)
}
func (dst *Tag) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Tag.Merge(dst, src)
}
func (m *Tag) XXX_Size() int {
	return xxx_messageInfo_Tag.Size(m)
}
func (m *Tag) XXX_DiscardUnknown() {
	xxx_messageInfo_Tag.DiscardUnknown(m)
}

var xxx_messageInfo_Tag proto.InternalMessageInfo

func (m *Tag) GetTag() string {
	if m != nil {
		return m.Tag
	}
	return ""
}

func (m *Tag) GetHot() uint64 {
	if m != nil {
		return m.Hot
	}
	return 0
}

func (m *Tag) GetFirstShow() *timestamp.Timestamp {
	if m != nil {
		return m.FirstShow
	}
	return nil
}

func (m *Tag) GetLastUse() *timestamp.Timestamp {
	if m != nil {
		return m.LastUse
	}
	return nil
}

func (m *Tag) GetLastUseBy() string {
	if m != nil {
		return m.LastUseBy
	}
	return ""
}

type Info struct {
	Uid                  string               `protobuf:"bytes,1,opt,name=uid,proto3" json:"uid,omitempty"`
	Id                   string               `protobuf:"bytes,2,opt,name=id,proto3" json:"id,omitempty"`
	Title                string               `protobuf:"bytes,3,opt,name=title,proto3" json:"title,omitempty"`
	Type                 uint32               `protobuf:"varint,4,opt,name=type,proto3" json:"type,omitempty"`
	Content              string               `protobuf:"bytes,5,opt,name=content,proto3" json:"content,omitempty"`
	CoverList            []string             `protobuf:"bytes,6,rep,name=coverList,proto3" json:"coverList,omitempty"`
	CoverVideo           string               `protobuf:"bytes,7,opt,name=coverVideo,proto3" json:"coverVideo,omitempty"`
	PublishTime          *timestamp.Timestamp `protobuf:"bytes,8,opt,name=publishTime,proto3" json:"publishTime,omitempty"`
	LastReviewTime       *timestamp.Timestamp `protobuf:"bytes,9,opt,name=lastReviewTime,proto3" json:"lastReviewTime,omitempty"`
	Valid                bool                 `protobuf:"varint,10,opt,name=valid,proto3" json:"valid,omitempty"`
	WatchCount           int64                `protobuf:"varint,11,opt,name=watchCount,proto3" json:"watchCount,omitempty"`
	Tags                 []string             `protobuf:"bytes,12,rep,name=tags,proto3" json:"tags,omitempty"`
	Likes                int64                `protobuf:"varint,13,opt,name=likes,proto3" json:"likes,omitempty"`
	IsLike               bool                 `protobuf:"varint,14,opt,name=isLike,proto3" json:"isLike,omitempty"`
	LikeList             []string             `protobuf:"bytes,15,rep,name=likeList,proto3" json:"likeList,omitempty"`
	Unlike               int64                `protobuf:"varint,16,opt,name=unlike,proto3" json:"unlike,omitempty"`
	IsUnlike             bool                 `protobuf:"varint,17,opt,name=isUnlike,proto3" json:"isUnlike,omitempty"`
	UnlikeList           []string             `protobuf:"bytes,18,rep,name=unlikeList,proto3" json:"unlikeList,omitempty"`
	Favorites            int64                `protobuf:"varint,19,opt,name=favorites,proto3" json:"favorites,omitempty"`
	IsFavorite           bool                 `protobuf:"varint,20,opt,name=isFavorite,proto3" json:"isFavorite,omitempty"`
	FavoriteList         []string             `protobuf:"bytes,21,rep,name=favoriteList,proto3" json:"favoriteList,omitempty"`
	LastModifyTime       *timestamp.Timestamp `protobuf:"bytes,22,opt,name=lastModifyTime,proto3" json:"lastModifyTime,omitempty"`
	CanReview            bool                 `protobuf:"varint,23,opt,name=canReview,proto3" json:"canReview,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *Info) Reset()         { *m = Info{} }
func (m *Info) String() string { return proto.CompactTextString(m) }
func (*Info) ProtoMessage()    {}
func (*Info) Descriptor() ([]byte, []int) {
	return fileDescriptor_content_85356c321e0b0cd2, []int{1}
}
func (m *Info) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Info.Unmarshal(m, b)
}
func (m *Info) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Info.Marshal(b, m, deterministic)
}
func (dst *Info) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Info.Merge(dst, src)
}
func (m *Info) XXX_Size() int {
	return xxx_messageInfo_Info.Size(m)
}
func (m *Info) XXX_DiscardUnknown() {
	xxx_messageInfo_Info.DiscardUnknown(m)
}

var xxx_messageInfo_Info proto.InternalMessageInfo

func (m *Info) GetUid() string {
	if m != nil {
		return m.Uid
	}
	return ""
}

func (m *Info) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *Info) GetTitle() string {
	if m != nil {
		return m.Title
	}
	return ""
}

func (m *Info) GetType() uint32 {
	if m != nil {
		return m.Type
	}
	return 0
}

func (m *Info) GetContent() string {
	if m != nil {
		return m.Content
	}
	return ""
}

func (m *Info) GetCoverList() []string {
	if m != nil {
		return m.CoverList
	}
	return nil
}

func (m *Info) GetCoverVideo() string {
	if m != nil {
		return m.CoverVideo
	}
	return ""
}

func (m *Info) GetPublishTime() *timestamp.Timestamp {
	if m != nil {
		return m.PublishTime
	}
	return nil
}

func (m *Info) GetLastReviewTime() *timestamp.Timestamp {
	if m != nil {
		return m.LastReviewTime
	}
	return nil
}

func (m *Info) GetValid() bool {
	if m != nil {
		return m.Valid
	}
	return false
}

func (m *Info) GetWatchCount() int64 {
	if m != nil {
		return m.WatchCount
	}
	return 0
}

func (m *Info) GetTags() []string {
	if m != nil {
		return m.Tags
	}
	return nil
}

func (m *Info) GetLikes() int64 {
	if m != nil {
		return m.Likes
	}
	return 0
}

func (m *Info) GetIsLike() bool {
	if m != nil {
		return m.IsLike
	}
	return false
}

func (m *Info) GetLikeList() []string {
	if m != nil {
		return m.LikeList
	}
	return nil
}

func (m *Info) GetUnlike() int64 {
	if m != nil {
		return m.Unlike
	}
	return 0
}

func (m *Info) GetIsUnlike() bool {
	if m != nil {
		return m.IsUnlike
	}
	return false
}

func (m *Info) GetUnlikeList() []string {
	if m != nil {
		return m.UnlikeList
	}
	return nil
}

func (m *Info) GetFavorites() int64 {
	if m != nil {
		return m.Favorites
	}
	return 0
}

func (m *Info) GetIsFavorite() bool {
	if m != nil {
		return m.IsFavorite
	}
	return false
}

func (m *Info) GetFavoriteList() []string {
	if m != nil {
		return m.FavoriteList
	}
	return nil
}

func (m *Info) GetLastModifyTime() *timestamp.Timestamp {
	if m != nil {
		return m.LastModifyTime
	}
	return nil
}

func (m *Info) GetCanReview() bool {
	if m != nil {
		return m.CanReview
	}
	return false
}

type GetTagReq struct {
	Page                 uint32   `protobuf:"varint,1,opt,name=page,proto3" json:"page,omitempty"`
	Size                 uint32   `protobuf:"varint,2,opt,name=size,proto3" json:"size,omitempty"`
	Uid                  string   `protobuf:"bytes,3,opt,name=uid,proto3" json:"uid,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetTagReq) Reset()         { *m = GetTagReq{} }
func (m *GetTagReq) String() string { return proto.CompactTextString(m) }
func (*GetTagReq) ProtoMessage()    {}
func (*GetTagReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_content_85356c321e0b0cd2, []int{2}
}
func (m *GetTagReq) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetTagReq.Unmarshal(m, b)
}
func (m *GetTagReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetTagReq.Marshal(b, m, deterministic)
}
func (dst *GetTagReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetTagReq.Merge(dst, src)
}
func (m *GetTagReq) XXX_Size() int {
	return xxx_messageInfo_GetTagReq.Size(m)
}
func (m *GetTagReq) XXX_DiscardUnknown() {
	xxx_messageInfo_GetTagReq.DiscardUnknown(m)
}

var xxx_messageInfo_GetTagReq proto.InternalMessageInfo

func (m *GetTagReq) GetPage() uint32 {
	if m != nil {
		return m.Page
	}
	return 0
}

func (m *GetTagReq) GetSize() uint32 {
	if m != nil {
		return m.Size
	}
	return 0
}

func (m *GetTagReq) GetUid() string {
	if m != nil {
		return m.Uid
	}
	return ""
}

type UidPageReq struct {
	Page                 uint32   `protobuf:"varint,1,opt,name=page,proto3" json:"page,omitempty"`
	Size                 uint32   `protobuf:"varint,2,opt,name=size,proto3" json:"size,omitempty"`
	Uid                  string   `protobuf:"bytes,3,opt,name=uid,proto3" json:"uid,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *UidPageReq) Reset()         { *m = UidPageReq{} }
func (m *UidPageReq) String() string { return proto.CompactTextString(m) }
func (*UidPageReq) ProtoMessage()    {}
func (*UidPageReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_content_85356c321e0b0cd2, []int{3}
}
func (m *UidPageReq) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UidPageReq.Unmarshal(m, b)
}
func (m *UidPageReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UidPageReq.Marshal(b, m, deterministic)
}
func (dst *UidPageReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UidPageReq.Merge(dst, src)
}
func (m *UidPageReq) XXX_Size() int {
	return xxx_messageInfo_UidPageReq.Size(m)
}
func (m *UidPageReq) XXX_DiscardUnknown() {
	xxx_messageInfo_UidPageReq.DiscardUnknown(m)
}

var xxx_messageInfo_UidPageReq proto.InternalMessageInfo

func (m *UidPageReq) GetPage() uint32 {
	if m != nil {
		return m.Page
	}
	return 0
}

func (m *UidPageReq) GetSize() uint32 {
	if m != nil {
		return m.Size
	}
	return 0
}

func (m *UidPageReq) GetUid() string {
	if m != nil {
		return m.Uid
	}
	return ""
}

type GetTagsResp struct {
	Tags                 []*Tag   `protobuf:"bytes,1,rep,name=tags,proto3" json:"tags,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetTagsResp) Reset()         { *m = GetTagsResp{} }
func (m *GetTagsResp) String() string { return proto.CompactTextString(m) }
func (*GetTagsResp) ProtoMessage()    {}
func (*GetTagsResp) Descriptor() ([]byte, []int) {
	return fileDescriptor_content_85356c321e0b0cd2, []int{4}
}
func (m *GetTagsResp) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetTagsResp.Unmarshal(m, b)
}
func (m *GetTagsResp) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetTagsResp.Marshal(b, m, deterministic)
}
func (dst *GetTagsResp) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetTagsResp.Merge(dst, src)
}
func (m *GetTagsResp) XXX_Size() int {
	return xxx_messageInfo_GetTagsResp.Size(m)
}
func (m *GetTagsResp) XXX_DiscardUnknown() {
	xxx_messageInfo_GetTagsResp.DiscardUnknown(m)
}

var xxx_messageInfo_GetTagsResp proto.InternalMessageInfo

func (m *GetTagsResp) GetTags() []*Tag {
	if m != nil {
		return m.Tags
	}
	return nil
}

type PublishInfoReq struct {
	Uid                  string            `protobuf:"bytes,1,opt,name=uid,proto3" json:"uid,omitempty"`
	Title                string            `protobuf:"bytes,2,opt,name=title,proto3" json:"title,omitempty"`
	Content              string            `protobuf:"bytes,3,opt,name=content,proto3" json:"content,omitempty"`
	Tags                 []string          `protobuf:"bytes,4,rep,name=tags,proto3" json:"tags,omitempty"`
	ImgList              map[string]string `protobuf:"bytes,5,rep,name=imgList,proto3" json:"imgList,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	VideoList            map[string]string `protobuf:"bytes,6,rep,name=videoList,proto3" json:"videoList,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	External             bool              `protobuf:"varint,7,opt,name=external,proto3" json:"external,omitempty"`
	CanReview            bool              `protobuf:"varint,8,opt,name=canReview,proto3" json:"canReview,omitempty"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *PublishInfoReq) Reset()         { *m = PublishInfoReq{} }
func (m *PublishInfoReq) String() string { return proto.CompactTextString(m) }
func (*PublishInfoReq) ProtoMessage()    {}
func (*PublishInfoReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_content_85356c321e0b0cd2, []int{5}
}
func (m *PublishInfoReq) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PublishInfoReq.Unmarshal(m, b)
}
func (m *PublishInfoReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PublishInfoReq.Marshal(b, m, deterministic)
}
func (dst *PublishInfoReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PublishInfoReq.Merge(dst, src)
}
func (m *PublishInfoReq) XXX_Size() int {
	return xxx_messageInfo_PublishInfoReq.Size(m)
}
func (m *PublishInfoReq) XXX_DiscardUnknown() {
	xxx_messageInfo_PublishInfoReq.DiscardUnknown(m)
}

var xxx_messageInfo_PublishInfoReq proto.InternalMessageInfo

func (m *PublishInfoReq) GetUid() string {
	if m != nil {
		return m.Uid
	}
	return ""
}

func (m *PublishInfoReq) GetTitle() string {
	if m != nil {
		return m.Title
	}
	return ""
}

func (m *PublishInfoReq) GetContent() string {
	if m != nil {
		return m.Content
	}
	return ""
}

func (m *PublishInfoReq) GetTags() []string {
	if m != nil {
		return m.Tags
	}
	return nil
}

func (m *PublishInfoReq) GetImgList() map[string]string {
	if m != nil {
		return m.ImgList
	}
	return nil
}

func (m *PublishInfoReq) GetVideoList() map[string]string {
	if m != nil {
		return m.VideoList
	}
	return nil
}

func (m *PublishInfoReq) GetExternal() bool {
	if m != nil {
		return m.External
	}
	return false
}

func (m *PublishInfoReq) GetCanReview() bool {
	if m != nil {
		return m.CanReview
	}
	return false
}

type EditInfoReq struct {
	Id                   string            `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Uid                  string            `protobuf:"bytes,2,opt,name=uid,proto3" json:"uid,omitempty"`
	Title                string            `protobuf:"bytes,3,opt,name=title,proto3" json:"title,omitempty"`
	Content              string            `protobuf:"bytes,4,opt,name=content,proto3" json:"content,omitempty"`
	Tags                 []string          `protobuf:"bytes,5,rep,name=tags,proto3" json:"tags,omitempty"`
	ImgList              map[string]string `protobuf:"bytes,6,rep,name=imgList,proto3" json:"imgList,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	VideoList            map[string]string `protobuf:"bytes,7,rep,name=videoList,proto3" json:"videoList,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	External             bool              `protobuf:"varint,8,opt,name=external,proto3" json:"external,omitempty"`
	CanReview            bool              `protobuf:"varint,9,opt,name=canReview,proto3" json:"canReview,omitempty"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *EditInfoReq) Reset()         { *m = EditInfoReq{} }
func (m *EditInfoReq) String() string { return proto.CompactTextString(m) }
func (*EditInfoReq) ProtoMessage()    {}
func (*EditInfoReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_content_85356c321e0b0cd2, []int{6}
}
func (m *EditInfoReq) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_EditInfoReq.Unmarshal(m, b)
}
func (m *EditInfoReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_EditInfoReq.Marshal(b, m, deterministic)
}
func (dst *EditInfoReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EditInfoReq.Merge(dst, src)
}
func (m *EditInfoReq) XXX_Size() int {
	return xxx_messageInfo_EditInfoReq.Size(m)
}
func (m *EditInfoReq) XXX_DiscardUnknown() {
	xxx_messageInfo_EditInfoReq.DiscardUnknown(m)
}

var xxx_messageInfo_EditInfoReq proto.InternalMessageInfo

func (m *EditInfoReq) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *EditInfoReq) GetUid() string {
	if m != nil {
		return m.Uid
	}
	return ""
}

func (m *EditInfoReq) GetTitle() string {
	if m != nil {
		return m.Title
	}
	return ""
}

func (m *EditInfoReq) GetContent() string {
	if m != nil {
		return m.Content
	}
	return ""
}

func (m *EditInfoReq) GetTags() []string {
	if m != nil {
		return m.Tags
	}
	return nil
}

func (m *EditInfoReq) GetImgList() map[string]string {
	if m != nil {
		return m.ImgList
	}
	return nil
}

func (m *EditInfoReq) GetVideoList() map[string]string {
	if m != nil {
		return m.VideoList
	}
	return nil
}

func (m *EditInfoReq) GetExternal() bool {
	if m != nil {
		return m.External
	}
	return false
}

func (m *EditInfoReq) GetCanReview() bool {
	if m != nil {
		return m.CanReview
	}
	return false
}

type GetInfosReq struct {
	Page                 uint32   `protobuf:"varint,1,opt,name=page,proto3" json:"page,omitempty"`
	Size                 uint32   `protobuf:"varint,2,opt,name=size,proto3" json:"size,omitempty"`
	Uid                  string   `protobuf:"bytes,3,opt,name=uid,proto3" json:"uid,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetInfosReq) Reset()         { *m = GetInfosReq{} }
func (m *GetInfosReq) String() string { return proto.CompactTextString(m) }
func (*GetInfosReq) ProtoMessage()    {}
func (*GetInfosReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_content_85356c321e0b0cd2, []int{7}
}
func (m *GetInfosReq) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetInfosReq.Unmarshal(m, b)
}
func (m *GetInfosReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetInfosReq.Marshal(b, m, deterministic)
}
func (dst *GetInfosReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetInfosReq.Merge(dst, src)
}
func (m *GetInfosReq) XXX_Size() int {
	return xxx_messageInfo_GetInfosReq.Size(m)
}
func (m *GetInfosReq) XXX_DiscardUnknown() {
	xxx_messageInfo_GetInfosReq.DiscardUnknown(m)
}

var xxx_messageInfo_GetInfosReq proto.InternalMessageInfo

func (m *GetInfosReq) GetPage() uint32 {
	if m != nil {
		return m.Page
	}
	return 0
}

func (m *GetInfosReq) GetSize() uint32 {
	if m != nil {
		return m.Size
	}
	return 0
}

func (m *GetInfosReq) GetUid() string {
	if m != nil {
		return m.Uid
	}
	return ""
}

type GetInfosResp struct {
	Infos                []*Info  `protobuf:"bytes,1,rep,name=infos,proto3" json:"infos,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetInfosResp) Reset()         { *m = GetInfosResp{} }
func (m *GetInfosResp) String() string { return proto.CompactTextString(m) }
func (*GetInfosResp) ProtoMessage()    {}
func (*GetInfosResp) Descriptor() ([]byte, []int) {
	return fileDescriptor_content_85356c321e0b0cd2, []int{8}
}
func (m *GetInfosResp) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetInfosResp.Unmarshal(m, b)
}
func (m *GetInfosResp) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetInfosResp.Marshal(b, m, deterministic)
}
func (dst *GetInfosResp) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetInfosResp.Merge(dst, src)
}
func (m *GetInfosResp) XXX_Size() int {
	return xxx_messageInfo_GetInfosResp.Size(m)
}
func (m *GetInfosResp) XXX_DiscardUnknown() {
	xxx_messageInfo_GetInfosResp.DiscardUnknown(m)
}

var xxx_messageInfo_GetInfosResp proto.InternalMessageInfo

func (m *GetInfosResp) GetInfos() []*Info {
	if m != nil {
		return m.Infos
	}
	return nil
}

type InfoIdReq struct {
	Id                   string   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *InfoIdReq) Reset()         { *m = InfoIdReq{} }
func (m *InfoIdReq) String() string { return proto.CompactTextString(m) }
func (*InfoIdReq) ProtoMessage()    {}
func (*InfoIdReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_content_85356c321e0b0cd2, []int{9}
}
func (m *InfoIdReq) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_InfoIdReq.Unmarshal(m, b)
}
func (m *InfoIdReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_InfoIdReq.Marshal(b, m, deterministic)
}
func (dst *InfoIdReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_InfoIdReq.Merge(dst, src)
}
func (m *InfoIdReq) XXX_Size() int {
	return xxx_messageInfo_InfoIdReq.Size(m)
}
func (m *InfoIdReq) XXX_DiscardUnknown() {
	xxx_messageInfo_InfoIdReq.DiscardUnknown(m)
}

var xxx_messageInfo_InfoIdReq proto.InternalMessageInfo

func (m *InfoIdReq) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

type InfoIdPageReq struct {
	Id                   string   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Size                 uint32   `protobuf:"varint,2,opt,name=size,proto3" json:"size,omitempty"`
	Uid                  string   `protobuf:"bytes,3,opt,name=uid,proto3" json:"uid,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *InfoIdPageReq) Reset()         { *m = InfoIdPageReq{} }
func (m *InfoIdPageReq) String() string { return proto.CompactTextString(m) }
func (*InfoIdPageReq) ProtoMessage()    {}
func (*InfoIdPageReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_content_85356c321e0b0cd2, []int{10}
}
func (m *InfoIdPageReq) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_InfoIdPageReq.Unmarshal(m, b)
}
func (m *InfoIdPageReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_InfoIdPageReq.Marshal(b, m, deterministic)
}
func (dst *InfoIdPageReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_InfoIdPageReq.Merge(dst, src)
}
func (m *InfoIdPageReq) XXX_Size() int {
	return xxx_messageInfo_InfoIdPageReq.Size(m)
}
func (m *InfoIdPageReq) XXX_DiscardUnknown() {
	xxx_messageInfo_InfoIdPageReq.DiscardUnknown(m)
}

var xxx_messageInfo_InfoIdPageReq proto.InternalMessageInfo

func (m *InfoIdPageReq) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *InfoIdPageReq) GetSize() uint32 {
	if m != nil {
		return m.Size
	}
	return 0
}

func (m *InfoIdPageReq) GetUid() string {
	if m != nil {
		return m.Uid
	}
	return ""
}

type UserIdsResp struct {
	Uids                 []string `protobuf:"bytes,1,rep,name=uids,proto3" json:"uids,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *UserIdsResp) Reset()         { *m = UserIdsResp{} }
func (m *UserIdsResp) String() string { return proto.CompactTextString(m) }
func (*UserIdsResp) ProtoMessage()    {}
func (*UserIdsResp) Descriptor() ([]byte, []int) {
	return fileDescriptor_content_85356c321e0b0cd2, []int{11}
}
func (m *UserIdsResp) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UserIdsResp.Unmarshal(m, b)
}
func (m *UserIdsResp) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UserIdsResp.Marshal(b, m, deterministic)
}
func (dst *UserIdsResp) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UserIdsResp.Merge(dst, src)
}
func (m *UserIdsResp) XXX_Size() int {
	return xxx_messageInfo_UserIdsResp.Size(m)
}
func (m *UserIdsResp) XXX_DiscardUnknown() {
	xxx_messageInfo_UserIdsResp.DiscardUnknown(m)
}

var xxx_messageInfo_UserIdsResp proto.InternalMessageInfo

func (m *UserIdsResp) GetUids() []string {
	if m != nil {
		return m.Uids
	}
	return nil
}

type InfoIdsResp struct {
	Ids                  []string `protobuf:"bytes,1,rep,name=ids,proto3" json:"ids,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *InfoIdsResp) Reset()         { *m = InfoIdsResp{} }
func (m *InfoIdsResp) String() string { return proto.CompactTextString(m) }
func (*InfoIdsResp) ProtoMessage()    {}
func (*InfoIdsResp) Descriptor() ([]byte, []int) {
	return fileDescriptor_content_85356c321e0b0cd2, []int{12}
}
func (m *InfoIdsResp) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_InfoIdsResp.Unmarshal(m, b)
}
func (m *InfoIdsResp) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_InfoIdsResp.Marshal(b, m, deterministic)
}
func (dst *InfoIdsResp) XXX_Merge(src proto.Message) {
	xxx_messageInfo_InfoIdsResp.Merge(dst, src)
}
func (m *InfoIdsResp) XXX_Size() int {
	return xxx_messageInfo_InfoIdsResp.Size(m)
}
func (m *InfoIdsResp) XXX_DiscardUnknown() {
	xxx_messageInfo_InfoIdsResp.DiscardUnknown(m)
}

var xxx_messageInfo_InfoIdsResp proto.InternalMessageInfo

func (m *InfoIdsResp) GetIds() []string {
	if m != nil {
		return m.Ids
	}
	return nil
}

func init() {
	proto.RegisterType((*Tag)(nil), "com.teddy.srv.content.Tag")
	proto.RegisterType((*Info)(nil), "com.teddy.srv.content.Info")
	proto.RegisterType((*GetTagReq)(nil), "com.teddy.srv.content.GetTagReq")
	proto.RegisterType((*UidPageReq)(nil), "com.teddy.srv.content.UidPageReq")
	proto.RegisterType((*GetTagsResp)(nil), "com.teddy.srv.content.GetTagsResp")
	proto.RegisterType((*PublishInfoReq)(nil), "com.teddy.srv.content.PublishInfoReq")
	proto.RegisterMapType((map[string]string)(nil), "com.teddy.srv.content.PublishInfoReq.ImgListEntry")
	proto.RegisterMapType((map[string]string)(nil), "com.teddy.srv.content.PublishInfoReq.VideoListEntry")
	proto.RegisterType((*EditInfoReq)(nil), "com.teddy.srv.content.EditInfoReq")
	proto.RegisterMapType((map[string]string)(nil), "com.teddy.srv.content.EditInfoReq.ImgListEntry")
	proto.RegisterMapType((map[string]string)(nil), "com.teddy.srv.content.EditInfoReq.VideoListEntry")
	proto.RegisterType((*GetInfosReq)(nil), "com.teddy.srv.content.GetInfosReq")
	proto.RegisterType((*GetInfosResp)(nil), "com.teddy.srv.content.GetInfosResp")
	proto.RegisterType((*InfoIdReq)(nil), "com.teddy.srv.content.InfoIdReq")
	proto.RegisterType((*InfoIdPageReq)(nil), "com.teddy.srv.content.InfoIdPageReq")
	proto.RegisterType((*UserIdsResp)(nil), "com.teddy.srv.content.UserIdsResp")
	proto.RegisterType((*InfoIdsResp)(nil), "com.teddy.srv.content.InfoIdsResp")
}

func init() { proto.RegisterFile("proto/content.proto", fileDescriptor_content_85356c321e0b0cd2) }

var fileDescriptor_content_85356c321e0b0cd2 = []byte{
	// 1072 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xcc, 0x56, 0x5f, 0x73, 0xdb, 0x44,
	0x10, 0xc7, 0x96, 0x5d, 0x5b, 0xeb, 0xc4, 0x0d, 0xd7, 0x36, 0x68, 0x5c, 0x06, 0x52, 0x01, 0x33,
	0x3c, 0xc9, 0xd3, 0xb4, 0x0f, 0x9d, 0x52, 0x1e, 0x48, 0x49, 0xd2, 0x30, 0x01, 0x8a, 0x68, 0xda,
	0x21, 0xc3, 0x30, 0xa3, 0x44, 0x67, 0xf9, 0xa6, 0xb6, 0x65, 0xac, 0xb3, 0x83, 0xf9, 0x62, 0xbc,
	0xf0, 0xc6, 0x97, 0xe0, 0x23, 0xf0, 0x31, 0xb8, 0xdd, 0xd3, 0x49, 0x72, 0x1a, 0xff, 0x61, 0xec,
	0x07, 0x5e, 0x92, 0xdd, 0xd5, 0xee, 0xef, 0x7e, 0xb7, 0xff, 0x7c, 0x70, 0x67, 0x38, 0x8a, 0x65,
	0xdc, 0xbe, 0x8c, 0x07, 0x92, 0x0f, 0xa4, 0x47, 0x1a, 0xbb, 0x77, 0x19, 0xf7, 0x3d, 0xc9, 0xc3,
	0x70, 0xea, 0x25, 0xa3, 0x89, 0x97, 0x7e, 0x6c, 0x3d, 0x8a, 0x84, 0xec, 0x8e, 0x2f, 0x94, 0xde,
	0x6f, 0x47, 0x71, 0x2f, 0x18, 0x44, 0x6d, 0xf2, 0xbf, 0x18, 0x77, 0xda, 0x43, 0x39, 0x1d, 0xf2,
	0xa4, 0xcd, 0xfb, 0x4a, 0xd0, 0x7f, 0x35, 0x56, 0xeb, 0x8b, 0xe5, 0x41, 0x52, 0xf4, 0x79, 0x22,
	0x83, 0xfe, 0x30, 0x97, 0x74, 0xb0, 0xfb, 0x47, 0x09, 0xac, 0x57, 0x41, 0xc4, 0x76, 0xc0, 0x92,
	0x41, 0xe4, 0x94, 0xf6, 0x4a, 0x9f, 0xdb, 0x3e, 0x8a, 0x68, 0xe9, 0xc6, 0xd2, 0x29, 0x2b, 0x4b,
	0xc5, 0x47, 0x91, 0x3d, 0x01, 0xbb, 0x23, 0x46, 0x89, 0xfc, 0xb1, 0x1b, 0x5f, 0x39, 0x96, 0xb2,
	0x37, 0xf6, 0x5b, 0x5e, 0x14, 0xc7, 0x51, 0x8f, 0x7b, 0xe6, 0x44, 0xef, 0x95, 0x39, 0xc0, 0xcf,
	0x9d, 0xd9, 0x63, 0xa8, 0xf5, 0x82, 0x44, 0x9e, 0x25, 0xdc, 0xa9, 0x2c, 0x8d, 0x33, 0xae, 0xec,
	0x43, 0xb0, 0x53, 0xf1, 0x60, 0xea, 0x54, 0x89, 0x59, 0x6e, 0x70, 0xff, 0xae, 0x42, 0xe5, 0x64,
	0xd0, 0x89, 0x91, 0xe8, 0x58, 0x84, 0x86, 0xba, 0x12, 0x59, 0x13, 0xca, 0xca, 0x50, 0x26, 0x83,
	0x92, 0xd8, 0x5d, 0xa8, 0x4a, 0x21, 0x7b, 0x9c, 0x48, 0xdb, 0xbe, 0x56, 0x18, 0x83, 0x0a, 0x66,
	0x87, 0x18, 0x6d, 0xfb, 0x24, 0x33, 0x07, 0x6a, 0x69, 0x2d, 0xd2, 0x03, 0x8d, 0x8a, 0x64, 0x2e,
	0xe3, 0x09, 0x1f, 0x9d, 0x8a, 0x44, 0x3a, 0xb7, 0xf6, 0x2c, 0x24, 0x93, 0x19, 0xd8, 0x47, 0x00,
	0xa4, 0xbc, 0x16, 0x21, 0x8f, 0x9d, 0x1a, 0x85, 0x16, 0x2c, 0xec, 0x19, 0x34, 0x86, 0xe3, 0x8b,
	0x9e, 0x48, 0xba, 0x78, 0x4f, 0xa7, 0xbe, 0x34, 0x09, 0x45, 0x77, 0x76, 0x00, 0x4d, 0xbc, 0xb7,
	0xcf, 0x27, 0x82, 0x5f, 0x11, 0x80, 0xbd, 0x14, 0xe0, 0x5a, 0x04, 0xe6, 0x60, 0x12, 0xf4, 0x54,
	0x5a, 0x40, 0x85, 0xd6, 0x7d, 0xad, 0x20, 0xef, 0xab, 0x40, 0x5e, 0x76, 0x9f, 0xc7, 0x63, 0x75,
	0xe5, 0x86, 0xfa, 0x64, 0xf9, 0x05, 0x0b, 0xe5, 0x28, 0x88, 0x12, 0x67, 0x8b, 0x2e, 0x4c, 0x32,
	0x22, 0xf5, 0xc4, 0x5b, 0x9e, 0x38, 0xdb, 0xe4, 0xae, 0x15, 0xb6, 0x0b, 0xb7, 0x44, 0x72, 0xaa,
	0x44, 0xa7, 0x49, 0x07, 0xa4, 0x1a, 0x6b, 0x41, 0x1d, 0x1d, 0x28, 0x6d, 0xb7, 0x09, 0x25, 0xd3,
	0x31, 0x66, 0x3c, 0x40, 0xcd, 0xd9, 0x21, 0xa8, 0x54, 0xc3, 0x18, 0x91, 0x9c, 0xe9, 0x2f, 0xef,
	0x13, 0x5a, 0xa6, 0x23, 0x63, 0xed, 0x45, 0x88, 0x8c, 0x10, 0x0b, 0x16, 0xac, 0x53, 0x27, 0x98,
	0xc4, 0x23, 0x21, 0x15, 0xc3, 0x3b, 0x04, 0x9b, 0x1b, 0x30, 0x5a, 0x24, 0x47, 0xa9, 0xea, 0xdc,
	0x25, 0xec, 0x82, 0x85, 0xb9, 0xb0, 0x65, 0x9c, 0x09, 0xff, 0x1e, 0xe1, 0xcf, 0xd8, 0x4c, 0x35,
	0xbe, 0x8d, 0x43, 0xd1, 0x99, 0x52, 0x35, 0x76, 0x57, 0xab, 0x46, 0x1e, 0x41, 0xdd, 0x14, 0x0c,
	0x74, 0x79, 0x9c, 0x0f, 0x88, 0x46, 0x6e, 0x70, 0x0f, 0xc1, 0x3e, 0xe6, 0x52, 0x8d, 0xa5, 0xcf,
	0x7f, 0xc5, 0x12, 0x0c, 0x83, 0x88, 0x53, 0x7f, 0xab, 0x36, 0x45, 0x19, 0x6d, 0x89, 0xf8, 0x9d,
	0x53, 0x8b, 0x2b, 0x1b, 0xca, 0x66, 0x0c, 0xac, 0x6c, 0x0c, 0xdc, 0x23, 0x80, 0x33, 0x11, 0xbe,
	0x54, 0x01, 0xeb, 0xe1, 0x7c, 0x09, 0x0d, 0x4d, 0x27, 0xf1, 0x79, 0x32, 0x64, 0x5e, 0xda, 0x13,
	0x25, 0x95, 0x1b, 0xbc, 0xf5, 0x8d, 0xab, 0xcc, 0x43, 0xf6, 0xe4, 0xe7, 0xfe, 0x69, 0x41, 0xf3,
	0xa5, 0xee, 0x66, 0x9c, 0x57, 0xe4, 0xf2, 0xee, 0xc8, 0x66, 0x23, 0x5a, 0x2e, 0x8e, 0x68, 0x61,
	0x1c, 0xad, 0xd9, 0x71, 0x34, 0x8d, 0x59, 0x29, 0x34, 0xe6, 0x29, 0xd4, 0x44, 0x3f, 0xa2, 0xba,
	0x55, 0x89, 0xdb, 0xfe, 0x1c, 0x6e, 0xb3, 0x6c, 0xbc, 0x13, 0x1d, 0x74, 0x38, 0x90, 0xa3, 0xa9,
	0x6f, 0x20, 0x98, 0x0f, 0xf6, 0x04, 0x67, 0x37, 0x1b, 0xf8, 0xc6, 0xfe, 0xe3, 0xd5, 0xf0, 0x5e,
	0x9b, 0x30, 0x8d, 0x98, 0xc3, 0x60, 0x63, 0xf3, 0xdf, 0x24, 0x1f, 0x0d, 0x82, 0x1e, 0x2d, 0x09,
	0xd5, 0xd8, 0x46, 0x9f, 0x6d, 0x89, 0xfa, 0xb5, 0x96, 0x68, 0x3d, 0x85, 0xad, 0x22, 0x4d, 0xcc,
	0xe0, 0x5b, 0x3e, 0x35, 0x19, 0x54, 0x62, 0x3a, 0xe0, 0xe3, 0x2c, 0x83, 0xa4, 0x3c, 0x2d, 0x3f,
	0x29, 0xb5, 0x9e, 0x41, 0x73, 0x96, 0xd2, 0x7f, 0x89, 0x76, 0xff, 0xb2, 0xa0, 0x71, 0x18, 0x0a,
	0x69, 0x6a, 0xa7, 0x97, 0x6b, 0x29, 0x5b, 0xae, 0x69, 0x2d, 0xcb, 0x37, 0xd4, 0xd2, 0x9a, 0x53,
	0xcb, 0xca, 0xcd, 0xb5, 0xac, 0x16, 0x6a, 0x79, 0x92, 0xd7, 0x52, 0xe7, 0xbe, 0x3d, 0x27, 0xf7,
	0x05, 0x6a, 0x73, 0x0a, 0xf9, 0x7d, 0xb1, 0x90, 0x35, 0x02, 0x7b, 0xb8, 0x02, 0xd8, 0x6a, 0x55,
	0xac, 0x2f, 0xaa, 0xa2, 0xfd, 0xff, 0xa9, 0xe2, 0x31, 0xcd, 0x30, 0xde, 0x2d, 0x59, 0x6f, 0x19,
	0x7c, 0x05, 0x5b, 0x39, 0x90, 0xda, 0x06, 0x0f, 0xa1, 0x2a, 0x50, 0x49, 0xd7, 0xc1, 0xfd, 0x39,
	0x99, 0xa5, 0xac, 0x6a, 0x4f, 0xf7, 0x3e, 0xd8, 0xa8, 0x9e, 0x84, 0x37, 0xb4, 0x93, 0xda, 0x7d,
	0xdb, 0xfa, 0xa3, 0xd9, 0x5b, 0xd7, 0xfb, 0x6d, 0x35, 0x9a, 0x0f, 0xa0, 0xa1, 0x9e, 0x09, 0xa3,
	0x93, 0x50, 0xb3, 0x54, 0x41, 0xca, 0xaa, 0x49, 0xaa, 0x16, 0x43, 0xd9, 0xfd, 0x18, 0x1a, 0xfa,
	0x24, 0xed, 0xa2, 0x30, 0x72, 0x0f, 0x14, 0xf7, 0xff, 0xb1, 0xa1, 0xf6, 0x3c, 0xed, 0xd1, 0x1f,
	0xa0, 0x96, 0xee, 0x40, 0xb6, 0x37, 0xe7, 0x8a, 0xd9, 0xca, 0x6e, 0xb9, 0x0b, 0x3d, 0xe8, 0x38,
	0xf7, 0x3d, 0xf6, 0x1d, 0x34, 0x0a, 0x8b, 0x83, 0x7d, 0xb6, 0xd2, 0x72, 0x69, 0xed, 0xbe, 0xf3,
	0x2b, 0x73, 0x88, 0x6f, 0x41, 0x85, 0xf7, 0x02, 0xea, 0xa6, 0x7f, 0x99, 0xbb, 0xbc, 0xc1, 0x17,
	0x20, 0x9d, 0x41, 0xdd, 0xd4, 0x98, 0x2d, 0xb8, 0x8b, 0xe9, 0xa6, 0xd6, 0x27, 0x4b, 0x7d, 0xe8,
	0xc2, 0x2f, 0x00, 0xbe, 0xe6, 0x3d, 0x2e, 0x39, 0x51, 0xdc, 0x5b, 0xd0, 0x29, 0xd4, 0x1a, 0x0b,
	0x08, 0x1e, 0x83, 0xfd, 0x06, 0x1f, 0x29, 0x6b, 0x03, 0x1d, 0x41, 0x1d, 0x5f, 0x29, 0x6b, 0xe3,
	0xa8, 0xab, 0x9d, 0x0d, 0x36, 0x82, 0xf4, 0x86, 0xe6, 0x0b, 0x7b, 0xf7, 0x94, 0xde, 0x55, 0x0f,
	0xe6, 0x60, 0xe5, 0xbf, 0xec, 0x73, 0xdb, 0xad, 0xd0, 0xdd, 0x0a, 0xf8, 0x27, 0x68, 0xa6, 0xc0,
	0xfa, 0x25, 0xb5, 0x41, 0xe8, 0xf3, 0x6c, 0x27, 0x20, 0xe7, 0x90, 0x7d, 0xba, 0x30, 0x6a, 0x19,
	0x76, 0x61, 0x6e, 0x15, 0xf6, 0xcf, 0x44, 0x1b, 0x23, 0x35, 0xed, 0xcd, 0xa2, 0x7f, 0x03, 0x5b,
	0xe6, 0xed, 0xb7, 0x76, 0xe5, 0xce, 0xe1, 0x76, 0x9a, 0xe0, 0xec, 0x39, 0xb9, 0xb1, 0x0c, 0xff,
	0x02, 0x3b, 0x69, 0x16, 0x0c, 0xf6, 0x46, 0xf3, 0x70, 0x50, 0x3b, 0xaf, 0xea, 0xfb, 0xdc, 0xa2,
	0x7f, 0x8f, 0xfe, 0x0d, 0x00, 0x00, 0xff, 0xff, 0x4d, 0x53, 0xb2, 0xc7, 0xb6, 0x0e, 0x00, 0x00,
}
