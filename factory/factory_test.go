package factory

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
)

type User struct {
	ID, Name, Email string
}

func TestFactory_MakeAndCreate(t *testing.T) {
	f := New(func(seq int64) User {
		return User{
			Name:  fmt.Sprintf("User %d", seq),
			Email: fmt.Sprintf("user%d@example.com", seq),
		}
	})

	// Make only
	u := f.Make()
	if u.Name == "" {
		t.Fatal("expected name to be set")
	}

	// Create with a fake persist
	f = f.WithPersist(func(ctx context.Context, u *User) (*User, error) {
		u.ID = "uuid-123"
		return u, nil
	})

	ctx := context.Background()
	saved, err := f.Create(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if saved.ID == "" {
		t.Fatal("expected ID to be assigned")
	}
}

func TestFactory_WithDefaults(t *testing.T) {
	// Simulate a faker library or default value generator
	fakeName := "John Doe"
	fakeEmail := "john.doe@faker.com"

	f := New(func(seq int64) User {
		return User{} // Empty, will be filled by defaults
	}).WithDefaults(func(u *User) {
		// This could use a faker library like go-faker
		u.Name = fakeName
		u.Email = fakeEmail
	})

	u := f.Make()
	if u.Name != fakeName {
		t.Fatalf("expected name %q, got %q", fakeName, u.Name)
	}
	if u.Email != fakeEmail {
		t.Fatalf("expected email %q, got %q", fakeEmail, u.Email)
	}

	// Per-call traits override defaults
	customName := "Jane Smith"
	u2 := f.Make(func(u *User) {
		u.Name = customName
	})
	if u2.Name != customName {
		t.Fatalf("expected name %q, got %q", customName, u2.Name)
	}
	if u2.Email != fakeEmail {
		t.Fatalf("expected email to remain %q, got %q", fakeEmail, u2.Email)
	}
}

func TestFactory_ManyMethods(t *testing.T) {
	f := New(func(seq int64) User {
		return User{
			Name:  fmt.Sprintf("User %d", seq),
			Email: fmt.Sprintf("user%d@example.com", seq),
		}
	})

	// MakeMany
	users := f.MakeMany(3)
	if len(users) != 3 {
		t.Fatalf("expected 3 users, got %d", len(users))
	}
	if users[0].Name == users[1].Name {
		t.Fatal("expected unique names for each user")
	}

	// CreateMany
	f = f.WithPersist(func(ctx context.Context, u *User) (*User, error) {
		u.ID = fmt.Sprintf("id-%s", u.Name)
		return u, nil
	})

	ctx := context.Background()
	saved, err := f.CreateMany(ctx, 5)
	if err != nil {
		t.Fatal(err)
	}
	if len(saved) != 5 {
		t.Fatalf("expected 5 saved users, got %d", len(saved))
	}
	for _, u := range saved {
		if u.ID == "" {
			t.Fatal("expected all users to have IDs")
		}
	}
}

func TestFactory_Sequence(t *testing.T) {
	// Test simple alternating sequence
	f := New(func(seq int64) User {
		return User{
			Name:  fmt.Sprintf("User %d", seq),
			Email: fmt.Sprintf("user%d@example.com", seq),
		}
	}).Sequence(
		func(u *User) { u.ID = "admin" },
		func(u *User) { u.ID = "user" },
	)

	users := f.MakeMany(6)
	expected := []string{"admin", "user", "admin", "user", "admin", "user"}
	for i, u := range users {
		if u.ID != expected[i] {
			t.Fatalf("user %d: expected ID %q, got %q", i, expected[i], u.ID)
		}
	}
}

func TestFactory_SequenceWithThreeStates(t *testing.T) {
	// Test sequence with 3 different states
	f := New(func(seq int64) User {
		return User{Name: "Test"}
	}).Sequence(
		func(u *User) { u.ID = "pending" },
		func(u *User) { u.ID = "active" },
		func(u *User) { u.ID = "inactive" },
	)

	users := f.MakeMany(7)
	expected := []string{"pending", "active", "inactive", "pending", "active", "inactive", "pending"}
	for i, u := range users {
		if u.ID != expected[i] {
			t.Fatalf("user %d: expected ID %q, got %q", i, expected[i], u.ID)
		}
	}
}

func TestFactory_SequenceWithOverride(t *testing.T) {
	// Test that per-call traits override sequence
	f := New(func(seq int64) User {
		return User{Name: "Test"}
	}).Sequence(
		func(u *User) { u.ID = "admin" },
		func(u *User) { u.ID = "user" },
	)

	// First item should be admin (from sequence)
	u1 := f.Make()
	if u1.ID != "admin" {
		t.Fatalf("expected ID 'admin', got %q", u1.ID)
	}

	// Second item should be user (from sequence), but override it
	u2 := f.Make(func(u *User) { u.ID = "superuser" })
	if u2.ID != "superuser" {
		t.Fatalf("expected ID 'superuser', got %q", u2.ID)
	}

	// Third should be admin again (sequence continues)
	u3 := f.Make()
	if u3.ID != "admin" {
		t.Fatalf("expected ID 'admin', got %q", u3.ID)
	}
}

func TestFactory_SequenceWithCreateMany(t *testing.T) {
	// Test sequence works with CreateMany
	f := New(func(seq int64) User {
		return User{
			Name:  fmt.Sprintf("User %d", seq),
			Email: fmt.Sprintf("user%d@example.com", seq),
		}
	}).Sequence(
		func(u *User) { u.ID = "role-A" },
		func(u *User) { u.ID = "role-B" },
	).WithPersist(func(ctx context.Context, u *User) (*User, error) {
		// Persist keeps the ID from sequence
		return u, nil
	})

	ctx := context.Background()
	users, err := f.CreateMany(ctx, 5)
	if err != nil {
		t.Fatal(err)
	}

	expected := []string{"role-A", "role-B", "role-A", "role-B", "role-A"}
	for i, u := range users {
		if u.ID != expected[i] {
			t.Fatalf("user %d: expected ID %q, got %q", i, expected[i], u.ID)
		}
	}
}

func TestFactory_NamedStates(t *testing.T) {
	// Define a factory with named states
	f := New(func(seq int64) User {
		return User{
			Name:  fmt.Sprintf("User %d", seq),
			Email: fmt.Sprintf("user%d@example.com", seq),
			ID:    "user",
		}
	}).DefineState("admin", func(u *User) {
		u.ID = "admin"
	}).DefineState("moderator", func(u *User) {
		u.ID = "moderator"
	})

	// Test using State() to apply named state
	admin := f.State("admin").Make()
	if admin.ID != "admin" {
		t.Fatalf("expected ID 'admin', got %q", admin.ID)
	}

	moderator := f.State("moderator").Make()
	if moderator.ID != "moderator" {
		t.Fatalf("expected ID 'moderator', got %q", moderator.ID)
	}

	// Original factory should still produce regular users
	user := f.Make()
	if user.ID != "user" {
		t.Fatalf("expected ID 'user', got %q", user.ID)
	}
}

func TestFactory_NamedStatesWithCreate(t *testing.T) {
	// Test named states work with Create
	f := New(func(seq int64) User {
		return User{
			Name:  fmt.Sprintf("User %d", seq),
			Email: fmt.Sprintf("user%d@example.com", seq),
		}
	}).DefineState("verified", func(u *User) {
		u.ID = "verified"
	}).WithPersist(func(ctx context.Context, u *User) (*User, error) {
		return u, nil
	})

	ctx := context.Background()
	user, err := f.State("verified").Create(ctx)
	if err != nil {
		t.Fatal(err)
	}

	if user.ID != "verified" {
		t.Fatalf("expected ID 'verified', got %q", user.ID)
	}
}

func TestFactory_NamedStatesChaining(t *testing.T) {
	// Test that State() returns a factory that can be chained
	f := New(func(seq int64) User {
		return User{Name: "Test"}
	}).DefineState("admin", func(u *User) {
		u.ID = "admin"
	}).DefineState("active", func(u *User) {
		u.Email = "active@test.com"
	})

	// Chain multiple states
	user := f.State("admin").State("active").Make()
	if user.ID != "admin" {
		t.Fatalf("expected ID 'admin', got %q", user.ID)
	}
	if user.Email != "active@test.com" {
		t.Fatalf("expected Email 'active@test.com', got %q", user.Email)
	}
}

func TestFactory_NamedStatesWithMakeMany(t *testing.T) {
	// Test named states work with MakeMany
	f := New(func(seq int64) User {
		return User{
			Name:  fmt.Sprintf("User %d", seq),
			Email: fmt.Sprintf("user%d@example.com", seq),
		}
	}).DefineState("premium", func(u *User) {
		u.ID = "premium"
	})

	users := f.State("premium").MakeMany(3)
	if len(users) != 3 {
		t.Fatalf("expected 3 users, got %d", len(users))
	}

	for i, u := range users {
		if u.ID != "premium" {
			t.Fatalf("user %d: expected ID 'premium', got %q", i, u.ID)
		}
	}
}

func TestFactory_NamedStatesWithTraitOverride(t *testing.T) {
	// Test that per-call traits override named states
	f := New(func(seq int64) User {
		return User{Name: "Test"}
	}).DefineState("admin", func(u *User) {
		u.ID = "admin"
	})

	// Apply state, then override with trait
	user := f.State("admin").Make(func(u *User) {
		u.ID = "superadmin"
	})

	if user.ID != "superadmin" {
		t.Fatalf("expected ID 'superadmin', got %q", user.ID)
	}
}

func TestFactory_UnknownStatePanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for unknown state")
		}
	}()

	f := New(func(seq int64) User {
		return User{Name: "Test"}
	})

	// This should panic
	f.State("nonexistent").Make()
}

