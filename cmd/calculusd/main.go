package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/afrolovskiy/calculus/Godeps/_workspace/src/github.com/apache/thrift/lib/go/thrift"

	"github.com/afrolovskiy/calculus/calculus"
	"github.com/afrolovskiy/calculus/calculus/proto"
)

var (
	addr = flag.String("addr", "localhost:9090", "<addr>:<port> of the server")
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) Calculate(expr string) (retval2 float64, err error) {
	rpn, err := calculus.NewRPN(expr)
	if err != nil {
		return 0, err
	}

	return rpn.Calculate()
}

func GetServer(addr string) (*thrift.TSimpleServer, error) {
	st, err := thrift.NewTServerSocket(addr)
	if err != nil {
		return nil, err
	}

	pr := proto.NewCalculusProcessor(NewHandler())
	pf := thrift.NewTBinaryProtocolFactoryDefault()
	tf := thrift.NewTTransportFactory()
	return thrift.NewTSimpleServer4(pr, st, tf, pf), nil
}

func main() {
	flag.Parse()

	server, err := GetServer(*addr)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	log.Printf("Starting server on %s", *addr)
	if err := server.Serve(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
