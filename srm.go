package main

import (
	"bufio"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"time"

	"golang.org/x/crypto/blake2b"
)

// DEBUG to be set to enable verbose comments
const DEBUG = false

func delFile(file string, fileSize int64, random1 []byte, random2 []byte) {
	f, err := os.OpenFile(file, os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	if DEBUG {
		fmt.Printf("\n[*] Deleting all %d bytes from file %s:", fileSize, file)
	}

	defer println("File securely deleted.")

	// Create a buffered writer from the file
	bufferedWriter := bufio.NewWriter(f)
	// Write bytes to buffer
	bytesWritten, _ := bufferedWriter.Write(
		random1,
	)
	if DEBUG {
		fmt.Printf("\n==> Pass1: random %d bytes written (hash:%x)", bytesWritten, blake2b.Sum256(random1))
	}
	f.Close()

	f, err = os.OpenFile(file, os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	// Create a buffered writer from the file
	bufferedWriter = bufio.NewWriter(f)
	// Write bytes to buffer
	bytesWritten, _ = bufferedWriter.Write(
		random2,
	)
	if DEBUG {
		fmt.Printf("\n==> Pass2: random %d bytes written (hash:%x)", bytesWritten, blake2b.Sum256(random2))
	}
	f.Close()

	hasher, _ := blake2b.New256(nil)
	f, err = os.Open(file)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	if _, err := io.Copy(hasher, f); err != nil {
		log.Fatal(err)
	}

	hash := hasher.Sum(nil)
	encodedHex := hex.EncodeToString(hash[:])
	if DEBUG {
		fmt.Printf("\n==> Pass3: verified matching pass2 (hash:%s)", encodedHex)
	}
	f.Close()

	err = os.Remove(file)
	if err != nil {
		fmt.Println(err)
		return
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
