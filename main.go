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
var onDisplay bool

func init() {
	flag.BoolVar(&showWithFiles, "f", false, "with files")
	flag.BoolVar(&help, "h", false, "Show commands")
	flag.BoolVar(&outputFile, "out", false, "Output to file ")
	flag.StringVar(&startDir, "path", ".", "Start dir ")
	flag.BoolVar(&onDisplay, "d", false, "Show on display")
	flag.Parse()
	if help {

		fmt.Println(">> Commands:")
		fmt.Println("-f\tWith files (Default only directory)")
		fmt.Println("-out\tOutput to file DirTree.txt")
		fmt.Println("-path\tSet start dir for scaning (default current directory)")
		fmt.Println("-d\tShow on display")
		os.Exit(1)
	}
}

func main() {

	var scan scandir.Scan
	scan.Dir.Path = startDir
	scan.OutputToFile = outputFile
	scan.ShowWithFiles = showWithFiles
	scan.OnDisplay = onDisplay

	tStart := time.Now()
	scan.GoScan()
	tEnd := time.Now()

	fmt.Println("Runtime: ", tEnd.Sub(tStart))

}
