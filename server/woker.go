package main

import("fmt")
import("net/rpc")
import("os")
import("time")
import("rc/shared")
import("crypto/tls")



type Worker struct{}

// RPC method for matrix operations
func (w *Worker) MatrixOp(param shared.StructureReq, res *shared.StructureResponse) error {
	fmt.Println("Received request, processing...")

	// Simulating processing delay
	time.Sleep(60 * time.Second)

	switch param.OperationType {
	case "add":
		res.Res = addMatrices(param.Mat1, param.Mat2)
		res.Worker=os.Getenv("PORT")
	case "transpose":
		res.Res = subMatrices(param.Mat1, param.Mat2)
		res.Worker=os.Getenv("PORT")
	case "multiply":
		res.Res = mulMatrices(param.Mat1, param.Mat2)
		res.Worker=os.Getenv("PORT")
	default:
		return fmt.Errorf(param.OperationType + " is not a valid operation")
		res.Worker=os.Getenv("PORT")
	}
	return nil
}

// Helper functions
func addMatrices(a, b [][]float64) [][]float64 {
	rows, cols := len(a), len(a[0])
	result := make([][]float64, rows)
	for i := range result {
		result[i] = make([]float64, cols)
		for j := range result[i] {
			result[i][j] = a[i][j] + b[i][j]
		}
	}
	return result
}

func subMatrices(a, b [][]float64) [][]float64 {
	// Transpose matrix a
	rows, cols := len(a[0]), len(a) // Swap dimensions for transpose
	transposedA := make([][]float64, rows)
	for i := range transposedA {
		transposedA[i] = make([]float64, cols)
		for j := range transposedA[i] {
			transposedA[i][j] = a[j][i]
		}
	}
	return transposedA
}

func mulMatrices(a, b [][]float64) [][]float64 {
	rows, cols, common := len(a), len(b[0]), len(b)
	result := make([][]float64, rows)
	for i := range result {
		result[i] = make([]float64, cols)
		for j := range result[i] {
			for k := 0; k < common; k++ {
				result[i][j] += a[i][k] * b[k][j]
			}
		}
	}
	return result
}

// Start the worker RPC server
func main() {
	worker := new(Worker)
	rpc.Register(worker)

	port := os.Getenv("PORT")

	// Load TLS certificate
	cert, err := tls.LoadX509KeyPair("server.crt", "server.key")
	if err != nil {
		fmt.Println("Failed to load certificate:", err)
		return
	}

	config := &tls.Config{Certificates: []tls.Certificate{cert}}

	listener, err := tls.Listen("tcp", fmt.Sprintf(":%s", port), config)
	if err != nil {
		fmt.Println("Error starting worker:", err)
		return
	}

	fmt.Printf("Secure Worker running on port %s...\n", port)
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Connection error:", err)
			continue
		}
		go rpc.ServeConn(conn)
	}

}