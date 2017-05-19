### bufio.NewScanner Limitations
```
// Programs that need more control over error handling or large tokens,
// or must run sequential scans on a reader, should use bufio.Reader instead.
```

### TODO
- [ ] Analyze body of commits (added/removed lines)
- [ ] Add concurrency (parallelize requests to github API)
- [ ] Add context + timeout to requests to github API
