package jira

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

// Client represents a Jira API client
type Client struct {
	BaseURL    string
	Username   string
	APIToken   string
	ProjectKey string
	HTTPClient *http.Client
}

// NewClient creates a new Jira API client
func NewClient(baseURL, username, apiToken, projectKey string) *Client {
	return &Client{
		BaseURL:    strings.TrimSuffix(baseURL, "/"),
		Username:   username,
		APIToken:   apiToken,
		ProjectKey: projectKey,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// makeRequest makes an HTTP request to the Jira API
func (c *Client) makeRequest(method, endpoint string, body interface{}) (*http.Response, error) {
	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonBody)
	}

	url := fmt.Sprintf("%s/rest/api/3/%s", c.BaseURL, endpoint)
	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// Set basic authentication
	auth := base64.StdEncoding.EncodeToString([]byte(c.Username + ":" + c.APIToken))
	req.Header.Set("Authorization", "Basic "+auth)

	log.Printf("Making %s request to: %s", method, url)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}

	return resp, nil
}

// handleResponse handles the HTTP response and checks for errors
func (c *Client) handleResponse(resp *http.Response, target interface{}) error {
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	log.Printf("Response status: %d, body length: %d", resp.StatusCode, len(body))

	if resp.StatusCode >= 400 {
		var errorResp ErrorResponse
		if err := json.Unmarshal(body, &errorResp); err != nil {
			return fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
		}
		return fmt.Errorf("Jira API error (HTTP %d): %v, %v", resp.StatusCode, errorResp.ErrorMessages, errorResp.Errors)
	}

	if target != nil && len(body) > 0 {
		if err := json.Unmarshal(body, target); err != nil {
			return fmt.Errorf("failed to unmarshal response: %w", err)
		}
	}

	return nil
}

// ListTestCases retrieves test cases from Jira
func (c *Client) ListTestCases() ([]TestCase, error) {
	log.Println("Fetching test cases from Jira...")

	// JQL query to find test cases (assuming Test issue type exists)
	jql := fmt.Sprintf("project = %s AND issuetype = Test", c.ProjectKey)
	endpoint := fmt.Sprintf("search?jql=%s&maxResults=100", jql)

	resp, err := c.makeRequest("GET", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch test cases: %w", err)
	}

	var jiraResp JiraResponse
	if err := c.handleResponse(resp, &jiraResp); err != nil {
		// If using demo credentials, return mock data
		if c.isDemoCredentials() {
			log.Println("Using demo credentials, returning mock test cases")
			return c.getMockTestCases(), nil
		}
		return nil, err
	}

	// Convert Jira issues to TestCase structs
	testCases := make([]TestCase, len(jiraResp.Issues))
	for i, issue := range jiraResp.Issues {
		testCases[i] = TestCase{
			ID:          issue.ID,
			Key:         issue.Key,
			Summary:     issue.Fields.Summary,
			Description: issue.Fields.Description,
			Status:      issue.Fields.Status.Name,
			Priority:    issue.Fields.Priority.Name,
			Labels:      issue.Fields.Labels,
			Reporter:    issue.Fields.Reporter.DisplayName,
			Assignee:    issue.Fields.Assignee.DisplayName,
		}
	}

	log.Printf("Successfully fetched %d test cases", len(testCases))
	return testCases, nil
}

// CreateTestCase creates a new test case in Jira
func (c *Client) CreateTestCase(tc *TestCase) (*TestCase, error) {
	log.Printf("Creating test case: %s", tc.Summary)

	// If using demo credentials, return mock response
	if c.isDemoCredentials() {
		log.Println("Using demo credentials, returning mock test case creation")
		return c.createMockTestCase(tc), nil
	}

	createReq := CreateIssueRequest{
		Fields: IssueFields{
			Summary:     tc.Summary,
			Description: tc.Description,
			IssueType: IssueType{
				Name: "Test", // Assuming Test issue type exists
			},
			Project: Project{
				Key: c.ProjectKey,
			},
			Labels: tc.Labels,
		},
	}

	if tc.Priority != "" {
		createReq.Fields.Priority = Priority{Name: tc.Priority}
	}

	resp, err := c.makeRequest("POST", "issue", createReq)
	if err != nil {
		return nil, fmt.Errorf("failed to create test case: %w", err)
	}

	var createResp CreateIssueResponse
	if err := c.handleResponse(resp, &createResp); err != nil {
		return nil, err
	}

	// Return the created test case with updated information
	createdTC := *tc
	createdTC.ID = createResp.ID
	createdTC.Key = createResp.Key
	createdTC.Status = "To Do"
	createdTC.CreatedDate = time.Now()

	log.Printf("Successfully created test case: %s", createdTC.Key)
	return &createdTC, nil
}

