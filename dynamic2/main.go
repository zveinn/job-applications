// https://forum.lowyat.net/topic/3765924/+0
package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"sync"
	"time"
)

type Item struct {
	name   string
	weight float64
	value  float64
}

type Solution struct {
	items       *[]Item // To avoid unneccessary copies, we just reference the array
	totalValue  float64
	totalWeight float64
}

// Implementation of https://en.wikipedia.org/wiki/Knapsack_problem#0.2F1_knapsack_problem
func knapsack(items []Item, knapsackSize int) Solution {
	// solution[x] is for the solutions for knapsack size x.
	// solution[x][y] is for the solution for knapsack size x, before adding item[y]
	// So the final solution will be in solution[knapsackSize][number_of_items + 1]
	solutions := make([][]Solution, knapsackSize+1)

	empty := make([]Item, 0)

	ilenL := len(items)
	for j := 0; j <= knapsackSize; j++ {
		solutions[j] = make([]Solution, len(items)+1)
		for i := 0; i <= ilenL; i++ {
			solutions[j][i].items = &empty
		}
	}

	// We build up our possible solutions by adding the items one by one
	for i := 0; i < ilenL; i++ {
		// ... for all possible knapsack sizes up to the size we care about
		for j := 1; j <= knapsackSize+1; j++ {
			if items[i].weight < float64(j) {
				solution := &solutions[j-1][i]

				altSolution := &solutions[j-1-int(items[i].weight)][i]

				if solution.totalValue < altSolution.totalValue+items[i].value {
					// extra work here, because we used an array ref for Solution.items
					// newItems := make([]Item, len(*altSolution.items)+1)
					// copy(newItems, *altSolution.items)
					// newItems[len(*altSolution.items)] = items[i]
					*altSolution.items = append(*altSolution.items, items[i])

					// newSolution :=
					solution = &Solution{
						// &newItems,
						altSolution.items,
						altSolution.totalValue + items[i].value,
						altSolution.totalWeight + items[i].weight,
					}
				}
				solutions[j-1][i+1] = *solution
			} else {
				solutions[j-1][i+1] = solutions[j-1][i]
			}
		}
	}

	// Since this is an iterative process, we would have solved the knapsack for all sizes up to the
	// given size. Show them.
	// for j := 1; j <= knapsackSize; j++ {
	// 	fmt.Println("Solution for size", j,
	// 		"is", *(solutions[j][len(items)].items),
	// 		"weight", solutions[j][len(items)].totalWeight,
	// 		"value", solutions[j][len(items)].totalValue,
	// 	)
	// }

	return solutions[knapsackSize][len(items)]
}

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

	items := []Item{}
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
		items = append(items, Item{
			value:  amountAsFloat,
			weight: v,
		})
	}

	knapsacksize := 1000

	// fmt.Println("For items", items)
	ss := time.Now()
	s := knapsack(items, knapsacksize)
	log.Println(time.Since(ss).Microseconds())
	log.Println(s)
}
