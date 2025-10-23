package factory

import (
	"context"
	"encoding/json"
	"sync/atomic"
)

// Trait mutates a T before persistence (like Laravel "state").
type Trait[T any] func(*T)

// BeforeCreate runs before persistence (e.g., validation, setup).
type BeforeCreate[T any] func(ctx context.Context, t *T) error

// AfterCreate runs after persistence (e.g., create related rows).
type AfterCreate[T any] func(ctx context.Context, t *T) error

// PersistFn saves *T (user provides DB logic) and returns possibly updated *T.
type PersistFn[T any] func(ctx context.Context, t *T) (*T, error)

// Factory builds Ts with defaults, traits, and optional persistence.
type Factory[T any] struct {
	makeFn      func(seq int64) T
	defaults    []Trait[T]          // Applied first (for faker/defaults)
	rawDefaults []Trait[T]          // Applied only for Raw/RawJSON methods
	traits      []Trait[T]          // Applied second (global traits)
	sequences   []Trait[T]          // Cycled through for each item
	states      map[string]Trait[T] // Named states (like Laravel state methods)
	persist     PersistFn[T]
	before      []BeforeCreate[T] // Hooks before persistence
	after       []AfterCreate[T]  // Hooks after persistence
	tapFn       func(T)           // Tap function for debugging
	seq         int64
	count       int // Count for fluent API (0 means not set)
}

// CountedFactory is a fluent wrapper that knows how many items to create.
type CountedFactory[T any] struct {
	factory *Factory[T]
	count   int
}

// New constructs a factory with a default make function (receives a sequence number).
func New[T any](makeFn func(seq int64) T) *Factory[T] {
	return &Factory[T]{
		makeFn: makeFn,
		states: make(map[string]Trait[T]),
	}
}

// WithDefaults sets default traits applied first (ideal for faker/default values).
// These are applied before WithTraits and per-call traits.
func (f *Factory[T]) WithDefaults(ts ...Trait[T]) *Factory[T] {
	f.defaults = append(f.defaults, ts...)
	return f
}

// WithRawDefaults sets traits applied ONLY when using Raw/RawJSON methods.
// Useful for adding fields needed for API testing but not for persistence.
// Example: Add validation fields, computed fields, or API-specific attributes.
func (f *Factory[T]) WithRawDefaults(ts ...Trait[T]) *Factory[T] {
	f.rawDefaults = append(f.rawDefaults, ts...)
	return f
}

// WithTraits appends global traits applied to every Make/Create call.
func (f *Factory[T]) WithTraits(ts ...Trait[T]) *Factory[T] {
	f.traits = append(f.traits, ts...)
	return f
}

// Sequence sets traits that cycle through for each created item (like Laravel's sequence()).
// Example: Sequence(trait1, trait2) will alternate: trait1, trait2, trait1, trait2...
func (f *Factory[T]) Sequence(ts ...Trait[T]) *Factory[T] {
	f.sequences = ts
	return f
}

// DefineState registers a named state that can be applied later (like Laravel state methods).
// Example: factory.DefineState("admin", func(u *User) { u.Role = "admin" })
func (f *Factory[T]) DefineState(name string, trait Trait[T]) *Factory[T] {
	f.states[name] = trait
	return f
}

// State applies a previously defined named state by adding it as a trait.
// Returns a new factory instance with the state applied.
// Example: factory.State("admin").Make()
func (f *Factory[T]) State(name string) *Factory[T] {
	trait, ok := f.states[name]
	if !ok {
		panic("factory: unknown state '" + name + "'")
	}
	// Create a shallow copy with the state trait added
	copy := *f
	copy.traits = append([]Trait[T]{}, f.traits...)
	copy.traits = append(copy.traits, trait)
	return &copy
}

// WithPersist sets how to save T (optional; required for Create()).
func (f *Factory[T]) WithPersist(p PersistFn[T]) *Factory[T] {
	f.persist = p
	return f
}

// BeforeCreate adds hooks executed before persistence.
func (f *Factory[T]) BeforeCreate(h BeforeCreate[T]) *Factory[T] {
	f.before = append(f.before, h)
	return f
}

// AfterCreate adds hooks executed after persistence.
func (f *Factory[T]) AfterCreate(h AfterCreate[T]) *Factory[T] {
	f.after = append(f.after, h)
	return f
}

