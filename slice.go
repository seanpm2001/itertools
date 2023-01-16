package itertools

import (
	optionext "github.com/go-playground/pkg/v5/values/option"
	"sort"
)

// WrapSlice accepts and turns a sliceWrapper into an iterator.
//
// The default the Map type to struct{} when none is required. See WrapSliceMap if one is needed.
func WrapSlice[T any](slice []T) sliceWrapper[T, struct{}] {
	return WrapSliceMap[T, struct{}](slice)
}

// WrapSliceMap accepts and turns a sliceWrapper into an iterator with a map type specified for IterPar() to allow the Map helper
// function.
func WrapSliceMap[T, MAP any](slice []T) sliceWrapper[T, MAP] {
	return sliceWrapper[T, MAP]{
		slice: slice,
	}
}

type sliceWrapper[T, MAP any] struct {
	slice []T
}

func (i *sliceWrapper[T, MAP]) Next() optionext.Option[T] {
	if len(i.slice) == 0 {
		return optionext.None[T]()
	}
	v := i.slice[0]
	i.slice = i.slice[1:]
	return optionext.Some(v)
}

// IntoIter turns the slice wrapper into an `Iterator[T]`
func (i sliceWrapper[T, MAP]) IntoIter() *sliceWrapper[T, MAP] {
	return &i
}

// Iter is a convenience function that converts the sliceWrapper iterator into an `*Iterate[T]`.
func (i sliceWrapper[T, MAP]) Iter() Iterate[T, *sliceWrapper[T, MAP], MAP] {
	return IterMap[T, *sliceWrapper[T, MAP], MAP](i.IntoIter())
}

// Slice returns the underlying sliceWrapper wrapped by the *sliceWrapper.
func (i sliceWrapper[T, MAP]) Slice() []T {
	return i.slice
}

// Len returns the length of the underlying sliceWrapper.
func (i sliceWrapper[T, MAP]) Len() int {
	return len(i.slice)
}

// Cap returns the capacity of the underlying sliceWrapper.
func (i sliceWrapper[T, MAP]) Cap() int {
	return cap(i.slice)
}

// Sort sorts the sliceWrapper x given the provided less function.
//
// The sort is not guaranteed to be stable: equal elements
// may be reversed from their original order.
// For a stable sort, use SortStable.
//
// `T` must be comparable.
func (i sliceWrapper[T, MAP]) Sort(less func(i T, j T) bool) sliceWrapper[T, MAP] {
	sort.Slice(i.slice, func(j, k int) bool {
		return less(i.slice[j], i.slice[k])
	})
	return WrapSliceMap[T, MAP](i.slice)
}

// SortStable sorts the sliceWrapper x using the provided less
// function, keeping equal elements in their original order.
func (i sliceWrapper[T, MAP]) SortStable(less func(i T, j T) bool) sliceWrapper[T, MAP] {
	sort.SliceStable(i.slice, func(j, k int) bool {
		return less(i.slice[j], i.slice[k])
	})
	return WrapSliceMap[T, MAP](i.slice)
}

// Retain retains only the elements specified by the function.
//
// This shuffles and returns the retained values of the slice.
func (i sliceWrapper[T, MAP]) Retain(fn func(v T) bool) sliceWrapper[T, MAP] {
	return WrapSliceMap[T, MAP](RetainSlice(i.slice, fn))
}

// RetainSlice retains only the elements specified by the function.
//
// This shuffles and returns the retained values of the slice.
func RetainSlice[T any](slice []T, fn func(v T) bool) []T {
	var j int
	for _, v := range slice {
		if fn(v) {
			slice[j] = v
			j++
		}
	}
	return slice[:j]
}
