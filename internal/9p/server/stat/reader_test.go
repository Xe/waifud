package stat // import "go.rbn.im/neinp/stat"

import (
	"bytes"
	"io"
	"testing"
	"time"
)

var stats = []Stat{
	Stat{Name: "a", Atime: time.Unix(1234567, 0), Mtime: time.Unix(1234567, 0)},
	Stat{Name: "b", Atime: time.Unix(1234567, 0), Mtime: time.Unix(1234567, 0)},
}

func TestReaderRead(t *testing.T) {
	r := NewReader(stats...)

	var buf bytes.Buffer
	buf.ReadFrom(r)

	for _, oldstat := range stats {
		var newstat Stat
		newstat.Decode(&buf)

		if newstat != oldstat {
			t.Fail()
		}
	}
}

func TestReaderSeek(t *testing.T) {
	r := NewReader(stats...)

	// We need to read into a buffer because stat.decode does multiple reads,
	// and StatReader only supports reading a whole stat at once.
	var buf bytes.Buffer
	var a0 Stat
	var a1 Stat

	// Read the encoded stats to the buffer
	buf.ReadFrom(r)

	_, err := a0.Decode(&buf)
	if err != nil {
		t.Error(err)
	}

	// Seek to start.
	_, err = r.Seek(0, io.SeekStart)
	if err != nil {
		t.Error(err)
	}

	// Reset the buffer, then read again.
	// The read now should contain the first stat again.
	buf.Reset()
	buf.ReadFrom(r)

	_, err = a1.Decode(&buf)
	if err != nil {
		t.Error(err)
	}

	if a0 != a1 {
		t.Logf("\n%v\n%v", a0, a1)
		t.Fail()
	}

}
