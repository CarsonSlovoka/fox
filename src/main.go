package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

func init() {
	var workDir string
	flag.StringVar(&workDir, "wDir", ".", "working directory")
	flag.Parse()
	if err := os.Chdir(workDir); err != nil {
		log.Fatal(err)
	}
	workingDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(fmt.Sprintf("Working Directory:%s", workingDir))
}

func main() {
	quitChan := make(chan error)
	go startCMD(&quitChan)
	for {
		select {
		case err := <-quitChan:
			log.Printf("Close App. %+v\n", err)
			return
		}
	}
}
