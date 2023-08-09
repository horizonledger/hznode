package main

import (
	"fmt"
	"log"
	"stack"
)

const (
	OP_TRANSFER byte = 0x00
	OP_PUSH     byte = 0x10
	OP_EQUALVERIFY byte = 0x11
	OP_CHECKSIGCALL byte = 0x12
	OP_KHASH       byte = 0x14
	OP_SHA256     byte = 0x15
	OP_NMREG      byte = 0x20
	OP_NMLOAD     byte = 0x21
)

type State struct {
	bal  map[string]int
	names map[string]string
}

type Runtime struct {
}

func (r *Runtime) evaltx(state State, stack *stack.Stack, caller string, sscript []byte, exescript []byte) {
	fmt.Printf("State: %+v\n", state)

	valid := r.evalscript(state, stack, caller, sscript)
	if valid {
		log.Println("execute " + string(exescript))
		r.evalscript(state, stack, caller, exescript)
	}
}

func (r *Runtime) evalscript(state State, stack *stack.Stack, caller string, sscript []byte) bool {
	log.Println("eval " + string(sscript))
	log.Println("caller " + caller)

	ptr := 0

	for ptr < len(sscript) {
		operation := sscript[ptr]
		opcode := operation
		ptr++

		switch opcode {
		case OP_TRANSFER:
			// TODO: Implement logic
		case OP_NMREG:
			// TODO: Implement logic
		case OP_NMLOAD:
			// TODO: Implement logic
		case OP_SHA256:
			// TODO: Implement logic
		case OP_PUSH:
			// TODO: Implement logic
		case OP_EQUALVERIFY:
			// TODO: Implement logic
		case OP_CHECKSIGCALL:
			// TODO: Implement logic
		case OP_KHASH:
			// TODO: Implement logic
		default:
			panic("Unrecognized opcode: " + string(opcode))
		}
	}

	log.Println("stack length at the end: ", stack.Len())
	log.Println("stack at the end: ", stack)

	// You will need to implement your own logic to determine whether the script is valid
	isValid := true

	return isValid
}

func main() {
	r := &Runtime{}
	stack := stack.New()

	state := State{
		bal: map[string]int{
			"user1": 100,
			"user2": 200,
		},
		names: map[string]string{
			"name1": "user1",
			"name2": "user2",
		},
	}

	sscript := []byte{OP_PUSH, OP_SHA256}
	exescript := []byte{OP_PUSH, OP_SHA256}

	r.evaltx(state, stack, "user1", sscript, exescript)
}
