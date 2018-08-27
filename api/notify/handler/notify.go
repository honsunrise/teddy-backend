package handler

import (
	"context"
	"encoding/json"
	"github.com/micro/go-log"

	api "github.com/micro/go-api/proto"
	"github.com/micro/go-micro/errors"
	"github.com/zhsyourai/teddy-backend/api/notify/client"
	"github.com/zhsyourai/teddy-backend/notify/proto"
)

type Notify struct{}

// Notify.Register is called by the API as /uaa/Register with post body
func (e *Notify) GetInbox(ctx context.Context, req *api.Request, rsp *api.Response) error {
	log.Log("Received Notify.Register request")

	// check method
	if req.Method != "POST" {
		return errors.BadRequest("com.micro.api.notify", "require post")
	}

	// let's make sure we get json
	ct, ok := req.Header["Content-Type"]
	if !ok || len(ct.Values) == 0 {
		return errors.BadRequest("go.micro.api.notify", "need content-type")
	}

	if ct.Values[0] != "application/json" {
		return errors.BadRequest("go.micro.api.notify", "expect application/json")
	}

	// parse body
	var body map[string]interface{}
	json.Unmarshal([]byte(req.Body), &body)

	// extract the client from the context
	notifyClient, ok := client.NotifyFromContext(ctx)
	if !ok {
		return errors.InternalServerError("com.teddy.api.notify.getinbox", "uaa client not found")
	}

	// make request
	response, err := notifyClient.GetInBox(ctx, &proto.GetInBoxReq{
		Username: body["username"].(string),
	})
	if err != nil {
		return errors.InternalServerError("com.teddy.api.notify.getinbox", err.Error())
	}

	b, _ := json.Marshal(response)

	rsp.StatusCode = 200
	rsp.Body = string(b)

	return nil
}
