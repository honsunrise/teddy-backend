// +build local

package main

import (
	"fmt"
)

func captchaSrvAddrFunc() (string, error) {
	const captchaSrvDomain = "srv-captcha"
	return fmt.Sprintf("%s:%d", captchaSrvDomain, 9090), nil
}

func contentSrvAddrFunc() (string, error) {
	const contentSrvDomain = "srv-content"
	return fmt.Sprintf("%s:%d", contentSrvDomain, 9091), nil
}
