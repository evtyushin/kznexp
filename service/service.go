package service

import (
	"context"
	"errors"
	"time"
)

type Service struct {
	Count       int
	IntervalSec int
}

type Item struct{}

type Batch []*Item

func (s *Service) Process(ctx context.Context, batch *Batch, tasks chan *Item) error {

	errs := make(chan error)
	go func() {
		for _, b := range *batch {
			select {
			case <-ctx.Done():
				errs <- errors.New("context canceled")
				return
			default:
				tasks <- b
			}
		}
		errs <- nil
	}()

	return <-errs
}

func (s *Service) CreateChanTasks() chan *Item {
	return make(chan *Item, s.Count)
}

func (s *Service) StartProcessing(tasks chan *Item) {

	// по условию на обработку пакета тратится определенное время
	time.Sleep(time.Duration(s.IntervalSec) * time.Second)

	for {
		<-tasks
		// do something
	}
}

func New(count int, interval int) *Service {

	return &Service{
		Count:       count,
		IntervalSec: interval,
	}
}
