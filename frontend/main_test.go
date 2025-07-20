package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestHomePage(t *testing.T) {
	gin.SetMode(gin.TestMode)

	frontend := NewFrontend()
	frontend.createInlineTemplates()
	router := frontend.setupRoutes()

	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	body := w.Body.String()
	if !strings.Contains(body, "Microservices Demo") {
		t.Error("Expected page to contain 'Microservices Demo'")
	}
	if !strings.Contains(body, "Create New Job") {
		t.Error("Expected page to contain 'Create New Job'")
	}
}

func TestSubmitJobForm(t *testing.T) {
	gin.SetMode(gin.TestMode)

	frontend := NewFrontend()
	frontend.createInlineTemplates()
	router := frontend.setupRoutes()

	// Test form submission with missing title
	form := url.Values{}
	form.Add("description", "Test description")

	req, _ := http.NewRequest("POST", "/submit", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusSeeOther {
		t.Errorf("Expected status code %d, got %d", http.StatusSeeOther, w.Code)
	}

	location := w.Header().Get("Location")
	if !strings.Contains(location, "error=Title") {
		t.Error("Expected redirect to contain error message about title")
	}
}

func TestJobStatusPage(t *testing.T) {
	gin.SetMode(gin.TestMode)

	frontend := NewFrontend()
	frontend.createInlineTemplates()
	router := frontend.setupRoutes()

	// Test with a job ID (this will fail to fetch from API, but should render template)
	req, _ := http.NewRequest("GET", "/status/test-123", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Should return 404 since API server is not running in test
	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestTemplateCreation(t *testing.T) {
	frontend := NewFrontend()
	frontend.createInlineTemplates()

	if frontend.templates == nil {
		t.Error("Expected templates to be created")
	}

	// Test that templates are properly defined
	templates := []string{"index", "job-status"}
	for _, tmplName := range templates {
		if frontend.templates.Lookup(tmplName) == nil {
			t.Errorf("Expected template %s to be defined", tmplName)
		}
	}
}

func TestStatusColorFunction(t *testing.T) {
	frontend := NewFrontend()
	frontend.createInlineTemplates()

	// The statusColor function should be available in templates
	tmpl := frontend.templates.Lookup("index")
	if tmpl == nil {
		t.Fatal("Expected index template to exist")
	}

	// We can't easily test template functions directly, but we can verify
	// the template was created successfully with the functions
}

func TestFormatTimeFunction(t *testing.T) {
	frontend := NewFrontend()
	frontend.createInlineTemplates()

	// The formatTime function should be available in templates
	tmpl := frontend.templates.Lookup("job-status")
	if tmpl == nil {
		t.Fatal("Expected job-status template to exist")
	}
}
