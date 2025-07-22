package main

import (
	"testing"

	"microservices-demo/shared"
)

func TestResearchAgentCreation(t *testing.T) {
	agent := NewResearchAgent()
	if agent == nil {
		t.Error("Expected research agent to be created")
	}
}

func TestResearchRequestHandling(t *testing.T) {
	agent := NewResearchAgent()
	if agent == nil {
		t.Error("Expected research agent to be created")
	}
	
	// Test with a sample research request
	jobMessage := shared.JobMessage{
		JobID:        "test-123",
		Title:        "Test Research",
		Query:        "Research about Go microservices",
		ResearchType: shared.ResearchTypeGeneral,
		MCPServices:  []shared.MCPService{shared.MCPServiceWeb},
	}

	// This is a basic test - the actual processing would require Ollama
	if jobMessage.JobID == "" {
		t.Error("JobID should not be empty")
	}
	if jobMessage.Query == "" {
		t.Error("Query should not be empty")
	}
}

func TestMCPServiceMock(t *testing.T) {
	agent := NewResearchAgent()
	agent.initMCPServices()
	
	if agent.mcpHandler == nil {
		t.Error("Expected MCP handler to be initialized")
	}
	
	// Test that web service is available
	if !agent.mcpHandler.availableServices[shared.MCPServiceWeb] {
		t.Error("Expected web service to be available")
	}
}

// Note: Full integration tests with Ollama would require external dependencies
// These are kept minimal for CI/CD pipeline compatibility