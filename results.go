package gorethink

import (
	"github.com/dancannon/gorethink/encoding"
	p "github.com/dancannon/gorethink/ql2"
)

type ResultRow struct {
	err  error
	rows *ResultRows
}

// Scan copies the result from the matched row into the value pointed at by dest.
// If more than one row is returned by the query then Scan returns the first and
// ignores the rest. If no row is found then Scan returns an error.
func (r *ResultRow) Scan(dest interface{}) error {
	if r.err != nil {
		return r.err
	}

	defer r.rows.Close()
	if !r.rows.Next() {
		return RqlDriverError{"No rows in the result set"}
	}
	err := r.rows.Scan(dest)
	if err != nil {
		return err
	}

	return nil
}

// ResultRows contains the result of a query. Its cursor starts before the first row
// of the result set. Use Next to advance through the rows.
type ResultRows struct {
	session      *Session
	query        *p.Query
	term         RqlTerm
	opts         map[string]interface{}
	buffer       []*p.Datum
	current      *p.Datum
	start        int
	end          int
	token        int64
	err          error
	closed       bool
	responseType p.Response_ResponseType
}

// Close closes the Rows, preventing further enumeration. If the end is
// encountered, the Rows are closed automatically. Close is idempotent.
func (r *ResultRows) Close() error {
	var err error

	if !r.closed {
		_, err = r.session.stopQuery(r.query, r.term, r.opts)
		r.closed = true
	}

	return err
}

// Err returns the error, if any, that was encountered during iteration.
func (r *ResultRows) Err() error {
	return r.err
}

// Next prepares the next row for reading. It returns true on success or false
// if there are no more rows. Every call to scan must be preceeded by a call
// to next. If all rows in the buffer have been read and a partial sequence was
// returned then Next will load more from the database
func (r *ResultRows) Next() bool {
	if r.closed {
		return false
	}

	if r.err != nil {
		return false
	}

	// Attempt to get a result in the buffer
	if r.end > r.start {
		row := r.buffer[r.start]

		if !r.advance() {
			return false
		}

		r.current = row
		if row != nil {
			return true
		}
	}

	// Check if all rows have been loaded
	if r.responseType == p.Response_SUCCESS_SEQUENCE {
		r.closed = true
		r.start = 0
		r.end = 0
		return false
	}

	// Load more data from the database

	// First, shift data to beginning of buffer if there's lots of empty space
	// or space is neded.
	if r.start > 0 && (r.end == len(r.buffer) || r.start > len(r.buffer)/2) {
		copy(r.buffer, r.buffer[r.start:r.end])
		r.end -= r.start
		r.start = 0
	}

	// Continue the query
	newResult, err := r.session.continueQuery(r.query, r.term, r.opts)
	if err != nil {
		r.err = err
		return false
	}

	r.buffer = append(r.buffer, newResult.buffer...)
	r.end += len(newResult.buffer)

	r.advance()
	r.current = r.buffer[r.start]

	return true
}

// advance moves the internal buffer pointer ahead to point to the next row
func (r *ResultRows) advance() bool {
	if r.end <= r.start {
		return false
	}

	r.start++
	return true
}

// Scan copies the result in the current row into the value pointed at by dest.
//
// If an argument as type *interface{}, Scan copies the value provided by the
// database without conversion.
//
// If the value is a struct then Scan traverses
// the result recursively and attempts to match the keys returned by the database
// to the name used by the structs field (either the struct field name or its
// key).
func (r *ResultRows) Scan(dest interface{}) error {
	if r.err != nil {
		return r.err
	}
	if r.current == nil {
		return RqlDriverError{"Scan called without calling Next"}
	}

	data, err := deconstructDatum(r.current, r.opts)
	if err != nil {
		return err
	}

	err = encoding.Decode(dest, data)
	if err != nil {
		return err
	}

	return nil
}

// All is a helper method for returning a slice containing all rows. The slice
// is of type []interface{}.
func (r *ResultRows) All() ([]interface{}, error) {
	rows := []interface{}{}

	for r.Next() {
		var row interface{}

		err := r.Scan(&row)
		if err != nil {
			return rows, err
		}

		rows = append(rows, row)
	}

	return rows, nil
}
