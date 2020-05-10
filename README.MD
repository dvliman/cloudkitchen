## Design

Files are separated by functionality within the main package. `main.go` is the program
entrace. It parses arguments/flags for configuration purpose. It reads input, initialize
states such as: rate limit (using Go's `time.Ticker`), channel for queueing up pick up requests.
On each tick, kitchen updates its state by accepting, processing, and expiring the orders

```
    Kitchen:
        * AcceptOrder(OrderReceived)
        * PickupOrderByID()
        * RemoveExpiredOrders()
        * ShelvesContent()
        * canMoveOneOrderToAnotherShelf()

    Shelf:
        * IsFull()
        * PlaceOrder(OrderReceived)
        * RemoveOrderByID()
        * RemoveOrderAtIndex()
        * GetRandomOrderIndex()
        * RemoveExpiredOrders()
        * GetOrderIDs()
        * decayModifier()
        * computeShelfLife()
```
## Usage
```
❯ ./build/cloudkitchen

CloudKitchen - a system that emulates the fulfillment of delivery orders for a kitchen

Usage:
        cloudkitchen --orders /home/david/orders.json

        cloudkitchen --orders /home/david/orders.json --rate 2 --min-pickup 2 --max-pickup 6

Flags:
  -max-pickup int
        maximum order pickup time in seconds (default 6)
  -min-pickup int
        minimum order pickup time in seconds (default 2)
  -orders string
        filepath to orders.json (i.e $PWD/orders.json)
  -rate int
        orders ingestion rate per seconds (default 2)
  -verbose
        verbosely print logs

```

## Run Tests and Build
```
❯ make all
go test -v -race ./...
=== RUN   TestKitchenPlaceOrderToShelfAccordingly
2020/05/10 05:49:31 AcceptOrder: can not selectShelfByTemperature: order.ID=1 order.Temperature=unknown
--- PASS: TestKitchenPlaceOrderToShelfAccordingly (0.00s)
=== RUN   TestKitchenPlaceOrderToOverflowIfFull
--- PASS: TestKitchenPlaceOrderToOverflowIfFull (0.00s)
=== RUN   TestKitchenRandomlyDiscardFromOverflowIfNoMoreRoom
--- PASS: TestKitchenRandomlyDiscardFromOverflowIfNoMoreRoom (0.00s)
=== RUN   TestKitchenPickupOrderByID
--- PASS: TestKitchenPickupOrderByID (0.00s)
=== RUN   TestSelectUnknownTemperature
--- PASS: TestSelectUnknownTemperature (0.00s)
=== RUN   TestKitchenRemoveExpiredOrders
--- PASS: TestKitchenRemoveExpiredOrders (0.00s)
=== RUN   TestTakeFirstTwo
--- PASS: TestTakeFirstTwo (0.00s)
=== RUN   TestNewShelf
--- PASS: TestNewShelf (0.00s)
=== RUN   TestShelfPlaceOrder
--- PASS: TestShelfPlaceOrder (0.00s)
=== RUN   TestShelfPlaceOrderWhenFull
--- PASS: TestShelfPlaceOrderWhenFull (0.00s)
=== RUN   TestShelfDecayModifier
--- PASS: TestShelfDecayModifier (0.00s)
=== RUN   TestShelfComputeShellLife
--- PASS: TestShelfComputeShellLife (0.00s)
=== RUN   TestShelfRemoveOrderAtIndex
--- PASS: TestShelfRemoveOrderAtIndex (0.00s)
=== RUN   TestShelfRemoveOrderByID
--- PASS: TestShelfRemoveOrderByID (0.00s)
=== RUN   TestShelfGetRandomOrderIndex
--- PASS: TestShelfGetRandomOrderIndex (0.00s)
PASS
ok  	github.com/dvliman/cloudkitchen	1.786s

```

## TestCoverage
```
make coverage
```   

* kitchen.go: 78.9% coverage
* main.go: 6% coverage
* shelf.go: 100% coverage
* overall: 55.6% coverage

## Sample Run
```
❯ ./build/cloudkitchen --orders $PWD/orders.json
Event=OrderReceived Order.ID=a8cfcb76-7f24-4420-a5ba-d46dd77bdffd, Order.Name=Banana Split, picking up in 3 seconds
Event=OrderReceived Order.ID=58e9b5fe-3fde-4a27-8e98-682e58a4a65d, Order.Name=McFlury, picking up in 5 seconds
Event=OrderReceived Order.ID=2ec069e3-576f-48eb-869f-74a540ef840c, Order.Name=Acai Bowl, picking up in 5 seconds
Event=OrderReceived Order.ID=690b85f7-8c7d-4337-bd02-04e04454c826, Order.Name=Yogurt, picking up in 5 seconds
Event=OrderReceived Order.ID=972aa5b8-5d83-4d5e-8cf3-8a1a1437b18a, Order.Name=Chocolate Gelato, picking up in 3 seconds
Event=OrderReceived Order.ID=c18e1242-0856-4203-a98c-7066ead3bd6b, Order.Name=Cobb Salad, picking up in 4 seconds
Event=OrderReceived Order.ID=66a2611c-9a93-4ccd-bb85-98f423247bf9, Order.Name=Cottage Cheese, picking up in 3 seconds
Event=OrderReceived Order.ID=4cc9d503-4e0e-42a3-b200-87c785468df9, Order.Name=Coke, picking up in 2 seconds
Event=OrderReceived Order.ID=0a6537d7-9568-400e-9857-7aceb4e7347c, Order.Name=Snow Cone, picking up in 2 seconds
Event=OrderReceived Order.ID=74e0893f-8544-4298-b507-9cf5ab847d83, Order.Name=Pad See Ew, picking up in 2 seconds
Event=OrderReceived Order.ID=7a5ea4ed-e378-4354-8ab3-a09cf563f621, Order.Name=Chunky Monkey, picking up in 4 seconds
Event=OrderReceived Order.ID=60fbe431-7f7f-46c5-8ecf-6a26ae898685, Order.Name=Beef Stew, picking up in 5 seconds
Event=OrderPickedUp, Order.ID=58e9b5fe-3fde-4a27-8e98-682e58a4a65d
Event=OrderReceived Order.ID=2e7f4153-92cb-43e2-96f4-d77e2f25422b, Order.Name=Cheese, picking up in 4 seconds
Event=OrderReceived Order.ID=d4ead98f-a158-429c-9c0d-28ab0fa830a6, Order.Name=Spinach Omelet, picking up in 3 seconds
Event=OrderPickedUp, Order.ID=4cc9d503-4e0e-42a3-b200-87c785468df9
Event=OrderPickedUp, Order.ID=c18e1242-0856-4203-a98c-7066ead3bd6b
Event=OrderPickedUp, Order.ID=74e0893f-8544-4298-b507-9cf5ab847d83
Event=OrderPickedUp, Order.ID=74e0893f-8544-4298-b507-9cf5ab847d83
Event=OrderPickedUp, Order.ID=4cc9d503-4e0e-42a3-b200-87c785468df9
Event=OrderPickedUp, Order.ID=58e9b5fe-3fde-4a27-8e98-682e58a4a65d
Event=OrderPickedUp, Order.ID=c18e1242-0856-4203-a98c-7066ead3bd6b
Event=OrderReceived Order.ID=e868e485-c759-411b-b439-ca7086326bf6, Order.Name=Beef Hash, picking up in 2 seconds
Event=OrderReceived Order.ID=2bfc7ca9-d8ea-4625-9a03-6aa4b5fb635e, Order.Name=Pork Chop, picking up in 4 seconds
Event=OrderPickedUp, Order.ID=690b85f7-8c7d-4337-bd02-04e04454c826
Event=OrderPickedUp, Order.ID=690b85f7-8c7d-4337-bd02-04e04454c826
Event=OrderReceived Order.ID=1a4d0e6d-e50e-4ec0-ad5d-0f33e6a724e6, Order.Name=Kale Salad, picking up in 5 seconds
Event=OrderReceived Order.ID=a9f8f42f-7a58-4629-8b66-1acbf362f3b8, Order.Name=Fresh Fruit, picking up in 3 seconds
Event=OrderPickedUp, Order.ID=d4ead98f-a158-429c-9c0d-28ab0fa830a6
Event=OrderPickedUp, Order.ID=60fbe431-7f7f-46c5-8ecf-6a26ae898685
Event=OrderReceived Order.ID=6d2c9eb3-6ff8-408a-a60d-e60eaa467dae, Order.Name=Cranberry Salad, picking up in 3 seconds
Event=OrderReceived Order.ID=916ddba3-9c25-4065-8a90-ecf5c40c8394, Order.Name=Fudge Ice Cream Cake, picking up in 4 seconds
```

## Sample Run (verbose logging)
```
❯ ./build/cloudkitchen --orders $PWD/orders.json --verbose true
Event=OrderReceived Order.ID=a8cfcb76-7f24-4420-a5ba-d46dd77bdffd, Order.Name=Banana Split, picking up in 3 seconds
Kitchen:
	Shelf: name=Hot Shelf capacity=10 allowableTemperatures=[hot] ordersCount=0 orders=[]
	Shelf: name=Cold Shelf capacity=10 allowableTemperatures=[cold] ordersCount=0 orders=[]
	Shelf: name=Frozen Shelf capacity=10 allowableTemperatures=[frozen] ordersCount=1 orders=[a8cfcb76-7f24-4420-a5ba-d46dd77bdffd]
	Shelf: name=Overflow Shelf capacity=15 allowableTemperatures=[] ordersCount=0 orders=[]

Event=OrderReceived Order.ID=58e9b5fe-3fde-4a27-8e98-682e58a4a65d, Order.Name=McFlury, picking up in 5 seconds
Kitchen:
	Shelf: name=Cold Shelf capacity=10 allowableTemperatures=[cold] ordersCount=0 orders=[]
	Shelf: name=Frozen Shelf capacity=10 allowableTemperatures=[frozen] ordersCount=2 orders=[a8cfcb76-7f24-4420-a5ba-d46dd77bdffd 58e9b5fe-3fde-4a27-8e98-682e58a4a65d]
	Shelf: name=Overflow Shelf capacity=15 allowableTemperatures=[] ordersCount=0 orders=[]
	Shelf: name=Hot Shelf capacity=10 allowableTemperatures=[hot] ordersCount=0 orders=[]

Event=OrderReceived Order.ID=2ec069e3-576f-48eb-869f-74a540ef840c, Order.Name=Acai Bowl, picking up in 5 seconds
Kitchen:
	Shelf: name=Hot Shelf capacity=10 allowableTemperatures=[hot] ordersCount=0 orders=[]
	Shelf: name=Cold Shelf capacity=10 allowableTemperatures=[cold] ordersCount=1 orders=[2ec069e3-576f-48eb-869f-74a540ef840c]
	Shelf: name=Frozen Shelf capacity=10 allowableTemperatures=[frozen] ordersCount=2 orders=[a8cfcb76-7f24-4420-a5ba-d46dd77bdffd 58e9b5fe-3fde-4a27-8e98-682e58a4a65d]
	Shelf: name=Overflow Shelf capacity=15 allowableTemperatures=[] ordersCount=0 orders=[]

Event=OrderReceived Order.ID=690b85f7-8c7d-4337-bd02-04e04454c826, Order.Name=Yogurt, picking up in 5 seconds
Kitchen:
	Shelf: name=Hot Shelf capacity=10 allowableTemperatures=[hot] ordersCount=0 orders=[]
	Shelf: name=Cold Shelf capacity=10 allowableTemperatures=[cold] ordersCount=2 orders=[2ec069e3-576f-48eb-869f-74a540ef840c 690b85f7-8c7d-4337-bd02-04e04454c826]
	Shelf: name=Frozen Shelf capacity=10 allowableTemperatures=[frozen] ordersCount=2 orders=[a8cfcb76-7f24-4420-a5ba-d46dd77bdffd 58e9b5fe-3fde-4a27-8e98-682e58a4a65d]
	Shelf: name=Overflow Shelf capacity=15 allowableTemperatures=[] ordersCount=0 orders=[]

Event=OrderReceived Order.ID=972aa5b8-5d83-4d5e-8cf3-8a1a1437b18a, Order.Name=Chocolate Gelato, picking up in 3 seconds
Kitchen:
	Shelf: name=Hot Shelf capacity=10 allowableTemperatures=[hot] ordersCount=0 orders=[]
	Shelf: name=Cold Shelf capacity=10 allowableTemperatures=[cold] ordersCount=2 orders=[2ec069e3-576f-48eb-869f-74a540ef840c 690b85f7-8c7d-4337-bd02-04e04454c826]
	Shelf: name=Frozen Shelf capacity=10 allowableTemperatures=[frozen] ordersCount=3 orders=[a8cfcb76-7f24-4420-a5ba-d46dd77bdffd 58e9b5fe-3fde-4a27-8e98-682e58a4a65d 972aa5b8-5d83-4d5e-8cf3-8a1a1437b18a]
	Shelf: name=Overflow Shelf capacity=15 allowableTemperatures=[] ordersCount=0 orders=[]

Event=OrderReceived Order.ID=c18e1242-0856-4203-a98c-7066ead3bd6b, Order.Name=Cobb Salad, picking up in 4 seconds
Kitchen:
	Shelf: name=Hot Shelf capacity=10 allowableTemperatures=[hot] ordersCount=0 orders=[]
	Shelf: name=Cold Shelf capacity=10 allowableTemperatures=[cold] ordersCount=3 orders=[2ec069e3-576f-48eb-869f-74a540ef840c 690b85f7-8c7d-4337-bd02-04e04454c826 c18e1242-0856-4203-a98c-7066ead3bd6b]
	Shelf: name=Frozen Shelf capacity=10 allowableTemperatures=[frozen] ordersCount=3 orders=[a8cfcb76-7f24-4420-a5ba-d46dd77bdffd 58e9b5fe-3fde-4a27-8e98-682e58a4a65d 972aa5b8-5d83-4d5e-8cf3-8a1a1437b18a]
	Shelf: name=Overflow Shelf capacity=15 allowableTemperatures=[] ordersCount=0 orders=[]

Event=OrderReceived Order.ID=66a2611c-9a93-4ccd-bb85-98f423247bf9, Order.Name=Cottage Cheese, picking up in 3 seconds
Kitchen:
	Shelf: name=Hot Shelf capacity=10 allowableTemperatures=[hot] ordersCount=0 orders=[]
	Shelf: name=Cold Shelf capacity=10 allowableTemperatures=[cold] ordersCount=4 orders=[2ec069e3-576f-48eb-869f-74a540ef840c 690b85f7-8c7d-4337-bd02-04e04454c826 c18e1242-0856-4203-a98c-7066ead3bd6b 66a2611c-9a93-4ccd-bb85-98f423247bf9]
	Shelf: name=Frozen Shelf capacity=10 allowableTemperatures=[frozen] ordersCount=3 orders=[a8cfcb76-7f24-4420-a5ba-d46dd77bdffd 58e9b5fe-3fde-4a27-8e98-682e58a4a65d 972aa5b8-5d83-4d5e-8cf3-8a1a1437b18a]
	Shelf: name=Overflow Shelf capacity=15 allowableTemperatures=[] ordersCount=0 orders=[]

Event=OrderReceived Order.ID=4cc9d503-4e0e-42a3-b200-87c785468df9, Order.Name=Coke, picking up in 2 seconds
Kitchen:
	Shelf: name=Cold Shelf capacity=10 allowableTemperatures=[cold] ordersCount=5 orders=[2ec069e3-576f-48eb-869f-74a540ef840c 690b85f7-8c7d-4337-bd02-04e04454c826 c18e1242-0856-4203-a98c-7066ead3bd6b 66a2611c-9a93-4ccd-bb85-98f423247bf9 4cc9d503-4e0e-42a3-b200-87c785468df9]
	Shelf: name=Frozen Shelf capacity=10 allowableTemperatures=[frozen] ordersCount=3 orders=[a8cfcb76-7f24-4420-a5ba-d46dd77bdffd 58e9b5fe-3fde-4a27-8e98-682e58a4a65d 972aa5b8-5d83-4d5e-8cf3-8a1a1437b18a]
	Shelf: name=Overflow Shelf capacity=15 allowableTemperatures=[] ordersCount=0 orders=[]
	Shelf: name=Hot Shelf capacity=10 allowableTemperatures=[hot] ordersCount=0 orders=[]

Event=OrderReceived Order.ID=0a6537d7-9568-400e-9857-7aceb4e7347c, Order.Name=Snow Cone, picking up in 2 seconds
Kitchen:
	Shelf: name=Cold Shelf capacity=10 allowableTemperatures=[cold] ordersCount=5 orders=[2ec069e3-576f-48eb-869f-74a540ef840c 690b85f7-8c7d-4337-bd02-04e04454c826 c18e1242-0856-4203-a98c-7066ead3bd6b 66a2611c-9a93-4ccd-bb85-98f423247bf9 4cc9d503-4e0e-42a3-b200-87c785468df9]
	Shelf: name=Frozen Shelf capacity=10 allowableTemperatures=[frozen] ordersCount=4 orders=[a8cfcb76-7f24-4420-a5ba-d46dd77bdffd 58e9b5fe-3fde-4a27-8e98-682e58a4a65d 972aa5b8-5d83-4d5e-8cf3-8a1a1437b18a 0a6537d7-9568-400e-9857-7aceb4e7347c]
	Shelf: name=Overflow Shelf capacity=15 allowableTemperatures=[] ordersCount=0 orders=[]
	Shelf: name=Hot Shelf capacity=10 allowableTemperatures=[hot] ordersCount=0 orders=[]

Event=OrderReceived Order.ID=74e0893f-8544-4298-b507-9cf5ab847d83, Order.Name=Pad See Ew, picking up in 2 seconds
Kitchen:
	Shelf: name=Overflow Shelf capacity=15 allowableTemperatures=[] ordersCount=0 orders=[]
	Shelf: name=Hot Shelf capacity=10 allowableTemperatures=[hot] ordersCount=1 orders=[74e0893f-8544-4298-b507-9cf5ab847d83]
	Shelf: name=Cold Shelf capacity=10 allowableTemperatures=[cold] ordersCount=5 orders=[2ec069e3-576f-48eb-869f-74a540ef840c 690b85f7-8c7d-4337-bd02-04e04454c826 c18e1242-0856-4203-a98c-7066ead3bd6b 66a2611c-9a93-4ccd-bb85-98f423247bf9 4cc9d503-4e0e-42a3-b200-87c785468df9]
	Shelf: name=Frozen Shelf capacity=10 allowableTemperatures=[frozen] ordersCount=4 orders=[a8cfcb76-7f24-4420-a5ba-d46dd77bdffd 58e9b5fe-3fde-4a27-8e98-682e58a4a65d 972aa5b8-5d83-4d5e-8cf3-8a1a1437b18a 0a6537d7-9568-400e-9857-7aceb4e7347c]

Event=OrderReceived Order.ID=7a5ea4ed-e378-4354-8ab3-a09cf563f621, Order.Name=Chunky Monkey, picking up in 4 seconds
Kitchen:
	Shelf: name=Hot Shelf capacity=10 allowableTemperatures=[hot] ordersCount=1 orders=[74e0893f-8544-4298-b507-9cf5ab847d83]
	Shelf: name=Cold Shelf capacity=10 allowableTemperatures=[cold] ordersCount=5 orders=[2ec069e3-576f-48eb-869f-74a540ef840c 690b85f7-8c7d-4337-bd02-04e04454c826 c18e1242-0856-4203-a98c-7066ead3bd6b 66a2611c-9a93-4ccd-bb85-98f423247bf9 4cc9d503-4e0e-42a3-b200-87c785468df9]
	Shelf: name=Frozen Shelf capacity=10 allowableTemperatures=[frozen] ordersCount=5 orders=[a8cfcb76-7f24-4420-a5ba-d46dd77bdffd 58e9b5fe-3fde-4a27-8e98-682e58a4a65d 972aa5b8-5d83-4d5e-8cf3-8a1a1437b18a 0a6537d7-9568-400e-9857-7aceb4e7347c 7a5ea4ed-e378-4354-8ab3-a09cf563f621]
	Shelf: name=Overflow Shelf capacity=15 allowableTemperatures=[] ordersCount=0 orders=[]

Event=OrderReceived Order.ID=60fbe431-7f7f-46c5-8ecf-6a26ae898685, Order.Name=Beef Stew, picking up in 5 seconds
Kitchen:
	Shelf: name=Frozen Shelf capacity=10 allowableTemperatures=[frozen] ordersCount=5 orders=[a8cfcb76-7f24-4420-a5ba-d46dd77bdffd 58e9b5fe-3fde-4a27-8e98-682e58a4a65d 972aa5b8-5d83-4d5e-8cf3-8a1a1437b18a 0a6537d7-9568-400e-9857-7aceb4e7347c 7a5ea4ed-e378-4354-8ab3-a09cf563f621]
	Shelf: name=Overflow Shelf capacity=15 allowableTemperatures=[] ordersCount=0 orders=[]
	Shelf: name=Hot Shelf capacity=10 allowableTemperatures=[hot] ordersCount=2 orders=[74e0893f-8544-4298-b507-9cf5ab847d83 60fbe431-7f7f-46c5-8ecf-6a26ae898685]
	Shelf: name=Cold Shelf capacity=10 allowableTemperatures=[cold] ordersCount=5 orders=[2ec069e3-576f-48eb-869f-74a540ef840c 690b85f7-8c7d-4337-bd02-04e04454c826 c18e1242-0856-4203-a98c-7066ead3bd6b 66a2611c-9a93-4ccd-bb85-98f423247bf9 4cc9d503-4e0e-42a3-b200-87c785468df9]

^C
```
## Improvements
* What unit is the decay rate? shelf life?
* What is the precision of that?
