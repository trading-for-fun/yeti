package book

import "testing"

func TestPlacingOrders(t *testing.T) {
	book := NewInMemoryOrderBook()

	sorder, err := book.GetOrder("foobar")
	if err == nil {
		t.Fatal("Expected getting an non-existent order to return an error")
	}

	order := Order{ID: "foobar", Price: 100, Side: SIDE_BUY}
	book.PlaceOrder(order, 10)

	sorder, err = book.GetOrder("foobar")
	if err != nil {
		t.Fatalf("Expected getting an order after placing it to return an order, instead got %s", err.Error())
	}
	if sorder.Order != order {
		t.Fatalf("Expected placed order %s to equal retrieved order %s", sorder.Order, order)
	}
	if sorder.State != STATE_PENDING {
		t.Fatalf("Expected just placed order %s to have pending state", sorder)
	}
	if sorder.Size != 10 {
		t.Fatalf("Expected order size %d to be 10", sorder.Size)
	}
}

func TestMutatingSingleOrder(t *testing.T) {
	book := NewInMemoryOrderBook()
	order := Order{ID: "foobar", Price: 100, Side: SIDE_BUY}
	book.PlaceOrder(order, 10)

	mut := OrderStateChange{
		State: STATE_OPEN,
	}
	errs := book.MutateOrder("foobar", []OrderMutation{mut})
	if errs != nil {
		t.Fatalf("Unexpected error mutating order book: %s", errs)
	}
	sorder, err := book.GetOrder("foobar")
	if sorder.State != STATE_OPEN {
		t.Fatalf("Mutation failed to apply. Expected state %s to be %s", sorder.State, STATE_OPEN)
	}

	mut = OrderStateChange{
		State: STATE_OPEN,
	}
	errs = book.MutateOrders("bazbar", []OrderMutation{mut})
	if errs == nil {
		t.Fatal("Expected state mutation on non-existent order to be invalid")
	}

	mut = OrderStateChange{
		State: "kjfdslakfdjsalfkjdslkfdsa",
	}
	errs = book.MutateOrders("foobar", []OrderMutation{mut})
	if err == nil {
		t.Fatal("Expected state mutation to an invalid order state to be invalid")
	}

	sizemut := OrderSizeChange{
		NewSize: 11,
	}
	errs = book.MutateOrders("foobar", []OrderMutation{sizemut})
	if errs != nil {
		t.Fatalf("Unexpected errors mutating order: %s", errs)
	}
	sorder, err = book.GetOrder("foobar")
	if err != nil {
		t.Fatalf("Failed to get mutated order: %s", err.Error())
	}
	if sorder.Size != 11 {
		t.Fatalf("Expected mutated order size %d to be 11", sorder.Size)
	}

	sizemut_new := OrderSizeChange{
		NewSize: 20,
		Time:    time.Unix(1),
	}
	sizemut_old := OrderSizeChange{
		NewSize: 15,
		Time:    time.Unix(0),
	}
	errs = book.MutateOrders("foobar", []OrderMutation{sizemut_new, sizemut_old})
	if errs != nil {
		t.Fatalf("Unexpected errors mutating order: %s", errs)
	}
	sorder, err = book.GetOrder("foobar")
	if err != nil {
		t.Fatalf("Failed to get mutated order: %s", err.Error())
	}
	if sorder.Size != 20 {
		t.Fatalf("Mutations failed to respect time ordering. Expected order size %d to be 20", sorder.Size)
	}
}