package main

import (
	"reflect"
	"sort"
	"strings"
	"testing"
)

// Test basic functionality
func TestNewStaticSearchTree(t *testing.T) {
	words := []string{"apple", "app", "application"}
	sst := NewStaticSearchTree(words)
	
	if sst == nil {
		t.Fatal("NewStaticSearchTree returned nil")
	}
	
	if sst.tree == nil {
		t.Fatal("tree map is nil")
	}
	
	if sst.Size() == 0 {
		t.Fatal("tree should not be empty")
	}
}

func TestEmptyWordList(t *testing.T) {
	words := []string{}
	sst := NewStaticSearchTree(words)
	
	if sst.Size() != 0 {
		t.Errorf("Expected size 0 for empty word list, got %d", sst.Size())
	}
	
	results := sst.Search("test")
	if len(results) != 0 {
		t.Errorf("Expected no results for empty tree, got %v", results)
	}
}

func TestSingleWord(t *testing.T) {
	words := []string{"hello"}
	sst := NewStaticSearchTree(words)
	
	// Should have prefixes: h, he, hel, hell, hello
	expectedSize := 5
	if sst.Size() != expectedSize {
		t.Errorf("Expected size %d, got %d", expectedSize, sst.Size())
	}
	
	// Test each prefix
	testCases := []struct {
		query    string
		expected []string
	}{
		{"h", []string{"hello"}},
		{"he", []string{"hello"}},
		{"hel", []string{"hello"}},
		{"hell", []string{"hello"}},
		{"hello", []string{"hello"}},
	}
	
	for _, tc := range testCases {
		results := sst.Search(tc.query)
		if !reflect.DeepEqual(results, tc.expected) {
			t.Errorf("Search('%s'): expected %v, got %v", tc.query, tc.expected, results)
		}
	}
}

func TestBasicSearch(t *testing.T) {
	words := []string{"apple", "app", "application", "banana", "band"}
	sst := NewStaticSearchTree(words)
	
	testCases := []struct {
		query    string
		expected []string
	}{
		{"app", []string{"app", "apple", "application"}},
		{"appl", []string{"apple", "application"}},
		{"ban", []string{"banana", "band"}},
		{"bana", []string{"banana"}},
		{"xyz", []string{}},
		{"", []string{}},
	}
	
	for _, tc := range testCases {
		results := sst.Search(tc.query)
		sort.Strings(results)
		sort.Strings(tc.expected)
		
		if !reflect.DeepEqual(results, tc.expected) {
			t.Errorf("Search('%s'): expected %v, got %v", tc.query, tc.expected, results)
		}
	}
}

func TestCaseInsensitivity(t *testing.T) {
	words := []string{"Apple", "BANANA", "CaR"}
	sst := NewStaticSearchTree(words)
	
	testCases := []struct {
		query    string
		expected []string
	}{
		{"app", []string{"Apple"}},
		{"APP", []string{"Apple"}},
		{"ApP", []string{"Apple"}},
		{"ban", []string{"BANANA"}},
		{"BAN", []string{"BANANA"}},
		{"car", []string{"CaR"}},
		{"CAR", []string{"CaR"}},
		{"Ca", []string{"CaR"}},
	}
	
	for _, tc := range testCases {
		results := sst.Search(tc.query)
		if !reflect.DeepEqual(results, tc.expected) {
			t.Errorf("Search('%s'): expected %v, got %v", tc.query, tc.expected, results)
		}
	}
}

func TestDuplicateWords(t *testing.T) {
	words := []string{"apple", "apple", "banana", "apple"}
	sst := NewStaticSearchTree(words)
	
	results := sst.Search("app")
	// Should only contain "apple" once, despite being in the input multiple times
	expected := []string{"apple"}
	
	if !reflect.DeepEqual(results, expected) {
		t.Errorf("Search with duplicates: expected %v, got %v", expected, results)
	}
}

func TestSearchWithLimit(t *testing.T) {
	words := []string{"app", "apple", "application", "apply", "approach"}
	sst := NewStaticSearchTree(words)
	
	testCases := []struct {
		query    string
		limit    int
		maxLen   int
	}{
		{"app", 2, 2},
		{"app", 10, 5}, // Should return all 5 matches
		{"app", 0, 0},
		{"xyz", 5, 0}, // No matches
	}
	
	for _, tc := range testCases {
		results := sst.SearchWithLimit(tc.query, tc.limit)
		if len(results) > tc.maxLen {
			t.Errorf("SearchWithLimit('%s', %d): expected max %d results, got %d", 
				tc.query, tc.limit, tc.maxLen, len(results))
		}
		if tc.limit > 0 && len(results) > tc.limit {
			t.Errorf("SearchWithLimit('%s', %d): exceeded limit, got %d results", 
				tc.query, tc.limit, len(results))
		}
	}
}

func TestGetAllPrefixes(t *testing.T) {
	words := []string{"hi", "hello"}
	sst := NewStaticSearchTree(words)
	
	prefixes := sst.GetAllPrefixes()
	
	// Expected prefixes: h, he, hel, hell, hello, hi
	expectedPrefixes := []string{"h", "he", "hel", "hell", "hello", "hi"}
	sort.Strings(prefixes)
	sort.Strings(expectedPrefixes)
	
	if !reflect.DeepEqual(prefixes, expectedPrefixes) {
		t.Errorf("GetAllPrefixes(): expected %v, got %v", expectedPrefixes, prefixes)
	}
}

