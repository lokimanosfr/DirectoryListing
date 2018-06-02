package scandir

import (
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
	"strings"
)

const (
	skipSVI = "System Volume Information"

	prefMiddle = "├───" //Разделитель между неконечными файлами
	prefEnd    = "└───" //Разделитель для конечного файла
	prefSpace  = "\t"   //Отступ между уровнями
	prefLine   = "|"    //Ведущая линия между уровнями
)

var (
	filesCount int //Количество найденых файлов
	dirCount   int //Количество найденых каталогов
)

type dir struct {
	Path     string
	isSubdir bool
	weight   int
}

//Scan dsds
type Scan struct {
	Dir           dir
	ShowWithFiles bool
	OutputToFile  bool
	OnDisplay     bool
	prop          prop
}
type prop struct {
	prefix         string
	space          string
	outputFileName string
	outputFile     *os.File

	isLastFile    bool
	pathSeparator string
}

func initScaner(scan *Scan) {

	scan.prop.prefix = prefMiddle
	scan.prop.pathSeparator = string(os.PathSeparator)
	scan.prop.outputFileName = "." + scan.prop.pathSeparator + "DirTree.txt"
	if scan.OutputToFile {
		var err error
		scan.prop.outputFile, err = os.Create(scan.prop.outputFileName)
		if err != nil {
			fmt.Println(err)
		}
	}

	scan.Dir.isSubdir = false
	scan.prop.isLastFile = false
	scan.prop.space = ""

}

//GoScan Точка входа для пакета
func (scan Scan) GoScan() {
	initScaner(&scan)
	if !scan.Dir.Exist() {
		fmt.Printf("Directory %s is not exist\n", scan.Dir.Path)
		return
	}

	recursiveScan(scan)

	fmt.Println("Files count: ", filesCount)
	fmt.Println("Directories count: ", dirCount)

}

//Exist Возвращает true если файл существует
func (scan dir) Exist() bool {
	if _, err := os.Stat(scan.Path); os.IsNotExist(err) {
		return false
	}
	return true
}

//skipDir Возвращает true если файл находится в исключениях
func skipDir(scan Scan) bool {
	splitedPath := strings.Split(scan.Dir.Path, scan.prop.pathSeparator)
	if splitedPath[len(splitedPath)-1] == skipSVI {
		return true
	}
	return false
}

//chekFolder Возвращает:
//[]os.Fileinfo - возвращается только тогда, когда запрашиваются только директории.
//bool - признак, проваливаться ли в директорию. Возвращает false только когда файл
//содержится в искулючениях или директорию пуста.
//int - индекс последнего файла, для своевременной остановки цикла в рекурсивной функции.
func checkFolder(scan Scan) ([]os.FileInfo, bool, int) {
	if skipDir(scan) {
		return nil, false, 0
	}

	files, err := ioutil.ReadDir(scan.Dir.Path)
	var onlyDirs []os.FileInfo
	if err != nil {
		fmt.Println(err)
	}
	lastIndex := len(files)
	if lastIndex == 0 {
		return files, false, 0
	}
	for i, file := range files {
		if scan.ShowWithFiles {
			return files, true, lastIndex - 1
		}
		if file.IsDir() && !scan.ShowWithFiles {
			onlyDirs = append(onlyDirs, files[i])
		}

	}
	if scan.ShowWithFiles == false && len(onlyDirs) != 0 {
		return onlyDirs, true, len(onlyDirs) - 1
	}
	return files, false, lastIndex - 1

}

//recursiveScan Рекурсивный проход по каталогу, при нахождении подкаталога
//производится вызов функции с новым путем и isSubdir=true
//В случае если достигнут последний файл в каталоге, возвращается на уровень выше
func recursiveScan(scan Scan) {
	files, goInside, lastFileIndex := checkFolder(scan)
	if !goInside {
		return
	}
	if lastFileIndex != 0 {
		sort.Slice(files, func(i int, j int) bool { return files[i].Name() > files[j].Name() }) //Сортировка содержимого каталога A-Z А-Я
	}

	if scan.Dir.isSubdir {
		if scan.prop.isLastFile {
			scan.prop.space += prefSpace
		} else {
			scan.prop.space += prefLine + prefSpace
		}

	}
	for i, file := range files {

		if i == lastFileIndex {
			scan.prop.prefix = scan.prop.space + prefEnd
			scan.prop.isLastFile = true
		} else {
			scan.prop.prefix = scan.prop.space + prefMiddle
			scan.prop.isLastFile = false
		}

		if file.IsDir() {

			scan.Dir.isSubdir = true
			if scan.OnDisplay {
				fmt.Println(scan.prop.prefix + file.Name())
			}

			if scan.OutputToFile {
				scan.prop.outputFile.WriteString(scan.prop.prefix + file.Name() + "\n")
			}
			dirCount++
			scan.Dir.Path += scan.prop.pathSeparator + file.Name()
			recursiveScan(scan)
			scan.Dir.isSubdir = false
			var tmp = strings.Split(scan.Dir.Path, scan.prop.pathSeparator)
			scan.Dir.Path = strings.Join(tmp[:len(tmp)-1], scan.prop.pathSeparator)

		} else {
			if scan.ShowWithFiles {

				filesCount++
				if scan.OnDisplay {
					fmt.Println(scan.prop.prefix + file.Name() + " (" + strconv.FormatInt(file.Size(), 10) + "b)")
				}

				if scan.OutputToFile {
					scan.prop.outputFile.WriteString(scan.prop.prefix + file.Name() + " (" + strconv.FormatInt(file.Size(), 10) + "b)" + "\n")
				}
			}

		}
	}
}
