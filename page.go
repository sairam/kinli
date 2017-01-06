package kinli

//Page has all information to be fillable
type Page struct {
	// Title of the page
	Title string // TODO make this richer context
	// User information that can be accessed on any page
	User interface{}
	// List of Flash Messages
	Flashes []string // List of flash messages
	// Context is the Main Content of the page
	Context interface{}
	// Data not part of the context (one timers info)
	Data map[string]string
	// ClientConfig contains constant stuff like GoogleAnalytics attributes etc.,
	ClientConfig map[string]string
}

// hc.isAuthed() {
// userInfo = &Authentication{}
var (
	// ClientConfig can be set once so that all pages can use that static information
	ClientConfig map[string]string
)

// NewPage is used by all HTML contexts to display the template
// Emails do not use Pages. Only for views
func NewPage(hc *HttpContext, title string, user interface{}, ctx interface{}, data map[string]string) *Page {
	return &Page{
		Title:        title,
		User:         user,
		Flashes:      hc.GetFlashes(),
		Context:      ctx,
		Data:         data,
		ClientConfig: ClientConfig,
	}

}
