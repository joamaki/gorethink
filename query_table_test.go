package gorethink

import (
	test "launchpad.net/gocheck"
)

func (s *RethinkSuite) TestTableCreate(c *test.C) {
	var response interface{}

	Db("test").TableDrop("test").Exec(sess)

	// Test database creation
	query := Db("test").TableCreate("test")

	err := query.RunRow(sess).Scan(&response)

	c.Assert(err, test.IsNil)
	c.Assert(response, JsonEquals, map[string]interface{}{"created": 1})
}

func (s *RethinkSuite) TestTableCreatePrimaryKey(c *test.C) {
	var response interface{}

	Db("test").TableDrop("testOpts").Exec(sess)

	// Test database creation
	query := Db("test").TableCreate("testOpts",
		"primary_key", "it",
	)

	err := query.RunRow(sess).Scan(&response)

	c.Assert(err, test.IsNil)
	c.Assert(response, JsonEquals, map[string]interface{}{"created": 1})
}

func (s *RethinkSuite) TestTableCreateSoftDurability(c *test.C) {
	var response interface{}

	Db("test").TableDrop("testOpts").Exec(sess)

	// Test database creation
	query := Db("test").TableCreate("testOpts",
		"durability", "soft",
	)

	err := query.RunRow(sess).Scan(&response)

	c.Assert(err, test.IsNil)
	c.Assert(response, JsonEquals, map[string]interface{}{"created": 1})
}

func (s *RethinkSuite) TestTableCreateSoftMultipleOpts(c *test.C) {
	var response interface{}

	Db("test").TableDrop("testOpts").Exec(sess)

	// Test database creation
	query := Db("test").TableCreate("testOpts",
		"primary_key", "it",
		"durability", "soft",
	)

	err := query.RunRow(sess).Scan(&response)

	c.Assert(err, test.IsNil)
	c.Assert(response, JsonEquals, map[string]interface{}{"created": 1})

	Db("test").TableDrop("test").Exec(sess)
}

func (s *RethinkSuite) TestTableList(c *test.C) {
	var response interface{}

	Db("test").TableCreate("test").Exec(sess)

	// Try and find it in the list
	success := false
	err := Db("test").TableList().RunRow(sess).Scan(&response)

	c.Assert(err, test.IsNil)
	c.Assert(response, test.FitsTypeOf, []interface{}{})

	for _, db := range response.([]interface{}) {
		if db == "test" {
			success = true
		}
	}

	c.Assert(success, test.Equals, true)
}

func (s *RethinkSuite) TestTableDelete(c *test.C) {
	var response interface{}

	Db("test").TableCreate("test").Exec(sess)

	// Test database creation
	query := Db("test").TableDrop("test")

	err := query.RunRow(sess).Scan(&response)

	c.Assert(err, test.IsNil)
	c.Assert(response, JsonEquals, map[string]interface{}{"dropped": 1})
}

func (s *RethinkSuite) TestTableIndexCreate(c *test.C) {
	var response interface{}

	Db("test").TableCreate("test").Exec(sess)
	Db("test").Table("test").IndexDrop("test").Exec(sess)

	// Test database creation
	query := Db("test").Table("test").IndexCreate("test")

	err := query.RunRow(sess).Scan(&response)

	c.Assert(err, test.IsNil)
	c.Assert(response, JsonEquals, map[string]interface{}{"created": 1})
}

func (s *RethinkSuite) TestTableIndexList(c *test.C) {
	var response interface{}

	Db("test").TableCreate("test").Exec(sess)
	Db("test").Table("test").IndexCreate("test").Exec(sess)

	// Try and find it in the list
	success := false
	err := Db("test").Table("test").IndexList().RunRow(sess).Scan(&response)

	c.Assert(err, test.IsNil)
	c.Assert(response, test.FitsTypeOf, []interface{}{})

	for _, db := range response.([]interface{}) {
		if db == "test" {
			success = true
		}
	}

	c.Assert(success, test.Equals, true)
}

func (s *RethinkSuite) TestTableIndexDelete(c *test.C) {
	var response interface{}

	Db("test").TableCreate("test").Exec(sess)
	Db("test").Table("test").IndexCreate("test").Exec(sess)

	// Test database creation
	query := Db("test").Table("test").IndexDrop("test")

	err := query.RunRow(sess).Scan(&response)

	c.Assert(err, test.IsNil)
	c.Assert(response, JsonEquals, map[string]interface{}{"dropped": 1})
}
