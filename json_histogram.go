package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func getValue(m map[string]interface{}, l string) (interface{}, bool) {
	members := strings.Split(l, ".")
	for _, selector := range members[:len(members)-1] {
		next, ok := m[selector]
		if !ok {
			return nil, false
		}
		m = next.(map[string]interface{})
		if m == nil {
			return nil, false
		}
	}
	result, ok := m[members[len(members)-1]]
	return result, ok
}

func main() {
	fieldPtr := flag.String("field", "", "Selected field")
	flag.Parse()
	field := *fieldPtr
	if field == "" {
		log.Fatal("Missing field")
	}

	histogram := map[interface{}]uint{}

	for _, arg := range flag.Args() {
		files, err := filepath.Glob(arg)
		if err != nil {
			log.Println("Error in name or pattern", arg)
		}
		for _, filename := range files {
			file, err := os.Open(filename)
			if err != nil {
				log.Println("Cannot open file", arg)
				continue
			}
			decoder := json.NewDecoder(file)
			for {
				var m map[string]interface{}
				if err := decoder.Decode(&m); err == io.EOF {
					break
				} else if err != nil {
					log.Println("Error parsing file", arg)
				}
				v, ok := getValue(m, field)
				if ok {
					var oldCount uint
					if oldCount, ok = histogram[v]; !ok {
						oldCount = 0
					}
					histogram[v] = oldCount + 1
				}
			}

		}
	}

	histogramAr := make(histogramEntryArray, len(histogram))
	var i uint
	for k, v := range histogram {
		histogramAr[i].count = v
		histogramAr[i].value = k
		i++
	}
	if len(histogramAr) > 0 {
		sort.Sort(histogramAr)
		for _, value := range histogramAr[0:min(len(histogramAr)-1, 9)] {
			fmt.Println(value.value, value.count)
		}
	}

}

type histogramEntry struct {
	count uint
	value interface{}
}

type histogramEntryArray []histogramEntry

func (h histogramEntryArray) Len() int {
	return len(h)
}

func (h histogramEntryArray) Less(i, j int) bool {
	return h[i].count > h[j].count
}

func (h histogramEntryArray) Swap(i, j int) {
	var t = h[i]
	h[i] = h[j]
	h[j] = t
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
