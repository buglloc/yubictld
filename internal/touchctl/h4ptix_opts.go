package touchctl

type H4ptixOption interface {
	isH4ptixOption()
}
type optH4ptixSerial struct {
	H4ptixOption
	serial string
}

func H4ptixWithSerial(serial string) H4ptixOption {
	return optH4ptixSerial{
		serial: serial,
	}
}