// Tier 1 Features Tests

func TestFactory_Raw(t *testing.T) {
	// Test that Raw() works like Make()
	f := New(func(seq int64) User {
		return User{
			Name:  fmt.Sprintf("User %d", seq),
			Email: fmt.Sprintf("user%d@example.com", seq),
		}
	}).DefineState("admin", func(u *User) {
		u.ID = "admin"
	})

	// Raw should work exactly like Make
	user := f.Raw()
	if user.Name == "" || user.Email == "" {
		t.Fatal("expected Raw to generate user with name and email")
	}

	// Raw with state
	admin := f.State("admin").Raw()
	if admin.ID != "admin" {
		t.Fatalf("expected ID 'admin', got %q", admin.ID)
	}

	// Raw with traits
	custom := f.Raw(func(u *User) {
		u.ID = "custom"
	})
	if custom.ID != "custom" {
		t.Fatalf("expected ID 'custom', got %q", custom.ID)
	}
}

func TestFactory_RawMany(t *testing.T) {
	f := New(func(seq int64) User {
		return User{
			Name:  fmt.Sprintf("User %d", seq),
			Email: fmt.Sprintf("user%d@example.com", seq),
		}
	})

	users := f.RawMany(5)
	if len(users) != 5 {
		t.Fatalf("expected 5 users, got %d", len(users))
	}

	// Verify sequence numbers are unique
	if users[0].Name == users[1].Name {
		t.Fatal("expected unique names for each user")
	}
}

func TestFactory_ResetSequence(t *testing.T) {
	f := New(func(seq int64) User {
		return User{
			Name:  fmt.Sprintf("User %d", seq),
			Email: fmt.Sprintf("user%d@example.com", seq),
		}
	})

	// Create first user
	u1 := f.Make()
	expectedName1 := "User 1"
	if u1.Name != expectedName1 {
		t.Fatalf("expected name %q, got %q", expectedName1, u1.Name)
	}

	// Create second user
	u2 := f.Make()
	expectedName2 := "User 2"
	if u2.Name != expectedName2 {
		t.Fatalf("expected name %q, got %q", expectedName2, u2.Name)
	}

	// Reset sequence
	f.ResetSequence()

	// Should start from 1 again
	u3 := f.Make()
	if u3.Name != expectedName1 {
		t.Fatalf("after reset, expected name %q, got %q", expectedName1, u3.Name)
	}
}

