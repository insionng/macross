package macross

type (
	// Sessioner is the interface that wraps the Session.
	Sessioner interface {
		Set(key, value interface{}) error //set session value
		Get(key interface{}) interface{}  //get session value
		Delete(key interface{}) error     //delete session value
		SessionID() string                //back current sessionID
		SessionRelease(ctx *Context)      // release the resource & save data to provider & return the data
		Flush() error                     //delete all data
	}

	sessioner struct{}
)

// Set value to session
func (s *sessioner) Set(key, value interface{}) error {
	return nil
}

// Get value from session by key
func (s *sessioner) Get(key interface{}) interface{} {
	return nil
}

// Delete in session by key
func (s *sessioner) Delete(key interface{}) error {
	return nil
}

// Flush clear all values in session
func (s *sessioner) Flush() error {
	return nil
}

// SessionID get this id of session store
func (s *sessioner) SessionID() string {
	return ""
}

// SessionRelease Implement method, no used.
func (s *sessioner) SessionRelease(ctx *Context) {
}
