package gossie

// Mapping maps type of a Go object to/from (a slice of) a Cassandra row.
type Mapping interface {

	// Cf returns the column family name
	Cf() string

	// MinColumns returns the minimal number of columns required by the mapped Go object
	MinColumns() int

	// Map converts a Go object compatible with this Mapping into a Row
	Map(source interface{}) (*Row, error)

	// Ummap fills the passed Go object with data from the row, staring at the offset column. It
	// returns the count of the read columns
	Unmap(destination interface{}, offset int, row *Row) (int, error)
}