func TestFactory_ResetSequenceChaining(t *testing.T) {
	// Test that ResetSequence returns factory for chaining
	f := New(func(seq int64) User {
		return User{Name: fmt.Sprintf("User %d", seq)}
	})

	users := f.ResetSequence().MakeMany(3)
	if len(users) != 3 {
		t.Fatalf("expected 3 users, got %d", len(users))
	}

	if users[0].Name != "User 1" {
		t.Fatalf("expected 'User 1', got %q", users[0].Name)
	}
}

func TestFactory_BeforeCreate(t *testing.T) {
	beforeCalled := false
	beforeCallCount := 0

	f := New(func(seq int64) User {
		return User{
			Name:  fmt.Sprintf("User %d", seq),
			Email: fmt.Sprintf("user%d@example.com", seq),
		}
	}).BeforeCreate(func(ctx context.Context, u *User) error {
		beforeCalled = true
		beforeCallCount++
		// Modify user before persistence
		u.ID = "before-" + u.Name
		return nil
	}).WithPersist(func(ctx context.Context, u *User) (*User, error) {
		// Check that before hook already ran
		if u.ID == "" {
			t.Fatal("expected ID to be set by before hook")
		}
		return u, nil
	})

	ctx := context.Background()
	user, err := f.Create(ctx)
	if err != nil {
		t.Fatal(err)
	}

	if !beforeCalled {
		t.Fatal("expected before hook to be called")
	}

	if beforeCallCount != 1 {
		t.Fatalf("expected before hook to be called once, got %d", beforeCallCount)
	}

	if user.ID != "before-User 1" {
		t.Fatalf("expected ID 'before-User 1', got %q", user.ID)
	}
}

func TestFactory_BeforeCreateMultiple(t *testing.T) {
	callOrder := []string{}

	f := New(func(seq int64) User {
		return User{Name: "Test"}
	}).BeforeCreate(func(ctx context.Context, u *User) error {
		callOrder = append(callOrder, "before1")
		return nil
	}).BeforeCreate(func(ctx context.Context, u *User) error {
		callOrder = append(callOrder, "before2")
		return nil
	}).WithPersist(func(ctx context.Context, u *User) (*User, error) {
		callOrder = append(callOrder, "persist")
		return u, nil
	}).AfterCreate(func(ctx context.Context, u *User) error {
		callOrder = append(callOrder, "after")
		return nil
	})

	ctx := context.Background()
	_, err := f.Create(ctx)
	if err != nil {
		t.Fatal(err)
	}

	expected := []string{"before1", "before2", "persist", "after"}
	if len(callOrder) != len(expected) {
		t.Fatalf("expected %d calls, got %d", len(expected), len(callOrder))
	}

	for i, call := range expected {
		if callOrder[i] != call {
			t.Fatalf("call %d: expected %q, got %q", i, call, callOrder[i])
		}
	}
}

func TestFactory_BeforeCreateError(t *testing.T) {
	f := New(func(seq int64) User {
		return User{Name: "Test"}
	}).BeforeCreate(func(ctx context.Context, u *User) error {
		return fmt.Errorf("validation failed")
	}).WithPersist(func(ctx context.Context, u *User) (*User, error) {
		t.Fatal("persist should not be called when before hook fails")
		return u, nil
	})

	ctx := context.Background()
	_, err := f.Create(ctx)
	if err == nil {
		t.Fatal("expected error from before hook")
	}

	if err.Error() != "validation failed" {
		t.Fatalf("expected 'validation failed', got %q", err.Error())
	}
}

func TestFactory_BeforeCreateWithCreateMany(t *testing.T) {
	callCount := 0

	f := New(func(seq int64) User {
		return User{
			Name:  fmt.Sprintf("User %d", seq),
			Email: fmt.Sprintf("user%d@example.com", seq),
		}
	}).BeforeCreate(func(ctx context.Context, u *User) error {
		callCount++
		return nil
	}).WithPersist(func(ctx context.Context, u *User) (*User, error) {
		u.ID = "saved"
		return u, nil
	})

	ctx := context.Background()
	users, err := f.CreateMany(ctx, 3)
	if err != nil {
		t.Fatal(err)
	}

	if len(users) != 3 {
		t.Fatalf("expected 3 users, got %d", len(users))
	}

	if callCount != 3 {
		t.Fatalf("expected before hook to be called 3 times, got %d", callCount)
	}
}

// Tier 2 Features Tests

func TestFactory_Count(t *testing.T) {
	f := New(func(seq int64) User {
		return User{
			Name:  fmt.Sprintf("User %d", seq),
			Email: fmt.Sprintf("user%d@example.com", seq),
		}
	})

	// Count().Make() should return slice
	users := f.Count(5).Make()
	if len(users) != 5 {
		t.Fatalf("expected 5 users, got %d", len(users))
	}

	// Verify unique sequence numbers
	for i, u := range users {
		expectedName := fmt.Sprintf("User %d", i+1)
		if u.Name != expectedName {
			t.Fatalf("user %d: expected name %q, got %q", i, expectedName, u.Name)
		}
	}
}

func TestFactory_CountWithCreate(t *testing.T) {
	f := New(func(seq int64) User {
		return User{
			Name:  fmt.Sprintf("User %d", seq),
			Email: fmt.Sprintf("user%d@example.com", seq),
		}
	}).WithPersist(func(ctx context.Context, u *User) (*User, error) {
		u.ID = fmt.Sprintf("id-%s", u.Name)
		return u, nil
	})

	ctx := context.Background()
	users, err := f.Count(3).Create(ctx)
	if err != nil {
		t.Fatal(err)
	}

	if len(users) != 3 {
		t.Fatalf("expected 3 users, got %d", len(users))
	}

	for _, u := range users {
		if u.ID == "" {
			t.Fatal("expected all users to have IDs")
		}
	}
}

func TestFactory_CountWithState(t *testing.T) {
	f := New(func(seq int64) User {
		return User{
			Name:  fmt.Sprintf("User %d", seq),
			Email: fmt.Sprintf("user%d@example.com", seq),
		}
	}).DefineState("admin", func(u *User) {
		u.ID = "admin"
	})

	// Count with State should be chainable
	users := f.Count(4).State("admin").Make()
	if len(users) != 4 {
		t.Fatalf("expected 4 users, got %d", len(users))
	}

	for i, u := range users {
		if u.ID != "admin" {
			t.Fatalf("user %d: expected ID 'admin', got %q", i, u.ID)
		}
	}
}

