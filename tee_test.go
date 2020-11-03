package linkflow

import (
	"context"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Tee(t *testing.T) {
	voice := make(chan interface{})
	out1, out2 := Tee(context.Background(), voice)

	var wg sync.WaitGroup

	wg.Add(2)
	defer wg.Wait()

	max := 100
	need := make([]int, max)

	for i := 0; i < max; i++ {
		need[i] = i
	}

	go func() {
		for _, v := range need {
			voice <- v
		}

		close(voice)

	}()

	go func() {
		var got []int
		defer wg.Done()
		for {
			select {
			case d, ok := <-out1:
				if !ok {
					assert.Equal(t, need, got)
					return
				}
				got = append(got, d.(int))
			}
		}
	}()

	go func() {
		defer wg.Done()
		var got []int
		for {
			select {
			case d, ok := <-out2:
				if !ok {
					assert.Equal(t, need, got)
					return
				}
				got = append(got, d.(int))
			}
		}
	}()
}
