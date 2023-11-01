package main

import (
	"go-socialapp/internal/socialserver"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	socialserver.NewApp("socialserver").Run()
}
