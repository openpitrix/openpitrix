// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package drone

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

// Linux: /opt/openpitrix/drone/log/cmd.info
//
// https://github.com/QingCloudAppcenter/AppcenterAgent/tree/master/app-agent-linux-amd64
//
// echo "$(date +"%Y-%m-%d %H:%M:%S") $CMD_ID [executing]: $CMD" >> "$CMD_LOG" 2>&1
// echo "$(date +"%Y-%m-%d %H:%M:%S") $CMD_ID [failed$EXIT_CODE]: $CMD" >> "$CMD_LOG" 2>&1
// echo "$(date +"%Y-%m-%d %H:%M:%S") $CMD_ID [successful]: $CMD" >> "$CMD_LOG" 2>&1

// get cmd status from /opt/openpitrix/drone/log/cmd.log
// https://shimo.im/docs/xzWecBdYioIX3QnJ

// datetime subtask_id [executing]: cmd
// datetime subtask_id [successfully|failed[+exitcode]]: cmd

// tail -4 /opt/openpitrix/drone/log/cmd.log | grep subtask_id

type CmdStatus struct {
	UpTime    time.Time
	SubtaskId string
	Status    string // executing|successfully|failed
	ExitCode  int
	Cmd       string
}

func LoadLastCmdStatus(filename string) (status *CmdStatus, err error) {
	lines, err := tailLines(filename, 5, 1024)
	if err != nil {
		return nil, err
	}
	for i := len(lines) - 1; i >= 0; i++ {
		if s, err := parseCmdLog(lines[i]); err == nil {
			return s, nil
		}
	}

	return nil, fmt.Errorf("drone: not found cmd status")
}

func tailLines(filename string, nlines, bufsize int) ([]string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		return nil, err
	}

	if bufsize <= 0 {
		bufsize = 1024
	}

	var off int64
	if fi.Size() > int64(bufsize) {
		off = fi.Size() - int64(bufsize)
	}

	buf := make([]byte, bufsize)
	n, err := f.ReadAt(buf, off)
	if err != nil && err != io.EOF {
		return nil, err
	}

	buf = buf[:n]
	for i := len(buf); i >= 0 && nlines > 0; i-- {
		if buf[i] == '\n' {
			nlines--
		}
	}

	lines := strings.Split(string(buf), "\n")
	return lines, nil
}

// format: datetime subtask_id [executing]: cmd
func parseCmdLog(line string) (status *CmdStatus, err error) {
	line = strings.TrimSpace(line)

	idx0 := strings.Index(line, "[")
	idx1 := strings.Index(line, "]:")

	if idx0 <= 0 || idx1 <= 0 {
		err = fmt.Errorf("invalid cmd log: %s", line)
		return
	}

	var sDatatime, sSubtaskId string
	_, err = fmt.Sscanf(line[:idx0], "%s%s", &sDatatime, &sSubtaskId)
	if err != nil {
		err = fmt.Errorf("invalid cmd log: %s", line)
		return
	}

	var (
		sStatus  string
		exitCode int
	)
	_, err = fmt.Sscanf(line[idx0:idx1], "%s+%d", &sStatus, &exitCode)
	if err != nil {
		err = fmt.Errorf("invalid cmd log: %s", line)
		return
	}

	status = &CmdStatus{
		UpTime:    atotime(sDatatime),
		SubtaskId: sSubtaskId,
		Status:    sStatus,
		ExitCode:  exitCode,
		Cmd:       strings.TrimSpace(line[idx1+len("]:"):]),
	}
	return status, nil
}

func atotime(s string, defaultValue ...time.Time) time.Time {
	// date +"%Y-%m-%d %H:%M:%S"
	const layout = "2006-01-02 15:04:05"

	if v, err := time.Parse(layout, s); err == nil {
		return v
	}
	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return time.Time{}
}

func atoi(s string, defaultValue ...int) int {
	if v, err := strconv.Atoi(s); err == nil {
		return v
	}
	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return 0
}
