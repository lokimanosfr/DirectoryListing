package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/lokimanosfr/DirectoryListing/scandir"
)

//Массив строк для отправки в result файл
// var allStrings []string
// var outputFile os.File
var help bool
var showWithFiles bool
var outputFile bool

func main() {
	// FileOut, err := os.Create("result.txt")
	// if err != nil {
	// 	log.Fatal("Cannot create file", err)
	// }
	// defer FileOut.Close()

	// if contains(args, "-tofile") {

	// 	outputFile, err := os.OpenFile("result.txt", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	// 	outputFile.
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}

	// }

	flag.BoolVar(&showWithFiles, "f", false, "\tShow with files")
	flag.BoolVar(&help, "help", false, "\tShow commands")
	flag.BoolVar(&outputFile, "out", false, "\tOutput to file ")
	flag.Parse()
	if help {

		fmt.Println(">> Commands:")
		fmt.Println("-f\tShow with files")
		fmt.Println("-out\tOutput to file result.txt")
		os.Exit(1)
	}

	tStart := time.Now()
	scandir.GoScan(showWithFiles, outputFile)
	tEnd := time.Now()

	fmt.Println("Runtime: ", tEnd.Sub(tStart))
	// outputFile.Close()

	// for _, value := range allStrings {
	// 	FileOut.WriteString(value)
	// }

}
