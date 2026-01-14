package filePlayer

import "music-bot/audio/types"

type Queue struct {
	baseArray []types.Song
	head      int
	tail      int
	length    int
	capacity  int
}

const ScalingFactor = 2
const QueueStartSize = 64

func NewQueue() *Queue {
	return &Queue{
		baseArray: make([]types.Song, QueueStartSize),
		head:      0,
		tail:      0,
		length:    0,
		capacity:  QueueStartSize,
	}
}

func (q *Queue) Enqueue(song types.Song) {
	if q.length >= q.capacity {
		q.resize()
	}
	q.baseArray[q.tail] = song
	q.tail = (q.tail + 1) % len(q.baseArray)
	q.length++
}

func (q *Queue) Dequeue() types.Song {
	if q.length == 0 {
		return nil
	}
	current := q.baseArray[q.head]
	q.baseArray[q.head] = nil
	q.head = (q.head + 1) % len(q.baseArray)
	q.length--
	return current
}

func (q *Queue) DequeueLastAdded() types.Song {
	if q.length == 0 {
		return nil
	}
	lastAdded := q.baseArray[q.tail]
	q.baseArray[q.tail] = nil
	return lastAdded
}

func (q *Queue) Peek() types.Song {
	if q.length == 0 {
		return nil
	}
	return q.baseArray[q.head]
}

func (q *Queue) PeekLastAdded() types.Song {
	if q.length == 0 {
		return nil
	}
	return q.baseArray[q.tail]
}

func (q *Queue) Length() int {
	return q.length
}

func (q *Queue) Find(s string) (types.Song, bool) {
	//TODO implement me
	panic("implement me")
}

func (q *Queue) resize() {
	capacity := (q.capacity * ScalingFactor) + 1
	newBaseArray := make([]types.Song, capacity)
	for i := 0; i < q.length; i++ {
		newBaseArray[i] = q.baseArray[(q.head+i)%q.capacity]
	}
	q.capacity = capacity
	q.baseArray = newBaseArray
	q.head = 0
	q.tail = q.length
}

func (q *Queue) desize() {
	capacity := (q.capacity / ScalingFactor) + 1
	newBaseArray := make([]types.Song, capacity)
	for i := 0; i < q.length; i++ {
		newBaseArray[i] = q.baseArray[(q.head+i)%q.capacity]
	}
	q.capacity = capacity
	q.baseArray = newBaseArray
	q.head = 0
	q.tail = q.length
}
