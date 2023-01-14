package client

import (
	"errors"
	"io"
	"log"
)

// size[4] Rauth tag[2] aqid[13]
func readRauth(r io.Reader) (aqid QID, err error) {
	var size uint32
	if err = readUint32(r, &size); err != nil {
		return
	}
	var msgType uint8
	if err = readUint8(r, &msgType); err != nil {
		return
	}
	var tag uint16
	if err = readUint16(r, &tag); err != nil {
		return
	}
	if msgType == Rerror {
		var errmsg string
		if err = readString(r, &errmsg); err != nil {
			return
		}
		err = errors.New(errmsg)
		return
	}
	if msgType != Rauth {
		err = errUnexpectedMsg
		return
	}
	if err = readQID(r, &aqid); err != nil {
		return
	}
	if *debugLog {
		log.Println("->", "Rauth", "tag", tag, "aqid", aqid)
	}
	return
}

// size[4] Rattach tag[2] qid[13]
func readRattach(r io.Reader) (qid QID, err error) {
	var size uint32
	if err = readUint32(r, &size); err != nil {
		return
	}
	var msgType uint8
	if err = readUint8(r, &msgType); err != nil {
		return
	}
	var tag uint16
	if err = readUint16(r, &tag); err != nil {
		return
	}
	if msgType == Rerror {
		var errmsg string
		if err = readString(r, &errmsg); err != nil {
			return
		}
		err = errors.New(errmsg)
		return
	}
	if msgType != Rattach {
		err = errUnexpectedMsg
		return
	}
	if err = readQID(r, &qid); err != nil {
		return
	}
	if *debugLog {
		log.Println("->", "Rattach", "tag", tag, "qid", qid)
	}
	return
}

// size[4] Rclunk tag[2]
func readRclunk(r io.Reader) (err error) {
	var size uint32
	if err = readUint32(r, &size); err != nil {
		return
	}
	var msgType uint8
	if err = readUint8(r, &msgType); err != nil {
		return
	}
	var tag uint16
	if err = readUint16(r, &tag); err != nil {
		return
	}
	if msgType == Rerror {
		var errmsg string
		if err = readString(r, &errmsg); err != nil {
			return
		}
		err = errors.New(errmsg)
		return
	}
	if msgType != Rclunk {
		err = errUnexpectedMsg
		return
	}
	if *debugLog {
		log.Println("->", "Rclunk", "tag", tag)
	}
	return
}

// size[4] Rerror tag[2] ename[s]
func readRerror(r io.Reader) (ename string, err error) {
	var size uint32
	if err = readUint32(r, &size); err != nil {
		return
	}
	var msgType uint8
	if err = readUint8(r, &msgType); err != nil {
		return
	}
	var tag uint16
	if err = readUint16(r, &tag); err != nil {
		return
	}
	if msgType == Rerror {
		var errmsg string
		if err = readString(r, &errmsg); err != nil {
			return
		}
		err = errors.New(errmsg)
		return
	}
	if msgType != Rerror {
		err = errUnexpectedMsg
		return
	}
	if err = readString(r, &ename); err != nil {
		return
	}
	if *debugLog {
		log.Println("->", "Rerror", "tag", tag, "ename", ename)
	}
	return
}

// size[4] Rflush tag[2]
func readRflush(r io.Reader) (err error) {
	var size uint32
	if err = readUint32(r, &size); err != nil {
		return
	}
	var msgType uint8
	if err = readUint8(r, &msgType); err != nil {
		return
	}
	var tag uint16
	if err = readUint16(r, &tag); err != nil {
		return
	}
	if msgType == Rerror {
		var errmsg string
		if err = readString(r, &errmsg); err != nil {
			return
		}
		err = errors.New(errmsg)
		return
	}
	if msgType != Rflush {
		err = errUnexpectedMsg
		return
	}
	if *debugLog {
		log.Println("->", "Rflush", "tag", tag)
	}
	return
}

