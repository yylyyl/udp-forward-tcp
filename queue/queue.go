package queue

import (
	"container/list"
	"errors"
	"sync"
)

type Queue struct {
	l    *list.List
	m    sync.Mutex
	wg   sync.WaitGroup
	wait bool
}

func NewQueue() *Queue {
	return &Queue{
		l: list.New(),
	}
}

func (q *Queue) Push(item interface{}) error {
	q.m.Lock()
	defer q.m.Unlock()

	if q.l.Len() >= 1024 {
		q.l.Init()
		return errors.New("queue too long, cleared")
	}

	q.l.PushFront(item)

	if q.wait {
		q.wait = false
		q.wg.Done()
	}

	return nil
}

func (q *Queue) Pop() interface{} {
	q.m.Lock()
	defer q.m.Unlock()

	if q.l.Len() == 0 {
		q.wait = true
		q.m.Unlock()
		q.wg.Add(1)
		q.wg.Wait()
		q.m.Lock()
	}

	ele := q.l.Back()
	q.l.Remove(ele)

	return ele.Value
}
