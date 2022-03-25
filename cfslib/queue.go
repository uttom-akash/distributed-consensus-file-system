package cfslib

import (
	"container/list"
)

type Queue struct {
	linkedList *list.List
}

func NewQueue() *Queue {
	return &Queue{
		linkedList: list.New(),
	}
}

func (q *Queue) Push(element string) error {
	q.linkedList.PushBack(element)

	return nil
}

func (q *Queue) Front() interface{} {
	front := q.linkedList.Front()

	if front == nil {
		return nil
	}

	return front.Value
}

func (q *Queue) Pop() {
	front := q.linkedList.Front()

	if front != nil {
		q.linkedList.Remove(front)
	}

}

func (q *Queue) Size() int {
	return q.linkedList.Len()
}

func (q *Queue) IsEmpty() bool {
	return q.linkedList.Len() == 0
}
