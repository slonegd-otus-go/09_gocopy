package internal

import (
	"errors"
	"io"
)

func Process(reader io.Reader, writer io.Writer, offset, limit int, callback func(progress int)) error {
	if offset != 0 {
		seeker, ok := reader.(io.Seeker)
		if !ok {
			return errors.New("не могу выполнить смещение")
		}
		_, err := seeker.Seek(int64(offset), io.SeekStart)
		if err != nil {
			return err
		}
	}
	// TODO разделить limit(или размер) на 100
	// сделать цикл считывания по частям
	// в цикле вызывать callback
	callback(0)
	io.CopyN(writer, reader, int64(limit))
	return nil
}
