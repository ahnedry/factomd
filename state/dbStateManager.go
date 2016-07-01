// Copyright 2015 FactomProject Authors. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package state

import (
	"encoding/hex"
	"fmt"
	"time"

	"github.com/FactomProject/factomd/common/interfaces"
	"github.com/FactomProject/factomd/common/messages"
	"github.com/FactomProject/factomd/log"
)

var _ = hex.EncodeToString
var _ = fmt.Print
var _ = time.Now()
var _ = log.Print

type DBState struct {
	isNew bool

	DBHash interfaces.IHash
	ABHash interfaces.IHash
	FBHash interfaces.IHash
	ECHash interfaces.IHash

	DirectoryBlock   interfaces.IDirectoryBlock
	AdminBlock       interfaces.IAdminBlock
	FactoidBlock     interfaces.IFBlock
	EntryCreditBlock interfaces.IEntryCreditBlock
	Locked           bool
	ReadyToSave      bool
	Saved            bool
}

type DBStateList struct {
	SrcNetwork          bool // True if I got this block from the network.
	LastTime            interfaces.Timestamp
	SecondsBetweenTests int
	Lastreq             int
	State               *State
	Base                uint32
	Complete            uint32
	DBStates            []*DBState
}

const SecondsBetweenTests = 1 // Default

func (list *DBStateList) String() string {
	str := "\nDBStates\n"
	str = fmt.Sprintf("%s  Base      = %d\n", str, list.Base)
	str = fmt.Sprintf("%s  timestamp = %s\n", str, list.LastTime.String())
	str = fmt.Sprintf("%s  Complete  = %d\n", str, list.Complete)
	rec := "M"
	last := ""
	for i, ds := range list.DBStates {
		rec = "?"
		if ds != nil {
			rec = "nil"
			if ds.DirectoryBlock != nil {
				rec = "x"

				fmt.Println("Justin ddddddddddddddddddddddddddd DBStateList String()")
				dblk, _ := list.State.DB.FetchDBlock(ds.DirectoryBlock.GetKeyMR())
				if dblk != nil {
					rec = "s"
				}

				if ds.Locked {
					rec = rec + "L"
				}

				if ds.ReadyToSave {
					rec = rec + "R"
				}

				if ds.Saved {
					rec = rec + "S"
				}
			}
		}
		if last != "" {
			str = last
		}
		str = fmt.Sprintf("%s  %1s-DState ?-(DState nil) x-(Not in DB) s-(In DB) L-(Locked) R-(Ready to Save) S-(Saved)\n   DState Height: %d\n%s", str, rec, list.Base+uint32(i), ds.String())
		if rec == "?" && last == "" {
			last = str
		}
	}

	return str
}

func (ds *DBState) String() string {
	str := ""
	if ds == nil {
		str = "  DBState = <nil>\n"
	} else if ds.DirectoryBlock == nil {
		str = "  Directory Block = <nil>\n"
	} else {

		str = fmt.Sprintf("%s      DBlk Height   = %v \n", str, ds.DirectoryBlock.GetHeader().GetDBHeight())
		str = fmt.Sprintf("%s      DBlock        = %x \n", str, ds.DirectoryBlock.GetHash().Bytes()[:5])
		str = fmt.Sprintf("%s      ABlock        = %x \n", str, ds.AdminBlock.GetHash().Bytes()[:5])
		str = fmt.Sprintf("%s      FBlock        = %x \n", str, ds.FactoidBlock.GetHash().Bytes()[:5])
		str = fmt.Sprintf("%s      ECBlock       = %x \n", str, ds.EntryCreditBlock.GetHash().Bytes()[:5])
	}
	return str
}

func (list *DBStateList) GetHighestRecordedBlock() uint32 {
	ht := uint32(0)
	for i, dbstate := range list.DBStates {
		if dbstate != nil && dbstate.Locked {
			ht = list.Base + uint32(i)
		}
	}

	return ht
}

