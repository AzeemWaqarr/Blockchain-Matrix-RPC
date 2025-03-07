package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"log"
	"net/rpc"
	"os"
	"rc/shared"
	"strconv"
	"strings"
)

// Get valid integer input with error handling
func getValidIntegerInput(prompt string) int {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print(prompt)
		scanner.Scan()
		input := scanner.Text()
		val, err := strconv.Atoi(input)
		if err == nil && val > 0 {
			return val
		}
		fmt.Println("Invalid input! Please enter a positive integer.")
	}
}

// Get valid matrix input from the user
func getMatrixInput(rows, cols int, name string) [][]float64 {
	matrix := make([][]float64, rows)
	fmt.Printf("Enter values for %s matrix (%dx%d), space-separated:\n", name, rows, cols)
	scanner := bufio.NewScanner(os.Stdin)

	for i := 0; i < rows; i++ {
		for {
			fmt.Printf("Row %d: ", i+1)
			scanner.Scan()
			line := scanner.Text()
			values := strings.Fields(line)

			if len(values) != cols {
				fmt.Println("Error: Number of values must match the number of columns.")
				continue
			}

			matrix[i] = make([]float64, cols)
			validRow := true

			for j, v := range values {
				val, err := strconv.ParseFloat(v, 64)
				if err != nil {
					fmt.Println("Invalid number:", err, "Please re-enter the row.")
					validRow = false
					break
				}
				matrix[i][j] = val
			}

			if validRow {
				break
			}
		}
	}
	return matrix
}

// Get a valid operation input
func getValidOperation() string {
	validOperations := map[string]bool{"add": true, "multiply": true, "transpose": true}
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Println("Choose an operation: add,multiply, transpose")
		fmt.Print("Enter operation: ")
		scanner.Scan()
		operation := strings.ToLower(scanner.Text())

		if validOperations[operation] {
			return operation
		}
		fmt.Println("Invalid operation! Please enter one of: add, subtract, multiply, transpose.")
	}
}

// Establish a secure connection
func connectToServer(address string) (*rpc.Client, error) {
	cert, err := tls.LoadX509KeyPair("server.crt", "server.key")
	if err != nil {
		return nil, fmt.Errorf("failed to load certificate: %w", err)
	}

	config := &tls.Config{Certificates: []tls.Certificate{cert}}

	conn, err := tls.Dial("tcp", address, config)
	if err != nil {
		return nil, fmt.Errorf("error connecting to secure controller: %w", err)
	}

	return rpc.NewClient(conn), nil
}

func main() {
	serverAddress := "192.168.18.89:5000" // Replace with actual server IP

	client, err := connectToServer(serverAddress)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	operation := getValidOperation()

	var req shared.StructureReq
	var res shared.StructureResponse

	if operation == "transpose" {
		rows := getValidIntegerInput("Enter number of rows: ")
		cols := getValidIntegerInput("Enter number of columns: ")
		mat1 := getMatrixInput(rows, cols, "matrix")

		req = shared.StructureReq{
			OperationType: operation,
			Mat1:          mat1,
			Mat2:          nil, // Not needed for transpose
		}
	} else {
		rows := getValidIntegerInput("Enter number of rows: ")
		cols := getValidIntegerInput("Enter number of columns: ")

		mat1 := getMatrixInput(rows, cols, "first")
		mat2 := getMatrixInput(rows, cols, "second")

		req = shared.StructureReq{
			OperationType: operation,
			Mat1:          mat1,
			Mat2:          mat2,
		}
	}

	// Call RPC function
	err = client.Call("Controller.Process", req, &res)
	if err != nil {
		fmt.Println("Error calling RPC:", err)
		return
	}

	// Display result
	fmt.Println("Result:")
	for _, row := range res.Res {
		fmt.Println(row)
	}
	fmt.Println("Processed by Worker:", res.Worker)
}