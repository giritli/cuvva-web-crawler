package queue

import (
	"net/url"
	"reflect"
	"testing"
)

func TestQueue_Add(t *testing.T) {
	q := NewQueue()

	want := &url.URL{
		Scheme:     "https",
		Host:       "test",
		Path:       "hello",
	}

	q.Add(want)
	got := <- q.URL()

	if !reflect.DeepEqual(want, got) {
		t.Errorf("Add() = %v, want %v", got, want)
	}

	urls := q.URLs()

	if want := []*url.URL{want}; !reflect.DeepEqual(want, urls) {
		t.Errorf("Add() = %v, want %v", urls, want)
	}
}