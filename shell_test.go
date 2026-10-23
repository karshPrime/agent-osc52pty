package main

import (
	"bytes"
	"io"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestShell(t *testing.T) {
	var s shell
	stdinR, stdinW := io.Pipe()
	stdoutR, stdoutW := io.Pipe()
	err := s.Open(shellOptions{
		CmdLine: []string{"head", "-1"},
		Stdin:   stdinR,
		Stdout:  stdoutW,
	})
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	defer s.Close()
	// read as a terminal would, while the command runs
	outputs := make(chan []byte, 1)
	go func() {
		output, _ := io.ReadAll(stdoutR)
		outputs <- output
	}()
	_, err = stdinW.Write([]byte("hello world\n"))
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	exitCode, err := s.Wait()
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	if !assert.Equal(t, 0, exitCode) {
		t.FailNow()
	}
	output := <-outputs
	if !assert.Equal(t, "hello world\r\nhello world\r\n", string(output)) {
		t.FailNow()
	}
}

func TestShellWaitDeliversAllOutput(t *testing.T) {
	var s shell
	stdinR, _ := io.Pipe()
	var stdout lockedBuffer
	err := s.Open(shellOptions{
		CmdLine: []string{"seq", "1", "20000"},
		Stdin:   stdinR,
		Stdout:  &stdout,
	})
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	defer s.Close()
	exitCode, err := s.Wait()
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	if !assert.Equal(t, 0, exitCode) {
		t.FailNow()
	}
	output := stdout.String()
	assert.Equal(t, 20000, strings.Count(output, "\r\n"))
	assert.True(t, strings.HasSuffix(output, "\r\n20000\r\n"), "output ends with %q", output[max(0, len(output)-16):])
}

type lockedBuffer struct {
	mu  sync.Mutex
	buf bytes.Buffer
}

func (lb *lockedBuffer) Write(data []byte) (int, error) {
	// a slow terminal, so output is still queued in the pty when the command exits
	time.Sleep(time.Millisecond)
	lb.mu.Lock()
	defer lb.mu.Unlock()
	return lb.buf.Write(data)
}

func (lb *lockedBuffer) Close() error { return nil }

func (lb *lockedBuffer) String() string {
	lb.mu.Lock()
	defer lb.mu.Unlock()
	return lb.buf.String()
}
