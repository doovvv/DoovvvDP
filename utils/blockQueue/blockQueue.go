package blockqueue

type BlockingQueue struct {
	queue chan any
}

func NewBlockingQueue(capacity int) *BlockingQueue {
	return &BlockingQueue{
		queue: make(chan any, capacity),
	}
}

func (q *BlockingQueue) Put(item any) {
	q.queue <- item
}

func (q *BlockingQueue) Get() any {
	return <-q.queue
}

func (q *BlockingQueue) Size() int {
	return len(q.queue)
}
