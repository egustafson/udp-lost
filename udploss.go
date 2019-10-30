package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strings"
)

const (
	ProcUdpFile = "/proc/net/udp"
	ProcUdpCols = 15
)

type ProcUDP struct {
	Slot       uint
	LocalAddr  net.IP
	LocalPort  uint16
	RemoteAddr net.IP
	RemotePort uint16
	Status     uint
	TxQueUsed  uint32 // in Bytes
	RxQueUsed  uint32 // in Bytes
	uid        uint
	Inode      uint
	Drops      uint
}

// readFileNoStat reads the entire file into a []byte.
//  (copied from:  https://github.com/prometheus/procfs/blob/master/internal/util/readfile.go)
//
func readFileNoStat(filename string) ([]byte, error) {
	const MaxBufferSize = 1024 * 128 // assume file is < 128k

	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	reader := io.LimitReader(f, MaxBufferSize)
	return ioutil.ReadAll(reader)
}

type ParseProcUdpError struct{}

func (e *ParseProcUdpError) Error() string {
	return fmt.Sprintf("Proc file format error (%s)", ProcUdpFile)
}

func parseProcUDP(r io.Reader) ([]ProcUDP, error) {
	var (
		scanner = bufio.NewScanner(r)
		result  []ProcUDP
	)

	fields := strings.Fields(scanner.Text())
	if len(fields) != ProcUdpCols {
		return nil, new(ParseProcUdpError)
	}
	for scanner.Scan() {
		fields = strings.Fields(scanner.Text())
		if len(fields) != ProcUdpCols-2 {
			return nil, new(ParseProcUdpError)
		}
		//
		// TODO
		//
		result = append(result, ProcUDP{})
	}
	return result, nil
}

func parseIPPort(s string) (net.IP, uint16, error) {
	return nil, 0, nil
}

func main() {

	log.Print("starting...")
	defer log.Print("exiting.\n")
	for {
		//
		// sleep
		//
		data, err := readFileNoStat(ProcUdpFile)
		if err != nil {
			log.Fatalf("Could not access:  %s\n", ProcUdpFile)
		}
		status, err := parseProcUDP(bytes.NewReader(data))
		if err != nil {
			log.Fatalf("Problem reading proc file: %s\n", err.Error())
		}
		for _, v := range status {
			if v.Drops > 0 {
				log.Printf("%d bytes dropped from %s:%d -> %s:%d\n",
					v.Drops,
					v.RemoteAddr.String(), v.RemotePort,
					v.LocalAddr.String(), v.LocalPort)
			}
		}
	}
}
