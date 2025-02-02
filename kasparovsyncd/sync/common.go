package sync

import (
	rpcmodel "github.com/c4ei/c4exd/rpc/model"
)

func stringsSetToSlice(set map[string]struct{}) []string {
	ids := make([]string, len(set))
	i := 0
	for id := range set {
		ids[i] = id
		i++
	}
	return ids
}

// rawAndVerboseBlock is a type that holds either
// the block hexadecimal raw representation and
// either its verbose representation.
type rawAndVerboseBlock struct {
	Raw     string
	Verbose *rpcmodel.GetBlockVerboseResult
}

func (r *rawAndVerboseBlock) String() string {
	return r.hash()
}

// hash returns the block hash
func (r *rawAndVerboseBlock) hash() string {
	return r.Verbose.Hash
}
