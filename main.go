package main

import (
	"fmt"
	"io/ioutil"
	"time"

	vegeta "github.com/tsenart/vegeta/lib"
)

func main() {
	rate := uint64(250) // per second
	duration := 4 * time.Second
	body, err := ioutil.ReadFile("reqBody.json")
	if err != nil {
		fmt.Println("err:", err)
		return
	}
	//fmt.Println("body:", string(body))
	targeter := vegeta.NewStaticTargeter(vegeta.Target{
		Method: "POST",
		URL:    "http://localhost:8090/alerting/v1/partners/123/sites/23/alerts",
		Body:   body,
	})
	attacker := vegeta.NewAttacker()

	var metrics vegeta.Metrics
	for res := range attacker.Attack(targeter, rate, duration) {
		metrics.Add(res)
	}
	metrics.Close()

	fmt.Printf("99th percentile: %s\n", metrics.Latencies.P99)
	fmt.Printf("max: %s\n", metrics.Latencies.Max)
	fmt.Println("total no of request executed:", metrics.Requests)
	fmt.Println("per of non error response:", metrics.Success)
}
