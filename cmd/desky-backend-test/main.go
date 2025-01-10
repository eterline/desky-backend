package main

import (
	"fmt"
	"time"

	expire "github.com/eterline/desky-backend/pkg/proxm-ve-tool"
)

func main() {
	t := expire.ExpireIn(1)

	fmt.Println(t.IsExpired())

	time.Sleep(1 * time.Minute)

	fmt.Println(t.IsExpired())

	time.Sleep(10 * time.Second)

	fmt.Println(t.IsExpired())
}
