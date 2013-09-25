package mockgossie

import (
	"fmt"
	"github.com/wadey/gossie/src/gossie"
)

type ExampleData struct {
	ID    string `cf:"example" key:"ID"`
	Value string `name:"value"`
}

func Example() {
	mock := NewMockConnectionPool()

	data := &ExampleData{
		ID:    "foo",
		Value: "bar",
	}

	err := MyCassandraSave(mock, data)
	if err != nil {
		panic(err)
	}

	data, err = MyCassandraLoad(mock, "foo")
	if err != nil {
		panic(err)
	}

	fmt.Println(data.Value)
	// Output: bar
}

func MyCassandraSave(cp gossie.ConnectionPool, data *ExampleData) error {
	mapping, _ := gossie.NewMapping(&ExampleData{})

	b := cp.Batch()
	b.Insert(mapping, data)
	return b.Run()
}

func MyCassandraLoad(cp gossie.ConnectionPool, id string) (*ExampleData, error) {
	mapping, _ := gossie.NewMapping(&ExampleData{})

	query := cp.Query(mapping)
	result, err := query.Get(id)
	if err != nil {
		return nil, err
	}
	data := &ExampleData{}
	err = result.Next(data)
	return data, err
}
