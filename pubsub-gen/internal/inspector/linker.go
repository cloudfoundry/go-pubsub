package inspector

import (
	"log"
)

type Linker struct{}

func NewLinker() Linker {
	return Linker{}
}

func (l Linker) Link(m map[string]Struct, interfaceToStruct map[string][]string) {
	for n := range m {
		l.linkFields(n, m, interfaceToStruct)
		l.linkSliceTypes(n, m)
	}
}

func (l Linker) linkFields(n string, m map[string]Struct, mi map[string][]string) {
	s := m[n]
	if s.InterfaceTypeFields == nil {
		s.InterfaceTypeFields = make(map[Field][]string)
	}

	for i, f := range s.Fields {
		_, ok := m[f.Type]
		if ok && !f.Slice.IsSlice {
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

func (l Linker) linkSliceTypes(n string, m map[string]Struct) {
	s := m[n]

	for i, f := range s.Fields {
		if !f.Slice.IsSlice || f.Slice.IsBasicType {
			continue
		}

		peer, ok := m[f.Type]
		if !ok {
			log.Fatal("Unknown type for slice: %s", f.Type)
		}

		ff, ok := l.findFieldName(f.Slice.FieldName, peer.Fields)
		if !ok {
			log.Fatal("Unknown field name for slice: %s %s", f.Type, f.Slice.FieldName)
		}

		f.Type = ff.Type
		s.Fields[i] = f
	}
}

func (l Linker) findFieldName(name string, fields []Field) (Field, bool) {
	for _, ff := range fields {
		if ff.Name == name {
			return ff, true
		}
	}
	return Field{}, false
}
