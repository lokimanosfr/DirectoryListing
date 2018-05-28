package scandir

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strconv"
)

//GoScan is we
var filesCount int
var dirCount int

var withFiles bool
var startPrefix = "├───"
var startTabs = ""

var outToFile bool
var outputFile *os.File
var outputFileName = "DirTree.txt"

//GoScan start point
func GoScan(startDir string, showWithFiles, outputToFile bool) {
	// args := os.Args
	if _, err := os.Stat(startDir); os.IsNotExist(err) {
		fmt.Printf("Directory %s is not exist\n", startDir)
		return
	}
	withFiles = showWithFiles
	outToFile = outputToFile
	if outToFile {
		var err error
		outputFile, err = os.Create(outputFileName)
		defer outputFile.Close()
		if err != nil {
			fmt.Println(err)
		} else {
			dir := ""
			if startDir != "." {
				dir = startDir
			} else {
				dir, _ = filepath.Abs(filepath.Dir(os.Args[0]))
			}

			outputFile.WriteString(dir + "\n")
		}
	}

	recursiveScan(startDir, startPrefix, startTabs, false, false)

	fmt.Println("Files count: ", filesCount)
	fmt.Println("Directories count: ", dirCount)

}

func getLastFileIndex(files []os.FileInfo, withFiles bool) (lastFileIndex int) {
	for index, file := range files {
		if withFiles {
			lastFileIndex = index
		} else {
			if file.IsDir() {
				lastFileIndex = index
			}

		}
	}

	return lastFileIndex
}

//Scaning recursively all folders
func recursiveScan(dir, prefix, tab string, isSubdir, isLastFile bool) {
	files, _ := ioutil.ReadDir(dir)
	sort.Slice(files, func(i int, j int) bool { return files[i].Name() > files[j].Name() })

	lastFileIndex := getLastFileIndex(files, withFiles)

	if isSubdir {
		if isLastFile {
			tab += "\t"
		} else {
			tab += "|" + "\t"
		}

	}
	prefix = tab + "├───"
	for i, file := range files {
		if i == lastFileIndex {
			prefix = tab + "└───"
			isLastFile = true
		} else {
			isLastFile = false
		}

		if file.IsDir() {
			isSubdir = true
			// fmt.Println(prefix + file.Name())
			if outToFile {
				outputFile.WriteString(prefix + file.Name() + "\n")

			}
			dirCount++
			recursiveScan(dir+string(os.PathSeparator)+file.Name(), prefix, tab, isSubdir, isLastFile)
		} else {
			if withFiles {

				filesCount++
				// fmt.Println(prefix + file.Name() + " (" + strconv.FormatInt(file.Size(), 10) + "b)")
				if outToFile {
					outputFile.WriteString(prefix + file.Name() + " (" + strconv.FormatInt(file.Size(), 10) + "b)" + "\n")
				}
			}

		}
	}
}
