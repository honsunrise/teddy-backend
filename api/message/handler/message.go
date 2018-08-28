package handler

import (
	"context"
	"encoding/json"
	"github.com/micro/go-log"

	api "github.com/micro/go-api/proto"
	"github.com/micro/go-micro/errors"
	"github.com/zhsyourai/teddy-backend/api/message/client"
	"github.com/zhsyourai/teddy-backend/message/proto"
)

type Message struct{}

// Message.Register is called by the API as /notify/inbox with post body
func (e *Message) Inbox(ctx context.Context, req *api.Request, rsp *api.Response) error {
	log.Log("Received Message.Register request")

	// check method
	if req.Method != "POST" {
		return errors.BadRequest("com.micro.api.message", "require post")
	}

	// let's make sure we get json
	ct, ok := req.Header["Content-Type"]
	if !ok || len(ct.Values) == 0 {
		return errors.BadRequest("go.micro.api.message", "need content-type")
	}

	if ct.Values[0] != "application/json" {
		return errors.BadRequest("go.micro.api.message", "expect application/json")
	}

	// parse body
	var body map[string]interface{}
	json.Unmarshal([]byte(req.Body), &body)

	// extract the client from the context
	messageClient, ok := client.MessageFromContext(ctx)
	if !ok {
		return errors.InternalServerError("com.teddy.api.message.inbox", "uaa client not found")
	}

	// make request
	response, err := messageClient.GetInBox(ctx, &proto.GetInBoxReq{
		Username: body["username"].(string),
	})
	if err != nil {
		return errors.InternalServerError("com.teddy.api.message.inbox", err.Error())
	}

	b, _ := json.Marshal(response)

	rsp.StatusCode = 200
	rsp.Body = string(b)

	return nil
}

// Message.Register is called by the API as /notify/inbox with post body
func (e *Message) Notify(context.Context, *api.Request, *api.Response) error {
	panic("implement me")
}
