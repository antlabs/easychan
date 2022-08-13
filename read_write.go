package easychan

import "context"

// 无阻塞读
func AsyncRead[T any](c chan T) (v T, ok bool) {
	select {
	case v, ok = <-c:
	default:
	}
	return
}

// 无阻塞写
func AsyncWrite[T any](c chan T, v T) {
	select {
	case c <- v:
	default:
	}
}

// 带context.Context的读
func ReadContext[T any](ctx context.Context, c chan T) (v T, ok bool) {
	select {
	case v, ok = <-c:
	case <-ctx.Done():
	}
	return
}

// 带context.Context的写
func WriteContext[T any](ctx context.Context, c chan T, v T) {
	select {
	case c <- v:
	case <-ctx.Done():
	}
}
