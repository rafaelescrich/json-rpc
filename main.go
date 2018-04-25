// package main

// import (
// 	"fmt"
// 	"log"
// 	"net/rpc"
// 	"net/rpc/jsonrpc"
// 	"time"

// 	"github.com/rafaelescrich/json-rpc/services"
// 	kcp "github.com/xtaci/kcp-go"
// )

// func startServer() {
// 	arith := new(services.Arith)

// 	server := rpc.NewServer()
// 	server.Register(arith)

// 	server.HandleHTTP(rpc.DefaultRPCPath, rpc.DefaultDebugPath)

// 	// l, e := net.Listen("tcp", ":8080")
// 	l, e := kcp.Listen(":8080")

// 	if e != nil {
// 		log.Fatal("listen error:", e)
// 	}

// 	for {
// 		conn, err := l.Accept()
// 		if err != nil {
// 			log.Fatal(err)
// 		}

// 		go server.ServeCodec(jsonrpc.NewServerCodec(conn))
// 	}
// }

// func main() {
// 	go startServer()
// 	start := time.Now()

// 	type results struct {
// 		args  *services.Args
// 		reply uint64
// 	}

// 	for i := 0; i < 4096; i++ {

// 		ss := make(chan results)

// 		go func(i int) {
// 			// conn, err := net.Dial("tcp", "localhost:8080")
// 			conn, err := kcp.Dial("localhost:8080")
// 			defer conn.Close()

// 			if err != nil {
// 				panic(err)
// 			}
// 			c := jsonrpc.NewClient(conn)

// 			args := &services.Args{A: 7 * uint64(i), B: 8 * uint64(i)}
// 			var reply uint64
// 			err = c.Call("Arith.Multiply", args, &reply)
// 			if err != nil {
// 				log.Fatal("arith error:", err)
// 			}
// 			ss <- results{args, reply}
// 		}(i)

// 		rr := <-ss
// 		fmt.Printf("Arith: %d*%d=%d\n", rr.args.A, rr.args.B, rr.reply)
// 	}

// 	t := time.Now()
// 	elapsed := t.Sub(start)
// 	fmt.Printf("Time elapsed: %s\n", elapsed.String())
// }
package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"runtime"
	"time"
	// kcp "github.com/xtaci/kcp-go"
)

// Args is a struct to
type Args struct {
	A, B uint64
}

// Quotient struct
type Quotient struct {
	Quo, Rem uint64
}

// Arith int64 type
type Arith uint64

// Multiply two numbers
func (t *Arith) Multiply(args *Args, reply *uint64) error {
	*reply = args.A * args.B
	return nil
}

// Divide two numbers
func (t *Arith) Divide(args *Args, quo *Quotient) error {
	if args.B == 0 {
		return errors.New("divide by zero")
	}
	quo.Quo = args.A / args.B
	quo.Rem = args.A % args.B
	return nil
}

func startServer() {
	arith := new(Arith)

	server := rpc.NewServer()
	server.Register(arith)

	server.HandleHTTP(rpc.DefaultRPCPath, rpc.DefaultDebugPath)

	l, e := net.Listen("tcp", ":8080")
	// l, e := kcp.Listen(":8080")

	if e != nil {
		log.Fatal("listen error:", e)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}

		go server.ServeCodec(jsonrpc.NewServerCodec(conn))
	}
}

func main() {
	go startServer()
	start := time.Now()

	type results struct {
		args  *Args
		reply uint64
	}

	for i := 0; i < 8000; i++ {

		ss := make(chan results)

		go func(i int) {
			conn, err := net.Dial("tcp", "localhost:8080")
			// conn, err := kcp.Dial("localhost:8080")
			defer conn.Close()

			if err != nil {
				panic(err)
			}
			c := jsonrpc.NewClient(conn)

			args := &Args{A: 7 * uint64(i), B: 8 * uint64(i)}
			var reply uint64
			err = c.Call("Arith.Multiply", args, &reply)
			if err != nil {
				log.Fatal("arith error:", err)
			}
			ss <- results{args, reply}
		}(i)

		rr := <-ss
		fmt.Printf("Arith: %d*%d=%d\n", rr.args.A, rr.args.B, rr.reply)
	}

	t := time.Now()
	elapsed := t.Sub(start)
	fmt.Printf("Time elapsed: %s\n", elapsed.String())
	runtime.GC()
}
