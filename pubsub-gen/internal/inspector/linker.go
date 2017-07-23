package inspector

type Linker struct{}

func NewLinker() Linker {
	return Linker{}
}

func (l Linker) Link(m map[string]Struct) {
	for n := range m {
		l.linkFields(n, m)
	}
}

func (l Linker) linkFields(n string, m map[string]Struct) {
	s := m[n]
	for i, f := range s.Fields {
		_, ok := m[f.Type]
		if !ok {
			continue
		}

		s.PeerTypeFields = append(s.PeerTypeFields, f)
		s.Fields = append(s.Fields[:i], s.Fields[i+1:]...)
		m[n] = s

		// We want to restart the loop because we've messed with our indexes
		l.linkFields(n, m)
		return
	}
}
