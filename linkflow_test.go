package linkflow

import (
	"context"
	"fmt"
	"io"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
)

type reqNet struct{}

var total int32

func (net *reqNet) Read(p []byte) (n int, err error) {
	atomic.AddInt32(&total, 1)
	if atomic.LoadInt32(&total) == 100 {
		return 0, io.EOF
	}

	n = copy(p, "hello")
	return
}

type audioTran struct {
	rw     chan []byte
	ctx    context.Context
	cancel context.CancelFunc
}

func audioTranNew() *audioTran {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	return &audioTran{rw: make(chan []byte), ctx: ctx, cancel: cancel}
}

func (a *audioTran) Read(p []byte) (n int, err error) {
	defer fmt.Printf("read:%s, ok\n", p)
	select {
	case <-a.ctx.Done():
		return 0, io.EOF
	default:
	}
	p2 := <-a.rw
	return len(p2), nil
}

func (a *audioTran) Write(p []byte) (n int, err error) {
	defer fmt.Printf("write:%s\n", p)

	select {
	case <-a.ctx.Done():
		return 0, io.EOF
	default:
	}
	a.rw <- append([]byte{}, p...)
	return len(p), nil
}

func (a *audioTran) stop() {
	a.cancel()
}

type rspNet struct{}

func (a *rspNet) Write(p []byte) (n int, err error) {
	fmt.Printf("rspnet:%s\n", p)
	return len(p), nil
}

func Test_LinkFlow(t *testing.T) {
	audio := audioTranNew()

	err := New().ReadHead(&reqNet{}, func(err error) {
		audio.stop()
	}).Pipe(audio, func(err error) {

	}).WriteTail(&rspNet{}).Ok()

	assert.NoError(t, err)
}
