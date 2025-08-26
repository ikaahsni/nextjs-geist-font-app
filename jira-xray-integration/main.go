package main

import (
	"log"
	"net/http"

	"jira-xray-integration/jira"

	"github.com/gin-gonic/gin"
)

var (
	config     *Config
	jiraClient *jira.Client
)

func main() {
	// Load configuration
	var err error
	config, err = LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Validate and log configuration
	config.ValidateConfig()

	// Initialize Jira client
	jiraClient = jira.NewClient(
		config.JiraBaseURL,
		config.JiraUsername,
		config.JiraAPIToken,
		config.JiraProjectKey,
	)

	// Initialize Gin router
	router := gin.Default()

	// Add middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(corsMiddleware())

	// API routes
	api := router.Group("/api")
	{
		// Test Case routes
		api.GET("/testcases", getTestCases)
		api.POST("/testcases", createTestCase)
		api.GET("/testcases/:key", getTestCase)

		// Test Execution routes
		api.GET("/testexecutions", getTestExecutions)
		api.POST("/testexecutions", createTestExecution)
		api.GET("/testexecutions/:key", getTestExecution)

		// Health check
		api.GET("/health", healthCheck)

		// API info
		api.GET("/info", getAPIInfo)
	}

	// Root route
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message":     "Jira Xray-like Integration API",
			"version":     "1.0.0",
			"description": "A Go application for test management with Jira integration",
			"endpoints": gin.H{
				"health":          "/api/health",
				"info":            "/api/info",
				"testcases":       "/api/testcases",
				"testexecutions":  "/api/testexecutions",
			},
		})
	})

	// Start server
	port := ":" + config.Port
	log.Printf("🚀 Server starting on port %s", config.Port)
	log.Printf("📋 API Documentation available at: http://localhost%s/api/info", port)
	
	if err := router.Run(port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// corsMiddleware adds CORS headers
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// Health check endpoint
func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "healthy",
		"timestamp": gin.H{
			"unix": gin.H{
				"seconds": gin.H{
					"value": "current_time",
				},
			},
		},
		"jira": gin.H{
			"base_url":    config.JiraBaseURL,
			"project_key": config.JiraProjectKey,
			"demo_mode":   config.JiraUsername == "demo_user",
		},
	})
}

// API info endpoint
func getAPIInfo(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"name":        "Jira Xray-like Integration API",
		"version":     "1.0.0",
		"description": "A Go application for test management with Jira integration",
		"endpoints": gin.H{
			"GET /api/health":                    "Health check",
			"GET /api/info":                      "API information",
			"GET /api/testcases":                 "List all test cases",
			"POST /api/testcases":                "Create a new test case",
			"GET /api/testcases/:key":            "Get a specific test case",
			"GET /api/testexecutions":            "List all test executions",
			"POST /api/testexecutions":           "Create a new test execution",
			"GET /api/testexecutions/:key":       "Get a specific test execution",
		},
		"example_requests": gin.H{
			"create_test_case": gin.H{
				"method": "POST",
				"url":    "/api/testcases",
				"body": gin.H{
					"summary":     "Sample Test Case",
					"description": "This is a sample test case description",
					"priority":    "High",
					"labels":      []string{"api", "integration"},
					"testType":    "Manual",
				},
			},
			"create_test_execution": gin.H{
				"method": "POST",
				"url":    "/api/testexecutions",
				"body": gin.H{
					"summary":     "Sample Test Execution",
					"description": "This is a sample test execution",
					"testCases":   []string{"TEST-1", "TEST-2"},
					"environment": "QA",
				},
			},
		},
		"configuration": gin.H{
			"jira_base_url":  config.JiraBaseURL,
			"project_key":    config.JiraProjectKey,
			"demo_mode":      config.JiraUsername == "demo_user",
		},
	})
}

// Get all test cases
func getTestCases(c *gin.Context) {
	log.Println("Handling GET /api/testcases request")

	testCases, err := jiraClient.ListTestCases()
	if err != nil {
		log.Printf("Error fetching test cases: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch test cases",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"testCases": testCases,
		"count":     len(testCases),
		"message":   "Test cases retrieved successfully",
	})
}

// Create a new test case
func createTestCase(c *gin.Context) {
	log.Println("Handling POST /api/testcases request")

	var testCase jira.TestCase
	if err := c.ShouldBindJSON(&testCase); err != nil {
		log.Printf("Error binding JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	// Validate required fields
	if testCase.Summary == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Summary is required",
		})
		return
	}

	createdTestCase, err := jiraClient.CreateTestCase(&testCase)
	if err != nil {
		log.Printf("Error creating test case: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create test case",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"testCase": createdTestCase,
		"message":  "Test case created successfully",
	})
}

// Get a specific test case
func getTestCase(c *gin.Context) {
	key := c.Param("key")
	log.Printf("Handling GET /api/testcases/%s request", key)

	// For demo purposes, return a mock test case
	// In a real implementation, you would fetch from Jira
	c.JSON(http.StatusOK, gin.H{
		"testCase": gin.H{
			"id":          "10001",
			"key":         key,
			"summary":     "Sample Test Case",
			"description": "This is a sample test case",
			"status":      "To Do",
			"priority":    "Medium",
		},
		"message": "Test case retrieved successfully",
	})
}

// Get all test executions
func getTestExecutions(c *gin.Context) {
	log.Println("Handling GET /api/testexecutions request")

	// For demo purposes, return mock test executions
	// In a real implementation, you would fetch from Jira
	mockExecutions := []gin.H{
		{
			"id":              "10005",
			"key":             "EXEC-1",
			"summary":         "Sprint 1 Test Execution",
			"status":          "In Progress",
			"executionStatus": "EXECUTING",
			"testCases":       []string{"TEST-1", "TEST-2"},
			"environment":     "QA",
		},
		{
			"id":              "10006",
			"key":             "EXEC-2",
			"summary":         "Regression Test Execution",
			"status":          "Done",
			"executionStatus": "PASS",
			"testCases":       []string{"TEST-1", "TEST-2", "TEST-3"},
			"environment":     "Staging",
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"testExecutions": mockExecutions,
		"count":          len(mockExecutions),
		"message":        "Test executions retrieved successfully",
	})
}

// Create a new test execution
func createTestExecution(c *gin.Context) {
	log.Println("Handling POST /api/testexecutions request")

	var testExecution jira.TestExecution
	if err := c.ShouldBindJSON(&testExecution); err != nil {
		log.Printf("Error binding JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	// Validate required fields
	if testExecution.Summary == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Summary is required",
		})
		return
	}

	if len(testExecution.TestCases) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "At least one test case is required",
		})
		return
	}

	createdTestExecution, err := jiraClient.CreateTestExecution(&testExecution)
	if err != nil {
		log.Printf("Error creating test execution: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create test execution",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"testExecution": createdTestExecution,
		"message":       "Test execution created successfully",
	})
}

// Get a specific test execution
func getTestExecution(c *gin.Context) {
	key := c.Param("key")
	log.Printf("Handling GET /api/testexecutions/%s request", key)

	testExecution, err := jiraClient.GetTestExecution(key)
	if err != nil {
		log.Printf("Error fetching test execution: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch test execution",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"testExecution": testExecution,
		"message":       "Test execution retrieved successfully",
	})
}
