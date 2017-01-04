package macross

import (
	"net/url"
)

type (
	//  RawStore is the interface that operates the session data.
	RawStore interface {
		// Set sets value to given key in session.
		Set(interface{}, interface{}) error
		// Get gets value by given key in session.
		Get(interface{}) interface{}
		// Delete deletes a key from session.
		Delete(interface{}) error
		// ID returns current session ID.
		ID() string
		// Release releases session resource and save data to provider.
		Release(*Context) error
		// Flush deletes all session data.
		Flush() error
	}

	// Sessioner is the interface that contains all data for one session process with specific ID.
	Sessioner interface {
		RawStore
		//---------------------------------------------//
		// Read returns raw session store by session ID.
		Read(string) (RawStore, error)
		// Destory deletes a session.
		Destory(*Context) error
		// RegenerateId regenerates a session store from old session ID to new one.
		RegenerateId(*Context) (RawStore, error)
		// Count counts and returns number of sessions.
		Count() int
		// GC calls GC to clean expired sessions.
		GC()
	}

	sessioner struct{}

	Flash struct {
		FlashNow bool
		Ctx      *Context
		url.Values
		ErrorMsg, WarningMsg, InfoMsg, SuccessMsg string
	}
)

// Set value to session
func (s *sessioner) Set(key, value interface{}) error { return nil }

// Get value from session by key
func (s *sessioner) Get(key interface{}) interface{} { return nil }

// Delete in session by key
func (s *sessioner) Delete(key interface{}) error { return nil }

// ID get this id of session store
func (s *sessioner) ID() string { return "" }

// Release Implement method, no used.
func (s *sessioner) Release(*Context) error { return nil }

// Flush clear all values in session
func (s *sessioner) Flush() error { return nil }

// Read returns raw session store by session ID.
func (s *sessioner) Read(string) (RawStore, error) { return nil, nil }

// Destory deletes a session.
func (s *sessioner) Destory(*Context) error { return nil }

// RegenerateId regenerates a session store from old session ID to new one.
func (s *sessioner) RegenerateId(*Context) (RawStore, error) { return nil, nil }

// Count counts and returns number of sessions.
func (s *sessioner) Count() int { return 0 }

// GC calls GC to clean expired sessions.
func (s *sessioner) GC() {}

// ___________.____       _____    _________ ___ ___
// \_   _____/|    |     /  _  \  /   _____//   |   \
//  |    __)  |    |    /  /_\  \ \_____  \/    ~    \
//  |     \   |    |___/    |    \/        \    Y    /
//  \___  /   |_______ \____|__  /_______  /\___|_  /
//      \/            \/       \/        \/       \/

func (f *Flash) set(name, msg string, current ...bool) {
	isShow := false
	if (len(current) == 0 && FlashNow) || (len(current) > 0 && current[0]) {
		isShow = true
	}

	if isShow {
		f.Ctx.Set("Flash", f)
	} else {
		f.Set(name, msg)
	}
}

func (f *Flash) Error(msg string, current ...bool) {
	f.ErrorMsg = msg
	f.set("error", msg, current...)
}

func (f *Flash) Warning(msg string, current ...bool) {
	f.WarningMsg = msg
	f.set("warning", msg, current...)
}

func (f *Flash) Info(msg string, current ...bool) {
	f.InfoMsg = msg
	f.set("info", msg, current...)
}

func (f *Flash) Success(msg string, current ...bool) {
	f.SuccessMsg = msg
	f.set("success", msg, current...)
}
