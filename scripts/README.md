# Scripts Documentation

This directory contains utility scripts to help with development, testing, and demonstration of the microservices application.

## Available Scripts

### 🏥 health-check.sh

**Purpose**: Comprehensive health check for all microservices and dependencies.

**Usage**:
```bash
./scripts/health-check.sh
```

**What it checks**:
- ✅ RabbitMQ Management UI accessibility (port 15672)
- ✅ API Server health endpoint (port 8081)
- ✅ Frontend web interface (port 8080)
- ✅ Job Runner process (checks if running)

**Output**: Color-coded status with quick links to all services.

**When to use**:
- After starting services to verify everything is working
- When debugging connectivity issues
- As part of deployment verification

---

### 🚀 demo.sh

**Purpose**: Interactive demonstration of the job processing workflow.

**Usage**:
```bash
./scripts/demo.sh
```

**What it does**:
1. **Verification**: Checks if API server is running
2. **Job Creation**: Creates 4 different test jobs:
   - Data Analysis Task
   - Email Campaign  
   - Backup Database
   - Generate Report
3. **Monitoring**: Tracks job status in real-time with emoji indicators:
   - 🟡 Pending
   - 🔄 Processing
   - ✅ Completed
   - ❌ Failed
4. **Results**: Shows final job details with full JSON output

**Features**:
- Real-time status updates every 5 seconds
- 90-second timeout protection
- Detailed final job results
- Color-coded output for easy reading

**When to use**:
- To demonstrate the full job processing workflow
- For testing end-to-end functionality
- To showcase asynchronous job processing
- For demos and presentations

---

### ✅ pre-commit.sh

**Purpose**: Runs the same quality checks as CI to ensure code will pass before committing.

**Usage**:
```bash
./scripts/pre-commit.sh
```

**What it checks**:
- ✅ **Dependencies**: Downloads and verifies Go modules
- ✅ **Static Analysis**: Runs `go vet` for potential issues
- ✅ **Code Formatting**: Ensures all files are properly formatted with `gofmt`
- ✅ **Linting**: Runs `golangci-lint` for code quality (installs if missing)
- ✅ **Testing**: Comprehensive test suite including:
  - Root module tests
  - Individual service tests (api-server, frontend, job-runner, shared)
  - Race condition detection
  - Integration tests (if RabbitMQ is available)

**Smart Features**:
- 🟡 **Auto-starts RabbitMQ** if not running (via Docker)
- 🔄 **Fallback to unit tests** if RabbitMQ unavailable
- 📦 **Auto-installs golangci-lint** if missing
- 🎨 **Color-coded output** for easy reading
- ⚡ **Fast feedback** - stops on first failure

**When to use**:
- Before committing code to ensure CI will pass
- During development to catch issues early
- As part of IDE pre-commit hooks
- For local quality assurance

---

## Prerequisites

Both scripts require:
- All microservices running (API server, job runner, frontend)
- RabbitMQ running and accessible
- `curl` command available
- Optional: `jq` for pretty JSON formatting (demo.sh)

## Running Scripts

Make sure scripts are executable:
```bash
chmod +x scripts/*.sh
```

Run from project root:
```bash
# Health check
./scripts/health-check.sh

# Demo
./scripts/demo.sh
```

## Script Features

### Color-coded Output
- 🟢 **Green**: Success/OK
- 🔴 **Red**: Error/Failed  
- 🟡 **Yellow**: Warning/Pending
- 🔵 **Blue**: Info/Processing

### Error Handling
- Both scripts check prerequisites before running
- Provide clear error messages with next steps
- Graceful handling of timeouts and failures

### Integration
- Scripts work with the existing Make commands
- Can be used in CI/CD pipelines
- Support both development and production environments

## Troubleshooting

If scripts fail:

1. **Check services are running**:
   ```bash
   make run-api     # Terminal 1
   make run-job     # Terminal 2  
   make run-frontend # Terminal 3
   ```

2. **Verify RabbitMQ**:
   ```bash
   make rabbitmq-up
   ```

3. **Check network connectivity**:
   ```bash
   curl http://localhost:8081/api/health
   ```

4. **Make scripts executable**:
   ```bash
   chmod +x scripts/*.sh
   ```
