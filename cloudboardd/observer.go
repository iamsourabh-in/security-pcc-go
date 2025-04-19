package main

import (
   "sync"

   "github.com/iamsourabh-in/security-pcc-go/proto/cloudboard/controller"
)

// LoadSubject implements the observer pattern for load metrics.
type LoadSubject struct {
   mu               sync.Mutex
   maxBatchSize     uint32
   currentBatchSize uint32
   observers        []chan *controller.LoadResponse
}

// NewLoadSubject creates a new LoadSubject with the given maximum batch size.
func NewLoadSubject(maxBatchSize int) *LoadSubject {
   return &LoadSubject{
       maxBatchSize: uint32(maxBatchSize),
       observers:    make([]chan *controller.LoadResponse, 0),
   }
}

// Subscribe registers an observer and returns a channel to receive updates,
// and a function to unsubscribe.
func (l *LoadSubject) Subscribe() (<-chan *controller.LoadResponse, func()) {
   ch := make(chan *controller.LoadResponse, 1)
   l.mu.Lock()
   // send initial state
   resp := &controller.LoadResponse{
       MaxBatchSize:     l.maxBatchSize,
       CurrentBatchSize: l.currentBatchSize,
       OptimalBatchSize: l.maxBatchSize,
   }
   ch <- resp
   l.observers = append(l.observers, ch)
   l.mu.Unlock()
   return ch, func() {
       l.mu.Lock()
       defer l.mu.Unlock()
       for i, c := range l.observers {
           if c == ch {
               l.observers = append(l.observers[:i], l.observers[i+1:]...)
               close(c)
               break
           }
       }
   }
}

// notify sends the current state to all observers.
func (l *LoadSubject) notify() {
   l.mu.Lock()
   defer l.mu.Unlock()
   resp := &controller.LoadResponse{
       MaxBatchSize:     l.maxBatchSize,
       CurrentBatchSize: l.currentBatchSize,
       OptimalBatchSize: l.maxBatchSize,
   }
   for _, c := range l.observers {
       select {
       case c <- resp:
       default:
       }
   }
}

// Increment increases the current batch size and notifies observers.
func (l *LoadSubject) Increment() {
   l.mu.Lock()
   if l.currentBatchSize < l.maxBatchSize {
       l.currentBatchSize++
   }
   l.mu.Unlock()
   l.notify()
}

// Decrement decreases the current batch size and notifies observers.
func (l *LoadSubject) Decrement() {
   l.mu.Lock()
   if l.currentBatchSize > 0 {
       l.currentBatchSize--
   }
   l.mu.Unlock()
   l.notify()
}