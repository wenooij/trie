package trie

import "fmt"

type Set struct {
	children map[string]*Set
	count    int
}

func (a *Set) Count() int { return a.count }

func (a *Set) empty() bool {
	return len(a.children) == 0 && a.count == 0
}

func (a *Set) Contains(s string) bool { return a.containsAtLeast(s, 1) }

func (a *Set) containsAtLeast(s string, n int) bool {
	if s == "" {
		return a.count >= n
	}
	// Fast path: substring match.
	if len(s) < len(a.children) {
		for i := len(s); i > 0; i-- {
			if child, ok := a.children[s[:i]]; ok {
				return child.Contains(s[i:])
			}
		}
	}
	// Slow path: prefix match.
	for prefix, child := range a.children {
		if n := prefixLen(prefix, s); n > 0 {
			return child.Contains(s[n:])
		}
	}
	// Default case: no prefix matches.
	return false
}

func (a *Set) Merge(other *Set) {
	if a == nil || other == nil {
		return
	}
	for prefix, child := range other.children {
		a.children[prefix].Merge(child)
	}
}

func (a *Set) putChild(prefix string, child *Set) *Set {
	if a.children == nil {
		a.children = make(map[string]*Set)
	}
	a.children[prefix] = child
	return child
}

func (a *Set) Insert(s string) *Set {
	if s == "" {
		a.count++
		return a
	}
	if len(s) < len(a.children) {
		// Search prefixes of s when it's is shorter.
		for i := len(s); i > 0; i-- {
			if child, ok := a.children[s[:i]]; ok {
				return child.Insert(s[i:])
			}
		}
		// Insert new prefix.
		return a.putChild(s, &Set{count: 1})
	}
	var (
		maxPrefix    string
		maxPrefixLen int
		maxChild     *Set
	)
	for prefix, child := range a.children {
		if n := prefixLen(s, prefix); n > 0 {
			maxPrefix = prefix
			maxPrefixLen = n
			maxChild = child
			break
		}
	}
	if maxPrefixLen == 0 {
		// Insert new prefix.
		return a.putChild(s, &Set{count: 1})
	}
	child := maxChild
	if maxPrefixLen < len(maxPrefix) {
		// Split existing prefix.
		newPrefix, newPostfix := maxPrefix[:maxPrefixLen], maxPrefix[maxPrefixLen:]
		// Create a new child for newPrefix.
		child = &Set{}
		child.putChild(newPostfix, maxChild)
		// Delete the old prefix.
		delete(a.children, maxPrefix)
		// Insert the new one.
		a.putChild(newPrefix, child)
	}
	// Insert the new child.
	return child.Insert(s[maxPrefixLen:])
}

func prefixLen(s1, s2 string) int {
	if s1 == s2 {
		return len(s1)
	}
	n := len(s1)
	if m := len(s2); m < n {
		n = m
	}
	for i := 0; i < n; i++ {
		if s1[i] != s2[i] {
			return i
		}
	}
	return n
}

func (n *Set) Remove(s string) (old *Set) {
	return n.remove(s, func() {})
}

func (a *Set) remove(s string, deleteSelfFn func()) (old *Set) {
	if s == "" {
		if a.count > 0 {
			if a.count--; a.empty() {
				deleteSelfFn()
				old = a
			}
		}
		return old
	}
	if len(s) < len(a.children) {
		// Fast path: prefix substring search.
		for i := 1; i <= len(s); i++ {
			if child, ok := a.children[s[:i]]; ok {
				return child.remove(s[i:], func() { delete(a.children, s[:i]) })
			}
		}
	} else {
		// Slow path: children search.
		for prefix, child := range a.children {
			if n := prefixLen(prefix, s); n > 0 {
				return child.remove(s[n:], func() { delete(a.children, prefix) })
			}
		}
	}
	// Prefix not found.
	return nil
}

func (a *Set) debug(prefix string, depth int) {
	fmt.Printf("depth %d: %q: %d\n", depth, prefix, a.count)
	for childPrefix, child := range a.children {
		child.debug(prefix+childPrefix, depth+1)
	}
}
