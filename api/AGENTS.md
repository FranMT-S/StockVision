# AGENTS RULES

## Code Style

- Every function must include a comment.
- The comment must start with the function’s name.
(Example: // MyFunction does something important)

- Every function must be named in camelCase.
(Example: myFunction, handleRequest)

- Every file name is in camelCase starting with lowercase.
(Example: agentManager.go, userService.go)

## Testing

- use testify to test golang api
- use testify suit to order tests
- every test must be have test cases 
- Test files must be finished with _test: `*_test.go`
  example: `customClient_test.go`
- Mock external dependencies appropriately