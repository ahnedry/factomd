// Copyright 2015 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package primitives

import (
	"bytes"
	"encoding/binary"
	"time"

	"github.com/FactomProject/factomd/common/interfaces"
)

func GetTimeMilli() uint64 {
	return uint64(time.Now().UnixNano()) / 1000000 // 10^-9 >> 10^-3
}

func GetTime() uint64 {
	return uint64(time.Now().Unix())
}

//A structure for handling timestamps for messages
type Timestamp uint64 //in miliseconds

var _ interfaces.ITimestamp = (*Timestamp)(nil)

func NewTimestampNow() Timestamp {
	var t Timestamp
	t.SetTimeNow()
	return t
}

func NewTimestampFromSeconds(s uint32) Timestamp {
	return Timestamp(int64(s) * 1000)
}

func NewTimestampFromMinutes(s uint32) Timestamp {
	return Timestamp(int64(s) * 60000)
}

func NewTimestampFromMilliseconds(s uint64) Timestamp {
	return Timestamp(s)
}

func (t Timestamp) SetTimeNow() {
	t = Timestamp(GetTimeMilli())
}

func (t Timestamp) SetTime(miliseconds uint64) {
	t = Timestamp(miliseconds)
}

func (t Timestamp) SetTimestamp(timestamp interfaces.ITimestamp) {
	t.SetTime(uint64(timestamp.GetTimeMilli()))
}

func (t Timestamp) GetTime() time.Time {
	return time.Unix(int64(t/1000), int64(((t)%1000)*1000))
}

func (t Timestamp) UnmarshalBinaryData(data []byte) (newData []byte, err error) {
	hd, data := binary.BigEndian.Uint32(data[:]), data[4:]
	ld, data := binary.BigEndian.Uint16(data[:]), data[2:]
	t = Timestamp((uint64(hd) << 16) + uint64(ld))
	return data, nil
}

func (t Timestamp) UnmarshalBinary(data []byte) error {
	_, err := t.UnmarshalBinaryData(data)
	return err
}

func (t Timestamp) GetTimeSeconds() int64 {
	return int64(t / 1000)
}

func (t Timestamp) GetTimeMinutesUInt32() uint32 {
	return uint32(t / 60000)
}

func (t Timestamp) GetTimeMilli() int64 {
	return int64(t)
}

func (t Timestamp) GetTimeMilliUInt64() uint64 {
	return uint64(t)
}

func (t Timestamp) GetTimeSecondsUInt32() uint32 {
	return uint32(t / 1000)
}

func (t Timestamp) MarshalBinary() ([]byte, error) {
	var out bytes.Buffer
	hd := uint32(t >> 16)
	ld := uint16(t & 0xFFFF)
	binary.Write(&out, binary.BigEndian, uint32(hd))
	binary.Write(&out, binary.BigEndian, uint16(ld))
	return out.Bytes(), nil
}

func (t Timestamp) String() string {
	return t.GetTime().Format("2006-01-02 15:04:05")
}
