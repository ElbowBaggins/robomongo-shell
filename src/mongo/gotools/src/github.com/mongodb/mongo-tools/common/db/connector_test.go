// Copyright (C) MongoDB, Inc. 2014-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package db_test

import (
	"testing"
	"time"

	"github.com/mongodb/mongo-tools/common/db"
	"github.com/mongodb/mongo-tools/common/options"
	"github.com/mongodb/mongo-tools/common/testtype"
	"github.com/mongodb/mongo-tools/common/testutil"
	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/mgo.v2"
)

func TestVanillaDBConnector(t *testing.T) {

	testtype.VerifyTestType(t, "db")

	Convey("With a vanilla db connector", t, func() {

		var connector *db.VanillaDBConnector

		Convey("calling Configure should populate the addrs and dial timeout"+
			" appropriately with no error", func() {

			connector = &db.VanillaDBConnector{}

			opts := options.ToolOptions{
				Connection: &options.Connection{
					Host: "host1,host2",
					Port: "20000",
				},
				Auth: &options.Auth{},
			}
			So(connector.Configure(opts), ShouldBeNil)
			info := db.VanillaTestWrapper(*connector).DialInfo()
			So(info.Addrs, ShouldResemble,
				[]string{"host1:20000", "host2:20000"})
			So(info.Timeout, ShouldResemble, time.Duration(opts.Timeout)*time.Second)

		})

		Convey("calling GetNewSession with a running mongod should connect"+
			" successfully", func() {

			connector = &db.VanillaDBConnector{}

			opts := options.ToolOptions{
				Connection: &options.Connection{
					Host: "localhost",
					Port: db.DefaultTestPort,
				},
				Auth: &options.Auth{},
			}
			So(connector.Configure(opts), ShouldBeNil)

			session, err := connector.GetNewSession()
			So(err, ShouldBeNil)
			So(session, ShouldNotBeNil)
			session.Close()

		})

	})

}

func TestVanillaDBConnectorWithAuth(t *testing.T) {
	testtype.VerifyTestType(t, "auth")
	session, err := mgo.Dial("localhost:33333")
	if err != nil {
		t.Fatalf("error dialing server: %v", err)
	}

	err = testutil.CreateUserAdmin(session)
	So(err, ShouldBeNil)
	err = testutil.CreateUserWithRole(session, "cAdmin", "password",
		mgo.RoleClusterAdmin, true)
	So(err, ShouldBeNil)
	session.Close()

	Convey("With a vanilla db connector and a mongod running with"+
		" auth", t, func() {

		var connector *db.VanillaDBConnector

		Convey("connecting without authentication should not be able"+
			" to run commands", func() {

			connector = &db.VanillaDBConnector{}

			opts := options.ToolOptions{
				Connection: &options.Connection{
					Host: "localhost",
					Port: db.DefaultTestPort,
				},
				Auth: &options.Auth{},
			}
			So(connector.Configure(opts), ShouldBeNil)

			session, err := connector.GetNewSession()
			So(err, ShouldBeNil)
			So(session, ShouldNotBeNil)

			So(session.DB("admin").Run("top", &struct{}{}), ShouldNotBeNil)
			session.Close()

		})

		Convey("connecting with authentication should succeed and"+
			" authenticate properly", func() {

			connector = &db.VanillaDBConnector{}

			opts := options.ToolOptions{
				Connection: &options.Connection{
					Host: "localhost",
					Port: db.DefaultTestPort,
				},
				Auth: &options.Auth{
					Username: "cAdmin",
					Password: "password",
				},
			}
			So(connector.Configure(opts), ShouldBeNil)

			session, err := connector.GetNewSession()
			So(err, ShouldBeNil)
			So(session, ShouldNotBeNil)

			So(session.DB("admin").Run("top", &struct{}{}), ShouldBeNil)
			session.Close()

		})

	})

}
