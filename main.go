package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

func main() {
	if err := run(os.Stdin, os.Stdout, os.Stderr); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(in io.Reader, out, errOut io.Writer) error {
	br := bufio.NewReader(in)

	for line := 1; ; line++ {
		text, err := br.ReadString('\n')
		if err != nil && err != io.EOF {
			return fmt.Errorf("read input: %w", err)
		}
		if text == "" && err == io.EOF {
			return nil
		}

		rect, parseErr := parseRectangle(text)
		if parseErr != nil {
			fmt.Fprintf(errOut, "line %d: %v\n", line, parseErr)
		} else {
			fmt.Fprintf(out, "%.2f %.2f\n", rect.area(), rect.perimeter())
		}

		if err == io.EOF {
			return nil
		}
	}
}

func parseRectangle(s string) (rectangle, error) {
	fields := strings.Fields(s)
	if len(fields) == 0 {
		return rectangle{}, errors.New("blank line")
	}
	if len(fields) != 2 {
		return rectangle{}, fmt.Errorf("got %d fields, want 2 (width height)", len(fields))
	}

	width, err := strconv.ParseFloat(fields[0], 64)
	if err != nil {
		return rectangle{}, fmt.Errorf("bad width %q: %w", fields[0], err)
	}

	height, err := strconv.ParseFloat(fields[1], 64)
	if err != nil {
		return rectangle{}, fmt.Errorf("bad height %q: %w", fields[1], err)
	}

	return newRectangle(width, height)
}
