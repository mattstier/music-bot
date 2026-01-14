package types

type Queue interface {
	Enqueue(Song)
	Dequeue() Song
	DequeueLastAdded() Song
	Peek() Song
	PeekLastAdded() Song
	Length() int
	Find(string) (Song, bool)
}
