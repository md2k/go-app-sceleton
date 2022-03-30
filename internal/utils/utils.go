package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

func Dump(cls interface{}) {
	data, err := json.MarshalIndent(cls, "", "    ")
	if err != nil {
		log.Println("[ERROR] Oh no! There was an error on Dump command: ", err)
		return
	}
	fmt.Println(string(data))
}

func Sdump(cls interface{}) string {
	data, err := json.MarshalIndent(cls, "", "    ")
	if err != nil {
		log.Println("[ERROR] Oh no! There was an error on Dump command: ", err)
		return ""
	}
	return fmt.Sprintln(string(data))
}

func IsStirngIn(s string, l []string) bool {
	for _, b := range l {
		if strings.ToLower(b) == strings.ToLower(s) {
			return true
		}
	}
	return false
}