// size[4] Ropen tag[2] qid[13] iounit[4]
func readRopen(r io.Reader) (qid QID, iounit uint32, err error) {
	var size uint32
	if err = readUint32(r, &size); err != nil {
		return
	}
	var msgType uint8
	if err = readUint8(r, &msgType); err != nil {
		return
	}
	var tag uint16
	if err = readUint16(r, &tag); err != nil {
		return
	}
	if msgType == Rerror {
		var errmsg string
		if err = readString(r, &errmsg); err != nil {
			return
		}
		err = errors.New(errmsg)
		return
	}
	if msgType != Ropen {
		err = errUnexpectedMsg
		return
	}
	if err = readQID(r, &qid); err != nil {
		return
	}
	if err = readUint32(r, &iounit); err != nil {
		return
	}
	if *debugLog {
		log.Println("->", "Ropen", "tag", tag, "qid", qid, "iounit", iounit)
	}
	return
}

// size[4] Rcreate tag[2] qid[13] iounit[4]
func readRcreate(r io.Reader) (qid QID, iounit uint32, err error) {
	var size uint32
	if err = readUint32(r, &size); err != nil {
		return
	}
	var msgType uint8
	if err = readUint8(r, &msgType); err != nil {
		return
	}
	var tag uint16
	if err = readUint16(r, &tag); err != nil {
		return
	}
	if msgType == Rerror {
		var errmsg string
		if err = readString(r, &errmsg); err != nil {
			return
		}
		err = errors.New(errmsg)
		return
	}
	if msgType != Rcreate {
		err = errUnexpectedMsg
		return
	}
	if err = readQID(r, &qid); err != nil {
		return
	}
	if err = readUint32(r, &iounit); err != nil {
		return
	}
	if *debugLog {
		log.Println("->", "Rcreate", "tag", tag, "qid", qid, "iounit", iounit)
	}
	return
}

// size[4] Ropenfd tag[2] qid[13] iounit[4] unixfd[4]
func readRopenfd(r io.Reader) (qid QID, iounit uint32, unixfd uint32, err error) {
	var size uint32
	if err = readUint32(r, &size); err != nil {
		return
	}
	var msgType uint8
	if err = readUint8(r, &msgType); err != nil {
		return
	}
	var tag uint16
	if err = readUint16(r, &tag); err != nil {
		return
	}
	if msgType == Rerror {
		var errmsg string
		if err = readString(r, &errmsg); err != nil {
			return
		}
		err = errors.New(errmsg)
		return
	}
	if msgType != Ropenfd {
		err = errUnexpectedMsg
		return
	}
	if err = readQID(r, &qid); err != nil {
		return
	}
	if err = readUint32(r, &iounit); err != nil {
		return
	}
	if err = readUint32(r, &unixfd); err != nil {
		return
	}
	if *debugLog {
		log.Println("->", "Ropenfd", "tag", tag, "qid", qid, "iounit", iounit, "unixfd", unixfd)
	}
	return
}

// size[4] Rread tag[2] data[count[4]]
func readRread(r io.Reader, data []byte) (n uint32, err error) {
	var size uint32
	if err = readUint32(r, &size); err != nil {
		return
	}
	var msgType uint8
	if err = readUint8(r, &msgType); err != nil {
		return
	}
	var tag uint16
	if err = readUint16(r, &tag); err != nil {
		return
	}
	if msgType == Rerror {
		var errmsg string
		if err = readString(r, &errmsg); err != nil {
			return
		}
		err = errors.New(errmsg)
		return
	}
	if msgType != Rread {
		err = errUnexpectedMsg
		return
	}
	if n, err = readAndFillByteSlice(r, data); err != nil {
		return
	}
	if *debugLog {
		log.Println("->", "Rread", "tag", tag, "data", data[:n])
	}
	return
}

// size[4] Rwrite tag[2] count[4]
func readRwrite(r io.Reader) (count uint32, err error) {
	var size uint32
	if err = readUint32(r, &size); err != nil {
		return
	}
	var msgType uint8
	if err = readUint8(r, &msgType); err != nil {
		return
	}
	var tag uint16
	if err = readUint16(r, &tag); err != nil {
		return
	}
	if msgType == Rerror {
		var errmsg string
		if err = readString(r, &errmsg); err != nil {
			return
		}
		err = errors.New(errmsg)
		return
	}
	if msgType != Rwrite {
		err = errUnexpectedMsg
		return
	}
	if err = readUint32(r, &count); err != nil {
		return
	}
	if *debugLog {
		log.Println("->", "Rwrite", "tag", tag, "count", count)
	}
	return
}

