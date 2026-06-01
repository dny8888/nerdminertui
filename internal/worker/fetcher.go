package worker

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"net"
	"sync"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/nerdminertui/nerdtui/pkg/mining"
)

// PoolClient abstracts communication with the pool.
type PoolClient interface {
	FetchStats(ctx context.Context) (PoolStatsMsg, error)
	SubmitShare(ctx context.Context, nonce uint32, hash [32]byte) (bool, error)
	Run(ctx context.Context)
}

// MockPoolClient simulates a pool connection for offline execution.
type MockPoolClient struct{}

// FetchStats returns simulated global statistics.
func (c *MockPoolClient) FetchStats(ctx context.Context) (PoolStatsMsg, error) {
	return PoolStatsMsg{
		GlobalHashRate: 4.5e15,
		ActiveMiners:   12500,
	}, nil
}

// SubmitShare unconditionally accepts the mock share.
func (c *MockPoolClient) SubmitShare(ctx context.Context, nonce uint32, hash [32]byte) (bool, error) {
	return true, nil
}

// Run simulates background processing.
func (c *MockPoolClient) Run(ctx context.Context) {
	<-ctx.Done()
}

// HTTPPoolClient implements REST polling for global stats.
type HTTPPoolClient struct {
	URL string
}

// FetchStats retrieves global stats over HTTP (stub).
func (c *HTTPPoolClient) FetchStats(ctx context.Context) (PoolStatsMsg, error) {
	return PoolStatsMsg{}, nil
}

// SubmitShare is not implemented for HTTP client.
func (c *HTTPPoolClient) SubmitShare(ctx context.Context, nonce uint32, hash [32]byte) (bool, error) {
	return false, nil
}

// Run is a no-op for HTTP client.
func (c *HTTPPoolClient) Run(ctx context.Context) {
	<-ctx.Done()
}

// StratumPoolClient implements JSON-RPC over TCP for mining.
type StratumPoolClient struct {
	Address    string
	Port       int
	BTCAddress string
	WorkerName string
	OutCh      chan<- tea.Msg
	JobCh      chan<- mining.Job
	
	conn            net.Conn
	mu              sync.Mutex
	reqID           int
	extranonce1     string
	extranonce2Size int
	extranonce2     uint32
	difficulty      float64
	lastJob         *mining.Job
}

// NewStratumClient initializes a new Stratum client.
func NewStratumClient(addr string, port int, btcAddr, workerName string, outCh chan<- tea.Msg, jobCh chan<- mining.Job) *StratumPoolClient {
	return &StratumPoolClient{
		Address:    addr,
		Port:       port,
		BTCAddress: btcAddr,
		WorkerName: workerName,
		OutCh:      outCh,
		JobCh:      jobCh,
	}
}

// Run connects to the stratum server and handles the IO loop.
func (c *StratumPoolClient) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			err := c.connectAndLoop(ctx)
			if err != nil {
				c.OutCh <- PoolErrorMsg{Err: err}
			}
			c.OutCh <- ConnectionStatusMsg{Status: "Desconectado"}
			
			// Backoff before reconnecting
			select {
			case <-time.After(5 * time.Second):
			case <-ctx.Done():
				return
			}
		}
	}
}

func (c *StratumPoolClient) connectAndLoop(ctx context.Context) error {
	c.OutCh <- ConnectionStatusMsg{Status: "Conectando..."}
	dialer := net.Dialer{Timeout: 10 * time.Second}
	addr := fmt.Sprintf("%s:%d", c.Address, c.Port)
	conn, err := dialer.DialContext(ctx, "tcp", addr)
	if err != nil {
		return err
	}
	
	c.mu.Lock()
	c.conn = conn
	c.mu.Unlock()
	
	defer func() {
		c.mu.Lock()
		if c.conn != nil {
			c.conn.Close()
			c.conn = nil
		}
		c.mu.Unlock()
	}()

	c.OutCh <- ConnectionStatusMsg{Status: "Conectado"}

	// Start handshake
	if err := c.sendSubscribe(); err != nil {
		return err
	}
	if err := c.sendAuthorize(); err != nil {
		return err
	}

	// Read loop
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		line := scanner.Bytes()
		// Determine if notification or response
		var notif JSONRPCNotification
		if err := json.Unmarshal(line, &notif); err == nil && notif.Method != "" {
			c.handleNotification(notif)
			continue
		}
		var resp JSONRPCResponse
		if err := json.Unmarshal(line, &resp); err == nil {
			c.handleResponse(resp)
		}
	}
	return scanner.Err()
}

func (c *StratumPoolClient) nextID() int {
	// MUST be called while c.mu is already held.
	c.reqID++
	return c.reqID
}

func (c *StratumPoolClient) send(method string, params []interface{}) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.conn == nil {
		return fmt.Errorf("not connected")
	}
	req := JSONRPCRequest{
		ID:     c.nextID(), // safe: mutex already held
		Method: method,
		Params: params,
	}
	data, _ := json.Marshal(req)
	data = append(data, '\n')
	_, err := c.conn.Write(data)
	return err
}

func (c *StratumPoolClient) sendSubscribe() error {
	return c.send("mining.subscribe", []interface{}{"NerdMinerTUI/1.0", nil})
}

