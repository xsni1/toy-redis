package parser

import (
	"testing"
)

// TODO: table tests could be used here
func TestParse(t *testing.T) {
	t.Run("Simple string - whole in single packet", func(t *testing.T) {
		expect := "SIMPLESTRING"
		in := make(chan []byte)
		go func() {
			in <- []byte("+SIMPLESTRING\r\n")
			close(in)
		}()
		out := Parse(in)
		res := <-out
		if res.Args[0] != "SIMPLESTRING" {
			t.Errorf("expected: %v, got: %v\n", expect, res.Args[0])
		}
	})
	t.Run("Simple string - segmented into multiple packets", func(t *testing.T) {
		expect := "SIMPLESTRING"
		in := make(chan []byte)
		go func() {
			in <- []byte("+SIMPLEST")
			in <- []byte("RING\r\n")
			close(in)
		}()
		out := Parse(in)
		res := <-out
		if res.Args[0] != "SIMPLESTRING" {
			t.Errorf("expected: %v, got: %v\n", expect, res.Args[0])
		}
	})
	t.Run("Bulk string - whole", func(t *testing.T) {
		expect := "SIMPLESTRING"
		in := make(chan []byte)
		go func() {
			in <- []byte("$12\r\nSIMPLESTRING\r\n")
			close(in)
		}()
		out := Parse(in)
		res := <-out
		if res.Args[0] != "SIMPLESTRING" {
			t.Errorf("expected: %v, got: %v\n", expect, res.Args[0])
		}
	})

	t.Run("Bulk string - segmented", func(t *testing.T) {
		expect := "SIMPLESTRING"
		in := make(chan []byte)
		go func() {
			in <- []byte("$12\r\nSIMPLEST")
			in <- []byte("RING\r\n")
			close(in)
		}()
		out := Parse(in)
		res := <-out
		if res.Args[0] != "SIMPLESTRING" {
			t.Errorf("expected: %v, got: %v\n", expect, res.Args[0])
		}
	})

	t.Run("Array - whole", func(t *testing.T) {
		expect := []string{"hello", "worl"}
		in := make(chan []byte)
		go func() {
			in <- []byte("*2\r\n$5\r\nhello\r\n$4\r\nworl\r\n")
			close(in)
		}()
		out := Parse(in)
		res := <-out
		for i, v := range expect {
			if v != res.Args[i] {
				t.Errorf("expected: %v, got: %v\n", v, res.Args[i])
			}
		}
	})

	t.Run("Whole message delivered in single slice", func(t *testing.T) {
		expect := []string{"hello", "world"}
		in := make(chan []byte)
		go func() {
			in <- []byte("*2\r\n$5\r\nhello\r\n$5\r\nworld\r\n")
			close(in)
		}()
		out := Parse(in)
		res := <-out
		for i, v := range expect {
			if v != res.Args[i] {
				t.Errorf("expected: %v, got: %v", v, res.Args[i])
			}
		}
	})

	// t.Run("Message has invalid eol", func(t *testing.T) {
	// 	expect := []string{"hello", "world"}
	// 	in := make(chan []byte)
	// 	go func() {
	// 		in <- []byte("*2\r\n$5\r\nh")
	// 		in <- []byte("ello\r\n$5\r\nworld\r")
	// 		close(in)
	// 	}()
	// 	out := Parse(in)
	// 	res := <-out
	// 	for i, v := range expect {
	// 		if v != res.args[i] {
	// 			t.Errorf("expected: %v, got: %v", v, res.args[i])
	// 		}
	// 	}
	// })
}