// size[4] Rremove tag[2]
func readRremove(r io.Reader) (err error) {
	var size uint32
	if err = readUint32(r, &size); err != nil {
		return
	}
	var msgType uint8
	if err = readUint8(r, &msgType); err != nil {
		return
	}
	var tag uint16
	if err = readUint16(r, &tag); err != nil {
		return
	}
	if msgType == Rerror {
		var errmsg string
		if err = readString(r, &errmsg); err != nil {
			return
		}
		err = errors.New(errmsg)
		return
	}
	if msgType != Rremove {
		err = errUnexpectedMsg
		return
	}
	if *debugLog {
		log.Println("->", "Rremove", "tag", tag)
	}
	return
}

// size[4] Rstat tag[2] stat[n]
func readRstat(r io.Reader) (stat Stat, err error) {
	var size uint32
	if err = readUint32(r, &size); err != nil {
		return
	}
	var msgType uint8
	if err = readUint8(r, &msgType); err != nil {
		return
	}
	var tag uint16
	if err = readUint16(r, &tag); err != nil {
		return
	}
	if msgType == Rerror {
		var errmsg string
		if err = readString(r, &errmsg); err != nil {
			return
		}
		err = errors.New(errmsg)
		return
	}
	if msgType != Rstat {
		err = errUnexpectedMsg
		return
	}
	// TODO: Why is this doubly size delimited?
	var outerStatSize uint16
	if err = readUint16(r, &outerStatSize); err != nil {
		return
	}
	if err = readStat(r, &stat); err != nil {
		return
	}
	if *debugLog {
		log.Println("->", "Rstat", "tag", tag, "stat", stat)
	}
	return
}

// size[4] Rwstat tag[2]
func readRwstat(r io.Reader) (err error) {
	var size uint32
	if err = readUint32(r, &size); err != nil {
		return
	}
	var msgType uint8
	if err = readUint8(r, &msgType); err != nil {
		return
	}
	var tag uint16
	if err = readUint16(r, &tag); err != nil {
		return
	}
	if msgType == Rerror {
		var errmsg string
		if err = readString(r, &errmsg); err != nil {
			return
		}
		err = errors.New(errmsg)
		return
	}
	if msgType != Rwstat {
		err = errUnexpectedMsg
		return
	}
	if *debugLog {
		log.Println("->", "Rwstat", "tag", tag)
	}
	return
}

// size[4] Rversion tag[2] msize[4] version[s]
func readRversion(r io.Reader) (msize uint32, version string, err error) {
	var size uint32
	if err = readUint32(r, &size); err != nil {
		return
	}
	var msgType uint8
	if err = readUint8(r, &msgType); err != nil {
		return
	}
	var tag uint16
	if err = readUint16(r, &tag); err != nil {
		return
	}
	if msgType == Rerror {
		var errmsg string
		if err = readString(r, &errmsg); err != nil {
			return
		}
		err = errors.New(errmsg)
		return
	}
	if msgType != Rversion {
		err = errUnexpectedMsg
		return
	}
	if err = readUint32(r, &msize); err != nil {
		return
	}
	if err = readString(r, &version); err != nil {
		return
	}
	if *debugLog {
		log.Println("->", "Rversion", "tag", tag, "msize", msize, "version", version)
	}
	return
}

// size[4] Rwalk tag[2] nwqid*(qid[13])
func readRwalk(r io.Reader) (qids []QID, err error) {
	var size uint32
	if err = readUint32(r, &size); err != nil {
		return
	}
	var msgType uint8
	if err = readUint8(r, &msgType); err != nil {
		return
	}
	var tag uint16
	if err = readUint16(r, &tag); err != nil {
		return
	}
	if msgType == Rerror {
		var errmsg string
		if err = readString(r, &errmsg); err != nil {
			return
		}
		err = errors.New(errmsg)
		return
	}
	if msgType != Rwalk {
		err = errUnexpectedMsg
		return
	}
	if err = readQIDSlice(r, &qids); err != nil {
		return
	}
	if *debugLog {
		log.Println("->", "Rwalk", "tag", tag, "qids", qids)
	}
	return
}
