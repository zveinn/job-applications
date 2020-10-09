package main

import (
	"log"
	"sort"
	"sync"
	"time"
)

type Transaction struct {
	ID              string
	Amount          float64
	BankName        string
	BankCountryCode string
	USDPerSecond    float64
	// MS              float64
}

type Result struct {
	ID         string
	Fraudulent bool
}

var mx = make(map[string]float64)
var mxm = sync.Mutex{}

func main() {

	transx := []Transaction{}
	for i := 0; i < 100000; i++ {
		transx = append(transx, Transaction{
			ID:              "inx",
			Amount:          float64(i + 200),
			BankName:        "X",
			BankCountryCode: "in",
		})
	}
	transx = append(transx, Transaction{
		ID:              "us-300",
		Amount:          float64(300),
		BankName:        "X",
		BankCountryCode: "us",
	})
	transx = append(transx, Transaction{
		ID:              "us-3mill",
		Amount:          float64(3000000),
		BankName:        "X",
		BankCountryCode: "us",
	})
	transx = append(transx, Transaction{
		ID:              "us-2",
		Amount:          float64(2),
		BankName:        "X",
		BankCountryCode: "us",
	})

	mx["us"] = 5
	mx["in"] = 300
	mx["uk"] = 50

	ss := time.Now()
	prioritize(transx)
	log.Println(time.Since(ss).Seconds())
	// for i, v := range transx {
	// 	log.Println(i, v)
	// }
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
		t[i].USDPerSecond = t[i].Amount / av
	}
	sort.Slice(t, func(a int, b int) bool {
		return t[a].USDPerSecond > t[b].USDPerSecond
	})
}

func ProcessTransactions(transaction []Transaction) (results []Result) {

	return
}

func isTransactionFraudulent(transaction *Transaction) {

}