// CreateTestExecution creates a new test execution in Jira
func (c *Client) CreateTestExecution(te *TestExecution) (*TestExecution, error) {
	log.Printf("Creating test execution: %s", te.Summary)

	// If using demo credentials, return mock response
	if c.isDemoCredentials() {
		log.Println("Using demo credentials, returning mock test execution creation")
		return c.createMockTestExecution(te), nil
	}

	createReq := CreateIssueRequest{
		Fields: IssueFields{
			Summary:     te.Summary,
			Description: te.Description,
			IssueType: IssueType{
				Name: "Test Execution", // Assuming Test Execution issue type exists
			},
			Project: Project{
				Key: c.ProjectKey,
			},
		},
	}

	resp, err := c.makeRequest("POST", "issue", createReq)
	if err != nil {
		return nil, fmt.Errorf("failed to create test execution: %w", err)
	}

	var createResp CreateIssueResponse
	if err := c.handleResponse(resp, &createResp); err != nil {
		return nil, err
	}

	// Return the created test execution with updated information
	createdTE := *te
	createdTE.ID = createResp.ID
	createdTE.Key = createResp.Key
	createdTE.Status = "To Do"
	createdTE.ExecutionStatus = "TODO"
	createdTE.StartDate = time.Now()

	log.Printf("Successfully created test execution: %s", createdTE.Key)
	return &createdTE, nil
}

// GetTestExecution retrieves a test execution by key
func (c *Client) GetTestExecution(key string) (*TestExecution, error) {
	log.Printf("Fetching test execution: %s", key)

	if c.isDemoCredentials() {
		log.Println("Using demo credentials, returning mock test execution")
		return c.getMockTestExecution(key), nil
	}

	endpoint := fmt.Sprintf("issue/%s", key)
	resp, err := c.makeRequest("GET", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch test execution: %w", err)
	}

	var issue JiraIssue
	if err := c.handleResponse(resp, &issue); err != nil {
		return nil, err
	}

	testExecution := &TestExecution{
		ID:          issue.ID,
		Key:         issue.Key,
		Summary:     issue.Fields.Summary,
		Description: issue.Fields.Description,
		Status:      issue.Fields.Status.Name,
	}

	log.Printf("Successfully fetched test execution: %s", testExecution.Key)
	return testExecution, nil
}

// isDemoCredentials checks if demo credentials are being used
func (c *Client) isDemoCredentials() bool {
	return c.Username == "demo_user" || c.APIToken == "demo_token_replace_with_actual"
}

// Mock data methods for demo purposes
func (c *Client) getMockTestCases() []TestCase {
	return []TestCase{
		{
			ID:          "10001",
			Key:         "TEST-1",
			Summary:     "Login functionality test",
			Description: "Test user login with valid credentials",
			Status:      "To Do",
			Priority:    "High",
			Labels:      []string{"login", "authentication"},
			TestType:    "Manual",
			CreatedDate: time.Now().AddDate(0, 0, -7),
			Reporter:    "Demo User",
		},
		{
			ID:          "10002",
			Key:         "TEST-2",
			Summary:     "Password reset functionality",
			Description: "Test password reset flow",
			Status:      "In Progress",
			Priority:    "Medium",
			Labels:      []string{"password", "reset"},
			TestType:    "Automated",
			CreatedDate: time.Now().AddDate(0, 0, -5),
			Reporter:    "Demo User",
		},
		{
			ID:          "10003",
			Key:         "TEST-3",
			Summary:     "User registration validation",
			Description: "Test user registration with various input validations",
			Status:      "Done",
			Priority:    "Medium",
			Labels:      []string{"registration", "validation"},
			TestType:    "Manual",
			CreatedDate: time.Now().AddDate(0, 0, -3),
			Reporter:    "Demo User",
		},
	}
}

func (c *Client) createMockTestCase(tc *TestCase) *TestCase {
	mockTC := *tc
	mockTC.ID = "10004"
	mockTC.Key = "TEST-4"
	mockTC.Status = "To Do"
	mockTC.CreatedDate = time.Now()
	mockTC.Reporter = "Demo User"
	return &mockTC
}

func (c *Client) createMockTestExecution(te *TestExecution) *TestExecution {
	mockTE := *te
	mockTE.ID = "10005"
	mockTE.Key = "EXEC-1"
	mockTE.Status = "To Do"
	mockTE.ExecutionStatus = "TODO"
	mockTE.StartDate = time.Now()
	mockTE.ExecutedBy = "Demo User"
	return &mockTE
}

func (c *Client) getMockTestExecution(key string) *TestExecution {
	return &TestExecution{
		ID:              "10005",
		Key:             key,
		Summary:         "Demo Test Execution",
		Description:     "This is a demo test execution",
		Status:          "In Progress",
		TestCases:       []string{"TEST-1", "TEST-2"},
		ExecutionStatus: "EXECUTING",
		StartDate:       time.Now().AddDate(0, 0, -1),
		ExecutedBy:      "Demo User",
		Environment:     "QA",
		TestResults: []TestResult{
			{
				TestCaseKey:   "TEST-1",
				Status:        "PASS",
				Comment:       "Test passed successfully",
				ExecutionTime: 5000,
				ExecutedBy:    "Demo User",
				ExecutedOn:    time.Now(),
			},
			{
				TestCaseKey:   "TEST-2",
				Status:        "FAIL",
				Comment:       "Test failed due to timeout",
				ExecutionTime: 10000,
				ExecutedBy:    "Demo User",
				ExecutedOn:    time.Now(),
				Defects:       []string{"BUG-123"},
			},
		},
	}
}
