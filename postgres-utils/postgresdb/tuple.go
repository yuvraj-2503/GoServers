package postgresdb

import "encoding/json"

type Tuple struct {
	params []interface{}
}

func NewTuple() *Tuple {
	return &Tuple{params: []interface{}{}}
}

// AddString adds a string value to the tuple.
func (t *Tuple) AddString(value string) {
	t.params = append(t.params, value)
}

// AddInt adds an integer value to the tuple.
func (t *Tuple) AddInt(value int) {
	t.params = append(t.params, value)
}

// AddInt64 adds a 64-bit integer value to the tuple.
func (t *Tuple) AddInt64(value int64) {
	t.params = append(t.params, value)
}

// AddFloat adds a float value to the tuple.
func (t *Tuple) AddFloat(value float64) {
	t.params = append(t.params, value)
}

// AddBool adds a boolean value to the tuple.
func (t *Tuple) AddBool(value bool) {
	t.params = append(t.params, value)
}

// AddValue adds a generic value to the tuple.
func (t *Tuple) AddValue(value interface{}) {
	t.params = append(t.params, value)
}

func (t *Tuple) AddObject(value interface{}) {
	jsonValue, _ := json.Marshal(value)
	t.params = append(t.params, string(jsonValue))
}
