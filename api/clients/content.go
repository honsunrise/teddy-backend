package clients

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"net"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/zhsyourai/teddy-backend/content/proto"
	"google.golang.org/grpc"
)

var contentKey = "__teddy_content_client_key__"

const contentSrvDomain = "srv-content"

// FromContext retrieves the client from the Context
func ContentFromContext(ctx *gin.Context) (proto.ContentClient, bool) {
	c, ok := ctx.Value(contentKey).(proto.ContentClient)
	return c, ok
}

// Client returns a wrapper for the UaaClient
func ContentNew() gin.HandlerFunc {
	var client proto.ContentClient = nil
	lock := sync.Mutex{}
	return func(ctx *gin.Context) {
		if client == nil {
			lock.Lock()
			defer lock.Unlock()
			_, addrs, err := net.LookupSRV("http", "tcp", contentSrvDomain)
			if err != nil {
				log.Errorf("Dial to captcha server error %v", err)
				ctx.Next()
				return
			}
			for _, addr := range addrs {
				log.Infof("%s SRV is %v", contentSrvDomain, addr)
			}
			conn, err := grpc.Dial(fmt.Sprintf("%s:%d", contentSrvDomain, addrs[0].Port), grpc.WithInsecure())
			if err != nil {
				log.Errorf("Dial to captcha server error %v", err)
				ctx.Next()
				return
			}
			client = proto.NewContentClient(conn)
		}
		ctx.Set(contentKey, client)
		ctx.Next()
	}
}
