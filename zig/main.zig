const std = @import("std");
const print = std.debug.print;
const ArrayList = std.ArrayList;
const HashMap = std.HashMap;
const Allocator = std.mem.Allocator;
const testing = std.testing;

/// StaticSearchTree represents a precomputed search tree for efficient prefix matching
pub const StaticSearchTree = struct {
    const Self = @This();
    
    // HashMap with string keys and ArrayList of strings as values
    tree: HashMap([]const u8, ArrayList([]const u8), StringContext, std.hash_map.default_max_load_percentage),
    allocator: Allocator,
    
    // String context for HashMap
    const StringContext = struct {
        pub fn hash(self: @This(), s: []const u8) u64 {
            _ = self;
            return std.hash_map.hashString(s);
        }
        pub fn eql(self: @This(), a: []const u8, b: []const u8) bool {
            _ = self;
            return std.mem.eql(u8, a, b);
        }
    };
    
    /// Initialize a new StaticSearchTree
    pub fn init(allocator: Allocator) Self {
        return Self{
            .tree = HashMap([]const u8, ArrayList([]const u8), StringContext, std.hash_map.default_max_load_percentage).init(allocator),
            .allocator = allocator,
        };
    }
    
    /// Deinitialize the StaticSearchTree and free all memory
    pub fn deinit(self: *Self) void {
        var iterator = self.tree.iterator();
        while (iterator.next()) |entry| {
            // Free the key (prefix)
            self.allocator.free(entry.key_ptr.*);
            // Free the value (ArrayList and its strings)
            for (entry.value_ptr.items) |word| {
                self.allocator.free(word);
            }
            entry.value_ptr.deinit();
        }
        self.tree.deinit();
    }
    
    /// Build the static search tree from a list of words
    pub fn build(self: *Self, words: []const []const u8) !void {
        // Sort words for consistent ordering
        var sorted_words = try self.allocator.alloc([]const u8, words.len);
        defer self.allocator.free(sorted_words);
        
        for (words, 0..) |word, i| {
            sorted_words[i] = word;
        }
        std.mem.sort([]const u8, sorted_words, {}, stringLessThan);
        
        // For each word, generate all possible prefixes and their matching results
        for (sorted_words) |word| {
            // Generate all prefixes of the word
            var i: usize = 1;
            while (i <= word.len) : (i += 1) {
                const prefix_slice = word[0..i];
                const prefix = try toLowerCaseAlloc(self.allocator, prefix_slice);
                
                // Find all words that match this prefix
                var matches = ArrayList([]const u8).init(self.allocator);
                
                for (sorted_words) |candidate| {
                    const candidate_lower = try toLowerCaseAlloc(self.allocator, candidate);
                    defer self.allocator.free(candidate_lower);
                    
                    if (std.mem.startsWith(u8, candidate_lower, prefix)) {
                        const candidate_copy = try self.allocator.dupe(u8, candidate);
                        try matches.append(candidate_copy);
                    }
                }
                
                // Store the matches for this prefix (handling duplicates)
                const result = try self.tree.getOrPut(prefix);
                if (result.found_existing) {
                    // Merge and deduplicate
                    const merged = try mergeDeduplicate(self.allocator, result.value_ptr.*, matches);
                    
                    // Free old list
                    for (result.value_ptr.items) |word| {
                        self.allocator.free(word);
                    }
                    result.value_ptr.deinit();
                    
                    // Free temporary matches list
                    for (matches.items) |word| {
                        self.allocator.free(word);
                    }
                    matches.deinit();
                    
                    result.value_ptr.* = merged;
                } else {
                    result.value_ptr.* = matches;
                }
            }
        }
    }
    
    /// Search for all words matching the given prefix
    pub fn search(self: *Self, query: []const u8) !ArrayList([]const u8) {
        const query_lower = try toLowerCaseAlloc(self.allocator, query);
        defer self.allocator.free(query_lower);
        
        var result = ArrayList([]const u8).init(self.allocator);
        
        if (self.tree.get(query_lower)) |matches| {
            for (matches.items) |word| {
                const word_copy = try self.allocator.dupe(u8, word);
                try result.append(word_copy);
            }
        }
        
        return result;
    }
    
    /// Search with a limit on the number of results
    pub fn searchWithLimit(self: *Self, query: []const u8, limit: usize) !ArrayList([]const u8) {
        var matches = try self.search(query);
        
        if (matches.items.len <= limit) {
            return matches;
        }
        
        // Create a new list with only the first 'limit' items
        var limited = ArrayList([]const u8).init(self.allocator);
        for (matches.items[0..limit]) |word| {
            const word_copy = try self.allocator.dupe(u8, word);
            try limited.append(word_copy);
        }
        
        // Free the original matches
        for (matches.items) |word| {
            self.allocator.free(word);
        }
        matches.deinit();
        
        return limited;
    }
    
    /// Get all stored prefixes (useful for debugging)
    pub fn getAllPrefixes(self: *Self) !ArrayList([]const u8) {
        var prefixes = ArrayList([]const u8).init(self.allocator);
        
        var iterator = self.tree.iterator();
        while (iterator.next()) |entry| {
            const prefix_copy = try self.allocator.dupe(u8, entry.key_ptr.*);
            try prefixes.append(prefix_copy);
        }
        
        // Sort the prefixes
        std.mem.sort([]const u8, prefixes.items, {}, stringLessThan);
        
        return prefixes;
    }
    
    /// Get the number of stored prefixes
    pub fn size(self: *Self) usize {
        return self.tree.count();
    }
    
    /// Print the entire tree structure (for debugging)
    pub fn printTree(self: *Self) !void {
        var prefixes = try self.getAllPrefixes();
        defer {
            for (prefixes.items) |prefix| {
                self.allocator.free(prefix);
            }
            prefixes.deinit();
        }
        
        for (prefixes.items) |prefix| {
            print("'{s}' -> [", .{prefix});
            if (self.tree.get(prefix)) |matches| {
                for (matches.items, 0..) |word, i| {
                    if (i > 0) print(", ");
                    print("'{s}'", .{word});
                }
            }
            print("]\n");
        }
    }
};

