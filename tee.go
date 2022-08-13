package easychan

import "context"

// Tee 函数有点类似
func Tee[T any](ctx context.Context, in chan T) (out1, out2 chan T) {
	out1 = make(chan T)
	out2 = make(chan T)

	go func() {

		defer func() {
			close(out1)
			close(out2)
		}()

		for {
			out1 := out1
			out2 := out2
			select {
			case data, ok := <-in:
				if !ok {
					return
				}

				for i := 0; i < 2; i++ {
					select {
					case out1 <- data:
						out1 = nil
					case out2 <- data:
						out2 = nil
					case <-ctx.Done():
						return
					}
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	return out1, out2
}
