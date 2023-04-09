package batchprocessor

import (
	"context"
	"sync"

	"github.com/pkg/errors"
)

type aggFun[In any, Out any] func([]In) ([]Out, error)
type resFun[In any, Out any] func(In, []Out) (Out, error)
type batchesList[In any, Out any] map[int]*batch[In, Out]

type Processor[In any, Out any] struct {
	options            Options
	aggregateFunctions map[string]aggFun[In, Out]
	resultFunctions    map[string]resFun[In, Out]
	batches            map[string]batchesList[In, Out]
	stop               chan struct{}
	mx                 sync.RWMutex
	finished           sync.WaitGroup
}

func NewProcessor[In any, Out any]() *Processor[In, Out] {
	p := Processor[In, Out]{
		aggregateFunctions: make(map[string]aggFun[In, Out]),
		resultFunctions:    make(map[string]resFun[In, Out]),
		stop:               make(chan struct{}),
		batches:            make(map[string]batchesList[In, Out]),
	}

	p.options = p.options.defaults()

	return &p
}

func (p *Processor[In, Out]) WithOptions(o Options) *Processor[In, Out] {
	if o.MaxDuration != 0 {
		p.options.MaxDuration = o.MaxDuration
	}
	if o.MaxBatchSize != 0 {
		p.options.MaxBatchSize = o.MaxBatchSize
	}

	return p
}

// AssignAggregateFunction function that will process batch operation, giving all results
func (p *Processor[In, Out]) AssignAggregateFunction(op string, fun aggFun[In, Out]) {
	p.mx.Lock()
	defer p.mx.Unlock()
	p.aggregateFunctions[op] = fun
}

// AssignAggregateFunction function that will process batch operation, giving all results
func (p *Processor[In, Out]) AssignResultFunction(op string, fun resFun[In, Out]) {
	p.mx.Lock()
	defer p.mx.Unlock()
	p.resultFunctions[op] = fun
}

// Run runs concurrently operation in batch
func (p *Processor[In, Out]) Run(ctx context.Context, op string, in In) ([]Out, error) {
	select {
	case <-ctx.Done():
		p.Stop()
	default:
	}
	select {
	case <-p.stop:
		return nil, ErrProcessorStopped
	default:
	}

	result := operation[In, Out]{
		result: make(chan BatchResult[Out]),
		input:  in,
	}

	if _, ok := p.aggregateFunctions[op]; !ok {
		return nil, ErrNoAggregateFunctionFound
	}
	p.process(result, p.aggregateFunctions[op], op)
	res := <-result.result

	if res.Error != nil {
		return nil, res.Error
	}

	return res.Value, nil
}

// Run runs concurrently operation in batch, returns only one result, same backend
func (p *Processor[In, Out]) RunSingleResult(ctx context.Context, op string, in In) (*Out, error) {
	select {
	case <-ctx.Done():
		p.Stop()
	default:
	}
	select {
	case <-p.stop:
		return nil, ErrProcessorStopped
	default:
	}

	result := operation[In, Out]{
		result: make(chan BatchResult[Out]),
		input:  in,
	}

	if _, ok := p.aggregateFunctions[op]; !ok {
		return nil, ErrNoAggregateFunctionFound
	}
	p.process(result, p.aggregateFunctions[op], op)

	res := <-result.result

	if res.Error != nil {
		return nil, errors.Wrap(res.Error, "error from aggregate")
	}

	if _, ok := p.resultFunctions[op]; !ok {
		return nil, ErrNoAggregateFunctionFound
	}
	singleResult, err := p.resultFunctions[op](in, res.Value)
	if err != nil {
		return nil, errors.Wrap(err, "error from result")
	}

	return &singleResult, nil
}

func (p *Processor[In, Out]) deleteRecordCallback(opName string, batchIdx int) {
	p.mx.Lock()
	defer p.mx.Unlock()

	if batches, ok := p.batches[opName]; ok {
		delete(batches, batchIdx)
	}
}

// process managing batches, removes old and assign new
func (p *Processor[In, Out]) process(op operation[In, Out], fun aggFun[In, Out], opName string) {
	p.mx.Lock()
	defer p.mx.Unlock()
	batches, ok := p.batches[opName]
	if !ok {
		p.batches[opName] = make(batchesList[In, Out])
		batches = p.batches[opName]
	}

	for k, v := range batches {
		select {
		case <-v.closed:
			delete(batches, k)
			continue
		default:
		}

		if v.isAvailable() {
			err := v.addElement(op)
			if err != nil {
				continue
			}
			return
		}
	}

	w := newBatch(p.options.MaxDuration, p.options.MaxBatchSize, fun, opName)

	i := 0
	for {
		// first available place
		if _, ok := batches[i]; ok {
			i++
			continue
		}

		batches[i] = w
		batches[i].inChan <- op
		go p.startBatch(w)

		w.setCallback(func() {
			p.deleteRecordCallback(opName, i)
		})

		return
	}
}

// startBatch batch start await
func (p *Processor[In, Out]) startBatch(w *batch[In, Out]) {
	p.finished.Add(1)
	defer p.finished.Done()

	w.process()
	w.timer.Stop()
}

func (p *Processor[In, Out]) Stop() {
	close(p.stop)
	p.finished.Wait()
}
