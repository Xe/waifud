package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

var (
	outfile = flag.String("o", "/dev/stdout", "output file")
	prefix  = flag.String("prefix", "", "Prefix for function names to print")
)

// Message specs, extracted from plan9port.
var msgSpecs [][]string = [][]string{
	{"size[4]", "Tauth", "tag[2]", "afid[4]", "uname[s]", "aname[s]"},
	{"size[4]", "Rauth", "tag[2]", "aqid[13]"},
	{"size[4]", "Tattach", "tag[2]", "fid[4]", "afid[4]", "uname[s]", "aname[s]"},
	{"size[4]", "Rattach", "tag[2]", "qid[13]"},
	{"size[4]", "Tclunk", "tag[2]", "fid[4]"},
	{"size[4]", "Rclunk", "tag[2]"},
	{"size[4]", "Rerror", "tag[2]", "ename[s]"},
	{"size[4]", "Tflush", "tag[2]", "oldtag[2]"},
	{"size[4]", "Rflush", "tag[2]"},
	{"size[4]", "Topen", "tag[2]", "fid[4]", "mode[1]"},
	{"size[4]", "Ropen", "tag[2]", "qid[13]", "iounit[4]"},
	{"size[4]", "Tcreate", "tag[2]", "fid[4]", "name[s]", "perm[4]", "mode[1]"},
	{"size[4]", "Rcreate", "tag[2]", "qid[13]", "iounit[4]"},
	{"size[4]", "Topenfd", "tag[2]", "fid[4]", "mode[1]"},
	{"size[4]", "Ropenfd", "tag[2]", "qid[13]", "iounit[4]", "unixfd[4]"},
	{"size[4]", "Tread", "tag[2]", "fid[4]", "offset[8]", "count[4]"},
	{"size[4]", "Rread", "tag[2]", "count[4]", "data[count]"},
	{"size[4]", "Twrite", "tag[2]", "fid[4]", "offset[8]", "count[4]", "data[count]"},
	{"size[4]", "Rwrite", "tag[2]", "count[4]"},
	{"size[4]", "Tremove", "tag[2]", "fid[4]"},
	{"size[4]", "Rremove", "tag[2]"},
	{"size[4]", "Tstat", "tag[2]", "fid[4]"},
	{"size[4]", "Rstat", "tag[2]", "stat[n]"},
	{"size[4]", "Twstat", "tag[2]", "fid[4]", "stat[n]"},
	{"size[4]", "Rwstat", "tag[2]"},
	{"size[4]", "Tversion", "tag[2]", "msize[4]", "version[s]"},
	{"size[4]", "Rversion", "tag[2]", "msize[4]", "version[s]"},
	{"size[4]", "Twalk", "tag[2]", "fid[4]", "newfid[4]", "nwname[2]", "nwname*(wname[s])"},
	{"size[4]", "Rwalk", "tag[2]", "nwqid[2]", "nwqid*(qid[13])"},
}

func printComment(ss []string) {
	fmt.Println()
	fmt.Print("//")
	for _, s := range ss {
		fmt.Print(" " + s)
	}
	fmt.Println()
}

// returns type, variable name, size calculation code
func getInfo(s string) (string, string, string) {
	name, _, _ := strings.Cut(s, "[")
	switch {
	case strings.HasSuffix(s, "[1]"):
		return "uint8", name, "1"
	case strings.HasSuffix(s, "[2]"):
		return "uint16", name, "2"
	case strings.HasSuffix(s, "[4]"):
		return "uint32", name, "4"
	case strings.HasSuffix(s, "[8]"):
		return "uint64", name, "8"
	case strings.HasSuffix(s, "[13]"):
		return "QID", name, "13"
	case strings.HasSuffix(s, "[s]"):
		return "string", name, fmt.Sprintf("(2 + len(%v))", name)
	case strings.HasPrefix(s, "T") || strings.HasPrefix(s, "R"):
		return "uint8", "msgType", "1"
	case s == "stat[n]":
		return "Stat", name, fmt.Sprintf("(39 + 8 + len(%v.Name) + len(%v.UID) + len(%v.GID) + len(%v.MUID))", name, name, name, name)
	case strings.HasSuffix(s, "[count[4]]"):
		return "[]byte", name, fmt.Sprintf("(4 + len(%v))", name)
	case s == "nwname*(wname[s])":
		return "[]string", "nwnames", "stringSliceSize(nwnames)"
	case s == "nwqid*(qid[13])":
		return "[]QID", "qids", "(2 + 13*len(qids))"
	default:
		log.Fatalf("unknown type: %q", s)
	}
	return "", "", ""
}

