package client

import (
	"io"
	"log"
)

// size[4] Tauth tag[2] afid[4] uname[s] aname[s]
func writeTauth(w io.Writer, tag uint16, afid uint32, uname string, aname string) error {
	if *debugLog {
		log.Println("<-", "Tauth", "tag", tag, "afid", afid, "uname", uname, "aname", aname)
	}
	size := uint32(4 + 1 + 2 + 4 + (2 + len(uname)) + (2 + len(aname)))
	if err := writeUint32(w, size); err != nil {
		return err
	}
	if err := writeUint8(w, Tauth); err != nil {
		return err
	}
	if err := writeUint16(w, tag); err != nil {
		return err
	}
	if err := writeUint32(w, afid); err != nil {
		return err
	}
	if err := writeString(w, uname); err != nil {
		return err
	}
	if err := writeString(w, aname); err != nil {
		return err
	}
	return nil
}

// size[4] Tattach tag[2] fid[4] afid[4] uname[s] aname[s]
func writeTattach(w io.Writer, tag uint16, fid uint32, afid uint32, uname string, aname string) error {
	if *debugLog {
		log.Println("<-", "Tattach", "tag", tag, "fid", fid, "afid", afid, "uname", uname, "aname", aname)
	}
	size := uint32(4 + 1 + 2 + 4 + 4 + (2 + len(uname)) + (2 + len(aname)))
	if err := writeUint32(w, size); err != nil {
		return err
	}
	if err := writeUint8(w, Tattach); err != nil {
		return err
	}
	if err := writeUint16(w, tag); err != nil {
		return err
	}
	if err := writeUint32(w, fid); err != nil {
		return err
	}
	if err := writeUint32(w, afid); err != nil {
		return err
	}
	if err := writeString(w, uname); err != nil {
		return err
	}
	if err := writeString(w, aname); err != nil {
		return err
	}
	return nil
}

// size[4] Tclunk tag[2] fid[4]
func writeTclunk(w io.Writer, tag uint16, fid uint32) error {
	if *debugLog {
		log.Println("<-", "Tclunk", "tag", tag, "fid", fid)
	}
	size := uint32(4 + 1 + 2 + 4)
	if err := writeUint32(w, size); err != nil {
		return err
	}
	if err := writeUint8(w, Tclunk); err != nil {
		return err
	}
	if err := writeUint16(w, tag); err != nil {
		return err
	}
	if err := writeUint32(w, fid); err != nil {
		return err
	}
	return nil
}

// size[4] Tflush tag[2] oldtag[2]
func writeTflush(w io.Writer, tag uint16, oldtag uint16) error {
	if *debugLog {
		log.Println("<-", "Tflush", "tag", tag, "oldtag", oldtag)
	}
	size := uint32(4 + 1 + 2 + 2)
	if err := writeUint32(w, size); err != nil {
		return err
	}
	if err := writeUint8(w, Tflush); err != nil {
		return err
	}
	if err := writeUint16(w, tag); err != nil {
		return err
	}
	if err := writeUint16(w, oldtag); err != nil {
		return err
	}
	return nil
}

// size[4] Topen tag[2] fid[4] mode[1]
func writeTopen(w io.Writer, tag uint16, fid uint32, mode uint8) error {
	if *debugLog {
		log.Println("<-", "Topen", "tag", tag, "fid", fid, "mode", mode)
	}
	size := uint32(4 + 1 + 2 + 4 + 1)
	if err := writeUint32(w, size); err != nil {
		return err
	}
	if err := writeUint8(w, Topen); err != nil {
		return err
	}
	if err := writeUint16(w, tag); err != nil {
		return err
	}
	if err := writeUint32(w, fid); err != nil {
		return err
	}
	if err := writeUint8(w, mode); err != nil {
		return err
	}
	return nil
}

// size[4] Tcreate tag[2] fid[4] name[s] perm[4] mode[1]
func writeTcreate(w io.Writer, tag uint16, fid uint32, name string, perm uint32, mode uint8) error {
	if *debugLog {
		log.Println("<-", "Tcreate", "tag", tag, "fid", fid, "name", name, "perm", perm, "mode", mode)
	}
	size := uint32(4 + 1 + 2 + 4 + (2 + len(name)) + 4 + 1)
	if err := writeUint32(w, size); err != nil {
		return err
	}
	if err := writeUint8(w, Tcreate); err != nil {
		return err
	}
	if err := writeUint16(w, tag); err != nil {
		return err
	}
	if err := writeUint32(w, fid); err != nil {
		return err
	}
	if err := writeString(w, name); err != nil {
		return err
	}
	if err := writeUint32(w, perm); err != nil {
		return err
	}
	if err := writeUint8(w, mode); err != nil {
		return err
	}
	return nil
}

// size[4] Topenfd tag[2] fid[4] mode[1]
func writeTopenfd(w io.Writer, tag uint16, fid uint32, mode uint8) error {
	if *debugLog {
		log.Println("<-", "Topenfd", "tag", tag, "fid", fid, "mode", mode)
	}
	size := uint32(4 + 1 + 2 + 4 + 1)
	if err := writeUint32(w, size); err != nil {
		return err
	}
	if err := writeUint8(w, Topenfd); err != nil {
		return err
	}
	if err := writeUint16(w, tag); err != nil {
		return err
	}
	if err := writeUint32(w, fid); err != nil {
		return err
	}
	if err := writeUint8(w, mode); err != nil {
		return err
	}
	return nil
}

