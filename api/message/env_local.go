// +build local

package main

import (
	"fmt"
)

func messageSrvAddrFunc() (string, error) {
	const messageSrvDomain = "srv-message"
	return fmt.Sprintf("%s:%d", messageSrvDomain, 9090), nil
}
