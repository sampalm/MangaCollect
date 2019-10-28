package sort

import (
	"log"
	"regexp"
	"strconv"
)

func GetNumber(name string) (num float64) {
	re := regexp.MustCompile(`[0-9]+([.0-9]+)?`)
	match := re.FindAllString(name, 1)
	if len(match) < 1 {
		log.Fatalln("GetNumber: regexp Failed")
	}
	num, err := strconv.ParseFloat(match[0], 64)
	if err != nil {
		log.Fatalln(err)
	}
	return num
}
