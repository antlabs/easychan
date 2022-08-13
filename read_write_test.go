package easychan

import (
	"context"
	"testing"
)

func Test_AsyncRead(t *testing.T) {
	c := make(chan bool)
	_, _ = AsyncRead(c)
}

func Test_AsyncWrite(t *testing.T) {
	c := make(chan bool)
	AsyncWrite(c, true)
}

func Test_ReadContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.TODO())
	cancel()
	_, _ = ReadContext(ctx, make(chan bool))
}

func Test_WriteContext(t *testing.T) {

	ctx, cancel := context.WithCancel(context.TODO())
	cancel()
	WriteContext(ctx, make(chan bool), true)
}
