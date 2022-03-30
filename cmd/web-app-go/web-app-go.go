package main

import (
	"fmt"
	"log"
	"strconv"
	"time"
	"web-app-go/internal/cmd"
)

var version string
var buildTime string

func main() {
	var i int64
	i, err := strconv.ParseInt(buildTime, 10, 64)
	if err != nil {
		log.Fatalf("[FATAL   ] Unable to convert buildTime string into int64: %s", err.Error())
	}
	t := time.Unix(i, 0)
	cmd.Execute(fmt.Sprintf("%s\nBuild time: %s", version, t.Format("02/01/2006, 15:04:05")))
}
