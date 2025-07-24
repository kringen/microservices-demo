package main

const indexTemplate = `
{{define "index"}}
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.Title}}</title>
    <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.1.3/dist/css/bootstrap.min.css" rel="stylesheet">
    <style>
        .job-card {
            margin-bottom: 1rem;
        }
        .status-badge {
            font-size: 0.8rem;
        }
    </style>
</head>
<body>
    <nav class="navbar navbar-dark bg-dark">
        <div class="container">
            <a class="navbar-brand" href="/">
                <strong>Microservices Demo</strong>
            </a>
            <span class="navbar-text">
                Go + RabbitMQ
            </span>
        </div>
    </nav>

    <div class="container mt-4">
        <div class="row">
            <div class="col-md-6">
                <div class="card">
                    <div class="card-header">
                        <h5 class="card-title">Create New Job</h5>
                    </div>
                    <div class="card-body">
                        <form id="jobForm" action="/submit" method="POST">
                            <div class="mb-3">
                                <label for="title" class="form-label">Job Title *</label>
                                <input type="text" class="form-control" id="title" name="title" required>
                            </div>
                            <div class="mb-3">
                                <label for="description" class="form-label">Description</label>
                                <textarea class="form-control" id="description" name="description" rows="3" 
                                    placeholder="Try: 'Data Analysis', 'Email Campaign', 'Backup Task', etc."></textarea>
                            </div>
                            <button type="submit" class="btn btn-primary" id="submitBtn">Submit Job</button>
                        </form>
                        
                        <!-- Success/Error messages -->
                        <div id="jobCreateMessage" class="mt-3" style="display: none;"></div>
                    </div>
                </div>
            </div>
            
            <div class="col-md-6">
                <div class="card">
                    <div class="card-header d-flex justify-content-between align-items-center">
                        <h5 class="card-title mb-0">Recent Jobs</h5>
                        <button class="btn btn-sm btn-outline-secondary" onclick="window.location.reload()">
                            Refresh
                        </button>
                    </div>
                    <div class="card-body">
                        {{if .Jobs}}
                            {{range .Jobs}}
                            <div class="job-card">
                                <div class="d-flex justify-content-between align-items-start">
                                    <div>
                                        <h6 class="mb-1">
                                            <a href="/status/{{.ID}}" class="text-decoration-none">{{.Title}}</a>
                                        </h6>
                                        <p class="mb-1 text-muted small">{{.Description}}</p>
                                        <small class="text-muted">{{formatTime .CreatedAt}}</small>
                                    </div>
                                    <span class="badge bg-{{statusColor .Status}} status-badge">
                                        {{.Status}}
                                    </span>
                                </div>
                                <hr class="my-2">
                            </div>
                            {{end}}
                        {{else}}
                            <p class="text-muted">No jobs yet. Create your first job!</p>
                        {{end}}
                    </div>
                </div>
            </div>
        </div>

        <div class="row mt-4">
            <div class="col-12">
                <div class="card">
                    <div class="card-header">
                        <h5 class="card-title">System Status</h5>
                    </div>
                    <div class="card-body">
                        <div id="system-status">
                            <div class="spinner-border spinner-border-sm" role="status">
                                <span class="visually-hidden">Loading...</span>
                            </div>
                            Checking system status...
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>

    <script src="https://cdn.jsdelivr.net/npm/bootstrap@5.1.3/dist/js/bootstrap.bundle.min.js"></script>
    <script>
        let jobRefreshInterval;

        // Function to refresh job list via AJAX
        function refreshJobs() {
            fetch('/api/jobs')
                .then(response => response.json())
                .then(data => {
                    if (data.jobs) {
                        updateJobsList(data.jobs);
                    }
                })
                .catch(error => {
                    console.log('Error refreshing jobs:', error);
                });
        }

        // Function to update the jobs list in the DOM
        function updateJobsList(jobs) {
            const jobsContainer = document.querySelector('.col-md-6:nth-child(2) .card-body');
            if (!jobs || jobs.length === 0) {
                jobsContainer.innerHTML = '<p class="text-muted">No jobs yet. Create your first job!</p>';
                return;
            }

            let jobsHTML = '';
            jobs.forEach(job => {
                const statusColor = getStatusColor(job.status);
                const formattedTime = new Date(job.created_at).toLocaleString();
                
                jobsHTML += 
                    '<div class="job-card">' +
                        '<div class="d-flex justify-content-between align-items-start">' +
                            '<div>' +
                                '<h6 class="mb-1">' +
                                    '<a href="/status/' + job.id + '" class="text-decoration-none">' + job.title + '</a>' +
                                '</h6>' +
                                '<p class="mb-1 text-muted small">' + (job.description || '') + '</p>' +
                                '<small class="text-muted">' + formattedTime + '</small>' +
                            '</div>' +
                            '<span class="badge bg-' + statusColor + ' status-badge">' +
                                job.status +
                            '</span>' +
                        '</div>' +
                        '<hr class="my-2">' +
                    '</div>';
            });
            
            jobsContainer.innerHTML = jobsHTML;
        }

        // Function to get status color
        function getStatusColor(status) {
            switch (status) {
                case 'pending': return 'warning';
                case 'processing': return 'info';
                case 'completed': return 'success';
                case 'failed': return 'danger';
                default: return 'secondary';
            }
        }

        // Function to show message
        function showMessage(message, type = 'info') {
            const messageDiv = document.getElementById('jobCreateMessage');
            messageDiv.className = 'alert alert-' + type;
            messageDiv.innerHTML = message;
            messageDiv.style.display = 'block';
            
            // Hide message after 5 seconds
            setTimeout(() => {
                messageDiv.style.display = 'none';
            }, 5000);
        }

        // Handle form submission with AJAX
        document.getElementById('jobForm').addEventListener('submit', function(e) {
            e.preventDefault();
            
            const submitBtn = document.getElementById('submitBtn');
            const title = document.getElementById('title').value;
            const description = document.getElementById('description').value;
            
            // Disable button and show loading
            submitBtn.disabled = true;
            submitBtn.innerHTML = '<span class="spinner-border spinner-border-sm me-2"></span>Creating job...';
            
            // Create job via AJAX
            fetch('/api/jobs', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    title: title,
                    description: description
                })
            })
            .then(response => response.json())
            .then(data => {
                if (data.success) {
                    showMessage('Job "' + data.job.title + '" created successfully! <a href="/status/' + data.job.id + '">View status</a>', 'success');
                    
                    // Clear form
                    document.getElementById('title').value = '';
                    document.getElementById('description').value = '';
                    
                    // Refresh jobs list immediately
                    refreshJobs();
                } else {
                    showMessage(data.error || 'Failed to create job', 'danger');
                }
            })
            .catch(error => {
                console.error('Error creating job:', error);
                showMessage('Network error. Please try again.', 'danger');
            })
            .finally(() => {
                // Re-enable button
                submitBtn.disabled = false;
                submitBtn.innerHTML = 'Submit Job';
            });
        });

        // Check system status
        fetch('/api/status')
            .then(response => response.json())
            .then(data => {
                const statusDiv = document.getElementById('system-status');
                const statusClass = data.status === 'healthy' ? 'text-success' : 'text-danger';
                const rabbitStatus = data.rabbitmq === 'connected' ? 'Connected' : 'Disconnected';
                
                statusDiv.innerHTML = '<div class="' + statusClass + '">' +
                    '<strong>API Server:</strong> ' + data.status + '<br>' +
                    '<strong>RabbitMQ:</strong> ' + rabbitStatus + '<br>' +
                    '<small>Last checked: ' + new Date(data.timestamp).toLocaleString() + '</small>' +
                    '</div>';
            })
            .catch(error => {
                document.getElementById('system-status').innerHTML = 
                    '<div class="text-danger">Failed to check system status</div>';
            });

        // Auto-refresh job list every 3 seconds
        jobRefreshInterval = setInterval(refreshJobs, 3000);
        
        // Initial load of jobs
        refreshJobs();
    </script>
</body>
</html>
{{end}}
`

