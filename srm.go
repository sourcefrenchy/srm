package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"time"

	"golang.org/x/crypto/blake2b"
)

// DEBUG to be set to enable verbose comments
const DEBUG = false

func getHash(dat []byte) string {
	hasher, _ := blake2b.New256(nil)
	hasher.Write(dat)
	return hex.EncodeToString(hasher.Sum(nil))
}

func delFile(file string, fileSize int64, random1 []byte, random2 []byte) {

	if !DEBUG {
		defer println("\nFile securely deleted.")
	}

	// ********* PASS1 WRITE - NO CHECK
	err := ioutil.WriteFile(file, random1, 0)
	if err != nil {
		log.Fatal(err)
	}

	if DEBUG {
		fmt.Printf("\n==> Pass1: random bytes overwritten: %x...", getHash(random1)[0:10])
	}
	b, err := ioutil.ReadFile(file) // b has type []byte
	if err != nil {
		log.Fatal(err)
	}
	// ********* PASS2 WRITE - NO CHECK
	err = ioutil.WriteFile(file, random2, 0)
	if err != nil {
		log.Fatal(err)
	}
	if DEBUG {
		fmt.Printf("\n==> Pass2: random bytes overwritten: %x...", getHash(random2)[0:10])
	}

	// ********* PASS3 CHECK PASS2
	b, err = ioutil.ReadFile(file) // b has type []byte
	if err != nil {
		log.Fatal(err)
	}

	if DEBUG {
		fmt.Printf("\n==> Pass3: verification:\n\tpass 2 hash: %x... <==>\n\tCurrent file content hash: %x...)", getHash(random2)[0:10], getHash(b)[0:10])
	}

	if DEBUG {
		fmt.Println("\n==> file still on the filesystem (DEBUG mode on)")
	} else {
		err = os.Remove(file)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}

func main() {
	filePath := flag.String("file", "", "exact location path to the file (Required)")
	flag.Parse()

	if *filePath == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	// NSA 130-2 applied to delete securely file on disk:
	// The write head passes over each sector two times (Random, Random).
	// There is one final pass to verify random characters by reading
	if fi, err := os.Stat(*filePath); err == nil {
		size := fi.Size()
		writePass1 := make([]byte, size)
		rand.Seed(time.Now().UTC().UnixNano())
		rand.Read(writePass1)
		writePass2 := make([]byte, size)
		rand.Seed(time.Now().UTC().UnixNano())
		rand.Read(writePass2)
		delFile(*filePath, size, writePass1, writePass2)
	} else {
		fmt.Printf("File %s does not exist\n", *filePath)
	}
}
