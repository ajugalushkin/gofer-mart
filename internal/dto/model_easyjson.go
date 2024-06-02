// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package dto

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

func easyjsonC80ae7adDecodeGithubComAjugalushkinGoferMartInternalDto(in *jlexer.Lexer, out *Withdrawal) {
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
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "Number":
			out.Number = string(in.String())
		case "Sum":
			out.Sum = float64(in.Float64())
		case "ProcessedAt":
			if data := in.Raw(); in.Ok() {
				in.AddError((out.ProcessedAt).UnmarshalJSON(data))
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
func easyjsonC80ae7adEncodeGithubComAjugalushkinGoferMartInternalDto(out *jwriter.Writer, in Withdrawal) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"Number\":"
		out.RawString(prefix[1:])
		out.String(string(in.Number))
	}
	{
		const prefix string = ",\"Sum\":"
		out.RawString(prefix)
		out.Float64(float64(in.Sum))
	}
	{
		const prefix string = ",\"ProcessedAt\":"
		out.RawString(prefix)
		out.Raw((in.ProcessedAt).MarshalJSON())
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Withdrawal) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonC80ae7adEncodeGithubComAjugalushkinGoferMartInternalDto(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Withdrawal) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonC80ae7adEncodeGithubComAjugalushkinGoferMartInternalDto(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Withdrawal) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonC80ae7adDecodeGithubComAjugalushkinGoferMartInternalDto(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Withdrawal) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonC80ae7adDecodeGithubComAjugalushkinGoferMartInternalDto(l, v)
}
func easyjsonC80ae7adDecodeGithubComAjugalushkinGoferMartInternalDto1(in *jlexer.Lexer, out *Withdraw) {
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
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "order":
			out.Order = string(in.String())
		case "sum":
			out.Sum = float64(in.Float64())
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
func easyjsonC80ae7adEncodeGithubComAjugalushkinGoferMartInternalDto1(out *jwriter.Writer, in Withdraw) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"order\":"
		out.RawString(prefix[1:])
		out.String(string(in.Order))
	}
	{
		const prefix string = ",\"sum\":"
		out.RawString(prefix)
		out.Float64(float64(in.Sum))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Withdraw) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonC80ae7adEncodeGithubComAjugalushkinGoferMartInternalDto1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Withdraw) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonC80ae7adEncodeGithubComAjugalushkinGoferMartInternalDto1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Withdraw) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonC80ae7adDecodeGithubComAjugalushkinGoferMartInternalDto1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Withdraw) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonC80ae7adDecodeGithubComAjugalushkinGoferMartInternalDto1(l, v)
}
func easyjsonC80ae7adDecodeGithubComAjugalushkinGoferMartInternalDto2(in *jlexer.Lexer, out *User) {
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
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "login":
			out.Login = string(in.String())
		case "password":
			out.Password = string(in.String())
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
func easyjsonC80ae7adEncodeGithubComAjugalushkinGoferMartInternalDto2(out *jwriter.Writer, in User) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"login\":"
		out.RawString(prefix[1:])
		out.String(string(in.Login))
	}
	{
		const prefix string = ",\"password\":"
		out.RawString(prefix)
		out.String(string(in.Password))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v User) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonC80ae7adEncodeGithubComAjugalushkinGoferMartInternalDto2(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v User) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonC80ae7adEncodeGithubComAjugalushkinGoferMartInternalDto2(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *User) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonC80ae7adDecodeGithubComAjugalushkinGoferMartInternalDto2(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *User) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonC80ae7adDecodeGithubComAjugalushkinGoferMartInternalDto2(l, v)
}
func easyjsonC80ae7adDecodeGithubComAjugalushkinGoferMartInternalDto3(in *jlexer.Lexer, out *OrderList) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		in.Skip()
		*out = nil
	} else {
		in.Delim('[')
		if *out == nil {
			if !in.IsDelim(']') {
				*out = make(OrderList, 0, 0)
			} else {
				*out = OrderList{}
			}
		} else {
			*out = (*out)[:0]
		}
		for !in.IsDelim(']') {
			var v1 Order
			(v1).UnmarshalEasyJSON(in)
			*out = append(*out, v1)
			in.WantComma()
		}
		in.Delim(']')
	}
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonC80ae7adEncodeGithubComAjugalushkinGoferMartInternalDto3(out *jwriter.Writer, in OrderList) {
	if in == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
		out.RawString("null")
	} else {
		out.RawByte('[')
		for v2, v3 := range in {
			if v2 > 0 {
				out.RawByte(',')
			}
			(v3).MarshalEasyJSON(out)
		}
		out.RawByte(']')
	}
}

