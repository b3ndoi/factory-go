package main

import (
	"fmt"

	"github.com/b3ndoi/factory-go/factory"
)

// User model
type User struct {
	ID    string
	Name  string
	Email string
	Role  string
}

func main() {
	fmt.Println("=== Factory-Go Basic Example ===\n")

	// Create a simple user factory
	userFactory := factory.New(func(seq int64) User {
		return User{
			Name:  fmt.Sprintf("User %d", seq),
			Email: fmt.Sprintf("user%d@example.com", seq),
			Role:  "user",
		}
	})

	// 1. Create a single user
	fmt.Println("1. Single user:")
	user := userFactory.Make()
	fmt.Printf("   %+v\n\n", user)

	// 2. Create multiple users
	fmt.Println("2. Multiple users:")
	users := userFactory.MakeMany(3)
	for i, u := range users {
		fmt.Printf("   [%d] %s (%s)\n", i+1, u.Name, u.Email)
	}

	// 3. Create with traits (customization)
	fmt.Println("\n3. User with custom role:")
	admin := userFactory.Make(func(u *User) {
		u.Role = "admin"
	})
	fmt.Printf("   %+v\n", admin)

	// 4. Using Count() fluent API
	fmt.Println("\n4. Fluent API:")
	moreUsers := userFactory.Count(5).Make()
	fmt.Printf("   Created %d users using Count(5).Make()\n", len(moreUsers))

	// 5. Named states
	fmt.Println("\n5. Named states:")
	userFactory = userFactory.
		DefineState("admin", func(u *User) {
			u.Role = "admin"
		}).
		DefineState("moderator", func(u *User) {
			u.Role = "moderator"
		})

	admin = userFactory.State("admin").Make()
	moderator := userFactory.State("moderator").Make()
	fmt.Printf("   Admin: %s (Role: %s)\n", admin.Name, admin.Role)
	fmt.Printf("   Moderator: %s (Role: %s)\n", moderator.Name, moderator.Role)

	// 6. Sequences
	fmt.Println("\n6. Sequences (alternating roles):")
	userFactory = userFactory.Sequence(
		func(u *User) { u.Role = "free" },
		func(u *User) { u.Role = "pro" },
		func(u *User) { u.Role = "enterprise" },
	)

	sequencedUsers := userFactory.Count(9).Make()
	for i, u := range sequencedUsers {
		fmt.Printf("   [%d] %s - Plan: %s\n", i+1, u.Name, u.Role)
	}

	// 7. Reset sequence for predictable data
	fmt.Println("\n7. Reset sequence:")
	userFactory.ResetSequence()
	resetUser := userFactory.Make()
	fmt.Printf("   After reset: %s (back to User 1)\n", resetUser.Name)

	fmt.Println("\n✅ Basic example complete!")
}
