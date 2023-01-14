package neinp

import (
	"context"
	"fmt"
	"io"
	"log"
	"sync"
	"sync/atomic"

	"github.com/Xe/waifud/internal/9p/server/message"
)

// P2000 is the interface which file systems must implement to be used with Server.
// To ease implementation of new filesystems, NopP2000 can be embedded. See the
// types defined in message for documentation of what they do.
type P2000 interface {
	Version(context.Context, message.TVersion) (message.RVersion, error)
	Auth(context.Context, message.TAuth) (message.RAuth, error)
	Attach(context.Context, message.TAttach) (message.RAttach, error)
	Walk(context.Context, message.TWalk) (message.RWalk, error)
	Open(context.Context, message.TOpen) (message.ROpen, error)
	Create(context.Context, message.TCreate) (message.RCreate, error)
	Read(context.Context, message.TRead) (message.RRead, error)
	Write(context.Context, message.TWrite) (message.RWrite, error)
	Clunk(context.Context, message.TClunk) (message.RClunk, error)
	Remove(context.Context, message.TRemove) (message.RRemove, error)
	Stat(context.Context, message.TStat) (message.RStat, error)
	Wstat(context.Context, message.TWstat) (message.RWstat, error)
	Close() error
}

// Server muxes and demuxes 9p messages from a connection,
// calling the corresponding handlers of a P2000 interface.
type Server struct {
	fs    P2000
	msize uint32

	tags  map[uint16]context.CancelFunc
	tagsm sync.Mutex

	done chan struct{}

	Debug bool
}

// NewServer returns a Server initialized to use fs for message handling.
func NewServer(fs P2000) *Server {
	return &Server{
		fs:    fs,
		msize: 4096, // initial msize, shoud suffice to read the version message
		tags:  make(map[uint16]context.CancelFunc),
		done:  make(chan struct{}),
	}
}

// rcv reads from an io.Reader and sends the decoded messages to a channel.
func (s *Server) rcv(done <-chan struct{}, r io.Reader) (<-chan message.Message, <-chan error) {
	in := make(chan message.Message)
	rcverr := make(chan error)

	go func() {
		defer close(in)
		defer close(rcverr)
		for {
			select {
			case <-done:
				return
			default:
				req := message.Message{}
				_, err := req.Decode(io.LimitReader(r, int64(atomic.LoadUint32(&s.msize))))
				if err != nil {
					rcverr <- err
					return
				}

				if s.Debug {
					log.Printf("<- %#v", req.Content)
				}

				in <- req
			}
		}
	}()

	return in, rcverr
}

// process receives messages from a channel and calls the corresponding handler of the
// implementing P2000 interface.
//
// process is supposed to be called as goroutine from handle. the methods of the implementing
// P2000 interface are called together with a context usable for cancelation, so that they can be canceled by
// a later flush message. NB: it may be better to create the context in the caller instead
// of creating it here, but as they are put in a mutexed map here this _should_ work.
func (s *Server) process(req message.Message) (message.Message, error) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	s.tagsm.Lock()
	s.tags[req.Tag] = cancel
	s.tagsm.Unlock()

	defer func() {
		s.tagsm.Lock()
		delete(s.tags, req.Tag)
		s.tagsm.Unlock()
	}()

	var c message.Content
	var err error

	switch t := req.Content.(type) {
	case *message.TVersion:
		// set msize to what the implementing file system wants
		var rVersion message.RVersion
		rVersion, err = s.fs.Version(ctx, *t)
		atomic.StoreUint32(&s.msize, rVersion.Msize)
		c = &rVersion
	case *message.TAuth:
		var rAuth message.RAuth
		rAuth, err = s.fs.Auth(ctx, *t)
		c = &rAuth
	case *message.TAttach:
		var rAttach message.RAttach
		rAttach, err = s.fs.Attach(ctx, *t)
		c = &rAttach
	case *message.TFlush:
		s.tagsm.Lock()
		cancel, ok := s.tags[t.Oldtag]
		s.tagsm.Unlock()
		if ok {
			cancel()
		}
		c = &message.RFlush{}
	case *message.TWalk:
		var rWalk message.RWalk
		rWalk, err = s.fs.Walk(ctx, *t)
		c = &rWalk
	case *message.TOpen:
		var rOpen message.ROpen
		rOpen, err = s.fs.Open(ctx, *t)
		c = &rOpen
	case *message.TCreate:
		var rCreate message.RCreate
		rCreate, err = s.fs.Create(ctx, *t)
		c = &rCreate
	case *message.TRead:
		var rRead message.RRead
		rRead, err = s.fs.Read(ctx, *t)
		c = &rRead
	case *message.TWrite:
		var rWrite message.RWrite
		rWrite, err = s.fs.Write(ctx, *t)
		c = &rWrite
	case *message.TClunk:
		var rClunk message.RClunk
		rClunk, err = s.fs.Clunk(ctx, *t)
		c = &rClunk
	case *message.TRemove:
		var rRemove message.RRemove
		rRemove, err = s.fs.Remove(ctx, *t)
		c = &rRemove
	case *message.TStat:
		var rStat message.RStat
		rStat, err = s.fs.Stat(ctx, *t)
		c = &rStat
	case *message.TWstat:
		var rWstat message.RWstat
		rWstat, err = s.fs.Wstat(ctx, *t)
		c = &rWstat
	default:
		// Return empty message and error for unexpected types. The empty message will
		// trigger the error being send to the handleErr channel.
		return message.Message{}, fmt.Errorf("unexpected message type %T", req.Content)
	}

	res := req.Response()
	res.Content = c

	// put errors from fs handers into the message
	if err != nil {
		res.Content = &message.RError{Ename: err.Error()}
	}

	return res, nil
}

