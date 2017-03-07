package util

import (
	"encoding/json"
	"os"
)

type Explain interface {
	Name() string
}

func PrettyPrint(object interface{}) {
	b, err := json.MarshalIndent(object, "", "  ")
	if err != nil {
		println("Failed to pretty-print:", err)
		println(object)
	} else {
		os.Stdout.Write(b)
		println("")
	}
}
