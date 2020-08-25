package linkflow

import (
	"io"

	"golang.org/x/sync/errgroup"
)

type LinkFlow struct {
	readHead  io.Reader
	writeTail io.Writer
	stop      []func(err error)
	pipe      []io.ReadWriter
}

func New() *LinkFlow {
	return &LinkFlow{}
}

func (l *LinkFlow) ReadHead(r io.Reader, stop func(err error)) *LinkFlow {
	l.readHead = r
	l.stop = append(l.stop, stop)
	return l
}

func (l *LinkFlow) Pipe(rw io.ReadWriter, stop func(err error)) *LinkFlow {
	l.pipe = append(l.pipe, rw)
	l.stop = append(l.stop, stop)
	return l
}

func (l *LinkFlow) WriteTail(w io.Writer) *LinkFlow {
	l.writeTail = w
	return l
}

func (l *LinkFlow) callStop(err error, readIndex int) {
	if l.stop[readIndex] != nil {
		l.stop[readIndex](err)
	}

}

func (l *LinkFlow) Ok() error {
	var g errgroup.Group

	if l.readHead == nil {
		panic("not call ReadHead function")
	}

	if l.writeTail == nil {
		panic("not call WriteTail function")
	}

	if len(l.pipe) == 0 {
		g.Go(func() error {
			readIndex := 0
			_, err := io.Copy(l.writeTail, l.readHead)
			l.callStop(err, readIndex)
			return err
		})
	}

	for k, v := range l.pipe {
		g.Go(func() error {
			var r io.Reader = l.readHead
			var w io.Writer = v

			var readIndex int
			if k > 0 {
				readIndex = k // K-1+1
				r = l.pipe[k-1]
			}

			_, err := io.Copy(w, r)
			l.callStop(err, readIndex)
			return err
		})
	}

	if len(l.pipe) > 0 {
		g.Go(func() error {
			var r io.Reader = l.pipe[len(l.pipe)-1]
			var w io.Writer = l.writeTail

			length := len(l.stop)

			_, err := io.Copy(w, r)
			l.callStop(err, length-2)
			return err

		})
	}

	return g.Wait()
}
