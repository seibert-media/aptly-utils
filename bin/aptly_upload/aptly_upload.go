package main

import (
	"flag"
	"os"

	"runtime"

	"io"

	"fmt"

	"github.com/bborbe/log"
)

var logger = log.DefaultLogger

const (
	PARAMETER_FILE     = "file"
	PARAMETER_LOGLEVEL = "loglevel"
)

func main() {
	defer logger.Close()
	logLevelPtr := flag.String(PARAMETER_LOGLEVEL, log.INFO_STRING, "one of OFF,TRACE,DEBUG,INFO,WARN,ERROR")
	filePtr := flag.String(PARAMETER_FILE, "", "file")
	flag.Parse()
	logger.SetLevelThreshold(log.LogStringToLevel(*logLevelPtr))
	logger.Debugf("set log level to %s", *logLevelPtr)

	runtime.GOMAXPROCS(runtime.NumCPU())

	writer := os.Stdout
	err := do(writer, *filePtr)
	if err != nil {
		logger.Fatal(err)
		logger.Close()
		os.Exit(1)
	}
}

func do(writer io.Writer, file string) error {
	if len(file) == 0 {
		return fmt.Errorf("parameter file missing")
	}
	return nil
}
