package queue

import (
	"net/url"
	"sync"
)

type Queue struct {
	queue chan *url.URL
	urls []*url.URL
	wg sync.WaitGroup
}

func NewQueue() *Queue {
	return &Queue{
		queue: make(chan *url.URL),
	}
}

// AddBulk is a wrapper around Add to allow for bulk adding of
// URL's to the queue.
func (q *Queue) AddBulk(us []*url.URL) bool {
	added := false
	for _, u := range us {
		if q.Add(u) {
			added = true
		}
	}

	return added
}

// Add will add a URL to the queue for processing and exclude any
// urls that have already been added to the queue. Will return true
// for any unique URL added.
func (q *Queue) Add(u *url.URL) bool {
	if u == nil {
		return false
	}

	for _, eu := range q.urls {
		if eu.String() == u.String() {
			return false
		}
	}

	q.wg.Add(1)
	q.urls = append(q.urls, u)

	go func() {
		q.queue <- u
	}()

	return true
}

// URL will return a channel for consuming URLs
func (q *Queue) URL() <-chan *url.URL {
	return q.queue
}

// URLs will return all current URL's that have been successfully
// added to the queue.
func (q *Queue) URLs() []*url.URL {
	return q.urls
}

func (q *Queue) Done()  {
	q.wg.Done()
}

func (q *Queue) Wait() {
	q.wg.Wait()
}

func (q *Queue) Close() {
	close(q.queue)
}