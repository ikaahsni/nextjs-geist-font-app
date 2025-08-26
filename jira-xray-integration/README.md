# Jira Xray-like Integration API

A comprehensive Go application that provides Xray-like test management functionality with Jira integration. This application allows you to manage test cases, test executions, and test results through a RESTful API that integrates with Atlassian Jira.

## Features

- 🧪 **Test Case Management**: Create, read, and manage test cases
- 🚀 **Test Execution Tracking**: Create and track test executions
- 📊 **Test Results**: Record and manage test results
- 🔗 **Jira Integration**: Seamless integration with Jira REST API
- 🎯 **RESTful API**: Clean and intuitive REST endpoints
- 🔒 **Authentication**: Secure Jira API authentication
- 📝 **Comprehensive Logging**: Detailed logging for debugging
- 🎭 **Demo Mode**: Works with mock data for testing

## Prerequisites

- Go 1.21 or higher
- Access to a Jira instance (Cloud or Server)
- Jira API token (for authentication)

## Installation

1. **Clone or download the project**:
   ```bash
   cd jira-xray-integration
   ```

2. **Install dependencies**:
   ```bash
   go mod tidy
   ```

3. **Set up environment configuration**:
   ```bash
   cp .env.sample .env
   ```

4. **Edit the `.env` file with your actual Jira credentials**:
   ```env
   JIRA_BASE_URL=https://yourcompany.atlassian.net
   JIRA_USERNAME=your-email@company.com
   JIRA_API_TOKEN=your-api-token-here
   JIRA_PROJECT_KEY=YOUR_PROJECT_KEY
   PORT=8080
   ```

## Getting Jira API Token

1. Go to [Atlassian Account Settings](https://id.atlassian.com/manage-profile/security/api-tokens)
2. Click "Create API token"
3. Give it a label (e.g., "Xray Integration")
4. Copy the generated token and use it in your `.env` file

## Running the Application

1. **Start the server**:
   ```bash
   go run main.go
   ```

2. **The server will start on port 8080** (or the port specified in your `.env` file)

3. **Access the API documentation**:
   ```
   http://localhost:8080/api/info
   ```

## API Endpoints

### Health Check
```bash
GET /api/health
```

### Test Cases

#### List all test cases
```bash
curl -X GET http://localhost:8080/api/testcases
```

#### Create a new test case
```bash
curl -X POST http://localhost:8080/api/testcases \
  -H "Content-Type: application/json" \
  -d '{
    "summary": "Login functionality test",
    "description": "Test user login with valid credentials",
    "priority": "High",
    "labels": ["login", "authentication"],
    "testType": "Manual"
  }'
```

#### Get a specific test case
```bash
curl -X GET http://localhost:8080/api/testcases/TEST-1
```

### Test Executions

#### List all test executions
```bash
curl -X GET http://localhost:8080/api/testexecutions
```

#### Create a new test execution
```bash
curl -X POST http://localhost:8080/api/testexecutions \
  -H "Content-Type: application/json" \
  -d '{
    "summary": "Sprint 1 Test Execution",
    "description": "Execute all test cases for Sprint 1",
    "testCases": ["TEST-1", "TEST-2"],
    "environment": "QA"
  }'
```

#### Get a specific test execution
```bash
curl -X GET http://localhost:8080/api/testexecutions/EXEC-1
```

## API Response Examples

### Test Case Response
```json
{
  "testCase": {
    "id": "10001",
    "key": "TEST-1",
    "summary": "Login functionality test",
    "description": "Test user login with valid credentials",
    "status": "To Do",
    "priority": "High",
    "labels": ["login", "authentication"],
    "testType": "Manual",
    "createdDate": "2024-01-15T10:30:00Z",
    "reporter": "John Doe"
  },
  "message": "Test case created successfully"
}
```

### Test Execution Response
```json
{
  "testExecution": {
    "id": "10005",
    "key": "EXEC-1",
    "summary": "Sprint 1 Test Execution",
    "description": "Execute all test cases for Sprint 1",
    "status": "In Progress",
    "testCases": ["TEST-1", "TEST-2"],
    "executionStatus": "EXECUTING",
    "environment": "QA",
    "executedBy": "Jane Smith",
    "startDate": "2024-01-15T14:00:00Z"
  },
  "message": "Test execution created successfully"
}
```

## Demo Mode

If you're using the demo credentials (demo_user/demo_token_replace_with_actual), the application will:

- ⚠️ Display warnings about using demo credentials
- 📊 Return mock data for all API calls
- 🎭 Simulate Jira API responses
- ✅ Allow you to test the API without a real Jira connection

## Configuration

### Environment Variables

| Variable | Description | Required | Default |
|----------|-------------|----------|---------|
| `JIRA_BASE_URL` | Your Jira instance URL | Yes | - |
| `JIRA_USERNAME` | Your Jira email address | Yes | - |
| `JIRA_API_TOKEN` | Your Jira API token | Yes | - |
| `JIRA_PROJECT_KEY` | Jira project key for tests | Yes | - |
| `PORT` | Server port | No | 8080 |

### Jira Issue Types

This application assumes your Jira instance has the following issue types:
- **Test**: For test cases
- **Test Execution**: For test executions

If these don't exist, you may need to:
1. Install Xray for Jira, or
2. Create custom issue types, or
3. Modify the code to use existing issue types

## Project Structure

```
jira-xray-integration/
├── main.go              # Main application entry point
├── config.go            # Configuration management
├── go.mod              # Go module dependencies
├── .env.sample         # Sample environment configuration
├── README.md           # This file
└── jira/
    ├── models.go       # Jira data models
    └── client.go       # Jira API client
```

## Development

### Adding New Endpoints

1. Define new models in `jira/models.go`
2. Add client methods in `jira/client.go`
3. Create handlers in `main.go`
4. Add routes to the router

### Testing with curl

Test all endpoints with the provided curl commands. For JSON responses, you can pipe through `jq` for better formatting:

```bash
curl -X GET http://localhost:8080/api/testcases | jq '.'
```

## Error Handling

The application provides comprehensive error handling:

- **400 Bad Request**: Invalid JSON or missing required fields
- **500 Internal Server Error**: Jira API errors or server issues
- **Detailed error messages**: All errors include descriptive messages

## Logging

The application logs:
- 📝 All HTTP requests and responses
- 🔍 Jira API calls and responses
- ⚠️ Configuration warnings
- ❌ Error details for debugging

## Security Considerations

- 🔒 **Never commit your `.env` file** to version control
- 🔑 **Use API tokens instead of passwords** for Jira authentication
- 🛡️ **Rotate API tokens regularly**
- 🌐 **Use HTTPS in production**
- 🔐 **Consider implementing rate limiting** for production use

## Troubleshooting

### Common Issues

1. **"Failed to load configuration" error**:
   - Ensure your `.env` file exists and contains all required variables

2. **Jira API authentication errors**:
   - Verify your API token is correct
   - Check that your username (email) is correct
   - Ensure your Jira URL is correct (include https://)

3. **"Test" or "Test Execution" issue type not found**:
   - Install Xray for Jira, or
   - Modify the code to use existing issue types

4. **Connection timeout errors**:
   - Check your network connection
   - Verify the Jira URL is accessible

### Debug Mode

To enable more detailed logging, you can modify the Gin mode:

```go
gin.SetMode(gin.DebugMode) // Add this in main.go
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

This project is open source and available under the MIT License.

## Support

For issues and questions:
1. Check the troubleshooting section
2. Review the logs for error details
3. Ensure your Jira configuration is correct
4. Test with demo mode first

---

**Happy Testing! 🧪✨**
