package batchprocessor

import (
	"time"

	"github.com/pkg/errors"
)

type batch[In any, Out any] struct {
	timer *time.Timer
	// inputs channel
	inChan chan operation[In, Out]
	// informs if write to inChan is available
	closed chan struct{}
	// function to execute in batch
	aggregateFun   aggFun[In, Out]
	operation      string
	maxBatchSize   int
	deleteCallback func()
}

func newBatch[In any, Out any](
	maxDuration time.Duration,
	maxBatchSize uint64,
	fun aggFun[In, Out],
	opName string,
) *batch[In, Out] {
	return &batch[In, Out]{
		aggregateFun: fun,
		inChan:       make(chan operation[In, Out], maxBatchSize),
		timer:        time.NewTimer(maxDuration),
		closed:       make(chan struct{}),
		operation:    opName,
		maxBatchSize: int(maxBatchSize),
	}
}

func (b *batch[In, Out]) isAvailable() bool {
	return len(b.inChan) < b.maxBatchSize
}

func (b *batch[In, Out]) setCallback(fun func()) {
	b.deleteCallback = fun
}

func (b *batch[In, Out]) addElement(op operation[In, Out]) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New("write to closed chan")
		}
	}()

	b.inChan <- op

	return nil
}

// waits until timer ends or batch will be full
func (b *batch[In, Out]) process() {
	for {
		select {
		case <-b.timer.C:
			close(b.closed)
			close(b.inChan)
			if b.deleteCallback != nil {
				b.deleteCallback()
			}

			b.processBatch()
			return
		default:
		}

		if len(b.inChan) >= b.maxBatchSize {
			close(b.closed)
			close(b.inChan)
			if b.deleteCallback != nil {
				b.deleteCallback()
			}

			b.processBatch()
			return
		}
	}
}

// executes aggregate function on batch
func (b *batch[In, Out]) processBatch() {
	if len(b.inChan) == 0 {
		return
	}

	// input for batch processing function
	in := make([]In, 0, len(b.inChan))
	ops := make([]operation[In, Out], 0, len(b.inChan))
	for chanOut := range b.inChan {
		ops = append(ops, chanOut)
		in = append(in, chanOut.input)
	}

	// process batch
	results, err := b.aggregateFun(in)
	if err != nil {
		for _, v := range ops {
			v.result <- BatchResult[Out]{
				Error: err,
			}
			close(v.result)
		}

		return
	}

	for _, v := range ops {
		v.result <- BatchResult[Out]{
			Value: results,
		}
		close(v.result)
	}
}
