package state

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/crypto"
)

// Tests that the node iterator indeed walks over the entire database contents.
func TestNodeIteratorCoverage(t *testing.T) {
	testNodeIteratorCoverage(t, rawdb.HashScheme)
	testNodeIteratorCoverage(t, rawdb.PathScheme)
}

func testNodeIteratorCoverage(t *testing.T, scheme string) {
	// Create some arbitrary test state to iterate
	db, sdb, ndb, root, _ := makeTestState(scheme)
	ndb.Commit(root, false)

	state, err := New(root, sdb, nil)
	if err != nil {
		t.Fatalf("failed to create state trie at %x: %v", root, err)
	}
	// Gather all the node hashes found by the iterator
	hashes := make(map[common.Hash]struct{})
	for it := newNodeIterator(state); it.Next(); {
		if it.Hash != (common.Hash{}) {
			hashes[it.Hash] = struct{}{}
		}
	}
	// Check in-disk nodes
	var (
		seenNodes = make(map[common.Hash]struct{})
		seenCodes = make(map[common.Hash]struct{})
	)
	it := db.NewIterator(nil, nil)
	for it.Next() {
		ok, hash := isTrieNode(scheme, it.Key(), it.Value())
		if !ok {
			continue
		}
		seenNodes[hash] = struct{}{}
	}
	it.Release()

	// Check in-disk codes
	it = db.NewIterator(nil, nil)
	for it.Next() {
		ok, hash := rawdb.IsCodeKey(it.Key())
		if !ok {
			continue
		}
		if _, ok := hashes[common.BytesToHash(hash)]; !ok {
			t.Errorf("state entry not reported %x", it.Key())
		}
		seenCodes[common.BytesToHash(hash)] = struct{}{}
	}
	it.Release()

	// Cross check the iterated hashes and the database/nodepool content
	for hash := range hashes {
		_, ok := seenNodes[hash]
		if !ok {
			_, ok = seenCodes[hash]
		}
		if !ok {
			t.Errorf("failed to retrieve reported node %x", hash)
		}
	}
}

// isTrieNode is a helper function which reports if the provided
// database entry belongs to a trie node or not.
func isTrieNode(scheme string, key, val []byte) (bool, common.Hash) {
	if scheme == rawdb.HashScheme {
		if rawdb.IsLegacyTrieNode(key, val) {
			return true, common.BytesToHash(key)
		}
	} else {
		ok := rawdb.IsAccountTrieNode(key)
		if ok {
			return true, crypto.Keccak256Hash(val)
		}
		ok = rawdb.IsStorageTrieNode(key)
		if ok {
			return true, crypto.Keccak256Hash(val)
		}
	}
	return false, common.Hash{}
}
