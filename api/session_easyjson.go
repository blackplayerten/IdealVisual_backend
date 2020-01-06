// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package api

import (
	json "encoding/json"
	account "github.com/blackplayerten/IdealVisual_backend/account"
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

func easyjsonA818f49aDecodeGithubComBlackplayertenIdealVisualBackendApi(in *jlexer.Lexer, out *AccWithToken) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	out.Account = new(account.Account)
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
		case "token":
			out.Token = string(in.String())
		case "id":
			out.ID = uint64(in.Uint64())
		case "email":
			out.Email = string(in.String())
		case "username":
			out.Username = string(in.String())
		case "avatar":
			if in.IsNull() {
				in.Skip()
				out.Avatar = nil
			} else {
				if out.Avatar == nil {
					out.Avatar = new(string)
				}
				*out.Avatar = string(in.String())
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
func easyjsonA818f49aEncodeGithubComBlackplayertenIdealVisualBackendApi(out *jwriter.Writer, in AccWithToken) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"token\":"
		out.RawString(prefix[1:])
		out.String(string(in.Token))
	}
	{
		const prefix string = ",\"id\":"
		out.RawString(prefix)
		out.Uint64(uint64(in.ID))
	}
	{
		const prefix string = ",\"email\":"
		out.RawString(prefix)
		out.String(string(in.Email))
	}
	{
		const prefix string = ",\"username\":"
		out.RawString(prefix)
		out.String(string(in.Username))
	}
	if in.Avatar != nil {
		const prefix string = ",\"avatar\":"
		out.RawString(prefix)
		out.String(string(*in.Avatar))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v AccWithToken) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonA818f49aEncodeGithubComBlackplayertenIdealVisualBackendApi(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v AccWithToken) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonA818f49aEncodeGithubComBlackplayertenIdealVisualBackendApi(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *AccWithToken) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonA818f49aDecodeGithubComBlackplayertenIdealVisualBackendApi(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *AccWithToken) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonA818f49aDecodeGithubComBlackplayertenIdealVisualBackendApi(l, v)
}