func TestFactory_CountWithTraits(t *testing.T) {
	f := New(func(seq int64) User {
		return User{
			Name:  fmt.Sprintf("User %d", seq),
			Email: fmt.Sprintf("user%d@example.com", seq),
		}
	})

	users := f.Count(3).Make(func(u *User) {
		u.ID = "custom"
	})

	if len(users) != 3 {
		t.Fatalf("expected 3 users, got %d", len(users))
	}

	for i, u := range users {
		if u.ID != "custom" {
			t.Fatalf("user %d: expected ID 'custom', got %q", i, u.ID)
		}
	}
}

func TestFactory_Times(t *testing.T) {
	// Times is an alias for Count
	f := New(func(seq int64) User {
		return User{Name: fmt.Sprintf("User %d", seq)}
	})

	users := f.Times(3).Make()
	if len(users) != 3 {
		t.Fatalf("expected 3 users, got %d", len(users))
	}
}

func TestFactory_CountRaw(t *testing.T) {
	f := New(func(seq int64) User {
		return User{Name: fmt.Sprintf("User %d", seq)}
	})

	users := f.Count(2).Raw()
	if len(users) != 2 {
		t.Fatalf("expected 2 users, got %d", len(users))
	}
}

// For() relationship tests

type Post struct {
	ID       string
	Title    string
	AuthorID string
}

func TestFactory_For(t *testing.T) {
	userFactory := New(func(seq int64) User {
		return User{
			ID:    fmt.Sprintf("user-%d", seq),
			Name:  fmt.Sprintf("User %d", seq),
			Email: fmt.Sprintf("user%d@example.com", seq),
		}
	})

	postFactory := New(func(seq int64) Post {
		return Post{
			Title: fmt.Sprintf("Post %d", seq),
		}
	})

	// Use For to create a post with a user
	post := For(postFactory, userFactory, func(p *Post, u *User) {
		p.AuthorID = u.ID
	}).Make()

	if post.AuthorID == "" {
		t.Fatal("expected AuthorID to be set")
	}

	if post.AuthorID[:5] != "user-" {
		t.Fatalf("expected AuthorID to start with 'user-', got %q", post.AuthorID)
	}

	if post.Title == "" {
		t.Fatal("expected Title to be set")
	}
}

func TestFactory_ForModel(t *testing.T) {
	user := User{
		ID:    "existing-user",
		Name:  "Existing User",
		Email: "existing@example.com",
	}

	postFactory := New(func(seq int64) Post {
		return Post{
			Title: fmt.Sprintf("Post %d", seq),
		}
	})

	// Use ForModel with an existing user
	post := ForModel(postFactory, &user, func(p *Post, u *User) {
		p.AuthorID = u.ID
	}).Make()

	if post.AuthorID != "existing-user" {
		t.Fatalf("expected AuthorID 'existing-user', got %q", post.AuthorID)
	}
}

func TestFactory_ForMultiple(t *testing.T) {
	userFactory := New(func(seq int64) User {
		return User{
			ID:   fmt.Sprintf("user-%d", seq),
			Name: fmt.Sprintf("User %d", seq),
		}
	})

	postFactory := New(func(seq int64) Post {
		return Post{
			Title: fmt.Sprintf("Post %d", seq),
		}
	})

	// Create multiple posts, each with their own user
	posts := For(postFactory, userFactory, func(p *Post, u *User) {
		p.AuthorID = u.ID
	}).MakeMany(3)

	if len(posts) != 3 {
		t.Fatalf("expected 3 posts, got %d", len(posts))
	}

	// Each post should have a different author (different user created each time)
	authorIDs := make(map[string]bool)
	for i, p := range posts {
		if p.AuthorID == "" {
			t.Fatalf("post %d: expected AuthorID to be set", i)
		}
		authorIDs[p.AuthorID] = true
	}

	if len(authorIDs) != 3 {
		t.Fatalf("expected 3 unique author IDs, got %d", len(authorIDs))
	}
}

func TestFactory_ForModelMultiple(t *testing.T) {
	user := User{
		ID:   "shared-user",
		Name: "Shared User",
	}

	postFactory := New(func(seq int64) Post {
		return Post{
			Title: fmt.Sprintf("Post %d", seq),
		}
	})

	// Create multiple posts, all with the same user
	posts := ForModel(postFactory, &user, func(p *Post, u *User) {
		p.AuthorID = u.ID
	}).MakeMany(3)

	if len(posts) != 3 {
		t.Fatalf("expected 3 posts, got %d", len(posts))
	}

	// All posts should have the same author
	for i, p := range posts {
		if p.AuthorID != "shared-user" {
			t.Fatalf("post %d: expected AuthorID 'shared-user', got %q", i, p.AuthorID)
		}
	}
}

// RawJSON and WithRawDefaults Tests

func TestFactory_RawJSON(t *testing.T) {
	f := New(func(seq int64) User {
		return User{
			ID:    fmt.Sprintf("user-%d", seq),
			Name:  fmt.Sprintf("User %d", seq),
			Email: fmt.Sprintf("user%d@example.com", seq),
		}
	})

	jsonData, err := f.RawJSON()
	if err != nil {
		t.Fatalf("RawJSON failed: %v", err)
	}

	if len(jsonData) == 0 {
		t.Fatal("expected non-empty JSON data")
	}

	// Verify it's valid JSON
	var user User
	if err := json.Unmarshal(jsonData, &user); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}

	if user.ID != "user-1" {
		t.Fatalf("expected ID 'user-1', got %q", user.ID)
	}
}