// Tap sets a function to be called with each created item (useful for debugging/logging).
func (f *Factory[T]) Tap(fn func(T)) *Factory[T] {
	f.tapFn = fn
	return f
}

// When applies traits only if the condition is true.
func (f *Factory[T]) When(condition bool, ts ...Trait[T]) *Factory[T] {
	if condition {
		f.traits = append(f.traits, ts...)
	}
	return f
}

// Unless applies traits only if the condition is false.
func (f *Factory[T]) Unless(condition bool, ts ...Trait[T]) *Factory[T] {
	if !condition {
		f.traits = append(f.traits, ts...)
	}
	return f
}

// Clone creates a deep copy of the factory for creating variations.
func (f *Factory[T]) Clone() *Factory[T] {
	clone := &Factory[T]{
		makeFn:      f.makeFn,
		defaults:    append([]Trait[T]{}, f.defaults...),
		rawDefaults: append([]Trait[T]{}, f.rawDefaults...),
		traits:      append([]Trait[T]{}, f.traits...),
		sequences:   append([]Trait[T]{}, f.sequences...),
		states:      make(map[string]Trait[T]),
		persist:     f.persist,
		before:      append([]BeforeCreate[T]{}, f.before...),
		after:       append([]AfterCreate[T]{}, f.after...),
		tapFn:       f.tapFn,
		seq:         0, // Reset sequence for clone
		count:       f.count,
	}
	// Deep copy states map
	for k, v := range f.states {
		clone.states[k] = v
	}
	return clone
}

func (f *Factory[T]) nextSeq() int64 {
	return atomic.AddInt64(&f.seq, 1)
}

// ResetSequence resets the sequence counter to 0.
// Useful for test isolation to get predictable sequence numbers.
func (f *Factory[T]) ResetSequence() *Factory[T] {
	atomic.StoreInt64(&f.seq, 0)
	return f
}

// Count sets the number of items to create (fluent API like Laravel).
// Returns a CountedFactory that has Make() and Create() methods for multiple items.
// Example: factory.Count(10).Make() or factory.Count(5).State("admin").Create(ctx)
func (f *Factory[T]) Count(n int) *CountedFactory[T] {
	return &CountedFactory[T]{
		factory: f,
		count:   n,
	}
}

// Times is an alias for Count (more semantic in some contexts).
func (f *Factory[T]) Times(n int) *CountedFactory[T] {
	return f.Count(n)
}

// Make builds but does not persist (like Laravel's make()).
// Applies traits in order: defaults → global traits → sequence → per-call traits.
func (f *Factory[T]) Make(ts ...Trait[T]) T {
	seq := f.nextSeq()
	t := f.makeFn(seq)

	// Apply defaults first (faker/default values)
	for _, tr := range f.defaults {
		tr(&t)
	}
	// Then global traits
	for _, tr := range f.traits {
		tr(&t)
	}
	// Then sequence trait (cycles through)
	if len(f.sequences) > 0 {
		idx := int((seq - 1) % int64(len(f.sequences)))
		f.sequences[idx](&t)
	}
	// Finally per-call traits
	for _, tr := range ts {
		tr(&t)
	}
	// Call tap function if set
	if f.tapFn != nil {
		f.tapFn(t)
	}
	return t
}

// Raw builds but does not persist, with rawDefaults applied (like Laravel's raw()).
// Applies: defaults → rawDefaults → global traits → sequence → per-call traits.
// Useful for getting attribute values for testing validation or API requests.
func (f *Factory[T]) Raw(ts ...Trait[T]) T {
	seq := f.nextSeq()
	t := f.makeFn(seq)

	// Apply defaults first (faker/default values)
	for _, tr := range f.defaults {
		tr(&t)
	}
	// Then raw-specific defaults
	for _, tr := range f.rawDefaults {
		tr(&t)
	}
	// Then global traits
	for _, tr := range f.traits {
		tr(&t)
	}
	// Then sequence trait (cycles through)
	if len(f.sequences) > 0 {
		idx := int((seq - 1) % int64(len(f.sequences)))
		f.sequences[idx](&t)
	}
	// Finally per-call traits
	for _, tr := range ts {
		tr(&t)
	}
	// Call tap function if set
	if f.tapFn != nil {
		f.tapFn(t)
	}
	return t
}

