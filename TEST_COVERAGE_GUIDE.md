# Test Coverage Guide

## Current Status

**Coverage: 1.6%** (Threshold: 75%)

The low coverage is because the mock-based tests don't exercise the actual implementation code in `mongo.go` files. The mocks return hardcoded values without calling the actual database methods.

## Why Mocks Don't Provide Coverage

When you test a mock implementation directly:
```go
mockStore := &MockUrlStore{
    GetFunc: func(ctx *context.Context, key, env string) (*UrlData, error) {
        return &UrlData{...}, nil
    },
}
got, err := mockStore.Get(...) // Tests the mock, not the implementation
```

The actual `UrlMongoStore.Get()` method in `mongo.go` is **never called**, so it has 0% coverage.

## Solution: Test Real Implementation with Mocked Mongo Collection

Instead of testing mocks directly, test the real implementation by:
1. Creating a `MockCollection` that implements `*mongo.Collection` interface methods
2. Passing the mock collection to the real store
3. Testing the real store's logic against mock data

### Example Refactor

**Before (mock testing):**
```go
mockStore := &MockUrlStore{
    GetFunc: func(...) { return &UrlData{...}, nil },
}
got, err := mockStore.Get(ctx, key, env) // 0% coverage of mongo.go
```

**After (implementation testing with mocked dependency):**
```go
mockColl := &MockMongoCollection{
    FindOneFunc: func(...) *mongo.SingleResult {
        return mockSingleResult(&UrlData{Key: "user-server", Url: "..."})
    },
}
store := NewUrlMongoStore(mockColl)
got, err := store.Get(ctx, "user-server", "LOCAL") // Tests real implementation!
```

## Recommended Approach for CI

1. **Unit Tests** (with mocks): Verify db layer implementation logic
2. **Integration Tests** (with real MongoDB in CI): Run in CI pipeline with Docker Mongo service
3. **Coverage**: Both combined should reach 75%+

## Files Affected

- `user-server/endpoints/db/mongo_test.go` - Tests `UrlStore` implementation
- `user-server/profile/db/mongo_test.go` - Tests `ProfileStore` implementation  
- `user-server/signup/db/mongo_test.go` - Tests `UserStore` implementation

## Next Steps

Choose one of the following approaches:

### Option 1: Keep Mocks + Add Integration Tests (Recommended)
Keep mock tests for fast unit testing, add integration tests that use real MongoDB container in CI to get coverage.

### Option 2: Refactor to Test Real Implementation
Refactor the test files to test the real store implementations with mocked MongoDB collections.

### Option 3: Adjust Coverage Threshold
Lower the coverage requirement for `_test.go` files since they're already testing the interface contracts.

## CI/CD Pipeline Status

✅ **Build Pipeline**: Working - `user-server-build`
✅ **Test Pipeline**: Working - `user-server-test` 
✅ **Quality Gate Pipeline**: Updated - `user-server-quality-gate`
✅ **Integrated Pipeline**: Updated - `user-server-integrated` (now checks all 3)

### Branch Protection Rules to Set in GitHub

Go to: **Settings > Branches > main > Branch protection rules**

Add required status checks:
- ✓ `user-server-build`
- ✓ `user-server-test`
- ✓ `user-server-quality-gate`
- ✓ `user-server-integrated`

Enable:
- ✓ Require branches to be up to date before merging
- ✓ Require status checks to pass before merging
- ✓ Dismiss stale pull request approvals
