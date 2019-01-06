package clients

import (
	"fmt"
	"google.golang.org/grpc"
	"net"
)

func getGRPCConn(addr string, srv bool) (*grpc.ClientConn, error) {
	finalAddr := addr
	if srv {
		_, addrs, err := net.LookupSRV("grpc", "tcp", addr)
		if err != nil {
			return nil, err
		}
		finalAddr = fmt.Sprintf("%s:%d", addr, addrs[0].Port)
	}
	return grpc.Dial(finalAddr, grpc.WithInsecure())
}
