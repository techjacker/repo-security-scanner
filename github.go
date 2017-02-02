package main

type githubResponseFull struct {
	Body struct {
		After   string      `json:"after"`
		BaseRef interface{} `json:"base_ref"`
		Before  string      `json:"before"`
		Commits []struct {
			Added  []string `json:"added"`
			Author struct {
				Email    string `json:"email"`
				Name     string `json:"name"`
				Username string `json:"username"`
			} `json:"author"`
			Committer struct {
				Email    string `json:"email"`
				Name     string `json:"name"`
				Username string `json:"username"`
			} `json:"committer"`
			Distinct  bool          `json:"distinct"`
			ID        string        `json:"id"`
			Message   string        `json:"message"`
			Modified  []interface{} `json:"modified"`
			Removed   []interface{} `json:"removed"`
			Timestamp string        `json:"timestamp"`
			TreeID    string        `json:"tree_id"`
			URL       string        `json:"url"`
		} `json:"commits"`
		Compare    string `json:"compare"`
		Created    bool   `json:"created"`
		Deleted    bool   `json:"deleted"`
		Forced     bool   `json:"forced"`
		HeadCommit struct {
			Added  []string `json:"added"`
			Author struct {
				Email    string `json:"email"`
				Name     string `json:"name"`
				Username string `json:"username"`
			} `json:"author"`
			Committer struct {
				Email    string `json:"email"`
				Name     string `json:"name"`
				Username string `json:"username"`
			} `json:"committer"`
			Distinct  bool          `json:"distinct"`
			ID        string        `json:"id"`
			Message   string        `json:"message"`
			Modified  []interface{} `json:"modified"`
			Removed   []interface{} `json:"removed"`
			Timestamp string        `json:"timestamp"`
			TreeID    string        `json:"tree_id"`
			URL       string        `json:"url"`
		} `json:"head_commit"`
		Installation struct {
			ID int64 `json:"id"`
		} `json:"installation"`
		Organization struct {
			AvatarURL        string      `json:"avatar_url"`
			Description      interface{} `json:"description"`
			EventsURL        string      `json:"events_url"`
			HooksURL         string      `json:"hooks_url"`
			ID               int64       `json:"id"`
			IssuesURL        string      `json:"issues_url"`
			Login            string      `json:"login"`
			MembersURL       string      `json:"members_url"`
			PublicMembersURL string      `json:"public_members_url"`
			ReposURL         string      `json:"repos_url"`
			URL              string      `json:"url"`
		} `json:"organization"`
		Pusher struct {
			Email string `json:"email"`
			Name  string `json:"name"`
		} `json:"pusher"`
		Ref        string `json:"ref"`
		Repository struct {
			ArchiveURL       string      `json:"archive_url"`
			AssigneesURL     string      `json:"assignees_url"`
			BlobsURL         string      `json:"blobs_url"`
			BranchesURL      string      `json:"branches_url"`
			CloneURL         string      `json:"clone_url"`
			CollaboratorsURL string      `json:"collaborators_url"`
			CommentsURL      string      `json:"comments_url"`
			CommitsURL       string      `json:"commits_url"`
			CompareURL       string      `json:"compare_url"`
			ContentsURL      string      `json:"contents_url"`
			ContributorsURL  string      `json:"contributors_url"`
			CreatedAt        int64       `json:"created_at"`
			DefaultBranch    string      `json:"default_branch"`
			DeploymentsURL   string      `json:"deployments_url"`
			Description      interface{} `json:"description"`
			DownloadsURL     string      `json:"downloads_url"`
			EventsURL        string      `json:"events_url"`
			Fork             bool        `json:"fork"`
			Forks            int64       `json:"forks"`
			ForksCount       int64       `json:"forks_count"`
			ForksURL         string      `json:"forks_url"`
			FullName         string      `json:"full_name"`
			GitCommitsURL    string      `json:"git_commits_url"`
			GitRefsURL       string      `json:"git_refs_url"`
			GitTagsURL       string      `json:"git_tags_url"`
			GitURL           string      `json:"git_url"`
			HasDownloads     bool        `json:"has_downloads"`
			HasIssues        bool        `json:"has_issues"`
			HasPages         bool        `json:"has_pages"`
			HasWiki          bool        `json:"has_wiki"`
			Homepage         interface{} `json:"homepage"`
			HooksURL         string      `json:"hooks_url"`
			HTMLURL          string      `json:"html_url"`
			ID               int64       `json:"id"`
			IssueCommentURL  string      `json:"issue_comment_url"`
			IssueEventsURL   string      `json:"issue_events_url"`
			IssuesURL        string      `json:"issues_url"`
			KeysURL          string      `json:"keys_url"`
			LabelsURL        string      `json:"labels_url"`
			Language         interface{} `json:"language"`
			LanguagesURL     string      `json:"languages_url"`
			MasterBranch     string      `json:"master_branch"`
			MergesURL        string      `json:"merges_url"`
			MilestonesURL    string      `json:"milestones_url"`
			MirrorURL        interface{} `json:"mirror_url"`
			Name             string      `json:"name"`
			NotificationsURL string      `json:"notifications_url"`
			OpenIssues       int64       `json:"open_issues"`
			OpenIssuesCount  int64       `json:"open_issues_count"`
			Organization     string      `json:"organization"`
			Owner            struct {
				Email interface{} `json:"email"`
				Name  string      `json:"name"`
			} `json:"owner"`
			Private         bool   `json:"private"`
			PullsURL        string `json:"pulls_url"`
			PushedAt        int64  `json:"pushed_at"`
			ReleasesURL     string `json:"releases_url"`
			Size            int64  `json:"size"`
			SSHURL          string `json:"ssh_url"`
			Stargazers      int64  `json:"stargazers"`
			StargazersCount int64  `json:"stargazers_count"`
			StargazersURL   string `json:"stargazers_url"`
			StatusesURL     string `json:"statuses_url"`
			SubscribersURL  string `json:"subscribers_url"`
			SubscriptionURL string `json:"subscription_url"`
			SvnURL          string `json:"svn_url"`
			TagsURL         string `json:"tags_url"`
			TeamsURL        string `json:"teams_url"`
			TreesURL        string `json:"trees_url"`
			UpdatedAt       string `json:"updated_at"`
			URL             string `json:"url"`
			Watchers        int64  `json:"watchers"`
			WatchersCount   int64  `json:"watchers_count"`
		} `json:"repository"`
		Sender struct {
			AvatarURL         string `json:"avatar_url"`
			EventsURL         string `json:"events_url"`
			FollowersURL      string `json:"followers_url"`
			FollowingURL      string `json:"following_url"`
			GistsURL          string `json:"gists_url"`
			GravatarID        string `json:"gravatar_id"`
			HTMLURL           string `json:"html_url"`
			ID                int64  `json:"id"`
			Login             string `json:"login"`
			OrganizationsURL  string `json:"organizations_url"`
			ReceivedEventsURL string `json:"received_events_url"`
			ReposURL          string `json:"repos_url"`
			SiteAdmin         bool   `json:"site_admin"`
			StarredURL        string `json:"starred_url"`
			SubscriptionsURL  string `json:"subscriptions_url"`
			Type              string `json:"type"`
			URL               string `json:"url"`
		} `json:"sender"`
	} `json:"body"`
	Headers struct {
		Accept            string `json:"accept"`
		Content_length    string `json:"content-length"`
		Content_type      string `json:"content-type"`
		Host              string `json:"host"`
		User_agent        string `json:"user-agent"`
		X_forwarded_for   string `json:"x-forwarded-for"`
		X_github_delivery string `json:"x-github-delivery"`
		X_github_event    string `json:"x-github-event"`
	} `json:"headers"`
}