/// Helper function to convert string to lowercase
fn toLowerCaseAlloc(allocator: Allocator, input: []const u8) ![]u8 {
    var result = try allocator.alloc(u8, input.len);
    for (input, 0..) |char, i| {
        result[i] = std.ascii.toLower(char);
    }
    return result;
}

/// Helper function for string comparison (for sorting)
fn stringLessThan(context: void, a: []const u8, b: []const u8) bool {
    _ = context;
    return std.mem.lessThan(u8, a, b);
}

/// Merge two ArrayLists and remove duplicates
fn mergeDeduplicate(allocator: Allocator, list1: ArrayList([]const u8), list2: ArrayList([]const u8)) !ArrayList([]const u8) {
    var seen = HashMap([]const u8, void, StaticSearchTree.StringContext, std.hash_map.default_max_load_percentage).init(allocator);
    defer seen.deinit();
    
    var result = ArrayList([]const u8).init(allocator);
    
    // Add items from first list
    for (list1.items) |item| {
        const item_copy = try allocator.dupe(u8, item);
        const get_result = try seen.getOrPut(item_copy);
        if (!get_result.found_existing) {
            try result.append(item_copy);
        } else {
            allocator.free(item_copy);
        }
    }
    
    // Add items from second list
    for (list2.items) |item| {
        const item_copy = try allocator.dupe(u8, item);
        const get_result = try seen.getOrPut(item_copy);
        if (!get_result.found_existing) {
            try result.append(item_copy);
        } else {
            allocator.free(item_copy);
        }
    }
    
    // Free the seen map keys
    var iterator = seen.iterator();
    while (iterator.next()) |entry| {
        allocator.free(entry.key_ptr.*);
    }
    
    return result;
}

/// Helper function to free an ArrayList of strings
fn freeStringList(allocator: Allocator, list: ArrayList([]const u8)) void {
    for (list.items) |string| {
        allocator.free(string);
    }
}

// Example usage and demonstration
pub fn main() !void {
    var gpa = std.heap.GeneralPurposeAllocator(.{}){};
    defer _ = gpa.deinit();
    const allocator = gpa.allocator();
    
    // Example word list
    const words = [_][]const u8{
        "apple", "application", "apply", "apricot",
        "banana", "band", "bandana", "bank",
        "cat", "car", "card", "care", "careful",
        "dog", "door", "double",
        "elephant", "eleven", "elevator",
    };
    
    print("Building Static Search Tree...\n");
    var sst = StaticSearchTree.init(allocator);
    defer sst.deinit();
    
    try sst.build(&words);
    print("Tree built with {} prefixes\n\n", .{sst.size()});
    
    // Example searches
    const queries = [_][]const u8{ "app", "ban", "car", "el", "z", "do" };
    
    for (queries) |query| {
        var results = try sst.search(query);
        defer {
            freeStringList(allocator, results);
            results.deinit();
        }
        
        print("Search '{s}': [", .{query});
        for (results.items, 0..) |word, i| {
            if (i > 0) print(", ");
            print("'{s}'", .{word});
        }
        print("]\n");
    }
    
    print("\n--- Limited Results (max 3) ---\n");
    for (queries) |query| {
        var results = try sst.searchWithLimit(query, 3);
        defer {
            freeStringList(allocator, results);
            results.deinit();
        }
        
        print("Search '{s}' (limit 3): [", .{query});
        for (results.items, 0..) |word, i| {
            if (i > 0) print(", ");
            print("'{s}'", .{word});
        }
        print("]\n");
    }
    
    // Demonstrate case insensitivity
    print("\n--- Case Insensitive Search ---\n");
    const case_queries = [_][]const u8{ "APP", "Car", "EL" };
    for (case_queries) |query| {
        var results = try sst.search(query);
        defer {
            freeStringList(allocator, results);
            results.deinit();
        }
        
        print("Search '{s}': [", .{query});
        for (results.items, 0..) |word, i| {
            if (i > 0) print(", ");
            print("'{s}'", .{word});
        }
        print("]\n");
    }
    
    // Show some tree structure for debugging
    print("\n--- Sample Tree Structure ---\n");
    const sample_prefixes = [_][]const u8{ "a", "ap", "app", "car", "el" };
    for (sample_prefixes) |prefix| {
        if (sst.tree.get(prefix)) |matches| {
            print("'{s}' -> [", .{prefix});
            for (matches.items, 0..) |word, i| {
                if (i > 0) print(", ");
                print("'{s}'", .{word});
            }
            print("]\n");
        }
    }
}