// Once a second at most, we check to see if we need to pull down some blocks to catch up.
func (list *DBStateList) Catchup() {

	now := list.State.GetTimestamp()

	dbsHeight := list.GetHighestRecordedBlock()

	// We only check if we need updates once every so often.
	if int(now)/1000-int(list.LastTime)/1000 < SecondsBetweenTests {
		return
	}
	list.LastTime = now

	begin := -1
	end := -1

	// Find the first range of blocks that we don't have.
	for i, v := range list.DBStates {
		if (v == nil || v.DirectoryBlock == nil) && begin < 0 {
			begin = i
		}
		if v == nil {
			end = i
		}
	}

	if begin > 0 {
		begin += int(list.Base)
		end += int(list.Base)
	} else {
		plHeight := list.State.GetHighestKnownBlock()
		// Don't worry about the block initialization case.
		if plHeight < 1 {
			return
		}

		if plHeight >= dbsHeight && plHeight-dbsHeight > 1 {
			begin = int(dbsHeight + 1)
			end = int(plHeight - 1)
		} else {
			return
		}
	}

	list.Lastreq = begin

	end2 := begin + 400
	if end < end2 {
		end2 = end
	}

	msg := messages.NewDBStateMissing(list.State, uint32(begin), uint32(end2))

	if msg != nil {

		fmt.Println("dddd ======================Ask for blocks", begin, end2)

		list.State.RunLeader = false
		list.State.StartDelay = list.State.GetTimestamp()
		list.State.NetworkOutMsgQueue() <- msg
	}

}

func (list *DBStateList) FixupLinks(p *DBState, d *DBState) (progress bool) {

	// If this block is new, then make sure all hashes are fully computed.
	if !d.isNew {
		return
	}

	hash, _ := p.EntryCreditBlock.HeaderHash()
	fmt.Println("")
	fmt.Println("mmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmm")
	fmt.Println(hash.String())
	fmt.Println("mmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmm")
	fmt.Println("")

	d.EntryCreditBlock.GetHeader().SetPrevHeaderHash(hash)

	hash, _ = p.EntryCreditBlock.GetFullHash()
	d.EntryCreditBlock.GetHeader().SetPrevFullHash(hash)
	fmt.Println("")
	fmt.Println("mmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmm2")
	fmt.Println(hash.String())
	fmt.Println("mmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmm2")
	fmt.Println("")

	d.AdminBlock.GetHeader().SetPrevFullHash(hash)

	p.FactoidBlock.SetDBHeight(p.DirectoryBlock.GetHeader().GetDBHeight())
	d.FactoidBlock.SetDBHeight(d.DirectoryBlock.GetHeader().GetDBHeight())
	d.FactoidBlock.SetPrevKeyMR(p.FactoidBlock.GetKeyMR().Bytes())
	d.FactoidBlock.SetPrevFullHash(p.FactoidBlock.GetPrevFullHash().Bytes())

	d.DirectoryBlock.GetHeader().SetPrevFullHash(p.DirectoryBlock.GetFullHash())
	d.DirectoryBlock.GetHeader().SetPrevKeyMR(p.DirectoryBlock.GetKeyMR())
	d.DirectoryBlock.GetHeader().SetTimestamp(list.State.GetLeaderTimestamp())

	fmt.Println("")
	fmt.Println("mmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmm3")
	fmt.Println(d.EntryCreditBlock.String())
	fmt.Println("mmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmm3")
	fmt.Println("")

	d.DirectoryBlock.SetABlockHash(d.AdminBlock)
	d.DirectoryBlock.SetECBlockHash(d.EntryCreditBlock)
	d.DirectoryBlock.SetFBlockHash(d.FactoidBlock)

	pl := list.State.ProcessLists.Get(d.DirectoryBlock.GetHeader().GetDBHeight())

	//for _, eb := range pl.NewEBlocks {
	//	eb.BuildHeader()
	//	eb.BodyKeyMR()
	//	eb.KeyMR()
	//}

	for _, eb := range pl.NewEBlocks {
		key, err := eb.KeyMR()
		if err != nil {
			panic(err.Error())
		}
		d.DirectoryBlock.AddEntry(eb.GetChainID(), key)
		fmt.Println("Justin bbbbbbbbbbbbbbbbbbbb:", key.String())

	}

	fmt.Println("Justin aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa: FixupLinks")
	d.DirectoryBlock.BuildBodyMR()
	d.DirectoryBlock.MarshalBinary()

	progress = true
	d.isNew = false
	return
}

