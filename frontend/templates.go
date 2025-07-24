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
        .research-card {
            margin-bottom: 1rem;
        }
        .status-badge {
            font-size: 0.8rem;
        }
        .confidence-bar {
            height: 6px;
        }
        .mcp-service {
            display: inline-block;
            margin: 2px;
            padding: 2px 8px;
            font-size: 0.75rem;
            border-radius: 12px;
        }
        .sources-list {
            max-height: 100px;
            overflow-y: auto;
        }
    </style>
</head>
<body>
    <nav class="navbar navbar-dark bg-dark">
        <div class="container">
            <a class="navbar-brand" href="/">
                <strong>AI Research Agent</strong>
            </a>
            <span class="navbar-text">
                Dapr + Ollama + MCP
            </span>
        </div>
    </nav>

    <div class="container mt-4">
        <div class="row">
            <div class="col-md-6">
                <div class="card">
                    <div class="card-header">
                        <h5 class="card-title">New Research Request</h5>
                    </div>
                    <div class="card-body">
                        <form id="researchForm" action="/submit" method="POST">
                            <div class="mb-3">
                                <label for="title" class="form-label">Job Title *</label>
                                <input type="text" class="form-control" id="title" name="title" required
                                    placeholder="e.g., 'Market Analysis for AI Tools'">
                            </div>
                            <div class="mb-3">
                                <label for="query" class="form-label">Instructions *</label>
                                <textarea class="form-control" id="query" name="query" rows="3" required
                                    placeholder="What would you like to research? Be specific about your information needs."></textarea>
                            </div>
                            <div class="mb-3">
                                <label for="research_type" class="form-label">Research Type</label>
                                <select class="form-select" id="research_type" name="research_type">
                                    <option value="general">General Research</option>
                                    <option value="technical">Technical Analysis</option>
                                    <option value="market">Market Research</option>
                                    <option value="competitive">Competitive Analysis</option>
                                    <option value="code">Code & Development</option>
                                    <option value="data">Data Analysis</option>
                                </select>
                            </div>
                            <div class="mb-3">
                                <label class="form-label">MCP Services to Use</label>
                                <div>
                                    <div class="form-check form-check-inline">
                                        <input class="form-check-input" type="checkbox" id="mcp_web" name="mcp_services" value="web" checked>
                                        <label class="form-check-label" for="mcp_web">Web Search</label>
                                    </div>
                                    <div class="form-check form-check-inline">
                                        <input class="form-check-input" type="checkbox" id="mcp_github" name="mcp_services" value="github">
                                        <label class="form-check-label" for="mcp_github">GitHub</label>
                                    </div>
                                    <div class="form-check form-check-inline">
                                        <input class="form-check-input" type="checkbox" id="mcp_files" name="mcp_services" value="files">
                                        <label class="form-check-label" for="mcp_files">Local Files</label>
                                    </div>
                                </div>
                                <small class="form-text text-muted">Select which MCP services to use for data gathering</small>
                            </div>
                            <button type="submit" class="btn btn-primary" id="submitBtn">Start Research</button>
                        </form>
                        
                        <!-- Success/Error messages -->
                        <div id="researchCreateMessage" class="mt-3" style="display: none;"></div>
                    </div>
                </div>
            </div>
            
            <div class="col-md-6">
                <div class="card">
                    <div class="card-header d-flex justify-content-between align-items-center">
                        <h5 class="card-title mb-0">Recent Research</h5>
                        <button class="btn btn-sm btn-outline-secondary" onclick="window.location.reload()">
                            Refresh
                        </button>
                    </div>
                    <div class="card-body">
                        {{if .Jobs}}
                            {{range .Jobs}}
                            <div class="research-card">
                                <div class="d-flex justify-content-between align-items-start">
                                    <div class="flex-grow-1">
                                        <h6 class="mb-1">
                                            <a href="/status/{{.ID}}" class="text-decoration-none">{{.Title}}</a>
                                        </h6>
                                        <p class="mb-1 text-muted small">{{.Query}}</p>
                                        {{if .MCPServices}}
                                        <div class="mb-2">
                                            {{range .MCPServices}}
                                            <span class="mcp-service bg-light text-dark border">{{.}}</span>
                                            {{end}}
                                        </div>
                                        {{end}}
                                        {{if .Confidence}}
                                        <div class="mb-2">
                                            <small class="text-muted">Confidence: {{printf "%.0f%%" (multiply .Confidence 100)}}</small>
                                            <div class="progress confidence-bar">
                                                <div class="progress-bar" style="width: {{printf "%.0f%%" (multiply .Confidence 100)}}"></div>
                                            </div>
                                        </div>
                                        {{end}}
                                        <small class="text-muted">{{formatTime .CreatedAt}}</small>
                                        {{if .TokensUsed}}
                                        <small class="text-muted"> • {{.TokensUsed}} tokens</small>
                                        {{end}}
                                    </div>
                                    <span class="badge bg-{{statusColor .Status}} status-badge">
                                        {{.Status}}
                                    </span>
                                </div>
                                <hr class="my-2">
                            </div>
                            {{end}}
                        {{else}}
                            <p class="text-muted">No research requests yet. Start your first research!</p>
                        {{end}}
                    </div>
                </div>
            </div>
        </div>

        <div class="row mt-4">
            <div class="col-12">
                <div class="card">
                    <div class="card-header">
                        <h5 class="card-title">AI Research Agent Status</h5>
                    </div>
                    <div class="card-body">
                        <div id="system-status">
                            <div class="spinner-border spinner-border-sm" role="status">
                                <span class="visually-hidden">Loading...</span>
                            </div>
                            Checking AI agent status...
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>

    <script src="https://cdn.jsdelivr.net/npm/bootstrap@5.1.3/dist/js/bootstrap.bundle.min.js"></script>
    <script src="https://cdn.jsdelivr.net/npm/marked@4.3.0/marked.min.js"></script>
    <script>
        let researchRefreshInterval;

        // Function to refresh research list via AJAX
        function refreshResearch() {
            fetch('/api/jobs')
                .then(response => response.json())
                .then(data => {
                    if (data.jobs) {
                        updateResearchList(data.jobs);
                    }
                })
                .catch(error => {
                    console.log('Error refreshing research:', error);
                });
        }

        // Function to update the research list in the DOM
        function updateResearchList(research) {
            const researchContainer = document.querySelector('.col-md-6:nth-child(2) .card-body');
            if (!research || research.length === 0) {
                researchContainer.innerHTML = '<p class="text-muted">No research requests yet. Start your first research!</p>';
                return;
            }

            let researchHTML = '';
            research.forEach(item => {
                const statusColor = getStatusColor(item.status);
                const formattedTime = new Date(item.created_at).toLocaleString();
                
                let mcpServicesHTML = '';
                if (item.mcp_services && item.mcp_services.length > 0) {
                    mcpServicesHTML = '<div class="mb-2">';
                    item.mcp_services.forEach(service => {
                        mcpServicesHTML += '<span class="mcp-service bg-light text-dark border">' + service + '</span>';
                    });
                    mcpServicesHTML += '</div>';
                }

                let confidenceHTML = '';
                if (item.confidence) {
                    const confidencePercent = Math.round(item.confidence * 100);
                    confidenceHTML = 
                        '<div class="mb-2">' +
                        '<small class="text-muted">Confidence: ' + confidencePercent + '%</small>' +
                        '<div class="progress confidence-bar">' +
                        '<div class="progress-bar" style="width: ' + confidencePercent + '%"></div>' +
                        '</div>' +
                        '</div>';
                }

                let tokensHTML = '';
                if (item.tokens_used) {
                    tokensHTML = ' • ' + item.tokens_used + ' tokens';
                }
                
                researchHTML += 
                    '<div class="research-card">' +
                        '<div class="d-flex justify-content-between align-items-start">' +
                            '<div class="flex-grow-1">' +
                                '<h6 class="mb-1">' +
                                    '<a href="/status/' + item.id + '" class="text-decoration-none">' + item.title + '</a>' +
                                '</h6>' +
                                '<p class="mb-1 text-muted small">' + (item.query || '') + '</p>' +
                                mcpServicesHTML +
                                confidenceHTML +
                                '<small class="text-muted">' + formattedTime + tokensHTML + '</small>' +
                            '</div>' +
                            '<span class="badge bg-' + statusColor + ' status-badge">' +
                                item.status +
                            '</span>' +
                        '</div>' +
                        '<hr class="my-2">' +
                    '</div>';
            });
            
            researchContainer.innerHTML = researchHTML;
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
            const messageDiv = document.getElementById('researchCreateMessage');
            messageDiv.className = 'alert alert-' + type;
            messageDiv.innerHTML = message;
            messageDiv.style.display = 'block';
            
            // Hide message after 5 seconds
            setTimeout(() => {
                messageDiv.style.display = 'none';
            }, 5000);
        }

        // Handle form submission with AJAX
        document.getElementById('researchForm').addEventListener('submit', function(e) {
            e.preventDefault();
            
            const submitBtn = document.getElementById('submitBtn');
            const title = document.getElementById('title').value;
            const query = document.getElementById('query').value;
            const researchType = document.getElementById('research_type').value;
            
            // Get selected MCP services
            const mcpServices = [];
            document.querySelectorAll('input[name="mcp_services"]:checked').forEach(checkbox => {
                mcpServices.push(checkbox.value);
            });
            
            // Disable button and show loading
            submitBtn.disabled = true;
            submitBtn.innerHTML = '<span class="spinner-border spinner-border-sm me-2"></span>Starting research...';
            
            // Create research request via AJAX
            fetch('/api/jobs', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    title: title,
                    query: query,
                    research_type: researchType,
                    mcp_services: mcpServices
                })
            })
            .then(response => response.json())
            .then(data => {
                if (data.success) {
                    showMessage('Research "' + data.job.title + '" started successfully! <a href="/status/' + data.job.id + '">View progress</a>', 'success');
                    
                    // Clear form
                    document.getElementById('title').value = '';
                    document.getElementById('query').value = '';
                    document.getElementById('research_type').selectedIndex = 0;
                    document.querySelectorAll('input[name="mcp_services"]').forEach(checkbox => {
                        checkbox.checked = checkbox.value === 'web'; // Reset to just web checked
                    });
                    
                    // Refresh research list immediately
                    refreshResearch();
                } else {
                    showMessage(data.error || 'Failed to start research', 'danger');
                }
            })
            .catch(error => {
                console.error('Error starting research:', error);
                showMessage('Network error. Please try again.', 'danger');
            })
            .finally(() => {
                // Re-enable button
                submitBtn.disabled = false;
                submitBtn.innerHTML = 'Start Research';
            });
        });

        // Check AI agent status
        fetch('/api/status')
            .then(response => response.json())
            .then(data => {
                const statusDiv = document.getElementById('system-status');
                const statusClass = data.status === 'healthy' ? 'text-success' : 'text-danger';
                const rabbitStatus = data.rabbitmq === 'connected' ? 'Connected' : 'Disconnected';
                
                let statusHTML = '<div class="' + statusClass + '">' +
                    '<strong>API Server:</strong> ' + data.status + '<br>' +
                    '<strong>RabbitMQ:</strong> ' + rabbitStatus + '<br>';
                
                // Add Ollama information
                if (data.ollama) {
                    statusHTML += '<strong>Ollama AI:</strong> ' + data.ollama.model + 
                                 ' <small>(' + data.ollama.endpoint + ')</small><br>';
                } else {
                    statusHTML += '<strong>Ollama AI:</strong> Unknown<br>';
                }
                
                // Add MCP information
                if (data.mcp) {
                    if (data.mcp.test_mode) {
                        statusHTML += '<strong>MCP Services:</strong> Test Mode (Simulated)<br>';
                    } else {
                        statusHTML += '<strong>MCP Services:</strong> Production Mode<br>';
                        if (data.mcp.endpoints) {
                            statusHTML += '<div class="ms-3 small">';
                            statusHTML += '• Web Search: ' + data.mcp.endpoints.web_search + '<br>';
                            statusHTML += '• GitHub: ' + data.mcp.endpoints.github + '<br>';
                            statusHTML += '• Files: ' + data.mcp.endpoints.files + '<br>';
                            statusHTML += '</div>';
                        }
                    }
                } else {
                    statusHTML += '<strong>MCP Services:</strong> Unknown<br>';
                }
                
                statusHTML += '<small>Last checked: ' + new Date(data.timestamp).toLocaleString() + '</small>' +
                    '</div>';
                
                statusDiv.innerHTML = statusHTML;
            })
            .catch(error => {
                document.getElementById('system-status').innerHTML = 
                    '<div class="text-warning">Unable to check AI agent status</div>';
            });

        // Auto-refresh research list every 3 seconds
        researchRefreshInterval = setInterval(refreshResearch, 3000);
        
        // Initial load of research
        refreshResearch();
    </script>
