package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"unicode"

	sorter "github.com/sampalm/comicJoin/sort"
	"github.com/sampalm/comicJoin/worker"
)

var dirname string
var ffolder string

func init() {
	if len(os.Args) < 2 {
		log.Fatalln("You must set a path!")
	}
	dirname = os.Args[1]
	ffolder = filepath.Base(dirname)
}

func main() {
	fmt.Println("::: Running new task :::")

	if len(os.Args) > 2 {
		if os.Args[2] == "reset" {
			files := worker.ReadPath(dirname)
			for _, file := range files {
				reg := regexp.MustCompile(`^((Vol|vol|volume|Volume)(.|)[0-9]+)`).FindString(file.Name())
				if reg != "" {
					newp := strings.Split(file.Name(), reg)[1]
					if unicode.IsSpace(rune(newp[0])) {
						newp = strings.Replace(newp, " ", "", 1)
					}
					fmt.Printf("Renaming File: %s\n", newp)
					if err := os.Rename(filepath.Join(dirname, file.Name()), filepath.Join(dirname, newp)); err != nil {
						log.Fatal(err)
					}
				}
			}
			fmt.Println("::: Reset finished :::")
		}
	}

	files := worker.ReadPath(dirname)
	// manga volume counter
	ct := 0
	tt := 1

	sort.Slice(files,
		func(i, j int) bool {
			return sorter.GetNumber(files[i].Name()) < sorter.GetNumber(files[j].Name())
		},
	)

	counter := make(map[string][]string, 0)

	for _, file := range files {

		// This is the default folder to save files
		if file.Name() == "Volumes" {
			continue
		}
		reg := regexp.MustCompile(`^((Vol|vol|volume|Volume)(.|)[0-9]+)`).FindString(file.Name())
		if reg == "" {
			if ct == 7 {
				ct = 0
				tt++
			}
			reg = fmt.Sprintf("Vol %d", tt)
			ct++
		}
		counter[reg] = append(counter[reg], filepath.Join(dirname, file.Name()))
	}

	if len(counter) < 1 {
		fmt.Println("No chapter was found")
		return
	}
	fmt.Printf("::: Working in %s now :::\n", ffolder)

	volpath := filepath.Join(dirname, "Volumes")
	if _, err := os.Stat(volpath); os.IsNotExist(err) {
		err := os.MkdirAll(volpath, 555)
		if err != nil {
			log.Fatalln(err)
		}
	}

	for vol, files := range counter {
		fmt.Printf("::: Creating %s :::\n", vol)
		if err := worker.CompressFiles(filepath.Join(dirname, "Volumes", ffolder+" "+vol+".zip"), files); err != nil {
			log.Println(err)
		}
	}
	fmt.Println("::: Work completed :::")
}
