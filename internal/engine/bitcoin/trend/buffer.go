package trend

import "sync"

type Buffer struct {
	mu        sync.Mutex
	snapshots []Snapshot
	maxSize   int
}

func NewBuffer(maxSize int) *Buffer {
	return &Buffer{
		snapshots: make([]Snapshot, 0, maxSize),
		maxSize:   maxSize,
	}
}

func (b *Buffer) Add(s Snapshot) *Snapshot {
	b.mu.Lock()
	defer b.mu.Unlock()

	var prev *Snapshot

	if len(b.snapshots) > 0 {
		tmp := b.snapshots[len(b.snapshots)-1]
		prev = &tmp
	}

	if len(b.snapshots) >= b.maxSize {
		b.snapshots = b.snapshots[1:]
	}

	b.snapshots = append(b.snapshots, s)

	return prev
}

func (b *Buffer) All() []Snapshot {
	b.mu.Lock()
	defer b.mu.Unlock()

	out := make([]Snapshot, len(b.snapshots))
	copy(out, b.snapshots)
	return out
}
