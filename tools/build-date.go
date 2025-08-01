//go:build ignore

package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Print(time.Now().Format("January 2, 2006"))
}
