package utils

import (
	"fmt"
	"github.com/zhsyourai/teddy-backend/common/types"
	"strings"
)

func BuildMongodbURI(databases []types.Database) string {
	parts := make([]string, len(databases))
	for i, database := range databases {
		authPart := fmt.Sprintf("%s:%s@", database.Username, database.Password)
		if authPart == ":@" {
			parts[i] = fmt.Sprint(database.Address)
		} else {
			authBase := fmt.Sprintf("/%s", database.AuthDB)
			if authBase == "/" {
				authBase = ""
			}
			parts[i] = fmt.Sprint(authPart, database.Address, authBase)
		}
	}
	result := "mongodb://"
	for _, part := range parts {
		result += fmt.Sprint(part + ",")
	}
	if result == "mongodb://" {
		return ""
	}
	return strings.TrimRight(result, ",")
}
