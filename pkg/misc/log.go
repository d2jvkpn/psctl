package misc

import (
	"fmt"
	"log"
	"time"
)

func SetLogTimeFmt() {
	log.SetFlags(0)
	log.SetOutput(new(logWriter))
}

type logWriter struct{}

func (writer *logWriter) Write(bts []byte) (int, error) {
	return fmt.Print(time.Now().Format(time.RFC3339) + " " + string(bts))
}
