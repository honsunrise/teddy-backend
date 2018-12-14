package clients

import (
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/zhsyourai/teddy-backend/common/proto"
	"google.golang.org/grpc"
	"net"
	"sync"
)

var uaaKey = "__teddy_uaa_client_key__"

const uaaSrvDomain = "srv-uaa"

// FromContext retrieves the client from the Context
func UaaFromContext(ctx *gin.Context) (proto.UAAClient, bool) {
	c, ok := ctx.Value(uaaKey).(proto.UAAClient)
	return c, ok
}

// Client returns a wrapper for the UaaClient
func UaaNew() gin.HandlerFunc {
	var client proto.UAAClient = nil
	lock := sync.Mutex{}
	return func(ctx *gin.Context) {
		if client == nil {
			lock.Lock()
			defer lock.Unlock()
			_, addrs, err := net.LookupSRV("grpc", "tcp", uaaSrvDomain)
			if err != nil {
				log.Errorf("Lookup uaa srv error %v", err)
				ctx.Next()
				return
			}
			for _, addr := range addrs {
				log.Infof("%s SRV is %v", uaaSrvDomain, addr)
			}
			conn, err := grpc.Dial(fmt.Sprintf("%s:%d", uaaSrvDomain, addrs[0].Port), grpc.WithInsecure())
			if err != nil {
				log.Errorf("Dial to uaa server error %v", err)
				ctx.Next()
				return
			}
			client = proto.NewUAAClient(conn)
		}
		ctx.Set(uaaKey, client)
		ctx.Next()
	}
}
