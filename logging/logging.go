package logging

import "fmt"

const FILE_NAME_FORMAT = "%d_%d_%d"
const LOG_DIR = "_logs"

var debug bool

func Debug(value bool) {
	debug = value
}

func Logging(message string) {
	if !debug {
		return
	}
	fmt.Println(message)
}
