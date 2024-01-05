package runtime

import (
	"context"
	"github.com/samber/lo"
	"sync/atomic"
	"time"
)

type (
	Task interface {
		IsPriority() bool
		// Init tries to perform the initial initialization of the service, the logic of the function must make sure
		// that all created connections to remote services are in working order and are pinging. Otherwise, the
		// application will need additional error handling.
		Init(ctx context.Context) error
		// Ping will be called by the service controller at regular intervals, it is important that a response with
		// any error will be regarded as an unrecoverable state of the service and will lead to an emergency stop of
		// the application. If the service is not critical for the application, like a memcached, then try to implement
		// the logic of self-diagnosis and service recovery inside Ping, and return the nil as a response even if the
		// recovery failed.
		Ping(ctx context.Context) error
		// Close will be executed when the service controller receives a stop command. Normally, this happens after the
		// main thread of the application has already finished. That is, no more requests from the outside are expected.
		Close() error
	}
	TaskKeeper struct {
		Tasks           []Task
		PingPeriod      time.Duration
		PingTimeout     time.Duration
		ShutdownTimeout time.Duration

		stop  chan struct{}
		state int32
	}
)

const (
	srvStateInit int32 = iota
	srvStateReady
	srvStateRunning
	srvStateShutdown
	srvStateOff

	defaultPingPeriod      = time.Second * 5
	defaultPingTimeout     = time.Millisecond * 1500
	defaultShutdownTimeout = time.Millisecond * 15000
)

// Init will initialize all registered services. Will return an error if at least one of the initialization functions
// returned an error. It is very important that after the first error, the context with which the initialization
// functions of all services are performed will be immediately canceled.
func (s *TaskKeeper) Init(ctx context.Context) error {
	if !s.checkState(srvStateInit, srvStateReady) {
		return ErrWrongState
	}
	if err := s.initAllServices(ctx); err != nil {
		return err
	}
	s.stop = make(chan struct{})
	if s.PingPeriod == 0 {
		s.PingPeriod = defaultPingPeriod
	}
	if s.PingTimeout == 0 {
		s.PingTimeout = defaultPingTimeout
	}
	if s.ShutdownTimeout == 0 {
		s.ShutdownTimeout = defaultShutdownTimeout
	}
	return nil
}

func (s *TaskKeeper) checkState(old, new int32) bool {
	return atomic.CompareAndSwapInt32(&s.state, old, new)
}

func (s *TaskKeeper) initAllServices(ctx context.Context) (initError error) {
	initCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	var p parallelRun
	priorityLength := len(lo.Filter(s.Tasks, func(item Task, index int) bool {
		return item.IsPriority()
	}))
	for i := range lo.Filter(s.Tasks, func(item Task, index int) bool {
		return item.IsPriority()
	}) {
		p.do(initCtx, s.Tasks[i].Init)
	}
	if p.wait() != nil {
		return p.err
	}
	var pp parallelRun
	for i := range lo.Filter(s.Tasks, func(item Task, index int) bool {
		return !item.IsPriority()
	}) {
		pp.do(initCtx, s.Tasks[i+priorityLength].Init)
	}
	return pp.wait()
}

func (s *TaskKeeper) testServices(ctx context.Context) error {
	var ctxPing, cancel = context.WithTimeout(ctx, s.PingTimeout)
	defer cancel()
	var p parallelRun
	for i := range s.Tasks {
		p.do(ctxPing, s.Tasks[i].Ping)
	}
	return p.wait()
}

// Watch monitors the health of resources. At a given frequency, all services will receive a Ping command,
// and if any of the responses contains an error, all execution will immediately stop and the error will be
// transmitted as a result of the Watch procedure.
//
// This procedure is synchronous, which means that control of the routine will be returned only when service monitoring is interrupted.
func (s *TaskKeeper) Watch(ctx context.Context) error {
	if !s.checkState(srvStateReady, srvStateRunning) {
		return ErrWrongState
	}
	if err := s.cycleTestServices(ctx); err != nil && err != ErrShutdown {
		return err
	}
	return nil
}

func (s *TaskKeeper) cycleTestServices(ctx context.Context) error {
	for {
		select {
		case <-s.stop:
			return nil
		case <-time.After(s.PingPeriod):
			if err := s.testServices(ctx); err != nil {
				return err
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

// Stop sends a signal that monitoring should be stopped. Stops execution of the Watch procedure
func (s *TaskKeeper) Stop() {
	if s.checkState(srvStateRunning, srvStateShutdown) {
		close(s.stop)
	}
}

func (s *TaskKeeper) Release() error {
	if s.checkState(srvStateShutdown, srvStateOff) {
		return s.release()
	}
	return ErrWrongState
}

func (s *TaskKeeper) release() error {
	shCtx, cancel := context.WithTimeout(context.Background(), s.ShutdownTimeout)
	defer cancel()
	var p parallelRun
	for i := range s.Tasks {
		var service = s.Tasks[i]
		p.do(shCtx, func(_ context.Context) error {
			return service.Close()
		})
	}
	var errCh = make(chan error)
	go func() {
		defer close(errCh)
		if err := p.wait(); err != nil {
			errCh <- err
		}
	}()
	for {
		select {
		case err, ok := <-errCh:
			if ok {
				return err
			}
			return nil
		case <-shCtx.Done():
			return shCtx.Err()
		}
	}
}
