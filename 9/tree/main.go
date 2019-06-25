package main

import (
	"io"
	"os"
	"sort"
	"strconv"
)

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}


func dirTree(out io.Writer, path string, printFiles bool) error {
	lvl1 := 0
	lvl2 := 0
	err := dir(out, path, printFiles, &lvl1, &lvl2)
	return err
}

func dir(out io.Writer, path string, printFiles bool, lvl1 *int, lvl2 *int) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	info, _ := file.Stat()
	if info.IsDir() {
		names, _ := file.Readdirnames(-1)
		if !printFiles {
			temp := make([]string, len(names))
			copy(temp, names)
			names = nil
			for _, fileName := range temp {
				newName := path + "/" + fileName
				tempF, err := os.Open(newName)
				if err != nil {
					return err
				}
				info, _ := tempF.Stat()
				tempF.Close()
				if info.IsDir() {
					names = append(names, fileName)
				}
			}
		}
		sort.Strings(names)
		for j, name := range names {
			var prefix string
			for i := 0; i < *lvl1; i++ {
				prefix += "│\t"
			}
			for i := 0; i < *lvl2; i++ {
				prefix += "\t"
			}
			size := len(names)
			if j < size-1 {
				prefix += "├───"
			} else {
				prefix += "└───"
			}
			newName := path + "/" + name
			f, _ := os.Open(newName)
			info, _ := f.Stat()
			if !info.IsDir() {
				out.Write([]byte(prefix + name))
				size := int(info.Size())
				var res string
				if size == 0 {
					res = " (empty)\n"
				} else {
					res = " (" + strconv.Itoa(size) + "b)\n"
				}
				out.Write([]byte(res))
			} else {
				out.Write([]byte(prefix + name))
				out.Write([]byte("\n"))
			}
			f.Close()
			if j < size-1 {
				*lvl1++
			} else {
				*lvl2++
			}
			err := dir(out, newName, printFiles, lvl1, lvl2)
			if j < size-1 {
				*lvl1--
			} else {
				*lvl2--
			}
			if err != nil {
				return err
			}
		}
	}
	return nil
}
