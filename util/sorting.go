package util

import (
	. "github.com/FactomProject/factomd/common/AdminBlock"
	. "github.com/FactomProject/factomd/common/DirectoryBlock"
	. "github.com/FactomProject/factomd/common/EntryBlock"
	. "github.com/FactomProject/factomd/common/EntryCreditBlock"
	. "github.com/FactomProject/factomd/common/interfaces"
)

//------------------------------------------------
// DBlock array sorting implementation - accending
type ByDBlockIDAccending []DirectoryBlock

func (f ByDBlockIDAccending) Len() int {
	return len(f)
}
func (f ByDBlockIDAccending) Less(i, j int) bool {
	return f[i].Header.DBHeight < f[j].Header.DBHeight
}
func (f ByDBlockIDAccending) Swap(i, j int) {
	f[i], f[j] = f[j], f[i]
}

//------------------------------------------------
// CBlock array sorting implementation - accending
type ByECBlockIDAccending []ECBlock

func (f ByECBlockIDAccending) Len() int {
	return len(f)
}
func (f ByECBlockIDAccending) Less(i, j int) bool {
	return f[i].Header.DBHeight < f[j].Header.DBHeight
}
func (f ByECBlockIDAccending) Swap(i, j int) {
	f[i], f[j] = f[j], f[i]
}

//------------------------------------------------
// ABlock array sorting implementation - accending
type ByABlockIDAccending []AdminBlock

func (f ByABlockIDAccending) Len() int {
	return len(f)
}
func (f ByABlockIDAccending) Less(i, j int) bool {
	return f[i].Header.DBHeight < f[j].Header.DBHeight
}
func (f ByABlockIDAccending) Swap(i, j int) {
	f[i], f[j] = f[j], f[i]
}

//------------------------------------------------
// ABlock array sorting implementation - accending
type ByFBlockIDAccending []IFBlock

func (f ByFBlockIDAccending) Len() int {
	return len(f)
}
func (f ByFBlockIDAccending) Less(i, j int) bool {
	return f[i].GetDBHeight() < f[j].GetDBHeight()
}
func (f ByFBlockIDAccending) Swap(i, j int) {
	f[i], f[j] = f[j], f[i]
}

//------------------------------------------------
// EBlock array sorting implementation - accending
type ByEBlockIDAccending []EBlock

func (f ByEBlockIDAccending) Len() int {
	return len(f)
}
func (f ByEBlockIDAccending) Less(i, j int) bool {
	return f[i].Header.EBSequence < f[j].Header.EBSequence
}
func (f ByEBlockIDAccending) Swap(i, j int) {
	f[i], f[j] = f[j], f[i]
}