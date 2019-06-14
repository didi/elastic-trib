package main

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
)

// Display use to output something on screen with table format.
type Display struct {
	w *tabwriter.Writer
}

// AddRow add a row of data.
func (d *Display) AddRow(row []string) {
	fmt.Fprintln(d.w, strings.Join(row, "\t"))
}

// Flush output all rows on screen.
func (d *Display) Flush() error {
	return d.w.Flush()
}

// NewTableDisplay creates a display instance, and uses to format output with table.
func NewTableDisplay() *Display {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	return &Display{w}
}