func TestFactory_RawManyJSON(t *testing.T) {
	f := New(func(seq int64) User {
		return User{
			ID:   fmt.Sprintf("user-%d", seq),
			Name: fmt.Sprintf("User %d", seq),
		}
	})

	jsonData, err := f.RawManyJSON(3)
	if err != nil {
		t.Fatalf("RawManyJSON failed: %v", err)
	}

	var users []User
	if err := json.Unmarshal(jsonData, &users); err != nil {
		t.Fatalf("invalid JSON array: %v", err)
	}

	if len(users) != 3 {
		t.Fatalf("expected 3 users in JSON, got %d", len(users))
	}

	for i, u := range users {
		expectedID := fmt.Sprintf("user-%d", i+1)
		if u.ID != expectedID {
			t.Fatalf("user %d: expected ID %q, got %q", i, expectedID, u.ID)
		}
	}
}

func TestFactory_WithRawDefaults(t *testing.T) {
	type APIUser struct {
		ID       string
		Name     string
		Email    string
		Password string // Only for raw/API, not for persistence
	}

	f := New(func(seq int64) APIUser {
		return APIUser{
			ID:    fmt.Sprintf("user-%d", seq),
			Name:  fmt.Sprintf("User %d", seq),
			Email: fmt.Sprintf("user%d@example.com", seq),
		}
	}).WithRawDefaults(func(u *APIUser) {
		// Add password only for API/raw output
		u.Password = "test-password"
	})

	// Make should NOT have rawDefaults
	user := f.Make()
	if user.Password != "" {
		t.Fatalf("Make() should not apply rawDefaults, got password %q", user.Password)
	}

	// Raw SHOULD have rawDefaults
	rawUser := f.Raw()
	if rawUser.Password != "test-password" {
		t.Fatalf("Raw() should apply rawDefaults, expected 'test-password', got %q", rawUser.Password)
	}
}

func TestFactory_RawDefaultsWithJSON(t *testing.T) {
	type APIRequest struct {
		Username string
		Email    string
		Token    string // Only for API requests
	}

	f := New(func(seq int64) APIRequest {
		return APIRequest{
			Username: fmt.Sprintf("user%d", seq),
			Email:    fmt.Sprintf("user%d@example.com", seq),
		}
	}).WithRawDefaults(func(r *APIRequest) {
		r.Token = "auth-token-123"
	})

	jsonData, err := f.RawJSON()
	if err != nil {
		t.Fatalf("RawJSON failed: %v", err)
	}

	var req APIRequest
	if err := json.Unmarshal(jsonData, &req); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}

	if req.Token != "auth-token-123" {
		t.Fatalf("expected Token 'auth-token-123', got %q", req.Token)
	}
}

func TestFactory_CountedFactoryRawJSON(t *testing.T) {
	f := New(func(seq int64) User {
		return User{
			ID:   fmt.Sprintf("user-%d", seq),
			Name: fmt.Sprintf("User %d", seq),
		}
	})

	jsonData, err := f.Count(5).RawJSON()
	if err != nil {
		t.Fatalf("Count().RawJSON() failed: %v", err)
	}

	var users []User
	if err := json.Unmarshal(jsonData, &users); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}

	if len(users) != 5 {
		t.Fatalf("expected 5 users, got %d", len(users))
	}
}

func TestFactory_RawDefaultsWithTraits(t *testing.T) {
	type APIData struct {
		Name      string
		Value     int
		Timestamp string // Only for raw
	}

	f := New(func(seq int64) APIData {
		return APIData{
			Name:  "test",
			Value: int(seq),
		}
	}).WithRawDefaults(func(d *APIData) {
		d.Timestamp = "2024-01-01T00:00:00Z"
	})

	// Per-call traits should still override rawDefaults
	data := f.Raw(func(d *APIData) {
		d.Timestamp = "custom-timestamp"
	})

	if data.Timestamp != "custom-timestamp" {
		t.Fatalf("per-call traits should override rawDefaults, got %q", data.Timestamp)
	}
}

func TestFactory_RawJSONWithState(t *testing.T) {
	f := New(func(seq int64) User {
		return User{
			ID:   fmt.Sprintf("user-%d", seq),
			Name: "Regular User",
		}
	}).DefineState("admin", func(u *User) {
		u.Name = "Admin User"
	})

	jsonData, err := f.State("admin").RawJSON()
	if err != nil {
		t.Fatalf("State().RawJSON() failed: %v", err)
	}

	var user User
	if err := json.Unmarshal(jsonData, &user); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}

	if user.Name != "Admin User" {
		t.Fatalf("expected Name 'Admin User', got %q", user.Name)
	}
}

// Must* Variants Tests

func TestFactory_MustCreate(t *testing.T) {
	f := New(func(seq int64) User {
		return User{
			Name:  fmt.Sprintf("User %d", seq),
			Email: fmt.Sprintf("user%d@example.com", seq),
		}
	}).WithPersist(func(ctx context.Context, u *User) (*User, error) {
		u.ID = "saved-id"
		return u, nil
	})

	ctx := context.Background()
	user := f.MustCreate(ctx)

	if user.ID != "saved-id" {
		t.Fatalf("expected ID 'saved-id', got %q", user.ID)
	}
}

func TestFactory_MustCreatePanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected MustCreate to panic on error")
		}
	}()

	f := New(func(seq int64) User {
		return User{Name: "Test"}
	}).WithPersist(func(ctx context.Context, u *User) (*User, error) {
		return nil, fmt.Errorf("persist error")
	})

	ctx := context.Background()
	_ = f.MustCreate(ctx) // Should panic
}

func TestFactory_MustCreateMany(t *testing.T) {
	f := New(func(seq int64) User {
		return User{
			Name:  fmt.Sprintf("User %d", seq),
			Email: fmt.Sprintf("user%d@example.com", seq),
		}
	}).WithPersist(func(ctx context.Context, u *User) (*User, error) {
		u.ID = "saved"
		return u, nil
	})

	ctx := context.Background()
	users := f.MustCreateMany(ctx, 3)

	if len(users) != 3 {
		t.Fatalf("expected 3 users, got %d", len(users))
	}
}

func TestFactory_MustRawJSON(t *testing.T) {
	f := New(func(seq int64) User {
		return User{
			ID:   fmt.Sprintf("user-%d", seq),
			Name: "Test User",
		}
	})

	jsonData := f.MustRawJSON()
	if len(jsonData) == 0 {
		t.Fatal("expected non-empty JSON data")
	}

	var user User
	if err := json.Unmarshal(jsonData, &user); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
}

