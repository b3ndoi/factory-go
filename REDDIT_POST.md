# Reddit Post for r/golang

## Title
```
Factory-Go: Type-safe factory pattern library for Go
```

## Post Body

```markdown
I've released Factory-Go v1.0.0, a type-safe factory pattern library for Go that uses generics.

**Purpose:** Announcing a new library and seeking feedback on API design and implementation.

**AI Usage:** This project was built with significant AI assistance (Claude via Cursor). I directed the architecture, features, and API design (inspired by Laravel's factory pattern), but approximately 95% of the code, tests, and documentation were AI-generated. This was primarily a learning project to understand Go generics and the factory pattern.

**The Problem:**
Existing Go factory libraries (like bluele/factory) use `interface{}` and require type assertions. Most Go projects end up writing hundreds of lines of manual test helper functions that aren't reusable.

**What Factory-Go Does:**
Provides type-safe test data generation using Go 1.21+ generics. Includes Laravel-inspired patterns: named states, sequences, and relationship helpers for belongs-to, has-many, and many-to-many patterns.

**Example:**
```go
userFactory := factory.New(func(seq int64) User {
    return User{Name: fmt.Sprintf("User %d", seq)}
})

// Type-safe, no assertions needed
users := userFactory.Count(10).Make()

// With relationships
user, posts := factory.Has(userFactory, postFactory, 5, linkFn).Create(ctx)
```

**Current Status:**
- 621 lines core library, 1,741 lines tests (AI-generated)
- 60 tests with 89% coverage
- Zero dependencies (stdlib only)
- 5 working examples
- Full CI/CD with GitHub Actions
- No production usage yet (just released)

**Links:**
- Repository: https://github.com/b3ndoi/factory-go
- Documentation: https://pkg.go.dev/github.com/b3ndoi/factory-go
- Examples: https://github.com/b3ndoi/factory-go/tree/main/examples

**Seeking Feedback:**
- Is the API intuitive for Go developers?
- Are there Laravel factory features I'm missing?
- Would you use this over manual test helpers?
- Given the heavy AI assistance, is this useful for the community or would you prefer human-written libraries?

Honest feedback appreciated, especially regarding whether AI-generated libraries provide value or if the community prefers entirely human-written code.
```

---

## Alternative: Small Projects Thread (Safer Option)

Given the AI disclosure and "no production usage yet," you might want to post in the **weekly "Small Projects" thread** instead of the front page.

The post above would work for either, but the Small Projects thread has:
- Lower stakes
- Less scrutiny  
- Still gets visibility
- More appropriate for AI-assisted projects

## When to Post

**Best times for r/golang:**
- Weekday mornings: 8-10am Eastern Time
- Avoid weekends
- Monday-Wednesday best

## What to Expect

**Possible reactions:**
- Questions about specific implementation details
- Requests for benchmarks vs other libraries
- Suggestions for missing features
- Some skepticism about AI-generated code
- Constructive feedback on API design

**Be ready to:**
- Respond honestly about AI usage
- Acknowledge limitations (no production use yet)
- Take feedback gracefully
- Explain design decisions you made
- Link to specific code examples

---

Would you like me to adjust the post for the Small Projects thread, or do you want to go for the front page with this version?


