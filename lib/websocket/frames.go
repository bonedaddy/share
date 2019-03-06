// Copyright 2019, Shulhan <ms@kilabit.info>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package websocket

//
// Frames represent continuous (fragmented) frame.
//
// A fragmented message consists of a single frame with the FIN bit clear and
// an opcode other than 0, followed by zero or more frames with the FIN bit
// clear and the opcode set to 0, and terminated by a single frame with the
// FIN bit set and an opcode of 0.
//
type Frames struct {
	v []*Frame
}

//
// Unpack websocket data protocol from raw bytes to one or more frames.
//
// When receiving packet from client, the underlying protocol or operating
// system may buffered the packet.  Client may send a single frame one at
// time, but server may receive one or more frame in one packet; and vice
// versa.  That's the reason why the Unpack return multiple frame instead of
// single frame.
//
// On success it will return one or more frames.
// On fail it will return zero frame.
//
func Unpack(in []byte) (frames *Frames) {
	if len(in) == 0 {
		return
	}

	frames = &Frames{}
	for {
		f, x := frameUnpack(in)
		if f == nil {
			break
		}

		frames.Append(f)
		if x >= uint64(len(in)) {
			break
		}
		in = in[x:]
	}

	return
}

//
// Append a frame as part of continuous frame.
// This function does not check if the appended frame is valid (i.e. zero
// operation code on second or later frame).
//
func (frames *Frames) Append(f *Frame) {
	if f != nil {
		frames.v = append(frames.v, f)
	}
}

//
// Get frame at specific index or nil if index out of range.
//
func (frames *Frames) Get(x int) *Frame {
	if x < 0 || x >= len(frames.v) {
		return nil
	}
	return frames.v[x]
}

//
// IsClosed will return true if one of the frame is control CLOSE frame.
//
func (frames *Frames) IsClosed() bool {
	if len(frames.v) == 0 {
		return false
	}
	for x := 0; x < len(frames.v); x++ {
		if frames.v[x].opcode == opcodeClose {
			return true
		}
	}
	return false
}

//
// Len return the number of frame.
//
func (frames *Frames) Len() int {
	return len(frames.v)
}

//
// Payload return the concatenation of continuous data frame's payload.
//
// The first frame must be a data frame, either text or binary, otherwise it
// will be considered empty payload, even if frames list is not empty.
//
// Any control CLOSE frame of frame with fin set will considered the last
// frame.
//
func (frames *Frames) Payload() (payload []byte) {
	if len(frames.v) == 0 {
		return
	}
	if !frames.v[0].IsData() {
		return
	}

	for x := 0; x < len(frames.v); x++ {
		if frames.v[x].opcode == opcodeClose {
			break
		}

		// Ignore control PING or PONG frame.
		if frames.v[x].opcode == opcodePing || frames.v[x].opcode == opcodePong {
			continue
		}

		payload = append(payload, frames.v[x].payload...)

		if frames.v[x].fin == frameIsFinished {
			break
		}
	}

	return
}
