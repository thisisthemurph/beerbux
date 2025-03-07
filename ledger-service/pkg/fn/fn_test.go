package fn

import (
	"fmt"
	"testing"
)

func TestContains(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		testCases := []struct {
			name     string
			slice    []int
			value    int
			expected bool
		}{
			{"value in slice", []int{1, 2, 3, 4, 5}, 3, true},
			{"value not in slice", []int{1, 2, 3, 4, 5}, 6, false},
			{"empty slice", []int{}, 1, false},
			{"single element match", []int{42}, 42, true},
			{"single element no match", []int{42}, 99, false},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				result := Contains(tc.slice, tc.value)
				if result != tc.expected {
					t.Errorf("Contains(%v, %d) = %v, expected %v", tc.slice, tc.value, result, tc.expected)
				}
			})
		}
	})

	t.Run("string", func(t *testing.T) {
		testCases := []struct {
			name     string
			slice    []string
			value    string
			expected bool
		}{
			{"value in slice", []string{"apple", "banana", "cherry"}, "banana", true},
			{"value not in slice", []string{"apple", "banana", "cherry"}, "orange", false},
			{"empty slice", []string{}, "apple", false},
			{"single element match", []string{"hello"}, "hello", true},
			{"single element no match", []string{"hello"}, "world", false},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				result := Contains(tc.slice, tc.value)
				if result != tc.expected {
					t.Errorf("Contains(%v, %q) = %v, expected %v", tc.slice, tc.value, result, tc.expected)
				}
			})
		}
	})

	t.Run("float64", func(t *testing.T) {
		testCases := []struct {
			name     string
			slice    []float64
			value    float64
			expected bool
		}{
			{"value in slice", []float64{1.1, 2.2, 3.3}, 2.2, true},
			{"value not in slice", []float64{1.1, 2.2, 3.3}, 4.4, false},
			{"empty slice", []float64{}, 1.1, false},
			{"single element match", []float64{42.0}, 42.0, true},
			{"single element no match", []float64{42.0}, 99.9, false},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				result := Contains(tc.slice, tc.value)
				if result != tc.expected {
					t.Errorf("Contains(%v, %f) = %v, expected %v", tc.slice, tc.value, result, tc.expected)
				}
			})
		}
	})
}

func TestContainsPointers(t *testing.T) {
	t.Run("int pointers", func(t *testing.T) {
		a, b, c := 1, 2, 3
		slice := []*int{&a, &b, &c}

		testCases := []struct {
			name     string
			slice    []*int
			value    *int
			expected bool
		}{
			{"value in slice", slice, &b, true},
			{"value not in slice", slice, new(int), false},
			{"nil pointer", slice, nil, false},
			{"empty slice", []*int{}, &a, false},
			{"matching pointer", slice, &a, true},
			{"slice containing nil", []*int{nil}, nil, true},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				result := Contains(tc.slice, tc.value)
				if result != tc.expected {
					t.Errorf("Contains(%v, %v) = %v, expected %v", tc.slice, tc.value, result, tc.expected)
				}
			})
		}
	})

	t.Run("string pointers", func(t *testing.T) {
		hello, world := "hello", "world"
		slice := []*string{&hello, &world}

		testCases := []struct {
			name     string
			slice    []*string
			value    *string
			expected bool
		}{
			{"value in slice", slice, &hello, true},
			{"value not in slice", slice, new(string), false},
			{"nil pointer", slice, nil, false},
			{"empty slice", []*string{}, &hello, false},
			{"slice containing nil", []*string{nil}, nil, true},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				result := Contains(tc.slice, tc.value)
				if result != tc.expected {
					t.Errorf("Contains(%v, %v) = %v, expected %v", tc.slice, tc.value, result, tc.expected)
				}
			})
		}
	})
}

