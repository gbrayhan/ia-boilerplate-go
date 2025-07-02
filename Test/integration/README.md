# Integration Tests

This directory contains the integration tests for the ia-boilerplate-go API. The tests are written in Cucumber/Gherkin and executed with the Godog framework.

## Testing Framework Philosophy

### 1. **Automatic Resource Management**
- Every resource created during the tests is automatically tracked.
- Automatic cleanup is executed at the end of each scenario.
- This prevents contamination between tests.

### 2. **Automatic Authentication**
- Each scenario performs a login at startup.
- Access tokens are managed globally.
- Authorization headers are added automatically.

### 3. **Dynamic Variables**
- Unique values are generated to avoid conflicts.
- Variables can be substituted in URLs and payloads.
- Values persist across steps in a scenario.

### 4. **Robust Validation**
- HTTP status codes are verified.
- JSON response structure is validated.
- Edge cases and error handling are covered.

## File Structure

```
Test/integration/
├── main_test.go              # Main test configuration
├── steps.go                  # Gherkin step implementations
├── README.md                 # This file
└── features/                 # Gherkin feature files
    ├── auth.feature          # Authentication tests
    ├── users.feature         # Users, roles and devices
    ├── medicine.feature      # Medicines
    ├── icd-cie.feature       # ICD-CIE codes
    ├── device-info.feature   # Device information
    └── error-handling.feature # Error handling tests
```

## Feature Files

### 1. **auth.feature**
Authentication and authorization tests:
- Login with valid/invalid credentials
- Token refresh
- Access to protected endpoints without authentication

### 2. **users.feature**
Comprehensive user management tests:
- CRUD for user roles
- CRUD for users
- CRUD for devices linked to users
- Paginated and property based searches

### 3. **medicine.feature**
Medicine management tests:
- CRUD for medicines
- Unique EAN code validation
- Advanced searches
- Required field handling

### 4. **icd-cie.feature**
ICD-CIE code tests:
- CRUD for ICD-CIE records
- Searches with multiple filters
- Search property validation
- Pagination and edge cases

### 5. **device-info.feature**
Device information tests:
- Device information endpoint
- Authenticated health check
- Device middleware verification

### 6. **error-handling.feature**
Error handling tests:
- Failed authentication cases
- Invalid IDs
- Missing required fields
- Malformed JSON payloads
- Pagination edge cases

## Running the Tests

### Option 1: Automated Script (Recommended)

```bash
# Run all tests
./scripts/run-integration-test.bash

# Run specific feature files
./scripts/run-integration-test.bash -f auth.feature
./scripts/run-integration-test.bash -f users.feature

# Run with Docker
./scripts/run-integration-test.bash -d -v

# Run with specific tags
./scripts/run-integration-test.bash -t @smoke

# Verbose mode
./scripts/run-integration-test.bash -v
```

### Option 2: Direct Command

```bash
# Run all tests
go test -tags=integration ./Test/integration/...

# Verbose output
go test -v -tags=integration ./Test/integration/...

# Run a specific feature
INTEGRATION_FEATURE_FILE=auth.feature go test -tags=integration ./Test/integration/...

# Run with scenario tags
INTEGRATION_SCENARIO_TAGS=@smoke go test -tags=integration ./Test/integration/...
```

### Option 3: Docker Compose

```bash
# Run tests with Docker
docker compose run --rm ia-boilerplate go test -tags=integration ./Test/integration/...

# Verbose output
docker compose run --rm ia-boilerplate go test -v -tags=integration ./Test/integration/...
```

## Environment Variables

Before running the tests make sure all variables from your `.env` file are loaded.
Only a few optional variables are shown below:

| Variable | Description | Example |
|----------|-------------|---------|
| `INTEGRATION_FEATURE_FILE` | Run only one feature file | `auth.feature` |
| `INTEGRATION_SCENARIO_TAGS` | Run scenarios with specific tags | `@smoke` |
| `INTEGRATION_TEST` | Enable testing mode | `true` |

## Example Scenario Structure

