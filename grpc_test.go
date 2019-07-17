package main

import "github.com/stretchr/testify/suite"

type TestSuiteGrpc struct {
	suite.Suite
	db Accessor
}

func (suite *TestSuiteGrpc) SetupSuite() {

	if true {
		db, err := createSQLDB(DIALECT, "testdb.db")
		if err != nil {
			panic(err)
		}
		suite.db = &db
	} else {
		db := createInMemoryDB()
		suite.db = &db
	}
}

func (suite *TestSuiteSQL) TestCreateGRPC() {

}
