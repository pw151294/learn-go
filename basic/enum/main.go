package main

import "fmt"

type Status int
type OrderStatus int

const (
	Pending Status = iota
	Approved
	Rejected
)

const (
	StatusPending OrderStatus = iota
	StatusProcessing
	StatusShipped
	StatusCompleted
	StatusCancelled
)

func (s *Status) String() string {
	return [...]string{"Pending", "Approved", "Rejected"}[*s]
}

func (os *OrderStatus) String() string {
	return [...]string{"Pending", "Processing", "Shipped", "Completed", "Cancelled"}[*os]
}

func (os *OrderStatus) IsValid() bool {
	return *os >= StatusPending && *os <= StatusCompleted
}

func main() {
	s := Pending
	fmt.Println(s, s.String())

	os := StatusCancelled
	fmt.Println(os.IsValid())
}
