package main

import (
	"bytes"
	"fmt"
	"testing"
)

var testOkInput = []string{`2 3 + =`, `4 2 / =`, `5 4 * =`, `1 5 - =`, `2 3 +`}

var testOkResult = []string{
	`Result = 5
`, `Result = 2
`, `Result = 20
`, `Result = -4
`, `Stack[0] = 5
`,
}

var testFailInput = []string{`2 0 /`, `asd`}

func TestOK(t *testing.T) {
	for i, val := range testOkInput {
		in := bytes.NewBufferString(val)
		out := bytes.NewBuffer(nil)
		err := calc(in, out)
		if err != nil {
			t.Errorf("Test OK failed: %s", err)
		}
		result := out.String()
		if result != testOkResult[i] {
			t.Errorf("Test OK failed, result not match")
			fmt.Println(result)
		}
	}
}

func TestFail(t *testing.T) {
	for _, val := range testFailInput {
		in := bytes.NewBufferString(val)
		out := bytes.NewBuffer(nil)
		err := calc(in, out)
		if err == nil {
			t.Errorf("Test FAIL failed: expected error")
		}
	}
}
