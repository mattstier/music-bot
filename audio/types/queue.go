package audio

type Queue interface {
	Enqueue() Song
	Dequeue() Song
	Peek() Song
	Length() int
}