// RawMany builds count items without persisting, with rawDefaults applied.
func (f *Factory[T]) RawMany(count int, ts ...Trait[T]) []T {
	items := make([]T, count)
	for i := 0; i < count; i++ {
		items[i] = f.Raw(ts...)
	}
	return items
}

// RawJSON builds and returns JSON representation (like Laravel's raw()).
// Useful for testing API endpoints without persistence.
func (f *Factory[T]) RawJSON(ts ...Trait[T]) ([]byte, error) {
	obj := f.Raw(ts...)
	return json.Marshal(obj)
}

// RawManyJSON builds count items and returns JSON array.
func (f *Factory[T]) RawManyJSON(count int, ts ...Trait[T]) ([]byte, error) {
	items := f.RawMany(count, ts...)
	return json.Marshal(items)
}

// Create builds, persists, runs hooks, and returns *T (like Laravel's create()).
func (f *Factory[T]) Create(ctx context.Context, ts ...Trait[T]) (*T, error) {
	if f.persist == nil {
		panic("factory: Create called without persist function; use WithPersist")
	}
	obj := f.Make(ts...)

	// Run before hooks
	for _, h := range f.before {
		if err := h(ctx, &obj); err != nil {
			return nil, err
		}
	}

	// Persist
	out, err := f.persist(ctx, &obj)
	if err != nil {
		return nil, err
	}

	// Run after hooks
	for _, h := range f.after {
		if err := h(ctx, out); err != nil {
			return nil, err
		}
	}
	return out, nil
}

// MakeMany builds count items without persisting (like Laravel's count()->make()).
func (f *Factory[T]) MakeMany(count int, ts ...Trait[T]) []T {
	items := make([]T, count)
	for i := 0; i < count; i++ {
		items[i] = f.Make(ts...)
	}
	return items
}

// CreateMany builds, persists, and runs hooks for count items (like Laravel's count()->create()).
func (f *Factory[T]) CreateMany(ctx context.Context, count int, ts ...Trait[T]) ([]*T, error) {
	if f.persist == nil {
		panic("factory: CreateMany called without persist function; use WithPersist")
	}
	items := make([]*T, 0, count)
	for i := 0; i < count; i++ {
		item, err := f.Create(ctx, ts...)
		if err != nil {
			return items, err
		}
		items = append(items, item)
	}
	return items, nil
}

// Must* Variants (panic on error instead of returning error)

// MustCreate builds, persists, and returns *T. Panics on error (useful in tests).
func (f *Factory[T]) MustCreate(ctx context.Context, ts ...Trait[T]) *T {
	item, err := f.Create(ctx, ts...)
	if err != nil {
		panic("factory: MustCreate failed: " + err.Error())
	}
	return item
}

// MustCreateMany builds, persists, and returns []*T. Panics on error (useful in tests).
func (f *Factory[T]) MustCreateMany(ctx context.Context, count int, ts ...Trait[T]) []*T {
	items, err := f.CreateMany(ctx, count, ts...)
	if err != nil {
		panic("factory: MustCreateMany failed: " + err.Error())
	}
	return items
}

// MustRawJSON builds and returns JSON. Panics on error (useful in tests).
func (f *Factory[T]) MustRawJSON(ts ...Trait[T]) []byte {
	data, err := f.RawJSON(ts...)
	if err != nil {
		panic("factory: MustRawJSON failed: " + err.Error())
	}
	return data
}

// MustRawManyJSON builds multiple items and returns JSON array. Panics on error (useful in tests).
func (f *Factory[T]) MustRawManyJSON(count int, ts ...Trait[T]) []byte {
	data, err := f.RawManyJSON(count, ts...)
	if err != nil {
		panic("factory: MustRawManyJSON failed: " + err.Error())
	}
	return data
}

// CountedFactory Methods

// Make builds count items without persisting.
func (cf *CountedFactory[T]) Make(ts ...Trait[T]) []T {
	return cf.factory.MakeMany(cf.count, ts...)
}

// Create builds, persists, and runs hooks for count items.
func (cf *CountedFactory[T]) Create(ctx context.Context, ts ...Trait[T]) ([]*T, error) {
	return cf.factory.CreateMany(ctx, cf.count, ts...)
}

