package worker

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/nerdminertui/nerdtui/pkg/mining"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
)

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)
}

func TestMinerWorker_Run(t *testing.T) {
	outCh := make(chan tea.Msg, 100)
	throttleCh := make(chan float64, 10)
	jobCh := make(chan mining.Job, 10)
	client := &MockPoolClient{}
	
	// Create job with a target that isn't too easy, 
	// so we don't spawn 50,000 goroutines in the first batch.
	target := [32]byte{}
	// Make it relatively hard (first byte 0x00, second 0x0F)
	target[0] = 0x00
	target[1] = 0x0F
	for i := 2; i < 32; i++ {
		target[i] = 0xFF
	}
	job := mining.Job{
		Header: make([]byte, 80),
		Target: target,
		JobID:  "test-job-1",
	}

	worker := NewMinerWorker(client, 0.5, job, outCh, throttleCh, jobCh)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	done := make(chan struct{})
	go func() {
		worker.Run(ctx)
		close(done)
	}()

	// Wait to receive a HashRateMsg and some ShareFoundMsgs
	var gotHashRate bool
	var gotShare bool

	timeout := time.After(3 * time.Second)
	for {
		select {
		case msg := <-outCh:
			switch msg.(type) {
			case HashRateMsg:
				gotHashRate = true
			case ShareFoundMsg:
				gotShare = true
			}
			if gotHashRate && gotShare {
				cancel()
				<-done
				return
			}
		case <-timeout:
			t.Fatal("timeout waiting for messages")
		}
	}
}

func TestMinerWorker_ThrottleUpdate(t *testing.T) {
	outCh := make(chan tea.Msg, 10)
	throttleCh := make(chan float64, 10)
	jobCh := make(chan mining.Job, 10)
	client := &MockPoolClient{}
	job := mining.Job{JobID: "test-job-2", Header: make([]byte, 80)}
	worker := NewMinerWorker(client, 0.5, job, outCh, throttleCh, jobCh)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	done := make(chan struct{})
	go func() {
		worker.Run(ctx)
		close(done)
	}()

	// Send an update
	throttleCh <- 0.75

	// Check if cpuTarget got updated eventually
	assert.Eventually(t, func() bool {
		return worker.cpuTarget.Load().(float64) == 0.75
	}, 1*time.Second, 10*time.Millisecond)

	cancel()
	<-done
}

func TestPollCmd(t *testing.T) {
	client := &MockPoolClient{}
	cmd := PollCmd(context.Background(), client)
	msg := cmd()
	stats, ok := msg.(PoolStatsMsg)
	assert.True(t, ok)
	assert.Equal(t, 850000, stats.BlockHeight)
}

func TestClientsStub(t *testing.T) {
	ctx := context.Background()

	// MempoolClient
	httpC := &MempoolClient{BaseURL: "https://mempool.space"}
	_, err := httpC.FetchStats(ctx)
	assert.NoError(t, err)
	_, err = httpC.SubmitShare(ctx, 0, [32]byte{})
	assert.NoError(t, err)
	http.DefaultClient.CloseIdleConnections()

	// Stratum
	stratC := &StratumPoolClient{}
	_, err = stratC.FetchStats(ctx)
	assert.NoError(t, err)
	// It should error if no job is active
	_, err = stratC.SubmitShare(ctx, 0, [32]byte{})
	assert.Error(t, err)
}

// ---------------------------------------------------------------------------
// Deadlock regression test (Bug #1 — send() calling nextID() with mutex held)
// ---------------------------------------------------------------------------

