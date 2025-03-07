package main

import (
	"crypto/tls"
	"fmt"
	"net/rpc"
	"sync"
	"time"

	"rc/shared"
)

type Controller struct {
	workers   []string
	clients   map[string]*rpc.Client
	jobCounts map[string]int
	mu        sync.Mutex
	startTime time.Time
	allDown   bool
}

// Initialize and register workers
func RegisterWorkers() *Controller {
	return &Controller{
		workers:   []string{"192.168.18.89:5001", "192.168.18.89:5002", "192.168.18.89:5003"},
		clients:   make(map[string]*rpc.Client),
		jobCounts: make(map[string]int),
	}
}

// Establish connections with all workers
func (c *Controller) RpcRegister() {
	for _, worker := range c.workers {
		c.connectWorker(worker)
	}
}

// Connect to a worker securely
func (c *Controller) connectWorker(worker string) {
	cert, err := tls.LoadX509KeyPair("server.crt", "server.key")
	if err != nil {
		fmt.Println("Failed to load certificate:", err)
		return
	}

	tlsConfig := &tls.Config{Certificates: []tls.Certificate{cert}, InsecureSkipVerify: true}
	conn, err := tls.Dial("tcp", worker, tlsConfig)
	if err != nil {
		fmt.Printf("Failed to connect to worker %s: %v\n", worker, err)
		return
	}

	rpcClient := rpc.NewClient(conn)

	c.mu.Lock()
	c.clients[worker] = rpcClient
	c.jobCounts[worker] = 0
	if c.allDown {
		c.allDown = false
		c.startTime = time.Time{}
	}
	c.mu.Unlock()
}

// Get the worker with the least number of jobs
func (c *Controller) getLeastBusyWorker() (string, *rpc.Client) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if len(c.clients) == 0 {
		if c.startTime.IsZero() {
			c.startTime = time.Now()
			c.allDown = true
		}
		return "", nil
	}

	c.allDown = false
	c.startTime = time.Time{}

	var leastBusyWorker string
	minJobs := 1000000

	for worker, count := range c.jobCounts {
		if count < minJobs {
			minJobs = count
			leastBusyWorker = worker
		}
	}

	if leastBusyWorker == "" {
		return "", nil
	}
	return leastBusyWorker, c.clients[leastBusyWorker]
}

// Process matrix operation request with automatic retry
func (c *Controller) Process(req shared.StructureReq, res *shared.StructureResponse) error {
	for {
		worker, client := c.getLeastBusyWorker()

		if client == nil {
			if c.allDown && time.Since(c.startTime) > 1*time.Minute {
				return fmt.Errorf("All workers are down for more than 1 minute. No nodes available.")
			}

			fmt.Println("No available workers. Retrying after 2 seconds...")
			time.Sleep(2 * time.Second)
			c.RpcRegister()
			continue
		}

		c.mu.Lock()
		c.jobCounts[worker]++
		c.mu.Unlock()

		err := client.Call("Worker.MatrixOp", req, res)

		c.mu.Lock()
		c.jobCounts[worker]--
		c.mu.Unlock()

		if err == nil {
			return nil
		}

		fmt.Printf("Worker %s failed. Retrying...\n", worker)
		c.mu.Lock()
		delete(c.clients, worker)
		c.mu.Unlock()

		c.connectWorker(worker)

		time.Sleep(2 * time.Second)
	}
}

// Start the Controller RPC server
func main() {
	controller := RegisterWorkers()
	controller.RpcRegister()

	err := rpc.Register(controller)
	if err != nil {
		fmt.Println("Error registering RPC:", err)
		return
	}

	cert, err := tls.LoadX509KeyPair("server.crt", "server.key")
	if err != nil {
		fmt.Printf("Failed to load certificate: %v", err)
		return
	}

	config := &tls.Config{Certificates: []tls.Certificate{cert}}

	listener, err := tls.Listen("tcp", ":5000", config)
	if err != nil {
		fmt.Printf("Error starting secure controller: %v", err)
		return
	}

	fmt.Println("Secure Controller running on port 5000...")
	for {
		conn, _ := listener.Accept()
		go rpc.ServeConn(conn)
	}
}