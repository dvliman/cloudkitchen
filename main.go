package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"time"
)

func main() {
	flag.Usage = usage

	ordersFilePath := flag.String("orders", "", "filepath to orders.json (i.e $PWD/orders.json)")
	ordersIngestionRate := flag.Int("rate", 2, "orders ingestion rate per seconds")
	minPickupTime := flag.Int("min-pickup", 2, "minimum order pickup time in seconds")
	maxPickupTime := flag.Int("max-pickup", 6, "maximum order pickup time in seconds")
	verbose := flag.Bool("verbose", false, "verbosely print logs")

	flag.Parse()

	if *ordersFilePath == "" {
		usage()
		flag.PrintDefaults()
		os.Exit(1)
	}

	orders, err := readOrders(*ordersFilePath)
	if err != nil {
		panic(err)
	}

	kitchen := NewKitchen()
	pickupRequests := make(chan string)

	for {
		select {
		// ingestion rate limit
		case <-time.Tick(time.Second / time.Duration(*ordersIngestionRate)):
			firstTwo, remaining := takeFirstTwo(orders)
			orders = remaining

			for _, order := range firstTwo {
				kitchen.AcceptOrder(OrderReceived{
					Order:      order,
					QueuedTime: time.Now().UTC(),
				})

				p := randomPickupTime(*minPickupTime, *maxPickupTime)
				time.AfterFunc(p, func() {
					pickupRequests <- order.ID
				})

				fmt.Printf("Event=OrderReceived Order.ID=%s, Order.Name=%s, picking up in %d seconds\n",
					order.ID, order.Name, int(p.Seconds()))
				if *verbose {
					fmt.Printf("Kitchen:\n%s\n", kitchen.ShelvesContent())
				}
			}

		case orderID := <-pickupRequests:
			kitchen.PickupOrderByID(orderID)

			fmt.Printf("Event=OrderPickedUp, Order.ID=%s\n", orderID)
			if *verbose {
				fmt.Printf("Kitchen:\n%s\n", kitchen.ShelvesContent())
			}

		case <-time.Tick(time.Second):
			kitchen.RemoveExpiredOrders()

			fmt.Println("Event=OrderDiscarded")
			if *verbose {
				fmt.Printf("Kitchen:\n%s\n", kitchen.ShelvesContent())
			}
		}
	}
}

func usage() {
	banner := `
CloudKitchen - a system that emulates the fulfillment of delivery orders for a kitchen

Usage:
	cloudkitchen --orders /home/david/orders.json

	cloudkitchen --orders /home/david/orders.json --rate 2 --min-pickup 2 --max-pickup 6

Flags:
`
	_, _ = fmt.Fprintf(os.Stderr, banner)
}

func readOrders(filepath string) ([]Order, error) {
	bytes, err := ioutil.ReadFile(filepath)
	if err != nil {
		return []Order{}, err
	}

	var orders []Order
	if err := json.Unmarshal(bytes, &orders); err != nil {
		return []Order{}, err
	}

	return orders, nil
}

func randomPickupTime(min, max int) time.Duration {
	random := rand.Intn(max-min) + min
	return time.Second * time.Duration(random)
}

func takeFirstTwo(xs []Order) ([]Order, []Order) {
	if len(xs) > 2 {
		return xs[0:2], xs[2:]
	}
	return xs, []Order{}
}
