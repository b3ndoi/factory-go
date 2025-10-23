# Factory-Go Examples

This directory contains comprehensive examples demonstrating all features of Factory-Go.

## Running the Examples

Each example is a standalone Go program. Run them with:

```bash
# Navigate to an example directory
cd examples/basic

# Run the example
go run main.go
```

---

## 📂 Available Examples

### 1. `basic/` - Getting Started ⭐ START HERE

**What it demonstrates:**
- ✅ Creating a simple factory
- ✅ Make() and MakeMany()
- ✅ Custom traits
- ✅ Count() fluent API
- ✅ Named states (DefineState, State)
- ✅ Sequences (cycling through values)
- ✅ ResetSequence()

**Perfect for:** First-time users learning the basics

```bash
cd examples/basic && go run main.go
```

---

### 2. `api_testing/` - HTTP API Testing

**What it demonstrates:**
- ✅ RawJSON() for API payloads
- ✅ WithRawDefaults() for API-specific fields
- ✅ Testing valid/invalid requests
- ✅ Bulk API testing
- ✅ Separation of domain vs API models
- ✅ Mock HTTP server testing

**Perfect for:** Testing REST APIs without a database

```bash
cd examples/api_testing && go run main.go
```

**Key Feature:** Shows how to add password fields for API testing without affecting domain models!

---

### 3. `database_seeding/` - Database Relationships

**What it demonstrates:**
- ✅ Has() - One-to-many relationships
- ✅ For() - Belongs-to with unique parents
- ✅ Recycle() - Belongs-to with shared parent
- ✅ HasAttached() - Many-to-many with pivot tables
- ✅ MustCreate() for clean error handling
- ✅ Tap() for logging
- ✅ Complex nested relationships

**Perfect for:** Seeding databases with related data

```bash
cd examples/database_seeding && go run main.go
```

**Key Feature:** Shows all 4 relationship patterns in action!

---

### 4. `complete_app/` - Full-Featured Blog Application

**What it demonstrates:**
- ✅ ALL factory features in one example
- ✅ Multiple models (User, Post, Comment, Tag)
- ✅ All relationship types
- ✅ BeforeCreate/AfterCreate hooks
- ✅ When/Unless conditionals
- ✅ Tap() debugging
- ✅ Clone() for factory variations
- ✅ Sequences for data variety
- ✅ Named states with chaining
- ✅ RawJSON for API testing
- ✅ WithRawDefaults separation

**Perfect for:** Understanding how everything works together

```bash
cd examples/complete_app && go run main.go
```

**Key Feature:** Comprehensive demonstration of ALL 35+ features!

---

### 5. `faker_integration/` - Realistic Data Generation

**What it demonstrates:**
- ✅ Integration with faker libraries
- ✅ WithDefaults() for faker data
- ✅ Realistic test data generation
- ✅ Overriding faker values
- ✅ Faker + Sequences
- ✅ Faker + States
- ✅ WithRawDefaults() for API fields
- ✅ Clone() for variations

**Perfect for:** Generating realistic, varied test data

```bash
cd examples/faker_integration && go run main.go
```

**Note:** Uses a simple faker simulation. In production, use:
```bash
go get github.com/brianvoe/gofakeit/v6
```

---

## 🎓 Learning Path

**Recommended order for learning:**

1. **Start with `basic/`** - Understand core concepts
2. **Then `api_testing/`** - Learn RawJSON and testing
3. **Then `database_seeding/`** - Master relationships
4. **Then `faker_integration/`** - Add realistic data
5. **Finally `complete_app/`** - See everything together

---

## 📚 What Each Example Teaches

### Basic Concepts
- `basic/` - Core factory operations
- `faker_integration/` - Realistic data generation

### Testing
- `api_testing/` - HTTP/REST API testing
- `database_seeding/` - Database testing with relationships

### Advanced
- `complete_app/` - Production-ready patterns with all features

---

## 🚀 Next Steps

After reviewing these examples:

1. **Try modifying them** - Change the models, add fields
2. **Combine patterns** - Mix features from different examples
3. **Build your own** - Create factories for your domain
4. **Read the main README** - For complete API reference

---

## 💡 Quick Tips from Examples

### From `basic/`:
- Use `Count()` for fluent API: `factory.Count(10).Make()`
- Define states once, use everywhere
- Sequences cycle automatically

### From `api_testing/`:
- `WithRawDefaults()` for passwords/tokens
- `MustRawJSON()` for test payloads
- Per-call traits override defaults

### From `database_seeding/`:
- `Has()` creates parent with children
- `Recycle()` shares same parent
- `HasAttached()` handles many-to-many

### From `complete_app/`:
- Combine features for powerful workflows
- `When/Unless` for environment-specific behavior
- `Clone()` for factory variations

### From `faker_integration/`:
- `WithDefaults()` perfect for faker
- Faker + sequences = varied data
- Override faker when needed

---

## 📖 Additional Resources

- **Main README:** `../README.md` - Complete API documentation
- **CHANGELOG:** `../CHANGELOG.md` - All features explained
- **Tests:** `../factory/factory_test.go` - 60 test examples
- **Comparison:** See how Factory-Go compares to alternatives

---

## 🎯 Common Use Cases

| Use Case | Example to Check |
|----------|------------------|
| Simple test data | `basic/` |
| API endpoint testing | `api_testing/` |
| Database seeding | `database_seeding/` |
| Realistic data | `faker_integration/` |
| Production setup | `complete_app/` |
| Relationships | `database_seeding/` + `complete_app/` |
| State management | `basic/` + `complete_app/` |
| Complex structures | `complete_app/` |

---

## ✨ Features Coverage

| Feature | basic | api | db | complete | faker |
|---------|-------|-----|-------|----------|-------|
| Make/Create | ✅ | ✅ | ✅ | ✅ | ✅ |
| Count() | ✅ | ✅ | ✅ | ✅ | ✅ |
| States | ✅ | - | ✅ | ✅ | ✅ |
| Sequences | ✅ | - | - | ✅ | ✅ |
| Relationships | - | - | ✅ | ✅ | - |
| RawJSON | - | ✅ | - | ✅ | ✅ |
| WithRawDefaults | - | ✅ | - | ✅ | ✅ |
| Hooks | - | - | - | ✅ | - |
| Tap() | - | - | ✅ | ✅ | - |
| When/Unless | - | - | - | ✅ | - |
| Clone() | - | - | - | ✅ | ✅ |
| Must* | - | ✅ | ✅ | ✅ | - |

---

Happy coding! 🚀

