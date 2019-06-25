package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"github.com/mailru/easyjson"
	"github.com/mailru/easyjson/jlexer"
	"github.com/mailru/easyjson/jwriter"
	"io"
	"os"
	"strconv"
	"strings"
)

type user struct {
	Browsers []string `json:"browsers"`
	Email    string   `json:"email"`
	Name     string   `json:"name"`
}

// вам надо написать более быструю оптимальную этой функции
func FastSearch(out io.Writer) {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	fileReader := bufio.NewReader(file)
	seenBrowsers := []string{}
	var buf bytes.Buffer
	buf.WriteString("found users:\n")
	out.Write(buf.Bytes())
	curUser := &user{}
	i := 0
	for {
		line, _, err := fileReader.ReadLine()
		if err != nil {
			break
		}
		isAndroid := false
		isMSIE := false
		errJson := curUser.UnmarshalJSON(line)
		if errJson != nil {
			break
		}
		for _, browser := range curUser.Browsers {
			androidFlag := strings.Contains(browser, "Android")
			MSIEFlag := strings.Contains(browser, "MSIE")
			if androidFlag || MSIEFlag {
				if androidFlag {
					isAndroid = true
				}
				if MSIEFlag {
					isMSIE = true
				}
				notSeenBefore := true
				for _, item := range seenBrowsers {
					if item == browser {
						notSeenBefore = false
					}
				}
				if notSeenBefore {
					seenBrowsers = append(seenBrowsers, browser)
				}
			}
		}
		i++
		if !(isAndroid && isMSIE) {
			continue
		}
		email := strings.ReplaceAll(curUser.Email, "@", " [at] ")
		buf.Reset()
		buf.WriteString("[" + strconv.Itoa(i-1) + "] " + curUser.Name + " <" + email + ">\n")
		out.Write(buf.Bytes())
	}
	out.Write([]byte("\nTotal unique browsers " + strconv.Itoa(len(seenBrowsers)) + "\n"))
}

var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjsonCc477bf9DecodeGolang20191599HwEasyJson(in *jlexer.Lexer, out *user) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "browsers":
			if in.IsNull() {
				in.Skip()
				out.Browsers = nil
			} else {
				in.Delim('[')
				if out.Browsers == nil {
					if !in.IsDelim(']') {
						out.Browsers = make([]string, 0, 4)
					} else {
						out.Browsers = []string{}
					}
				} else {
					out.Browsers = (out.Browsers)[:0]
				}
				for !in.IsDelim(']') {
					var v1 string
					v1 = string(in.String())
					out.Browsers = append(out.Browsers, v1)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "email":
			out.Email = string(in.String())
		case "name":
			out.Name = string(in.String())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonCc477bf9EncodeGolang20191599HwEasyJson(out *jwriter.Writer, in user) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"browsers\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		if in.Browsers == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v2, v3 := range in.Browsers {
				if v2 > 0 {
					out.RawByte(',')
				}
				out.String(string(v3))
			}
			out.RawByte(']')
		}
	}
	{
		const prefix string = ",\"email\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Email))
	}
	{
		const prefix string = ",\"name\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Name))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v user) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonCc477bf9EncodeGolang20191599HwEasyJson(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v user) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonCc477bf9EncodeGolang20191599HwEasyJson(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *user) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonCc477bf9DecodeGolang20191599HwEasyJson(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *user) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonCc477bf9DecodeGolang20191599HwEasyJson(l, v)
}
