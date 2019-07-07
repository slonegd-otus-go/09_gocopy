package internal

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type notSeeker struct{}

func (notSeeker) Read([]byte) (int, error) {
	return 0, nil
}

type errorSeeker struct{}

func (errorSeeker) Read([]byte) (int, error) {
	return 0, nil
}
func (errorSeeker) Seek(int64, int) (int64, error) {
	return 0, errors.New("seek error")
}

type errorRead struct{}

func (errorRead) Read([]byte) (int, error) {
	return 0, errors.New("read error")
}
func (errorRead) Seek(int64, int) (int64, error) {
	return 0, nil
}

func TestProcess(t *testing.T) {
	var gotProgress string
	callback := func(progress int) {
		gotProgress += strconv.Itoa(progress) + " "
	}
	tests := []struct {
		name         string
		reader       io.Reader
		writer       io.Writer
		offset       int
		limit        int
		callback     func(progress int)
		wantWriter   string
		wantProgress string
		wantErr      string
	}{
		{
			name:         "happy path",
			reader:       strings.NewReader("test"),
			writer:       &bytes.Buffer{},
			limit:        4,
			callback:     callback,
			wantWriter:   "test",
			wantProgress: "0 1 2 3 4 ",
		},
		{
			name:         "happy path with offset",
			reader:       strings.NewReader("test"),
			writer:       &bytes.Buffer{},
			offset:       2,
			limit:        2,
			callback:     callback,
			wantWriter:   "st",
			wantProgress: "0 1 2 ",
		},
		{
			name:         "limit < lehgth",
			reader:       strings.NewReader("test"),
			writer:       &bytes.Buffer{},
			limit:        3,
			callback:     callback,
			wantWriter:   "tes",
			wantProgress: "0 1 2 3 ",
		},
		{
			name:         "limit > lehgth",
			reader:       strings.NewReader("test"),
			writer:       &bytes.Buffer{},
			limit:        5,
			callback:     callback,
			wantWriter:   "test",
			wantProgress: "0 1 2 3 4 4 ", // дополнительная 4, потому что файл меньше limit
		},
		{
			name:         "offset > lehgth",
			reader:       strings.NewReader("test"),
			writer:       &bytes.Buffer{},
			limit:        5,
			offset:       5,
			callback:     callback,
			wantWriter:   "",
			wantProgress: "0 0 ",
		},
		{
			name:         "no Seeker error",
			reader:       notSeeker{},
			writer:       &bytes.Buffer{},
			limit:        5,
			offset:       5,
			callback:     callback,
			wantWriter:   "",
			wantProgress: "",
			wantErr:      "не могу выполнить смещение",
		},
		{
			name:         "seeker error",
			reader:       errorSeeker{},
			writer:       &bytes.Buffer{},
			limit:        5,
			offset:       5,
			callback:     callback,
			wantWriter:   "",
			wantProgress: "",
			wantErr:      "seek error",
		},
		{
			name:         "read error",
			reader:       errorRead{},
			writer:       &bytes.Buffer{},
			limit:        5,
			offset:       5,
			callback:     callback,
			wantWriter:   "",
			wantProgress: "0 ",
			wantErr:      "ошибка копирования: read error",
		},
	}
	for _, tt := range tests {
		//nolint:scopelint
		t.Run(tt.name, func(t *testing.T) {
			gotProgress = ""
			err := Process(tt.reader, tt.writer, tt.offset, tt.limit, tt.callback)
			if len(tt.wantErr) != 0 {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr, err.Error())
			}
			assert.Equal(t, tt.wantWriter, tt.writer.(fmt.Stringer).String())
			assert.Equal(t, tt.wantProgress, gotProgress)
		})
	}
}
