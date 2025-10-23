package main

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/b3ndoi/factory-go/factory"
)

// Models
type User struct {
	ID    string
	Name  string
	Email string
}

type Post struct {
	ID       string
	Title    string
	Content  string
	AuthorID string
}

type Comment struct {
	ID       string
	Content  string
	PostID   string
	AuthorID string
}

type Role struct {
	ID   string
	Name string
}

type UserRole struct {
	UserID    string
	RoleID    string
	GrantedAt string
}

// Mock database (in real app, this would be your actual DB)
type MockDB struct {
	users     []*User
	posts     []*Post
	comments  []*Comment
	roles     []*Role
	userRoles []*UserRole
}

func (db *MockDB) CreateUser(ctx context.Context, u *User) (*User, error) {
	u.ID = fmt.Sprintf("user-%d", len(db.users)+1)
	db.users = append(db.users, u)
	return u, nil
}

func (db *MockDB) CreatePost(ctx context.Context, p *Post) (*Post, error) {
	p.ID = fmt.Sprintf("post-%d", len(db.posts)+1)
	db.posts = append(db.posts, p)
	return p, nil
}

func (db *MockDB) CreateComment(ctx context.Context, c *Comment) (*Comment, error) {
	c.ID = fmt.Sprintf("comment-%d", len(db.comments)+1)
	db.comments = append(db.comments, c)
	return c, nil
}

func (db *MockDB) CreateRole(ctx context.Context, r *Role) (*Role, error) {
	r.ID = fmt.Sprintf("role-%d", len(db.roles)+1)
	db.roles = append(db.roles, r)
	return r, nil
}

func (db *MockDB) CreateUserRole(ctx context.Context, ur *UserRole) (*UserRole, error) {
	db.userRoles = append(db.userRoles, ur)
	return ur, nil
}

