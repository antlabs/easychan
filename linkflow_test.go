package linkflow

import (
	"context"
	"fmt"
	"io"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
)

type reqNet struct {
	total int32
}

var total int32

func (net *reqNet) Read(p []byte) (n int, err error) {
	atomic.AddInt32(&total, 1)
	if atomic.LoadInt32(&total) >= 100 {
		fmt.Printf("network bye bye\n")
		return 0, io.EOF
	}

	n = copy(p, "hello")
	atomic.AddInt32(&net.total, int32(n))
	fmt.Printf("0.wirte total:%d\n", net.total)
	return
}

type audioTran struct {
	rw         chan []byte
	ctx        context.Context
	cancel     context.CancelFunc
	writeTotal int32
	readTotal  int32
}

func audioTranNew() *audioTran {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	return &audioTran{rw: make(chan []byte), ctx: ctx, cancel: cancel}
}

func (a *audioTran) Write(p []byte) (n int, err error) {
	defer func() {
		fmt.Printf("2.write:%s, ok:%d, err:%v\n", p, atomic.AddInt32(&a.writeTotal, int32(len(p))), err)
	}()

	select {
	case <-a.ctx.Done():
		return 0, io.EOF
	case a.rw <- append([]byte{}, p...):
		return len(p), nil
	}

	return 0, io.EOF
}

func (a *audioTran) Read(p []byte) (n int, err error) {
	defer func() {
		fmt.Printf("2.read:%s, ok:%d, err:%v\n", p, atomic.AddInt32(&a.readTotal, int32(n)), err)
	}()

	select {
	case <-a.ctx.Done():
		return 0, io.EOF
	case p2, ok := <-a.rw:
		if !ok {
			return 0, io.EOF
		}
		copy(p, p2)
		return len(p2), nil
	}

	return 0, io.EOF
}

func (a *audioTran) stop() {
	a.cancel()
	fmt.Printf("audioTran stop \n")
}

type rspNet struct {
	total int32
}

func (a *rspNet) Write(p []byte) (n int, err error) {
	atomic.AddInt32(&a.total, int32(len(p)))
	fmt.Printf("3.rspnet:%s:%d\n", p, a.total)
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