func dontReturnTag(name string) bool {
	// We don't want to return the tag when reading reply data;
	// the tags are already peeked in advance of reading by the 9p
	// protocol layer.
	return name[0] == 'R'
}

func printDebugLine(name string, ss []string) {
	request := name[0] == 'T'

	fmt.Println("\tif *debugLog {")
	if request {
		fmt.Print("\t\tlog.Println(\"<-\"")
	} else {
		fmt.Print("\t\tlog.Println(\"->\"")
	}
	for _, s := range ss {
		_, n, _ := getInfo(s)
		if n == "size" {
			continue
		}
		if n == "msgType" {
			fmt.Printf(", \"%v\"", name)
			continue
		}
		if n == "data" && name == "Rread" {
			fmt.Print(", \"data\", data[:n]")
			continue
		}
		fmt.Printf(", \"%v\", %v", n, n)
	}
	fmt.Println(")")
	fmt.Println("\t}")
}

func printReadFunc(ss []string) {
	name := ss[1]
	funcname := "read" + name
	if !strings.HasPrefix(funcname, *prefix) {
		return
	}
	printComment(ss)

	switch funcname {
	case "readRread":
		fmt.Println("func readRread(r io.Reader, data []byte) (n uint32, err error) {")
	case "readTwrite":
		fmt.Println("func readTwrite(r io.Reader, data []byte) (tag uint16, fid uint32, offset uint64, err error) {")
	default:
		fmt.Print("func " + funcname + "(r io.Reader) (")
		for _, s := range ss {
			t, n, _ := getInfo(s)
			if n == "msgType" || n == "size" {
				continue
			}
			if n == "tag" && dontReturnTag(name) {
				continue
			}
			fmt.Print(n + " " + t + ", ")
		}
		fmt.Println("err error) {")
	}

	// Reading
	fmt.Println("\tvar size uint32")
	for _, s := range ss {
		t, n, _ := getInfo(s)
		fname := fmt.Sprintf("read%v", strings.Title(t))
		if t == "[]string" {
			fname = "readStringSlice"
		}
		if t == "[]QID" {
			fname = "readQIDSlice"
		}
		if t == "[]byte" {
			if funcname != "readRread" {
				log.Fatal("[]byte reading is only implemented for readRread right now, tried ", funcname)
			}
			fmt.Printf("\tif n, err = readAndFillByteSlice(r, %v); err != nil {\n", n)
			fmt.Println("\t\treturn")
			fmt.Println("\t}")
			continue
		}
		if s == "stat[n]" {
			fmt.Println("\t// TODO: Why is this doubly size delimited?")
			fmt.Println("\tvar outerStatSize uint16")
			fmt.Println("\tif err = readUint16(r, &outerStatSize); err != nil {")
			fmt.Println("\t\treturn")
			fmt.Println("\t}")
		}
		if n == "msgType" {
			fmt.Println("\tvar msgType uint8")
		}
		if n == "tag" && dontReturnTag(name) {
			fmt.Println("\tvar tag uint16")
		}
		fmt.Printf("\tif err = %v(r, &%v); err != nil {\n", fname, n)
		fmt.Println("\t\treturn")
		fmt.Println("\t}")
		if n == "tag" {
			if name[0] == 'R' {
				// XXX: Check whether this reads the full error message. (Unix extensions?)
				fmt.Println("\tif msgType == Rerror {")
				fmt.Println("\t\tvar errmsg string")
				fmt.Println("\t\tif err = readString(r, &errmsg); err != nil {")
				fmt.Println("\t\t\treturn")
				fmt.Println("\t\t}")
				fmt.Println("\t\terr = errors.New(errmsg)")
				fmt.Println("\t\treturn")
				fmt.Println("\t}")
			}
			fmt.Println("\tif msgType !=", name, "{")
			fmt.Println("\t\terr = errUnexpectedMsg")
			fmt.Println("\t\treturn")
			fmt.Println("\t}")
		}
	}
	printDebugLine(name, ss)

	fmt.Println("\treturn")
	fmt.Println("}")
}

