package main

import (
	"fmt"
	// time.LoadLocation needs tzdata on the host, and some minimal images might not have it.
	// This will embed tzdata in built binary, so no problems on such hosts.
	_ "time/tzdata"
)

func main() {
	fmt.Println("HTTP Runner")
}
