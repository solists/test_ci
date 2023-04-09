package batchprocessor

import (
	"context"
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestProcessorRun(t *testing.T) {
	p := NewProcessor[int, int]().WithOptions(Options{MaxDuration: maxTime})
	p.AssignAggregateFunction("add", func(in []int) ([]int, error) {
		res := make([]int, 0, len(in))
		for _, v := range in {
			res = append(res, v+1)
		}
		return res, nil
	})

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := p.Run(ctx, "add", 1)
	assert.NoError(t, err)
	assert.ElementsMatch(t, []int{2}, res)

	p.Stop()
}

func TestProcessorStopped(t *testing.T) {
	p := NewProcessor[int, int]().WithOptions(Options{MaxDuration: maxTime})
	p.Stop()

	_, err := p.Run(context.Background(), "add", 1)
	assert.EqualError(t, err, ErrProcessorStopped.Error())
}

func TestProcessorCancelled(t *testing.T) {
	p := NewProcessor[int, int]().WithOptions(Options{MaxDuration: maxTime})
	p.AssignAggregateFunction("add", func(in []int) ([]int, error) {
		res := make([]int, 0, len(in))
		for _, v := range in {
			res = append(res, v+1)
		}
		return res, nil
	})

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := p.Run(ctx, "add", 1)
	assert.EqualError(t, err, ErrProcessorStopped.Error())
}

func TestEqual(t *testing.T) {
	p := NewProcessor[int, int]().WithOptions(Options{MaxDuration: maxTime * 10})
	p.AssignAggregateFunction("add", func(in []int) ([]int, error) {
		res := make([]int, 0, len(in))
		for _, v := range in {
			res = append(res, v+10)
		}
		return res, nil
	})

	var a []int
	var b []int
	for k := 0; k < 10; k++ {
		// nolint
		a = append(a, rand.Int())
	}

	for i := range a {
		b = append(b, a[i]+10)
		if i == len(a)-1 {
			time.Sleep(10 * time.Millisecond)
			result, err := p.Run(context.Background(), "add", a[i])
			assert.NoError(t, err)
			assert.ElementsMatch(t, b, result)
			continue
		}

		go func(o int) {
			r, err := p.Run(context.Background(), "add", o)
			assert.NoError(t, err)
			// assure, every goroutine gets the same res
			assert.ElementsMatch(t, b, r)
		}(a[i])
	}

	time.Sleep(10 * time.Millisecond)
}

func TestEqualSingle(t *testing.T) {
	p := NewProcessor[int, int]().WithOptions(Options{MaxDuration: maxTime * 10})
	p.AssignAggregateFunction("add", func(in []int) ([]int, error) {
		res := make([]int, 0, len(in))
		for _, v := range in {
			res = append(res, v+10)
		}
		return res, nil
	})

	p.AssignResultFunction("add", func(in int, outs []int) (int, error) {
		for _, v := range outs {
			if v-10 == in {
				return in, nil
			}
		}

		return 0, errors.New("not found")
	})

	var a []int
	var b []int
	for k := 0; k < 10; k++ {
		// nolint
		a = append(a, rand.Int())
	}

	for i := range a {
		b = append(b, a[i]+10)
		if i == len(a)-1 {
			time.Sleep(10 * time.Millisecond)
			result, err := p.Run(context.Background(), "add", a[i])

			assert.NoError(t, err)
			assert.ElementsMatch(t, b, result)
			continue
		}

		go func(o int) {
			r, err := p.RunSingleResult(context.Background(), "add", o)
			assert.NoError(t, err)
			// assure, every goroutine gets the corresponding result
			assert.Equal(t, o, *r)
		}(a[i])
	}

	time.Sleep(10 * time.Millisecond)
}

// many with batches
func TestEqualMaxSize(t *testing.T) {
	p := NewProcessor[int, int]().WithOptions(Options{MaxDuration: maxTime * 10})
	p.AssignAggregateFunction("add", func(in []int) ([]int, error) {
		res := make([]int, 0, len(in))
		for _, v := range in {
			res = append(res, v+10)
		}
		return res, nil
	})

	var a []int
	var modifiedA []int
	for k := 0; k < 1000; k++ {
		// nolint
		a = append(a, rand.Int())
		modifiedA = append(modifiedA, a[k]+10)
	}

	maxBatchSize := 20
	var b = modifiedA[0:min(len(modifiedA), maxBatchSize)]

	for i := range a {
		if i%maxBatchSize == maxBatchSize-1 {
			time.Sleep(10 * time.Millisecond)
			result, err := p.Run(context.Background(), "add", a[i])
			assert.NoError(t, err)
			assert.ElementsMatch(t, b, result)
			b = modifiedA[min(i+1, len(modifiedA)):min(len(modifiedA), i+1+maxBatchSize)]
			continue
		}

		copySlice := make([]int, len(b), cap(b))
		copy(copySlice, b)
		go func(o int, b []int) {
			r, err := p.Run(context.Background(), "add", o)
			assert.NoError(t, err)
			// assure, every goroutine gets the same res
			assert.ElementsMatch(t, b, r)
		}(a[i], b)
	}

	time.Sleep(10 * time.Millisecond)
}

// many
func TestFullMap(t *testing.T) {
	p := NewProcessor[int, int]().WithOptions(Options{MaxDuration: maxTime * 10})
	p.AssignAggregateFunction("add", func(in []int) ([]int, error) {
		res := make([]int, 0, len(in))
		for _, v := range in {
			res = append(res, v+10)
		}
		return res, nil
	})

	var a []int
	for k := 0; k < 10015; k++ {
		// nolint
		a = append(a, rand.Int())
	}

	wg := &sync.WaitGroup{}

	ch := make(chan int, 1000)

	for i := range a {
		wg.Add(1)
		go func(o int, ch chan int, wg *sync.WaitGroup) {
			defer wg.Done()
			r, err := p.Run(context.Background(), "add", o)
			assert.NoError(t, err)
			for _, v := range r {
				ch <- v
			}
		}(a[i], ch, wg)
	}

	wgOuter := &sync.WaitGroup{}
	wgOuter.Add(1)
	go func(ch chan int, a []int, wg *sync.WaitGroup) {
		defer wg.Done()
		var checkMap = make(map[int]struct{}, len(a))
		for i := range ch {
			if _, ok := checkMap[i-10]; !ok {
				checkMap[i-10] = struct{}{}
			}
		}

		missCounter := 0
		for _, v := range a {
			if _, ok := checkMap[v]; !ok {
				missCounter++
			}
		}

		if missCounter > 0 {
			t.Errorf("missed %v items", missCounter)
			t.Fail()
		}
	}(ch, a, wgOuter)

	wg.Wait()
	close(ch)

	wgOuter.Wait()
	time.Sleep(10 * time.Millisecond)
}

// we wait until timer is up
func TestEqualByTime(t *testing.T) {
	p := NewProcessor[int, int]()
	p.AssignAggregateFunction("add", func(in []int) ([]int, error) {
		res := make([]int, 0, len(in))
		for _, v := range in {
			res = append(res, v+10)
		}
		return res, nil
	})

	var a []int
	var b []int
	for k := 0; k < 100; k++ {
		// nolint
		a = append(a, rand.Int())
	}

	for i := range a {
		b = append(b, a[i]+10)
		if i%5 == 0 {
			result, err := p.Run(context.Background(), "add", a[i])
			time.Sleep(maxTime * 2)
			assert.NoError(t, err)
			assert.ElementsMatch(t, b, result)
			b = []int{}
			continue
		}

		go func(o int) {
			_, err := p.Run(context.Background(), "add", o)
			assert.NoError(t, err)
		}(a[i])
	}

	time.Sleep(10 * time.Millisecond)
}

func TestFuzzyProcessor(t *testing.T) {
	t.Run("TestRun", func(t *testing.T) {
		p := NewProcessor[int, int]().WithOptions(Options{MaxDuration: maxTime})
		aggFun := func(inputs []int) ([]int, error) {
			outputs := []int{}
			for _, input := range inputs {
				outputs = append(outputs, input+1)
			}
			return outputs, nil
		}
		p.AssignAggregateFunction("addOne", aggFun)
		ctx := context.Background()
		outputs, err := p.Run(ctx, "addOne", 1)
		if err != nil {
			t.Error(err)
		}
		if outputs[0] != 2 {
			t.Errorf("expected 2, got %d", outputs[0])
		}
	})

	t.Run("TestStop", func(t *testing.T) {
		p := NewProcessor[int, int]().WithOptions(Options{MaxDuration: maxTime})
		p.Stop()
		_, err := p.Run(context.Background(), "addOne", 1)
		if err != ErrProcessorStopped {
			t.Errorf("expected %s, got %s", ErrProcessorStopped, err)
		}
	})

	time.Sleep(10 * time.Millisecond)
}

func min(a, b int) int {
	if a > b {
		return b
	}
	return a
}
