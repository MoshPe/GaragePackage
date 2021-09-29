package Garage

import (
	"log"
	"os"
)

func closeFile(importFile *os.File) {
	err := importFile.Close()
	if err != nil {
		log.Fatalln("Error in closing the import file")
	}
}
