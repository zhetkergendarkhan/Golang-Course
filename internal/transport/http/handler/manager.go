package handler

type Manager struct {
	*ProductHandler
	*CategoryHandler
}

func NewManager(ph *ProductHandler, ch *CategoryHandler) *Manager {
	return &Manager{ProductHandler: ph, CategoryHandler: ch}
}