// Tests
test "basic functionality" {
    const allocator = testing.allocator;
    
    var sst = StaticSearchTree.init(allocator);
    defer sst.deinit();
    
    const words = [_][]const u8{ "apple", "app", "application" };
    try sst.build(&words);
    
    try testing.expect(sst.size() > 0);
}

test "empty word list" {
    const allocator = testing.allocator;
    
    var sst = StaticSearchTree.init(allocator);
    defer sst.deinit();
    
    const words = [_][]const u8{};
    try sst.build(&words);
    
    try testing.expect(sst.size() == 0);
    
    var results = try sst.search("test");
    defer {
        freeStringList(allocator, results);
        results.deinit();
    }
    
    try testing.expect(results.items.len == 0);
}

test "single word" {
    const allocator = testing.allocator;
    
    var sst = StaticSearchTree.init(allocator);
    defer sst.deinit();
    
    const words = [_][]const u8{"hello"};
    try sst.build(&words);
    
    // Should have prefixes: h, he, hel, hell, hello
    try testing.expect(sst.size() == 5);
    
    // Test each prefix
    const test_cases = [_]struct { query: []const u8, expected_count: usize }{
        .{ .query = "h", .expected_count = 1 },
        .{ .query = "he", .expected_count = 1 },
        .{ .query = "hel", .expected_count = 1 },
        .{ .query = "hell", .expected_count = 1 },
        .{ .query = "hello", .expected_count = 1 },
    };
    
    for (test_cases) |tc| {
        var results = try sst.search(tc.query);
        defer {
            freeStringList(allocator, results);
            results.deinit();
        }
        
        try testing.expect(results.items.len == tc.expected_count);
        if (results.items.len > 0) {
            try testing.expect(std.mem.eql(u8, results.items[0], "hello"));
        }
    }
}

test "basic search" {
    const allocator = testing.allocator;
    
    var sst = StaticSearchTree.init(allocator);
    defer sst.deinit();
    
    const words = [_][]const u8{ "apple", "app", "application", "banana", "band" };
    try sst.build(&words);
    
    const test_cases = [_]struct { query: []const u8, min_expected: usize }{
        .{ .query = "app", .min_expected = 3 },
        .{ .query = "appl", .min_expected = 2 },
        .{ .query = "ban", .min_expected = 2 },
        .{ .query = "bana", .min_expected = 1 },
        .{ .query = "xyz", .min_expected = 0 },
        .{ .query = "", .min_expected = 0 },
    };
    
    for (test_cases) |tc| {
        var results = try sst.search(tc.query);
        defer {
            freeStringList(allocator, results);
            results.deinit();
        }
        
        try testing.expect(results.items.len >= tc.min_expected);
    }
}

test "case insensitivity" {
    const allocator = testing.allocator;
    
    var sst = StaticSearchTree.init(allocator);
    defer sst.deinit();
    
    const words = [_][]const u8{ "Apple", "BANANA", "CaR" };
    try sst.build(&words);
    
    const test_cases = [_]struct { query: []const u8, expected_count: usize }{
        .{ .query = "app", .expected_count = 1 },
        .{ .query = "APP", .expected_count = 1 },
        .{ .query = "ApP", .expected_count = 1 },
        .{ .query = "ban", .expected_count = 1 },
        .{ .query = "BAN", .expected_count = 1 },
        .{ .query = "car", .expected_count = 1 },
        .{ .query = "CAR", .expected_count = 1 },
        .{ .query = "Ca", .expected_count = 1 },
    };
    
    for (test_cases) |tc| {
        var results = try sst.search(tc.query);
        defer {
            freeStringList(allocator, results);
            results.deinit();
        }
        
        try testing.expect(results.items.len == tc.expected_count);
    }
}

test "search with limit" {
    const allocator = testing.allocator;
    
    var sst = StaticSearchTree.init(allocator);
    defer sst.deinit();
    
    const words = [_][]const u8{ "app", "apple", "application", "apply", "approach" };
    try sst.build(&words);
    
    const test_cases = [_]struct { query: []const u8, limit: usize, max_expected: usize }{
        .{ .query = "app", .limit = 2, .max_expected = 2 },
        .{ .query = "app", .limit = 10, .max_expected = 5 },
        .{ .query = "app", .limit = 0, .max_expected = 0 },
        .{ .query = "xyz", .limit = 5, .max_expected = 0 },
    };
    
    for (test_cases) |tc| {
        var results = try sst.searchWithLimit(tc.query, tc.limit);
        defer {
            freeStringList(allocator, results);
            results.deinit();
        }
        
        try testing.expect(results.items.len <= tc.max_expected);
        if (tc.limit > 0) {
            try testing.expect(results.items.len <= tc.limit);
        }
    }
}