func TestSize(t *testing.T) {
	testCases := []struct {
		words        []string
		expectedSize int
	}{
		{[]string{}, 0},
		{[]string{"a"}, 1},
		{[]string{"ab"}, 2}, // "a", "ab"
		{[]string{"abc"}, 3}, // "a", "ab", "abc"
		{[]string{"a", "ab"}, 2}, // "a", "ab" (no duplicates)
		{[]string{"cat", "car"}, 4}, // "c", "ca", "cat", "car"
	}
	
	for _, tc := range testCases {
		sst := NewStaticSearchTree(tc.words)
		if sst.Size() != tc.expectedSize {
			t.Errorf("Size() for words %v: expected %d, got %d", 
				tc.words, tc.expectedSize, sst.Size())
		}
	}
}

func TestMergeDeduplicate(t *testing.T) {
	testCases := []struct {
		slice1   []string
		slice2   []string
		expected []string
	}{
		{[]string{"a", "b"}, []string{"c", "d"}, []string{"a", "b", "c", "d"}},
		{[]string{"a", "b"}, []string{"b", "c"}, []string{"a", "b", "c"}},
		{[]string{}, []string{"a", "b"}, []string{"a", "b"}},
		{[]string{"a", "b"}, []string{}, []string{"a", "b"}},
		{[]string{}, []string{}, []string{}},
		{[]string{"a", "a", "b"}, []string{"b", "c", "c"}, []string{"a", "b", "c"}},
	}
	
	for _, tc := range testCases {
		result := mergeDeduplicate(tc.slice1, tc.slice2)
		sort.Strings(result)
		sort.Strings(tc.expected)
		
		if !reflect.DeepEqual(result, tc.expected) {
			t.Errorf("mergeDeduplicate(%v, %v): expected %v, got %v", 
				tc.slice1, tc.slice2, tc.expected, result)
		}
	}
}

func TestLargeDataset(t *testing.T) {
	// Generate a larger dataset
	words := make([]string, 1000)
	for i := 0; i < 1000; i++ {
		words[i] = fmt.Sprintf("word%d", i)
	}
	
	sst := NewStaticSearchTree(words)
	
	// Test that it builds successfully
	if sst.Size() == 0 {
		t.Error("Large dataset should produce non-empty tree")
	}
	
	// Test some searches
	results := sst.Search("word1")
	if len(results) == 0 {
		t.Error("Should find matches for 'word1' prefix")
	}
	
	// Test that results contain expected words
	found := false
	for _, word := range results {
		if strings.HasPrefix(word, "word1") {
			found = true
			break
		}
	}
	if !found {
		t.Error("Results should contain words starting with 'word1'")
	}
}

func TestSpecialCharacters(t *testing.T) {
	words := []string{"hello-world", "hello_world", "hello.world", "hello world"}
	sst := NewStaticSearchTree(words)
	
	// Test searching with special characters
	testCases := []struct {
		query    string
		minCount int // Minimum expected results
	}{
		{"hello", 4}, // Should match all variants
		{"hello-", 1},
		{"hello_", 1},
		{"hello.", 1},
		{"hello ", 1},
	}
	
	for _, tc := range testCases {
		results := sst.Search(tc.query)
		if len(results) < tc.minCount {
			t.Errorf("Search('%s'): expected at least %d results, got %d (%v)", 
				tc.query, tc.minCount, len(results), results)
		}
	}
}

func TestUnicodeCharacters(t *testing.T) {
	words := []string{"café", "naïve", "résumé", "jalapeño"}
	sst := NewStaticSearchTree(words)
	
	testCases := []struct {
		query    string
		expected []string
	}{
		{"caf", []string{"café"}},
		{"naï", []string{"naïve"}},
		{"rés", []string{"résumé"}},
		{"jal", []string{"jalapeño"}},
	}
	
	for _, tc := range testCases {
		results := sst.Search(tc.query)
		if !reflect.DeepEqual(results, tc.expected) {
			t.Errorf("Search('%s'): expected %v, got %v", tc.query, tc.expected, results)
		}
	}
}

// Benchmark tests
func BenchmarkBuild(b *testing.B) {
	words := make([]string, 1000)
	for i := 0; i < 1000; i++ {
		words[i] = fmt.Sprintf("word%d", i)
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		NewStaticSearchTree(words)
	}
}

func BenchmarkSearch(b *testing.B) {
	words := make([]string, 1000)
	for i := 0; i < 1000; i++ {
		words[i] = fmt.Sprintf("word%d", i)
	}
	sst := NewStaticSearchTree(words)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sst.Search("word1")
	}
}

func BenchmarkSearchWithLimit(b *testing.B) {
	words := make([]string, 1000)
	for i := 0; i < 1000; i++ {
		words[i] = fmt.Sprintf("word%d", i)
	}
	sst := NewStaticSearchTree(words)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sst.SearchWithLimit("word1", 10)
	}
}

// Example test demonstrating usage
func ExampleStaticSearchTree() {
	words := []string{"apple", "app", "application", "banana"}
	sst := NewStaticSearchTree(words)
	
	results := sst.Search("app")
	sort.Strings(results) // Sort for consistent output
	fmt.Println(results)
	// Output: [app apple application]
}

// Test that search results are not modifiable (defensive copying)
func TestSearchResultsImmutability(t *testing.T) {
	words := []string{"apple", "app"}
	sst := NewStaticSearchTree(words)
	
	results1 := sst.Search("app")
	results2 := sst.Search("app")
	
	// Modify first result set
	if len(results1) > 0 {
		results1[0] = "modified"
	}
	
	// Second result set should be unchanged
	if len(results2) > 0 && results2[0] == "modified" {
		t.Error("Search results should be independent copies")
	}
}
