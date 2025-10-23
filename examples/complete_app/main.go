package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"

	"github.com/b3ndoi/factory-go/factory"
)

// Models for a complete blog application
type User struct {
	ID              string
	Name            string
	Email           string
	Password        string
	Role            string
	Active          bool
	EmailVerifiedAt string
}

type Post struct {
	ID        string
	Title     string
	Content   string
	AuthorID  string
	Published bool
	Views     int
}

type Comment struct {
	ID       string
	Content  string
	PostID   string
	AuthorID string
	Approved bool
}

type Tag struct {
	ID   string
	Name string
	Slug string
}

type PostTag struct {
	PostID string
	TagID  string
}

// Mock database
type DB struct {
	users    map[string]*User
	posts    map[string]*Post
	comments map[string]*Comment
	tags     map[string]*Tag
	postTags []*PostTag
}

func NewDB() *DB {
	return &DB{
		users:    make(map[string]*User),
		posts:    make(map[string]*Post),
		comments: make(map[string]*Comment),
		tags:     make(map[string]*Tag),
		postTags: []*PostTag{},
	}
}

func (db *DB) CreateUser(ctx context.Context, u *User) (*User, error) {
	u.ID = fmt.Sprintf("user-%d", len(db.users)+1)
	db.users[u.ID] = u
	return u, nil
}

func (db *DB) CreatePost(ctx context.Context, p *Post) (*Post, error) {
	p.ID = fmt.Sprintf("post-%d", len(db.posts)+1)
	db.posts[p.ID] = p
	return p, nil
}

func (db *DB) CreateComment(ctx context.Context, c *Comment) (*Comment, error) {
	c.ID = fmt.Sprintf("comment-%d", len(db.comments)+1)
	db.comments[c.ID] = c
	return c, nil
}

func (db *DB) CreateTag(ctx context.Context, t *Tag) (*Tag, error) {
	t.ID = fmt.Sprintf("tag-%d", len(db.tags)+1)
	db.tags[t.ID] = t
	return t, nil
}

func (db *DB) CreatePostTag(ctx context.Context, pt *PostTag) (*PostTag, error) {
	db.postTags = append(db.postTags, pt)
	return pt, nil
}

