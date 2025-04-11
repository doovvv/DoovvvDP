package utils

import (
	"errors"
	"sync"
	"time"
)

const (
	workerBits uint8 = 10
	numberBits uint8 = 12
	workerMax  int64 = (1 << workerBits) - 1
	numberMax  int64 = (1 << numberBits) - 1
	epoch      int64 = 1609459200000 // 2021-01-01 00:00:00 UTC
)

type IDWorker struct {
	lock          sync.Mutex
	workerID      int64
	lastTimestamp int64
	sequence      int64
}

func NewIdWorker(workerID int64) (*IDWorker, error) {
	if workerID < 0 || workerID > workerMax {
		return nil, errors.New("worker ID out of range")
	}
	return &IDWorker{
		workerID:      workerID,
		lastTimestamp: -1,
		sequence:      0,
	}, nil
}
func (w *IDWorker) Generate() (int64, error) {
	w.lock.Lock()
	defer w.lock.Unlock()

	timestamp := time.Now().UnixMilli()
	if timestamp < w.lastTimestamp {
		return 0, errors.New("timestamp has gone back")
	} else if timestamp == w.lastTimestamp {
		w.sequence = (w.sequence + 1) & numberMax
		if w.sequence == 0 { // 序列号用完，等待下一毫秒
			for timestamp <= w.lastTimestamp {
				timestamp = time.Now().UnixMilli()
			}
		}
	}
	w.lastTimestamp = timestamp
	id := ((timestamp - epoch) << (workerBits + numberBits)) | (w.workerID << numberBits) | w.sequence
	return id, nil
}
