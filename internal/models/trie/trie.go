package trie

import (
	"sort"
	"strings"
	"sync"
	"time"

	lru "github.com/hashicorp/golang-lru"
)

type TrieNode struct {
	Children  map[rune]*TrieNode `json:"children,omitempty"`
	IsEnd     bool               `json:"isEnd,omitempty"`
	Word      string             `json:"word,omitempty"`
	Score     float64            `json:"score,omitempty"`
	Frequency int                `json:"frequency,omitempty"`
	LastUsed  time.Time          `json:"lastUsed,omitempty"`
	UpdatedAt time.Time          `json:"updatedAt,omitempty"`
}

type Trie struct {
	root       *TrieNode
	cache      *lru.Cache
	mutex      sync.RWMutex
	maxResults int
}

func NewTrie(cacheSize, maxResults int) (*Trie, error) {
	cache, err := lru.New(cacheSize)
	if err != nil {
		return nil, err
	}
	return &Trie{
		root: &TrieNode{
			Children: make(map[rune]*TrieNode),
		},
		cache:      cache,
		maxResults: maxResults,
	}, nil
}

func (t *Trie) Insert(word string, score float64) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	node := t.root
	for _, char := range strings.ToLower(word) {
		if node.Children == nil {
			node.Children = make(map[rune]*TrieNode)
		}
		if _, ok := node.Children[char]; !ok {
			node.Children[char] = &TrieNode{
				Children: make(map[rune]*TrieNode),
			}
		}
		node = node.Children[char]
	}
	node.IsEnd = true
	node.Word = word
	node.Score = score
	node.UpdatedAt = time.Now()
}

func (t *Trie) SearchWithSuggestions(prefix string, limit int) []SuggestionResult {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	if cached, ok := t.cache.Get(prefix); ok {
		return cached.([]SuggestionResult)
	}

	node := t.root

	for _, char := range strings.ToLower(prefix) {
		if node.Children == nil {
			node.Children = make(map[rune]*TrieNode)
		}
		if _, ok := node.Children[char]; !ok {
			return nil
		}
		node = node.Children[char]
	}

	suggestions := t.collectSuggestions(node, prefix)

	sort.Slice(suggestions, func(i, j int) bool {
		scoreI := suggestions[i].Score * float64(suggestions[i].Frequency)
		scoreJ := suggestions[j].Score * float64(suggestions[j].Frequency)

		recencyBoostI := time.Since(suggestions[i].LastUsed).Hours()
		recencyBoostJ := time.Since(suggestions[j].LastUsed).Hours()

		scoreI = scoreI * (1 + 1/recencyBoostI)
		scoreJ = scoreJ * (1 + 1/recencyBoostJ)

		return scoreI > scoreJ
	})

	if len(suggestions) > limit {
		suggestions = suggestions[:limit]
	}
	return suggestions
}

func (t *Trie) collectSuggestions(node *TrieNode, prefix string) []SuggestionResult {
	var results []SuggestionResult

	if node.IsEnd {
		results = append(results, SuggestionResult{
			Word:      node.Word,
			Score:     float64(node.Score),
			Frequency: node.Frequency,
			LastUsed:  node.LastUsed,
		})
	}

	for char, childNode := range node.Children {
		results = append(results, t.collectSuggestions(childNode, prefix+string(char))...)
	}

	return results
}

func (t *Trie) RecordSelection(word string) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	node := t.root

	for _, char := range strings.ToLower(word) {
		if node.Children == nil {
			node.Children = make(map[rune]*TrieNode)
		}
		if _, ok := node.Children[char]; !ok {
			node.Frequency++
			node.LastUsed = time.Now()
		}
	}
}

func (t *Trie) Cleanup(maxAge time.Duration) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.cleanupNode(t.root, maxAge)
}

func (t *Trie) cleanupNode(node *TrieNode, maxAge time.Duration) bool {
	if node == nil {
		return true
	}

	if !node.LastUsed.IsZero() && time.Since(node.LastUsed) > maxAge {
		return true
	}

	for char, childNode := range node.Children {
		if t.cleanupNode(childNode, maxAge) {
			delete(node.Children, char)
		}
	}

	return len(node.Children) == 0 && !node.IsEnd
}
