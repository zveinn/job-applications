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
	"sort"
	"strconv"
	"sync"
	"time"
)

type Transaction struct {
	ID                string
	Amount            float64
	BankName          string
	BankCountryCode   string
	USDPerMillisecond float64
	MS                float64
	// MS              float64
}

type Result struct {
	ID         string
	Fraudulent bool
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

	transx := []*Transaction{}
	csvFile, _ := os.Open("transactions.txt")
	reader := csv.NewReader(bufio.NewReader(csvFile))
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
		transx = append(transx, &Transaction{
			ID:              line[0],
			Amount:          amountAsFloat,
			BankCountryCode: line[2],
		})
	}

	prioritize(transx)

}

func prioritize(transx []*Transaction) []*Result {
	ss := time.Now()
	var v float64
	var ok2 bool
	tlength := len(transx) - 1

	for i := 0; i <= tlength; i++ {
		mxm.Lock()
		v, ok2 = mx[transx[i].BankCountryCode]
		mxm.Unlock()
		if !ok2 {
			// Can't find the country code ..
			// we should probably handle this in a more gracefull way
			continue
		}
		transx[i].MS = v
		transx[i].USDPerMillisecond = transx[i].Amount / v
	}

	sort.Slice(transx, func(a int, b int) bool {
		return transx[a].USDPerMillisecond > transx[b].USDPerMillisecond
	})
	// PROPERLY FIND THE OPTIMAL FILLING

	finalProcessingTime := time.Since(ss)
	fmt.Println("Prioritization / Microseconds:", finalProcessingTime.Microseconds())
	return ProcessTransactions(transx)
}

func ProcessTransactions(transx []*Transaction) (results []*Result) {
	ss := time.Now()
	tlength := len(transx) - 1
	var maxValue float64 = 1000
	var currentMS float64 = 0
	var totalAmount float64 = 0
	for i := 0; i <= tlength; i++ {
		if (currentMS + transx[i].MS) > maxValue {
			continue
		}
		currentMS += transx[i].MS
		totalAmount += transx[i].Amount
		results = append(results, &Result{
			ID:         transx[i].ID,
			Fraudulent: isTransactionFraudulent(transx[i]),
		})
	}
	postAssignTime := time.Since(ss)
	fmt.Println("Processing Transactions / Microseconds:", postAssignTime.Microseconds())
	fmt.Println("Total Amount Processed:", totalAmount)
	fmt.Println("Total Transactions Processed:", len(results))
	fmt.Println("Total Milliseconds assigned:", currentMS)
	fmt.Println("USD Per Millisecond:", totalAmount/currentMS)
	fmt.Println("Total USD Per 1000 Milliseconds:", 1000*(totalAmount/currentMS))
	return
}

func isTransactionFraudulent(transaction *Transaction) bool {
	return true
}
