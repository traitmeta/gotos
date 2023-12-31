package btc

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/btcsuite/btcd/wire"
)

// ReadMempoolFromPath reads a mempool file from a given path and returns a Mempool type
func ReadMempoolFromPath(path string) (mem Mempool, err error) {
	file, err := os.Open(path)
	if err != nil {
		return mem, fmt.Errorf("read mempool file fail: %s" + err.Error())
	}
	defer file.Close()

	r := bufio.NewReader(file)
	header, err := readFileHeader(r)
	if err != nil {
		return mem, fmt.Errorf("read header fail: %s", err.Error())
	}
	mem.header = header

	var entries []MempoolEntry
	for i := int64(0); i < header.numTx; i++ {
		mementry, err := readMempoolEntry(r)
		if err != nil {
			return mem, fmt.Errorf("read mempoolEntry fail at index %d : %s", i, err.Error())
		}
		entries = append(entries, mementry)
	}
	mem.entries = entries

	var feeDelta []byte
	for {
		remainingBytes, err := r.ReadByte()
		if err == io.EOF {
			break
		} else if err != nil {
			return mem, fmt.Errorf("read feeDelta fail: %s" + err.Error())
		}
		feeDelta = append(feeDelta, remainingBytes)
	}
	mem.mapDeltas = feeDelta

	return
}

func readFileHeader(r *bufio.Reader) (header FileHeader, err error) {
	fileVersion, err := readLEint64(r)
	if err != nil {
		return header, err
	}

	switch fileVersion {
	case MEMPOOL_DUMP_VERSION_NO_XOR_KEY:
		fmt.Println("Leave XOR-key empty")
	case MEMPOOL_DUMP_VERSION:
		_, err := readLEint64(r)
		if err != nil {
			return header, err
		}
	default:
		return header, errors.New("mempool dat file version error")
	}

	numberOfTx, err := readLEint64(r)
	if err != nil {
		return header, err
	}

	header = FileHeader{fileVersion, numberOfTx}
	return
}

func readMempoolEntry(r *bufio.Reader) (mementry MempoolEntry, err error) {
	msgTx := &wire.MsgTx{}
	err = msgTx.Deserialize(r)
	if err != nil {
		return mementry, err
	}

	timestamp, err := readLEint64(r)
	if err != nil {
		return mementry, err
	}

	feeDelta, err := readLEint64(r)
	if err != nil {
		return mementry, err
	}

	mementry = MempoolEntry{msgTx, timestamp, feeDelta}
	return
}

// reads the next 64bit in Little Endian and returns a int64
// get the mempoolEntry's timestamp and feeDelta
func readLEint64(r *bufio.Reader) (res int64, err error) {
	next64Bit := make([]byte, 8)
	_, err = io.ReadFull(r, next64Bit)
	if err != nil {
		return 0, err
	}

	res = int64(binary.LittleEndian.Uint64(next64Bit))
	return
}