const jobStatusTemplate = `
{{define "job-status"}}
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.Title}}</title>
    <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.1.3/dist/css/bootstrap.min.css" rel="stylesheet">
</head>
<body>
    <nav class="navbar navbar-dark bg-dark">
        <div class="container">
            <a class="navbar-brand" href="/">
                <strong>Microservices Demo</strong>
            </a>
            <span class="navbar-text">
                Go + RabbitMQ
            </span>
        </div>
    </nav>

    <div class="container mt-4">
        <div class="row">
            <div class="col-md-8 offset-md-2">
                <div class="card">
                    <div class="card-header d-flex justify-content-between align-items-center">
                        <h5 class="card-title mb-0">Job Details</h5>
                        <a href="/" class="btn btn-sm btn-outline-secondary">← Back to Home</a>
                    </div>
                    <div class="card-body">
                        <div class="row">
                            <div class="col-md-6">
                                <h6>{{.Job.Title}}</h6>
                                <p class="text-muted">{{.Job.Description}}</p>
                            </div>
                            <div class="col-md-6 text-end">
                                <span class="badge bg-{{statusColor .Job.Status}} fs-6">
                                    {{.Job.Status}}
                                </span>
                            </div>
                        </div>
                        
                        <hr>
                        
                        <div class="row">
                            <div class="col-md-4">
                                <strong>Job ID:</strong><br>
                                <code>{{.Job.ID}}</code>
                            </div>
                            <div class="col-md-4">
                                <strong>Created:</strong><br>
                                {{formatTime .Job.CreatedAt}}
                            </div>
                            <div class="col-md-4">
                                {{if .Job.CompletedAt}}
                                <strong>Completed:</strong><br>
                                {{formatTime .Job.CompletedAt}}
                                {{else}}
                                <strong>Duration:</strong><br>
                                <span class="text-muted">In progress...</span>
                                {{end}}
                            </div>
                        </div>
                        
                        {{if .Job.Result}}
                        <hr>
                        <div class="alert alert-success">
                            <strong>Result:</strong>
                            <div class="job-result">{{.Job.Result}}</div>
                        </div>
                        {{else}}
                        <div class="job-result"></div>
                        {{end}}
                        
                        {{if .Job.Error}}
                        <hr>
                        <div class="alert alert-danger">
                            <strong>Error:</strong><br>
                            {{.Job.Error}}
                        </div>
                        {{end}}
                        
                        {{if .Job.Sources}}
                        <hr>
                        <div class="sources-section">
                            <strong>Sources:</strong>
                            <ul class="sources-list">
                                {{range .Job.Sources}}
                                <li><a href="{{.}}" target="_blank" rel="noopener">{{.}}</a></li>
                                {{end}}
                            </ul>
                        </div>
                        {{else}}
                        <div class="sources-section" style="display: none;">
                            <strong>Sources:</strong>
                            <ul class="sources-list"></ul>
                        </div>
                        {{end}}
                        
                        {{if .Job.Confidence}}
                        <hr>
                        <div class="confidence-section">
                            <small class="confidence-text text-muted">Confidence: {{printf "%.0f%%" (multiply .Job.Confidence 100)}}</small>
                            <div class="progress confidence-bar mt-1">
                                <div class="progress-bar" style="width: {{printf "%.0f%%" (multiply .Job.Confidence 100)}}"></div>
                            </div>
                        </div>
                        {{else}}
                        <div class="confidence-section" style="display: none;">
                            <small class="confidence-text text-muted"></small>
                            <div class="progress confidence-bar mt-1">
                                <div class="progress-bar" style="width: 0%"></div>
                            </div>
                        </div>
                        {{end}}
                        
                        {{if eq .Job.Status "pending"}}
                        <div class="alert alert-info">
                            <div class="d-flex align-items-center">
                                <div class="spinner-border spinner-border-sm me-2" role="status">
                                    <span class="visually-hidden">Loading...</span>
                                </div>
                                Job is waiting to be processed...
                            </div>
                        </div>
                        {{else if eq .Job.Status "processing"}}
                        <div class="alert alert-warning">
                            <div class="d-flex align-items-center">
                                <div class="spinner-border spinner-border-sm me-2" role="status">
                                    <span class="visually-hidden">Loading...</span>
                                </div>
                                Job is currently being processed...
                            </div>
                        </div>
                        {{end}}
                    </div>
                </div>
                
                <div class="auto-refresh mt-3">
                    <small class="text-muted">
                        <em>Status updates automatically every 3 seconds (no page reload required).</em>
                    </small>
                </div>
            </div>
        </div>
    </div>

    <script src="https://cdn.jsdelivr.net/npm/bootstrap@5.1.3/dist/js/bootstrap.bundle.min.js"></script>
    <script>
        // Auto-refresh job status via AJAX instead of page reload
        if (window.location.pathname.includes('/status/')) {
            const jobId = window.location.pathname.split('/').pop();
            let statusRefreshInterval;
            
            function refreshJobStatus() {
                fetch('/api/jobs/' + jobId)
                    .then(response => response.json())
                    .then(job => {
                        // Update job status badge
                        const statusBadge = document.querySelector('.badge');
                        if (statusBadge) {
                            statusBadge.className = 'badge bg-' + getStatusColor(job.status) + ' fs-6';
                            statusBadge.textContent = job.status;
                        }
                        
                        // Update confidence if present
                        if (job.confidence) {
                            const confidenceText = document.querySelector('.confidence-text');
                            const confidenceBar = document.querySelector('.confidence-bar .progress-bar');
                            const confidenceSection = document.querySelector('.confidence-section');
                            if (confidenceText && confidenceBar && confidenceSection) {
                                const confidencePercent = Math.round(job.confidence * 100);
                                confidenceText.textContent = 'Confidence: ' + confidencePercent + '%';
                                confidenceBar.style.width = confidencePercent + '%';
                                confidenceSection.style.display = 'block';
                            }
                        }
                        
                        // Update result content if job is completed
                        if (job.status === 'completed' && job.result) {
                            const resultDiv = document.querySelector('.job-result');
                            if (resultDiv && !resultDiv.innerHTML.trim()) {
                                resultDiv.innerHTML = '<div class="research-result">' + 
                                    job.result.replace(/\n/g, '<br>') + '</div>';
                            }
                        }
                        
                        // Update sources if available
                        if (job.sources && job.sources.length > 0) {
                            const sourcesDiv = document.querySelector('.sources-list');
                            const sourcesSection = document.querySelector('.sources-section');
                            if (sourcesDiv && sourcesSection && !sourcesDiv.innerHTML.trim()) {
                                let sourcesHTML = '';
                                job.sources.forEach(source => {
                                    sourcesHTML += '<li><a href="' + source + '" target="_blank" rel="noopener">' + 
                                        source + '</a></li>';
                                });
                                sourcesDiv.innerHTML = sourcesHTML;
                                sourcesSection.style.display = 'block';
                            }
                        }
                        
                        // Stop refreshing if job is completed or failed
                        if (job.status === 'completed' || job.status === 'failed') {
                            clearInterval(statusRefreshInterval);
                            // Update the auto-refresh message
                            const autoRefreshDiv = document.querySelector('.auto-refresh em');
                            if (autoRefreshDiv) {
                                autoRefreshDiv.textContent = 'Job ' + job.status + '. Auto-refresh stopped.';
                            }
                        }
                    })
                    .catch(error => {
                        console.log('Error refreshing job status:', error);
                    });
            }
            
            function getStatusColor(status) {
                switch(status) {
                    case 'completed': return 'success';
                    case 'failed': return 'danger';
                    case 'processing': return 'warning';
                    case 'pending': return 'secondary';
                    default: return 'primary';
                }
            }
            
            // Start auto-refresh every 3 seconds, but only if job is not completed
            const currentStatus = document.querySelector('.badge').textContent.toLowerCase().trim();
            if (currentStatus !== 'completed' && currentStatus !== 'failed') {
                statusRefreshInterval = setInterval(refreshJobStatus, 3000);
            }
        }
    </script>
</body>
</html>
{{end}}
`
