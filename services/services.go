package services

import "errors"

// Args is a struct to
type Args struct {
	A, B int
}

// Quotient struct
type Quotient struct {
	Quo, Rem int
}

// Arith int type
type Arith int

// Multiply two numbers
func (t *Arith) Multiply(args *Args, reply *int) error {
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