func TestFactory_CountedFactoryMustCreate(t *testing.T) {
	f := New(func(seq int64) User {
		return User{
			Name:  fmt.Sprintf("User %d", seq),
			Email: fmt.Sprintf("user%d@example.com", seq),
		}
	}).WithPersist(func(ctx context.Context, u *User) (*User, error) {
		u.ID = "saved"
		return u, nil
	})

	ctx := context.Background()
	users := f.Count(5).MustCreate(ctx)

	if len(users) != 5 {
		t.Fatalf("expected 5 users, got %d", len(users))
	}
}

// Tap() Tests

func TestFactory_Tap(t *testing.T) {
	callCount := 0
	var capturedNames []string

	f := New(func(seq int64) User {
		return User{
			Name:  fmt.Sprintf("User %d", seq),
			Email: fmt.Sprintf("user%d@example.com", seq),
		}
	}).Tap(func(u User) {
		callCount++
		capturedNames = append(capturedNames, u.Name)
	})

	// Make should call tap
	user := f.Make()
	if callCount != 1 {
		t.Fatalf("expected tap to be called once, got %d", callCount)
	}

	// MakeMany should call tap for each
	users := f.MakeMany(3)
	if callCount != 4 { // 1 from Make + 3 from MakeMany
		t.Fatalf("expected tap to be called 4 times total, got %d", callCount)
	}

	if len(users) != 3 {
		t.Fatalf("expected 3 users, got %d", len(users))
	}

	// Verify captured names
	if capturedNames[0] != user.Name {
		t.Fatalf("expected first captured name %q, got %q", user.Name, capturedNames[0])
	}
}

func TestFactory_TapWithRaw(t *testing.T) {
	callCount := 0

	f := New(func(seq int64) User {
		return User{Name: fmt.Sprintf("User %d", seq)}
	}).Tap(func(u User) {
		callCount++
	})

	// Raw should also call tap
	_ = f.Raw()
	if callCount != 1 {
		t.Fatalf("expected tap to be called once for Raw, got %d", callCount)
	}

	// RawMany should call tap for each
	_ = f.RawMany(2)
	if callCount != 3 { // 1 from Raw + 2 from RawMany
		t.Fatalf("expected tap to be called 3 times total, got %d", callCount)
	}
}

// When() / Unless() Tests

func TestFactory_When(t *testing.T) {
	// When true - trait should apply
	f := New(func(seq int64) User {
		return User{Name: "User"}
	}).When(true, func(u *User) {
		u.ID = "applied"
	})

	user := f.Make()
	if user.ID != "applied" {
		t.Fatalf("expected ID 'applied', got %q", user.ID)
	}

	// When false - trait should NOT apply
	f2 := New(func(seq int64) User {
		return User{Name: "User"}
	}).When(false, func(u *User) {
		u.ID = "applied"
	})

	user2 := f2.Make()
	if user2.ID != "" {
		t.Fatalf("expected empty ID, got %q", user2.ID)
	}
}

func TestFactory_Unless(t *testing.T) {
	// Unless false - trait should apply
	f := New(func(seq int64) User {
		return User{Name: "User"}
	}).Unless(false, func(u *User) {
		u.ID = "applied"
	})

	user := f.Make()
	if user.ID != "applied" {
		t.Fatalf("expected ID 'applied', got %q", user.ID)
	}

	// Unless true - trait should NOT apply
	f2 := New(func(seq int64) User {
		return User{Name: "User"}
	}).Unless(true, func(u *User) {
		u.ID = "applied"
	})

	user2 := f2.Make()
	if user2.ID != "" {
		t.Fatalf("expected empty ID, got %q", user2.ID)
	}
}

func TestFactory_WhenUnlessChaining(t *testing.T) {
	isProd := false
	isTest := true

	f := New(func(seq int64) User {
		return User{Name: "User"}
	}).
		When(isProd, func(u *User) {
			u.Email = "prod@example.com"
		}).
		Unless(isProd, func(u *User) {
			u.Email = "test@example.com"
		}).
		When(isTest, func(u *User) {
			u.ID = "test-id"
		})

	user := f.Make()

	if user.Email != "test@example.com" {
		t.Fatalf("expected Email 'test@example.com', got %q", user.Email)
	}

	if user.ID != "test-id" {
		t.Fatalf("expected ID 'test-id', got %q", user.ID)
	}
}

// Clone() Tests

func TestFactory_Clone(t *testing.T) {
	baseFactory := New(func(seq int64) User {
		return User{
			Name:  fmt.Sprintf("User %d", seq),
			Email: fmt.Sprintf("user%d@example.com", seq),
		}
	}).WithDefaults(func(u *User) {
		u.ID = "default"
	})

	// Clone and modify
	adminFactory := baseFactory.Clone().WithTraits(func(u *User) {
		u.ID = "admin"
	})

	// Original should be unchanged
	user := baseFactory.Make()
	if user.ID != "default" {
		t.Fatalf("base factory: expected ID 'default', got %q", user.ID)
	}

	// Clone should have new trait
	admin := adminFactory.Make()
	if admin.ID != "admin" {
		t.Fatalf("admin factory: expected ID 'admin', got %q", admin.ID)
	}
}

func TestFactory_CloneWithStates(t *testing.T) {
	baseFactory := New(func(seq int64) User {
		return User{Name: "User"}
	}).DefineState("admin", func(u *User) {
		u.ID = "admin"
	})

	// Clone should have the states
	cloned := baseFactory.Clone()

	admin := cloned.State("admin").Make()
	if admin.ID != "admin" {
		t.Fatalf("expected ID 'admin', got %q", admin.ID)
	}
}

