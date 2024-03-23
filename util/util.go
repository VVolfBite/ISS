package util

import (
	"bytes"
	"encoding/gob"
	"math/rand"
	"time"
	// "fmt"
	// "github.com/kelindar/bitmap"
	// "math/rand"
	// "net"
	// "reflect"
	// "strconv"
	// "time"
)

func RandomPick(n int, f int) []int {
	var randomPick []int
	for i := 0; i < f; i++ {
		var randomID int
		exists := true
		for exists {
			s := rand.NewSource(time.Now().UnixNano())
			r := rand.New(s)
			randomID = r.Intn(n)
			exists = FindIntSlice(randomPick, randomID)
		}
		randomPick = append(randomPick, randomID)
	}
	return randomPick
}

func FindIntSlice(slice []int, val int) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

func SizeOf(v interface{}) int {
	var buffer bytes.Buffer
	enc := gob.NewEncoder(&buffer)
	enc.Encode(v)
	return buffer.Len()
}

func BytesToIdentifier(bytes []byte) Identifier {
    var id Identifier
    copy(id[:], bytes)
    return id
}

func IdentifierToBytes(id Identifier) []byte {
    return id[:]
}

type Signature []byte;
type Identifier [32]byte;
type NodeID int;
