// +build !local

package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"net"
)

func messageSrvAddrFunc() (string, error) {
	const messageSrvDomain = "srv-message"
	_, addrs, err := net.LookupSRV("grpc", "tcp", messageSrvDomain)
	if err != nil {
		log.Errorf("Lookup srv error %v", err)
		return "", err
	}
	for _, addr := range addrs {
		log.Infof("%s SRV is %v", messageSrvDomain, addr)
	}
	return fmt.Sprintf("%s:%d", messageSrvDomain, addrs[0].Port), nil
}