// MarshalJSON supports json.Marshaler interface
func (v OrderList) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonC80ae7adEncodeGithubComAjugalushkinGoferMartInternalDto3(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v OrderList) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonC80ae7adEncodeGithubComAjugalushkinGoferMartInternalDto3(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *OrderList) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonC80ae7adDecodeGithubComAjugalushkinGoferMartInternalDto3(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *OrderList) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonC80ae7adDecodeGithubComAjugalushkinGoferMartInternalDto3(l, v)
}
func easyjsonC80ae7adDecodeGithubComAjugalushkinGoferMartInternalDto4(in *jlexer.Lexer, out *Order) {
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
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "Number":
			out.Number = string(in.String())
		case "Status":
			out.Status = string(in.String())
		case "Accrual":
			if in.IsNull() {
				in.Skip()
				out.Accrual = nil
			} else {
				if out.Accrual == nil {
					out.Accrual = new(float64)
				}
				*out.Accrual = float64(in.Float64())
			}
		case "UploadedAt":
			if data := in.Raw(); in.Ok() {
				in.AddError((out.UploadedAt).UnmarshalJSON(data))
			}
		case "UserID":
			out.UserID = string(in.String())
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
func easyjsonC80ae7adEncodeGithubComAjugalushkinGoferMartInternalDto4(out *jwriter.Writer, in Order) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"Number\":"
		out.RawString(prefix[1:])
		out.String(string(in.Number))
	}
	{
		const prefix string = ",\"Status\":"
		out.RawString(prefix)
		out.String(string(in.Status))
	}
	{
		const prefix string = ",\"Accrual\":"
		out.RawString(prefix)
		if in.Accrual == nil {
			out.RawString("null")
		} else {
			out.Float64(float64(*in.Accrual))
		}
	}
	{
		const prefix string = ",\"UploadedAt\":"
		out.RawString(prefix)
		out.Raw((in.UploadedAt).MarshalJSON())
	}
	{
		const prefix string = ",\"UserID\":"
		out.RawString(prefix)
		out.String(string(in.UserID))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Order) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonC80ae7adEncodeGithubComAjugalushkinGoferMartInternalDto4(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Order) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonC80ae7adEncodeGithubComAjugalushkinGoferMartInternalDto4(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Order) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonC80ae7adDecodeGithubComAjugalushkinGoferMartInternalDto4(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Order) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonC80ae7adDecodeGithubComAjugalushkinGoferMartInternalDto4(l, v)
}
func easyjsonC80ae7adDecodeGithubComAjugalushkinGoferMartInternalDto5(in *jlexer.Lexer, out *Balance) {
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
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "current":
			out.Current = float64(in.Float64())
		case "withdrawn":
			out.Withdrawn = float64(in.Float64())
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
func easyjsonC80ae7adEncodeGithubComAjugalushkinGoferMartInternalDto5(out *jwriter.Writer, in Balance) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"current\":"
		out.RawString(prefix[1:])
		out.Float64(float64(in.Current))
	}
	{
		const prefix string = ",\"withdrawn\":"
		out.RawString(prefix)
		out.Float64(float64(in.Withdrawn))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Balance) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonC80ae7adEncodeGithubComAjugalushkinGoferMartInternalDto5(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Balance) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonC80ae7adEncodeGithubComAjugalushkinGoferMartInternalDto5(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Balance) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonC80ae7adDecodeGithubComAjugalushkinGoferMartInternalDto5(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Balance) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonC80ae7adDecodeGithubComAjugalushkinGoferMartInternalDto5(l, v)
}
func easyjsonC80ae7adDecodeGithubComAjugalushkinGoferMartInternalDto6(in *jlexer.Lexer, out *Accrual) {
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
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "order":
			out.Order = string(in.String())
		case "status":
			out.Status = string(in.String())
		case "accrual":
			out.Accrual = float64(in.Float64())
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
func easyjsonC80ae7adEncodeGithubComAjugalushkinGoferMartInternalDto6(out *jwriter.Writer, in Accrual) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"order\":"
		out.RawString(prefix[1:])
		out.String(string(in.Order))
	}
	{
		const prefix string = ",\"status\":"
		out.RawString(prefix)
		out.String(string(in.Status))
	}
	{
		const prefix string = ",\"accrual\":"
		out.RawString(prefix)
		out.Float64(float64(in.Accrual))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Accrual) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonC80ae7adEncodeGithubComAjugalushkinGoferMartInternalDto6(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Accrual) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonC80ae7adEncodeGithubComAjugalushkinGoferMartInternalDto6(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Accrual) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonC80ae7adDecodeGithubComAjugalushkinGoferMartInternalDto6(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Accrual) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonC80ae7adDecodeGithubComAjugalushkinGoferMartInternalDto6(l, v)
}