func (list *DBStateList) ProcessBlocks(i int, d *DBState) (progress bool) {
	fmt.Println("Justin hhhhhhhhhhhhhhhhhh ProcBlocks", i, d.Locked)
	if !d.Locked {
		list.LastTime = list.State.GetTimestamp() // If I saved or processed stuff, I'm good for a while

		// Any updates required to the state as established by the AdminBlock are applied here.
		d.AdminBlock.UpdateState(list.State)

		// Process the Factoid End of Block
		fs := list.State.GetFactoidState()
		fs.AddTransactionBlock(d.FactoidBlock)
		fs.AddECBlock(d.EntryCreditBlock)
		fs.ProcessEndOfBlock(list.State)

		// Promote the currently scheduled next FER
		list.State.ProcessRecentFERChainEntries()

		// Step my counter of Complete blocks
		if uint32(i) > list.Complete {
			fmt.Println("Justin hhhhhhhhhhhhhhhhhhhh up", i, list.Complete)
			list.Complete = uint32(i - 1)
		}
		progress = true
		d.Locked = true // Only after all is done will I admit this state has been saved.
	}
	return
}

func (list *DBStateList) SaveDBStateToDB(d *DBState) (progress bool) {
	fmt.Println("Justin gggggggggggb SaveDBStateToDB", d.DirectoryBlock.GetDatabaseHeight())
	dblk, _ := list.State.DB.FetchDBKeyMRByHash(d.DirectoryBlock.GetKeyMR())
	if d.Saved {
		if dblk == nil {
			panic(fmt.Sprintf("Claimed to be found on %s DBHeight %d Hash %x",
				list.State.FactomNodeName,
				d.DirectoryBlock.GetHeader().GetDBHeight(),
				d.DirectoryBlock.GetKeyMR().Bytes()))
		}
		return
	}

	if !d.ReadyToSave {
		return
	}

	// Take the height, and some function of the identity chain, and use that to decide to trim.  That
	// way, not all nodes in a simulation Trim() at the same time.
	v := int(d.DirectoryBlock.GetHeader().GetDBHeight()) + int(list.State.IdentityChainID.Bytes()[0])
	if v%4 == 0 {
		list.State.DB.Trim()
	}

	list.State.DB.StartMultiBatch()

	if err := list.State.DB.ProcessABlockMultiBatch(d.AdminBlock); err != nil {
		panic(err.Error())
	}

	if err := list.State.DB.ProcessFBlockMultiBatch(d.FactoidBlock); err != nil {
		panic(err.Error())
	}

	if err := list.State.DB.ProcessECBlockMultiBatch(d.EntryCreditBlock, false); err != nil {
		panic(err.Error())
	}
	pl := list.State.ProcessLists.Get(d.DirectoryBlock.GetHeader().GetDBHeight())
	if pl != nil {
		for _, eb := range pl.NewEBlocks {
			if err := list.State.DB.ProcessEBlockMultiBatch(eb, false); err != nil {
				panic(err.Error())
			}

			for _, e := range eb.GetBody().GetEBEntries() {
				if err := list.State.DB.InsertEntry(pl.NewEntries[e.Fixed()]); err != nil {
					panic(err.Error())
				}
			}
		}
	}

	if err := list.State.DB.ProcessDBlockMultiBatch(d.DirectoryBlock); err != nil {
		panic(err.Error())
	}

	if err := list.State.DB.ExecuteMultiBatch(); err != nil {
		panic(err.Error())
	}
	progress = true
	d.ReadyToSave = false
	d.Saved = true
	return
}

func (list *DBStateList) UpdateState() (progress bool) {

	list.Catchup()

	for i, d := range list.DBStates {

		//fmt.Printf("dddd %20s %10s --- %10s %10v %10s %10v \n", "DBStateList Update", list.State.FactomNodeName, "Looking at", i, "DBHeight", list.Base+uint32(i))

		// Must process blocks in sequence.  Missing a block says we must stop.
		if d == nil {
			return
		}

		if d.Saved {
			continue
		}

		if i > 0 {
			progress = list.FixupLinks(list.DBStates[i-1], d)
		}
		progress = list.ProcessBlocks(i, d) || progress

		progress = list.SaveDBStateToDB(d) || progress

		// Make sure the directory block is properly synced up with the prior block, if there
		// is one.

	}
	return
}