</body>
</html>
{{end}}
`

const researchStatusTemplate = `
{{define "research-status"}}
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.Title}}</title>
    <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.1.3/dist/css/bootstrap.min.css" rel="stylesheet">
    <style>
        .confidence-bar {
            height: 10px;
        }
        .mcp-service {
            display: inline-block;
            margin: 2px;
            padding: 4px 12px;
            font-size: 0.8rem;
            border-radius: 15px;
        }
        .sources-list {
            max-height: 200px;
            overflow-y: auto;
        }
        .research-result {
            line-height: 1.6;
        }
        .research-result h1, .research-result h2, .research-result h3, 
        .research-result h4, .research-result h5, .research-result h6 {
            margin-top: 1.5rem;
            margin-bottom: 1rem;
            color: #495057;
        }
        .research-result h1 { font-size: 1.5rem; }
        .research-result h2 { font-size: 1.3rem; }
        .research-result h3 { font-size: 1.1rem; }
        .research-result p {
            margin-bottom: 1rem;
        }
        .research-result ul, .research-result ol {
            margin-bottom: 1rem;
            padding-left: 1.5rem;
        }
        .research-result li {
            margin-bottom: 0.25rem;
        }
        .research-result code {
            background-color: #f8f9fa;
            padding: 0.2rem 0.4rem;
            border-radius: 0.25rem;
            font-size: 0.9em;
            color: #d63384;
        }
        .research-result pre {
            background-color: #f8f9fa;
            padding: 1rem;
            border-radius: 0.375rem;
            overflow-x: auto;
            margin-bottom: 1rem;
        }
        .research-result pre code {
            background: none;
            padding: 0;
            color: #212529;
        }
        .research-result blockquote {
            border-left: 4px solid #dee2e6;
            padding-left: 1rem;
            margin: 1rem 0;
            font-style: italic;
            color: #6c757d;
        }
        .research-result table {
            margin: 1rem 0;
        }
    </style>