func TestFactory_CloneSequenceReset(t *testing.T) {
	baseFactory := New(func(seq int64) User {
		return User{Name: fmt.Sprintf("User %d", seq)}
	})

	// Create from base
	u1 := baseFactory.Make()
	if u1.Name != "User 1" {
		t.Fatalf("expected 'User 1', got %q", u1.Name)
	}

	u2 := baseFactory.Make()
	if u2.Name != "User 2" {
		t.Fatalf("expected 'User 2', got %q", u2.Name)
	}

	// Clone resets sequence
	cloned := baseFactory.Clone()
	u3 := cloned.Make()
	if u3.Name != "User 1" {
		t.Fatalf("cloned factory: expected 'User 1', got %q", u3.Name)
	}

	// Base factory continues
	u4 := baseFactory.Make()
	if u4.Name != "User 3" {
		t.Fatalf("base factory: expected 'User 3', got %q", u4.Name)
	}
}

// Advanced Relationship Tests

func TestFactory_Has(t *testing.T) {
	userFactory := New(func(seq int64) User {
		return User{
			ID:   fmt.Sprintf("user-%d", seq),
			Name: fmt.Sprintf("User %d", seq),
		}
	})

	postFactory := New(func(seq int64) Post {
		return Post{
			Title: fmt.Sprintf("Post %d", seq),
		}
	})

	// Use Has to create user with 3 posts
	user, posts := Has(userFactory, postFactory, 3, func(u *User, p *Post) {
		p.AuthorID = u.ID
	}).Make()

	if user.ID == "" {
		t.Fatal("expected user ID to be set")
	}

	if len(posts) != 3 {
		t.Fatalf("expected 3 posts, got %d", len(posts))
	}

	// All posts should belong to the user
	for i, post := range posts {
		if post.AuthorID != user.ID {
			t.Fatalf("post %d: expected AuthorID %q, got %q", i, user.ID, post.AuthorID)
		}
	}
}

func TestFactory_HasWithCreate(t *testing.T) {
	userFactory := New(func(seq int64) User {
		return User{
			Name:  fmt.Sprintf("User %d", seq),
			Email: fmt.Sprintf("user%d@example.com", seq),
		}
	}).WithPersist(func(ctx context.Context, u *User) (*User, error) {
		u.ID = fmt.Sprintf("user-%s", u.Name)
		return u, nil
	})

	postFactory := New(func(seq int64) Post {
		return Post{
			Title: fmt.Sprintf("Post %d", seq),
		}
	}).WithPersist(func(ctx context.Context, p *Post) (*Post, error) {
		p.ID = fmt.Sprintf("post-%s", p.Title)
		return p, nil
	})

	ctx := context.Background()
	user, posts, err := Has(userFactory, postFactory, 5, func(u *User, p *Post) {
		p.AuthorID = u.ID
	}).Create(ctx)

	if err != nil {
		t.Fatal(err)
	}

	if user == nil {
		t.Fatal("expected user to be created")
	}

	if len(posts) != 5 {
		t.Fatalf("expected 5 posts, got %d", len(posts))
	}

	// Verify all posts have IDs (were persisted)
	for i, post := range posts {
		if post.ID == "" {
			t.Fatalf("post %d: expected ID to be set", i)
		}
		if post.AuthorID != user.ID {
			t.Fatalf("post %d: expected AuthorID %q, got %q", i, user.ID, post.AuthorID)
		}
	}
}

func TestFactory_HasMustCreate(t *testing.T) {
	userFactory := New(func(seq int64) User {
		return User{Name: fmt.Sprintf("User %d", seq)}
	}).WithPersist(func(ctx context.Context, u *User) (*User, error) {
		u.ID = "saved"
		return u, nil
	})

	postFactory := New(func(seq int64) Post {
		return Post{Title: fmt.Sprintf("Post %d", seq)}
	}).WithPersist(func(ctx context.Context, p *Post) (*Post, error) {
		p.ID = "saved"
		return p, nil
	})

	ctx := context.Background()
	user, posts := Has(userFactory, postFactory, 2, func(u *User, p *Post) {
		p.AuthorID = u.ID
	}).MustCreate(ctx)

	if user == nil {
		t.Fatal("expected user")
	}

	if len(posts) != 2 {
		t.Fatalf("expected 2 posts, got %d", len(posts))
	}
}

func TestFactory_Recycle(t *testing.T) {
	user := User{
		ID:   "recycled-user",
		Name: "Recycled User",
	}

	postFactory := New(func(seq int64) Post {
		return Post{
			Title: fmt.Sprintf("Post %d", seq),
		}
	})

	// All posts should use the same user
	posts := Recycle(postFactory, &user, func(p *Post, u *User) {
		p.AuthorID = u.ID
	}).MakeMany(5)

	if len(posts) != 5 {
		t.Fatalf("expected 5 posts, got %d", len(posts))
	}

	// All should have same author
	for i, post := range posts {
		if post.AuthorID != "recycled-user" {
			t.Fatalf("post %d: expected AuthorID 'recycled-user', got %q", i, post.AuthorID)
		}
	}
}

func TestFactory_RecycleWithCreate(t *testing.T) {
	userFactory := New(func(seq int64) User {
		return User{
			Name:  fmt.Sprintf("User %d", seq),
			Email: fmt.Sprintf("user%d@example.com", seq),
		}
	}).WithPersist(func(ctx context.Context, u *User) (*User, error) {
		u.ID = "user-saved"
		return u, nil
	})

	postFactory := New(func(seq int64) Post {
		return Post{
			Title: fmt.Sprintf("Post %d", seq),
		}
	}).WithPersist(func(ctx context.Context, p *Post) (*Post, error) {
		p.ID = fmt.Sprintf("post-%s", p.Title)
		return p, nil
	})

	ctx := context.Background()

	// Create one user
	user, err := userFactory.Create(ctx)
	if err != nil {
		t.Fatal(err)
	}

	// Create 3 posts, all recycling the same user
	posts, err := Recycle(postFactory, user, func(p *Post, u *User) {
		p.AuthorID = u.ID
	}).CreateMany(ctx, 3)

	if err != nil {
		t.Fatal(err)
	}

	if len(posts) != 3 {
		t.Fatalf("expected 3 posts, got %d", len(posts))
	}

	// All posts should have same author
	for i, post := range posts {
		if post.AuthorID != user.ID {
			t.Fatalf("post %d: expected AuthorID %q, got %q", i, user.ID, post.AuthorID)
		}
	}
}

