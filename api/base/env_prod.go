// +build !local

package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"net"
)

func captchaSrvAddrFunc() (string, error) {
	const captchaSrvDomain = "srv-captcha"
	_, addrs, err := net.LookupSRV("grpc", "tcp", captchaSrvDomain)
	if err != nil {
		log.Errorf("Lookup srv error %v", err)
		return "", err
	}
	for _, addr := range addrs {
		log.Infof("%s SRV is %v", captchaSrvDomain, addr)
	}
	return fmt.Sprintf("%s:%d", captchaSrvDomain, addrs[0].Port), nil
}