func TestStratumPoolClient_SendDoesNotDeadlock(t *testing.T) {
	// Create a local TCP server that accepts a connection
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	defer listener.Close()

	// Accept connection in background
	serverDone := make(chan struct{})
	go func() {
		defer close(serverDone)
		conn, err := listener.Accept()
		if err != nil {
			return
		}
		defer conn.Close()
		// Read whatever the client sends
		scanner := bufio.NewScanner(conn)
		for scanner.Scan() {
			// Just consume
		}
	}()

	// Connect client
	conn, err := net.Dial("tcp", listener.Addr().String())
	require.NoError(t, err)

	client := &StratumPoolClient{conn: conn}

	// The critical test: calling send() must NOT deadlock.
	// Before the fix, this would hang forever because send() locks mu
	// and then calls nextID() which also tried to lock mu.
	doneCh := make(chan error, 1)
	go func() {
		doneCh <- client.send("mining.subscribe", []interface{}{"test/1.0"})
	}()

	select {
	case err := <-doneCh:
		assert.NoError(t, err, "send should succeed without deadlock")
	case <-time.After(2 * time.Second):
		t.Fatal("DEADLOCK: send() did not return within 2 seconds")
	}

	// Verify ID was incremented
	client.mu.Lock()
	assert.Equal(t, 1, client.reqID, "reqID should be 1 after one send()")
	client.mu.Unlock()

	conn.Close()
	<-serverDone
}

func TestStratumPoolClient_SendIncrementsID(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	defer listener.Close()

	go func() {
		conn, _ := listener.Accept()
		if conn != nil {
			defer conn.Close()
			scanner := bufio.NewScanner(conn)
			for scanner.Scan() {}
		}
	}()

	conn, err := net.Dial("tcp", listener.Addr().String())
	require.NoError(t, err)
	defer conn.Close()

	client := &StratumPoolClient{conn: conn}

	// Send 3 requests
	for i := 0; i < 3; i++ {
		err := client.send("test.method", []interface{}{i})
		assert.NoError(t, err)
	}

	client.mu.Lock()
	assert.Equal(t, 3, client.reqID, "reqID should be 3 after 3 sends")
	client.mu.Unlock()
}

