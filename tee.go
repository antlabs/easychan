package linkflow

import "context"

// Tee 函数有点类似
func Tee(ctx context.Context, in chan interface{}) (out1, out2 chan interface{}) {
	out1 = make(chan interface{})
	out2 = make(chan interface{})

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
