package bufinsert

import (
	"gorm.io/gorm"
	"sync"
)

type BufferedInserter struct {
	db     *gorm.DB
	buffer []interface{}
	count  int
	locker sync.Mutex
}

func NewInserter(db *gorm.DB, bufferSize int) *BufferedInserter {
	return &BufferedInserter{
		db:     db,
		buffer: make([]interface{}, bufferSize),
		count:  0,
	}
}

// Insert inserts the data into the buffer or into the database(when the buffer is full).
func (receiver *BufferedInserter) Insert(data ...interface{}) {
	receiver.locker.Lock()
	defer receiver.locker.Unlock()

	for _, d := range data {
		receiver.buffer[receiver.count] = d
		receiver.count++

		if full := receiver.count == len(receiver.buffer); full {
			receiver.db.CreateInBatches(receiver.buffer, len(receiver.buffer))
			receiver.count = 0
		}
	}
}

// Flush inserts any buffered data to the database.
func (receiver *BufferedInserter) Flush() {
	receiver.locker.Lock()
	defer receiver.locker.Unlock()

	buffered := receiver.buffer[:receiver.count]
	receiver.db.CreateInBatches(buffered, len(buffered))
	receiver.count = 0
}

// Size returns the size of the underlying buffer in bytes.
func (receiver *BufferedInserter) Size() int {
	receiver.locker.Lock()
	defer receiver.locker.Unlock()

	return len(receiver.buffer)
}

// Available returns how many bytes are unused in the buffer.
func (receiver *BufferedInserter) Available() int {
	receiver.locker.Lock()
	defer receiver.locker.Unlock()

	return len(receiver.buffer) - receiver.count
}

// Buffered returns the number of bytes that have been written into the current buffer.
func (receiver *BufferedInserter) Buffered() int {
	receiver.locker.Lock()
	defer receiver.locker.Unlock()

	return receiver.count
}
