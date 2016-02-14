package trie

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

type Trie struct {
	Root *Branch
}

/*
NewTrie returns the pointer to a new Trie with an initialized root Branch
*/
func NewTrie() *Trie {
	t := &Trie{
		Root: &Branch{
			Branches: make(map[rune]*Branch),
		},
	}
	return t
}

/*
Add adds an entry to the trie and returns the branch node that the insertion
was made at - or rather where the end of the entry was marked.
*/
func (t *Trie) Add(entry string) *Branch {
	t.Root.Lock()
	b := t.Root.add([]rune(entry))
	t.Root.Unlock()
	return b
}

/*
Delete decrements the count of an existing entry by one. If the count equals
zero it removes an the entry from the trie. Returns true if the entry existed,
false otherwise. Note that the return value says something about the previous
existence of the entry - not whether it has been completely removed or just
its count decremented.
*/
func (t *Trie) Delete(entry string) bool {
	if len(entry) == 0 {
		return false
	}
	t.Root.Lock()
	deleted := t.Root.delete([]rune(entry))
	t.Root.Unlock()
	return deleted
}

/*
GetBranch returns the branch end if the `entry` exists in the `Trie`
*/
func (t *Trie) GetBranch(entry string) *Branch {
	return t.Root.getBranch([]rune(entry))
}

/*
Has returns true if the `entry` exists in the `Trie`
*/
func (t *Trie) Has(entry string) bool {
	return t.Root.has([]rune(entry))
}

/*
HasCount returns true  if the `entry` exists in the `Trie`. The second returned
value is the count how often the entry has been set.
*/
func (t *Trie) HasCount(entry string) (exists bool, count int64) {
	return t.Root.hasCount([]rune(entry))
}

/*
HasPrefix returns true if the the `Trie` contains entries with the given prefix
*/
func (t *Trie) HasPrefix(prefix string) bool {
	return t.Root.hasPrefix([]rune(prefix))
}

/*
HasPrefixCount returns true if the the `Trie` contains entries with the given
prefix. The second returned value is the count how often the entry has been set.
*/
func (t *Trie) HasPrefixCount(prefix string) (exists bool, count int64) {
	return t.Root.hasPrefixCount([]rune(prefix))
}

/*
Members returns all entries of the Trie with their counts as MemberInfo
*/
func (t *Trie) Members() []*MemberInfo {
	return t.Root.members([]rune{})
}

/*
Members returns a Slice of all entries of the Trie
*/
func (t *Trie) MembersList() (members []string) {
	for _, mi := range t.Root.members([]rune{}) {
		members = append(members, mi.Value)
	}
	return
}

/*
PrefixMembers returns all entries of the Trie that have the given prefix
with their counts as MemberInfo
*/
func (t *Trie) PrefixMembers(prefix string) []*MemberInfo {
	return t.Root.prefixMembers([]rune{}, []rune(prefix))
}

/*
PrefixMembers returns a List of all entries of the Trie that have the
given prefix
*/
func (t *Trie) PrefixMembersList(prefix string) (members []string) {
	for _, mi := range t.Root.prefixMembers([]rune{}, []rune(prefix)) {
		members = append(members, mi.Value)
	}
	return
}

/*
Dump returns a string representation of the `Trie`
*/
func (t *Trie) Dump() string {
	return t.Root.Dump(0)
}

/*
 */
func (t *Trie) PrintDump() {
	t.Root.PrintDump()
}

/*
DumpToFile dumps all values into a slice of strings and writes that to a file
using encoding/gob.

The Trie itself can currently not be encoded directly because gob does not
directly support structs with a sync.Mutex on them.
*/
func (t *Trie) DumpToFile(fname string) (err error) {
	t.Root.Lock()
	entries := t.Members()
	t.Root.Unlock()

	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)
	if err = enc.Encode(entries); err != nil {
		err = errors.New(fmt.Sprintf("Could encode Trie entries for dump file: %v", err))
		return
	}

	f, err := os.Create(fname)
	if err != nil {
		err = errors.New(fmt.Sprintf("Could not save dump file: %v", err))
		return
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	_, err = w.Write(buf.Bytes())
	if err != nil {
		err = errors.New(fmt.Sprintf("Error writing to dump file: %v", err))
		return
	}
	// log.Printf("wrote %d bytes to dumpfile %s\n", bl, fname)
	w.Flush()
	return
}

/*
LoadFromFile loads a plain wordlist from a txt file and creates a new Trie
by Add()ing all of them.
*/
func (t *Trie) LoadFromFile(fname string) (tr *Trie, err error) {
	tr = NewTrie()
	startTime := time.Now()

    f, err := os.Open(fname)
    if err != nil {
        panic(err)
    }
    defer f.Close()

    rd := bufio.NewReader(f)
    count :=0
    for {
        line, err := rd.ReadString('\n')
        if err != nil || io.EOF == err {
            break
        }
        tr.Add(line)
        count++
    }


	log.Printf("adding %d words to index took: %v\n", count, time.Since(startTime))

	return
}

