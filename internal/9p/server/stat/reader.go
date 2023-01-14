package stat // import "go.rbn.im/neinp/stat"

import (
	"bytes"
	"fmt"
	"io"
)

/*Reader is a reader for multiple Stats, aka directories.*/
type Reader struct {
	stats   []Stat
	pos     int
	bytepos int64
}

/*NewReader returns a new Reader initialized with stats and ready for reading.*/
func NewReader(stats ...Stat) *Reader {
	return &Reader{stats: stats}
}

func (s *Reader) Read(p []byte) (int, error) {
	var outbuf bytes.Buffer
	var buf bytes.Buffer
	npos := s.pos

	if s.pos == len(s.stats) || len(s.stats) == 0 {
		return 0, io.EOF
	}

	for _, stat := range s.stats[s.pos:] {
		_, err := stat.Encode(&buf)
		if err != nil {
			return 0, err
		}

		if buf.Len()+outbuf.Len() > len(p) {
			break
		}

		_, err = buf.WriteTo(&outbuf)
		if err != nil {
			return 0, err
		}

		npos++
	}

	n := copy(p, outbuf.Bytes())

	s.bytepos += int64(n)
	s.pos = npos

	return n, nil
}

//Seek interface implementation.
func (s *Reader) Seek(offset int64, whence int) (int64, error) {
	if whence != io.SeekStart {
		return s.bytepos, fmt.Errorf("Unsupported whence %v", whence)
	}

	switch offset {
	case 0:
		s.pos = 0
		s.bytepos = 0
	case s.bytepos:
	default:
		return s.bytepos, fmt.Errorf("can't seek there")
	}

	return s.bytepos, nil
}
