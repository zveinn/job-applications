package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math"
	"os"
	"strconv"
	"sync"
	"time"
)

// type Item struct {
// 	name   string
// 	weight int
// 	value  int
// }

type Itemx struct {
	ID              string
	Amount          float64
	BankName        string
	BankCountryCode string
	// USDPerMillisecond float64
	Weight float64
	Value  float64
	MS     float64
	// MS              float64
}

var mx = make(map[string]float64)
var mxm = sync.Mutex{}

func main() {
	// knapsacksize := 50

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
			ID:     line[0],
			Value:  amountAsFloat,
			Weight: int(v),
			// Coef:   amountAsFloat / v,
		})
	}
	// DYNAMIC WEIGHT
	// start := time.Now()
	// x, _ := KnapsackDynamicWeight(50, items)
	// log.Println("DYNAMIC WEIGHT (Implemented in: https://github.com/zvadaadam/knapsack)")
	// log.Println(time.Since(start).Microseconds())
	// log.Println(x)

	// for i, v := range config1 {
	// 	if v == 1 {
	// 		log.Println("ITEM:", items[i].ID, "//", items[i].Value, "//", items[i].Weight)
	// 	}
	// }

	xt := time.Now()
	x2, _ := Custom(1000, items)
	fmt.Println("DYNAMIC WEIGHT: (refactored)")
	fmt.Println("Prioritization // Microseconds:", time.Since(xt).Microseconds())
	fmt.Println("TOTAL processed USD: ", math.Round(x2*100)/100)
	fmt.Println("USD per Millisecond:", (math.Round(x2*100)/100)/1000)
	// var totalValue float64 = 0
	// for _, v := range config2 {
	// 	if v != 0 {
	// 		totalValue += items[v].Value
	// 		fmt.Println("ITEM:", items[v].ID, "//", items[v].Value, "//", items[v].Weight)
	// 	}
	// }
}

type Item struct {
	ID     string
	Weight int
	Value  float64
	Coef   float64
	// FactorValue int
}

func Custom(capacity int, items []*Item) (float64, []int) {

	ilength := len(items)

	matrix := make([][]float64, capacity+1)
	for i := 0; i <= capacity; i++ {
		matrix[i] = make([]float64, ilength+1)
	}

	var left float64
	var right float64
	// var cw float64
	for c := 0; c <= capacity; c++ {
		// cw = float64(c)
		for i := 0; i <= ilength; i++ {

			// ignore 0 items
			if c == 0 || i == 0 {
				matrix[c][i] = 0
				// this if filters out items we can't use becaue
				// the weight is bigger then our capacity
			} else if items[i-1].Weight > c {
				matrix[c][i] = matrix[c][i-1]
			} else {
				// the value of the previous item in the item list
				// plus
				// left = matrix[current position - previous item weight][previous item position]
				// right = matrix[current positions][previous item position]
				left = items[i-1].Value + matrix[c-items[i-1].Weight][i-1]
				right = matrix[c][i-1]
				// assign whichever value is bigger
				matrix[c][i] = max(left, right)
			}
		}
	}

	// this is the list of items in a flat array
	// 0 == not using this item
	// 1 == using this item
	var config []int
	for c := capacity; c > 0; c-- {
		for i := ilength; i > 0; i-- {
			if i == 0 {
				continue
			}
			// if the last item index has a weight of 0
			// or the second last item in the matrix is not equal to the last item in the matrix
			// we leave a 1 in the config slice to indicate we want to use this item
			if items[i-1].Weight == 0 || matrix[c][i] != matrix[c][i-1] {
				c = c - items[i-1].Weight // we skip to the next weight if we found our match
				config = append(config, i-1)
			} else {
				config = append(config, 0)
			}
		}
	}

	// config = backtrack(items, capacity, matrix, capacity, ilength, config)

	return matrix[capacity][ilength], config
}
func max(a float64, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

// https://github.com/zvadaadam/knapsack
// My implementation is based on the one above, just refactored for a little more speed.
func KnapsackDynamicWeight(capacity int, items []*Item) (float64, []float64) {

	ilength := len(items)

	matrix := make([][]float64, capacity+1)
	for i := 0; i <= capacity; i++ {
		matrix[i] = make([]float64, ilength+1)
	}

	var left float64
	var right float64
	// var cw float64
	for c := 0; c <= capacity; c++ {
		// cw = float64(c)
		for i := 0; i <= ilength; i++ {

			// ignore 0 items
			if c == 0 || i == 0 {
				matrix[c][i] = 0
				// this if filters out items we can't use becaue
				// the weight is bigger then our capacity
			} else if items[i-1].Weight > c {
				matrix[c][i] = matrix[c][i-1]
			} else {
				// the value of the previous item in the item list
				// plus
				// left = matrix[current position - previous item weight][previous item position]
				// right = matrix[current positions][previous item position]
				left = items[i-1].Value + matrix[c-items[i-1].Weight][i-1]
				right = matrix[c][i-1]
				// assign whichever value is bigger
				matrix[c][i] = max(left, right)
			}
		}
	}

	// this is the list of items in a flat array
	// 0 == not using this item
	// 1 == using this item
	var config []float64
	config = backtrack(items, capacity, matrix, capacity, ilength, config)

	return matrix[capacity][ilength], config
}

// work through the matrix starting at the end.
func backtrack(items []*Item, capacity int, matrix [][]float64, indexWeight int, indexItem int, config []float64) []float64 {
	// if we have no items, we return
	if indexItem == 0 {
		return config
	}

	// if the last item index has a weight of 0
	// or the second last item in the matrix is not equal to the last item in the matrix
	// we leave a 1 in the config slice to indicate we want to use this item
	if items[indexItem-1].Weight == 0 || matrix[indexWeight][indexItem] != matrix[indexWeight][indexItem-1] {
		config = backtrack(items, capacity, matrix, indexWeight-int(items[indexItem-1].Weight), indexItem-1, config)
		config = append(config, 1)
	} else {
		config = backtrack(items, capacity, matrix, indexWeight, indexItem-1, config)
		config = append(config, 0)
	}

	return config
}