// Raw builds count items without persisting, with rawDefaults applied.
func (cf *CountedFactory[T]) Raw(ts ...Trait[T]) []T {
	return cf.factory.RawMany(cf.count, ts...)
}

// RawJSON builds count items and returns JSON array.
func (cf *CountedFactory[T]) RawJSON(ts ...Trait[T]) ([]byte, error) {
	return cf.factory.RawManyJSON(cf.count, ts...)
}

// State applies a named state to the underlying factory and returns a new CountedFactory.
func (cf *CountedFactory[T]) State(name string) *CountedFactory[T] {
	return &CountedFactory[T]{
		factory: cf.factory.State(name),
		count:   cf.count,
	}
}

// MustCreate builds, persists, and returns []*T. Panics on error (useful in tests).
func (cf *CountedFactory[T]) MustCreate(ctx context.Context, ts ...Trait[T]) []*T {
	return cf.factory.MustCreateMany(ctx, cf.count, ts...)
}

// MustRawJSON builds count items and returns JSON array. Panics on error (useful in tests).
func (cf *CountedFactory[T]) MustRawJSON(ts ...Trait[T]) []byte {
	return cf.factory.MustRawManyJSON(cf.count, ts...)
}

// Relationship Helpers

// For sets up a belongs-to relationship by creating a related model first.
// The linkFn receives the current model and the created related model to establish the relationship.
// Example: For(postFactory, userFactory, func(p *Post, u *User) { p.AuthorID = u.ID })
func For[T any, R any](f *Factory[T], relatedFactory *Factory[R], linkFn func(*T, *R)) *Factory[T] {
	// Create a copy of the factory with an added trait
	copy := *f
	copy.traits = append([]Trait[T]{}, f.traits...)

	// Add a trait that will create the related model when Make is called
	// Note: This only works for Make/Raw, not Create (which needs context)
	copy.defaults = append([]Trait[T]{}, f.defaults...)
	copy.defaults = append(copy.defaults, func(t *T) {
		related := relatedFactory.Make()
		linkFn(t, &related)
	})

	return &copy
}

// ForModel sets up a belongs-to relationship using an existing model instance.
// The linkFn receives the current model and the existing related model.
// Example: ForModel(postFactory, user, func(p *Post, u *User) { p.AuthorID = u.ID })
func ForModel[T any, R any](f *Factory[T], related *R, linkFn func(*T, *R)) *Factory[T] {
	// Create a copy with an added trait
	copy := *f
	copy.traits = append([]Trait[T]{}, f.traits...)
	copy.traits = append(copy.traits, func(t *T) {
		linkFn(t, related)
	})

	return &copy
}

// Recycle is an alias for ForModel - reuse the same related model across multiple creations.
// Example: Recycle(postFactory, user, func(p *Post, u *User) { p.AuthorID = u.ID })
func Recycle[T any, R any](f *Factory[T], related *R, linkFn func(*T, *R)) *Factory[T] {
	return ForModel(f, related, linkFn)
}

// Has creates a parent model with child models (inverse of For).
// Creates one parent, then creates 'count' children linked to that parent.
// Returns a factory that when Create() is called, will create parent + children.
// Example: Has(userFactory, postFactory, 3, func(u *User, p *Post) { p.AuthorID = u.ID })
func Has[T any, R any](
	parentFactory *Factory[T],
	childFactory *Factory[R],
	count int,
	linkFn func(parent *T, child *R),
) *HasFactory[T, R] {
	return &HasFactory[T, R]{
		parent: parentFactory,
		child:  childFactory,
		count:  count,
		linkFn: linkFn,
	}
}

// HasAttached creates a parent model with many-to-many relationships through a pivot table.
// Creates one parent, creates 'count' related models, and creates pivot records for each.
// Example: HasAttached(userFactory, roleFactory, pivotFactory, 3, linkFn)
func HasAttached[T any, R any, P any](
	parentFactory *Factory[T],
	relatedFactory *Factory[R],
	pivotFactory *Factory[P],
	count int,
	linkFn func(pivot *P, parent *T, related *R),
) *HasAttachedFactory[T, R, P] {
	return &HasAttachedFactory[T, R, P]{
		parent:       parentFactory,
		related:      relatedFactory,
		pivotFactory: pivotFactory,
		count:        count,
		linkFn:       linkFn,
	}
}

