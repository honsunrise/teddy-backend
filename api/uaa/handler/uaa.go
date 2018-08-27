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

type Uaa struct{}

// Uaa.Register is called by the API as /uaa/Register with post body
func (e *Uaa) Register(ctx context.Context, req *api.Request, rsp *api.Response) error {
	log.Log("Received Uaa.Register request")

	// check method
	if req.Method != "POST" {
		return errors.BadRequest("com.micro.api.uaa", "require post")
	}

	// let's make sure we get json
	ct, ok := req.Header["Content-Type"]
	if !ok || len(ct.Values) == 0 {
		return errors.BadRequest("go.micro.api.uaa", "need content-type")
	}

	if ct.Values[0] != "application/json" {
		return errors.BadRequest("go.micro.api.uaa", "expect application/json")
	}

	// parse body
	var body map[string]interface{}
	json.Unmarshal([]byte(req.Body), &body)

	// extract the client from the context
	uaaClient, ok := client.UaaFromContext(ctx)
	if !ok {
		return errors.InternalServerError("com.teddy.api.uaa.register", "uaa client not found")
	}

	// make request
	response, err := uaaClient.Register(ctx, &proto.RegisterReq{
		Username: body["username"].(string),
		Password: body["password"].(string),
	})
	if err != nil {
		return errors.InternalServerError("com.teddy.api.uaa.register", err.Error())
	}

	b, _ := json.Marshal(response)

	rsp.StatusCode = 200
	rsp.Body = string(b)

	return nil
}
