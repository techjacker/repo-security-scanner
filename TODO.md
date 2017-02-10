- [ ] Add middleware - validate = from github origin + tests
	- use gorilla middleware vs write own


- [ ] type []Stats == logger = strings.Stringer() interface -> for creating string for email
- [ ] Add email notifications (+ interface + tests)


- [ ] Enable analysis of private github repos (authenticate using integration ID + private key - add to secrets)


- [ ] Analyze body of commits (added/removed lines)
- [ ] Add concurrency (parallelize requests to github API)
- [ ] Add context + timeout to requests to github API
