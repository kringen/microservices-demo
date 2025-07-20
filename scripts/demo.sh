#!/bin/bash

# ============================================================================
# Microservices Job Processing Demo Script
# ============================================================================
# 
# This script demonstrates the complete job processing workflow by:
# 1. Creating 4 different test jobs via the API
# 2. Monitoring their status in real-time  
# 3. Displaying final results with full details
#
# Prerequisites:
# - API Server running on localhost:8081
# - Job Runner service consuming from RabbitMQ
# - Frontend web app on localhost:8080 (optional)
# - RabbitMQ server running on localhost:5672
#
# Usage: ./scripts/demo.sh
#
# The script will:
# - Verify API server accessibility
# - Create jobs: Data Analysis, Email Campaign, Backup, Report Generation
# - Monitor status every 5 seconds with emoji indicators
# - Show final job details (JSON output)
# - Timeout after 90 seconds if jobs don't complete
#
# ============================================================================

API_URL="http://localhost:8081"

echo "🚀 Microservices Job Processing Demo"
echo "====================================="

# Check if API server is running
if ! curl -s -f "${API_URL}/api/health" > /dev/null; then
    echo "❌ API server is not running at ${API_URL}"
    echo "Please start the services first:"
    echo "  Terminal 1: cd api-server && go run main.go"
    echo "  Terminal 2: cd job-runner && go run main.go"
    echo "  Terminal 3: cd frontend && go run main.go"
    exit 1
fi

echo "✅ API server is running"

# Function to create a job
create_job() {
    local title="$1"
    local description="$2"
    
    echo "📝 Creating job: $title"
    
    local response=$(curl -s -X POST "${API_URL}/api/jobs" \
        -H "Content-Type: application/json" \
        -d "{\"title\": \"$title\", \"description\": \"$description\"}")
    
    local job_id=$(echo "$response" | grep -o '"id":"[^"]*"' | cut -d'"' -f4)
    echo "   Job ID: $job_id"
    echo "$job_id"
}

# Function to check job status
check_job() {
    local job_id="$1"
    curl -s "${API_URL}/api/jobs/${job_id}" | \
        grep -o '"status":"[^"]*"' | cut -d'"' -f4
}

# Function to get job details
get_job_details() {
    local job_id="$1"
    curl -s "${API_URL}/api/jobs/${job_id}"
}

# Create test jobs
echo ""
echo "Creating test jobs..."

job1=$(create_job "Data Analysis Task" "Analyze customer data for monthly report")
job2=$(create_job "Email Campaign" "Send marketing emails to customer list")
job3=$(create_job "Backup Database" "Perform daily database backup")
job4=$(create_job "Generate Report" "Create quarterly financial report")

jobs=($job1 $job2 $job3 $job4)

echo ""
echo "Monitoring job progress (will check every 5 seconds for up to 90 seconds)..."

start_time=$(date +%s)
max_wait=90

while true; do
    current_time=$(date +%s)
    elapsed=$((current_time - start_time))
    
    if [ $elapsed -gt $max_wait ]; then
        echo "⏰ Timeout reached (90 seconds)"
        break
    fi
    
    echo ""
    echo "⏱️  Time elapsed: ${elapsed}s"
    
    all_completed=true
    
    for job_id in "${jobs[@]}"; do
        if [ -n "$job_id" ]; then
            status=$(check_job "$job_id")
            
            case "$status" in
                "pending")
                    echo "   $job_id: 🟡 Pending"
                    all_completed=false
                    ;;
                "processing")
                    echo "   $job_id: 🔄 Processing"
                    all_completed=false
                    ;;
                "completed")
                    echo "   $job_id: ✅ Completed"
                    ;;
                "failed")
                    echo "   $job_id: ❌ Failed"
                    ;;
                *)
                    echo "   $job_id: ❓ Unknown status: $status"
                    all_completed=false
                    ;;
            esac
        fi
    done
    
    if [ "$all_completed" = true ]; then
        echo ""
        echo "🎉 All jobs completed!"
        break
    fi
    
    sleep 5
done

echo ""
echo "📊 Final Job Details:"
echo "===================="

for job_id in "${jobs[@]}"; do
    if [ -n "$job_id" ]; then
        echo ""
        echo "Job ID: $job_id"
        get_job_details "$job_id" | jq . 2>/dev/null || get_job_details "$job_id"
    fi
done

echo ""
echo "✨ Demo completed!"
echo "💡 You can also view the jobs in the web interface at: http://localhost:8080"
