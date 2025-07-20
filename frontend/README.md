# Frontend

The Frontend is a lightweight web application that provides a user-friendly interface for job submission, status monitoring, and result viewing. Built with Go and Gin, it offers both traditional HTML forms and modern AJAX interactions.

## 🎯 Purpose

- **User Interface**: Web-based job submission and monitoring
- **Real-time Updates**: Live job status with automatic refresh
- **API Integration**: Seamless communication with API server
- **Responsive Design**: Modern Bootstrap-based UI

## 🏗️ Architecture

```
User Browser ↔ Frontend Web App ↔ API Server ↔ RabbitMQ ↔ Job Runner
                     ↓
                HTML Templates
                AJAX Polling
                Bootstrap UI
```

## ✨ Features

### Job Management
- **Job Submission**: Easy form-based job creation
- **Status Monitoring**: Real-time job status updates
- **Job History**: View all submitted jobs
- **Detailed View**: Individual job status pages

### User Experience
- **AJAX Integration**: No page refresh for job submission
- **Auto-refresh**: Job list updates every 3 seconds
- **Status Indicators**: Color-coded job status badges
- **Responsive Design**: Works on desktop and mobile

### System Monitoring
- **Health Dashboard**: API server and RabbitMQ status
- **Service Status**: Real-time system health checks
- **Quick Links**: Direct access to all system components

## 🚀 Quick Start

### Local Development
```bash
# Run with default settings
make run

# Run with custom API server
API_SERVER_URL=http://localhost:8081 make run-env

# Run with hot reload
make dev

# Open in browser
make open
```

### Docker
```bash
# Build and run standalone
make docker-run

# Run with Docker network (for full stack)
make docker-run-network
```

### Testing
```bash
# Run tests
make test

# Check if frontend is accessible
make health-check
```

## ⚙️ Configuration

### Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `API_SERVER_URL` | `http://localhost:8081` | API server base URL |
| `GIN_MODE` | `debug` | Gin framework mode (debug/release) |
| `PORT` | `8080` | Frontend server port |

### Example Configuration
```bash
export API_SERVER_URL="http://api-server:8081"
export GIN_MODE="release"
```

## 🎨 User Interface

### Main Dashboard
- **Job Creation Form**: Title and description input
- **Recent Jobs Panel**: Live-updating job list
- **System Status**: Health indicators for all services

### Job Status Page
- **Job Details**: ID, title, description, timestamps
- **Status Indicators**: Pending, Processing, Completed, Failed
- **Results Display**: Job output and error messages
- **Auto-refresh**: Updates every 3 seconds for active jobs

### UI Components
```html
<!-- Job Status Badges -->
<span class="badge bg-warning">Pending</span>
<span class="badge bg-info">Processing</span>
<span class="badge bg-success">Completed</span>
<span class="badge bg-danger">Failed</span>

<!-- Loading Spinners -->
<div class="spinner-border spinner-border-sm"></div>
```

## 🔄 API Integration

### Job Submission
```javascript
fetch('/api/jobs', {
    method: 'POST',
    headers: {'Content-Type': 'application/json'},
    body: JSON.stringify({
        title: 'Data Analysis',
        description: 'Process customer data'
    })
})
```

### Status Polling
```javascript
// Auto-refresh every 3 seconds
setInterval(() => {
    fetch('/api/jobs')
        .then(response => response.json())
        .then(data => updateJobsList(data.jobs));
}, 3000);
```

### Health Monitoring
```javascript
fetch('/api/status')
    .then(response => response.json())
    .then(status => updateSystemStatus(status));
```

## 📱 Responsive Design

### Bootstrap Integration
- **Grid System**: Responsive column layouts
- **Components**: Cards, badges, buttons, forms
- **Utilities**: Spacing, colors, typography

### Mobile Support
- **Viewport Meta**: Proper mobile scaling
- **Touch-friendly**: Large buttons and links
- **Responsive Tables**: Horizontal scrolling on small screens

## 🧪 Testing

