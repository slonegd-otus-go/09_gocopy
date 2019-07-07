package internal

import (
	"errors"
	"fmt"
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

	progress := 0
	for progress < limit {
		callback(progress)
		written, err := io.CopyN(writer, reader, int64(chunk))
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("ошибка копирования: %v", err)
		}
		progress += int(written)
	}
	callback(progress)

	return nil
}