// HasFactory manages has-many relationships.
type HasFactory[T any, R any] struct {
	parent *Factory[T]
	child  *Factory[R]
	count  int
	linkFn func(*T, *R)
}

// HasAttachedFactory manages many-to-many relationships with pivot tables.
type HasAttachedFactory[T any, R any, P any] struct {
	parent       *Factory[T]
	related      *Factory[R]
	pivotFactory *Factory[P]
	count        int
	linkFn       func(*P, *T, *R)
}

// Make creates parent with children (in-memory only).
func (hf *HasFactory[T, R]) Make() (T, []R) {
	parent := hf.parent.Make()
	children := make([]R, hf.count)
	for i := 0; i < hf.count; i++ {
		child := hf.child.Make()
		if hf.linkFn != nil {
			hf.linkFn(&parent, &child)
		}
		children[i] = child
	}
	return parent, children
}

// Create creates and persists parent with children.
// Returns the parent and all created children.
func (hf *HasFactory[T, R]) Create(ctx context.Context) (*T, []*R, error) {
	// Create parent first
	parent, err := hf.parent.Create(ctx)
	if err != nil {
		return nil, nil, err
	}

	// Create children linked to parent
	children := make([]*R, 0, hf.count)
	for i := 0; i < hf.count; i++ {
		var child *R
		var err error

		if hf.linkFn != nil {
			// Create wrapper function that swaps parameter order for Recycle
			child, err = Recycle(hf.child, parent, func(c *R, p *T) {
				hf.linkFn(p, c)
			}).Create(ctx)
		} else {
			// No link function - just create child
			child, err = hf.child.Create(ctx)
		}

		if err != nil {
			return parent, children, err
		}
		children = append(children, child)
	}

	return parent, children, nil
}

// MustCreate creates and persists parent with children, panics on error.
func (hf *HasFactory[T, R]) MustCreate(ctx context.Context) (*T, []*R) {
	parent, children, err := hf.Create(ctx)
	if err != nil {
		panic("factory: HasFactory.MustCreate failed: " + err.Error())
	}
	return parent, children
}

// HasAttachedFactory Methods

// Make creates parent with related models and pivot records (in-memory only).
func (haf *HasAttachedFactory[T, R, P]) Make() (T, []R, []P) {
	parent := haf.parent.Make()
	related := make([]R, haf.count)
	pivots := make([]P, haf.count)

	for i := 0; i < haf.count; i++ {
		rel := haf.related.Make()
		pivot := haf.pivotFactory.Make()
		haf.linkFn(&pivot, &parent, &rel)
		related[i] = rel
		pivots[i] = pivot
	}

	return parent, related, pivots
}

// Create creates and persists parent, related models, and pivot records.
func (haf *HasAttachedFactory[T, R, P]) Create(ctx context.Context) (*T, []*R, []*P, error) {
	// Create parent first
	parent, err := haf.parent.Create(ctx)
	if err != nil {
		return nil, nil, nil, err
	}

	// Create related models and pivot records
	relatedModels := make([]*R, 0, haf.count)
	pivotRecords := make([]*P, 0, haf.count)

	for i := 0; i < haf.count; i++ {
		// Create related model
		related, err := haf.related.Create(ctx)
		if err != nil {
			return parent, relatedModels, pivotRecords, err
		}
		relatedModels = append(relatedModels, related)

		// Create pivot record with link function
		pivot, err := haf.pivotFactory.Create(ctx, func(p *P) {
			haf.linkFn(p, parent, related)
		})
		if err != nil {
			return parent, relatedModels, pivotRecords, err
		}
		pivotRecords = append(pivotRecords, pivot)
	}

	return parent, relatedModels, pivotRecords, nil
}

// MustCreate creates and persists parent, related models, and pivot records. Panics on error.
func (haf *HasAttachedFactory[T, R, P]) MustCreate(ctx context.Context) (*T, []*R, []*P) {
	parent, related, pivots, err := haf.Create(ctx)
	if err != nil {
		panic("factory: HasAttachedFactory.MustCreate failed: " + err.Error())
	}
	return parent, related, pivots
}
