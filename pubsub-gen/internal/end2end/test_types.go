package end2end

type X struct {
	I  int
	J  string
	Y1 Y
	Y2 *Y
	M  message
}

type Y struct {
	I int
	J string
}

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
