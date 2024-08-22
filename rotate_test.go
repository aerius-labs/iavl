package iavl

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math"
	"testing"

	"github.com/cosmos/iavl/internal/encoding"
	"github.com/stretchr/testify/require"
)

func TestRotation(t *testing.T) {
	tree := getTestTree(0)
	tree.ImmutableTree.version = 1 << 60
	for i := 1; i < (1 << 22); i += 1 {
		key := make([]byte, 40)
		binary.BigEndian.PutUint32(key, uint32(i))
		tree.Set(key, []byte{0})
	}
	root := tree.ImmutableTree.Hash()
	t.Logf("initial root: %x", root)

	// get membership proof for key 38
	key := make([]byte, 40)
	binary.BigEndian.PutUint32(key, uint32((1<<21)-1))
	fmt.Println("key: ", hex.EncodeToString(key))
	proof, err := tree.GetMembershipProof(key)
	require.NoError(t, err)
	// verify membership proof
	valid, err := tree.ImmutableTree.VerifyMembership(proof, key)
	require.NoError(t, err)
	fmt.Println("valid: ", valid)
}

func TestPrefix(t *testing.T) {
	var varintBuf [binary.MaxVarintLen64]byte
	prefix := convertVarIntToBytes(int64(38), varintBuf)
	prefix = append(prefix, convertVarIntToBytes(1<<20, varintBuf)...)
	prefix = append(prefix, convertVarIntToBytes(math.MaxInt64, varintBuf)...)

	fmt.Println("length of prefix: ", len(prefix))
}

func TestEncodings(t *testing.T) {
	// var varintBuf [binary.MaxVarintLen64]byte

	w := new(bytes.Buffer)
	fmt.Println("length of prefix: ", w.Len())
	err := encoding.EncodeVarint(w, int64(math.MaxInt8))
	require.NoError(t, err)
	err = encoding.EncodeVarint(w, math.MaxInt64)
	require.NoError(t, err)
	err = encoding.EncodeVarint(w, math.MaxInt64)
	require.NoError(t, err)

	fmt.Println("length of prefix: ", w.Len())
}
