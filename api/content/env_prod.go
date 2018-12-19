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

func contentSrvAddrFunc() (string, error) {
	const contentSrvDomain = "srv-content"
	_, addrs, err := net.LookupSRV("grpc", "tcp", contentSrvDomain)
	if err != nil {
		log.Errorf("Lookup srv error %v", err)
		return "", err
	}
	for _, addr := range addrs {
		log.Infof("%s SRV is %v", contentSrvDomain, addr)
	}
	return fmt.Sprintf("%s:%d", contentSrvDomain, addrs[0].Port), nil
}
