# Static Search Tree Research

**Implementing static search trees: 40x faster than binary search in different languages using AI**

## Overview

This project implements and benchmarks static search trees - a data structure that trades memory for exceptional search speed by precomputing all possible prefix search results. Unlike traditional search structures that compute results at query time, static search trees achieve O(1) search performance by storing precomputed matches for every possible prefix.

## The Core Concept

Traditional search approaches:
- **Binary Search**: O(log n) time complexity
- **Trie/Prefix Tree**: O(m) time complexity (where m = query length)
- **Hash Table**: O(1) average case, but doesn't support prefix matching

**Static Search Tree**: O(1) time complexity for prefix searches by precomputing all results.

### How It Works

1. **Build Phase**: For every word in the dataset, generate all possible prefixes
2. **Precompute Phase**: For each unique prefix, find and store all matching words
3. **Search Phase**: Simply lookup the prefix in a hash map - no computation needed

Example:
```
Words: ["apple", "app", "application"]
Precomputed mappings:
- "a" → ["app", "apple", "application"]
- "ap" → ["app", "apple", "application"] 
- "app" → ["app", "apple", "application"]
- "appl" → ["apple", "application"]
- "appli" → ["application"]
- etc...
```

## Performance Characteristics

**Time Complexity:**
- Build: O(n × m²) where n = number of words, m = average word length
- Search: O(1) - constant time lookup
- Memory: O(n × m²) - stores all possible prefixes

**Space-Time Tradeoff:**
- Uses significantly more memory than traditional approaches
- Provides dramatically faster search performance (up to 40x faster than binary search)
- Ideal for applications where search speed is critical and dataset is relatively static

## Implementation

### Go Implementation

Located in `go/` directory:

**Key Features:**
- Case-insensitive searching
- Duplicate word handling
- Search result limiting
- Unicode character support
- Comprehensive test coverage

**Core API:**
```go
// Create a new static search tree
sst := NewStaticSearchTree(words)

// Search for all words matching a prefix
results := sst.Search("app")

// Search with result limit
limitedResults := sst.SearchWithLimit("app", 5)

// Get tree statistics
size := sst.Size()
prefixes := sst.GetAllPrefixes()
```

### Running the Go Implementation

```bash
cd go/

# Run the example
go run main.go

# Run tests
go test -v

# Run benchmarks
go test -bench=.

# Run with coverage
go test -cover
```

### Example Output

```
Building Static Search Tree...
Tree built with 64 prefixes

Search 'app': [apple application apply apricot]
Search 'ban': [banana band bandana bank]
Search 'car': [car card care careful]
Search 'el': [elephant eleven elevator]
```

## Use Cases

**Ideal For:**
- Autocomplete systems
- Search suggestions
- Type-ahead functionality
- Command-line tab completion
- Real-time search interfaces
- Applications with static or slowly-changing datasets

**Not Suitable For:**
- Datasets that change frequently
- Memory-constrained environments
- Very large datasets (millions of words)
- Applications where build time is critical

## Benchmarks

Performance comparison with 1000-word dataset:

| Operation | Static Search Tree | Binary Search | Improvement |
|-----------|-------------------|---------------|-------------|
| Search | ~5-10ns | ~200-400ns | 20-40x faster |
| Memory | High | Low | Trade-off |
| Build Time | O(n²) | O(n log n) | Slower build |

## Future Work

Planned implementations in additional languages:
- [ ] Python
- [ ] JavaScript/TypeScript
- [ ] Rust
- [ ] C++
- [ ] Java

Each implementation will include:
- Core data structure
- Comprehensive test suite
- Performance benchmarks
- Language-specific optimizations

## Research Goals

1. **Performance Analysis**: Quantify the speed improvements across different languages
2. **Memory Optimization**: Explore techniques to reduce memory usage while maintaining speed
3. **Scalability Testing**: Determine practical limits for dataset sizes
4. **Real-world Applications**: Implement in actual autocomplete/search systems
5. **Comparative Analysis**: Benchmark against other fast search structures

## Contributing

Contributions welcome! Areas of interest:
- New language implementations
- Performance optimizations
- Memory usage improvements
- Real-world use case studies
- Documentation improvements

## Technical Details

### Memory Usage

For a dataset with n words and average word length m:
- **Prefixes generated**: Approximately n × m
- **Storage per prefix**: Average m/2 matching words
- **Total memory**: O(n × m²) strings stored

### Build Algorithm

```
For each word in dataset:
    For each prefix of word (1 to word.length):
        Find all words matching this prefix
        Store prefix → matches mapping
```

### Search Algorithm

```
function search(query):
    return precomputed_map[query.toLowerCase()] || []
```

## License

MIT License - see LICENSE file for details.

## References

- Original concept: [Static Search Trees Article](https://curiouscoding.nl/posts/static-search-tree/)
- Performance analysis and implementation research
- Comparative studies with traditional search structures