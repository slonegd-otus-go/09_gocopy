package internal

import (
	"bytes"
	"io"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProcess(t *testing.T) {
	var gotProgress string
	callback := func(progress int) {
		gotProgress += strconv.Itoa(progress) + "% "
	}
	tests := []struct {
		name         string
		reader       io.Reader
		offset       int
		limit        int
		callback     func(progress int)
		wantWriter   string
		wantProgress string
		// wantErr      bool
	}{
		{
			name:   "happy path",
			reader: strings.NewReader("test"),
			limit:  4,
			callback: callback,
			wantWriter:   "test",
			wantProgress: "0% 25% 50% 75% 100% ",
		},
		{
			name:   "happy path with offset",
			reader: strings.NewReader("test"),
			offset: 2,
			limit:  2,
			callback: callback,
			wantWriter:   "st",
			wantProgress: "0% 50% 100% ",
		},
		{
			name:   "limit < lehgth",
			reader: strings.NewReader("test"),
			limit:  3,
			callback: callback,
			wantWriter:   "tes",
			wantProgress: "0% 33% 66% 100% ",
		},
		{
			name:   "limit > lehgth",
			reader: strings.NewReader("test"),
			limit:  5,
			callback: callback,
			wantWriter:   "test",
			wantProgress: "0% 20% 40% 60% 80% 100% ",
		},
		{
			name:   "offset > lehgth",
			reader: strings.NewReader("test"),
			limit:  5,
			offset: 5,
			callback: callback,
			wantWriter:   "",
			wantProgress: "0% 100% ",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotProgress = ""
			writer := &bytes.Buffer{}
			//nolint
			Process(tt.reader, writer, tt.offset, tt.limit, tt.callback)
			//nolint
			assert.Equal(t, tt.wantWriter, writer.String())
			//nolint
			assert.Equal(t, tt.wantProgress, gotProgress)
		})
	}
}
