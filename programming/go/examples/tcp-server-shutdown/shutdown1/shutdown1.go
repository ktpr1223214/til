package shutdown1

import (
	"io"
	"log"
	"net"
	"sync"
)

type Server struct {
	listener net.Listener
	quit     chan interface{}
	wg       sync.WaitGroup
}

func NewServer(addr string) *Server {
	s := &Server{
		quit: make(chan interface{}),
	}
	l, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}
	s.listener = l
	s.wg.Add(1)
	go s.serve()
	return s
}

func (s *Server) Stop() {
	// ここ 2 つで serve stopped accepting new clients
	close(s.quit)
	s.listener.Close()
	// s.wg.Wait 呼び出しで、Stop will block until all the handlers have returned
	// つまり、s.handleConection(conn) の終了を待つ
	s.wg.Wait()
}

func (s *Server) serve() {
	defer s.wg.Done()

	for {
		conn, err := s.listener.Accept()
		if err != nil {
			select {
			case <-s.quit:
				return
			default:
				log.Println("accept error", err)
			}
		} else {
			s.wg.Add(1)
			go func() {
				s.handleConection(conn)
				s.wg.Done()
			}()
		}
	}
}

func (s *Server) handleConection(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 2048)
	for {
		n, err := conn.Read(buf)
		if err != nil && err != io.EOF {
			log.Println("read error", err)
			return
		}
		if n == 0 {
			return
		}
		log.Printf("received from %v: %s", conn.RemoteAddr(), string(buf[:n]))
	}
}

func init() {
	log.SetFlags(log.Ltime | log.Lmicroseconds)
}
