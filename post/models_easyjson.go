// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package post

import (
	json "encoding/json"
	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
	time "time"
)

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjsonD2b7633eDecodeGithubComBlackplayertenIdealVisualBackendPost(in *jlexer.Lexer, out *Post) {
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
		case "id":
			out.ID = string(in.String())
		case "photo":
			out.Photo = string(in.String())
		case "photo_index":
			if in.IsNull() {
				in.Skip()
				out.PhotoIndex = nil
			} else {
				if out.PhotoIndex == nil {
					out.PhotoIndex = new(int)
				}
				*out.PhotoIndex = int(in.Int())
			}
		case "date":
			if in.IsNull() {
				in.Skip()
				out.Date = nil
			} else {
				if out.Date == nil {
					out.Date = new(time.Time)
				}
				if data := in.Raw(); in.Ok() {
					in.AddError((*out.Date).UnmarshalJSON(data))
				}
			}
		case "place":
			if in.IsNull() {
				in.Skip()
				out.Place = nil
			} else {
				if out.Place == nil {
					out.Place = new(string)
				}
				*out.Place = string(in.String())
			}
		case "text":
			if in.IsNull() {
				in.Skip()
				out.Text = nil
			} else {
				if out.Text == nil {
					out.Text = new(string)
				}
				*out.Text = string(in.String())
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
func easyjsonD2b7633eEncodeGithubComBlackplayertenIdealVisualBackendPost(out *jwriter.Writer, in Post) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"id\":"
		out.RawString(prefix[1:])
		out.String(string(in.ID))
	}
	{
		const prefix string = ",\"photo\":"
		out.RawString(prefix)
		out.String(string(in.Photo))
	}
	if in.PhotoIndex != nil {
		const prefix string = ",\"photo_index\":"
		out.RawString(prefix)
		out.Int(int(*in.PhotoIndex))
	}
	if in.Date != nil {
		const prefix string = ",\"date\":"
		out.RawString(prefix)
		out.Raw((*in.Date).MarshalJSON())
	}
	if in.Place != nil {
		const prefix string = ",\"place\":"
		out.RawString(prefix)
		out.String(string(*in.Place))
	}
	if in.Text != nil {
		const prefix string = ",\"text\":"
		out.RawString(prefix)
		out.String(string(*in.Text))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Post) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonD2b7633eEncodeGithubComBlackplayertenIdealVisualBackendPost(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Post) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonD2b7633eEncodeGithubComBlackplayertenIdealVisualBackendPost(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Post) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonD2b7633eDecodeGithubComBlackplayertenIdealVisualBackendPost(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Post) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonD2b7633eDecodeGithubComBlackplayertenIdealVisualBackendPost(l, v)
}
