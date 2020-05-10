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
	discardExpiredRate := flag.Int("discard-rate", 10, "discard expired orders every n seconds")
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

	ingestionRateLimit := time.NewTicker(time.Second / time.Duration(*ordersIngestionRate))
	discardExpiredOrder := time.NewTicker(time.Second * time.Duration(*discardExpiredRate))

	for {
		select {
		case <-ingestionRateLimit.C:
			first, rest := first(orders)
			order := *first
			orders = rest

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

		case orderID := <-pickupRequests:
			kitchen.PickupOrderByID(orderID)

			fmt.Printf("Event=OrderPickedUp, Order.ID=%s\n", orderID)
			if *verbose {
				fmt.Printf("Kitchen:\n%s\n", kitchen.ShelvesContent())
			}

		case <-discardExpiredOrder.C:
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
	cloudkitchen --orders $PWD/orders.json           (minimum required)

	cloudkitchen --orders $PWD/orders.json --verbose (with verbose logging)

	cloudkitchen --orders $PWD/orders.json --rate 2 --discard-rate 10 --min-pickup 2 --max-pickup 6 --verbose

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

func first(xs []Order) (*Order, []Order) {
	if len(xs) >= 1 {
		return &xs[:1][0], xs[1:]
	}
	return nil, []Order{}
}
