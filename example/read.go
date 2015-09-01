package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/bemasher/tsip"
)

const (
	LogFilename = "capture.bin"
)

func init() {
	log.SetFlags(log.Lshortfile)
}

func main() {
	logFile, err := os.Open(LogFilename)
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()

	for idx := 0; ; idx++ {
		var packet tsip.Packet
		err = packet.Read(logFile)
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println(err)
			continue
		}

		fmt.Printf("%+v\n", packet)
	}
}