// TestMap tests the Map function with multiple types.
func TestMap(t *testing.T) {
	t.Run("int to int", func(t *testing.T) {
		testCases := []struct {
			name     string
			input    []int
			mapFunc  func(int) int
			expected []int
		}{
			{"double values", []int{1, 2, 3}, func(n int) int { return n * 2 }, []int{2, 4, 6}},
			{"square values", []int{1, 2, 3}, func(n int) int { return n * n }, []int{1, 4, 9}},
			{"empty slice", []int{}, func(n int) int { return n * 2 }, []int{}},
			{"negate values", []int{1, -2, 3}, func(n int) int { return -n }, []int{-1, 2, -3}},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				result := Map(tc.input, tc.mapFunc)
				if len(result) != len(tc.expected) {
					t.Fatalf("Map(%v) length = %d, expected %d", tc.input, len(result), len(tc.expected))
				}
				for i, v := range result {
					if v != tc.expected[i] {
						t.Errorf("Map(%v)[%d] = %d, expected %d", tc.input, i, v, tc.expected[i])
					}
				}
			})
		}
	})

	t.Run("string to string", func(t *testing.T) {
		testCases := []struct {
			name     string
			input    []string
			mapFunc  func(string) string
			expected []string
		}{
			{"uppercase", []string{"a", "b", "c"}, func(s string) string { return s + "!" }, []string{"a!", "b!", "c!"}},
			{"reverse", []string{"go", "is", "fun"}, func(s string) string { return s + s }, []string{"gogo", "isis", "funfun"}},
			{"empty slice", []string{}, func(s string) string { return s + "!" }, []string{}},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				result := Map(tc.input, tc.mapFunc)
				if len(result) != len(tc.expected) {
					t.Fatalf("Map(%v) length = %d, expected %d", tc.input, len(result), len(tc.expected))
				}
				for i, v := range result {
					if v != tc.expected[i] {
						t.Errorf("Map(%v)[%d] = %q, expected %q", tc.input, i, v, tc.expected[i])
					}
				}
			})
		}
	})

	t.Run("struct to string", func(t *testing.T) {
		type User struct {
			Name string
			Age  int
		}

		users := []User{
			{"Alice", 30},
			{"Bob", 25},
		}

		expected := []string{"Alice (30)", "Bob (25)"}
		mapFunc := func(u User) string {
			return fmt.Sprintf("%s (%d)", u.Name, u.Age)
		}

		result := Map(users, mapFunc)

		if len(result) != len(expected) {
			t.Fatalf("Map(%v) length = %d, expected %d", users, len(result), len(expected))
		}

		for i, v := range result {
			if v != expected[i] {
				t.Errorf("Map(%v)[%d] = %q, expected %q", users, i, v, expected[i])
			}
		}
	})
}

func TestMap_NilSlice(t *testing.T) {
	t.Run("nil int slice", func(t *testing.T) {
		var nilSlice []int = nil // Explicitly nil
		result := Map(nilSlice, func(i int) int { return i * 2 })

		if result != nil {
			t.Errorf("Map(nil, f) = %v, expected nil", result)
		}
	})

	t.Run("nil string slice", func(t *testing.T) {
		var nilSlice []string = nil // Explicitly nil
		result := Map(nilSlice, func(s string) string { return s + "!" })

		if result != nil {
			t.Errorf("Map(nil, f) = %v, expected nil", result)
		}
	})

	t.Run("nil struct slice", func(t *testing.T) {
		type Person struct {
			Name string
		}
		var nilSlice []Person = nil // Explicitly nil
		result := Map(nilSlice, func(p Person) string { return p.Name })

		if result != nil {
			t.Errorf("Map(nil, f) = %v, expected nil", result)
		}
	})

	t.Run("nil pointer slice", func(t *testing.T) {
		var nilSlice []*int = nil // Explicitly nil
		result := Map(nilSlice, func(p *int) *int { return p })

		if result != nil {
			t.Errorf("Map(nil, f) = %v, expected nil", result)
		}
	})

	t.Run("nil and empty slice differentiation", func(t *testing.T) {
		var nilSlice []int = nil
		emptySlice := []int{} // Allocated empty slice

		nilResult := Map(nilSlice, func(i int) int { return i * 2 })
		emptyResult := Map(emptySlice, func(i int) int { return i * 2 })

		if nilResult != nil {
			t.Errorf("Map(nilSlice, f) = %v, expected nil", nilResult)
		}
		if emptyResult == nil {
			t.Errorf("Map(emptySlice, f) = nil, expected empty slice")
		}
	})
}
