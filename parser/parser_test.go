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
        if res.args[0] != "SIMPLESTRING" {
            t.Errorf("expected: %v, got: %v\n", expect, res.args[0])
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
        if res.args[0] != "SIMPLESTRING" {
            t.Errorf("expected: %v, got: %v\n", expect, res.args[0])
        }
	})

	// t.Run("Whole message delivered in single slice", func(t *testing.T) {
	// 	expect := []string{"hello", "world"}
	// 	in := make(chan []byte)
	// 	go func() {
	// 		in <- []byte("*2\r\n$5\r\nhello\r\n$5\r\nworld\r\n")
	// 		close(in)
	// 	}()
	// 	out := Parse(in)
	// 	res := <-out
        // for i, v := range expect {
            // if v != res[i] {
                // t.Errorf("expected: %v, got: %v", v, res[i])
            // }
        // }
	// })

	// t.Run("Message partitioned into multiple slices", func(t *testing.T) {
	// 	expect := []string{"hello", "world"}
	// 	in := make(chan []byte)
	// 	go func() {
	// 		in <- []byte("*2\r\n$5\r\nh")
	// 		in <- []byte("ello\r\n$5\r\nworld\r\n")
	// 		close(in)
	// 	}()
	// 	out := Parse(in)
	// 	res := <-out
        // for i, v := range expect {
            // if v != res[i] {
                // t.Errorf("expected: %v, got: %v", v, res[i])
            // }
        // }
	// })

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
        // for i, v := range expect {
            // if v != res[i] {
                // t.Errorf("expected: %v, got: %v", v, res[i])
            // }
        // }
	// })
}
