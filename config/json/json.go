package main

import (
	"encoding/json"
	"fmt"
)

type OrderItem struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

type Order struct {
	ID         string      `json:"id"`
	Name       string      `json:"name,omitempty"`
	Items      []OrderItem `json:"items"`
	Quantity   int         `json:"quantity"`
	TotalPrice float64     `json:"total_price"`
}

const jsonStr = `{"id":"1","items":[{"id":1,"name":"1","price":99.99},{"id":2,"name":"2","price":99.99}],"quantity":10,"total_price":99.99}`

func unMarshal() {
	var o Order
	err := json.Unmarshal([]byte(jsonStr), &o)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", o)
}

func marshal() {
	o := Order{
		ID:         "1",
		Quantity:   10,
		TotalPrice: 99.99,
		Items: []OrderItem{
			{
				ID:    1,
				Name:  "1",
				Price: 99.99,
			},
			{
				ID:    2,
				Name:  "2",
				Price: 99.99,
			},
		},
	}

	bytes, err := json.Marshal(o)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(bytes))
}

func main() {
	marshal()
	unMarshal()

	m := struct {
		ID    string `json:"id"`
		Items []struct {
			ID    int     `json:"id"`
			Name  string  `json:"name"`
			Price float64 `json:"price"`
		}
	}{}
	err := json.Unmarshal([]byte(jsonStr), &m)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v %+v\n", m.ID, m.Items[0])

	arr := []int{1, 2, 3, 4}
	bytes, _ := json.Marshal(arr)
	fmt.Println(string(bytes))
}
