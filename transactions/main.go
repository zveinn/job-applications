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
		ID:              "us-300",
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
	m.Set("uk1", float64(50))
	m.Set("uk2", float64(50))
	m.Set("uk3", float64(50))
	m.Set("uk4", float64(50))
	m.Set("uk5", float64(50))
	m.Set("uk6", float64(50))
	m.Set("uk7", float64(50))
	m.Set("uk8", float64(50))
	m.Set("uk9", float64(50))
	m.Set("uk0", float64(50))
	m.Set("uk2", float64(50))

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
	var bv interface{}
	var ok1 bool
	var ok2 bool
	sort.Slice(t, func(a int, b int) bool {
		av, ok1 = m.Get(t[a].BankCountryCode)
		bv, ok2 = m.Get(t[b].BankCountryCode)
		if !ok1 || !ok2 {
			return false
		}
		// t[a].USDPerSecond = t[a].Amount / av.(float64)
		// t[b].USDPerSecond = t[b].Amount / bv.(float64)
		return (t[a].Amount / av.(float64)) > (t[b].Amount / bv.(float64))
	})
}

func ProcessTransactions(transaction []Transaction) (results []Result) {

	return
}

func isTransactionFraudulent(transaction *Transaction) {

}
