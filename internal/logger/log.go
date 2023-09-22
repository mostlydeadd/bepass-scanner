package logger

import (
	"fmt"
)

var silent bool

func Silent() {
	silent = true
}

func Log(text string, prefix string) {
	if !silent {
		fmt.Printf("[%s] %s\n", prefix, text)
	}
}
