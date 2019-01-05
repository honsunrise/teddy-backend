// +build local

package main

import (
	"fmt"
)

func captchaSrvAddrFunc() (string, error) {
	const captchaSrvDomain = "srv-captcha"
	return fmt.Sprintf("%s:%d", captchaSrvDomain, 9090), nil
}

func uaaSrvAddrFunc() (string, error) {
	const uaaSrvDomain = "srv-uaa"
	return fmt.Sprintf("%s:%d", uaaSrvDomain, 9093), nil
}
