package itertools

import (
	optionext "github.com/go-playground/pkg/v5/values/option"
)

// Entry represents a single Map entry.
type Entry[K comparable, V any] struct {
	Key   K
	Value V
}

// WrapMap creates a new iterator for transformation of types.
func WrapMap[K comparable, V any](m map[K]V) *mapWrapper[K, V, struct{}] {
	return WrapMapMap[K, V, struct{}](m)
}

// WrapMapMap creates a new `mapWrapper` for use which also specifies a potential future `Map` operation.
func WrapMapMap[K comparable, V, MAP any](m map[K]V) *mapWrapper[K, V, MAP] {
	return &mapWrapper[K, V, MAP]{
		mp: m,
	}
}

// mapWrapper is used to transform elements from one type to another.
type mapWrapper[K comparable, V, MAP any] struct {
	mp map[K]V
}

// Next returns the next transformed element or None if at the end of the iterator.
//
// Warning: This consumes(removes) the map entries as it iterates.
func (i *mapWrapper[K, V, MAP]) Next() optionext.Option[Entry[K, V]] {
	for k, v := range i.mp {
		delete(i.mp, k)
		return optionext.Some(Entry[K, V]{
			Key:   k,
			Value: v,
		})
	}
	return optionext.None[Entry[K, V]]()
}

// Iter is a convenience function that converts the map iterator into an `*Iterate[T]`.
func (i *mapWrapper[K, V, MAP]) Iter() *Iterate[Entry[K, V], MAP] {
	return IterMap[Entry[K, V], MAP](i)
}

// Retain retains only the elements specified by the function and removes others.
func (i *mapWrapper[K, V, MAP]) Retain(fn func(k K, v V) bool) *mapWrapper[K, V, MAP] {
	for k, v := range i.mp {
		if fn(k, v) {
			continue
		}
		delete(i.mp, k)
	}
	return i
}

// Len returns the underlying map's length.
func (i *mapWrapper[K, V, MAP]) Len() int {
	return len(i.mp)
}
