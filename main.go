package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strconv"
)

//Массив строк для отправки в result файл
var allStrings []string

func main() {
	FileOut, err := os.Create("result.txt")
	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	defer FileOut.Close()
	args := os.Args
	fmt.Print()

	if contains(args, "-f") {
		recursiveIteration(".", "├───", "", false, false, true)
	} else {
		recursiveIteration(".", "├───", "", false, false, false)
	}

	for _, value := range allStrings {
		FileOut.WriteString(value)
	}

}

func contains(args []string, substring string) bool {
	for _, val := range args {
		if val == substring {
			return true
		}
	}
	return false
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

func recursiveIteration(dir, prefix, tab string, isSubdir, last, withFiles bool) {
	files, _ := ioutil.ReadDir(dir)
	sort.Slice(files, func(i int, j int) bool { return files[i].Name() > files[j].Name() })

	lastFileIndex := getLastFileIndex(files, withFiles)
	if isSubdir {
		if last {
			tab += "\t"
		} else {
			tab += "|" + "\t"
		}
		prefix = tab + "├───"
	}

	for i, file := range files {
		if i == lastFileIndex {
			prefix = tab + "└───"
			last = true
		} else {
			last = false
		}

		if file.IsDir() {
			isSubdir = true
			fmt.Println(prefix + file.Name())
			allStrings = append(allStrings, prefix+file.Name()+"\n")

			recursiveIteration(dir+string(os.PathSeparator)+file.Name(), prefix, tab, isSubdir, last, withFiles)
		} else {
			if withFiles {
				fmt.Println(prefix + file.Name() + " (" + strconv.FormatInt(file.Size(), 10) + "b)")
			}
			allStrings = append(allStrings, prefix+file.Name()+" ("+strconv.FormatInt(file.Size(), 10)+"b)"+"\n")

		}
	}
}
