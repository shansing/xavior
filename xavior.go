package main

import (
	"crypto/sha512"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"time"
)

const BUFFER = 16 * 1024
const TIMEOUT = 30 * time.Second
const DEBUG = true

func printError(err error) {
	if DEBUG {
		fmt.Println(err, err.Error())
	}
}

func computeKey(str string) []byte {
	bytes := []byte(str)
	hasher := sha512.New()
	hasher.Write(bytes)
	return hasher.Sum(nil)
}

func xorCopy(dst net.Conn, src net.Conn, key []byte) (err error) {
	defer func(dst net.Conn) {
		_ = dst.Close()
	}(dst)
	defer func(src net.Conn) {
		_ = src.Close()
	}(src)

	buf := make([]byte, BUFFER)
	xorKeyIndex := 0
	for {
		nr, er := src.Read(buf)
		if nr > 0 {
			goodBuf := buf[:nr]
			for i := 0; i < len(goodBuf); i++ {
				// don't just do `goodBuf[i] ^= key[i%len(key)]` cause `nr < len(buf)` happens
				goodBuf[i] ^= key[xorKeyIndex]
				xorKeyIndex = (xorKeyIndex + 1) % len(key)
			}
			nw, ew := dst.Write(goodBuf)
			if nw < 0 || nr < nw {
				nw = 0
				if ew == nil {
					err = errors.New("errInvalidWrite")
				}
			}
			if ew != nil {
				err = ew
				break
			}
			if nr != nw {
				err = io.ErrShortWrite
				break
			}
		}
		if er != nil {
			if er != io.EOF {
				err = er
			}
			break
		}
	}
	return err
}

func main() {
	listenHost := flag.String("l", "127.0.0.1:1234", "listened address")
	remoteHost := flag.String("r", "127.0.0.1:22", "remote address forwarding to")
	sendPassword := flag.String("sp", "SendP@ssw0rd", "password when sending data")
	receivePassword := flag.String("rp", "ReceiveP@ssw0rd", "password when receving data")
	flag.Parse()

	sendKey := computeKey(*sendPassword)
	receiveKey := computeKey(*receivePassword)
	if DEBUG {
		println(fmt.Sprintf("sendKey :%x", sendKey))
		println(fmt.Sprintf("receiveKey: %x", receiveKey))
	}

	log.Printf("Listening on %s ...", *listenHost)
	listen, err := net.Listen("tcp", *listenHost)
	if err != nil {
		printError(err)
		os.Exit(0)
	}
	for {
		localConnection, err := listen.Accept()
		if err != nil {
			printError(err)
			continue
		}

		log.Printf("Establishing connection to %s", *remoteHost)
		remoteConnection, err := net.DialTimeout("tcp", *remoteHost, TIMEOUT)
		if err != nil {
			printError(err)
			_ = localConnection.Close()
			continue
		}

		go func() {
			err := xorCopy(remoteConnection, localConnection, sendKey)
			if err != nil {
				printError(err)
			}
		}()
		go func() {
			err := xorCopy(localConnection, remoteConnection, receiveKey)
			if err != nil {
				printError(err)
			}
		}()
	}

}
