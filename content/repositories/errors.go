package repositories

import "errors"

var ErrUpdateInfo = errors.New("info update error")
var ErrTitleExisted = errors.New("info title has been existed")

var ErrUpdateTag = errors.New("tag update error")
var ErrTagExisted = errors.New("info tag has been existed")

var ErrUpdateFavorite = errors.New("favorite update error")
var ErrFavoriteExisted = errors.New("favorite has been existed")
var ErrFavoriteNotExisted = errors.New("favorite not exist")
var ErrFavoriteCount = errors.New("can't get favorite count")

var ErrUpdateThumb = errors.New("thumb update error")
var ErrThumbExisted = errors.New("thumb has been existed")
var ErrThumbNotExisted = errors.New("thumb not exist")
var ErrThumbCount = errors.New("can't get thumb count")
