package store

func NewMockStore(p PatientStorer, s SessionStorer, u UserStorer) *Store {
	return &Store{
		Patient: p,
		Session: s,
		User:    u,
	}
}