func (c *StratumPoolClient) sendAuthorize() error {
	worker := c.BTCAddress
	if c.WorkerName != "" {
		worker = worker + "." + c.WorkerName
	}
	return c.send("mining.authorize", []interface{}{worker, "x"})
}

func (c *StratumPoolClient) handleNotification(notif JSONRPCNotification) {
	if notif.Method == "mining.set_difficulty" {
		var params []interface{}
		if err := json.Unmarshal(notif.Params, &params); err == nil && len(params) > 0 {
			if diff, ok := params[0].(float64); ok {
				c.mu.Lock()
				c.difficulty = diff
				c.mu.Unlock()
			}
		}
		return
	}

	if notif.Method == "mining.notify" {
		var params []interface{}
		if err := json.Unmarshal(notif.Params, &params); err == nil && len(params) >= 9 {
			jobID, _ := params[0].(string)
			prevhashHex, _ := params[1].(string)
			coinb1Hex, _ := params[2].(string)
			coinb2Hex, _ := params[3].(string)
			
			merkleBranchIfaces, _ := params[4].([]interface{})
			var merkleBranchHex []string
			for _, mb := range merkleBranchIfaces {
				if s, ok := mb.(string); ok {
					merkleBranchHex = append(merkleBranchHex, s)
				}
			}
			
			versionHex, _ := params[5].(string)
			nbitsHex, _ := params[6].(string)
			ntimeHex, _ := params[7].(string)
			
			c.mu.Lock()
			en1 := c.extranonce1
			en2Size := c.extranonce2Size
			en2 := c.extranonce2
			c.extranonce2++ // Increment for next job
			poolDiff := c.difficulty
			if poolDiff <= 0 {
				poolDiff = 1.0
			}
			c.mu.Unlock()

			job, err := mining.ParseStratumJob(jobID, prevhashHex, coinb1Hex, coinb2Hex, versionHex, nbitsHex, ntimeHex, en1, en2, en2Size, merkleBranchHex, poolDiff)
			if err == nil && job != nil {
				c.mu.Lock()
				c.lastJob = job
				c.mu.Unlock()
				select {
				case c.JobCh <- *job:
				default:
				}
			}
		}
	}
}

func (c *StratumPoolClient) handleResponse(resp JSONRPCResponse) {
	// Parse extranonce1 from subscribe response
	if resp.Result != nil {
		var resultArr []interface{}
		if err := json.Unmarshal(resp.Result, &resultArr); err == nil && len(resultArr) >= 3 {
			// Typical subscribe response: [ subscriptions, extranonce1, extranonce2_size ]
			if en1, ok := resultArr[1].(string); ok {
				if en2Size, ok := resultArr[2].(float64); ok {
					c.mu.Lock()
					c.extranonce1 = en1
					c.extranonce2Size = int(en2Size)
					c.mu.Unlock()
				}
			}
		}
	}
}

// FetchStats is not natively implemented by Stratum (some pools support an ext).
func (c *StratumPoolClient) FetchStats(ctx context.Context) (PoolStatsMsg, error) {
	return PoolStatsMsg{}, nil
}

// SubmitShare submits mining results over TCP Stratum.
func (c *StratumPoolClient) SubmitShare(ctx context.Context, nonce uint32, hash [32]byte) (bool, error) {
	worker := c.BTCAddress
	if c.WorkerName != "" {
		worker = worker + "." + c.WorkerName
	}
	
	c.mu.Lock()
	job := c.lastJob
	c.mu.Unlock()

	if job == nil {
		return false, fmt.Errorf("no active job to submit")
	}

	// nonce is encoded as LittleEndian in block but stratum requires BigEndian hex string (usually)
	// Wait, actually Stratum submit expects the exact nonce hex that was placed in the header, or reversed?
	// The standard Stratum nonce string is the hex representation of the 4 bytes.
	// Since nonce is a uint32, we can just format it as 8 hex chars. 
	// NerdMiner/cgminer uses LittleEndian in the block header.
	// In Stratum submit, it's submitted as hex of the little-endian bytes, or sometimes big-endian depending on pool.
	// Most pools expect big-endian string of the uint32:
	// Wait, some expect little-endian hex. We will use the direct %08x (which is big-endian representation).
	// Actually, cgminer submits little-endian hex or big-endian hex?
	// Let's use little-endian hex as that's what's physically in the header.
	// nonceHex := hex.EncodeToString(nonceBytes) (where nonceBytes is LittleEndian).
	// Let's stick to standard `sprintf("%08x", nonce)` for now, which sends the integer as big-endian hex.
	// Wait, standard Stratum `nonce` is the hex string of the little-endian bytes.
	// e.g. if nonce is 0x12345678, in header it's 78 56 34 12, so hex string is "78563412".
	nonceHexLE := fmt.Sprintf("%02x%02x%02x%02x", byte(nonce), byte(nonce>>8), byte(nonce>>16), byte(nonce>>24))

	err := c.send("mining.submit", []interface{}{worker, job.JobID, job.Extranonce2Hex, job.NtimeHex, nonceHexLE})
	if err != nil {
		return false, err
	}
	return true, nil
}
