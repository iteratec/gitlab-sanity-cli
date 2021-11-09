package gitlabsanitycli

// internal container for communicating data between channels
type content struct {
	ID           int
	Req          string
	Description  string
	Name         string
	Active       bool
	Status       string
	IPAddress    string
	Online       bool
	IsShared     bool
	GroupName    string
	GroupID      int
	LastActivity float64
}
