package handler

import (
	"context"
	"encoding/json"
	"github.com/micro/go-log"

	api "github.com/micro/go-api/proto"
	"github.com/micro/go-micro/errors"
	"github.com/zhsyourai/teddy-backend/api/uaa/client"
	"github.com/zhsyourai/teddy-backend/uaa/proto"
)

type UaaApi struct{}

func extractValue(pair *api.Pair) string {
	if pair == nil {
		return ""
	}
	if len(pair.Values) == 0 {
		return ""
	}
	return pair.Values[0]
}

// UaaApi.Call is called by the API as /uaa/Register with post body {"name": "foo"}
func (e *UaaApi) Register(ctx context.Context, req *api.Request, rsp *api.Response) error {
	log.Log("Received UaaApi.Register request")

	// extract the client from the context
	uaaClient, ok := client.UaaFromContext(ctx)
	if !ok {
		return errors.InternalServerError("com.teddy.api.uaa.register", "uaa client not found")
	}

	// make request
	response, err := uaaClient.Register(ctx, &proto.RegisterReq{
		Username: extractValue(req.Post["username"]),
		Password: extractValue(req.Post["password"]),
	})
	if err != nil {
		return errors.InternalServerError("com.teddy.api.uaa.register", err.Error())
	}

	b, _ := json.Marshal(response)

	rsp.StatusCode = 200
	rsp.Body = string(b)

	return nil
}
