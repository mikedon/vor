package tests

import "time"

type happyMockOSHandler struct {
	CountConfirm  int
	CountExit     int
	CountExists   int
	CountExec     int
	CountExectime int
}

func (m happyMockOSHandler) Confirm(string) bool {
	m.CountConfirm += 1
	return true
}

func (m happyMockOSHandler) Exit(string) {
	m.CountExit += 1
}

func (m happyMockOSHandler) Exists(str string) (bool, error) {
	m.CountExists += 1
	return true, nil
}

func (m happyMockOSHandler) Exec(string) (string, error) {
	m.CountExec += 1
	return "", nil
}

func (m happyMockOSHandler) ExecutionTimer(time.Time, string) {
	m.CountExectime += 1
}

type sadMockOSHandler struct{}

func (m sadMockOSHandler) Confirm(string) bool {
	return false
}

func (m sadMockOSHandler) Exit(string) {}

func (m sadMockOSHandler) Exists(string) (bool, error) {
	return false, nil
}

func (m sadMockOSHandler) Exec(string) (string, error) {
	return "", *new(error)
}

func (m sadMockOSHandler) ExecutionTimer(time.Time, string) {}