func printWriteFunc(ss []string) {

	name := ss[1]
	var msgType string

	funcname := "write" + name
	if !strings.HasPrefix(funcname, *prefix) {
		return
	}
	printComment(ss)

	fmt.Print("func " + funcname + "(w io.Writer, ")
	for i, s := range ss {
		t, n, _ := getInfo(s)
		// msgType is fixed for each method
		if n == "msgType" {
			msgType = s
			continue
		}
		// size is calculated dynamically based on other parameters
		if n == "size" {
			continue
		}
		fmt.Print(n + " " + t)
		if i < len(ss)-1 {
			fmt.Print(", ")
		}
	}
	fmt.Println(") error {")

	printDebugLine(name, ss)

	// Size calculation
	fmt.Print("\tsize := uint32(")
	for i, s := range ss {
		_, _, sz := getInfo(s)
		fmt.Print(sz)
		if i < len(ss)-1 {
			fmt.Print(" + ")
		}
	}
	fmt.Println(")")

	for _, s := range ss {
		t, n, _ := getInfo(s)
		funcname := fmt.Sprintf("write%v", strings.Title(t))
		if t == "[]string" {
			funcname = "writeStringSlice"
		}
		if t == "[]QID" {
			funcname = "writeQIDSlice"
		}
		if t == "[]byte" {
			funcname = "writeByteSlice"
		}
		if n == "msgType" {
			n = msgType // resolve to constant directly
		}
		fmt.Printf("\tif err := %v(w, %v); err != nil {\n", funcname, n)
		fmt.Println("\t\treturn err")
		fmt.Println("\t}")
	}
	fmt.Println("\treturn nil")
	fmt.Println("}")
}

// Conflate:
// count[4] data[count] => data[count[4]]
func conflate(in []string) []string {
	var o []string
	for i := 0; i < len(in); i++ {
		// byte buffers
		if in[i] == "count[4]" {
			if i+1 < len(in) && in[i+1] == "data[count]" {
				o = append(o, "data[count[4]]")
				i++
				continue
			}
		}
		// directory names for walk
		if in[i] == "nwname[2]" {
			if i+1 < len(in) && in[i+1] == "nwname*(wname[s])" {
				o = append(o, "nwname*(wname[s])")
				i++
				continue
			}
		}
		// qids for walk response
		if in[i] == "nwqid[2]" {
			if i+1 < len(in) && in[i+1] == "nwqid*(qid[13])" {
				o = append(o, "nwqid*(qid[13])")
				i++
				continue
			}
		}
		o = append(o, in[i])
	}
	return o
}

func main() {
	flag.Parse()
	f, err := os.Create(*outfile)
	if err != nil {
		log.Fatal("Create:", err)
	}
	defer f.Close()
	os.Stdout = f

	if strings.HasPrefix(*prefix, "w") {
		fmt.Println(`package ninep

import (
	"io"
	"log"
)`)
	} else {
		fmt.Println(`package ninep

import (
	"errors"
	"io"
	"log"
)`)
	}

	for _, ss := range msgSpecs {
		ss = conflate(ss)
		printWriteFunc(ss)
		printReadFunc(ss)
	}
}
