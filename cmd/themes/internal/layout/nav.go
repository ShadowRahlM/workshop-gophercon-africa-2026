package layout

// Link is one entry in a themed page's navigation.
type Link struct {
	Path  string
	Title string
}

// Nav is the navigation every themed page shows, in order.
var Nav = []Link{
	{Path: "/antarctica", Title: "Antarctica"},
	{Path: "/kenya", Title: "Kenya"},
	{Path: "/humpback", Title: "Humpback"},
	{Path: "/uganda", Title: "Uganda"},
}

// Current returns the aria-current value for this link on the page at
// the given path. "false" is a valid value meaning "not the current
// page", so navigation markup needs no branch around it.
func (l Link) Current(path string) string {
	if l.Path == path {
		return "page"
	}
	return "false"
}
