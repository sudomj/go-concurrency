package main

import (
	"fmt"
	"math/rand"
	"time"
)

type Data struct {
	Name     string
	Price    float64
	Quantity int
}

type Result struct {
	Name  string
	Value float64
	Tax   float64
}

type Event struct {
	name       string
	totalPrice float64
	time       time.Time
}

type Pipeline struct {
	DataCh   chan Data
	ResultCh chan Result
	EventCh  chan Event
}

func NewPipeline(size int) *Pipeline {

	return &Pipeline{
		DataCh:   make(chan Data, size),
		ResultCh: make(chan Result, size),
		EventCh:  make(chan Event, size),
	}
}

func (p *Pipeline) Send(data Data) {
	p.DataCh <- data
}

func (p *Pipeline) Read() {

	defer close(p.ResultCh)

	for data := range p.DataCh {

		res := Result{
			Name:  data.Name,
			Value: data.Price * float64(data.Quantity),
			Tax:   data.Price * 0.15,
		}
		p.ResultCh <- res
	}
}

func (p *Pipeline) Transform() {

	defer close(p.EventCh)

	for result := range p.ResultCh {

		event := Event{
			name:       result.Name,
			totalPrice: result.Value + result.Tax,
			time:       time.Now(),
		}

		p.EventCh <- event
	}
}

func main() {

	products := []string{"Apple", "Rice", "Tomato", "Bean", "Sugar", "Lemon"}
	size := len(products)
	pipeline := NewPipeline(size)
	for index := 0; index < len(products); index++ {

		price := float64(rand.Intn(50)) + rand.Float64()
		data := Data{
			Name:     products[index],
			Quantity: index,
			Price:    price,
		}
		pipeline.Send(data)
	}
	close(pipeline.DataCh)

	go pipeline.Read()
	go pipeline.Transform()

	for event := range pipeline.EventCh {
		fmt.Printf("product %s has a total price of %v\n", event.name, event.totalPrice)
	}
}
