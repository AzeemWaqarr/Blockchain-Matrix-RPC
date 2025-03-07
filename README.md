# Secure Matrix Computation using Blockchain-Based Client-Server RPC Model 🔒

## 📌 Project Overview
This project implements a **secure client-server architecture** for **matrix computations** using **Remote Procedure Calls (RPC)** and **TLS encryption**. The system distributes computation across multiple worker nodes while ensuring **fault tolerance** and **load balancing**.

## 🛠️ Installation & Setup
### 1. Clone the Repository
```sh
git clone https://github.com/AzeemWaqarr/Blockchain-Matrix-RPC.git
cd Blockchain-Matrix-RPC
```

### 2. Install Dependencies
Ensure you have **Go installed** and initialize the module:
```sh
go mod tidy
```

### 3. Generate TLS Certificates (If Not Provided)
To create secure connections, generate self-signed certificates:
```sh
openssl req -x509 -newkey rsa:4096 -keyout security/server.key -out security/server.crt -days 365 -nodes
```

---
## 📂 Project Structure
```
📦 Blockchain-Matrix-RPC
├── 📄 README.md       # Project Documentation
├── 📄 Block Chain Assignment 1.pdf  # Assignment Instructions
├── 📂 shared          # Shared data structures
│   ├── struct.go
├── 📂 server
│   ├── controller.go  # Server-side logic (Coordinator)
│   ├── worker.go      # Worker Node Code
├── 📂 client
│   ├── client.go      # Client Program
├── 📂 security
│   ├── server.crt     # TLS Certificate
│   ├── server.key     # TLS Private Key
│   ├── ca.crt         # Certificate Authority (CA) Cert
│   ├── ca.key         # CA Private Key
└── 📄 go.mod          # Go Module File
```

---
## ⚡ System Architecture
1️⃣ **Client:** Requests matrix operations via **RPC**.  
2️⃣ **Coordinator (Server):** Assigns tasks using **First-Come, First-Served (FCFS)** scheduling and load balancing.  
3️⃣ **Workers:** Perform matrix computations (Addition, Multiplication, Transpose).  
4️⃣ **TLS Security:** Secure communication between client, server, and workers.  

---
## 🔗 Secure RPC Communication
- The system uses **TLS encryption** for secure data transfer.
- The server and workers exchange encrypted data to prevent unauthorized access.

---
## 🚀 Running the Application
### 1. Start the Controller (Server)
```sh
cd server
go run controller.go
```

### 2. Start Worker Nodes
```sh
cd server
go run worker.go
```

### 3. Run the Client
```sh
cd client
go run client.go
```

---
## 📊 Features
✅ Secure TLS communication 🔒  
✅ Load Balancing (Least Busy Worker) ⚖️  
✅ Fault Tolerance (Automatic Worker Reassignment) 🔄  
✅ Supports Matrix Addition, Multiplication & Transpose 🧮  

## 📬 Contact
- **Email:** azeem.waqarr@gmail.com  
- **GitHub:** [AzeemWaqar](https://github.com/AzeemWaqar)  
