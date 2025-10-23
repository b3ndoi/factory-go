# Changelog

## [Unreleased]

### Added

#### Raw() Method
- `Raw(...traits)` - Alias for `Make()` that emphasizes getting raw attributes
- `RawMany(count, ...traits)` - Alias for `MakeMany()`
- Useful for testing validations or API requests without persistence
- Works with all factory features (states, sequences, traits)

#### BeforeCreate Hooks
- `BeforeCreate(hookFn)` - Add hooks that run before persistence
- Similar to Laravel's `beforeCreating()` callback
- Useful for validation, computed fields, or setup logic
- Hooks run in order and can return errors to prevent persistence
- Multiple hooks can be chained
- Execution order: Make → BeforeCreate → Persist → AfterCreate

#### ResetSequence
- `ResetSequence()` - Reset the sequence counter to 0
- Essential for test isolation to get predictable sequence numbers
- Chainable for fluent API usage
- Thread-safe using atomic operations

#### Sequence Support
- `Sequence(...traits)` - Cycle through different attribute values when creating multiple models
- Similar to Laravel's `sequence()` method
- Automatically cycles through provided traits in order
- Works with `Make()`, `MakeMany()`, `Create()`, and `CreateMany()`
- Example: Alternate between admin/user roles when creating multiple users

#### Named States
- `DefineState(name, trait)` - Register reusable named states
- `State(name)` - Apply a named state to a factory
- Similar to Laravel's state methods (e.g., `User::factory()->admin()`)
- Chainable: `factory.State("admin").State("verified").Make()`
- Type-safe: Panics if referencing undefined states
- Makes code much more readable than inline trait functions

#### Enhanced Documentation
- Comprehensive README with examples
- Comparison table with Laravel factories
- Complete example showing faker integration with named states
- Trait application order documentation

### Features Summary

The factory now supports:
1. **Basic creation**: `Make()`, `Create()`, `MakeMany()`, `CreateMany()`
2. **Defaults**: `WithDefaults()` for faker/default values
3. **Global traits**: `WithTraits()` for common modifications
4. **Sequences**: `Sequence()` for cycling through patterns
5. **Named states**: `DefineState()` + `State()` for reusable configurations
6. **Persistence**: `WithPersist()` for database operations
7. **Hooks**: `AfterCreate()` for post-persistence operations

### Trait Application Order
1. Base struct (from `makeFn`)
2. Defaults (from `WithDefaults`)
3. Global traits (from `WithTraits`)
4. Sequence (from `Sequence`, cycles)
5. Named states (from `State()`)
6. Per-call traits (passed to `Make`/`Create`)

#### Count() - Fluent API
- `Count(n)` - Fluent API for creating multiple items (like Laravel's `count()`)
- Returns `CountedFactory` with chainable methods
- `Times(n)` - Semantic alias for `Count()`
- Fully chainable with states: `factory.Count(10).State("admin").Make()`
- Works with all creation methods: Make, Create, Raw
- More expressive than `MakeMany(10)` for Laravel users

#### For() - Relationship Helper
- `For[T, R](factory, relatedFactory, linkFn)` - Set up belongs-to relationships
- Creates a new related model for each item
- `ForModel[T, R](factory, model, linkFn)` - Use existing related model
- All items share the same related model
- Generic functions work with any two types
- Example: `factory.For(postFactory, userFactory, func(p *Post, u *User) { p.AuthorID = u.ID })`

#### RawJSON() - JSON Output for API Testing
- `RawJSON(...traits)` - Build object and marshal to JSON
- `RawManyJSON(count, ...traits)` - Build multiple objects and marshal to JSON array
- Perfect for testing API endpoints without database
- Works with all factory features (states, sequences, traits)
- CountedFactory also has `RawJSON()` method: `factory.Count(10).RawJSON()`

#### WithRawDefaults() - API-Specific Fields
- `WithRawDefaults(...traits)` - Set traits applied ONLY for Raw/RawJSON methods
- Useful for adding passwords, tokens, or computed fields for API testing
- Does NOT affect Make() or Create() methods
- Applied after WithDefaults but before WithTraits
- Example: Add password field for API requests but not for database persistence

#### Must* Variants - Cleaner Test Code
- `MustCreate(ctx, ...traits)` - Create and panic on error (no error handling)
- `MustCreateMany(ctx, count, ...traits)` - Create many and panic on error
- `MustRawJSON(...traits)` - Get JSON and panic on marshal error
- `MustRawManyJSON(count, ...traits)` - Get JSON array and panic on marshal error
- CountedFactory: `MustCreate(ctx)`, `MustRawJSON()`
- Perfect for tests where you want to fail fast
- Common Go idiom for test code

#### Tap() - Debugging & Inspection
- `Tap(fn func(T))` - Set function called for each created item
- Non-intrusive debugging without modifying items
- Useful for logging, counting, validation
- Works with Make(), Raw(), and Create() methods
- Example: `factory.Tap(func(u User) { fmt.Printf("%+v\n", u) })`

#### When() / Unless() - Conditional Logic
- `When(condition, ...traits)` - Apply traits only if condition is true
- `Unless(condition, ...traits)` - Apply traits only if condition is false
- Perfect for environment-specific behavior
- Example: `factory.When(isProd, func(u *User) { u.Email = faker.Email() })`
- Chainable for multiple conditions

#### Clone() - Factory Variations
- `Clone()` - Create deep copy of factory with reset sequence
- All traits, states, hooks are deep copied
- Sequence counter resets to 0 for each clone
- Original factory remains unchanged
- Perfect for creating test variations
- Example: `adminFactory := baseFactory.Clone().State("admin")`

#### Advanced Relationship Features
- `Has[T, R](parent, child, count, linkFn)` - Create parent with multiple children (one-to-many)
- Returns `HasFactory[T, R]` with `Make()`, `Create()`, `MustCreate()` methods
- Inverse of `For()` - creates one parent with many children
- Example: `factory.Has(userFactory, postFactory, 5, linkFn).Create(ctx)`
- Perfect for seeding data with related records

- `Recycle[T, R](factory, model, linkFn)` - Alias for `ForModel()` with semantic naming
- Reuse the same parent model across multiple children
- More readable than `ForModel()` in some contexts
- Example: `factory.Recycle(postFactory, user, linkFn).Count(10).Create(ctx)`

- `HasAttached[T, R, P](parent, related, pivot, count, linkFn)` - Many-to-many relationships
- Creates parent, related models, and pivot table records
- Supports pivot attributes (e.g., `active`, `created_at`)
- Returns `HasFactory[T, R]`
- Example: User with Roles through UserRole pivot table
- Most complex but most powerful relationship helper

### Test Coverage
- 60 comprehensive tests covering all features
- Tests for edge cases (overrides, chaining, panics, errors, relationships, JSON, conditionals, pivots)
- 100% passing test suite
- 89% code coverage
- Tier 1 tests: Raw, RawMany, ResetSequence, BeforeCreate (8 tests)
- Tier 2 tests: Count, Times, For, ForModel (10 tests)
- JSON/RawDefaults tests: RawJSON, RawManyJSON, WithRawDefaults (7 tests)
- Features 1-4 tests: Must* (5), Tap (2), When/Unless (3), Clone (3) = 13 tests
- Advanced Relationships: Has (3), Recycle (3), HasAttached (3) = 9 tests

