# Object Pool Pattern Demonstration in Go

## Overview
This repository showcases the implementation of the Object Pool design pattern in Go. The project demonstrates how to manage a pool of reusable resources — in this case, database connections — efficiently. The primary focus is on optimizing resource allocation, minimizing the overhead of resource initialization, and managing resource lifecycle to enhance application performance.

## Pattern Description
The Object Pool Pattern is used to manage a set of ready-to-use objects from a pool rather than allocating and freeing them individually. This is especially useful in applications where the cost of initializing a class instance is high, the rate of instantiation of a class is high, and the number of instantiations in use at any one time is low. In this project, the pattern is implemented to manage database connections, ensuring that connections are reused efficiently and that the pool dynamically adjusts to load, improving overall application throughput.

## Project Structure
- **cmd/**: Contains the application entry point (`main.go`), demonstrating the use of the Object Pool pattern to manage database connections.
- **pkg/**
    - **pool/**: Implements the connection pool, managing the lifecycle and availability of database connections.
- **internal/**
    - **database/**: Contains the `Connection` class definition which represents individual database connections.

## Getting Started

### Prerequisites
Ensure you have Go installed on your system. You can download it from [Go's official site](https://golang.org/dl/).

### Installation
Clone this repository to your local machine:
```bash
git clone https://github.com/arthurfp/Go_Object_Pool_Pattern.git
cd Go_Object_Pool_Pattern
```

### Running the Application
To run the application, execute:
```bash
go run cmd/main.go
```

### Running the Tests
To execute the tests and verify the functionality:
```bash
go test ./...
```

### Contributing
Contributions are welcome! Please feel free to submit pull requests or open issues to discuss proposed changes or enhancements.

### Author
Arthur Ferreira - github.com/arthurfp
