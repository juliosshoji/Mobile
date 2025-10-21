package model

type Repository[T any] interface {
	Get(string) (T, error)
	Save(T) error
	Update(T) error
	Delete(string) error
	GetAll() *[]T
}