### Manual Testing
```bash
# Start frontend
make run

# Open in browser
make open

# Test job submission
# 1. Fill out job form
# 2. Submit job
# 3. Verify AJAX response
# 4. Check job appears in list
# 5. Click job link to view details
```

### Automated Testing
```bash
# Run unit tests
make test

# Test template rendering
go test -run TestTemplates

# Test API integration
go test -run TestAPIIntegration
```

### Load Testing
```bash
# Test with multiple concurrent users
# Submit many jobs simultaneously
# Verify UI remains responsive
```

## 🎭 Templates

### Template Structure
```go
// Inline templates with Go template syntax
const indexTemplate = `{{define "index"}}...{{end}}`
const jobStatusTemplate = `{{define "job-status"}}...{{end}}`
```

### Template Functions
```go
template.FuncMap{
    "formatTime": func(t time.Time) string {
        return t.Format("2006-01-02 15:04:05")
    },
    "statusColor": func(status JobStatus) string {
        // Return Bootstrap color class
    },
}
```

## 🔧 Development

### Prerequisites
- Go 1.21+
- API server running
- Access to shared module

### Hot Reload Setup
```bash
# Install air for hot reload
go install github.com/cosmtrek/air@latest

# Run with hot reload
make dev
```

### Code Structure
```
frontend/
├── main.go           # Server and routes
├── templates.go      # HTML templates
├── main_test.go      # Unit tests
├── Dockerfile        # Container definition
├── Makefile          # Build commands
└── README.md         # This file
```

## 🐳 Docker

### Dockerfile Features
- **Multi-stage build** for smaller images
- **Health checks** for monitoring
- **Environment configuration**
- **Static asset handling**

### Container Health Check
```dockerfile
HEALTHCHECK --interval=30s --timeout=3s \
    CMD wget --spider http://localhost:8080/ || exit 1
```

## 📊 Monitoring

### Frontend Metrics
- Page load times
- AJAX response times
- User interaction events
- Error rates

### Browser Console
- Network requests
- JavaScript errors
- Performance timing
- API response inspection

## 🔍 Debugging

### Development Mode
```bash
# Enable debug logging
GIN_MODE=debug make run

# Check browser developer tools
# - Network tab for API calls
# - Console for JavaScript errors
# - Elements for DOM inspection
```

### Common Issues
```bash
# API server not accessible
curl http://localhost:8081/api/health

# CORS issues
# Check browser console for CORS errors

# Template errors
# Check server logs for template compilation errors
```

## 🛠️ Troubleshooting

### Connection Issues
```bash
# Check API server status
make health-check

# Verify API server URL
echo $API_SERVER_URL

# Test API connectivity
curl http://localhost:8081/api/jobs
```

### UI Issues
```bash
# Clear browser cache
# Check browser console for errors
# Verify Bootstrap CSS/JS loading
# Test with different browsers
```

### Performance Issues
```bash
# Monitor network requests
# Optimize AJAX polling frequency
# Check for memory leaks in browser
# Minimize template complexity
```

## 🎯 Use Cases

### Development
- Testing job submission workflows
- Monitoring job processing in real-time
- Debugging API integration
- Demonstrating microservices interaction

### Production
- Admin interface for job management
- User-facing job submission portal
- System monitoring dashboard
- Customer service tools

## 🔐 Security

### Best Practices
- Input validation on forms
- CSRF protection
- XSS prevention
- Secure API communication

### Environment Security
```bash
# Use HTTPS in production
# Secure API endpoints
# Validate environment variables
# Implement authentication if needed
```

## 📚 Related Services

- **[API Server](../api-server/README.md)** - Backend API
- **[Job Runner](../job-runner/README.md)** - Job processing
- **[Shared](../shared/README.md)** - Common utilities

## 🎨 Customization

### Styling
- Modify Bootstrap variables
- Add custom CSS classes
- Update color schemes
- Customize component layouts

### Functionality
- Add new job types
- Implement user authentication
- Add job filtering/search
- Extend monitoring capabilities