```gherkin
Scenario: TC01 - Create a new user successfully
  Given I generate a unique alias as "newUserUsername"
  And I generate a unique alias as "newUserEmail"
  When I send a POST request to "/api/users" with body:
    """
    {
      "username": "${newUserUsername}",
      "email": "${newUserEmail}@test.com",
      "password": "securePassword123",
      "roleId": 1,
      "enabled": true
    }
    """
  Then the response code should be 201
  And the JSON response should contain key "id"
  And I save the JSON response key "id" as "userID"
```

## Available Steps

### Given Steps (Setup)
- `I generate a unique alias as "varName"`
- `I generate a unique EAN code as "varName"`
- `I clear the authentication token`
- `I am authenticated as a user`

### When Steps (Actions)
- `I send a GET request to "path"`
- `I send a POST request to "path" with body:`
- `I send a PUT request to "path" with body:`
- `I send a DELETE request to "path"`

### Then Steps (Assertions)
- `the response code should be 200`
- `the JSON response should contain key "keyName"`
- `the JSON response should contain "field": "value"`
- `the JSON response should contain error "error": "message"`
- `I save the JSON response key "key" as "varName"`

## Resource Management

### Automatic Creation
Resources created during tests are automatically tracked:

```go
// In steps.go
func trackResource(path string) {
    // Track resources for later cleanup
}
```

### Automatic Cleanup
At the end of each scenario all created resources are removed:

```go
// In steps.go
func InitializeScenario(ctx *godog.ScenarioContext) {
    // Automatic setup and teardown
}
```

## Debugging

### Verbose Mode
```bash
go test -v -tags=integration ./Test/integration/...
```

### Detailed Logs
The tests include detailed logs showing:
- Request URLs
- Sent headers
- Response codes
- Response bodies
- Generated variables

### Debug Variables
```bash
# Enable debug logs
export DEBUG=true
go test -tags=integration ./Test/integration/...
```

## Best Practices

### 1. **Unique Names**
Always generate unique values:
```gherkin
Given I generate a unique alias as "testUser"
```

### 2. **Complete Validation**
Validate both the response code and its content:
```gherkin
Then the response code should be 201
And the JSON response should contain key "id"
And the JSON response should contain "username": "${testUser}"
```

### 3. **Error Handling**
Include tests for error cases:
```gherkin
Scenario: Attempt to create user with missing fields
  When I send a POST request to "/api/users" with body:
    """
    {
      "firstName": "John"
    }
    """
  Then the response code should be 400
  And the JSON response should contain key "error"
```

### 4. **Resource Cleanup**
Resources are cleaned automatically, but you can remove them manually:
```gherkin
When I send a DELETE request to "/api/users/${userID}"
Then the response code should be 200
```

## Troubleshooting

### Common Issues

1. **Database connection errors**
   - Ensure Docker Compose is running.
   - Check the database environment variables.

2. **Tests failing due to existing resources**
   - Run with the `-c` option to clean up before starting.
   - Make sure no tests are running in parallel.

3. **Authentication errors**
   - Verify the test credentials.
   - Check that the server is running.

4. **Test timeouts**
   - Increase the timeout in the configuration.
   - Verify network connectivity.

### Debug Logs
```bash
# Enable detailed logs
export GODOG_DEBUG=true
go test -v -tags=integration ./Test/integration/...
```

## Contributing

### Adding New Tests

1. **Create a feature file**
   ```bash
   touch Test/integration/features/new-feature.feature
   ```

2. **Implement steps** (if needed):
   - Add functions in `steps.go`.
   - Register them in `InitializeScenario`.

3. **Run the tests**:
   ```bash
   ./scripts/run-integration-test.bash -f new-feature.feature
   ```

### Naming Conventions

- **Feature files**: `kebab-case.feature`
- **Scenarios**: `TC01 - Test description`
- **Variables**: `camelCase` or `snake_case`
- **Tags**: `@smoke`, `@regression`, `@critical`

## Continuous Integration

### GitHub Actions
```yaml
- name: Run Integration Tests
  run: |
    docker compose up -d
    ./scripts/run-integration-test.bash -d -v
```

### Jenkins Pipeline
```groovy
stage('Integration Tests') {
    steps {
        sh './scripts/run-integration-test.bash -d -v'
    }
}
```

## Additional Resources

- [Godog Documentation](https://github.com/cucumber/godog)
- [Gherkin Syntax](https://cucumber.io/docs/gherkin/)
- [Testing in Go](https://golang.org/pkg/testing/)
