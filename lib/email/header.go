// Copyright 2019, Shulhan <ms@kilabit.info>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package email

import (
	"bytes"
	"fmt"
	"log"
)

//
// Header represent list of fields in message.
//
type Header struct {
	//
	// fields is ordered from top to bottom, the first field in message
	// header is equal to the first element in slice.
	//
	// We are not using map here it to prevent the header being reordeded
	// when packing the message back into raw format.
	//
	fields []*Field
}

//
// ParseHeader parse the raw header from top to bottom.
//
// Raw header that start with CRLF indicate an empty header.
// In this case, it will return nil Header, indicating that no header was
// parsed, and the leading CRLF is removed on returned "rest".
//
// The raw header may end with optional CRLF, an empty line that separate
// header from body of message.
//
// On success it will return the rest of raw input (possible message's body)
// without leading CRLF.
//
func ParseHeader(raw []byte) (hdr *Header, rest []byte, err error) {
	var (
		field *Field
	)
	if len(raw) == 0 {
		return nil, nil, nil
	}

	rest = raw
	for len(rest) >= 2 {
		if rest[0] == cr && rest[1] == lf {
			rest = rest[2:]
			return hdr, rest, nil
		}

		field, rest, err = ParseField(rest)
		if err != nil {
			return nil, rest, err
		}
		if hdr == nil {
			hdr = &Header{}
		}
		hdr.fields = append(hdr.fields, field)
	}
	if len(rest) == 0 {
		return hdr, rest, nil
	}

	err = fmt.Errorf("ParseHeader: invalid end of header: '%s'", rest)

	return nil, rest, err
}

//
// Boundary return the message body boundary defined in Content-Type.
// If no field Content-Type or no boundary it will return nil.
//
func (hdr *Header) Boundary() []byte {
	ct := hdr.ContentType()
	if ct == nil {
		return nil
	}
	return ct.GetParamValue(ParamNameBoundary)
}

//
// ContentType return the unpacked value of field "Content-Type", or nil if no
// field Content-Type exist or there is an error when unpacking.
//
func (hdr *Header) ContentType() *ContentType {
	for _, f := range hdr.fields {
		if f.Type != FieldTypeContentType {
			continue
		}
		if f.ContentType == nil {
			err := f.Unpack()
			if err != nil {
				log.Println("ContentType: ", err)
				return nil
			}
		}
		return f.ContentType
	}
	return nil
}

//
// DKIM return sub-header of the "n" DKIM-Signature, start from the top.
// If no DKIM-Signature found it will return nil.
//
// For example, to get the second DKIM-Signature from the top, call it with
// "n=2", but if no second DKIM-Signature it will return nil.
//
func (hdr *Header) DKIM(n int) (dkimHeader *Header) {
	if n == 0 || len(hdr.fields) == 0 {
		return nil
	}

	x := 0
	for ; x < len(hdr.fields); x++ {
		if hdr.fields[x].Type == FieldTypeDKIMSignature {
			n--
			if n == 0 {
				break
			}
		}
	}
	if x == len(hdr.fields) {
		return nil
	}
	dkimHeader = &Header{
		fields: make([]*Field, 0, len(hdr.fields)-x),
	}
	for ; x < len(hdr.fields); x++ {
		dkimHeader.fields = append(dkimHeader.fields, hdr.fields[x])
	}

	return dkimHeader
}

//
// Filter specific field type.  If multiple fields type exist it will
// return all of them.
//
func (hdr *Header) Filter(ft FieldType) (fields []*Field) {
	for x := len(hdr.fields) - 1; x >= 0; x-- {
		if hdr.fields[x].Type == ft {
			fields = append(fields, hdr.fields[x])
		}
	}
	return
}

//
// PushTop put the field at the top of header.
//
func (hdr *Header) PushTop(f *Field) {
	hdr.fields = append([]*Field{f}, hdr.fields...)
}

//
// Relaxed canonicalize the header using "relaxed" algorithm and return it.
//
func (hdr *Header) Relaxed() []byte {
	var bb bytes.Buffer

	for _, f := range hdr.fields {
		if len(f.Name) > 0 && len(f.Value) > 0 {
			bb.Write(f.Name)
			bb.WriteByte(':')
			bb.Write(f.Value)
		}
	}

	return bb.Bytes()
}

//
// Simple canonicalize the header using "simple" algorithm.
//
func (hdr *Header) Simple() []byte {
	var bb bytes.Buffer

	for _, f := range hdr.fields {
		if len(f.oriName) > 0 && len(f.oriValue) > 0 {
			bb.Write(f.oriName)
			bb.WriteByte(':')
			bb.Write(f.oriValue)
		}
	}

	return bb.Bytes()
}

//
// popByName remove the field where the name match from header.
//
func (hdr *Header) popByName(name []byte) (f *Field) {
	for x := len(hdr.fields) - 1; x >= 0; x-- {
		if bytes.Equal(hdr.fields[x].Name, name) {
			f = hdr.fields[x]
			hdr.fields = append(hdr.fields[:x], hdr.fields[x+1:]...)
		}
	}
	return f
}
