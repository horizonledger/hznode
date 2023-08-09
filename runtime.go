package main

import (
	"errors"
	"fmt"
	"log"
	"stack"
)

const (
	OP_TRANSFER     byte = 0x00
	OP_PUSH         byte = 0x10
	OP_EQUALVERIFY  byte = 0x11
	OP_CHECKSIGCALL byte = 0x12
	OP_KHASH        byte = 0x14
	OP_SHA256       byte = 0x15
	OP_NMREG        byte = 0x20
	OP_NMLOAD       byte = 0x21
)

type State struct {
	bal   map[string]int
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
			amount := stack[len(stack)-1].(int)
			stack = stack[:len(stack)-1]
			to := stack[len(stack)-1].(string)
			stack = stack[:len(stack)-1]
			log.Println("transfer ", to, amount, caller)
			if amount > state.bal[caller] {
				log.Println("out of funds")
			} else {
				log.Println("enough funds. transfer")
				x := state.bal[caller]
				state.bal[caller] = x - amount
				y := state.bal[to]
				state.bal[to] = y + amount
			}

		case OP_NMREG:
			// Ensure there's at least one value on the stack.
			if len(stack) < 1 {
				panic(errors.New("Not enough items on stack for OP_NMREG"))
			}

			// Pop the name from the stack.
			name := stack[len(stack)-1].(string) // assuming the name is a string; adjust type accordingly
			stack = stack[:len(stack)-1]

			// Register the name against the caller.
			if _, exists := state.names[name]; exists {
				panic(errors.New("Name already registered"))
			}

			// Optionally, you can check registration costs or any other business logic here

			state.names[name] = caller

			log.Printf("register %s %s", name, caller)

		case OP_NMLOAD:
			// Ensure there's at least one value on the stack.
			if len(stack) < 1 {
				panic(errors.New("Not enough items on stack for OP_NMLOAD"))
			}

			// Pop the name from the stack.
			name := stack[len(stack)-1].(string) // assuming the name is a string; adjust type accordingly
			stack = stack[:len(stack)-1]

			// Load the owner of the name from the state.
			owner, exists := state.names[name]
			if !exists {
				// If the name isn't registered, push a nil or an empty string.
				// Adjust this behavior based on your application's requirements.
				stack = append(stack, "")
			} else {
				stack = append(stack, owner)
			}

			log.Printf("load %s -> %s", name, owner)

		case OP_SHA256:
			// TODO: Implement logic
		case OP_PUSH:
			if ptr >= len(sscript) {
				panic(errors.New("OP_PUSH without a subsequent value"))
			}

			value := sscript[ptr]
			stack = append(stack, value)
			ptr++ // Increment the pointer to move past the value.

		case OP_EQUALVERIFY:
			if len(stack) < 2 {
				panic(errors.New("Not enough items on stack for OP_EQUALVERIFY"))
			}

			// Pop the top two values from the stack.
			value1 := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			value2 := stack[len(stack)-1]
			stack = stack[:len(stack)-1]

			// Check if the two values are equal.
			if value1 != value2 {
				panic(errors.New("OP_EQUALVERIFY failed: top two stack values are not equal"))
			}

			// Push 'true' onto the stack to indicate successful verification.
			stack = append(stack, true)

		case OP_CHECKSIGCALL:
			// Ensure there are at least two values on the stack.
			if len(stack) < 2 {
				panic(errors.New("Not enough items on stack for OP_CHECKSIGCALL"))
			}

			// Pop the message hash and DER encoded signature from the stack.
			derSign := stack[len(stack)-1].(string) // assuming the values are strings; adjust types accordingly
			stack = stack[:len(stack)-1]
			msgHash := stack[len(stack)-1].(string)
			stack = stack[:len(stack)-1]

			publicKey := caller

			// Perform signature verification (assuming a function like checkSig exists in Go).
			isValid, err := checkSig(publicKey, msgHash, derSign)
			if err != nil {
				// If there's an error in signature verification, push false onto the stack.
				stack = append(stack, false)
			} else {
				// Push the result of signature verification onto the stack.
				stack = append(stack, isValid)
			}

		case OP_KHASH:
			// Ensure there's at least one value on the stack.
			if len(stack) < 1 {
				panic(errors.New("Not enough items on stack for OP_KHASH"))
			}

			// Pop the top value from the stack for hashing.
			value := stack[len(stack)-1].(string) // assuming the value is a string; adjust type accordingly
			stack = stack[:len(stack)-1]

			// Compute the Keccak hash (assuming you have a function named keccak for this).
			hash := keccak(value)

			// Push the hash onto the stack.
			stack = append(stack, hash)
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