func main() {
	fmt.Println("=== Factory-Go Database Seeding Example ===")
	fmt.Println()

	db := &MockDB{}
	ctx := context.Background()

	// Setup factories
	userFactory := factory.New(func(seq int64) User {
		return User{
			Name:  fmt.Sprintf("User %d", seq),
			Email: fmt.Sprintf("user%d@example.com", seq),
		}
	}).WithPersist(db.CreateUser)

	postFactory := factory.New(func(seq int64) Post {
		return Post{
			Title:   fmt.Sprintf("Post %d - Interesting Title", seq),
			Content: fmt.Sprintf("This is the content for post %d...", seq),
		}
	}).WithPersist(db.CreatePost)

	commentFactory := factory.New(func(seq int64) Comment {
		return Comment{
			Content: fmt.Sprintf("Comment %d content", seq),
		}
	}).WithPersist(db.CreateComment)

	roleFactory := factory.New(func(seq int64) Role {
		return Role{
			Name: []string{"Admin", "Editor", "Viewer"}[seq%3],
		}
	}).WithPersist(db.CreateRole)

	userRoleFactory := factory.New(func(seq int64) UserRole {
		return UserRole{
			GrantedAt: "2024-01-01",
		}
	}).WithPersist(db.CreateUserRole)

	// === Seeding Examples ===

	// 1. FIRST: Create roles (they should exist before assignment)
	fmt.Println("1. Creating roles...")
	allRoles := roleFactory.Count(3).MustCreate(ctx)
	fmt.Printf("   ✅ Created %d roles: ", len(allRoles))
	for i, r := range allRoles {
		if i > 0 {
			fmt.Print(", ")
		}
		fmt.Print(r.Name)
	}
	fmt.Println()

	// 2. Setup AfterCreate hook for automatic role assignment
	fmt.Println("\n2. Setting up AfterCreate hook for automatic role assignment...")
	roleIndex := 0
	userFactory = userFactory.AfterCreate(func(ctx context.Context, u *User) error {
		// Assign a role to each user (cycling through available roles)
		role := allRoles[roleIndex%len(allRoles)]
		roleIndex++

		_, err := userRoleFactory.Create(ctx, func(ur *UserRole) {
			ur.UserID = u.ID
			ur.RoleID = role.ID
		})

		if err == nil {
			fmt.Printf("   [AfterCreate] Auto-assigned '%s' to user '%s'\n", role.Name, u.Name)
		}
		return err
	})
	fmt.Println("   ✅ AfterCreate hook configured - every user will get a role!")

	// 3. Create regular users (roles AUTO-ASSIGNED via AfterCreate!)
	fmt.Println("\n3. Creating 10 regular users...")
	regularUsers := userFactory.Count(10).MustCreate(ctx)
	fmt.Printf("   ✅ Created %d users (each with a role auto-assigned!)\n", len(regularUsers))

	// 4. Create 3 more users (also auto-assigned roles)
	fmt.Println("\n4. Creating 3 more users...")
	moreUsers := userFactory.Count(3).MustCreate(ctx)
	fmt.Printf("   ✅ Created %d more users (with roles auto-assigned!)\n", len(moreUsers))

	// 5. Has() - Create user with 5 posts (user also gets role via AfterCreate)
	fmt.Println("\n5. Creating user with 5 posts (Has relationship)...")
	userWithPosts, posts := factory.Has(userFactory, postFactory, 5, func(u *User, p *Post) {
		p.AuthorID = u.ID
	}).MustCreate(ctx)
	fmt.Printf("   ✅ Created user '%s' with %d posts (and auto-assigned role)\n", userWithPosts.Name, len(posts))

	// 6. Recycle() - Create 10 posts for same user
	fmt.Println("\n6. Creating 10 posts for existing user (Recycle)...")
	existingUser := regularUsers[0]
	recycledPosts := factory.Recycle(postFactory, existingUser, func(p *Post, u *User) {
		p.AuthorID = u.ID
	}).Count(10).MustCreate(ctx)
	fmt.Printf("   ✅ Created %d posts for user '%s'\n", len(recycledPosts), existingUser.Name)

	// 7. For() - Create posts, each with different author
	fmt.Println("\n7. Creating 5 posts, each with different author (For)...")
	postsWithAuthors := factory.For(postFactory, userFactory, func(p *Post, u *User) {
		p.AuthorID = u.ID
	}).Count(5).MustCreate(ctx)
	fmt.Printf("   ✅ Created %d posts with %d different authors (each author got a role too!)\n", len(postsWithAuthors), len(postsWithAuthors))

	// 8. Nested relationships - posts with comments
	fmt.Println("\n8. Adding comments to posts...")
	commentCount := 0
	for _, post := range posts[:3] { // Just first 3 posts
		comments := factory.Recycle(commentFactory, post, func(c *Comment, p *Post) {
			c.PostID = p.ID
			// Random author from our users
			c.AuthorID = regularUsers[rand.Intn(len(regularUsers))].ID
		}).Count(rand.Intn(5) + 1).MustCreate(ctx)
		commentCount += len(comments)
		fmt.Printf("   - Post '%s': %d comments\n", post.Title, len(comments))
	}
	fmt.Printf("   ✅ Created %d total comments\n", commentCount)

	// 9. Using Tap() for logging
	fmt.Println("\n9. Using Tap() for logging during creation...")
	creationCount := 0
	userFactory.Tap(func(u User) {
		creationCount++
		fmt.Printf("   [Tap] Creating user #%d: %s\n", creationCount, u.Name)
	}).Count(3).MustCreate(ctx)

	// Final stats
	fmt.Println("\n" + repeat("═", 50))
	fmt.Println("📊 FINAL DATABASE STATISTICS")
	fmt.Println(repeat("═", 50))
	fmt.Printf("Users:      %d\n", len(db.users))
	fmt.Printf("Posts:      %d\n", len(db.posts))
	fmt.Printf("Comments:   %d\n", len(db.comments))
	fmt.Printf("Roles:      %d\n", len(db.roles))
	fmt.Printf("UserRoles:  %d\n", len(db.userRoles))
	fmt.Println(repeat("═", 50))

	fmt.Println("\n✅ Database seeding example complete!")
	fmt.Println("\nKey Features Demonstrated:")
	fmt.Println("  ✅ AfterCreate() - Automatic role assignment hook")
	fmt.Println("  ✅ Has() - One-to-many relationships")
	fmt.Println("  ✅ For() - Belongs-to with unique parents")
	fmt.Println("  ✅ Recycle() - Belongs-to with shared parent")
	fmt.Println("  ✅ Many-to-many with pivot tables")
	fmt.Println("  ✅ Tap() - Logging/debugging")
	fmt.Println("  ✅ Count() - Fluent API")
	fmt.Println("  ✅ MustCreate() - Clean error handling")
	fmt.Println("  ✅ Per-call traits for customization")
	fmt.Println("\n💡 Notice: Every user automatically gets a role via AfterCreate hook!")
	fmt.Println("   UserRoles count = Users count (perfect 1:1 assignment!)")
}

// Helper to repeat strings
func repeat(s string, n int) string {
	result := ""
	for i := 0; i < n; i++ {
		result += s
	}
	return result
}
