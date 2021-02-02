package sbase

// New returns a new empty Sbase
func New() *SBase {
	s := &SBase{
		Creds:    &Creds{},
		Manifest: &Manifest{},
		Catalog:  new(map[string]*App),
	}
	return s
}
