# SwiftSearch

## Introduction
SwiftSearch is a robust and efficient tool designed to enhance file searching operations. Leveraging the power of Go, SwiftSearch aims to deliver fast and accurate search results, streamlining the file lookup process in large directories or complex data structures.

## Setup Instructions
1. **Clone the repository**:
   ```bash
   git clone https://github.com/ashborn3/SwiftSearch.git
   cd SwiftSearch
   ```

2. **Install dependencies**:
   Ensure Go is installed on your system (Go 1.16+ recommended).
   ```bash
   go mod download
   ```

3. **Run the application**:
   ```bash
   go run main.go
   ```

## Benchmark Results
*Note: The following results are for demonstration purposes.*

- **Test Environment**: 
  - CPU: Intel i7 10th Gen
  - RAM: 16GB
  - Storage: NVMe SSD

- **Performance Metrics**:
  - **Single File Search**: Avg. time: 120ms
  - **Directory (10,000 files)**: Avg. time: 350ms
  - **Large Dataset (1M+ files)**: Avg. time: 2.5s

