package ports

type AccrualDispatcher interface {
	Dispatch(order string)
}
