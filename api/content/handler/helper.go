package handler

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go"
	"io"
	"mime/multipart"
	"path/filepath"
	"strconv"
	"strings"
	"teddy-backend/common/proto/content"
)

func buildTags(tag string) ([]*content.TagAndType, error) {
	rawTags := strings.Split(tag, ",")
	tags := make([]*content.TagAndType, 0, len(rawTags))
	for _, item := range rawTags {
		typeAndTag := strings.Split(item, ":")
		if len(typeAndTag) == 2 {
			tags = append(tags, &content.TagAndType{
				Type: typeAndTag[0],
				Tag:  typeAndTag[1],
			})

		} else {
			return nil, ErrTagNotCorrect
		}
	}
	return tags, nil
}

func uploadFile(file *multipart.FileHeader, client *minio.Client, bucket string) (string, error) {
	filename := filepath.Base(file.Filename)
	ext := filepath.Ext(file.Filename)

	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, src); err != nil {
		return "", err
	}
	src.Seek(0, io.SeekStart)
	_, err = hash.Write([]byte(filename))
	if err != nil {
		return "", err
	}

	objectName := hex.EncodeToString(hash.Sum(nil)) + ext

	putLen, err := client.PutObject(bucket, objectName, src, -1,
		minio.PutObjectOptions{ContentType: file.Header["Content-Type"][0]})
	if err != nil {
		return "", err
	}

	if putLen != file.Size {
		return "", err
	}
	return objectName, err
}

func buildSort(sort string) ([]*content.Sort, error) {
	rawSorts := strings.Split(sort, ",")
	sorts := make([]*content.Sort, 0, len(rawSorts))
	for _, item := range rawSorts {
		nameAndOrder := strings.Split(item, ":")
		if len(nameAndOrder) == 1 {
			sorts = append(sorts, &content.Sort{
				Name: nameAndOrder[0],
				Asc:  false,
			})
		} else if len(nameAndOrder) == 2 {
			if strings.ToUpper(nameAndOrder[1]) == "ASC" {
				sorts = append(sorts, &content.Sort{
					Name: nameAndOrder[0],
					Asc:  true,
				})
			} else if strings.ToUpper(nameAndOrder[1]) == "DESC" {
				sorts = append(sorts, &content.Sort{
					Name: nameAndOrder[0],
					Asc:  false,
				})
			} else {
				return nil, ErrOrderNotCorrect
			}
		}
	}
	return sorts, nil
}

func extractPageSizeSort(ctx *gin.Context) (uint64, uint64, []*content.Sort, error) {
	page, err := strconv.ParseUint(ctx.DefaultQuery("page", "0"), 10, 32)
	if err != nil && err.(*strconv.NumError).Num != "" {
		return 0, 0, nil, err
	}

	size, err := strconv.ParseUint(ctx.DefaultQuery("size", "10"), 10, 32)
	if err != nil && err.(*strconv.NumError).Num != "" {
		return 0, 0, nil, err

	}

	var sorts []*content.Sort
	if ctx.Query("sorts") != "" {
		sorts, err = buildSort(ctx.Query("sorts"))
		if err != nil {
			return 0, 0, nil, err
		}
	}
	return page, size, sorts, nil
}
