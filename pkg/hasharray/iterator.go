package hasharray

type Iterator[T any] interface {
	Next() (T, bool)
}
