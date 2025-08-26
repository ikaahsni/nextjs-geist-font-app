package jira

import "time"

// TestCase represents a test case in Jira
type TestCase struct {
	ID          string            `json:"id,omitempty"`
	Key         string            `json:"key,omitempty"`
	Summary     string            `json:"summary" binding:"required"`
	Description string            `json:"description"`
	Status      string            `json:"status,omitempty"`
	Priority    string            `json:"priority,omitempty"`
	Labels      []string          `json:"labels,omitempty"`
	Components  []string          `json:"components,omitempty"`
	TestType    string            `json:"testType,omitempty"` // Manual, Automated, etc.
	CreatedDate time.Time         `json:"createdDate,omitempty"`
	UpdatedDate time.Time         `json:"updatedDate,omitempty"`
	Reporter    string            `json:"reporter,omitempty"`
	Assignee    string            `json:"assignee,omitempty"`
	CustomFields map[string]interface{} `json:"customFields,omitempty"`
}

// TestExecution represents a test execution in Jira
type TestExecution struct {
	ID              string                 `json:"id,omitempty"`
	Key             string                 `json:"key,omitempty"`
	Summary         string                 `json:"summary" binding:"required"`
	Description     string                 `json:"description"`
	Status          string                 `json:"status,omitempty"`
	TestCases       []string               `json:"testCases" binding:"required"` // Array of test case keys
	ExecutionStatus string                 `json:"executionStatus,omitempty"`    // PASS, FAIL, TODO, EXECUTING
	StartDate       time.Time              `json:"startDate,omitempty"`
	EndDate         time.Time              `json:"endDate,omitempty"`
	ExecutedBy      string                 `json:"executedBy,omitempty"`
	Environment     string                 `json:"environment,omitempty"`
	TestResults     []TestResult           `json:"testResults,omitempty"`
	CustomFields    map[string]interface{} `json:"customFields,omitempty"`
}

// TestResult represents the result of a single test case execution
type TestResult struct {
	TestCaseKey     string    `json:"testCaseKey"`
	Status          string    `json:"status"` // PASS, FAIL, TODO, EXECUTING
	Comment         string    `json:"comment,omitempty"`
	ExecutionTime   int       `json:"executionTime,omitempty"` // in milliseconds
	ExecutedBy      string    `json:"executedBy,omitempty"`
	ExecutedOn      time.Time `json:"executedOn,omitempty"`
	Defects         []string  `json:"defects,omitempty"` // Array of defect keys
	Evidence        []string  `json:"evidence,omitempty"` // Array of attachment URLs
}

// TestPlan represents a test plan in Jira
type TestPlan struct {
	ID           string            `json:"id,omitempty"`
	Key          string            `json:"key,omitempty"`
	Summary      string            `json:"summary" binding:"required"`
	Description  string            `json:"description"`
	Status       string            `json:"status,omitempty"`
	TestCases    []string          `json:"testCases,omitempty"` // Array of test case keys
	CreatedDate  time.Time         `json:"createdDate,omitempty"`
	UpdatedDate  time.Time         `json:"updatedDate,omitempty"`
	Owner        string            `json:"owner,omitempty"`
	CustomFields map[string]interface{} `json:"customFields,omitempty"`
}

// JiraIssue represents a generic Jira issue structure
type JiraIssue struct {
	ID     string     `json:"id,omitempty"`
	Key    string     `json:"key,omitempty"`
	Fields IssueFields `json:"fields"`
}

// IssueFields represents the fields of a Jira issue
type IssueFields struct {
	Summary     string      `json:"summary"`
	Description string      `json:"description,omitempty"`
	IssueType   IssueType   `json:"issuetype"`
	Project     Project     `json:"project"`
	Priority    Priority    `json:"priority,omitempty"`
	Status      Status      `json:"status,omitempty"`
	Reporter    User        `json:"reporter,omitempty"`
	Assignee    User        `json:"assignee,omitempty"`
	Labels      []string    `json:"labels,omitempty"`
	Components  []Component `json:"components,omitempty"`
}

// IssueType represents a Jira issue type
type IssueType struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name"`
}

// Project represents a Jira project
type Project struct {
	Key  string `json:"key"`
	Name string `json:"name,omitempty"`
}

// Priority represents a Jira priority
type Priority struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name"`
}

// Status represents a Jira status
type Status struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name"`
}

// User represents a Jira user
type User struct {
	AccountID    string `json:"accountId,omitempty"`
	EmailAddress string `json:"emailAddress,omitempty"`
	DisplayName  string `json:"displayName,omitempty"`
}

// Component represents a Jira component
type Component struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name"`
}

// JiraResponse represents a generic Jira API response
type JiraResponse struct {
	Issues     []JiraIssue `json:"issues,omitempty"`
	StartAt    int         `json:"startAt,omitempty"`
	MaxResults int         `json:"maxResults,omitempty"`
	Total      int         `json:"total,omitempty"`
}

// ErrorResponse represents an error response from Jira API
type ErrorResponse struct {
	ErrorMessages []string          `json:"errorMessages,omitempty"`
	Errors        map[string]string `json:"errors,omitempty"`
}

// CreateIssueRequest represents a request to create a Jira issue
type CreateIssueRequest struct {
	Fields IssueFields `json:"fields"`
}

// CreateIssueResponse represents a response from creating a Jira issue
type CreateIssueResponse struct {
	ID   string `json:"id"`
	Key  string `json:"key"`
	Self string `json:"self"`
}
