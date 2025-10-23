package main

import (
	"context"
	"fmt"

	"github.com/b3ndoi/factory-go/factory"
)

// User with realistic fields
type User struct {
	ID        string
	FirstName string
	LastName  string
	Email     string
	Phone     string
	Address   string
	City      string
	Country   string
	Bio       string
	Password  string `json:"password,omitempty"` // Only for API/Raw
}

// Simple faker simulation (in real app, use github.com/brianvoe/gofakeit/v6)
type SimpleFaker struct {
	seq int
}

func (f *SimpleFaker) FirstName() string {
	names := []string{"John", "Jane", "Alice", "Bob", "Charlie", "Diana"}
	f.seq++
	return names[f.seq%len(names)]
}

func (f *SimpleFaker) LastName() string {
	names := []string{"Smith", "Johnson", "Williams", "Brown", "Jones", "Garcia"}
	f.seq++
	return names[f.seq%len(names)]
}

func (f *SimpleFaker) Email() string {
	f.seq++
	return fmt.Sprintf("user%d@faker.com", f.seq)
}

func (f *SimpleFaker) Phone() string {
	f.seq++
	return fmt.Sprintf("+1-555-%04d", f.seq)
}

func (f *SimpleFaker) City() string {
	cities := []string{"New York", "Los Angeles", "Chicago", "Houston", "Phoenix"}
	f.seq++
	return cities[f.seq%len(cities)]
}

func (f *SimpleFaker) Country() string {
	return "USA"
}

func (f *SimpleFaker) Address() string {
	f.seq++
	return fmt.Sprintf("%d Main Street", f.seq*100)
}

func (f *SimpleFaker) Paragraph() string {
	return "This is a sample bio paragraph with interesting information about the user."
}

// Mock DB
type DB struct {
	users []*User
}

func (db *DB) CreateUser(ctx context.Context, u *User) (*User, error) {
	u.ID = fmt.Sprintf("user-%d", len(db.users)+1)
	db.users = append(db.users, u)
	return u, nil
}

func main() {
	fmt.Println("=== Factory-Go with Faker Integration ===")
	fmt.Println()

	db := &DB{}
	ctx := context.Background()
	faker := &SimpleFaker{}

	// Create factory with faker-generated realistic data
	userFactory := factory.New(func(seq int64) User {
		return User{} // Start with empty struct
	}).WithDefaults(func(u *User) {
		// Use faker for ALL fields - realistic test data!
		u.FirstName = faker.FirstName()
		u.LastName = faker.LastName()
		u.Email = faker.Email()
		u.Phone = faker.Phone()
		u.Address = faker.Address()
		u.City = faker.City()
		u.Country = faker.Country()
		u.Bio = faker.Paragraph()
	}).WithRawDefaults(func(u *User) {
		// Password only for API/Raw, not for domain models
		u.Password = "FakerPassword123!"
	}).WithPersist(db.CreateUser)

	// 1. Create users with realistic data
	fmt.Println("1. Creating users with faker data:")
	users := userFactory.Count(5).MustCreate(ctx)
	for i, u := range users {
		fmt.Printf("   [%d] %s %s\n", i+1, u.FirstName, u.LastName)
		fmt.Printf("       Email: %s\n", u.Email)
		fmt.Printf("       Phone: %s\n", u.Phone)
		fmt.Printf("       Location: %s, %s\n", u.City, u.Country)
		fmt.Printf("       Bio: %s\n\n", u.Bio[:50]+"...")
	}

	// 2. Override faker data with specific values
	fmt.Println("2. Overriding faker data for specific test:")
	specificUser := userFactory.Make(func(u *User) {
		u.FirstName = "John"
		u.LastName = "Doe"
		u.Email = "john.doe@custom.com"
		// Other fields still use faker!
	})
	fmt.Printf("   Custom: %s %s (%s)\n", specificUser.FirstName, specificUser.LastName, specificUser.Email)
	fmt.Printf("   Faker: Phone=%s, City=%s\n", specificUser.Phone, specificUser.City)

	// 3. Using WithRawDefaults for API testing
	fmt.Println("\n3. API testing with faker + rawDefaults:")
	apiUser := userFactory.Raw() // Includes faker data + password
	fmt.Printf("   Name: %s %s\n", apiUser.FirstName, apiUser.LastName)
	fmt.Printf("   Email: %s\n", apiUser.Email)
	fmt.Printf("   Password: %s (only in Raw, not in Make/Create!)\n", apiUser.Password)

	domainUser := userFactory.Make() // Faker data but NO password
	fmt.Printf("   Domain user password: '%s' (empty ✅)\n", domainUser.Password)

	// 4. Combining faker with sequences
	fmt.Println("\n4. Faker + Sequences:")
	userFactory.ResetSequence()
	tieredUsers := userFactory.Sequence(
		func(u *User) { u.Bio = "Free tier user" },
		func(u *User) { u.Bio = "Pro tier user" },
		func(u *User) { u.Bio = "Enterprise tier user" },
	).Count(9).Make()

	for i, u := range tieredUsers {
		fmt.Printf("   [%d] %s - %s\n", i+1, u.FirstName+" "+u.LastName, u.Bio)
	}

	// 5. Faker with states
	fmt.Println("\n5. Faker + Named States:")
	adminFactory := userFactory.Clone().DefineState("superadmin", func(u *User) {
		u.Email = "admin@company.com" // Override faker email for admins
		u.Bio = "System Administrator"
	})

	superadmin := adminFactory.State("superadmin").Make()
	fmt.Printf("   Superadmin: %s %s\n", superadmin.FirstName, superadmin.LastName)
	fmt.Printf("   Email: %s (custom, not faker)\n", superadmin.Email)
	fmt.Printf("   Bio: %s (custom)\n", superadmin.Bio)

	fmt.Println("\n" + repeat("═", 60))
	fmt.Println("✅ Faker Integration Example Complete!")
	fmt.Println(repeat("═", 60))

	fmt.Println("\n🎯 Key Takeaways:")
	fmt.Println("  ✅ WithDefaults() perfect for faker integration")
	fmt.Println("  ✅ Faker generates realistic test data")
	fmt.Println("  ✅ Can still override faker values with traits")
	fmt.Println("  ✅ WithRawDefaults() adds API-specific fields")
	fmt.Println("  ✅ Faker + Sequences = varied realistic data")
	fmt.Println("  ✅ Faker + States = controlled variations")
	fmt.Println("\n💡 In production, use: github.com/brianvoe/gofakeit/v6")
}

func repeat(s string, n int) string {
	result := ""
	for i := 0; i < n; i++ {
		result += s
	}
	return result
}
