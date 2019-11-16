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

	sorter "github.com/sampalm/MangaCollect/sort"
	"github.com/sampalm/MangaCollect/worker"
)

var dirname string
var ffolder string
var vfolder string
var lcap float64
var lvol float64
var defaultFolter = "_volumes"

func init() {
	if len(os.Args) < 2 {
		log.Fatalln("You must set a path!")
	}
	dirname = os.Args[1]
	ffolder = filepath.Base(dirname)
	vfolder = filepath.Join(dirname, defaultFolter)
}

func main() {
	fmt.Println("::: Checking existence of volumes :::")
	if _, err := os.Stat(vfolder); err == nil {

		list, err := filepath.Glob(vfolder + "/*.zip") // returns a list of files
		if err != nil {
			log.Fatal(err)
		}

		for i, l := range list {
			list[i] = filepath.Base(l)
		}

		sort.Slice(list,
			func(i, j int) bool {
				return sorter.GetNumber(list[i]) < sorter.GetNumber(list[j])
			},
		)

		cpt, err := worker.GetLastChapter(filepath.Join(vfolder, list[len(list)-1]))
		if err != nil {
			log.Fatalln(err)
		}
		lvol = sorter.GetNumber(list[len(list)-1])
		lcap = sorter.GetNumber(cpt)
		fmt.Printf("Found: Vol %.1f - Cap %.1f\n", lvol, lcap)
	}

	fmt.Println("::: Running new task :::")

	if len(os.Args) > 2 {
		if os.Args[2] == "reset" {
			files := worker.ReadPath(dirname, defaultFolter)
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

	mangafl := worker.ReadPath(dirname, defaultFolter)
	// manga volume counter
	ct := 0
	tt := 1
	if lvol > 0 {
		tt = int(lvol) + 1
	}

	sort.Slice(mangafl,
		func(i, j int) bool {
			return sorter.GetNumber(mangafl[i].Name()) < sorter.GetNumber(mangafl[j].Name())
		},
	)

	counter := make(map[string][]string, 0)
	for _, file := range mangafl {
		// This is the default folder to save files
		if file.Name() == defaultFolter {
			continue
		}
		if lcap > 0 {
			if sorter.GetNumber(file.Name()) <= lcap {
				continue
			}
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

	if _, err := os.Stat(vfolder); os.IsNotExist(err) {
		err := os.MkdirAll(vfolder, 555)
		if err != nil {
			log.Fatalln(err)
		}
	}

	for vol, files := range counter {
		fmt.Printf("::: Creating %s :::\n", vol)
		if err := worker.CompressFiles(filepath.Join(dirname, defaultFolter, ffolder+" "+vol+".zip"), files); err != nil {
			log.Println(err)
		}
	}
	fmt.Println("::: Work completed :::")
}
