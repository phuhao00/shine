package chanrpc

import (
	"context"
	"errors"
	"fmt"
	"runtime"

	"github.com/golang/gddo/log"
	"github.com/phuhao00/shine/conf"
)

// one server per goroutine (goroutine not safe)
// one client per goroutine (goroutine not safe)
type Server struct {
	// id -> function
	functions map[uint16]func(data []byte) error
	ChanCall  chan *CallInfo
}

type CallInfo struct {
	f       func([]byte) error
	args    []byte
	chanRet chan *RetInfo
	cb      func(interface{}, error)
}

type RetInfo struct {
	err error
	ret interface{}
	cb  func(interface{}, error)
}

type Client struct {
	s               *Server
	chanSyncRet     chan *RetInfo
	ChanAsynRet     chan *RetInfo
	pendingAsynCall int
}

func NewServer(l int) *Server {
	s := new(Server)
	s.functions = make(map[uint16]func(data []byte) error)
	s.ChanCall = make(chan *CallInfo, l)
	return s
}

func assert(i interface{}) []interface{} {
	if i == nil {
		return nil
	} else {
		return i.([]interface{})
	}
}

// Register you must call the function before calling Open and Go
func (s *Server) Register(id uint16, f func(data []byte) error) {
	if _, ok := s.functions[id]; ok {
		panic(fmt.Sprintf("function id %v: already registered", id))
	}

	s.functions[id] = f
}

func (s *Server) ret(ci *CallInfo, ri *RetInfo) (err error) {
	if ci.chanRet == nil {
		return
	}

	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()

	ri.cb = ci.cb
	ci.chanRet <- ri
	return
}

func (s *Server) exec(ci *CallInfo) (err error) {
	defer func() {
		if r := recover(); r != nil {
			if conf.LenStackBuf > 0 {
				buf := make([]byte, conf.LenStackBuf)
				l := runtime.Stack(buf, false)
				err = fmt.Errorf("%v: %s", r, buf[:l])
			} else {
				err = fmt.Errorf("%v", r)
			}
			s.ret(ci, &RetInfo{err: fmt.Errorf("%v", r)})
		}
	}()

	return ci.f(ci.args)
}

func (s *Server) Exec(ci *CallInfo) {
	err := s.exec(ci)
	if err != nil {
		log.Error(context.Background(), "%v", err)
	}
}

// Go goroutine safe
func (s *Server) Go(id uint16, args []byte) {
	f := s.functions[id]
	if f == nil {
		return
	}

	defer func() {
		recover()
	}()

	s.ChanCall <- &CallInfo{
		f:    f,
		args: args,
	}
}

// Call goroutine safe
func (s *Server) Call(id uint16, args []byte) error {
	return s.Open(0).Call(id, args)
}

func (s *Server) Close() {
	close(s.ChanCall)

	for ci := range s.ChanCall {
		s.ret(ci, &RetInfo{
			err: errors.New("chanrpc server closed"),
		})
	}
}

// Open goroutine safe
func (s *Server) Open(l int) *Client {
	c := NewClient(l)
	c.Attach(s)
	return c
}

func NewClient(l int) *Client {
	c := new(Client)
	c.chanSyncRet = make(chan *RetInfo, 1)
	c.ChanAsynRet = make(chan *RetInfo, l)
	return c
}

func (c *Client) Attach(s *Server) {
	c.s = s
}

func (c *Client) call(ci *CallInfo, block bool) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()

	if block {
		c.s.ChanCall <- ci
	} else {
		select {
		case c.s.ChanCall <- ci:
		default:
			err = errors.New("chanrpc channel full")
		}
	}
	return
}

func (c *Client) f(id uint16) (f func([]byte) error, err error) {
	if c.s == nil {
		err = errors.New("server not attached")
		return
	}

	f = c.s.functions[id]
	if f == nil {
		err = fmt.Errorf("function id %v: function not registered", id)
		return
	}

	return f, nil
}

func (c *Client) Call(id uint16, args []byte) error {
	f, err := c.f(id)
	if err != nil {
		return err
	}

	err = c.call(&CallInfo{
		f:       f,
		args:    args,
		chanRet: c.chanSyncRet,
	}, true)
	if err != nil {
		return err
	}

	ri := <-c.chanSyncRet
	return ri.err
}

func (c *Client) asynCall(id uint16, args []byte, cb func(interface{}, error)) {
	f, err := c.f(id)
	if err != nil {
		c.ChanAsynRet <- &RetInfo{err: err, cb: cb}
		return
	}

	err = c.call(&CallInfo{
		f:       f,
		args:    args,
		chanRet: c.ChanAsynRet,
		cb:      cb,
	}, false)
	if err != nil {
		c.ChanAsynRet <- &RetInfo{err: err, cb: cb}
		return
	}
}

func (c *Client) AsynCall(id uint16, data []byte, cb func(interface{}, error)) {
	if len(data) < 1 {
		panic("callback function not found")
	}

	// too many calls
	if c.pendingAsynCall >= cap(c.ChanAsynRet) {
		execCb(&RetInfo{err: errors.New("too many calls"), cb: cb})
		return
	}
	c.asynCall(id, data, cb)
	c.pendingAsynCall++
}

func execCb(ri *RetInfo) {
	defer func() {
		if r := recover(); r != nil {
			if conf.LenStackBuf > 0 {
				buf := make([]byte, conf.LenStackBuf)
				l := runtime.Stack(buf, false)
				log.Error(context.Background(), "%v: %s", r, buf[:l])
			} else {
				log.Error(context.Background(), "%v", r)
			}
		}
	}()

	ri.cb(ri.ret, ri.err)
	return
}

func (c *Client) Cb(ri *RetInfo) {
	c.pendingAsynCall--
	execCb(ri)
}

func (c *Client) Close() {
	for c.pendingAsynCall > 0 {
		c.Cb(<-c.ChanAsynRet)
	}
}

func (c *Client) Idle() bool {
	return c.pendingAsynCall == 0
}
