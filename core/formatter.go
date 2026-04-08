package core

type MessageFormatter[T any] interface {
	Format(entry *Entry) (T, error)
}

type JsonMessageFormatter MessageFormatter[map[string]any]
