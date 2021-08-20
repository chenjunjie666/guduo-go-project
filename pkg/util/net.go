package util

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

func GenRandomIP() string {
	rand.Seed(time.Now().UnixNano())
	ip := ""
	seg := 0
	for {
		n := rand.Intn(256)
		if seg == 0 && (n <= 0 || n == 127 || n == 168 || n == 172 || n == 192 || n == 255){
			continue
		}
		if n <= 0 || n > 255 {
			continue
		}

		ip += fmt.Sprintf("%d.", n)
		seg++
		if seg >= 4 {
			break
		}
	}

	return strings.Trim(ip, ".")

}
