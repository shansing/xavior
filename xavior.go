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
)

const BUFFER = 16 * 1024
const DEBUG = true

func printError(err error) {
	if DEBUG {
		fmt.Println(err, err.Error())
	}
}
func printString(str string) {
	if DEBUG {
		fmt.Println(str)
	}
}

func computeKey(str string) []byte {
	bytes := []byte(str)
	hasher := sha512.New512_224()
	hasher.Write(bytes)
	return hasher.Sum(nil)
}

func xorCopy(dst net.Conn, src net.Conn, key []byte) (err error) {
	if cw, ok := dst.(interface{ CloseWrite() error }); ok {
		defer cw.CloseWrite()
	} else {
		return fmt.Errorf("Connection doesn't implement CloseWrite method")
	}
	if cw, ok := src.(interface{ CloseRead() error }); ok {
		defer cw.CloseRead()
	} else {
		return fmt.Errorf("Connection doesn't implement CloseRead method")
	}

	buf := make([]byte, BUFFER)
	xorKeyIndex := 0
	for {
		nr, er := src.Read(buf)
		if nr > 0 {
			goodBuf := buf[:nr]
			for i := 0; i < len(goodBuf); i++ {
				// don't just do `goodBuf[i] ^= key[i%len(key)]` cause nr < buf happens
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
				err = errors.New("ErrShortWrite")
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
	sendPassword := flag.String("sp", "SendPassword", "password when sending data")
	receivePassword := flag.String("rp", "ReceivePassword", "password when receving data")
	flag.Parse()

	sendKey := computeKey(*sendPassword)
	receiveKey := computeKey(*receivePassword)
	if DEBUG {
		println(fmt.Sprintf("sendKey:%x", sendKey))
		println(fmt.Sprintf("receiveKey:%x", receiveKey))
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
		remoteConnection, err := net.Dial("tcp", *remoteHost)
		if err != nil {
			printError(err)
			_ = localConnection.Close()
			continue
		}

		pipe(localConnection, remoteConnection, sendKey, receiveKey)
	}

}
func pipe(localConnection net.Conn, remoteConnection net.Conn, sendKey []byte, receiveKey []byte) {
	errorNum := 0
	onClose := func() {
		if errorNum >= 2 {
			printString("Closing the connections...")
			_ = remoteConnection.Close()
			_ = localConnection.Close()
		}
	}
	go func() {
		err := xorCopy(remoteConnection, localConnection, sendKey)
		if err != nil {
			printError(err)
			errorNum++
			onClose()
		}
	}()
	go func() {
		err := xorCopy(localConnection, remoteConnection, receiveKey)
		if err != nil {
			printError(err)
			errorNum++
			onClose()
		}
	}()
}

//type XOREncryptWriter struct {
//	key []byte
//	w   io.Writer
//}
//
//func NewXOREncryptWriter(key []byte, w io.Writer) *XOREncryptWriter {
//	return &XOREncryptWriter{
//		key: key,
//		w:   w,
//	}
//}
//
//func (xw *XOREncryptWriter) Write(p []byte) (n int, err error) {
//	encrypted := make([]byte, len(p))
//	for i := 0; i < len(p); i++ {
//		encrypted[i] = p[i] ^ xw.key[i%len(xw.key)]
//	}
//	return xw.w.Write(encrypted)
//}
