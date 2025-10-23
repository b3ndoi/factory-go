# Factory-Go

[![CI](https://github.com/b3ndoi/factory-go/actions/workflows/ci.yml/badge.svg)](https://github.com/b3ndoi/factory-go/actions/workflows/ci.yml)
[![CodeQL](https://github.com/b3ndoi/factory-go/actions/workflows/codeql.yml/badge.svg)](https://github.com/b3ndoi/factory-go/actions/workflows/codeql.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/b3ndoi/factory-go.svg)](https://pkg.go.dev/github.com/b3ndoi/factory-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/b3ndoi/factory-go)](https://goreportcard.com/report/github.com/b3ndoi/factory-go)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Tests](https://img.shields.io/badge/tests-60%20passing-brightgreen)](factory/factory_test.go)
[![Coverage](https://img.shields.io/badge/coverage-89%25-brightgreen)](factory/factory_test.go)

A Laravel-inspired factory pattern for Go, providing a flexible and type-safe way to build test data with defaults, traits, and optional persistence.

## Features

- 🎯 **Type-safe** - Uses Go generics for full type safety
- 🔧 **Flexible** - Support for defaults, traits, and custom persistence
- 🚀 **Laravel-inspired** - Familiar API if you've used Laravel factories
- 🧪 **Test-friendly** - Perfect for seeding test databases or creating in-memory fixtures
- 🔄 **Faker integration** - Easy integration with faker libraries for realistic data
- 🐛 **Debugging tools** - Tap() for inspecting items during creation
- 🌍 **Environment-aware** - When()/Unless() for conditional behavior
- 🔁 **Factory variations** - Clone() for creating factory variations
- ⚡ **Must* variants** - Panic on error for cleaner test code
- 📦 **JSON support** - Direct JSON output for API testing
- 🔗 **Relationships** - Built-in support for model relationships

## Installation

```bash
go get github.com/b3ndoi/factory-go
```

## 📖 Examples

Check out the `/examples` directory for comprehensive examples:

- **[basic/](examples/basic)** - Getting started with Factory-Go ⭐
- **[api_testing/](examples/api_testing)** - HTTP API testing with RawJSON
- **[database_seeding/](examples/database_seeding)** - Seeding DBs with relationships
- **[complete_app/](examples/complete_app)** - Full-featured blog application
- **[faker_integration/](examples/faker_integration)** - Realistic data with faker

Each example is runnable: `cd examples/basic && go run main.go`

## Quick Reference

```go
// Setup
factory := factory.New(makeFn).
    WithDefaults(trait).           // Default values
    WithRawDefaults(trait).        // Only for Raw/JSON
    DefineState("admin", trait).   // Named states
    Sequence(trait1, trait2).      // Cycle patterns
    WithPersist(persistFn).        // DB persistence
    BeforeCreate(hookFn).          // Before hooks
    AfterCreate(hookFn).           // After hooks
    Tap(inspectFn).                // Debug/log
    When(condition, trait).        // Conditional
    Unless(condition, trait)       // Inverse conditional

// Create single
user := factory.Make()                    // In-memory
user := factory.Raw()                     // With rawDefaults
json := factory.MustRawJSON()             // As JSON (panic on error)
user := factory.MustCreate(ctx)           // Persist (panic on error)

// Create multiple  
users := factory.MakeMany(10)             // In-memory
users := factory.Count(10).Make()         // Fluent API
users := factory.Count(5).State("admin").MustCreate(ctx)

// Relationships
post := factory.For(postFactory, userFactory, linkFn).Make()                     // Each child gets own parent
posts := factory.Recycle(postFactory, user, linkFn).Count(5).Make()              // All share same parent
user, posts := factory.Has(userFactory, postFactory, 5, linkFn).MustCreate(ctx)  // Parent with children
user, roles, pivots := factory.HasAttached(userF, roleF, pivotF, 3, linkFn).MustCreate(ctx)  // Many-to-many

// Utilities
factory.Clone()          // Deep copy with reset sequence
factory.ResetSequence()  // Reset to seq=1
```

## Quick Start

```go
import "github.com/b3ndoi/factory-go/factory"

// Define your model
type User struct {
    ID    string
    Name  string
    Email string
    Role  string
}

// Create a factory
userFactory := factory.New(func(seq int64) User {
    return User{
        Name:  fmt.Sprintf("User %d", seq),
        Email: fmt.Sprintf("user%d@example.com", seq),
        Role:  "user",
    }
})

// Make an in-memory user (not persisted)
user := userFactory.Make()

// Make 10 users at once
users := userFactory.MakeMany(10)
// Or use the fluent Count() API
users = userFactory.Count(10).Make()

// Make with custom traits
admin := userFactory.Make(func(u *User) {
    u.Role = "admin"
})
```

## Sequence - Cycling Through Attributes

The `Sequence` method allows you to cycle through different attribute values when creating multiple models, just like [Laravel's sequence()](https://laravel.com/docs/12.x/eloquent-factories#sequences):

```go
// Alternate between admin and user roles
userFactory := factory.New(func(seq int64) User {
    return User{
        Name:  fmt.Sprintf("User %d", seq),
        Email: fmt.Sprintf("user%d@example.com", seq),
    }
}).Sequence(
    func(u *User) { u.Role = "admin" },
    func(u *User) { u.Role = "user" },
)

// Creates: admin, user, admin, user, admin
users := userFactory.MakeMany(5)
```

### Advanced Sequences

Sequences work with any number of states and automatically cycle:

```go
// Three-state sequence
statusFactory := factory.New(func(seq int64) Order {
    return Order{Number: seq}
}).Sequence(
    func(o *Order) { o.Status = "pending" },
    func(o *Order) { o.Status = "processing" },
    func(o *Order) { o.Status = "completed" },
)

// Creates 10 orders cycling through: pending, processing, completed, pending...
orders := statusFactory.MakeMany(10)
```

### Sequence with Per-Call Overrides

Per-call traits always override sequence values:

```go
factory := factory.New(makeFn).Sequence(
    func(u *User) { u.Role = "admin" },
    func(u *User) { u.Role = "user" },
)

u1 := factory.Make()                                      // Role: "admin" (from sequence)
u2 := factory.Make(func(u *User) { u.Role = "guest" })   // Role: "guest" (override)
u3 := factory.Make()                                      // Role: "admin" (sequence continues)
```

## Named States

Named states let you define reusable state configurations, similar to [Laravel's state methods](https://laravel.com/docs/12.x/eloquent-factories#factory-states):

```go
// Define named states
userFactory := factory.New(func(seq int64) User {
    return User{
        Name:  fmt.Sprintf("User %d", seq),
        Email: fmt.Sprintf("user%d@example.com", seq),
        Role:  "user",
    }
}).DefineState("admin", func(u *User) {
    u.Role = "admin"
    u.Permissions = []string{"read", "write", "delete"}
}).DefineState("moderator", func(u *User) {
    u.Role = "moderator"
    u.Permissions = []string{"read", "write"}
}).DefineState("verified", func(u *User) {
    u.EmailVerifiedAt = time.Now()
})

// Use named states - much cleaner than inline functions!
admin := userFactory.State("admin").Make()
verifiedAdmin := userFactory.State("admin").State("verified").Make()

// Works with all factory methods
admins, _ := userFactory.State("admin").CreateMany(ctx, 5)
```

### Benefits of Named States

1. **Reusable** - Define once, use everywhere
2. **Readable** - `State("admin")` is clearer than inline functions
3. **Chainable** - Combine multiple states easily
4. **Type-safe** - Panics if you reference an undefined state

```go
// Chain multiple states
user := factory.State("admin").State("verified").State("premium").Make()

// Override state with per-call traits
customAdmin := factory.State("admin").Make(func(u *User) {
    u.Name = "Custom Admin Name"
})
```

## Using WithDefaults (Faker Integration)

The `WithDefaults` method is perfect for integrating faker libraries or defining reusable default values:

```go
import (
    "github.com/b3ndoi/factory-go/factory"
    "github.com/brianvoe/gofakeit/v6"
)

// Create factory with faker defaults
userFactory := factory.New(func(seq int64) User {
    return User{} // Empty struct
}).WithDefaults(func(u *User) {
    // Use faker library for realistic data
    u.Name = gofakeit.Name()
    u.Email = gofakeit.Email()
    u.Role = "user"
})

// Each call generates unique fake data
user1 := userFactory.Make() // John Doe, john@example.com
user2 := userFactory.Make() // Jane Smith, jane@example.com

// Override specific fields
admin := userFactory.Make(func(u *User) {
    u.Role = "admin" // Keeps fake name and email
})
```

## Trait Application Order

Traits are applied in a specific order, allowing for flexible overrides:

### For Make() and Create():
1. **Base struct** (from `makeFn`)
2. **Defaults** (from `WithDefaults`) - for faker/default values
3. **Global traits** (from `WithTraits`) - for common modifications
4. **Sequence** (from `Sequence`) - cycles through on each item
5. **Per-call traits** (passed to `Make`/`Create`) - for specific customizations

### For Raw() and RawJSON():
1. **Base struct** (from `makeFn`)
2. **Defaults** (from `WithDefaults`) - for faker/default values
3. **RawDefaults** (from `WithRawDefaults`) - **only for Raw/RawJSON**
4. **Global traits** (from `WithTraits`) - for common modifications
5. **Sequence** (from `Sequence`) - cycles through on each item
6. **Per-call traits** (passed to `Raw`/`RawJSON`) - for specific customizations

```go
userFactory := factory.New(func(seq int64) User {
    return User{Role: "guest"} // 1. Base
}).WithDefaults(func(u *User) {
    u.Role = "user" // 2. Overrides base
    u.Name = gofakeit.Name()
}).WithTraits(func(u *User) {
    u.Email = strings.ToLower(u.Email) // 3. Modifies defaults
}).Sequence(
    func(u *User) { u.Role = "admin" },  // 4a. First item
    func(u *User) { u.Role = "moderator" }, // 4b. Second item (cycles)
)

// Per-call trait overrides everything
superuser := userFactory.Make(func(u *User) {
    u.Role = "superuser" // 5. Overrides all previous (including sequence)
})
```

## Raw Attributes & JSON (API Testing)

Get fully built objects without persisting - perfect for testing APIs:

### Raw() with Separate Defaults

Use `WithRawDefaults()` to add fields only for raw/JSON output (not persistence):

```go
type User struct {
    ID       string
    Name     string
    Email    string
    Password string `json:"password,omitempty"` // Only for API, not DB
}

userFactory := factory.New(func(seq int64) User {
    return User{
        Name:  fmt.Sprintf("User %d", seq),
        Email: fmt.Sprintf("user%d@example.com", seq),
    }
}).WithRawDefaults(func(u *User) {
    // This ONLY applies to Raw/RawJSON, not Make/Create
    u.Password = "test-password-123"
})

// Make() does NOT include rawDefaults
user := userFactory.Make()  // Password: ""

// Raw() DOES include rawDefaults
rawUser := userFactory.Raw() // Password: "test-password-123"
```

### RawJSON() for API Testing

Get JSON directly for HTTP tests:

```go
// Single object as JSON
jsonData, err := userFactory.RawJSON()
// POST to API: {"name":"User 1","email":"user1@example.com","password":"test-password-123"}

// Multiple objects as JSON array
jsonArray, err := userFactory.RawManyJSON(10)
// Bulk API request

// Works with Count() fluent API
jsonData, err := userFactory.Count(5).RawJSON()

// Works with states
jsonData, err := userFactory.State("admin").RawJSON()
```

### Real-World Example

```go
// Testing user registration endpoint
func TestUserRegistration(t *testing.T) {
    userFactory := factory.New(func(seq int64) UserRequest {
        return UserRequest{
            Username: fmt.Sprintf("user%d", seq),
            Email:    fmt.Sprintf("user%d@test.com", seq),
        }
    }).WithRawDefaults(func(r *UserRequest) {
        r.Password = "ValidPassword123!"
        r.PasswordConfirm = "ValidPassword123!"
    })
    
    // Get JSON payload for API test
    payload, _ := userFactory.RawJSON()
    
    // POST to registration endpoint
    resp := testClient.POST("/api/register", payload)
    assert.Equal(t, 201, resp.StatusCode)
}
```

## Persistence with Create

```go
// Set up persistence
userFactory := factory.New(func(seq int64) User {
    return User{
        Name:  fmt.Sprintf("User %d", seq),
        Email: fmt.Sprintf("user%d@example.com", seq),
    }
}).WithPersist(func(ctx context.Context, u *User) (*User, error) {
    // Your database logic
    u.ID = uuid.New().String()
    err := db.InsertUser(ctx, u)
    return u, err
})

// Create and persist a single user
user, err := userFactory.Create(context.Background())

// Create and persist multiple users
users, err := userFactory.CreateMany(context.Background(), 10)

// Create with custom traits
admin, err := userFactory.Create(context.Background(), func(u *User) {
    u.Role = "admin"
})
```

## Reset Sequence (Test Isolation)

Reset the sequence counter to get predictable data in tests:

```go
func TestUserCreation(t *testing.T) {
    // Reset before each test for predictable sequence numbers
    userFactory.ResetSequence()
    
    user := userFactory.Make()
    // Always "User 1" because sequence was reset
    assert.Equal(t, "User 1", user.Name)
}

// Chainable
users := userFactory.ResetSequence().MakeMany(5)
// Creates: User 1, User 2, User 3, User 4, User 5
```

## Count() - Fluent API

Use the fluent `Count()` API for a more Laravel-like syntax:

```go
// Instead of MakeMany(10)
users := userFactory.Count(10).Make()

// Works with Create too
users, err := userFactory.Count(5).Create(ctx)

// Fully chainable with states
admins := userFactory.Count(3).State("admin").Make()
verifiedUsers, _ := userFactory.Count(10).State("verified").Create(ctx)

// With per-call traits
customUsers := userFactory.Count(5).Make(func(u *User) {
    u.Active = true
})

// Raw() also works
rawData := userFactory.Count(20).Raw()

// Times() is an alias for Count()
users = userFactory.Times(3).Make()
```

The `CountedFactory` returned by `Count()` has these methods:
- `Make(...traits) []T` - Build count items in-memory
- `Create(ctx, ...traits) ([]*T, error)` - Build and persist count items
- `Raw(...traits) []T` - Alias for Make
- `State(name) *CountedFactory[T]` - Apply a named state (chainable)

## Relationships

Factory-Go provides powerful relationship helpers for all common database relationship patterns.

### For() - Belongs To (Each Item Gets Its Own Parent)

The `For()` function creates a new related model for each item:

```go
type Post struct {
    ID       string
    Title    string
    AuthorID string
}

// Each post gets its own newly created user
post := factory.For(postFactory, userFactory, func(p *Post, u *User) {
    p.AuthorID = u.ID
}).Make()

// Create multiple posts, each with their own user
posts := factory.For(postFactory, userFactory, func(p *Post, u *User) {
    p.AuthorID = u.ID
}).MakeMany(3)
// Creates 3 posts with 3 different users
```

### ForModel() / Recycle() - Belongs To (Shared Parent)

Use an existing model instance across multiple children:

```go
// Create/get an existing user
user, _ := userFactory.Create(ctx)

// All posts will belong to the same user
posts := factory.ForModel(postFactory, user, func(p *Post, u *User) {
    p.AuthorID = u.ID
}).MakeMany(10)

// Recycle() is an alias - more semantic name
posts := factory.Recycle(postFactory, user, func(p *Post, u *User) {
    p.AuthorID = u.ID
}).Count(10).MustCreate(ctx)
```

### Has() - One-to-Many (Parent with Children)

Create a parent with multiple children (inverse of `For`):

```go
// Create 1 user with 5 posts
user, posts := factory.Has(userFactory, postFactory, 5, func(u *User, p *Post) {
    p.AuthorID = u.ID
}).Make()
// Returns: user + []posts

// With persistence
user, posts, err := factory.Has(userFactory, postFactory, 3, func(u *User, p *Post) {
    p.AuthorID = u.ID
}).Create(ctx)
// Creates and saves 1 user + 3 posts

// MustCreate variant
user, posts := factory.Has(userFactory, postFactory, 10, linkFn).MustCreate(ctx)
```

### HasAttached() - Many-to-Many with Pivot

Create parent with many-to-many relationships through a pivot table:

```go
type UserRole struct {
    UserID string
    RoleID string
    Active bool // Pivot field
}

// Create user with 3 roles and pivot records
user, roles, err := factory.HasAttached(
    userFactory,
    roleFactory,
    userRoleFactory, // Pivot factory
    3,
    func(pivot *UserRole, user *User, role *Role) {
        pivot.UserID = user.ID
        pivot.RoleID = role.ID
        pivot.Active = true
    },
).Create(ctx)
// Creates: 1 user + 3 roles + 3 pivot records
```

### Relationship Pattern Summary

| Pattern | Function | Use Case | Example |
|---------|----------|----------|---------|
| Belongs To (unique) | `For()` | Each child has different parent | Post → User (each post by different user) |
| Belongs To (shared) | `ForModel()` / `Recycle()` | All children share same parent | Posts → User (all by same user) |
| Has Many | `Has()` | Parent with multiple children | User → Posts (user with multiple posts) |
| Many-to-Many | `HasAttached()` | Parent with many children + pivot | User → Roles (with pivot attributes) |

## Must* Variants (Clean Test Code)

Panic on error instead of returning - perfect for tests where you want to fail fast!

```go
func TestUserCreation(t *testing.T) {
    // No error handling needed - panics on failure
    user := userFactory.MustCreate(ctx)
    assert.Equal(t, "user@example.com", user.Email)
    
    // Works with Count() too
    users := userFactory.Count(10).MustCreate(ctx)
    assert.Len(t, users, 10)
    
    // JSON variants
    jsonData := userFactory.MustRawJSON()
    jsonArray := userFactory.Count(5).MustRawJSON()
}
```

**Available Must* methods:**
- `MustCreate(ctx, ...traits)` - Create and panic on error
- `MustCreateMany(ctx, count, ...traits)` - Create many and panic on error
- `MustRawJSON(...traits)` - Get JSON and panic on marshal error
- `MustRawManyJSON(count, ...traits)` - Get JSON array and panic on marshal error
- `Count(n).MustCreate(ctx)` - Fluent API with Must

## Tap() - Debugging & Inspection

Inspect or log items during creation without modifying them:

```go
// Debug what's being created
userFactory := factory.New(makeFn).Tap(func(u User) {
    fmt.Printf("Creating: %+v\n", u)
})

// Count items
count := 0
factory.Tap(func(u User) { count++ }).MakeMany(10)
fmt.Printf("Created %d users\n", count)

// Log to file
factory.Tap(func(u User) {
    log.Printf("User created: %s (%s)", u.Name, u.Email)
}).CreateMany(ctx, 5)

// Validate during creation
factory.Tap(func(u User) {
    if u.Email == "" {
        panic("Email is required!")
    }
}).Make()
```

**Key points:**
- Non-intrusive - doesn't modify the item
- Called for every Make(), Raw(), Create() operation
- Useful for debugging, logging, validation, counting

## When() / Unless() - Conditional Logic

Apply traits based on runtime conditions:

```go
isProd := os.Getenv("ENV") == "production"
isTest := !isProd

userFactory := factory.New(makeFn).
    When(isProd, func(u *User) {
        u.Email = faker.Email() // Real emails in production
    }).
    Unless(isProd, func(u *User) {
        u.Email = "test@example.com" // Fixed email in test
    }).
    When(isTest, func(u *User) {
        u.Active = false // Inactive users in tests
    })

user := userFactory.Make() // Traits applied based on environment
```

### Real-World Examples

```go
// Database-specific behavior
usePostgres := config.DB == "postgres"

factory.
    When(usePostgres, func(u *User) {
        u.CreatedAt = time.Now()
    }).
    Unless(usePostgres, func(u *User) {
        u.CreatedAt = time.Time{} // Let SQLite handle it
    })

// Feature flags
enableNewFeature := featureFlags.IsEnabled("new_feature")

factory.
    When(enableNewFeature, func(u *User) {
        u.NewField = "enabled"
    }).
    Unless(enableNewFeature, func(u *User) {
        u.NewField = ""
    })
```

## Clone() - Factory Variations

Create factory variations without affecting the original:

```go
// Base factory
baseFactory := factory.New(func(seq int64) User {
    return User{
        Name:  fmt.Sprintf("User %d", seq),
        Email: fmt.Sprintf("user%d@example.com", seq),
        Role:  "user",
    }
}).WithDefaults(func(u *User) {
    u.Active = true
})

// Create variations
adminFactory := baseFactory.Clone().WithTraits(func(u *User) {
    u.Role = "admin"
})

moderatorFactory := baseFactory.Clone().WithTraits(func(u *User) {
    u.Role = "moderator"
})

testFactory := baseFactory.Clone().WithTraits(func(u *User) {
    u.Email = "test@example.com"
})

// Each factory is independent
regularUser := baseFactory.Make()     // Role: "user"
admin := adminFactory.Make()          // Role: "admin"
moderator := moderatorFactory.Make()  // Role: "moderator"
testUser := testFactory.Make()        // Email: "test@example.com"
```

**Key features:**
- Deep copy of all traits, states, and hooks
- Sequence counter is reset for each clone
- Original factory remains unchanged
- Perfect for creating test variations

## Lifecycle Hooks

### Before Create Hooks

Run logic before persistence (e.g., validation, setup):

```go
userFactory := factory.New(func(seq int64) User {
    return User{Name: fmt.Sprintf("User %d", seq)}
}).WithPersist(persistFn).
  BeforeCreate(func(ctx context.Context, u *User) error {
    // Validate before saving
    if u.Email == "" {
        return errors.New("email is required")
    }
    return nil
}).BeforeCreate(func(ctx context.Context, u *User) error {
    // Set computed fields
    u.Slug = slugify(u.Name)
    return nil
})
```

If any `BeforeCreate` hook returns an error, persistence is skipped and the error is returned.

### After Create Hooks

Run logic after persistence (e.g., creating related records):

```go
userFactory := factory.New(func(seq int64) User {
    return User{Name: fmt.Sprintf("User %d", seq)}
}).WithPersist(persistFn).
  AfterCreate(func(ctx context.Context, u *User) error {
    // Create a profile for this user
    return db.CreateProfile(ctx, u.ID)
}).AfterCreate(func(ctx context.Context, u *User) error {
    // Send welcome email
    return emailService.SendWelcome(ctx, u.Email)
})
```

### Hook Execution Order

When calling `Create()`:
1. `Make()` - Build object with traits
2. **BeforeCreate hooks** - Run in order (can return error)
3. **Persist** - Save to database
4. **AfterCreate hooks** - Run in order (can return error)

## Complete Example with Faker and Named States

```go
package main

import (
    "context"
    "fmt"
    "time"
    
    "github.com/b3ndoi/factory-go/factory"
    "github.com/brianvoe/gofakeit/v6"
)

type User struct {
    ID              string
    FirstName       string
    LastName        string
    Email           string
    Role            string
    Active          bool
    EmailVerifiedAt *time.Time
}

func main() {
    userFactory := factory.New(func(seq int64) User {
        return User{}
    }).WithDefaults(func(u *User) {
        // Generate realistic fake data
        u.FirstName = gofakeit.FirstName()
        u.LastName = gofakeit.LastName()
        u.Email = gofakeit.Email()
        u.Role = "user"
        u.Active = true
    }).DefineState("admin", func(u *User) {
        u.Role = "admin"
    }).DefineState("moderator", func(u *User) {
        u.Role = "moderator"
    }).DefineState("verified", func(u *User) {
        now := time.Now()
        u.EmailVerifiedAt = &now
    }).DefineState("inactive", func(u *User) {
        u.Active = false
    }).WithPersist(func(ctx context.Context, u *User) (*User, error) {
        u.ID = gofakeit.UUID()
        // db.Insert(ctx, u)
        return u, nil
    })

    ctx := context.Background()

    // Create 10 regular verified users
    users, _ := userFactory.State("verified").CreateMany(ctx, 10)
    
    // Create 5 verified admins (chain multiple states!)
    admins, _ := userFactory.State("admin").State("verified").CreateMany(ctx, 5)
    
    // Create 3 inactive moderators
    mods, _ := userFactory.State("moderator").State("inactive").CreateMany(ctx, 3)
    
    // Create custom user with state + per-call override
    special, _ := userFactory.State("admin").Create(ctx, func(u *User) {
        u.FirstName = "Special"
        u.LastName = "Admin"
    })
    
    fmt.Printf("Created %d users, %d admins, %d moderators, 1 special admin\n", 
        len(users), len(admins), len(mods))
    fmt.Printf("Special admin: %s %s (%s)\n", 
        special.FirstName, special.LastName, special.Role)
}
```

## API Reference

### Factory Methods

#### Setup Methods
- `New(makeFn)` - Create a new factory with a base make function
- `WithDefaults(...traits)` - Set default traits (applied first, ideal for faker)
- `WithRawDefaults(...traits)` - Set traits applied ONLY for Raw/RawJSON methods
- `WithTraits(...traits)` - Add global traits (applied after defaults)
- `Sequence(...traits)` - Set traits that cycle through for each item created
- `DefineState(name, trait)` - Register a named state for reusable configurations
- `WithPersist(persistFn)` - Set persistence function (required for Create methods)
- `BeforeCreate(hookFn)` - Add hooks that run before persistence
- `AfterCreate(hookFn)` - Add hooks that run after persistence
- `Tap(fn func(T))` - Set function to inspect/log each created item
- `When(condition, ...traits)` - Apply traits only if condition is true
- `Unless(condition, ...traits)` - Apply traits only if condition is false
- `Clone()` - Create deep copy of factory with reset sequence

#### State Application
- `State(name)` - Apply a named state (returns new factory instance with state applied)

#### Fluent API
- `Count(n)` - Set count for fluent API (returns `CountedFactory`)
- `Times(n)` - Alias for `Count()`

#### Creation Methods (Single Item)
- `Make(...traits)` - Build object in-memory without persisting
- `Raw(...traits)` - Build with rawDefaults applied (for API testing)
- `RawJSON(...traits)` - Build and marshal to JSON
- `Create(ctx, ...traits)` - Build, persist, and run hooks for one object

#### Creation Methods (Multiple Items)
- `MakeMany(count, ...traits)` - Build multiple objects in-memory
- `RawMany(count, ...traits)` - Build multiple with rawDefaults applied
- `RawManyJSON(count, ...traits)` - Build multiple and marshal to JSON array
- `CreateMany(ctx, count, ...traits)` - Build, persist, and run hooks for multiple objects

#### Must* Variants (Panic on Error)
- `MustCreate(ctx, ...traits)` - Create and panic on error
- `MustCreateMany(ctx, count, ...traits)` - Create many and panic on error
- `MustRawJSON(...traits)` - Get JSON and panic on marshal error
- `MustRawManyJSON(count, ...traits)` - Get JSON array and panic on marshal error

#### Relationship Helpers
- `For[T, R](factory, relatedFactory, linkFn)` - Belongs-to: Each child gets its own parent
- `ForModel[T, R](factory, relatedModel, linkFn)` - Belongs-to: All children share same parent
- `Recycle[T, R](factory, relatedModel, linkFn)` - Alias for ForModel (semantic naming)
- `Has[T, R](parentFactory, childFactory, count, linkFn)` - Has-many: Create parent with children
- `HasAttached[T, R, P](parent, related, pivot, count, linkFn)` - Many-to-many: Parent with children + pivot

#### Utility Methods
- `ResetSequence()` - Reset sequence counter to 0 (useful for test isolation)
- `Clone()` - Create deep copy of factory (resets sequence)

### CountedFactory Methods

Returned by `Count()` or `Times()`:

- `Make(...traits) []T` - Build count items in-memory
- `Create(ctx, ...traits) ([]*T, error)` - Build, persist, and run hooks for count items
- `MustCreate(ctx, ...traits) []*T` - Create count items and panic on error
- `Raw(...traits) []T` - Build count items with rawDefaults applied
- `RawJSON(...traits) ([]byte, error)` - Build count items and marshal to JSON array
- `MustRawJSON(...traits) []byte` - Get JSON array and panic on marshal error
- `State(name) *CountedFactory[T]` - Apply a named state (chainable)

### HasFactory Methods

Returned by `Has()`:

- `Make() (T, []R)` - Build parent with children in-memory
- `Create(ctx) (*T, []*R, error)` - Create and persist parent with children
- `MustCreate(ctx) (*T, []*R)` - Create parent with children, panic on error

### HasAttachedFactory Methods

Returned by `HasAttached()`:

- `Make() (T, []R, []P)` - Build parent with related models and pivots in-memory
- `Create(ctx) (*T, []*R, []*P, error)` - Create and persist parent, related models, and pivot records
- `MustCreate(ctx) (*T, []*R, []*P)` - Create parent, related, and pivots, panic on error

### Type Definitions

```go
type Trait[T any] func(*T)
type BeforeCreate[T any] func(ctx context.Context, t *T) error
type AfterCreate[T any] func(ctx context.Context, t *T) error
type PersistFn[T any] func(ctx context.Context, t *T) (*T, error)
```

## Comparison with Laravel

| Laravel | Factory-Go |
|---------|------------|
| `User::factory()->make()` | `userFactory.Make()` |
| `User::factory()->raw()` | `userFactory.Raw()` |
| `User::factory()->count(10)->make()` | `userFactory.Count(10).Make()` or `MakeMany(10)` |
| `User::factory()->create()` | `userFactory.Create(ctx)` |
| `User::factory()->count(10)->create()` | `userFactory.Count(10).Create(ctx)` or `CreateMany(ctx, 10)` |
| `User::factory()->admin()->create()` | `userFactory.State("admin").Create(ctx)` |
| `User::factory()->count(5)->admin()->create()` | `userFactory.Count(5).State("admin").Create(ctx)` |
| `User::factory()->sequence(...)->create()` | `userFactory.Sequence(...).Create(ctx)` |
| `Post::factory()->for(User::factory())->create()` | `factory.For(postFactory, userFactory, linkFn).Create(ctx)` |
| `Post::factory()->for($user)->create()` | `factory.ForModel(postFactory, user, linkFn).Create(ctx)` |
| `Post::factory()->recycle($user)->create()` | `factory.Recycle(postFactory, user, linkFn).Create(ctx)` |
| `User::factory()->has(Post::factory()->count(3))` | `factory.Has(userFactory, postFactory, 3, linkFn).Create(ctx)` |
| `User::factory()->hasAttached(Role::factory(), pivot)` | `factory.HasAttached(userFactory, roleFactory, pivotFactory, 3, linkFn)` |
| `definition()` | `WithDefaults()` |
| `configure()` | `WithTraits()` |
| `public function admin() { return $this->state(...); }` | `DefineState("admin", trait)` |
| `beforeCreating()` | `BeforeCreate()` |
| `afterCreating()` | `AfterCreate()` |
| *(No equivalent)* | `ResetSequence()`, `Clone()`, `Tap()`, `When()`, `Unless()`, `Must*` |

## License

MIT License - See LICENSE file for details

