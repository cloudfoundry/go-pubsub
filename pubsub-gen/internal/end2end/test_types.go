package end2end

type X struct {
	I             int
	J             string
	Y1            Y
	Y2            *Y
	E1            Empty
	E2            *Empty
	M             message
	Repeated      []string
	RepeatedY     []Y
	RepeatedEmpty []Empty
	MapY          map[string]Y
}

type Y struct {
	I  int
	J  string
	E1 Empty
	E2 *Empty
}

type Empty struct{}

type message interface {
	message()
}

type M1 struct {
	A int
}

func (m M1) message() {}

type M2 struct {
	A int
	B int
}

func (m M2) message() {}

type M3 struct{}

func (m M3) message() {}
