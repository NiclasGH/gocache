package persistence

import (
	"bufio"
	"gocache/internal/resp"
	"io"
	"os"
	"sync"
	"time"
)

type Aof struct {
	file   *os.File
	reader *bufio.Reader
	mutex  sync.Mutex
}

func NewAof(path string) (*Aof, error) {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return nil, err
	}

	aof := &Aof{
		file:   file,
		reader: bufio.NewReader(file),
	}

	// ensuring data integrity, even if the program crashes
	go func() {
		for {
			aof.mutex.Lock()
			aof.file.Sync()
			aof.mutex.Unlock()

			time.Sleep(time.Second)
		}
	}()

	return aof, nil
}

func (aof *Aof) Initialize(fn func(resp.Value)) error {
	aof.mutex.Lock()
	defer aof.mutex.Unlock()

	// move current file buffer to start
	aof.file.Seek(0, io.SeekStart)

	reader := resp.NewReader(aof.file)

	for {
		value, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}

			return err
		}

		fn(value)
	}

	return nil
}

func (aof *Aof) Save(value resp.Value) error {
	bytes := value.Marshal()

	aof.mutex.Lock()
	defer aof.mutex.Unlock()

	if _, err := aof.file.Write(bytes); err != nil {
		return err
	}

	return nil
}

func (aof *Aof) Close() error {
	aof.mutex.Lock()
	defer aof.mutex.Unlock()

	return aof.file.Close()
}