// size[4] Tread tag[2] fid[4] offset[8] count[4]
func writeTread(w io.Writer, tag uint16, fid uint32, offset uint64, count uint32) error {
	if *debugLog {
		log.Println("<-", "Tread", "tag", tag, "fid", fid, "offset", offset, "count", count)
	}
	size := uint32(4 + 1 + 2 + 4 + 8 + 4)
	if err := writeUint32(w, size); err != nil {
		return err
	}
	if err := writeUint8(w, Tread); err != nil {
		return err
	}
	if err := writeUint16(w, tag); err != nil {
		return err
	}
	if err := writeUint32(w, fid); err != nil {
		return err
	}
	if err := writeUint64(w, offset); err != nil {
		return err
	}
	if err := writeUint32(w, count); err != nil {
		return err
	}
	return nil
}

// size[4] Twrite tag[2] fid[4] offset[8] data[count[4]]
func writeTwrite(w io.Writer, tag uint16, fid uint32, offset uint64, data []byte) error {
	if *debugLog {
		log.Println("<-", "Twrite", "tag", tag, "fid", fid, "offset", offset, "data", data)
	}
	size := uint32(4 + 1 + 2 + 4 + 8 + (4 + len(data)))
	if err := writeUint32(w, size); err != nil {
		return err
	}
	if err := writeUint8(w, Twrite); err != nil {
		return err
	}
	if err := writeUint16(w, tag); err != nil {
		return err
	}
	if err := writeUint32(w, fid); err != nil {
		return err
	}
	if err := writeUint64(w, offset); err != nil {
		return err
	}
	if err := writeByteSlice(w, data); err != nil {
		return err
	}
	return nil
}

// size[4] Tremove tag[2] fid[4]
func writeTremove(w io.Writer, tag uint16, fid uint32) error {
	if *debugLog {
		log.Println("<-", "Tremove", "tag", tag, "fid", fid)
	}
	size := uint32(4 + 1 + 2 + 4)
	if err := writeUint32(w, size); err != nil {
		return err
	}
	if err := writeUint8(w, Tremove); err != nil {
		return err
	}
	if err := writeUint16(w, tag); err != nil {
		return err
	}
	if err := writeUint32(w, fid); err != nil {
		return err
	}
	return nil
}

// size[4] Tstat tag[2] fid[4]
func writeTstat(w io.Writer, tag uint16, fid uint32) error {
	if *debugLog {
		log.Println("<-", "Tstat", "tag", tag, "fid", fid)
	}
	size := uint32(4 + 1 + 2 + 4)
	if err := writeUint32(w, size); err != nil {
		return err
	}
	if err := writeUint8(w, Tstat); err != nil {
		return err
	}
	if err := writeUint16(w, tag); err != nil {
		return err
	}
	if err := writeUint32(w, fid); err != nil {
		return err
	}
	return nil
}

// size[4] Twstat tag[2] fid[4] stat[n]
func writeTwstat(w io.Writer, tag uint16, fid uint32, stat Stat) error {
	if *debugLog {
		log.Println("<-", "Twstat", "tag", tag, "fid", fid, "stat", stat)
	}
	size := uint32(4 + 1 + 2 + 4 + (39 + 8 + len(stat.Name) + len(stat.UID) + len(stat.GID) + len(stat.MUID)))
	if err := writeUint32(w, size); err != nil {
		return err
	}
	if err := writeUint8(w, Twstat); err != nil {
		return err
	}
	if err := writeUint16(w, tag); err != nil {
		return err
	}
	if err := writeUint32(w, fid); err != nil {
		return err
	}
	if err := writeStat(w, stat); err != nil {
		return err
	}
	return nil
}

// size[4] Tversion tag[2] msize[4] version[s]
func writeTversion(w io.Writer, tag uint16, msize uint32, version string) error {
	if *debugLog {
		log.Println("<-", "Tversion", "tag", tag, "msize", msize, "version", version)
	}
	size := uint32(4 + 1 + 2 + 4 + (2 + len(version)))
	if err := writeUint32(w, size); err != nil {
		return err
	}
	if err := writeUint8(w, Tversion); err != nil {
		return err
	}
	if err := writeUint16(w, tag); err != nil {
		return err
	}
	if err := writeUint32(w, msize); err != nil {
		return err
	}
	if err := writeString(w, version); err != nil {
		return err
	}
	return nil
}

// size[4] Twalk tag[2] fid[4] newfid[4] nwname*(wname[s])
func writeTwalk(w io.Writer, tag uint16, fid uint32, newfid uint32, nwnames []string) error {
	if *debugLog {
		log.Println("<-", "Twalk", "tag", tag, "fid", fid, "newfid", newfid, "nwnames", nwnames)
	}
	size := uint32(4 + 1 + 2 + 4 + 4 + stringSliceSize(nwnames))
	if err := writeUint32(w, size); err != nil {
		return err
	}
	if err := writeUint8(w, Twalk); err != nil {
		return err
	}
	if err := writeUint16(w, tag); err != nil {
		return err
	}
	if err := writeUint32(w, fid); err != nil {
		return err
	}
	if err := writeUint32(w, newfid); err != nil {
		return err
	}
	if err := writeStringSlice(w, nwnames); err != nil {
		return err
	}
	return nil
}
