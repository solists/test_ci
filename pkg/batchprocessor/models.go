package batchprocessor

import (
	"time"

	"github.com/pkg/errors"
)

type operation[In any, Out any] struct {
	result chan BatchResult[Out]
	input  In
}

type Options struct {
	MaxDuration  time.Duration
	MaxBatchSize uint64
}

const maxTime = 30 * time.Millisecond

func (o Options) defaults() Options {
	o.MaxDuration = maxTime
	o.MaxBatchSize = 20

	return o
}

type BatchResult[Out any] struct {
	Value []Out
	Error error
}

var ErrProcessorStopped = errors.New("processor is stopped")
var ErrNoAggregateFunctionFound = errors.New("no aggregate function found")
var ErrNoResultFunctionFound = errors.New("no result function found")
var ErrOpCancelled = errors.New("operation was canceled ")
