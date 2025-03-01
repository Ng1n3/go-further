package mysql

import (
	"database/sql"

	"github.com/Ng1n3/go-further/pkg/models"
)

// Define a SnippeModel type which wraps a sql.DB connectin pool.
type SnippetModel struct {
	DB *sql.DB
}

// This will insert a  new snippet into the database
func (m *SnippetModel) Insert(title, content, expires string) (int, error) {
	stmt := `INSERT INTO snippets (title, content, created_at, expires) VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

	result, err := m.DB.Exec(stmt, title, content, expires)
	if err != nil {
		return 0, nil
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

// This will return a specific snippet based on its id.
func (m *SnippetModel) Get(id int) (*models.Snippet, error) {
	// stmt := `SELECT id, title, content, created_at, expires FROM snippets WHERE expires > UTC_TIMESTAMP() AND ID = ?`

	// Use the QueryRow() method on the connection pool to execute our
	// SQL statement, passing in the untrusted id variable as the value for the
	// placeholder parameter. This returns a pointer to a sql.Row object which
	// holds the result from the database.
	// row := m.DB.QueryRow(stmt, id)
	// err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created_at, &s.Expires)

	// Initialize a pointe rto new Zeroed Snippet struct.
	s := &models.Snippet{}

	// Use row.Scan() to copy the values from each field in sql.Row to the
	// corresponding field in the Snippet struct. Notice that the arguments
	// to row.Scan are *pointers* to the place you want to copy the data into,
	// and the number of arguments must be exactly the same as the number of
	// columns returned by your statement. If the query returns no rows, then
	// row.Scan() will return a sql.ErrNoRows error. We check for that and retu
	// our own models.ErrNoRecord error instead of a Snippet object
	// err := m.DB.QueryRow("SELECT ...", id).Scan(&s.ID, &s.Title, &s.Content, &s.Created_at, &s.Expires)
	err := m.DB.QueryRow("SELECT id, title, content, created_at, expires FROM snippets WHERE expires > UTC_TIMESTAMP() AND id = ?", id).Scan(&s.ID, &s.Title, &s.Content, &s.Created_at, &s.Expires)
	if err == sql.ErrNoRows {
		return nil, models.ErrNoRecord
	} else if err != nil {
		return nil, err
	}
	return s, nil
}

/*
 You might be wondering why we’re returning the
models.ErrNoRecord error instead of sql.ErrNoRows directly. The
reason is to help encapsulate the model completely, so that our
application isn’t concerned with the underlying datastore or reliant on
datastore-specific errors for its behavior.
*/

// This will return the 10 most recently created snippets.
func (m *SnippetModel) Latest() ([]*models.Snippet, error) {

	stmt := `SELECT id, title, content, created_at, expires FROM snippets
  WHERE expires > UTC_TIMESTAMP() ORDER BY created_at DESC LIMIT 10`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	snippets := []*models.Snippet{}

	for rows.Next() {
		s := &models.Snippet{}

		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created_at, &s.Expires)
		if err != nil {
			return nil, err
		}
		// Append it to the slice of snippets.
		snippets = append(snippets, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil
}
