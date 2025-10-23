package factory

import (
	"context"
	"fmt"
	"testing"
)

// Benchmark models
type BenchUser struct {
	ID    string
	Name  string
	Email string
	Role  string
}

type BenchPost struct {
	ID       string
	Title    string
	AuthorID string
}

// Benchmarks for core operations

func BenchmarkMake(b *testing.B) {
	f := New(func(seq int64) BenchUser {
		return BenchUser{
			Name:  fmt.Sprintf("User %d", seq),
			Email: fmt.Sprintf("user%d@example.com", seq),
			Role:  "user",
		}
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = f.Make()
	}
}

func BenchmarkMakeWithDefaults(b *testing.B) {
	f := New(func(seq int64) BenchUser {
		return BenchUser{}
	}).WithDefaults(func(u *BenchUser) {
		u.Name = "Default Name"
		u.Email = "default@example.com"
		u.Role = "user"
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = f.Make()
	}
}

func BenchmarkMakeWithTraits(b *testing.B) {
	f := New(func(seq int64) BenchUser {
		return BenchUser{
			Name:  fmt.Sprintf("User %d", seq),
			Email: fmt.Sprintf("user%d@example.com", seq),
		}
	}).WithTraits(func(u *BenchUser) {
		u.Role = "admin"
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = f.Make()
	}
}

func BenchmarkMakeWithState(b *testing.B) {
	f := New(func(seq int64) BenchUser {
		return BenchUser{
			Name:  fmt.Sprintf("User %d", seq),
			Email: fmt.Sprintf("user%d@example.com", seq),
		}
	}).DefineState("admin", func(u *BenchUser) {
		u.Role = "admin"
	})

	adminFactory := f.State("admin")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = adminFactory.Make()
	}
}

func BenchmarkMakeWithSequence(b *testing.B) {
	f := New(func(seq int64) BenchUser {
		return BenchUser{
			Name:  fmt.Sprintf("User %d", seq),
			Email: fmt.Sprintf("user%d@example.com", seq),
		}
	}).Sequence(
		func(u *BenchUser) { u.Role = "admin" },
		func(u *BenchUser) { u.Role = "user" },
		func(u *BenchUser) { u.Role = "moderator" },
	)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = f.Make()
	}
}

func BenchmarkMakeMany(b *testing.B) {
	f := New(func(seq int64) BenchUser {
		return BenchUser{
			Name:  fmt.Sprintf("User %d", seq),
			Email: fmt.Sprintf("user%d@example.com", seq),
		}
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = f.MakeMany(10)
	}
}

func BenchmarkMakeMany_100(b *testing.B) {
	f := New(func(seq int64) BenchUser {
		return BenchUser{
			Name:  fmt.Sprintf("User %d", seq),
			Email: fmt.Sprintf("user%d@example.com", seq),
		}
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = f.MakeMany(100)
	}
}

func BenchmarkCountFluentAPI(b *testing.B) {
	f := New(func(seq int64) BenchUser {
		return BenchUser{
			Name:  fmt.Sprintf("User %d", seq),
			Email: fmt.Sprintf("user%d@example.com", seq),
		}
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = f.Count(10).Make()
	}
}

func BenchmarkRaw(b *testing.B) {
	f := New(func(seq int64) BenchUser {
		return BenchUser{
			Name:  fmt.Sprintf("User %d", seq),
			Email: fmt.Sprintf("user%d@example.com", seq),
		}
	}).WithRawDefaults(func(u *BenchUser) {
		u.Role = "user"
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = f.Raw()
	}
}

func BenchmarkRawJSON(b *testing.B) {
	f := New(func(seq int64) BenchUser {
		return BenchUser{
			Name:  fmt.Sprintf("User %d", seq),
			Email: fmt.Sprintf("user%d@example.com", seq),
		}
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = f.RawJSON()
	}
}

func BenchmarkCreate(b *testing.B) {
	f := New(func(seq int64) BenchUser {
		return BenchUser{
			Name:  fmt.Sprintf("User %d", seq),
			Email: fmt.Sprintf("user%d@example.com", seq),
		}
	}).WithPersist(func(ctx context.Context, u *BenchUser) (*BenchUser, error) {
		u.ID = fmt.Sprintf("id-%d", len(u.Name))
		return u, nil
	})

	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = f.Create(ctx)
	}
}

func BenchmarkCreateWithHooks(b *testing.B) {
	f := New(func(seq int64) BenchUser {
		return BenchUser{
			Name:  fmt.Sprintf("User %d", seq),
			Email: fmt.Sprintf("user%d@example.com", seq),
		}
	}).BeforeCreate(func(ctx context.Context, u *BenchUser) error {
		return nil // No-op hook
	}).AfterCreate(func(ctx context.Context, u *BenchUser) error {
		return nil // No-op hook
	}).WithPersist(func(ctx context.Context, u *BenchUser) (*BenchUser, error) {
		u.ID = "saved"
		return u, nil
	})

	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = f.Create(ctx)
	}
}

func BenchmarkCreateMany(b *testing.B) {
	f := New(func(seq int64) BenchUser {
		return BenchUser{
			Name:  fmt.Sprintf("User %d", seq),
			Email: fmt.Sprintf("user%d@example.com", seq),
		}
	}).WithPersist(func(ctx context.Context, u *BenchUser) (*BenchUser, error) {
		u.ID = "saved"
		return u, nil
	})

	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = f.CreateMany(ctx, 10)
	}
}

func BenchmarkClone(b *testing.B) {
	f := New(func(seq int64) BenchUser {
		return BenchUser{
			Name:  fmt.Sprintf("User %d", seq),
			Email: fmt.Sprintf("user%d@example.com", seq),
		}
	}).WithDefaults(func(u *BenchUser) {
		u.Role = "user"
	}).DefineState("admin", func(u *BenchUser) {
		u.Role = "admin"
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = f.Clone()
	}
}

// Benchmarks for relationships

func BenchmarkFor(b *testing.B) {
	userFactory := New(func(seq int64) BenchUser {
		return BenchUser{
			ID:   fmt.Sprintf("user-%d", seq),
			Name: fmt.Sprintf("User %d", seq),
		}
	})

	postFactory := New(func(seq int64) BenchPost {
		return BenchPost{
			Title: fmt.Sprintf("Post %d", seq),
		}
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = For(postFactory, userFactory, func(p *BenchPost, u *BenchUser) {
			p.AuthorID = u.ID
		}).Make()
	}
}

func BenchmarkHas(b *testing.B) {
	userFactory := New(func(seq int64) BenchUser {
		return BenchUser{
			ID:   fmt.Sprintf("user-%d", seq),
			Name: fmt.Sprintf("User %d", seq),
		}
	})

	postFactory := New(func(seq int64) BenchPost {
		return BenchPost{
			Title: fmt.Sprintf("Post %d", seq),
		}
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = Has(userFactory, postFactory, 5, func(u *BenchUser, p *BenchPost) {
			p.AuthorID = u.ID
		}).Make()
	}
}

// Comparison: Factory vs Manual Helpers

func BenchmarkManualHelper(b *testing.B) {
	// Simulate manual test helper function
	createUser := func(seq int) BenchUser {
		return BenchUser{
			Name:  fmt.Sprintf("User %d", seq),
			Email: fmt.Sprintf("user%d@example.com", seq),
			Role:  "user",
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = createUser(i)
	}
}

func BenchmarkManualHelperMany(b *testing.B) {
	createUsers := func(count int) []BenchUser {
		users := make([]BenchUser, count)
		for i := 0; i < count; i++ {
			users[i] = BenchUser{
				Name:  fmt.Sprintf("User %d", i),
				Email: fmt.Sprintf("user%d@example.com", i),
				Role:  "user",
			}
		}
		return users
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = createUsers(10)
	}
}

// Parallel execution benchmarks

func BenchmarkMakeParallel(b *testing.B) {
	f := New(func(seq int64) BenchUser {
		return BenchUser{
			Name:  fmt.Sprintf("User %d", seq),
			Email: fmt.Sprintf("user%d@example.com", seq),
		}
	})

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = f.Make()
		}
	})
}

func BenchmarkClonePerGoroutine(b *testing.B) {
	baseFactory := New(func(seq int64) BenchUser {
		return BenchUser{
			Name:  fmt.Sprintf("User %d", seq),
			Email: fmt.Sprintf("user%d@example.com", seq),
		}
	})

	b.RunParallel(func(pb *testing.PB) {
		f := baseFactory.Clone() // Each goroutine gets its own clone
		for pb.Next() {
			_ = f.Make()
		}
	})
}

