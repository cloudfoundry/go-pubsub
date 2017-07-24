package inspector

type Linker struct{}

func NewLinker() Linker {
	return Linker{}
}

func (l Linker) Link(m map[string]Struct, interfaceToStruct map[string][]string) {
	for n := range m {
		l.linkFields(n, m, interfaceToStruct)
	}
}

func (l Linker) linkFields(n string, m map[string]Struct, mi map[string][]string) {
	s := m[n]
	if s.InterfaceTypeFields == nil {
		s.InterfaceTypeFields = make(map[Field][]string)
	}

	for i, f := range s.Fields {
		_, ok := m[f.Type]
		if ok {
			s.PeerTypeFields = append(s.PeerTypeFields, f)
			s.Fields = append(s.Fields[:i], s.Fields[i+1:]...)
			m[n] = s

			// We want to restart the loop because we've messed with our indexes
			l.linkFields(n, m, mi)
			return
		}

		t, ok := mi[f.Type]
		if ok {
			s.InterfaceTypeFields[f] = append(s.InterfaceTypeFields[f], t...)
			s.Fields = append(s.Fields[:i], s.Fields[i+1:]...)
			m[n] = s
		}
	}
}
