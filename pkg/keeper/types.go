package keeper

// window is what is created when new cmd is used
type window struct {
	Owner       string          // the user name who created the window
	Expire      int64           // the time after which submissions can be released
	Submissions map[string]user // the list of submissions 1 per user
}

// user information
type user struct {
	Id         string
	Name       string
	Submission string
}
