package main

import (
	"encoding/csv"
	"io"
	"os"
)

type csvHelper struct {
	f      io.Closer
	writer *csv.Writer
}

func newCSVHelper(filepath string) *csvHelper {
	f, err := os.Create(filepath)
	if err != nil {
		panic(err)
	}
	writer := csv.NewWriter(f)
	return &csvHelper{f, writer}
}

func (h *csvHelper) write(record []string) error {
	defer h.writer.Flush()
	return h.writer.Write(record)
}

func (h *csvHelper) close() {
	defer h.f.Close()

	h.writer.Flush()

	err := h.writer.Error()
	if err != nil {
		panic(err)
	}
}
