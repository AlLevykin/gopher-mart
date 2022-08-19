package ports

type AccrualDispatcher interface {
	Start()
	Dispatch(order string)
	Stop()
}
