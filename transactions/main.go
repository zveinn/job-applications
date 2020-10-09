package main

import (
	"log"
	"sort"
	"time"

	"github.com/cornelk/hashmap"
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

// Lock free concurrent hash map implementation
var m = &hashmap.HashMap{}

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

	m.Set("us", float64(5))
	m.Set("in", float64(200))
	m.Set("uk", float64(50))
	// m.Set("uk2", float64(50))
	// m.Set("uk3", float64(50))
	// m.Set("uk4", float64(50))
	// m.Set("uk5", float64(50))
	// m.Set("uk6", float64(50))
	// m.Set("uk7", float64(50))
	// m.Set("uk8", float64(50))
	// m.Set("uk9", float64(50))
	// m.Set("uk0", float64(50))
	// m.Set("uk2", float64(50))

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
	var av interface{}
	// var bv interface{}
	// var ok1 bool
	var ok2 bool
	for i, v := range t {
		av, ok2 = m.Get(v.BankCountryCode)
		if !ok2 {
			// Can't find any latenzy...
			continue
		}
		t[i].USDPerSecond = v.Amount / av.(float64)
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
