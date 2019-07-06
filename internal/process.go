package internal

import (
	"errors"
	"io"
)

func Process(reader io.Reader, writer io.Writer, offset, limit int) error {
	if offset != 0 {
		seeker, ok := reader.(io.Seeker)
		if !ok {
			return errors.New("не могу выполнить смещение")
		}
		err := seeker.Seek(int64(offset), io.SeekStart)
		if err != nil {
			return err
		}
	}
	io.CopyN(writer, reader, int64(limit))
	return nil
}
