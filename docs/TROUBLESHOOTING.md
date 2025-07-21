# Troubleshooting Guide

This guide provides comprehensive troubleshooting information for the Microservices Demo application.

## Table of Contents

- [Quick Diagnostics](#quick-diagnostics)
- [Common Issues](#common-issues)
- [Service-Specific Issues](#service-specific-issues)
- [Infrastructure Issues](#infrastructure-issues)
- [Performance Issues](#performance-issues)
- [Debugging Tools](#debugging-tools)
- [Log Analysis](#log-analysis)
- [Recovery Procedures](#recovery-procedures)

## Quick Diagnostics

### Health Check Script
```bash
#!/bin/bash
# Quick system health check

echo "🔍 Microservices Demo Health Check"
echo "=================================="

# Check if services are running
echo "📊 Service Status:"
curl -f http://localhost:8081/health 2>/dev/null && echo "✅ API Server: OK" || echo "❌ API Server: FAILED"
curl -f http://localhost:8080 2>/dev/null && echo "✅ Frontend: OK" || echo "❌ Frontend: FAILED"
curl -f http://localhost:15672 2>/dev/null && echo "✅ RabbitMQ Management: OK" || echo "❌ RabbitMQ Management: FAILED"

# Check ports
echo -e "\n🔌 Port Status:"
lsof -i :8080 >/dev/null && echo "✅ Port 8080: In use" || echo "❌ Port 8080: Available"
lsof -i :8081 >/dev/null && echo "✅ Port 8081: In use" || echo "❌ Port 8081: Available"
lsof -i :5672 >/dev/null && echo "✅ Port 5672: In use" || echo "❌ Port 5672: Available"

# Check Docker (if using Docker)
if command -v docker &> /dev/null; then
    echo -e "\n🐳 Docker Status:"
    docker ps --format "table {{.Names}}\t{{.Status}}" | grep -E "(rabbitmq|api-server|frontend|job-runner)"
fi

# Check Go processes
echo -e "\n🔧 Go Processes:"
pgrep -f "go run\|api-server\|frontend\|job-runner" || echo "No Go processes found"
```

### One-Line Diagnostics
```bash
# Quick service check
make health-check

# Check all containers
docker ps --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}"

# Check Kubernetes pods
kubectl get pods -o wide

# Check system resources
top -o cpu | head -10
```

## Common Issues

### 1. "Connection Refused" Errors

#### Symptoms
```
dial tcp [::1]:5672: connect: connection refused
curl: (7) Failed to connect to localhost port 8081: Connection refused
```

#### Diagnosis
```bash
# Check if services are running
ps aux | grep -E "(rabbitmq|api-server|frontend|job-runner)"

# Check port availability
netstat -tuln | grep -E "(5672|8080|8081|15672)"

# Check Docker containers
docker ps
```

#### Solutions

**For Local Development:**
```bash
# Start RabbitMQ
make rabbitmq-up

# Start services
make run-all-background

# Or start individually
cd api-server && go run main.go &
cd frontend && go run main.go &
cd job-runner && go run main.go &
```

**For Docker:**
```bash
# Restart Docker stack
make docker-down
make docker-up

# Check Docker network
docker network ls
docker network inspect microservices-demo_default
```

**For Kubernetes:**
```bash
# Check pod status
kubectl get pods
kubectl describe pod <pod-name>

# Check service endpoints
kubectl get endpoints
```

### 2. RabbitMQ Connection Issues

#### Symptoms
```
Exception (504) Reason: "channel/connection is not open"
dial tcp: lookup rabbitmq on 127.0.0.53:53: no such host
```

#### Diagnosis
```bash
# Check RabbitMQ status
curl http://localhost:15672

# Check RabbitMQ logs
docker logs microservices-demo_rabbitmq_1

# Test connection manually
telnet localhost 5672
```

#### Solutions

**RabbitMQ Not Started:**
```bash
# Docker Compose
docker-compose up rabbitmq -d

# Docker standalone
docker run -d --name rabbitmq \
  -p 5672:5672 -p 15672:15672 \
  -e RABBITMQ_DEFAULT_USER=guest \
  -e RABBITMQ_DEFAULT_PASS=guest \
  rabbitmq:3.12-management

# Wait for RabbitMQ to be ready
timeout 60s bash -c 'until curl -f http://localhost:15672; do sleep 2; done'
```

**Wrong Connection URL:**
```bash
# Check environment variables
echo $RABBITMQ_URL

# Correct formats:
export RABBITMQ_URL=amqp://guest:guest@localhost:5672/
export RABBITMQ_URL=amqp://guest:guest@rabbitmq:5672/  # Docker network
```

**Network Issues:**
```bash
# Docker network troubleshooting
docker network create microservices-net
docker run --network microservices-net ...

# Kubernetes DNS troubleshooting
kubectl run debug --image=busybox -it --rm --restart=Never -- nslookup rabbitmq
```

### 3. Job Processing Issues

#### Symptoms
```
Jobs stuck in "pending" status
Jobs created but never processed
Job runner not consuming messages
```

#### Diagnosis
```bash
# Check RabbitMQ queue status
curl -u guest:guest http://localhost:15672/api/queues

# Check job runner logs
docker logs microservices-demo_job-runner_1

# Test job creation manually
curl -X POST http://localhost:8081/api/jobs \
  -H "Content-Type: application/json" \
  -d '{"description": "Test job"}'
```

#### Solutions

**Job Runner Not Started:**
```bash
# Start job runner
cd job-runner && go run main.go

# Or with Docker
docker-compose up job-runner -d

# Check if multiple instances are running
docker-compose up --scale job-runner=3 -d
```

**Queue Configuration Issues:**
```bash
# Reset RabbitMQ queues
docker exec microservices-demo_rabbitmq_1 rabbitmqctl purge_queue job_queue
docker exec microservices-demo_rabbitmq_1 rabbitmqctl delete_queue job_queue

# Restart services to recreate queues
make docker-restart
```

### 4. Frontend Issues

#### Symptoms
```
404 Not Found for static assets
"Failed to fetch" errors in browser console
Blank pages or template errors
```

#### Diagnosis
```bash
# Check frontend logs
docker logs microservices-demo_frontend_1

# Test API connectivity from frontend
docker exec microservices-demo_frontend_1 wget -O- http://api-server:8081/health

# Check browser network tab for failed requests
```

#### Solutions

**Static Asset Issues:**
```bash
# Rebuild frontend with assets
cd frontend
go build -o bin/frontend .

# Check if assets are embedded
ls -la static/
```

**API Connection Issues:**
```bash
# Check API server URL configuration
echo $API_SERVER_URL

# Update configuration
export API_SERVER_URL=http://localhost:8081  # Local
export API_SERVER_URL=http://api-server:8081  # Docker
```

### 5. Build and Compilation Issues

#### Symptoms
```
go: module not found
undefined: SomeFunction
build constraints exclude all Go files
```

#### Solutions

**Module Issues:**
```bash
# Clean module cache
go clean -modcache
go mod download
go mod verify

# Reset go.sum if corrupted
rm go.sum
go mod tidy
```

**Build Tag Issues:**
```bash
# Check build tags
go list -tags=integration ./...

# Build with specific tags
go build -tags=integration .
```

**Version Conflicts:**
```bash
# Update dependencies
go get -u ./...
go mod tidy

# Downgrade if needed
go get github.com/some/package@v1.2.3
```

## Service-Specific Issues

### API Server Issues

#### High Memory Usage
```bash
# Check memory usage
docker stats microservices-demo_api-server_1

# Profile memory usage
go tool pprof http://localhost:8081/debug/pprof/heap
```

**Solutions:**
```bash
# Implement connection pooling
# Add memory limits in Docker
# Use sync.Pool for frequent allocations
```

#### Request Timeout
```bash
# Check for long-running requests
curl -w "@curl-format.txt" http://localhost:8081/api/jobs

# Increase timeout settings
export REQUEST_TIMEOUT=30s
```

### Job Runner Issues

#### Worker Exhaustion
```bash
# Check job queue depth
curl -u guest:guest http://localhost:15672/api/queues/%2F/job_queue

# Scale job runners
docker-compose up --scale job-runner=5 -d
```

#### Memory Leaks
```bash
# Monitor memory over time
while true; do docker stats --no-stream job-runner; sleep 10; done

# Profile memory
go tool pprof http://localhost:6060/debug/pprof/heap
```

### Frontend Issues

#### Template Errors
```bash
# Check template parsing
cd frontend
go run main.go -validate-templates

# Common template issues:
# - Missing template files
# - Syntax errors in templates
# - Missing template data
```

## Infrastructure Issues

### Docker Issues

#### Container Won't Start
```bash
# Check container logs
docker logs <container-name>

# Check Docker daemon
sudo systemctl status docker

# Check disk space
df -h
docker system df
```

#### Network Issues
```bash
# List Docker networks
docker network ls

# Inspect network
docker network inspect microservices-demo_default

# Recreate network
docker-compose down
docker network prune
docker-compose up -d
```

#### Volume Issues
```bash
# Check volumes
docker volume ls

# Clean volumes
docker volume prune

# Backup/restore volumes
docker run --rm -v rabbitmq_data:/data -v $(pwd):/backup busybox tar czf /backup/backup.tar.gz /data
```

### Kubernetes Issues

#### Pod Scheduling Issues
```bash
# Check node resources
kubectl describe nodes

# Check pod resource requests
kubectl describe pod <pod-name>

# Check for resource constraints
kubectl get events --sort-by=.metadata.creationTimestamp
```

#### Service Discovery Issues
```bash
# Test DNS resolution
kubectl run debug --image=busybox -it --rm --restart=Never -- nslookup api-server

# Check service endpoints
kubectl get endpoints

# Verify label selectors
kubectl get pods --show-labels
```

#### Ingress Issues
```bash
# Check ingress status
kubectl get ingress -o wide

# Check ingress controller logs
kubectl logs -n ingress-nginx deployment/nginx-ingress-controller

# Test ingress connectivity
curl -H "Host: microservices.example.com" http://<ingress-ip>/
```

## Performance Issues

### High CPU Usage

#### Diagnosis
```bash
# Monitor CPU usage
top -p $(pgrep -d',' -f 'api-server|job-runner|frontend')

# Profile CPU usage
go tool pprof http://localhost:8081/debug/pprof/profile?seconds=30
```

#### Solutions
```bash
# Optimize hot paths
# Add CPU limits
# Scale horizontally
# Use connection pooling
```

### High Memory Usage

#### Diagnosis
```bash
# Check memory usage
ps aux --sort=-%mem | head -10

# Profile memory
go tool pprof http://localhost:8081/debug/pprof/heap
```

#### Solutions
```bash
# Fix memory leaks
# Add memory limits
# Use object pooling
# Optimize data structures
```

### Slow Response Times

#### Diagnosis
```bash
# Measure response times
curl -w "@curl-format.txt" http://localhost:8081/api/jobs

# Check database query times (if using database)
# Profile application
go tool pprof http://localhost:8081/debug/pprof/profile
```

## Debugging Tools

### Go Debugging

#### pprof Profiling
```bash
# Add pprof endpoint to your services
import _ "net/http/pprof"

go func() {
    log.Println(http.ListenAndServe("localhost:6060", nil))
}()

# Profile CPU
go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30

# Profile memory
go tool pprof http://localhost:6060/debug/pprof/heap

# Profile goroutines
go tool pprof http://localhost:6060/debug/pprof/goroutine
```

#### Delve Debugger
```bash
# Install delve
go install github.com/go-delve/delve/cmd/dlv@latest

# Debug application
cd api-server
dlv debug . -- --port=8081

# Attach to running process
dlv attach $(pgrep api-server)
```

### Network Debugging

#### tcpdump/Wireshark
```bash
# Capture traffic on localhost
sudo tcpdump -i lo port 5672 -w rabbitmq.pcap

# Monitor HTTP traffic
sudo tcpdump -i any port 8081 -A
```

#### curl for API testing
```bash
# Test with timing
curl -w "@curl-format.txt" http://localhost:8081/api/jobs

# Test with verbose output
curl -v http://localhost:8081/api/jobs

# Test specific scenarios
curl -X POST http://localhost:8081/api/jobs \
  -H "Content-Type: application/json" \
  -d '{"description": "Debug test job"}' \
  -w "Response time: %{time_total}s\n"
```

### Container Debugging

#### Docker debugging
```bash
# Execute commands in running container
docker exec -it <container-name> /bin/sh

# Copy files from container
docker cp <container-name>:/app/logs ./logs

# Inspect container
docker inspect <container-name>
```

#### Kubernetes debugging
```bash
# Debug pod
kubectl exec -it <pod-name> -- /bin/sh

# Port forward for debugging
kubectl port-forward pod/<pod-name> 8080:8080

# Debug with a debug container
kubectl debug <pod-name> -it --image=busybox --target=<container-name>
```

## Log Analysis

### Log Locations

#### Local Development
```bash
# Service logs (if using systemd)
journalctl -u api-server -f

# Application logs (stdout/stderr)
./api-server 2>&1 | tee api-server.log
```

#### Docker
```bash
# Container logs
docker logs -f microservices-demo_api-server_1

# All services
docker-compose logs -f

# Specific time range
docker logs --since "2024-01-01T00:00:00" <container-name>
```

#### Kubernetes
```bash
# Pod logs
kubectl logs -f deployment/api-server

# Previous container logs
kubectl logs deployment/api-server --previous

# All containers in pod
kubectl logs -f pod/<pod-name> --all-containers=true
```

### Log Analysis Patterns

#### Error Patterns
```bash
# Find errors in logs
grep -i error api-server.log

# Count error types
grep -o "error.*" api-server.log | sort | uniq -c

# Find connection issues
grep -i "connection refused\|timeout\|network" *.log
```

#### Performance Analysis
```bash
# Find slow requests
grep "duration.*ms" api-server.log | awk '$NF > 1000'

# Analyze request patterns
grep "POST /api/jobs" api-server.log | wc -l
```

### Structured Log Analysis

#### JSON logs with jq
```bash
# Parse JSON logs
cat api-server.log | jq 'select(.level == "error")'

# Aggregate metrics
cat api-server.log | jq -r '.duration' | sort -n | tail -10
```

## Recovery Procedures

### Service Recovery

#### Restart Individual Services
```bash
# Docker Compose
docker-compose restart api-server
docker-compose restart job-runner
docker-compose restart frontend

# Kubernetes
kubectl rollout restart deployment/api-server
kubectl rollout restart deployment/job-runner
kubectl rollout restart deployment/frontend
```

#### Complete System Recovery
```bash
# Docker Compose
docker-compose down
docker-compose up -d

# Kubernetes
kubectl delete pods --all
# Pods will be recreated automatically

# Local development
make stop-all
make clean
make docker-up
```

### Data Recovery

#### RabbitMQ Queue Recovery
```bash
# Check queue status
curl -u guest:guest http://localhost:15672/api/queues

# Purge queues if needed
docker exec rabbitmq rabbitmqctl purge_queue job_queue

# Export/import queue definitions
docker exec rabbitmq rabbitmqctl export_definitions /tmp/definitions.json
docker exec rabbitmq rabbitmqctl import_definitions /tmp/definitions.json
```

#### Job State Recovery
```bash
# Since jobs are stored in memory, they're lost on restart
# For production, implement persistent storage recovery:

# 1. Database recovery procedures
# 2. Event sourcing replay
# 3. State reconstruction from message queue
```

### Disaster Recovery

#### Complete Environment Rebuild
```bash
# 1. Stop everything
make docker-down
docker system prune -f

# 2. Rebuild images
make docker-build

# 3. Start fresh
make docker-up

# 4. Verify functionality
make health-check
./scripts/demo.sh
```

#### Backup and Restore
```bash
# Backup Docker volumes
docker run --rm -v microservices-demo_rabbitmq_data:/data \
  -v $(pwd):/backup busybox \
  tar czf /backup/rabbitmq-backup.tar.gz /data

# Restore Docker volumes
docker run --rm -v microservices-demo_rabbitmq_data:/data \
  -v $(pwd):/backup busybox \
  tar xzf /backup/rabbitmq-backup.tar.gz -C /
```

---

## Emergency Contacts and Escalation

### Quick Reference Commands
```bash
# Stop everything immediately
make docker-down
pkill -f "go run\|api-server\|job-runner\|frontend"

# Emergency health check
./scripts/health-check.sh

# Get system information
docker version && docker-compose version && go version && kubectl version

# Check system resources
df -h && free -h && top -n1 -b
```

---

## Related Documentation

- [Development Guide](DEVELOPMENT.md) - Local development setup and debugging
- [Deployment Guide](DEPLOYMENT.md) - Production deployment troubleshooting  
- [Architecture Documentation](ARCHITECTURE.md) - System design understanding
- [API Documentation](API.md) - API endpoint testing and validation
