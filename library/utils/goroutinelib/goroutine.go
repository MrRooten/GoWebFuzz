package goroutinelib

import "sync"

type HandlerFunction func(vargs []interface{})
type Task struct {
	handler HandlerFunction
	params  []interface{}
}

type GoroutinePool struct {
	size               int
	capacity           int
	requestGoroutineCh chan int
	wg                 *sync.WaitGroup
}

func (pool *GoroutinePool) NewGoroutinePool(num int) {
	pool.capacity = num
	pool.requestGoroutineCh = make(chan int, num)
	pool.wg = new(sync.WaitGroup)

	for i := 0; i < num; i++ {
		pool.requestGoroutineCh <- i
	}
}

func (pool *GoroutinePool) RunTask(function HandlerFunction, vargs []interface{}) {

	poolId := <-pool.requestGoroutineCh
	pool.wg.Add(1)
	pool.size += 1

	go func() {
		function(vargs)
		defer func() {
			pool.wg.Done()
			pool.size -= 1
			pool.requestGoroutineCh <- poolId
		}()
	}()
}

func (pool *GoroutinePool) WaitTask() {
	pool.wg.Wait()
}

func (pool *GoroutinePool) RenewTaskSize(size int) {
	pool.wg.Wait()
	pool.requestGoroutineCh = make(chan int, size)
	pool.size = size
	pool.capacity = size
}
