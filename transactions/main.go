package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
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

	transx := []Transaction{}
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
		transx = append(transx, Transaction{
			ID:              line[0],
			Amount:          amountAsFloat,
			BankCountryCode: line[2],
		})
	}

	RunTest(transx)
}

func RunTest(transx []Transaction) {
	ss := time.Now()
	prioritize(transx)
	finalProcessingTime := time.Since(ss)

	var maxValue float64 = 1000
	var currentMS float64 = 0
	tlength := len(transx) - 1
	var totalAmount float64 = 0

	for i := 0; i <= tlength; i++ {
		if (currentMS + transx[i].MS) > maxValue {
			continue
		}
		currentMS += transx[i].MS
		totalAmount += transx[i].Amount
	}
	postAssignTime := time.Since(ss)
	log.Println("Prioritization / Microseconds:", finalProcessingTime.Microseconds())
	log.Println("Assigning Transactions / Microseconds:", postAssignTime.Microseconds())
	log.Println("Total Amount Processed:", totalAmount)
	log.Println("Total Milliseconds assigned:", currentMS)
	log.Println("USD Per Millisecond:", totalAmount/currentMS)
	log.Println("Total USD Per 1000 Milliseconds:", 1000*(totalAmount/currentMS))
}

// 1000 ms bucket
// varying response times from apis..
// TO GET USD/per MS processing rate per fransaction we have to
// SORT BY .. AMOUNT/respone time

func prioritize(t []Transaction) {
	var av float64
	var ok2 bool
	tlength := len(t) - 1
	for i := 0; i <= tlength; i++ {
		mxm.Lock()
		av, ok2 = mx[t[i].BankCountryCode]
		mxm.Unlock()
		if !ok2 {
			// Can't find any latenzy...
			continue
		}
		t[i].MS = av
		t[i].USDPerMillisecond = t[i].Amount / av
	}
	sort.Slice(t, func(a int, b int) bool {
		return t[a].USDPerMillisecond > t[b].USDPerMillisecond
	})
}

func ProcessTransactions(transaction []Transaction) (results []Result) {

	return
}

func isTransactionFraudulent(transaction *Transaction) {

}
