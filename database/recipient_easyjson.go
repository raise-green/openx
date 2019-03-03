// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package database

import (
	json "encoding/json"
	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
)

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjson5e116a33DecodeGithubComYaleOpenLabOpenxDatabase(in *jlexer.Lexer, out *Recipient) {
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
		case "U":
			(out.U).UnmarshalEasyJSON(in)
		case "ReceivedSolarProjects":
			if in.IsNull() {
				in.Skip()
				out.ReceivedSolarProjects = nil
			} else {
				in.Delim('[')
				if out.ReceivedSolarProjects == nil {
					if !in.IsDelim(']') {
						out.ReceivedSolarProjects = make([]string, 0, 4)
					} else {
						out.ReceivedSolarProjects = []string{}
					}
				} else {
					out.ReceivedSolarProjects = (out.ReceivedSolarProjects)[:0]
				}
				for !in.IsDelim(']') {
					var v1 string
					v1 = string(in.String())
					out.ReceivedSolarProjects = append(out.ReceivedSolarProjects, v1)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "ReceivedConstructionBonds":
			if in.IsNull() {
				in.Skip()
				out.ReceivedConstructionBonds = nil
			} else {
				in.Delim('[')
				if out.ReceivedConstructionBonds == nil {
					if !in.IsDelim(']') {
						out.ReceivedConstructionBonds = make([]string, 0, 4)
					} else {
						out.ReceivedConstructionBonds = []string{}
					}
				} else {
					out.ReceivedConstructionBonds = (out.ReceivedConstructionBonds)[:0]
				}
				for !in.IsDelim(']') {
					var v2 string
					v2 = string(in.String())
					out.ReceivedConstructionBonds = append(out.ReceivedConstructionBonds, v2)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "DeviceId":
			out.DeviceId = string(in.String())
		case "DeviceStarts":
			if in.IsNull() {
				in.Skip()
				out.DeviceStarts = nil
			} else {
				in.Delim('[')
				if out.DeviceStarts == nil {
					if !in.IsDelim(']') {
						out.DeviceStarts = make([]string, 0, 4)
					} else {
						out.DeviceStarts = []string{}
					}
				} else {
					out.DeviceStarts = (out.DeviceStarts)[:0]
				}
				for !in.IsDelim(']') {
					var v3 string
					v3 = string(in.String())
					out.DeviceStarts = append(out.DeviceStarts, v3)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "DeviceLocation":
			out.DeviceLocation = string(in.String())
		case "StateHashes":
			if in.IsNull() {
				in.Skip()
				out.StateHashes = nil
			} else {
				in.Delim('[')
				if out.StateHashes == nil {
					if !in.IsDelim(']') {
						out.StateHashes = make([]string, 0, 4)
					} else {
						out.StateHashes = []string{}
					}
				} else {
					out.StateHashes = (out.StateHashes)[:0]
				}
				for !in.IsDelim(']') {
					var v4 string
					v4 = string(in.String())
					out.StateHashes = append(out.StateHashes, v4)
					in.WantComma()
				}
				in.Delim(']')
			}
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
func easyjson5e116a33EncodeGithubComYaleOpenLabOpenxDatabase(out *jwriter.Writer, in Recipient) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"U\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		(in.U).MarshalEasyJSON(out)
	}
	{
		const prefix string = ",\"ReceivedSolarProjects\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		if in.ReceivedSolarProjects == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v5, v6 := range in.ReceivedSolarProjects {
				if v5 > 0 {
					out.RawByte(',')
				}
				out.String(string(v6))
			}
			out.RawByte(']')
		}
	}
	{
		const prefix string = ",\"ReceivedConstructionBonds\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		if in.ReceivedConstructionBonds == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v7, v8 := range in.ReceivedConstructionBonds {
				if v7 > 0 {
					out.RawByte(',')
				}
				out.String(string(v8))
			}
			out.RawByte(']')
		}
	}
	{
		const prefix string = ",\"DeviceId\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.DeviceId))
	}
	{
		const prefix string = ",\"DeviceStarts\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		if in.DeviceStarts == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v9, v10 := range in.DeviceStarts {
				if v9 > 0 {
					out.RawByte(',')
				}
				out.String(string(v10))
			}
			out.RawByte(']')
		}
	}
	{
		const prefix string = ",\"DeviceLocation\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.DeviceLocation))
	}
	{
		const prefix string = ",\"StateHashes\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		if in.StateHashes == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v11, v12 := range in.StateHashes {
				if v11 > 0 {
					out.RawByte(',')
				}
				out.String(string(v12))
			}
			out.RawByte(']')
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Recipient) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson5e116a33EncodeGithubComYaleOpenLabOpenxDatabase(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Recipient) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson5e116a33EncodeGithubComYaleOpenLabOpenxDatabase(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Recipient) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson5e116a33DecodeGithubComYaleOpenLabOpenxDatabase(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Recipient) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson5e116a33DecodeGithubComYaleOpenLabOpenxDatabase(l, v)
}