// handle creates a process-goroutine per incoming message.
func (s *Server) handle(done <-chan struct{}, in <-chan message.Message) (<-chan <-chan message.Message, <-chan <-chan error) {
	out := make(chan (<-chan message.Message))
	handleErr := make(chan (<-chan error))

	go func() {
		defer close(out)
		for {
			select {
			case <-done:
				return
			case req := <-in:
				processRes := make(chan message.Message)
				processErr := make(chan error)

				go func() {
					defer close(processRes)
					defer close(processErr)

					res, err := s.process(req)
					processRes <- res
					if err != nil {
						processErr <- err
					}
				}()

				out <- processRes
				handleErr <- processErr
			}
		}
	}()

	return out, handleErr
}

// send sends the replys to the client.
func (s *Server) send(done <-chan struct{}, out <-chan (<-chan message.Message), w io.Writer) <-chan error {
	senderr := make(chan error)
	sendmtx := new(sync.Mutex)

	go func() {
		for {
			select {
			case <-done:
				return
			case x := <-out:
				go func(ch <-chan message.Message) {
					for msg := range ch {
						sendmtx.Lock()
						defer sendmtx.Unlock()

						if s.Debug {
							log.Printf("-> %#v", msg.Content)
						}

						_, err := msg.Encode(w)
						if err != nil {
							senderr <- err
							return
						}
					}
				}(x)
			}
		}
	}()

	return senderr
}

// Serve the filesystem on the connection given by rw.
func (s *Server) Serve(rw io.ReadWriter) error {
	done := make(chan struct{})

	in, rcvErr := s.rcv(done, rw)
	out, handleErr := s.handle(done, in)
	sendErr := s.send(done, out, rw)

	defer s.cleanup()

	procErr := make(chan error)
	for {
		select {
		case err := <-rcvErr:
			close(done)
			return err
		case processErr := <-handleErr:
			// the channel may be closed or error nil
			go func() {
				err := <-processErr
				if err != nil {
					procErr <- err
				}
			}()
		case err := <-procErr:
			close(done)
			return err
		case err := <-sendErr:
			close(done)
			return err
		case <-s.done:
			close(done)
			return nil
		}
	}
}

// Shutdown serving the filesystem.
func (s *Server) Shutdown() error {
	select {
	case <-s.done:
		return fmt.Errorf("shutdown but not serving")
	default:
		close(s.done)
	}
	return nil
}

// close the implementing fs
func (s *Server) cleanup() {
	if s.Debug {
		log.Println("cleanup: start")
	}
	err := s.fs.Close()
	if err != nil {
		log.Println("cleanup:", err)
	}
}
