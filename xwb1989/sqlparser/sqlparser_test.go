package sqlparser

import "testing"

func TestParseQuery(t *testing.T) {
	ParseQuery("SELECT name, email FROM `users` WHERE name = `john` ORDER BY name DESC")
}
