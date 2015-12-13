package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"

	"github.com/afrolovskiy/calculus/Godeps/_workspace/src/github.com/apache/thrift/lib/go/thrift"

	"github.com/afrolovskiy/calculus/calculus/proto"
)

var (
	addr = flag.String("addr", "localhost:9090", "<addr>:<port> of the server")
)

func GetClient(addr string) (*proto.CalculusClient, error) {
	ts, err := thrift.NewTSocket(addr)
	if err != nil {
		return nil, err
	}

	tf := thrift.NewTTransportFactory()
	tr := tf.GetTransport(ts)
	if err := tr.Open(); err != nil {
		return nil, err
	}

	pf := thrift.NewTBinaryProtocolFactoryDefault()
	return proto.NewCalculusClientFactory(tr, pf), nil
}

func main() {
	flag.Parse()

	client, err := GetClient(*addr)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	reader := bufio.NewReader(os.Stdin)
	expr, err := reader.ReadString('\n')
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	res, err := client.Calculate(expr)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Println(res)
}
