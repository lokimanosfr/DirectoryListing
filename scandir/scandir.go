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
	skip = "System Volume Information"
)

var (
	filesCount int //Количество найденых файлов
	dirCount   int //Количество найденых каталогов

	startPrefix string //Разделитель для построения дерева
	startTabs   string //Начальный отступ

	outToFile bool //вывести в файл?
	//Файл для вывода
	outputFileName string //Имя выходного файла
)

var scaner Scan

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
	tab            string
	outputFileName string
	outputFile     *os.File

	isLastFile    bool
	pathSeparator string
}

func initScaner(p *Scan) {

	p.prop.prefix = "├───"
	p.prop.outputFileName = "DirTree.txt"
	if p.OutputToFile {
		p.prop.outputFile, _ = os.Create(scaner.prop.outputFileName)
	}

	p.Dir.isSubdir = false
	p.prop.isLastFile = false
	p.prop.tab = ""
	p.prop.pathSeparator = string(os.PathSeparator)

	p.Dir.Path = "."

}

//GoScan Точка входа для пакета
func (p Scan) GoScan() {
	initScaner(&p)
	if !p.Dir.Exist() {
		fmt.Printf("Directory %s is not exist\n", p.Dir.Path)
		return
	}

	recursiveScan(p)

	fmt.Println("Files count: ", filesCount)
	fmt.Println("Directories count: ", dirCount)

}

//Exist Возвращает true если существует
func (p dir) Exist() bool {
	if _, err := os.Stat(p.Path); os.IsNotExist(err) {
		return false
	}
	return true
}

func skipDir(p Scan) bool {
	splitedPath := strings.Split(p.Dir.Path, p.prop.pathSeparator)
	if splitedPath[len(splitedPath)-1] == skip {
		return true
	}
	return false
}

func checkFolder(p Scan) ([]os.FileInfo, bool, int) {
	if skipDir(p) {
		return nil, false, 0
	}

	files, err := ioutil.ReadDir(p.Dir.Path)
	var onlyDirs []os.FileInfo
	if err != nil {
		fmt.Println(err)
	}
	count := len(files)
	if count == 0 {
		return files, false, 0
	}
	for i, file := range files {
		if p.ShowWithFiles {
			return files, true, count - 1
		}
		if file.IsDir() && !p.ShowWithFiles {
			onlyDirs = append(onlyDirs, files[i])
		}

	}
	if p.ShowWithFiles == false && len(onlyDirs) != 0 {
		return onlyDirs, true, len(onlyDirs) - 1
	}
	return files, false, count - 1

}

//recursiveScan рекурсивный проход по каталогу, при нахождении подкаталога
//производится рекурсивный вызов функции с новым путем и аргументов isSubdir=true
//В случае если достигнут последний файл в каталоге, вызывается
// func recursiveScan(dir, prefix, tab string, isSubdir, isLastFile bool) {
func recursiveScan(p Scan) {

	files, _ := ioutil.ReadDir(p.Dir.Path)
	files, goInside, lastFileIndex := checkFolder(p)
	if !goInside {
		return
	}
	if lastFileIndex != 0 {
		sort.Slice(files, func(i int, j int) bool { return files[i].Name() > files[j].Name() }) //Сортировка содержимого каталога A-Z А-Я
	}

	// p.prop.prefix = p.prop.tab + "├───"
	if p.Dir.isSubdir {
		if p.prop.isLastFile {
			p.prop.tab += "\t"
		} else {
			p.prop.tab += "|" + "\t"
		}

	}
	for i, file := range files {

		if i == lastFileIndex {
			p.prop.prefix = p.prop.tab + "└───"
			p.prop.isLastFile = true
		} else {
			p.prop.prefix = p.prop.tab + "├───"
			p.prop.isLastFile = false
		}

		if file.IsDir() {

			p.Dir.isSubdir = true
			if p.OnDisplay {
				fmt.Println(p.prop.prefix + file.Name())
			}

			if p.OutputToFile {
				p.prop.outputFile.WriteString(p.prop.prefix + file.Name() + "\n")
			}
			dirCount++
			p.Dir.Path += p.prop.pathSeparator + file.Name()
			recursiveScan(p)
			p.Dir.isSubdir = false
			var tmp = strings.Split(p.Dir.Path, p.prop.pathSeparator)
			p.Dir.Path = strings.Join(tmp[:len(tmp)-1], p.prop.pathSeparator)

		} else {
			if p.ShowWithFiles {

				filesCount++
				if p.OnDisplay {
					fmt.Println(p.prop.prefix + file.Name() + " (" + strconv.FormatInt(file.Size(), 10) + "b)")
				}

				if p.OutputToFile {
					p.prop.outputFile.WriteString(p.prop.prefix + file.Name() + " (" + strconv.FormatInt(file.Size(), 10) + "b)" + "\n")
				}
			}

		}
	}
}
