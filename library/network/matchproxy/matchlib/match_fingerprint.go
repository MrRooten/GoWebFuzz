package matchlib

import (
	"gowebfuzz/library/utils/goroutinelib"
	"sync"
)

type MatchQueue struct {
	locker *sync.Mutex
	queue  []*Pair
}

func (queue *MatchQueue) Push(pair *Pair) {
	queue.locker.Lock()
	defer queue.locker.Unlock()
	queue.queue = append(queue.queue, pair)
}

func (queue *MatchQueue) PopFront() *Pair {
	queue.locker.Lock()
	defer queue.locker.Unlock()
	if len(queue.queue) == 0 {
		return nil
	}
	res := queue.queue[0]
	queue.queue = queue.queue[1:]
	return res
}

func (listener *MessageListener) SendMessage(message *Message) {
	listener.MessageChannel <- message
}

var FingerPrintMatchPool *goroutinelib.GoroutinePool

func MatchFingerPrint(records *_Records) {
	if FingerPrintMatchPool == nil {
		FingerPrintMatchPool = &goroutinelib.GoroutinePool{}
		FingerPrintMatchPool.NewGoroutinePool(100)
	}

}
