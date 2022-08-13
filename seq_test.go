package easychan

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Seq(t *testing.T) {
	ctx, cancel := context.WithCancel(context.TODO())
	in := Seq[int](ctx)
	i := 0

	defer cancel()

	got := make([]int, 0)
	need := []int{0, 1, 2, 3, 4}

	for v := range in {
		got = append(got, v)
		if i == 4 {
			cancel()
			assert.Equal(t, need, got)
			return
		}
		i++
	}
}
