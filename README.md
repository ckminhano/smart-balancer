# smart-balancer

This is a project to create a smart load balancer that can be used in various applications. 
It is designed to distribute workloads efficiently across multiple servers or resources using a set of criteria.
The criteria for balancing is based on the backend server's response time, CPU usage, and memory usage.

## Features
- Load balancing based on response time, CPU usage, and memory usage
- Support for multiple backend servers
- Dynamic adjustment of load balancing criteria
- Easy to integrate with existing applications using REST APIs and prometheus metrics

## Requirements
- Go 1.18 or later
- Docker for containerization
- Prometheus for monitoring and metrics in the backend

## Installation
1. Clone the repository:
   ```bash
   git clone 
   ```

2. cd smart-balancer
   ```
    cd smart-balancer
    ```
3. Build the project:
    ```bash
    go build -o smart-balancer .
    ```