func TestFactory_RecycleWithCount(t *testing.T) {
	user := User{ID: "shared", Name: "Shared User"}

	postFactory := New(func(seq int64) Post {
		return Post{Title: fmt.Sprintf("Post %d", seq)}
	})

	// Use Recycle with Count() fluent API
	posts := Recycle(postFactory, &user, func(p *Post, u *User) {
		p.AuthorID = u.ID
	}).Count(10).Make()

	if len(posts) != 10 {
		t.Fatalf("expected 10 posts, got %d", len(posts))
	}

	for i, post := range posts {
		if post.AuthorID != "shared" {
			t.Fatalf("post %d: expected shared author", i)
		}
	}
}

// HasAttached Tests

type UserRole struct {
	UserID string
	RoleID string
	Active bool
}

type Role struct {
	ID   string
	Name string
}

func TestFactory_HasAttached(t *testing.T) {
	userFactory := New(func(seq int64) User {
		return User{
			ID:   fmt.Sprintf("user-%d", seq),
			Name: fmt.Sprintf("User %d", seq),
		}
	}).WithPersist(func(ctx context.Context, u *User) (*User, error) {
		return u, nil
	})

	roleFactory := New(func(seq int64) Role {
		return Role{
			ID:   fmt.Sprintf("role-%d", seq),
			Name: fmt.Sprintf("Role %d", seq),
		}
	}).WithPersist(func(ctx context.Context, r *Role) (*Role, error) {
		return r, nil
	})

	userRoleFactory := New(func(seq int64) UserRole {
		return UserRole{Active: true}
	}).WithPersist(func(ctx context.Context, ur *UserRole) (*UserRole, error) {
		return ur, nil
	})

	ctx := context.Background()

	// Create user with 3 roles attached
	user, roles, pivots, err := HasAttached(
		userFactory,
		roleFactory,
		userRoleFactory,
		3,
		func(ur *UserRole, u *User, r *Role) {
			ur.UserID = u.ID
			ur.RoleID = r.ID
		},
	).Create(ctx)

	if err != nil {
		t.Fatal(err)
	}

	if user == nil {
		t.Fatal("expected user to be created")
	}

	if len(roles) != 3 {
		t.Fatalf("expected 3 roles, got %d", len(roles))
	}

	if len(pivots) != 3 {
		t.Fatalf("expected 3 pivot records, got %d", len(pivots))
	}

	// Verify pivot records have correct IDs
	for i, pivot := range pivots {
		if pivot.UserID != user.ID {
			t.Fatalf("pivot %d: expected UserID %q, got %q", i, user.ID, pivot.UserID)
		}
		if pivot.RoleID != roles[i].ID {
			t.Fatalf("pivot %d: expected RoleID %q, got %q", i, roles[i].ID, pivot.RoleID)
		}
		if !pivot.Active {
			t.Fatalf("pivot %d: expected Active to be true", i)
		}
	}
}

func TestFactory_HasAttachedMake(t *testing.T) {
	userFactory := New(func(seq int64) User {
		return User{
			ID:   fmt.Sprintf("user-%d", seq),
			Name: fmt.Sprintf("User %d", seq),
		}
	})

	roleFactory := New(func(seq int64) Role {
		return Role{
			ID:   fmt.Sprintf("role-%d", seq),
			Name: fmt.Sprintf("Role %d", seq),
		}
	})

	userRoleFactory := New(func(seq int64) UserRole {
		return UserRole{Active: false}
	})

	// Make (in-memory) user with 2 roles
	user, roles, pivots := HasAttached(
		userFactory,
		roleFactory,
		userRoleFactory,
		2,
		func(ur *UserRole, u *User, r *Role) {
			ur.UserID = u.ID
			ur.RoleID = r.ID
			ur.Active = true
		},
	).Make()

	if user.ID == "" {
		t.Fatal("expected user ID to be set")
	}

	if len(roles) != 2 {
		t.Fatalf("expected 2 roles, got %d", len(roles))
	}

	if len(pivots) != 2 {
		t.Fatalf("expected 2 pivots, got %d", len(pivots))
	}

	// Verify pivots are linked correctly
	for i := 0; i < 2; i++ {
		if pivots[i].UserID != user.ID {
			t.Fatalf("pivot %d: expected UserID %q, got %q", i, user.ID, pivots[i].UserID)
		}
		if pivots[i].RoleID != roles[i].ID {
			t.Fatalf("pivot %d: expected RoleID %q, got %q", i, roles[i].ID, pivots[i].RoleID)
		}
		if !pivots[i].Active {
			t.Fatalf("pivot %d: expected Active true, got false", i)
		}
	}
}

func TestFactory_HasAttachedMustCreate(t *testing.T) {
	userFactory := New(func(seq int64) User {
		return User{Name: fmt.Sprintf("User %d", seq)}
	}).WithPersist(func(ctx context.Context, u *User) (*User, error) {
		u.ID = "user-saved"
		return u, nil
	})

	roleFactory := New(func(seq int64) Role {
		return Role{Name: fmt.Sprintf("Role %d", seq)}
	}).WithPersist(func(ctx context.Context, r *Role) (*Role, error) {
		r.ID = fmt.Sprintf("role-%s", r.Name)
		return r, nil
	})

	userRoleFactory := New(func(seq int64) UserRole {
		return UserRole{}
	}).WithPersist(func(ctx context.Context, ur *UserRole) (*UserRole, error) {
		return ur, nil
	})

	ctx := context.Background()
	user, roles, pivots := HasAttached(
		userFactory,
		roleFactory,
		userRoleFactory,
		2,
		func(ur *UserRole, u *User, r *Role) {
			ur.UserID = u.ID
			ur.RoleID = r.ID
			ur.Active = true
		},
	).MustCreate(ctx)

	if user == nil {
		t.Fatal("expected user")
	}

	if len(roles) != 2 {
		t.Fatalf("expected 2 roles, got %d", len(roles))
	}

	if len(pivots) != 2 {
		t.Fatalf("expected 2 pivots, got %d", len(pivots))
	}
}
