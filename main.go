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
var startDir string

func main() {

	flag.BoolVar(&showWithFiles, "sf", false, "Show with files")
	flag.BoolVar(&help, "h", false, "Show commands")
	flag.BoolVar(&outputFile, "out", false, "Output to file ")
	flag.StringVar(&startDir, "sd", ".", "Start dir ")
	flag.Parse()
	if help {

		fmt.Println(">> Commands:")
		fmt.Println("-sf\tShow with files")
		fmt.Println("-out\tOutput to file DirTree.txt")
		fmt.Println("-sd\tSet start dir for scaning (default current directory)")
		os.Exit(1)
	}

	tStart := time.Now()
	scandir.GoScan(startDir, showWithFiles, outputFile)
	tEnd := time.Now()

	fmt.Println("Runtime: ", tEnd.Sub(tStart))

}
