package persistence

import (
	"strconv"
	"sync/atomic"

	"github.com/shauera/messages/model"
)

type MemoryRepository struct {
	personIDCounter int64
	personsStorage map[string]model.Person
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		personsStorage: make(map[string]model.Person),
	}
}

func (mr *MemoryRepository) CreatePerson(person model.Person) (*string, error) {
	id := strconv.FormatInt(atomic.AddInt64(&mr.personIDCounter, 1), 10)
	newPerson := person
	newPerson.ID = id
	mr.personsStorage[id] = newPerson

	return &id, nil
}

func (mr *MemoryRepository) ListPersons() (model.Persons, error) {

	persons := make(model.Persons, 0, len(mr.personsStorage))

	for  _, value := range mr.personsStorage {
		persons = append(persons, value)
	}

	return persons, nil
}

func (mr *MemoryRepository) FindPersonById(id string) (*model.Person, error) {
	person := mr.personsStorage[id]
	return &person, nil
}
