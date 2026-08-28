package manifest

// Entry is one deployable asset.
type Entry struct {
	Path string // deployment-relative path, forward slashes, e.g. "assets/logo.svg"
	Hash string // Cloudflare Pages content hash (see hash.go)
}

// Manifest is an ordered, dual-indexed collection of Entries.
type Manifest struct {
	entries []Entry
	byPath  map[string]int
	byHash  map[string]int
}

func New() *Manifest {
	return &Manifest{
		byPath: make(map[string]int),
		byHash: make(map[string]int),
	}
}

func (m *Manifest) Add(path, hash string) {
	idx := len(m.entries)
	m.entries = append(m.entries, Entry{Path: path, Hash: hash})
	m.byPath[path] = idx
	m.byHash[hash] = idx
}

func (m *Manifest) Len() int                      { return len(m.entries) }
func (m *Manifest) All() []Entry                  { return m.entries } // range-friendly
func (m *Manifest) ByPath(p string) (Entry, bool) { e, ok := m.lookup(m.byPath, p); return e, ok }
func (m *Manifest) ByHash(h string) (Entry, bool) { e, ok := m.lookup(m.byHash, h); return e, ok }

func (m *Manifest) lookup(idx map[string]int, key string) (Entry, bool) {
	i, ok := idx[key]
	if !ok {
		return Entry{}, false
	}
	return m.entries[i], true
}

// Hashes returns every hash, for the check-missing call.
func (m *Manifest) Hashes() []string {
	out := make([]string, len(m.entries))
	for i, e := range m.entries {
		out[i] = e.Hash
	}
	return out
}

// Subset returns a new Manifest containing only the given hashes —
// exactly what you need after check-missing tells you what to upload.
func (m *Manifest) Subset(hashes []string) *Manifest {
	sub := New()
	for _, h := range hashes {
		if e, ok := m.ByHash(h); ok {
			sub.Add(e.Path, e.Hash)
		}
	}
	return sub
}

// APIManifest returns the {path: hash} map the "create deployment" call
// expects in its `manifest` field.
func (m *Manifest) APIManifest() map[string]string {
	out := make(map[string]string, len(m.entries))
	for _, e := range m.entries {
		out[e.Path] = e.Hash
	}
	return out
}