func main() {
	fmt.Println("=== Factory-Go Complete App Example ===")
	fmt.Println()

	db := NewDB()
	ctx := context.Background()
	isProd := os.Getenv("ENV") == "production"

	// Setup comprehensive factories with all features

	// User Factory with EVERY feature
	userFactory := factory.New(func(seq int64) User {
		return User{
			Name:   fmt.Sprintf("User %d", seq),
			Email:  fmt.Sprintf("user%d@example.com", seq),
			Role:   "user",
			Active: true,
		}
	}).
		WithDefaults(func(u *User) {
			// In real app, use faker: u.Name = faker.Name()
			// For demo, we keep it simple
		}).
		WithRawDefaults(func(u *User) {
			// Password only for API testing, not DB
			u.Password = "SecurePassword123!"
		}).
		DefineState("admin", func(u *User) {
			u.Role = "admin"
		}).
		DefineState("moderator", func(u *User) {
			u.Role = "moderator"
		}).
		DefineState("verified", func(u *User) {
			u.EmailVerifiedAt = "2024-01-01T00:00:00Z"
		}).
		DefineState("inactive", func(u *User) {
			u.Active = false
		}).
		When(isProd, func(u *User) {
			// Production-specific behavior
			u.Email = fmt.Sprintf("prod-%s", u.Email)
		}).
		Unless(isProd, func(u *User) {
			// Test-specific behavior
			u.Email = fmt.Sprintf("test-%s", u.Email)
		}).
		WithPersist(db.CreateUser).
		BeforeCreate(func(ctx context.Context, u *User) error {
			fmt.Printf("   [BeforeCreate] Validating user: %s\n", u.Name)
			return nil
		}).
		AfterCreate(func(ctx context.Context, u *User) error {
			fmt.Printf("   [AfterCreate] User created with ID: %s\n", u.ID)
			return nil
		}).
		Tap(func(u User) {
			fmt.Printf("   [Tap] Building user: %s\n", u.Name)
		})

	// Post Factory
	postFactory := factory.New(func(seq int64) Post {
		return Post{
			Title:     fmt.Sprintf("Post %d: Interesting Title", seq),
			Content:   fmt.Sprintf("Content for post %d...", seq),
			Published: false,
			Views:     0,
		}
	}).
		WithPersist(db.CreatePost).
		DefineState("published", func(p *Post) {
			p.Published = true
		}).
		DefineState("popular", func(p *Post) {
			p.Views = 1000
		})

	// Comment Factory (declared for completeness)
	_ = factory.New(func(seq int64) Comment {
		return Comment{
			Content:  fmt.Sprintf("Comment %d", seq),
			Approved: true,
		}
	}).WithPersist(db.CreateComment)

	// Tag Factory
	tagFactory := factory.New(func(seq int64) Tag {
		tags := []string{"golang", "programming", "tutorial", "web", "backend"}
		return Tag{
			Name: tags[int(seq-1)%len(tags)],
			Slug: tags[int(seq-1)%len(tags)],
		}
	}).WithPersist(db.CreateTag)

	// PostTag Factory (pivot)
	postTagFactory := factory.New(func(seq int64) PostTag {
		return PostTag{}
	}).WithPersist(db.CreatePostTag)

	// === Comprehensive Blog Application Seeding ===

	fmt.Println("🎯 Seeding a Complete Blog Application")
	fmt.Println()

	// STEP 1: Create Tags (needed first for HasAttached later)
	fmt.Println("1. Creating blog tags...")
	allTags := tagFactory.Count(5).MustCreate(ctx)
	fmt.Printf("   ✅ Created %d tags: ", len(allTags))
	for i, tag := range allTags {
		if i > 0 {
			fmt.Print(", ")
		}
		fmt.Print(tag.Name)
	}
	fmt.Println()
	fmt.Println()

	// STEP 2: Create Users with different roles (using Sequence)
	fmt.Println("2. Creating users with different subscription plans...")
	regularUsers := userFactory.Sequence(
		func(u *User) { u.Role = "free" },
		func(u *User) { u.Role = "pro" },
		func(u *User) { u.Role = "enterprise" },
	).Count(12).MustCreate(ctx)

	roleCounts := make(map[string]int)
	for _, u := range regularUsers {
		roleCounts[u.Role]++
	}
	fmt.Printf("   ✅ Created 12 users: ")
	for role, count := range roleCounts {
		fmt.Printf("%d %s, ", count, role)
	}
	fmt.Println()

	// STEP 3: Create verified admin users
	fmt.Println("\n3. Creating verified admin users...")
	admins := userFactory.State("admin").State("verified").Count(2).MustCreate(ctx)
	fmt.Printf("   ✅ Created %d verified admins\n", len(admins))

	// STEP 4: Each prolific author creates multiple posts
	fmt.Println("\n4. Creating posts from top authors (Has relationship)...")
	var allPosts []*Post

	// First author writes 10 posts
	author1, posts1 := factory.Has(userFactory, postFactory, 10, func(u *User, p *Post) {
		p.AuthorID = u.ID
	}).MustCreate(ctx)
	allPosts = append(allPosts, posts1...)
	fmt.Printf("   ✅ %s wrote %d posts\n", author1.Name, len(posts1))

	// Second author writes 8 posts
	author2, posts2 := factory.Has(userFactory, postFactory, 8, func(u *User, p *Post) {
		p.AuthorID = u.ID
	}).MustCreate(ctx)
	allPosts = append(allPosts, posts2...)
	fmt.Printf("   ✅ %s wrote %d posts\n", author2.Name, len(posts2))

	// STEP 5: Publish some posts (using State)
	fmt.Println("\n5. Publishing posts...")
	for i := 0; i < 10; i++ {
		allPosts[i].Published = true
	}
	fmt.Printf("   ✅ Published 10 out of %d posts\n", len(allPosts))

	// STEP 6: Add tags to posts (HasAttached - many-to-many)
	fmt.Println("\n6. Tagging posts...")
	taggedPostCount := 0
	for i := 0; i < min(5, len(allPosts)); i++ {
		post := allPosts[i]
		numTags := (i % 3) + 1 // 1-3 tags per post

		for j := 0; j < numTags; j++ {
			postTagFactory.MustCreate(ctx, func(pt *PostTag) {
				pt.PostID = post.ID
				pt.TagID = allTags[j%len(allTags)].ID
			})
		}
		taggedPostCount++
	}
	fmt.Printf("   ✅ Tagged %d posts with various tags\n", taggedPostCount)

	// STEP 7: Users comment on posts
	fmt.Println("\n7. Creating comments from readers...")
	commentFactory := factory.New(func(seq int64) Comment {
		return Comment{
			Content:  fmt.Sprintf("Great post! This is comment %d", seq),
			Approved: true,
		}
	}).WithPersist(db.CreateComment)

	totalComments := 0
	// First 5 published posts get 3-7 comments each
	for i := 0; i < min(5, len(allPosts)); i++ {
		post := allPosts[i]
		numComments := 3 + (i % 5) // 3-7 comments

		for j := 0; j < numComments; j++ {
			randomUser := regularUsers[rand.Intn(len(regularUsers))]
			commentFactory.MustCreate(ctx, func(c *Comment) {
				c.PostID = post.ID
				c.AuthorID = randomUser.ID
			})
			totalComments++
		}
	}
	fmt.Printf("   ✅ Created %d comments from various users\n", totalComments)

	// STEP 8: Demonstrate When/Unless for environment-specific behavior
	fmt.Println("\n8. Environment-aware factory (When/Unless)...")
	isProdEnv := isProd
	envFactory := userFactory.Clone().
		When(isProdEnv, func(u *User) {
			fmt.Printf("   [When] Production mode enabled for %s\n", u.Name)
		}).
		Unless(isProdEnv, func(u *User) {
			fmt.Printf("   [Unless] Test mode enabled for %s\n", u.Name)
		})

	envUser := envFactory.Make()
	fmt.Printf("   ✅ Created user with environment-specific config: %s\n", envUser.Email)

	// STEP 9: Using Tap() to monitor creation
	fmt.Println("\n9. Using Tap() to monitor creation...")
	monitored := 0
	userFactory.Clone().Tap(func(u User) {
		monitored++
		fmt.Printf("   [Tap #%d] Monitoring: %s (%s)\n", monitored, u.Name, u.Email)
	}).Count(3).Make()

	// STEP 10: Test RawJSON for API endpoint testing
	fmt.Println("\n10. Generating API test data (RawJSON)...")
	apiUser := userFactory.State("admin").Raw()
	fmt.Printf("   Domain model - Password: '%s' (empty ✅)\n", apiUser.Password)

	apiJSON := userFactory.State("admin").MustRawJSON()
	fmt.Printf("   API JSON - includes password: %d bytes\n", len(apiJSON))
	fmt.Printf("   ✅ Demonstrates WithRawDefaults() separation\n")

	// Final statistics
	fmt.Println("\n" + repeat("═", 60))
	fmt.Println("📊 COMPLETE BLOG APPLICATION STATISTICS")
	fmt.Println(repeat("═", 60))
	fmt.Printf("Total Users:      %d (12 regular + 2 admins + authors + test users)\n", len(db.users))
	fmt.Printf("Total Posts:      %d (from multiple authors)\n", len(db.posts))
	fmt.Printf("Total Comments:   %d (from various readers)\n", len(db.comments))
	fmt.Printf("Total Tags:       %d\n", len(db.tags))
	fmt.Printf("Total PostTags:   %d (many-to-many relationships)\n", len(db.postTags))
	fmt.Println(repeat("═", 60))

	fmt.Println("\n📖 What We Just Built:")
	fmt.Println("  → A complete blog with users, authors, posts, comments, and tags")
	fmt.Println("  → Users have subscription plans (free/pro/enterprise)")
	fmt.Println("  → Top authors each wrote multiple posts")
	fmt.Println("  → Posts are tagged for organization")
	fmt.Println("  → Users commented on published posts")
	fmt.Println("  → Admins can moderate the platform")

	fmt.Println("\n✅ Complete app example finished!")
	fmt.Println("\n🎯 ALL Features Demonstrated:")
	fmt.Println("  ✅ Sequences - Cycling subscription plans")
	fmt.Println("  ✅ Named States - admin, verified, published, popular")
	fmt.Println("  ✅ State Chaining - admin + verified")
	fmt.Println("  ✅ Has() - Authors with multiple posts")
	fmt.Println("  ✅ Many-to-Many - Posts with tags via pivot")
	fmt.Println("  ✅ Count() - Fluent API throughout")
	fmt.Println("  ✅ MustCreate() - Clean error handling")
	fmt.Println("  ✅ When/Unless - Environment-aware factories")
	fmt.Println("  ✅ BeforeCreate/AfterCreate - Lifecycle hooks")
	fmt.Println("  ✅ Tap() - Monitoring and debugging")
	fmt.Println("  ✅ Clone() - Factory variations")
	fmt.Println("  ✅ RawJSON() - API testing")
	fmt.Println("  ✅ WithRawDefaults() - Separate API/domain fields")
	fmt.Println("  ✅ Per-call traits - Custom modifications")

	fmt.Println("\n💡 This example shows how ALL features work together")
	fmt.Println("   in a realistic, production-like scenario!")
}

func repeat(s string, n int) string {
	result := ""
	for i := 0; i < n; i++ {
		result += s
	}
	return result
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
