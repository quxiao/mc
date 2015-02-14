package metric

import (
	"bufio"
	"os/exec"
	"strconv"
	"strings"
)

const (
	NAME_IDX = iota
	TOTAL_IDX
	USED_IDX
	AVAILABLE_IDX
	USED_RATE_IDX
	MOUNT_IDX
	METRIC_DISK_FIELD_LEN
)

type MetricDisk struct {
	Name      string
	Mount     string
	Total     uint64
	Available uint64
}

func GetDiskInfo() (diskMap map[string]*MetricDisk, _ error) {
	diskMap = make(map[string]*MetricDisk)

	cmd := exec.Command("df", "-lk")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return diskMap, err
	}
	if err := cmd.Start(); err != nil {
		return diskMap, err
	}
	reader := bufio.NewReader(stdout)
	scanner := bufio.NewScanner(reader)
	lineno := 0

	/*
		Filesystem     1K-blocks     Used Available Use% Mounted on
		/dev/sda1       20158332  2446224  16688108  13% /
		tmpfs           66055224        0  66055224   0% /dev/shm
	*/
	for scanner.Scan() {
		lineno++
		if lineno <= 1 {
			continue
		}
		fields := strings.Fields(scanner.Text())
		if len(fields) < METRIC_DISK_FIELD_LEN {
			continue
		}
		disk := &MetricDisk{}

		var err error
		disk.Name = fields[NAME_IDX]
		disk.Mount = fields[MOUNT_IDX]
		disk.Total, err = strconv.ParseUint(fields[TOTAL_IDX], 10, 64)
		if err != nil {
			continue
		}
		disk.Available, err = strconv.ParseUint(fields[AVAILABLE_IDX], 10, 64)
		if err != nil {
			continue
		}

		diskMap[disk.Name] = disk
	}

	cmd.Wait()

	return diskMap, nil
}
