# LRU Cache

A simple and educational implementation of an LRU (Least Recently Used) Cache in Go. Designed for learning purposes and small-scale usage, this project demonstrates how an LRU eviction policy can be applied to cache key-value data efficiently.

## Features

- Least Recently Used (LRU) eviction policy
- O(1) complexity for `Get` and `Put` operations
- Thread-safe with `sync.Mutex`
- Customizable maximum number of entries
- Optional logging and eviction callbacks

## Getting Started

### Requirements

- Go 1.20+

### Installation

Clone the repository:

```bash
git clone https://github.com/seu-usuario/lru-cache-go.git
cd lru-cache-go
```

### Usage

```go
package main

import (
	"fmt"
	"lru"
)

func main() {
	cache := lru.NewLRUCache(3, lru.Options{
		Logs: true, // enable logging
	})

	cache.Set("a", 1)
	cache.Set("b", 2)
	cache.Set("c", 3)

	if val, ok := cache.Get("a"); ok {
		fmt.Println("Found:", val)
	}

	cache.Set("d", 4) // Insert "d" — this will evict the least recently used key, which is "b"

	// Try to get "b" — it should have been evicted
	if _, ok := cache.Get("b"); !ok {
		fmt.Println("b has been evicted")
	}
}
```

### Running tests

```bash
go test ./...
```

For concurrency tests:

```bash
go test ./... -v -race
```

## Technical Documentation

### What is a Cache?

A cache is a temporary storage mechanism that holds frequently accessed data to improve performance and reduce the cost of repeated data retrieval or computation.
What is an Eviction Policy?

When a cache reaches its storage limit, an eviction policy decides which entry should be removed to make room for new ones.

Common eviction strategies include:
- **LRU** (Least Recently Used): removes the entry that hasn’t been used for the longest time.
- **LFU** (Least Frequently Used): removes the entry with the lowest access count.
- **FIFO** (First In, First Out): removes the earliest added entry, regardless of usage.

### How LRU Cache Works

- **Purpose:** Keeps recently accessed data available for fast retrieval, evicting the least recently accessed data when the cache is full.
- **O(1) performance** for both Set and Get operations by combining:
  - A map[string]*list.Element for direct access by key.
  - A container/list (doubly-linked list) to track usage order.
- **Automatic Eviction:** When maxEntries is exceeded, the least recently used entry (tail of the list) is removed.
- **Eviction Callback:** Supports a user-defined function evictCallback triggered when an entry is evicted.
- **Optional Logging:** The logs flag in Options enables printing of the current cache order and actions taken.


<img width="1121" height="325" alt="Data Representation" src="https://github.com/user-attachments/assets/a176ef50-ca96-4560-a10b-2a1ebf7b03f1" />

### Asymptotic Complexity of LRU Cache Operations

This LRU cache implementation is optimized for performance, offering **constant time O(1)** operations for both retrieval and updates, thanks to the combination of a hashmap and a doubly linked list.

#### Time Complexity Table

| Method              | Time Complexity | Description |
|---------------------|------------------|-------------|
| `Set(key, value)`   | **O(1)**         | Inserts or updates a key-value pair. If the key exists, it updates the value and moves it to the front. If the cache exceeds its size, it evicts the least recently used item. |
| `Get(key)`          | **O(1)**         | Retrieves the value and moves the accessed item to the front of the list (marks as most recently used). |
| `removeElement(e)`  | **O(1)**         | Removes a specific element from the doubly linked list and the map. |
| `removeLastElement()`| **O(1)**        | Removes the least recently used item (at the end of the list). |
| `printList()`       | **O(n)**         | Iterates over the entire list to print all keys and values — only used for debugging. |

#### Why O(1) is Possible

- The **hash map** (`map[string]*list.Element`) allows constant-time access to elements by key.
- The **doubly linked list** (`container/list`) provides constant-time operations to move elements to the front or remove them from the middle/end.
- This structure ensures that every operation critical to LRU behavior (`Set`, `Get`, `Evict`) can be done in constant time, regardless of the cache size.

#### Note

- The only method with linear time complexity is `printList()`, used for debugging purposes. It does **not** affect runtime performance in production scenarios.
