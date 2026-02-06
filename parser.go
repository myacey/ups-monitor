package main

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"go.bug.st/serial"
)

func readLine(port serial.Port) (string, error) {
	var buf []byte
	tmp := make([]byte, 1)

	for {
		n, err := port.Read(tmp)
		if err != nil {
			return "", err
		}
		if n == 0 {
			continue
		}

		if tmp[0] == '\r' {
			return string(buf), nil
		}

		buf = append(buf, tmp[0])
	}
}

func parseUPSStatus(line string) (*UPSStatus, error) {
	line = strings.TrimSpace(line)
	line = strings.TrimPrefix(line, "(")

	fields := strings.Fields(line)
	if len(fields) < 8 {
		return nil, errors.New("invalid UPS response length")
	}

	inV, _ := strconv.ParseFloat(fields[0], 64)
	inFaultV, _ := strconv.ParseFloat(fields[1], 64)
	outV, _ := strconv.ParseFloat(fields[2], 64)
	outCurrent, _ := strconv.Atoi(fields[3])
	freq, _ := strconv.ParseFloat(fields[4], 64)
	battV, _ := strconv.ParseFloat(fields[5], 64)
	temp, _ := strconv.ParseFloat(fields[6], 64)
	statusBits := fields[7]

	status := &UPSStatus{
		InputVoltage:      inV,
		InputFaultVoltage: inFaultV,
		OutputVoltage:     outV,
		OutputCurrentPct:  outCurrent,
		InputFrequency:    freq,
		BatteryVoltage:    battV,
		Temperature:       temp,

		UtilityFail:  statusBits[0] == '1',
		BatteryLow:   statusBits[1] == '1',
		BypassActive: statusBits[2] == '1',
		UPSFailed:    statusBits[3] == '1',
		IsStandby:    statusBits[4] == '1',
		TestRunning:  statusBits[5] == '1',
		Shutdown:     statusBits[6] == '1',
		BeeperOn:     statusBits[7] == '1',
	}

	return status, nil
}

func readUPSStatus() {
	mode := &serial.Mode{
		BaudRate: 2400,
		DataBits: 8,
		Parity:   serial.NoParity,
		StopBits: serial.OneStopBit,
	}

	port, err := serial.Open(*portToConn, mode)
	if err != nil {
		log.Fatal("failed to open port:", err)
	}
	defer port.Close()

	for {
		_, err := port.Write([]byte("Q1\r"))
		if err != nil {
			log.Println("failed to write data to port:", err)
			continue
		}

		// wait for response
		time.Sleep(300 * time.Millisecond)

		line, err := readLine(port)
		if err != nil {
			log.Println("read error:", err)
			continue
		}
		log.Println("raw response:", line)

		status, err := parseUPSStatus(line)
		if err != nil {
			log.Println("failed to parse:", err)
			continue
		}

		fmt.Printf(
			"IN: %.1fV | OUT: %.1fV | LOAD: %d%% | BAT: %.2fV | TEMP: %.1f°C | FAIL: %v\n",
			status.InputVoltage,
			status.OutputVoltage,
			status.OutputCurrentPct,
			status.BatteryVoltage,
			status.Temperature,
			status.UtilityFail,
		)

		lastStatus.Store(status)

		time.Sleep(time.Duration(*frequency) * time.Millisecond)
	}
}
