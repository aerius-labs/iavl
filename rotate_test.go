package iavl

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math"
	"testing"
	"bytes"

	"github.com/stretchr/testify/require"
	"github.com/cosmos/iavl/internal/encoding"
)

func TestRotation(t *testing.T) {
	fmt.Println("Hello World")
	tree := getTestTree(0)
	for i := 1; i < 30000000; i += 6 {
		key := make([]byte, 4)
		binary.BigEndian.PutUint32(key, uint32(i))
		tree.Set(key, key)
	}
	root := tree.ImmutableTree.Hash()
	t.Logf("initial root: %x", root)

	proofs := []string{}

	for k := 253; k < 500; k++ {
		i := k % 350
		key := make([]byte, 4)
		binary.BigEndian.PutUint32(key, uint32(i))
		proof, err := tree.GetMembershipProof(key)
		if err != nil {
			proof, err = tree.GetNonMembershipProof(key)
			require.NoError(t, err)
		}
		marshalled, err := proof.Marshal()
		require.NoError(t, err)
		proofs = append(proofs, hex.EncodeToString(marshalled))
	}

	newvaluehashes := []string{}
	roots := []string{}
	for k := 253; k < 500; k++ {
		i := k % 350
		key := make([]byte, 4)
		binary.BigEndian.PutUint32(key, uint32(i))
		tree.Set(key, []byte{byte(i + 1)})

		root = tree.ImmutableTree.Hash()
		roots = append(roots, hex.EncodeToString(root))
		tree.root.traverse(tree.ImmutableTree, true, func(node *Node) bool {
			if string(node.key) == string(key) && node.subtreeHeight == 0 {
				newvaluehashes = append(newvaluehashes, hex.EncodeToString(node.hash))
				return true
			}
			return false
		})
	}
	fmt.Println("final root: ", hex.EncodeToString(root))
	fmt.Println("Proofs: ", proofs)
	fmt.Println("New Value Hashes after updating : ", newvaluehashes)
	fmt.Println("Root Hashes after updating : ", roots)
}

func TestPrefix(t *testing.T) {
	var varintBuf [binary.MaxVarintLen64]byte
	prefix := convertVarIntToBytes(int64(math.MaxInt8), varintBuf)
	prefix = append(prefix, convertVarIntToBytes(math.MaxInt64, varintBuf)...)
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
