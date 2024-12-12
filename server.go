package main

import (
	"fmt"
	"go-distributed-cache/cache"
	"log"
	"net"
	"strings"
)

type ServerOpts struct {
	ListenAddr string
	IsLeader   bool
}

type Server struct {
	ServerOpts
	cache cache.Cacher
}

func NewServer(opts ServerOpts, c cache.Cacher) *Server {
	return &Server{
		ServerOpts: opts,
		cache:      c,
	}

}

func (s *Server) start() error {
	ln, err := net.Listen("tcp", s.ListenAddr)
	if err != nil {
		return fmt.Errorf("listen err: %s", err)
	}
	fmt.Printf("server starting on port [%s]\n", s.ListenAddr)
	for {
		conn, err := ln.Accept()
		if err != nil {

			log.Printf("accept err: %s\n", err)
			continue
		}
		go s.handleConn(conn)

	}
}

func (s *Server) handleConn(conn net.Conn) {
	defer func() {
		conn.Close()
	}()
	buf := make([]byte, 2048)

	for {
		n, err := conn.Read(buf)
		if err != nil {
			log.Printf("conn reading err: %s\n", err)
		}
		msg := buf[:n]
		fmt.Println(string(msg))

		go s.handleCommand(conn, msg)
	}

}

func (s *Server) handleCommand(conn net.Conn, rawCmd []byte) {
	var (
		rawStr = string(rawCmd)
		parts  = strings.Split(rawStr, " ")
	)

	if len(parts) == 0 {
		//respond
		log.Println("empty command")
		return
	}

	cmd := parts[0]

	if cmd == CMDSet {
		if len(parts) != 4 {
			//respond
			log.Println("invalid SET command")
			return
		}
		msg := &MSGSet{
			Key:   []byte(parts[1]),
			Value: []byte(parts[2]),
			TTL:   parts[3],
		}
		if err := s.handleSetCmd(conn, msg)

}

func (s *Server) handleSetCmd(conn net.Conn) error {
	return nil
}
