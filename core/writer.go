package core

type Writer interface {
	Write(entry *Entry) error
}
