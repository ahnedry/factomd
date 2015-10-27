// Copyright 2015 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package entryCreditBlock

import (
	"bytes"
	"github.com/FactomProject/factomd/common/interfaces"
	"github.com/FactomProject/factomd/common/primitives"
)

//var IncreaseBalanceSize int = 32 + 4 + 32

type IncreaseBalance struct {
	ECPubKey *primitives.ByteSlice32
	TXID     interfaces.IHash
	Index    uint64
	NumEC    uint64
}

var _ interfaces.Printable = (*IncreaseBalance)(nil)

//var _ interfaces.BinaryMarshallable = (*IncreaseBalance)(nil)
var _ interfaces.ShortInterpretable = (*IncreaseBalance)(nil)
var _ ECBlockEntry = (*IncreaseBalance)(nil)

//func (c *IncreaseBalance) MarshalledSize() uint64 {
//	return uint64(IncreaseBalanceSize)
//}

func NewIncreaseBalance() *IncreaseBalance {
	r := new(IncreaseBalance)
	r.TXID = primitives.NewZeroHash()
	return r
}

func (e *IncreaseBalance) Hash() interfaces.IHash {
	bin, err := e.MarshalBinary()
	if err != nil {
		panic(err)
	}
	return primitives.Sha(bin)
}

func (b *IncreaseBalance) ECID() byte {
	return ECIDBalanceIncrease
}

func (b *IncreaseBalance) IsInterpretable() bool {
	return false
}

func (b *IncreaseBalance) Interpret() string {
	return ""
}

func (b *IncreaseBalance) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)

	buf.Write(b.ECPubKey[:])

	buf.Write(b.TXID.Bytes())

	primitives.EncodeVarInt(buf, b.Index)

	primitives.EncodeVarInt(buf, b.NumEC)

	return buf.Bytes(), nil
}

func (b *IncreaseBalance) UnmarshalBinaryData(data []byte) (newData []byte, err error) {
	buf := bytes.NewBuffer(data)
	hash := make([]byte, 32)

	_, err = buf.Read(hash)
	if err != nil {
		return
	}
	b.ECPubKey = new(primitives.ByteSlice32)
	copy(b.ECPubKey[:], hash)

	_, err = buf.Read(hash)
	if err != nil {
		return
	}
	if b.TXID == nil {
		b.TXID = primitives.NewZeroHash()
	}
	b.TXID.SetBytes(hash)

	tmp := make([]byte, 0)
	b.Index, tmp = primitives.DecodeVarInt(buf.Bytes())

	b.NumEC, tmp = primitives.DecodeVarInt(tmp)

	newData = tmp
	return
}

func (b *IncreaseBalance) UnmarshalBinary(data []byte) (err error) {
	_, err = b.UnmarshalBinaryData(data)
	return
}

func (e *IncreaseBalance) JSONByte() ([]byte, error) {
	return primitives.EncodeJSON(e)
}

func (e *IncreaseBalance) JSONString() (string, error) {
	return primitives.EncodeJSONString(e)
}

func (e *IncreaseBalance) JSONBuffer(b *bytes.Buffer) error {
	return primitives.EncodeJSONToBuffer(e, b)
}

func (e *IncreaseBalance) String() string {
	str, _ := e.JSONString()
	return str
}