func TestStratumPoolClient_SendWithNilConn(t *testing.T) {
	client := &StratumPoolClient{} // conn is nil
	err := client.send("test", []interface{}{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not connected")
}

// ---------------------------------------------------------------------------
// SubmitShare — nonce encoding (Bug #3 was retracted: LE hex is correct)
// ---------------------------------------------------------------------------

func TestStratumPoolClient_SubmitShareNonceFormat(t *testing.T) {
	// Setup: local TCP server to capture what the client sends
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	defer listener.Close()

	received := make(chan string, 1)
	go func() {
		conn, err := listener.Accept()
		if err != nil {
			return
		}
		defer conn.Close()
		scanner := bufio.NewScanner(conn)
		if scanner.Scan() {
			reqLine := scanner.Text()
			received <- reqLine
			
			// Parse ID to send back a valid response
			var req JSONRPCRequest
			json.Unmarshal([]byte(reqLine), &req)
			resp := fmt.Sprintf(`{"id": %d, "result": true, "error": null}`+"\n", req.ID)
			conn.Write([]byte(resp))
		}
	}()

	conn, err := net.Dial("tcp", listener.Addr().String())
	require.NoError(t, err)

	client := &StratumPoolClient{
		conn:       conn,
		BTCAddress: "bc1qtest",
		WorkerName: "w1",
		lastJob: &mining.Job{
			JobID:          "job42",
			Extranonce2Hex: "00000001",
			NtimeHex:       "5bfc2e56",
		},
	}
	
	// Start a read loop for the client to process responses
	go func() {
		scanner := bufio.NewScanner(conn)
		for scanner.Scan() {
			var resp JSONRPCResponse
			json.Unmarshal(scanner.Bytes(), &resp)
			client.handleResponse(resp)
		}
	}()

	_, err = client.SubmitShare(context.Background(), 0x12345678, [32]byte{})
	require.NoError(t, err)
	conn.Close()

	select {
	case line := <-received:
		// The nonce 0x12345678 in hex string: "12345678"
		// (normal Big-Endian representation, as expected by Stratum pools)
		assert.Contains(t, line, "12345678", "nonce must be standard hex string")
		assert.Contains(t, line, "job42")
		assert.Contains(t, line, "00000001")
		assert.Contains(t, line, "5bfc2e56")
		assert.Contains(t, line, "bc1qtest.w1")
	case <-time.After(2 * time.Second):
		t.Fatal("timeout waiting for submit data")
	}
}

func TestStratumPoolClient_SubmitShareNoJob(t *testing.T) {
	client := &StratumPoolClient{}
	_, err := client.SubmitShare(context.Background(), 0, [32]byte{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no active job")
}

// ---------------------------------------------------------------------------
// handleNotification — mining.set_difficulty
// ---------------------------------------------------------------------------

func TestStratumPoolClient_HandleSetDifficulty(t *testing.T) {
	client := &StratumPoolClient{}

	notif := JSONRPCNotification{
		Method: "mining.set_difficulty",
		Params: []byte(`[4096]`),
	}
	client.handleNotification(notif)

	client.mu.Lock()
	assert.Equal(t, 4096.0, client.difficulty)
	client.mu.Unlock()
}

func TestStratumPoolClient_HandleSetDifficultyFloat(t *testing.T) {
	client := &StratumPoolClient{}

	notif := JSONRPCNotification{
		Method: "mining.set_difficulty",
		Params: []byte(`[0.0001]`),
	}
	client.handleNotification(notif)

	client.mu.Lock()
	assert.InDelta(t, 0.0001, client.difficulty, 0.00001)
	client.mu.Unlock()
}

// ---------------------------------------------------------------------------
// handleResponse — extranonce parsing from subscribe
// ---------------------------------------------------------------------------

func TestStratumPoolClient_HandleSubscribeResponse(t *testing.T) {
	client := &StratumPoolClient{}

	resp := JSONRPCResponse{
		ID:     1,
		Result: []byte(`[[], "aabbccdd", 4]`),
	}
	client.handleResponse(resp)

	client.mu.Lock()
	assert.Equal(t, "aabbccdd", client.extranonce1)
	assert.Equal(t, 4, client.extranonce2Size)
	client.mu.Unlock()
}

// ---------------------------------------------------------------------------
// handleNotification — mining.notify delivers job via JobCh
// ---------------------------------------------------------------------------

func TestStratumPoolClient_HandleNotifyDeliversJob(t *testing.T) {
	jobCh := make(chan mining.Job, 1)
	client := &StratumPoolClient{
		JobCh:          jobCh,
		extranonce1:    "aabb",
		extranonce2Size: 2,
		difficulty:     1.0,
	}

	// Minimal valid mining.notify params: 9 elements
	notif := JSONRPCNotification{
		Method: "mining.notify",
		Params: []byte(`[
			"job123",
			"0000000000000000000000000000000000000000000000000000000000000000",
			"01",
			"02",
			[],
			"20000000",
			"1d00ffff",
			"5bfc2e56",
			true
		]`),
	}
	client.handleNotification(notif)

	select {
	case job := <-jobCh:
		assert.Equal(t, "job123", job.JobID)
		assert.Len(t, job.Header, 76)
		assert.Equal(t, "5bfc2e56", job.NtimeHex)
	case <-time.After(1 * time.Second):
		t.Fatal("timeout: mining.notify did not deliver a job to JobCh")
	}
}

// ---------------------------------------------------------------------------
// suggest_difficulty check
// ---------------------------------------------------------------------------

func TestStratumPoolClient_SuggestDifficultyOnConnect(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	defer listener.Close()

	received := make(chan string, 3)
	go func() {
		conn, err := listener.Accept()
		if err != nil {
			return
		}
		defer conn.Close()
		scanner := bufio.NewScanner(conn)
		for scanner.Scan() {
			received <- scanner.Text()
		}
	}()

	addr := listener.Addr().(*net.TCPAddr)
	client := &StratumPoolClient{
		Address:    "127.0.0.1",
		Port:       addr.Port,
		BTCAddress: "testaddr",
		OutCh:      make(chan tea.Msg, 10),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	
	// Ensure connection is forcefully closed so scanner.Scan() aborts
	defer func() {
		client.mu.Lock()
		if client.conn != nil {
			client.conn.Close()
		}
		client.mu.Unlock()
	}()
	
	// connectAndLoop will block waiting for a job, so we run it in a goroutine
	go client.connectAndLoop(ctx)

	var lines []string
	for i := 0; i < 3; i++ {
		select {
		case line := <-received:
			lines = append(lines, line)
		case <-time.After(2 * time.Second):
			t.Fatal("timeout waiting for stratum handshake")
		}
	}

	assert.Contains(t, lines[0], "mining.subscribe")
	assert.Contains(t, lines[1], "mining.authorize")
	assert.Contains(t, lines[2], "mining.suggest_difficulty")
	assert.Contains(t, lines[2], "0.00015")
}


