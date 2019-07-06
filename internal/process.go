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

	chunk := limit / 100
	if chunk == 0 {
		chunk = 1
	}

	for i := 0; i < limit; i += chunk {
		callback(i * 100 / limit)
		io.CopyN(writer, reader, int64(chunk))
	}
	callback(100)

	return nil
}
