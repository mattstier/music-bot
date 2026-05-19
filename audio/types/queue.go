package types

import "iter"

// iterable queue interface for song playlists
type Queue interface {
	Enqueue(Song)
	Dequeue() Song
	DequeueLastAdded() Song
	Peek() Song
	PeekLastAdded() Song
	Length() int
	Find(string) (Song, bool)
	List() []Song
	All() iter.Seq[Song]
}
