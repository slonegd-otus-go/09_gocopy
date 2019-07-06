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
	tests := []struct {
		name         string
		reader       io.Reader
		offset       int
		limit        int
		callback     func(progress int)
		wantWriter   string
		wantProgress string
		wantErr      bool
	}{
		{
			name:   "happy path",
			reader: strings.NewReader("test"),
			limit:  4,
			callback: func(progress int) {
				gotProgress += strconv.Itoa(progress) + "% "
			},
			wantWriter:   "test",
			wantProgress: "0% 25% 50% 75% 100% ",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotProgress = ""
			writer := &bytes.Buffer{}
			Process(tt.reader, writer, tt.offset, tt.limit, tt.callback)
			assert.Equal(t, tt.wantWriter, writer.String())
			assert.Equal(t, tt.wantProgress, gotProgress)
		})
	}
}