func (list *DBStateList) Last() *DBState {
	last := (*DBState)(nil)
	for _, ds := range list.DBStates {
		if ds == nil || ds.DirectoryBlock == nil {
			return last
		}
		last = ds
	}
	return last
}

func (list *DBStateList) Highest() uint32 {
	high := list.Base + uint32(len(list.DBStates)) - 1
	if high == 0 && len(list.DBStates) == 1 {
		return 1
	}
	return high
}

func (list *DBStateList) Put(dbState *DBState) {

	// Hold off on any requests if I'm actually processing...
	list.LastTime = list.State.GetTimestamp()

	dblk := dbState.DirectoryBlock
	dbheight := dblk.GetHeader().GetDBHeight()

	// Count completed states, starting from the beginning (since base starts at
	// zero.
	cnt := 0
searchLoop:
	for i, v := range list.DBStates {
		if v == nil || v.DirectoryBlock == nil || !v.Locked {
			if v != nil && v.DirectoryBlock == nil {
				list.DBStates[i] = nil
			}
			break searchLoop
		}
		cnt++
	}

	keep := uint32(2) // How many states to keep around; debugging helps with more.
	if uint32(cnt) > keep {
		var dbstates []*DBState
		dbstates = append(dbstates, list.DBStates[cnt-int(keep):]...)
		list.DBStates = dbstates
		list.Base = list.Base + uint32(cnt) - keep
		list.Complete = list.Complete - uint32(cnt) + keep
		fmt.Println("Justin hhhhhhhhhhhhhhhhhhhhhhhhhhhhf", list.Complete)
	}

	index := int(dbheight) - int(list.Base)

	fmt.Println("Justin, index is:", index, "and list.Complete is", list.Complete, "and cnt is ", cnt)
	// If we have already processed this State, ignore it.
	if index < int(list.Complete) {
		return
	}

	fmt.Println("Justin rrrrrrrrrrrrrrrrrrrrrr Put", dbState.DirectoryBlock.GetHeader().GetDBHeight())
	if dbState.isNew {
		dbState.DirectoryBlock.SetABlockHash(dbState.AdminBlock)
		dbState.DirectoryBlock.SetECBlockHash(dbState.EntryCreditBlock)
		dbState.DirectoryBlock.SetFBlockHash(dbState.FactoidBlock)
	}

	// make room for this entry.
	for len(list.DBStates) <= index {
		list.DBStates = append(list.DBStates, nil)
	}
	if list.DBStates[index] == nil {
		list.DBStates[index] = dbState
	}
	list.Complete = uint32(index)
}

func (list *DBStateList) Get(height int) *DBState {
	if height < 0 {
		return nil
	}
	i := height - int(list.Base)
	if i < 0 || i >= len(list.DBStates) {
		return nil
	}
	return list.DBStates[i]
}

func (list *DBStateList) NewDBState(isNew bool,
	directoryBlock interfaces.IDirectoryBlock,
	adminBlock interfaces.IAdminBlock,
	factoidBlock interfaces.IFBlock,
	entryCreditBlock interfaces.IEntryCreditBlock) *DBState {

	dbState := new(DBState)

	dbState.DBHash = directoryBlock.DatabasePrimaryIndex()
	dbState.ABHash = adminBlock.DatabasePrimaryIndex()
	dbState.FBHash = factoidBlock.DatabasePrimaryIndex()
	dbState.ECHash = entryCreditBlock.DatabasePrimaryIndex()

	dbState.isNew = isNew
	dbState.DirectoryBlock = directoryBlock
	dbState.AdminBlock = adminBlock
	dbState.FactoidBlock = factoidBlock
	dbState.EntryCreditBlock = entryCreditBlock

	fmt.Println("Justin rrrrrrrrrrrrrrrrrrrrrrrrrrrrrx NewDBState", isNew, directoryBlock.GetDatabaseHeight())
	list.Put(dbState)

	return dbState
}