</head>
<body>
    <nav class="navbar navbar-dark bg-dark">
        <div class="container">
            <a class="navbar-brand" href="/">
                <strong>AI Research Agent</strong>
            </a>
            <span class="navbar-text">
                Dapr + Ollama + MCP
            </span>
        </div>
    </nav>

    <div class="container mt-4">
        <div class="row">
            <div class="col-md-10 offset-md-1">
                <div class="card">
                    <div class="card-header d-flex justify-content-between align-items-center">
                        <h5 class="card-title mb-0">Research Details</h5>
                        <a href="/" class="btn btn-sm btn-outline-secondary">← Back to Home</a>
                    </div>
                    <div class="card-body">
                        <div class="row">
                            <div class="col-md-8">
                                <h6>{{.Job.Title}}</h6>
                                <p class="text-muted">{{.Job.Query}}</p>
                                {{if .Job.MCPServices}}
                                <div class="mb-3">
                                    <small class="text-muted">MCP Services Used:</small><br>
                                    {{range .Job.MCPServices}}
                                    <span class="mcp-service bg-light text-dark border">{{.}}</span>
                                    {{end}}
                                </div>
                                {{end}}
                            </div>
                            <div class="col-md-4 text-end">
                                <span class="badge bg-{{statusColor .Job.Status}} fs-6">
                                    {{.Job.Status}}
                                </span>
                                {{if .Job.Confidence}}
                                <div class="mt-2">
                                    <small class="text-muted">Confidence: {{printf "%.0f%%" (multiply .Job.Confidence 100)}}</small>
                                    <div class="progress confidence-bar">
                                        <div class="progress-bar" style="width: {{printf "%.0f%%" (multiply .Job.Confidence 100)}}"></div>
                                    </div>
                                </div>
                                {{end}}
                            </div>
                        </div>
                        
                        <hr>
                        
                        <div class="row">
                            <div class="col-md-3">
                                <strong>Research ID:</strong><br>
                                <code>{{.Job.ID}}</code>
                            </div>
                            <div class="col-md-3">
                                <strong>Created:</strong><br>
                                {{formatTime .Job.CreatedAt}}
                            </div>
                            <div class="col-md-3">
                                {{if .Job.CompletedAt}}
                                <strong>Completed:</strong><br>
                                {{formatTime .Job.CompletedAt}}
                                {{else}}
                                <strong>Duration:</strong><br>
                                <span class="text-muted">In progress...</span>
                                {{end}}
                            </div>
                            <div class="col-md-3">
                                {{if .Job.TokensUsed}}
                                <strong>Tokens Used:</strong><br>
                                {{.Job.TokensUsed}}
                                {{else}}
                                <strong>Research Type:</strong><br>
                                {{.Job.ResearchType}}
                                {{end}}
                            </div>
                        </div>
                        
                        {{if .Job.Result}}
                        <hr>
                        <div class="card">
                            <div class="card-header">
                                <h6 class="mb-0">Research Results</h6>
                            </div>
                            <div class="card-body">
                                <div class="research-result" id="research-result-content">{{.Job.Result}}</div>
                            </div>
                        </div>
                        {{end}}
                        
                        {{if .Job.Sources}}
                        <hr>
                        <div class="card">
                            <div class="card-header">
                                <h6 class="mb-0">Sources ({{len .Job.Sources}})</h6>
                            </div>
                            <div class="card-body sources-list">
                                {{range $index, $source := .Job.Sources}}
                                <div class="mb-2">
                                    <strong>{{add $index 1}}.</strong>
                                    {{if hasPrefix $source "http"}}
                                    <a href="{{$source}}" target="_blank" class="text-decoration-none">{{$source}}</a>
                                    {{else}}
                                    <code>{{$source}}</code>
                                    {{end}}
                                </div>
                                {{end}}
                            </div>
                        </div>
                                {{end}}
                            </div>
                        </div>
                        
                        {{if .Job.Result}}
                        <hr>
                        <div class="alert alert-success">
                            <strong>Result:</strong><br>
                            {{.Job.Result}}
                        </div>
                        {{end}}
                        
                        {{if .Job.Error}}
                        <hr>
                        <div class="alert alert-danger">
                            <strong>Research Failed:</strong><br>
                            {{.Job.Error}}
                        </div>
                        {{end}}
                        
                        {{if eq .Job.Status "pending"}}
                        <div class="alert alert-info">
                            <div class="d-flex align-items-center">
                                <div class="spinner-border spinner-border-sm me-2" role="status">
                                    <span class="visually-hidden">Loading...</span>
                                </div>
                                Research request is waiting to be processed by AI agent...
                            </div>
                        </div>
                        {{else if eq .Job.Status "processing"}}
                        <div class="alert alert-warning">
                            <div class="d-flex align-items-center">
                                <div class="spinner-border spinner-border-sm me-2" role="status">
                                    <span class="visually-hidden">Loading...</span>
                                </div>
                                AI agent is gathering information and analyzing data...
                            </div>
                        </div>
                        {{end}}
                    </div>
                </div>
                
                <div class="auto-refresh mt-3">
                    <small class="text-muted">
                        <em>This page auto-refreshes every 3 seconds while research is in progress.</em>
                    </small>
                </div>
            </div>
        </div>
    </div>

    <script src="https://cdn.jsdelivr.net/npm/bootstrap@5.1.3/dist/js/bootstrap.bundle.min.js"></script>
    <script src="https://cdn.jsdelivr.net/npm/marked@4.3.0/marked.min.js"></script>
    <script>
        // Render markdown in research results
        document.addEventListener('DOMContentLoaded', function() {
            const resultElement = document.getElementById('research-result-content');
            if (resultElement && resultElement.textContent.trim()) {
                const markdownText = resultElement.textContent;
                // Configure marked with safe defaults
                marked.setOptions({
                    breaks: true,        // Convert line breaks to <br>
                    gfm: true,          // GitHub Flavored Markdown
                    sanitize: false,    // We trust the AI output, but you might want to sanitize
                    smartLists: true,   // Use smarter list behavior
                    smartypants: true   // Use smart quotes and dashes
                });
                
                // Render markdown to HTML
                const htmlContent = marked.parse(markdownText);
                resultElement.innerHTML = htmlContent;
                
                // Add some custom styling to the rendered content
                resultElement.style.lineHeight = '1.6';
                resultElement.style.fontSize = '15px';
                
                // Style tables if any
                const tables = resultElement.querySelectorAll('table');
                tables.forEach(table => {
                    table.classList.add('table', 'table-striped', 'table-sm');
                    table.style.marginTop = '1rem';
                });
                
                // Style code blocks
                const codeBlocks = resultElement.querySelectorAll('pre code');
                codeBlocks.forEach(block => {
                    block.style.backgroundColor = '#f8f9fa';
                    block.style.padding = '1rem';
                    block.style.borderRadius = '0.375rem';
                    block.style.fontSize = '14px';
                });
                
                // Style blockquotes
                const blockquotes = resultElement.querySelectorAll('blockquote');
                blockquotes.forEach(quote => {
                    quote.style.borderLeft = '4px solid #dee2e6';
                    quote.style.paddingLeft = '1rem';
                    quote.style.marginLeft = '0';
                    quote.style.fontStyle = 'italic';
                    quote.style.color = '#6c757d';
                });
            }
        });

        // Auto-refresh job status via AJAX every 3 seconds
        if (window.location.pathname.includes('/status/')) {
            const jobId = window.location.pathname.split('/status/')[1];
            
            function refreshJobStatus() {
                fetch('/api/jobs/' + jobId)
                    .then(response => response.json())
                    .then(job => {
                        // Update status badge
                        const statusBadge = document.querySelector('.badge');
                        if (statusBadge) {
                            statusBadge.className = 'badge bg-' + getStatusColor(job.status);
                            statusBadge.textContent = job.status.charAt(0).toUpperCase() + job.status.slice(1);
                        }
                        
                        // Update confidence if present
                        const confidenceElement = document.querySelector('.progress-bar');
                        if (confidenceElement && job.confidence) {
                            confidenceElement.style.width = (job.confidence * 100) + '%';
                        }
                        
                        // Update result if completed
                        const resultContainer = document.getElementById('research-result-content');
                        if (job.result && resultContainer) {
                            // Only update if content has changed
                            if (resultContainer.textContent !== job.result) {
                                resultContainer.textContent = job.result;
                                // Re-render markdown
                                if (typeof marked !== 'undefined') {
                                    resultContainer.innerHTML = marked.parse(job.result);
                                }
                            }
                        }
                        
                        // Update sources
                        const sourcesContainer = document.querySelector('.sources-list');
                        if (job.sources && sourcesContainer && job.sources.length > 0) {
                            // Update sources count in header
                            const sourcesHeader = document.querySelector('.card-header h6');
                            if (sourcesHeader) {
                                sourcesHeader.textContent = 'Sources (' + job.sources.length + ')';
                            }
                            
                            // Update sources list
                            let sourcesList = '';
                            job.sources.forEach(function(source, index) {
                                sourcesList += '<div class="mb-2"><strong>' + (index + 1) + '.</strong> <a href="' + source + '" target="_blank" class="text-decoration-none">' + source + '</a></div>';
                            });
                            sourcesContainer.innerHTML = sourcesList;
                        }
                        
                        // Update completion time if job is done
                        if (job.completed_at) {
                            const durationElements = document.querySelectorAll('strong');
                            durationElements.forEach(function(el) {
                                if (el.textContent === 'Duration:') {
                                    const nextElement = el.parentNode.querySelector('span, br + *');
                                    if (nextElement && nextElement.textContent.includes('In progress')) {
                                        const completedTime = new Date(job.completed_at).toLocaleString();
                                        el.textContent = 'Completed:';
                                        el.nextSibling.textContent = '\n' + completedTime;
                                    }
                                }
                            });
                        }
                        
                        // Stop refreshing if job is completed or failed
                        if (job.status === 'completed' || job.status === 'failed') {
                            clearInterval(statusRefreshInterval);
                            
                            // Update auto-refresh message
                            const autoRefreshMsg = document.querySelector('.auto-refresh');
                            if (autoRefreshMsg) {
                                autoRefreshMsg.innerHTML = '<small class="text-muted"><em>Research completed. Auto-refresh stopped.</em></small>';
                            }
                        }
                    })
                    .catch(function(error) {
                        console.error('Error refreshing job status:', error);
                    });
            }
            
            function getStatusColor(status) {
                switch (status) {
                    case 'pending': return 'warning';
                    case 'processing': return 'info';
                    case 'completed': return 'success';
                    case 'failed': return 'danger';
                    default: return 'secondary';
                }
            }
            
            // Start auto-refresh interval
            const statusRefreshInterval = setInterval(refreshJobStatus, 3000);
            
            // Initial refresh
            refreshJobStatus();
        }
    </script>
</body>
</html>
{{end}}
`
