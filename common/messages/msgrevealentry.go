// Copyright (c) 2013-2014 Conformal Systems LLC.
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package messages

import (
	"bytes"
	. "github.com/FactomProject/factomd/common/entryBlock"
	. "github.com/FactomProject/factomd/common/interfaces"
	"io"
)

// MsgRevealEntry implements the Message interface and represents a factom
// Reveal-Entry message.  It is used by client to reveal the entry.
type MsgRevealEntry struct {
	Entry IEBEntry
}

// BtcEncode encodes the receiver to w using the bitcoin protocol encoding.
// This is part of the Message interface implementation.
func (msg *MsgRevealEntry) BtcEncode(w io.Writer, pver uint32) error {

	//Entry
	bytes, err := msg.Entry.MarshalBinary()
	if err != nil {
		return err
	}

	err = writeVarBytes(w, pver, bytes)
	if err != nil {
		return err
	}

	return nil
}

// BtcDecode decodes r using the bitcoin protocol encoding into the receiver.
// This is part of the Message interface implementation.
func (msg *MsgRevealEntry) BtcDecode(r io.Reader, pver uint32) error {
	//Entry
	bytes, err := readVarBytes(r, pver, uint32(MaxAppMsgPayload), CmdRevealEntry)
	if err != nil {
		return err
	}

	msg.Entry = new(Entry)
	err = msg.Entry.UnmarshalBinary(bytes)
	if err != nil {
		return err
	}

	return nil
}

// Command returns the protocol command string for the message.  This is part
// of the Message interface implementation.
func (msg *MsgRevealEntry) Command() string {
	return CmdRevealEntry
}

// MaxPayloadLength returns the maximum length the payload can be for the
// receiver.  This is part of the Message interface implementation.
func (msg *MsgRevealEntry) MaxPayloadLength(pver uint32) uint32 {
	return MaxAppMsgPayload
}

// NewMsgInv returns a new bitcoin inv message that conforms to the Message
// interface.  See MsgInv for details.
func NewMsgRevealEntry() *MsgRevealEntry {
	return &MsgRevealEntry{}
}

// Create a sha hash from the message binary (output of BtcEncode)
func (msg *MsgRevealEntry) Sha() (IHash, error) {

	buf := bytes.NewBuffer(nil)
	msg.BtcEncode(buf, ProtocolVersion)
	var sha IHash
	_ = sha.SetBytes(Sha256(buf.Bytes()))

	return sha, nil
}

// Check whether the msg can pass the message level validations
func (msg *MsgRevealEntry) IsValid() bool {
	return true
	//return msg.Entry.IsValid()
}