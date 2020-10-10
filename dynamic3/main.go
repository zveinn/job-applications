package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"sync"
)

var mx = make(map[string]float64)
var mxm = sync.Mutex{}

func main() {

	fileBytes, err := ioutil.ReadFile("api_latencies.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(fileBytes, &mx)
	if err != nil {
		panic(err)
	}

	items := []*Item{}
	csvFile, _ := os.Open("transactions.txt")
	reader := csv.NewReader(bufio.NewReader(csvFile))
	var v float64
	var ok2 bool
	for {
		line, error := reader.Read()
		if error == io.EOF {
			break
		} else if error != nil {
			log.Fatal(error)
		}
		if line[1] == "amount" {
			continue
		}
		amountAsFloat, err := strconv.ParseFloat(line[1], 64)
		if err != nil {
			log.Println(err, line[1])
			continue
		}
		mxm.Lock()
		v, ok2 = mx[line[2]]
		mxm.Unlock()
		if !ok2 {
			// Can't find the country code ..
			// we should probably handle this in a more gracefull way
			continue
		}
		items = append(items, &Item{
			v: int(amountAsFloat),
			w: int(v),
		})
	}

	wants = items
	log.Println("starting...")
	xx, w, v2 := m(len(items)-1, 1000)
	fmt.Println(items)
	log.Println(xx, v2)
	fmt.Println("weight:", w)
	fmt.Println("value:", v)
}

var wants []*Item

type Item struct {
	string
	w, v int
}

func m(i, w int) ([]string, int, int) {
	if i < 0 || w == 0 {
		return nil, 0, 0
	} else if wants[i].w > w {
		return m(i-1, w)
	}
	i0, w0, v0 := m(i-1, w)
	i1, w1, v1 := m(i-1, w-wants[i].w)
	v1 += wants[i].v
	if v1 > v0 {
		return append(i1, wants[i].string), w1 + wants[i].w, v1
	}
	return i0, w0, v0
}
