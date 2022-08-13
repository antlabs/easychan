package easychan

import (
	"context"

	"golang.org/x/exp/constraints"
)

func Seq[T constraints.Integer](ctx context.Context) (c chan T) {
	c = make(chan T)
	go func() {
		var i T
		for ; ; i++ {
			c <- i
		}
	}()
	return c
}
