## Design

Files grouped by functionality within the main package. `main.go` is the program
entrance. It parses argument flags for configuration purposes. It then reads input, initializes
states such as rate limit (using Go's built-in ticker), a channel for pick up requests.
On each tick, kitchen updates its state by accepting, processing, and expiring the orders

```
    Kitchen:
        * NewKitchen
        * AcceptOrder
        * PickupOrderByID
        * RemoveExpiredOrders
        * ShelvesContent
        * canMoveOneOrderToAnotherShelf
        * selectShelfByTemperature

    Shelf:
        * IsFull
        * PlaceOrder
        * RemoveOrderByID
        * RemoveOrderAtIndex
        * GetRandomOrderIndex
        * RemoveExpiredOrders
        * GetOrderIDs
        * decayModifier
        * computeShelfLife
```
## Usage
```
❯ ./build/cloudkitchen

CloudKitchen - a system that emulates the fulfillment of delivery orders for a kitchen

Usage:
        cloudkitchen --orders $PWD/orders.json           (minimum required)

        cloudkitchen --orders $PWD/orders.json --verbose (with verbose logging)

        cloudkitchen --orders $PWD/orders.json --rate 2 --discard-rate 10 --min-pickup 2 --max-pickup 6 --verbose

Flags:
  -discard-rate int
        discard expired orders every n seconds (default 10)
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

## Run Tests
* Assuming Go is installed
```
❯ make all
go test -v -race ./...
=== RUN   TestKitchenPlaceOrderToShelfAccordingly
2020/05/11 14:43:13 AcceptOrder: can not selectShelfByTemperature: order.ID=1 order.Temperature=unknown
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
=== RUN   TestKitchenShelvesContent
--- PASS: TestKitchenShelvesContent (0.00s)
=== RUN   TestTakeFirst
--- PASS: TestTakeFirst (0.00s)
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
=== RUN   TestRemoveExpiredOrders
--- PASS: TestRemoveExpiredOrders (0.00s)
=== RUN   TestShelfRemoveOrderAtIndex
--- PASS: TestShelfRemoveOrderAtIndex (0.00s)
=== RUN   TestShelfRemoveOrderByID
--- PASS: TestShelfRemoveOrderByID (0.00s)
=== RUN   TestShelfGetRandomOrderIndex
--- PASS: TestShelfGetRandomOrderIndex (0.00s)
=== RUN   TestShelfGetOrderIDs
--- PASS: TestShelfGetOrderIDs (0.00s)
PASS
ok  	github.com/dvliman/cloudkitchen	1.426s
go build -o build/cloudkitchen
```
## Build
```
make all
./build/cloudkitchen --orders $PWD/orders.json
./build/cloudkitchen --orders $PWD/orders.json --verbose
```
Note: The current binary in `build/cloudkitchen` is built
against `darwin/amd64`. If you use windows / linux, you'll
want to `make build` or `go build -o <output>`; I have not tested on
different platform/architecture

## Test Coverage
```
make coverage
```   

* kitchen.go: 89.5% coverage
* main.go: 5.7% coverage
* shelf.go: 100% coverage
* overall: 58.1% coverage

## Sample Run
```
❯ ./build/cloudkitchen --orders $PWD/orders.json
Event=OrderReceived Order.ID=a8cfcb76-7f24-4420-a5ba-d46dd77bdffd, Order.Name=Banana Split, picking up in 3 seconds
Event=OrderReceived Order.ID=58e9b5fe-3fde-4a27-8e98-682e58a4a65d, Order.Name=McFlury, picking up in 5 seconds
Event=OrderReceived Order.ID=2ec069e3-576f-48eb-869f-74a540ef840c, Order.Name=Acai Bowl, picking up in 5 seconds
Event=OrderReceived Order.ID=690b85f7-8c7d-4337-bd02-04e04454c826, Order.Name=Yogurt, picking up in 5 seconds
Event=OrderDiscarded
Event=OrderReceived Order.ID=972aa5b8-5d83-4d5e-8cf3-8a1a1437b18a, Order.Name=Chocolate Gelato, picking up in 3 seconds
Event=OrderReceived Order.ID=c18e1242-0856-4203-a98c-7066ead3bd6b, Order.Name=Cobb Salad, picking up in 4 seconds
Event=OrderReceived Order.ID=66a2611c-9a93-4ccd-bb85-98f423247bf9, Order.Name=Cottage Cheese, picking up in 3 seconds
Event=OrderReceived Order.ID=4cc9d503-4e0e-42a3-b200-87c785468df9, Order.Name=Coke, picking up in 2 seconds
Event=OrderDiscarded
Event=OrderReceived Order.ID=0a6537d7-9568-400e-9857-7aceb4e7347c, Order.Name=Snow Cone, picking up in 2 seconds
Event=OrderReceived Order.ID=74e0893f-8544-4298-b507-9cf5ab847d83, Order.Name=Pad See Ew, picking up in 2 seconds
Event=OrderReceived Order.ID=7a5ea4ed-e378-4354-8ab3-a09cf563f621, Order.Name=Chunky Monkey, picking up in 4 seconds
Event=OrderReceived Order.ID=60fbe431-7f7f-46c5-8ecf-6a26ae898685, Order.Name=Beef Stew, picking up in 5 seconds
Event=OrderDiscarded
Event=OrderReceived Order.ID=2e7f4153-92cb-43e2-96f4-d77e2f25422b, Order.Name=Cheese, picking up in 4 seconds
Event=OrderReceived Order.ID=d4ead98f-a158-429c-9c0d-28ab0fa830a6, Order.Name=Spinach Omelet, picking up in 3 seconds
Event=OrderPickedUp, Order.ID=58e9b5fe-3fde-4a27-8e98-682e58a4a65d
Event=OrderReceived Order.ID=e868e485-c759-411b-b439-ca7086326bf6, Order.Name=Beef Hash, picking up in 2 seconds
Event=OrderReceived Order.ID=2bfc7ca9-d8ea-4625-9a03-6aa4b5fb635e, Order.Name=Pork Chop, picking up in 4 seconds
Event=OrderDiscarded
Event=OrderPickedUp, Order.ID=4cc9d503-4e0e-42a3-b200-87c785468df9
Event=OrderReceived Order.ID=1a4d0e6d-e50e-4ec0-ad5d-0f33e6a724e6, Order.Name=Kale Salad, picking up in 5 seconds
Event=OrderReceived Order.ID=a9f8f42f-7a58-4629-8b66-1acbf362f3b8, Order.Name=Fresh Fruit, picking up in 3 seconds
Event=OrderPickedUp, Order.ID=c18e1242-0856-4203-a98c-7066ead3bd6b
Event=OrderPickedUp, Order.ID=74e0893f-8544-4298-b507-9cf5ab847d83
Event=OrderPickedUp, Order.ID=74e0893f-8544-4298-b507-9cf5ab847d83
Event=OrderReceived Order.ID=6d2c9eb3-6ff8-408a-a60d-e60eaa467dae, Order.Name=Cranberry Salad, picking up in 3 seconds
Event=OrderReceived Order.ID=916ddba3-9c25-4065-8a90-ecf5c40c8394, Order.Name=Fudge Ice Cream Cake, picking up in 4 seconds
Event=OrderDiscarded
Event=OrderPickedUp, Order.ID=4cc9d503-4e0e-42a3-b200-87c785468df9
Event=OrderReceived Order.ID=23f5cb3a-fa51-49ad-9e29-83e98d6e7d6f, Order.Name=Mint Chocolate Ice Cream, picking up in 5 seconds
Event=OrderReceived Order.ID=c700ab30-d7bd-4e45-aafb-4329876c716c, Order.Name=Vegan Pizza, picking up in 4 seconds
Event=OrderPickedUp, Order.ID=c18e1242-0856-4203-a98c-7066ead3bd6b
Event=OrderPickedUp, Order.ID=58e9b5fe-3fde-4a27-8e98-682e58a4a65d
Event=OrderReceived Order.ID=cf0932a9-533c-4603-bb1e-512c6e697b92, Order.Name=Orange Chicken, picking up in 2 seconds
Event=OrderReceived Order.ID=ed05dafe-d928-4f1a-856d-1dab3868acbd, Order.Name=MeatLoaf, picking up in 4 seconds
Event=OrderDiscarded
Event=OrderPickedUp, Order.ID=690b85f7-8c7d-4337-bd02-04e04454c826
Event=OrderPickedUp, Order.ID=690b85f7-8c7d-4337-bd02-04e04454c826
Event=OrderPickedUp, Order.ID=2bfc7ca9-d8ea-4625-9a03-6aa4b5fb635e
Event=OrderReceived Order.ID=53f012e3-da55-4314-b2e6-28e63f0418ad, Order.Name=Milk, picking up in 5 seconds
Event=OrderReceived Order.ID=5dd79771-4536-4eaa-b64c-63ac606f1357, Order.Name=Pastrami Sandwich, picking up in 5 seconds
Event=OrderPickedUp, Order.ID=d4ead98f-a158-429c-9c0d-28ab0fa830a6
Event=OrderReceived Order.ID=f84d8263-7aae-4254-85b7-27865ffd76b7, Order.Name=Arugula, picking up in 5 seconds
Event=OrderReceived Order.ID=f7074980-3424-43d4-9166-f495077a273e, Order.Name=Pickles, picking up in 2 seconds
Event=OrderDiscarded
```

## Sample Run (verbose logging)
```
❯ ./build/cloudkitchen --orders $PWD/orders.json --verbose
Event=OrderReceived Order.ID=a8cfcb76-7f24-4420-a5ba-d46dd77bdffd, Order.Name=Banana Split, picking up in 3 seconds
Kitchen:
        Shelf: name=Hot Shelf capacity=10 allowableTemperatures=[hot] ordersCount=0 orders=[]
        Shelf: name=Cold Shelf capacity=10 allowableTemperatures=[cold] ordersCount=0 orders=[]
        Shelf: name=Frozen Shelf capacity=10 allowableTemperatures=[frozen] ordersCount=1 orders=[a8cfcb76-7f24-4420-a5ba-d46dd77bdffd]
        Shelf: name=Overflow Shelf capacity=15 allowableTemperatures=[] ordersCount=0 orders=[]

Event=OrderReceived Order.ID=58e9b5fe-3fde-4a27-8e98-682e58a4a65d, Order.Name=McFlury, picking up in 5 seconds
Kitchen:
        Shelf: name=Hot Shelf capacity=10 allowableTemperatures=[hot] ordersCount=0 orders=[]
        Shelf: name=Cold Shelf capacity=10 allowableTemperatures=[cold] ordersCount=0 orders=[]
        Shelf: name=Frozen Shelf capacity=10 allowableTemperatures=[frozen] ordersCount=2 orders=[a8cfcb76-7f24-4420-a5ba-d46dd77bdffd 58e9b5fe-3fde-4a27-8e98-682e58a4a65d]
        Shelf: name=Overflow Shelf capacity=15 allowableTemperatures=[] ordersCount=0 orders=[]

Event=OrderReceived Order.ID=2ec069e3-576f-48eb-869f-74a540ef840c, Order.Name=Acai Bowl, picking up in 5 seconds
Kitchen:
        Shelf: name=Cold Shelf capacity=10 allowableTemperatures=[cold] ordersCount=1 orders=[2ec069e3-576f-48eb-869f-74a540ef840c]
        Shelf: name=Frozen Shelf capacity=10 allowableTemperatures=[frozen] ordersCount=2 orders=[a8cfcb76-7f24-4420-a5ba-d46dd77bdffd 58e9b5fe-3fde-4a27-8e98-682e58a4a65d]
        Shelf: name=Overflow Shelf capacity=15 allowableTemperatures=[] ordersCount=0 orders=[]
        Shelf: name=Hot Shelf capacity=10 allowableTemperatures=[hot] ordersCount=0 orders=[]

Event=OrderReceived Order.ID=690b85f7-8c7d-4337-bd02-04e04454c826, Order.Name=Yogurt, picking up in 5 seconds
Kitchen:
        Shelf: name=Hot Shelf capacity=10 allowableTemperatures=[hot] ordersCount=0 orders=[]
        Shelf: name=Cold Shelf capacity=10 allowableTemperatures=[cold] ordersCount=2 orders=[2ec069e3-576f-48eb-869f-74a540ef840c 690b85f7-8c7d-4337-bd02-04e04454c826]
        Shelf: name=Frozen Shelf capacity=10 allowableTemperatures=[frozen] ordersCount=2 orders=[a8cfcb76-7f24-4420-a5ba-d46dd77bdffd 58e9b5fe-3fde-4a27-8e98-682e58a4a65d]
        Shelf: name=Overflow Shelf capacity=15 allowableTemperatures=[] ordersCount=0 orders=[]

Event=OrderDiscarded
Kitchen:
        Shelf: name=Frozen Shelf capacity=10 allowableTemperatures=[frozen] ordersCount=0 orders=[]
        Shelf: name=Overflow Shelf capacity=15 allowableTemperatures=[] ordersCount=0 orders=[]
        Shelf: name=Hot Shelf capacity=10 allowableTemperatures=[hot] ordersCount=0 orders=[]
        Shelf: name=Cold Shelf capacity=10 allowableTemperatures=[cold] ordersCount=0 orders=[]

Event=OrderReceived Order.ID=972aa5b8-5d83-4d5e-8cf3-8a1a1437b18a, Order.Name=Chocolate Gelato, picking up in 3 seconds
Kitchen:
        Shelf: name=Hot Shelf capacity=10 allowableTemperatures=[hot] ordersCount=0 orders=[]
        Shelf: name=Cold Shelf capacity=10 allowableTemperatures=[cold] ordersCount=0 orders=[]
        Shelf: name=Frozen Shelf capacity=10 allowableTemperatures=[frozen] ordersCount=1 orders=[972aa5b8-5d83-4d5e-8cf3-8a1a1437b18a]
        Shelf: name=Overflow Shelf capacity=15 allowableTemperatures=[] ordersCount=0 orders=[]

Event=OrderReceived Order.ID=c18e1242-0856-4203-a98c-7066ead3bd6b, Order.Name=Cobb Salad, picking up in 4 seconds
Kitchen:
        Shelf: name=Cold Shelf capacity=10 allowableTemperatures=[cold] ordersCount=1 orders=[c18e1242-0856-4203-a98c-7066ead3bd6b]
        Shelf: name=Frozen Shelf capacity=10 allowableTemperatures=[frozen] ordersCount=1 orders=[972aa5b8-5d83-4d5e-8cf3-8a1a1437b18a]
        Shelf: name=Overflow Shelf capacity=15 allowableTemperatures=[] ordersCount=0 orders=[]
        Shelf: name=Hot Shelf capacity=10 allowableTemperatures=[hot] ordersCount=0 orders=[]

Event=OrderReceived Order.ID=66a2611c-9a93-4ccd-bb85-98f423247bf9, Order.Name=Cottage Cheese, picking up in 3 seconds
Kitchen:
        Shelf: name=Hot Shelf capacity=10 allowableTemperatures=[hot] ordersCount=0 orders=[]
        Shelf: name=Cold Shelf capacity=10 allowableTemperatures=[cold] ordersCount=2 orders=[c18e1242-0856-4203-a98c-7066ead3bd6b 66a2611c-9a93-4ccd-bb85-98f423247bf9]
        Shelf: name=Frozen Shelf capacity=10 allowableTemperatures=[frozen] ordersCount=1 orders=[972aa5b8-5d83-4d5e-8cf3-8a1a1437b18a]
        Shelf: name=Overflow Shelf capacity=15 allowableTemperatures=[] ordersCount=0 orders=[]

Event=OrderReceived Order.ID=4cc9d503-4e0e-42a3-b200-87c785468df9, Order.Name=Coke, picking up in 2 seconds
Kitchen:
        Shelf: name=Hot Shelf capacity=10 allowableTemperatures=[hot] ordersCount=0 orders=[]
        Shelf: name=Cold Shelf capacity=10 allowableTemperatures=[cold] ordersCount=3 orders=[c18e1242-0856-4203-a98c-7066ead3bd6b 66a2611c-9a93-4ccd-bb85-98f423247bf9 4cc9d503-4e0e-42a3-b200-87c785468df9]
        Shelf: name=Frozen Shelf capacity=10 allowableTemperatures=[frozen] ordersCount=1 orders=[972aa5b8-5d83-4d5e-8cf3-8a1a1437b18a]
        Shelf: name=Overflow Shelf capacity=15 allowableTemperatures=[] ordersCount=0 orders=[]

Event=OrderDiscarded
Kitchen:
        Shelf: name=Hot Shelf capacity=10 allowableTemperatures=[hot] ordersCount=0 orders=[]
        Shelf: name=Cold Shelf capacity=10 allowableTemperatures=[cold] ordersCount=0 orders=[]
        Shelf: name=Frozen Shelf capacity=10 allowableTemperatures=[frozen] ordersCount=0 orders=[]
        Shelf: name=Overflow Shelf capacity=15 allowableTemperatures=[] ordersCount=0 orders=[]

Event=OrderReceived Order.ID=0a6537d7-9568-400e-9857-7aceb4e7347c, Order.Name=Snow Cone, picking up in 2 seconds
Kitchen:
        Shelf: name=Frozen Shelf capacity=10 allowableTemperatures=[frozen] ordersCount=1 orders=[0a6537d7-9568-400e-9857-7aceb4e7347c]
        Shelf: name=Overflow Shelf capacity=15 allowableTemperatures=[] ordersCount=0 orders=[]
        Shelf: name=Hot Shelf capacity=10 allowableTemperatures=[hot] ordersCount=0 orders=[]
        Shelf: name=Cold Shelf capacity=10 allowableTemperatures=[cold] ordersCount=0 orders=[]

Event=OrderReceived Order.ID=74e0893f-8544-4298-b507-9cf5ab847d83, Order.Name=Pad See Ew, picking up in 2 seconds
Kitchen:
        Shelf: name=Hot Shelf capacity=10 allowableTemperatures=[hot] ordersCount=1 orders=[74e0893f-8544-4298-b507-9cf5ab847d83]
        Shelf: name=Cold Shelf capacity=10 allowableTemperatures=[cold] ordersCount=0 orders=[]
        Shelf: name=Frozen Shelf capacity=10 allowableTemperatures=[frozen] ordersCount=1 orders=[0a6537d7-9568-400e-9857-7aceb4e7347c]
        Shelf: name=Overflow Shelf capacity=15 allowableTemperatures=[] ordersCount=0 orders=[]

Event=OrderReceived Order.ID=7a5ea4ed-e378-4354-8ab3-a09cf563f621, Order.Name=Chunky Monkey, picking up in 4 seconds
Kitchen:
        Shelf: name=Hot Shelf capacity=10 allowableTemperatures=[hot] ordersCount=1 orders=[74e0893f-8544-4298-b507-9cf5ab847d83]
        Shelf: name=Cold Shelf capacity=10 allowableTemperatures=[cold] ordersCount=0 orders=[]
        Shelf: name=Frozen Shelf capacity=10 allowableTemperatures=[frozen] ordersCount=2 orders=[0a6537d7-9568-400e-9857-7aceb4e7347c 7a5ea4ed-e378-4354-8ab3-a09cf563f621]
        Shelf: name=Overflow Shelf capacity=15 allowableTemperatures=[] ordersCount=0 orders=[]

Event=OrderReceived Order.ID=60fbe431-7f7f-46c5-8ecf-6a26ae898685, Order.Name=Beef Stew, picking up in 5 seconds
Kitchen:
        Shelf: name=Hot Shelf capacity=10 allowableTemperatures=[hot] ordersCount=2 orders=[74e0893f-8544-4298-b507-9cf5ab847d83 60fbe431-7f7f-46c5-8ecf-6a26ae898685]
        Shelf: name=Cold Shelf capacity=10 allowableTemperatures=[cold] ordersCount=0 orders=[]
        Shelf: name=Frozen Shelf capacity=10 allowableTemperatures=[frozen] ordersCount=2 orders=[0a6537d7-9568-400e-9857-7aceb4e7347c 7a5ea4ed-e378-4354-8ab3-a09cf563f621]
        Shelf: name=Overflow Shelf capacity=15 allowableTemperatures=[] ordersCount=0 orders=[]

Event=OrderDiscarded
Kitchen:
        Shelf: name=Overflow Shelf capacity=15 allowableTemperatures=[] ordersCount=0 orders=[]
        Shelf: name=Hot Shelf capacity=10 allowableTemperatures=[hot] ordersCount=0 orders=[]
        Shelf: name=Cold Shelf capacity=10 allowableTemperatures=[cold] ordersCount=0 orders=[]
        Shelf: name=Frozen Shelf capacity=10 allowableTemperatures=[frozen] ordersCount=0 orders=[]

^C
```
## Improvements
* Write more tests for the main loop
* Fuzzy tests the shelf life formula?
* 
