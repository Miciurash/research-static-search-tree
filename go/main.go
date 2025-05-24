package main

import (
	"fmt"
	"sort"
	"strings"
)

// StaticSearchTree represents a precomputed search tree for efficient prefix matching
type StaticSearchTree struct {
	tree map[string][]string
}

// NewStaticSearchTree creates a new static search tree from a list of words
func NewStaticSearchTree(words []string) *StaticSearchTree {
	sst := &StaticSearchTree{
		tree: make(map[string][]string),
	}
	sst.build(words)
	return sst
}

// build constructs the static search tree by precomputing all prefix combinations
func (sst *StaticSearchTree) build(words []string) {
	// Sort words to ensure consistent ordering
	sort.Strings(words)
	
	// For each word, generate all possible prefixes and their matching results
	for _, word := range words {
		// Generate all prefixes of the word
		for i := 1; i <= len(word); i++ {
			prefix := strings.ToLower(word[:i])
			
			// Find all words that match this prefix
			var matches []string
			for _, candidate := range words {
				if strings.HasPrefix(strings.ToLower(candidate), prefix) {
					matches = append(matches, candidate)
				}
			}
			
			// Store the matches for this prefix (avoiding duplicates)
			if existing, exists := sst.tree[prefix]; exists {
				// Merge and deduplicate
				merged := mergeDeduplicate(existing, matches)
				sst.tree[prefix] = merged
			} else {
				sst.tree[prefix] = matches
			}
		}
	}
}

// mergeDeduplicate merges two slices and removes duplicates
func mergeDeduplicate(slice1, slice2 []string) []string {
	seen := make(map[string]bool)
	var result []string
	
	// Add all items from both slices
	for _, item := range slice1 {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}
	
	for _, item := range slice2 {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}
	
	return result
}

// Search performs a prefix search and returns all matching words
func (sst *StaticSearchTree) Search(query string) []string {
	query = strings.ToLower(query)
	if matches, exists := sst.tree[query]; exists {
		// Return a copy to prevent external modification
		result := make([]string, len(matches))
		copy(result, matches)
		return result
	}
	return []string{}
}

// SearchWithLimit performs a prefix search with a maximum number of results
func (sst *StaticSearchTree) SearchWithLimit(query string, limit int) []string {
	matches := sst.Search(query)
	if len(matches) <= limit {
		return matches
	}
	return matches[:limit]
}

// GetAllPrefixes returns all stored prefixes (useful for debugging)
func (sst *StaticSearchTree) GetAllPrefixes() []string {
	var prefixes []string
	for prefix := range sst.tree {
		prefixes = append(prefixes, prefix)
	}
	sort.Strings(prefixes)
	return prefixes
}

// Size returns the number of stored prefixes
func (sst *StaticSearchTree) Size() int {
	return len(sst.tree)
}

// PrintTree prints the entire tree structure (for debugging)
func (sst *StaticSearchTree) PrintTree() {
	prefixes := sst.GetAllPrefixes()
	for _, prefix := range prefixes {
		fmt.Printf("'%s' -> %v\n", prefix, sst.tree[prefix])
	}
}

// Example usage and demonstration
func main() {
	// Example word list - could be loaded from a file or database
	words := []string{
		"apple", "application", "apply", "apricot",
		"banana", "band", "bandana", "bank",
		"cat", "car", "card", "care", "careful",
		"dog", "door", "double",
		"elephant", "eleven", "elevator",
	}
	
	fmt.Println("Building Static Search Tree...")
	sst := NewStaticSearchTree(words)
	
	fmt.Printf("Tree built with %d prefixes\n\n", sst.Size())
	
	// Example searches
	queries := []string{"app", "ban", "car", "el", "z", "do"}
	
	for _, query := range queries {
		results := sst.Search(query)
		fmt.Printf("Search '%s': %v\n", query, results)
	}
	
	fmt.Println("\n--- Limited Results (max 3) ---")
	for _, query := range queries {
		results := sst.SearchWithLimit(query, 3)
		fmt.Printf("Search '%s' (limit 3): %v\n", query, results)
	}
	
	// Demonstrate case insensitivity
	fmt.Println("\n--- Case Insensitive Search ---")
	caseQueries := []string{"APP", "Car", "EL"}
	for _, query := range caseQueries {
		results := sst.Search(query)
		fmt.Printf("Search '%s': %v\n", query, results)
	}
	
	// Show some tree structure for debugging
	fmt.Println("\n--- Sample Tree Structure ---")
	samplePrefixes := []string{"a", "ap", "app", "car", "el"}
	for _, prefix := range samplePrefixes {
		if matches, exists := sst.tree[prefix]; exists {
			fmt.Printf("'%s' -> %v\n", prefix, matches)
		}
	}
}

// Benchmark function to test performance
func BenchmarkSearch(sst *StaticSearchTree, queries []string, iterations int) {
	fmt.Printf("\n--- Performance Test (%d iterations) ---\n", iterations)
	
	for _, query := range queries {
		// Time the search operations
		results := sst.Search(query)
		fmt.Printf("Query '%s': %d results\n", query, len(results))
	}
}
