package mc

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

const (
	PROTOCAL_NAME_IDX = iota
	R_TO_ALG_IDX
	R_TO_MIN_IDX
	R_TO_MAX_IDX
	MAX_CONN_IDX
	ACTIVE_OPENS_IDX
	PASSIVE_OPENS_IDX
	ATTEMPT_FAILS_IDX
	ESTAB_RESETS_IDX
	CURR_ESTAB_IDX
	IN_SEGS_IDX
	OUT_SEGS_IDX
	RETRANS_SEGS_IDX
	IN_ERRS_IDX
	OUT_ERRS_IDX
	NET_SNMP_FIELDS_LEN
)

type MetricConn struct {
	Timestamp   int64
	Established uint64
	TimeWait    uint64
	// TODO
}

func GetConnInfo() (conn MetricConn, _ error) {
	{
		/*
			cat /proc/net/snmp
			Ip: Forwarding DefaultTTL InReceives InHdrErrors InAddrErrors ForwDatagrams InUnknownProtos InDiscards InDelivers OutRequests OutDiscards OutNoRoutes ReasmTimeout ReasmReqds ReasmOKs ReasmFails FragOKs FragFails FragCreates
			Ip: 1 64 1984201988 0 200 0 52 0 1983901334 1936944569 43427 0 0 34 17 0 0 0 0
			Icmp: InMsgs InErrors InDestUnreachs InTimeExcds InParmProbs InSrcQuenchs InRedirects InEchos InEchoReps InTimestamps InTimestampReps InAddrMasks InAddrMaskReps OutMsgs OutErrors OutDestUnreachs OutTimeExcds OutParmProbs OutSrcQuenchs OutRedirects OutEchos OutEchoReps OutTimestamps OutTimestampReps OutAddrMasks OutAddrMaskReps
			Icmp: 10579557 18868 979611 19436 0 6 181 603629 8976614 6 11 0 0 15270630 0 5444760 0 0 0 0 9204369 603609 0 6 0 0
			IcmpMsg: InType0 InType3 InType4 InType5 InType8 InType11 InType13 InType14 OutType0 OutType3 OutType8 OutType14 OutType69
			IcmpMsg: 8976614 979611 6 181 603629 19436 6 11 603609 5444760 9204369 6 17886
			Tcp: RtoAlgorithm RtoMin RtoMax MaxConn ActiveOpens PassiveOpens AttemptFails EstabResets CurrEstab InSegs OutSegs RetransSegs InErrs OutRsts
			Tcp: 1 200 120000 -1 108778201 78477668 18427452 354910 139 1860407343 1805237086 3218436 7974 20251942
			Udp: InDatagrams NoPorts InErrors OutDatagrams RcvbufErrors SndbufErrors
			Udp: 98311688 14215298 70 113741916 0 0
			UdpLite: InDatagrams NoPorts InErrors OutDatagrams RcvbufErrors SndbufErrors
			UdpLite: 0 0 0 0 0 0
		*/
		file, err := os.Open("/proc/net/snmp")
		if err != nil {
			return conn, err
		}
		defer file.Close()

		reader := bufio.NewReader(file)
		scanner := bufio.NewScanner(reader)
		tcp_flag := 0

		for scanner.Scan() {
			if !strings.HasPrefix(scanner.Text(), "Tcp:") {
				continue
			}
			tcp_flag++
			if tcp_flag <= 1 {
				continue
			}
			fields := strings.Fields(scanner.Text())
			if len(fields) < NET_SNMP_FIELDS_LEN {
				continue
			}
			conn.Established, err = strconv.ParseUint(fields[CURR_ESTAB_IDX], 10, 64)
			if err != nil {
				return conn, err
			}
			break
		}
	}

	{
		/*
			cat /proc/net/sockstat
			sockets: used 25460
			TCP: inuse 20 orphan 0 tw 69 alloc 473 mem 2
			UDP: inuse 2 mem 1
			UDPLITE: inuse 0
			RAW: inuse 0
			FRAG: inuse 0 memory 0
		*/
		file, err := os.Open("/proc/net/sockstat")
		if err != nil {
			return conn, err
		}
		defer file.Close()

		reader := bufio.NewReader(file)
		scanner := bufio.NewScanner(reader)

		for scanner.Scan() {
			if !strings.HasPrefix(scanner.Text(), "TCP:") {
				continue
			}
			fields := strings.Fields(scanner.Text())
			fieldsLen := len(fields)
			if fieldsLen%2 == 0 || fieldsLen <= 1 {
				continue
			}
			for i := 1; i+1 < fieldsLen; i += 2 {
				if fields[i] == "tw" {
					conn.TimeWait, err = strconv.ParseUint(fields[i+1], 10, 64)
					if err != nil {
						return conn, err
					}
					break
				}
			}
			break
		}
	}

	return conn, nil
}
