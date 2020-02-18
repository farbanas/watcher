package main

import (
	"log"
	"path/filepath"
	"sort"
	"strings"
)


func parseFiles(files string) (watchList []string){
	tokens := strings.Split(files, " ")
	for _, token := range tokens {
		matches, err := filepath.Glob(token)
		if err != nil {
			log.Fatalln(err)
		}
		if matches != nil {
			watchList = append(watchList, matches...)
		}
	}
	deduplicate(watchList)
	return
}

func postprocessFlags() {
	// TODO: function that will check if flags that should exclude each other are set
}

func absPath(matches []string) (expanded []string) {
	for _, path := range matches {
		f, err := filepath.Abs(path)
		if err != nil {
			log.Println(err)
			expanded = append(expanded, path)
		} else {
			expanded = append(expanded, f)
		}
	}
	return

}

func deduplicate(watchList []string) (uniques []string) {
	sort.Strings(watchList)
	j := 0
	for i := 1; i < len(watchList); i++ {
		f1, _ := filepath.Abs(watchList[i])
		f2, _ := filepath.Abs(watchList[j])
		if f1 == f2 {
			continue
		}
		j++
		watchList[j] = watchList[i]
	}
	uniques = watchList[:j+1]
	return
}

func convertEvents(eventsString string) (events int) {
	elems := strings.Split(eventsString, " ")
	for _, el := range elems {
		if e, ok := eventsMap[el]; ok {
			events |= e
		}
	}
	return
}
