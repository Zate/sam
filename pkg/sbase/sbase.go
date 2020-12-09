package sbase

// New returns a new empty Sbase
func New() *SBase {
	s := &SBase{
		Creds:    new(Creds),
		Manifest: new(Manifest),
		Catalog:  make([]*App, 0),
	}
	return s
}
