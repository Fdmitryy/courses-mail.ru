package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
)

func main() {
	err := calc(os.Stdin, os.Stdout)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func calc(input io.Reader, output io.Writer) error {
	var sp = new(int)
	*sp = 0
	stack := make([]int, 1000)
	in := bufio.NewScanner(input)
	in.Split(bufio.ScanWords)
	for in.Scan() {
		c := in.Text()
		switch c {
		case "\n":
			fallthrough
		case " ":
		case "=":
			if _, err := fmt.Fprintf(output, "Result = %d\n", pop(sp, stack)); err != nil {
				return err
			}
		case "+":
			push(pop(sp, stack)+pop(sp, stack), sp, stack)
		case "-":
			push(-pop(sp, stack)+pop(sp, stack), sp, stack)
		case "*":
			push(pop(sp, stack)*pop(sp, stack), sp, stack)
		case "/":
			if a := pop(sp, stack); a == 0 {
				return fmt.Errorf("DIVISION BY ZERO")
			} else {
				push(pop(sp, stack)/a, sp, stack)
			}
		default:
			if x, err := strconv.Atoi(c); err != nil {
				return fmt.Errorf("CAN'T READ INTEGER")
			} else {
				push(x, sp, stack)
			}
		}
	}
	i := 0
	for *sp != 0 {
		if _, err := fmt.Fprintf(output, "Stack[%d] = %d\n", i, pop(sp, stack)); err != nil {
			return err
		}
		i++
	}
	return nil
}

func pop(sp *int, stack []int) int {
	if *sp > 0 {
		*sp--
		return stack[*sp]
	} else {
		fmt.Println("Impossible to do pop() for empty stack")
		return 0
	}
}

func push(a int, sp *int, stack []int) {
	stack[*sp] = a
	*sp++
}
