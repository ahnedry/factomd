// Copyright 2015 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package interfaces

import (
	"fmt"
)

// This object will hold the public keys for servers that are not
// us, and maybe other information about servers.
type IFctServer interface {
	GetChainID() IHash
	String() string
	GetName() string
}

type FctServer struct {
	ChainID IHash
	Name    string
}

var _ IFctServer = (*FctServer)(nil)

func (s *FctServer) GetChainID() IHash {
	return s.ChainID
}

func (s *FctServer) String() string {
	return fmt.Sprintf("Server[:4]: %x", s.GetChainID().Bytes()[:10])
}

func (s *FctServer) GetName() string {
	return s.Name
}
