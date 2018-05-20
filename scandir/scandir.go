package scandir

import (
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
)

//GoScan is we
var filesCount = 0
var dirCount = 0

//GoScan gf
func GoScan(showWithFiles, outputFilePath bool) {
	// args := os.Args

	if showWithFiles {
		recursiveIteration(".", "├───", "", false, false, true)
	} else {
		recursiveIteration(".", "├───", "", false, false, false)
	}
	fmt.Println("Files count: ", filesCount)
	fmt.Println("Directories count: ", dirCount)

}

// //Содержит ли каталог подкаталог
// func contains(args []string, substring string) bool {
// 	for _, val := range args {
// 		if val == substring {
// 			return true
// 		}
// 	}
// 	return false
// }

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
func recursiveIteration(dir, prefix, tab string, isSubdir, isLastFile, withFiles bool) {
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
			fmt.Println(prefix + file.Name())
			// allStrings = append(allStrings, prefix+file.Name()+"\n")
			// if _, err := outputFile.WriteString(prefix + file.Name() + "\n"); err != nil {
			// _, err := outputFile.WriteString("DOWN")
			// if err != nil {
			// 	panic(err)
			// }
			dirCount++
			recursiveIteration(dir+string(os.PathSeparator)+file.Name(), prefix, tab, isSubdir, isLastFile, withFiles)
		} else {
			if withFiles {
				// if _, err := outputFile.WriteString(prefix + file.Name() + " (" + strconv.FormatInt(file.Size(), 10) + "b)"); err != nil {
				// _, err := outputFile.WriteString("HI")
				// if err != nil {
				// 	panic(err)
				// }
				filesCount++
				fmt.Println(prefix + file.Name() + " (" + strconv.FormatInt(file.Size(), 10) + "b)")
				// allStrings = append(allStrings, prefix+file.Name()+" ("+strconv.FormatInt(file.Size(), 10)+"b)"+"\n")
			}

		}
	}
}
