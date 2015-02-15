package mc

import (
	"bufio"
	"time"
	//"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

/*
cat /proc/net/dev
Inter-|   Receive                                                |  Transmit
 face |bytes    packets errs drop fifo frame compressed multicast|bytes    packets errs drop fifo colls carrier compressed
     lo:120309825211 699253920    0    0    0     0          0         0 120309825211 699253920    0    0    0     0       0          0
*/
const (
	BYTES_IN_IDX = iota
	PACKETS_IN_IDX
	ERRS_IN_IDX
	DROP_IN_IDX
	FIFO_IN_IDX
	FRAME_IN_IDX
	COMPRESSED_IN_IDX
	MULTICAST_IN_IDX

	BYTES_OUT_IDX
	PACKETS_OUT_IDX
	ERRS_OUT_IDX
	DROP_OUT_IDX
	FIFO_OUT_IDX
	COLLS_OUT_IDX
	CARRIER_OUT_IDX
	COMPRESSED_OUT_IDX
)

var INTER_NAME_ADDRS_MAP = make(map[string][]string)

func init() {
	//get all interface name and ip address when start up
	interfaces, err := net.Interfaces()
	if err != nil {
		os.Exit(-1)
	}
	for _, inter := range interfaces {
		addrsArr := make([]string, 0)
		addrs, err := inter.Addrs()
		if err != nil {
			continue
		}
		for _, addr := range addrs {
			addrsArr = append(addrsArr, addr.String())
		}
		INTER_NAME_ADDRS_MAP[inter.Name] = addrsArr
	}

}

type MetricInterface struct {
	Timestamp      int64
	Name           string
	Ip             string
	TotalByteIn    uint64
	TotalPacketIn  uint64
	TotalErrorIn   uint64
	TotalDropIn    uint64
	TotalByteOut   uint64
	TotalPacketOut uint64
	TotalErrorOut  uint64
	TotalDropOut   uint64
}

func GetInterfaceInfo() (map[string]*MetricInterface, error) {
	interMap := make(map[string]*MetricInterface)

	file, err := os.Open("/proc/net/dev")
	if err != nil {
		return interMap, err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	scanner := bufio.NewScanner(reader)
	lineno := 0

	for scanner.Scan() {
		lineno++
		if lineno <= 2 {
			continue
		}
		fields := strings.Split(scanner.Text(), ":")
		if len(fields) < 2 {
			continue
		}
		inter := &MetricInterface{Timestamp: time.Now().Unix(), Name: strings.TrimSpace(fields[0])}
		metrics := strings.Fields(fields[1])

		inter.TotalByteIn, err = strconv.ParseUint(metrics[BYTES_IN_IDX], 10, 64)
		if err != nil {
			return interMap, err
		}
		inter.TotalPacketIn, err = strconv.ParseUint(metrics[PACKETS_IN_IDX], 10, 64)
		if err != nil {
			return interMap, err
		}
		inter.TotalErrorIn, err = strconv.ParseUint(metrics[ERRS_IN_IDX], 10, 64)
		if err != nil {
			return interMap, err
		}
		inter.TotalDropIn, err = strconv.ParseUint(metrics[DROP_IN_IDX], 10, 64)
		if err != nil {
			return interMap, err
		}

		inter.TotalByteOut, err = strconv.ParseUint(metrics[BYTES_OUT_IDX], 10, 64)
		if err != nil {
			return interMap, err
		}
		inter.TotalPacketOut, err = strconv.ParseUint(metrics[PACKETS_OUT_IDX], 10, 64)
		if err != nil {
			return interMap, err
		}
		inter.TotalErrorOut, err = strconv.ParseUint(metrics[ERRS_OUT_IDX], 10, 64)
		if err != nil {
			return interMap, err
		}
		inter.TotalDropOut, err = strconv.ParseUint(metrics[DROP_OUT_IDX], 10, 64)
		if err != nil {
			return interMap, err
		}

		interMap[inter.Name] = inter
	}

	return interMap, nil
}
