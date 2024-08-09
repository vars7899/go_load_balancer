# Go Load Balancer

A simple and efficient load balancer implemented in Go. This project aims to distribute incoming network traffic across multiple servers to ensure high availability and reliability.

## Features

- **Dynamic Server Pool**: Automatically manage a pool of servers.
- **Health Checks**: Regularly check the health of servers to ensure they are available.
- **Retry Mechanism**: Automatically retry requests to healthy servers if a server fails.
- **Configurable Strategies**: Choose from various load balancing strategies.

## Load Balancing Strategies

The following load balancing strategies are currently implemented:

1. **Round Robin**: Distributes requests evenly across all servers in a circular order.
2. **Least Connections**: Directs traffic to the server with the fewest active connections.
3. **Weighted Round Robin**: Allow servers to have different weights to influence request distribution.

### Future Plans

- **Random**: Selects a server at random for each request.
- **Session Persistence**: Implement sticky sessions to route requests from the same client to the same server.
- **Advanced Health Checks**: Introduce more sophisticated health check mechanisms, such as response time monitoring.
- **Metrics and Monitoring**: Add metrics collection and monitoring capabilities for better insights into server performance.

## Getting Started

### Prerequisites

- Go 1.16 or later
- Dependencies managed via Go modules

### Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/yourusername/go_load_balancer.git
   cd go_load_balancer
   ```

2. Install dependencies:
   ```bash
   go mod tidy
   ```

### Configuration

Before running the load balancer, configure the `config.yaml` file to specify the servers and load balancing strategy.

### Running the Load Balancer

To start the load balancer, run:

```bash
  git clone https://github.com/vars7899/go_load_balancer.git
  cd go_load_balancer
  go run .
```

#stress test load balancer

```bash
  hey -n 20000 -c 100 http://localhost:8080/lb
```

# mock server cluster

```go
  package main

  import (
    "fmt"
    "log"
    "net/http"
  )

  func main() {
    go startServer(":8081", "s1")
    go startServer(":8082", "s2")
    go startServer(":8083", "s3")
    go startServer(":8084", "s4")
    go startServer(":8085", "s5")

    // Block the main goroutine to keep the program running
    select {}
  }

  func startServer(port, serverName string) {
    fmt.Printf("%s: started\n", serverName)
    mux := http.NewServeMux()
    mux.HandleFunc("/", func(rw http.ResponseWriter, rq *http.Request) {
      fmt.Printf("%s: %s %s %s\n", serverName, rq.Method, rq.Host, rq.RequestURI)
      // time.Sleep(1 * time.Second)
      rw.WriteHeader(http.StatusOK)
      rw.Write([]byte("<--response-->"))
    })

    server := &http.Server{
      Addr:    port,
      Handler: mux,
    }

    if err := server.ListenAndServe(); err != nil {
      log.Fatalf("%s: %v", serverName, err)
    }
  }
```
