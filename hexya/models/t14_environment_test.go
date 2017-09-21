// Copyright 2016 NDP Systèmes. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package models

import (
	"testing"

	"github.com/hexya-erp/hexya/hexya/models/security"
	"github.com/hexya-erp/hexya/hexya/models/types"
	. "github.com/smartystreets/goconvey/convey"
)

func TestEnvironment(t *testing.T) {
	Convey("Testing Environment Modifications", t, func() {
		SimulateInNewEnvironment(security.SuperUserID, func(env Environment) {
			env.context = types.NewContext().WithKey("key", "context value")
			users := env.Pool("User")
			userJane := users.Search(users.Model().Field("Email").Equals("jane.smith@example.com"))
			Convey("Checking WithEnv", func() {
				env2 := newEnvironment(2)
				userJane1 := userJane.Call("WithEnv", env2).(RecordCollection)
				So(userJane1.Env().Uid(), ShouldEqual, 2)
				So(userJane.Env().Uid(), ShouldEqual, 1)
				So(userJane.Env().Context().HasKey("key"), ShouldBeTrue)
				So(userJane1.Env().Context().IsEmpty(), ShouldBeTrue)
				So(userJane.Env().callStack, ShouldBeEmpty)
				So(userJane1.Env().callStack, ShouldBeEmpty)
				env2.rollback()
			})
			Convey("Checking WithContext", func() {
				userJane1 := userJane.Call("WithContext", "newKey", "This is a different key").(RecordCollection)
				So(userJane1.Env().Context().HasKey("key"), ShouldBeTrue)
				So(userJane1.Env().Context().HasKey("newKey"), ShouldBeTrue)
				So(userJane1.Env().Context().Get("key"), ShouldEqual, "context value")
				So(userJane1.Env().Context().Get("newKey"), ShouldEqual, "This is a different key")
				So(userJane1.Env().Uid(), ShouldEqual, security.SuperUserID)
				So(userJane1.Env().callStack, ShouldBeEmpty)
				So(userJane.Env().Context().HasKey("key"), ShouldBeTrue)
				So(userJane.Env().Context().HasKey("newKey"), ShouldBeFalse)
				So(userJane.Env().Context().Get("key"), ShouldEqual, "context value")
				So(userJane.Env().Uid(), ShouldEqual, security.SuperUserID)
				So(userJane.Env().callStack, ShouldBeEmpty)
			})
			Convey("Checking WithNewContext", func() {
				newCtx := types.NewContext().WithKey("newKey", "This is a different key")
				userJane1 := userJane.Call("WithNewContext", newCtx).(RecordCollection)
				So(userJane1.Env().Context().HasKey("key"), ShouldBeFalse)
				So(userJane1.Env().Context().HasKey("newKey"), ShouldBeTrue)
				So(userJane1.Env().Context().Get("newKey"), ShouldEqual, "This is a different key")
				So(userJane1.Env().Uid(), ShouldEqual, security.SuperUserID)
				So(userJane1.Env().callStack, ShouldBeEmpty)
				So(userJane.Env().Context().HasKey("key"), ShouldBeTrue)
				So(userJane.Env().Context().HasKey("newKey"), ShouldBeFalse)
				So(userJane.Env().Context().Get("key"), ShouldEqual, "context value")
				So(userJane.Env().Uid(), ShouldEqual, security.SuperUserID)
				So(userJane.Env().callStack, ShouldBeEmpty)
			})
			Convey("Checking Sudo", func() {
				userJane1 := userJane.Sudo(2)
				userJane2 := userJane1.Call("Sudo").(RecordCollection)
				So(userJane1.Env().Uid(), ShouldEqual, 2)
				So(userJane1.Env().callStack, ShouldBeEmpty)
				So(userJane.Env().Uid(), ShouldEqual, security.SuperUserID)
				So(userJane.Env().callStack, ShouldBeEmpty)
				So(userJane2.Env().Uid(), ShouldEqual, security.SuperUserID)
				So(userJane2.Env().callStack, ShouldBeEmpty)
			})
			Convey("Checking combined modifications", func() {
				userJane1 := userJane.Sudo(2)
				userJane2 := userJane1.Sudo()
				userJane = userJane.WithContext("key", "modified value")
				So(userJane.Env().Context().Get("key"), ShouldEqual, "modified value")
				So(userJane1.Env().Context().Get("key"), ShouldEqual, "context value")
				So(userJane1.Env().Uid(), ShouldEqual, 2)
				So(userJane2.Env().Context().Get("key"), ShouldEqual, "context value")
				So(userJane2.Env().Uid(), ShouldEqual, security.SuperUserID)
			})
			Convey("Checking overridden WithContext", func() {
				posts := env.Pool("Post").FetchAll()
				posts1 := posts.WithContext("foo", "bar")
				So(posts1.Env().Context().HasKey("foo"), ShouldBeTrue)
				So(posts1.Env().Context().GetString("foo"), ShouldEqual, "bar")
				So(posts1.Env().callStack, ShouldBeEmpty)
				So(posts.Env().Context().HasKey("foo"), ShouldBeFalse)
				So(posts.Env().callStack, ShouldBeEmpty)
			})
		})
	})
	Convey("Testing cache operation", t, func() {
		SimulateInNewEnvironment(security.SuperUserID, func(env Environment) {
			users := env.Pool("User")
			userJane := users.Search(users.Model().Field("Email").Equals("jane.smith@example.com"))
			Convey("Cache should be empty at startup", func() {
				So(env.cache.data, ShouldBeEmpty)
				So(env.cache.m2mLinks, ShouldBeEmpty)
			})
			Convey("Loading a RecordSet should populate the cache", func() {
				userJane = userJane.Load()
				So(env.cache.m2mLinks, ShouldBeEmpty)
				So(env.cache.data, ShouldHaveLength, 1)
				janeCacheRef := cacheRef{model: users.model, id: userJane.ids[0]}
				So(env.cache.data, ShouldContainKey, janeCacheRef)
				So(env.cache.data[janeCacheRef], ShouldContainKey, "id")
				So(env.cache.data[janeCacheRef]["id"], ShouldEqual, userJane.ids[0])
				So(env.cache.data[janeCacheRef], ShouldContainKey, "name")
				So(env.cache.data[janeCacheRef]["name"], ShouldEqual, "Jane A. Smith")
				So(env.cache.data[janeCacheRef], ShouldContainKey, "email")
				So(env.cache.data[janeCacheRef]["email"], ShouldEqual, "jane.smith@example.com")
				So(env.cache.checkIfInCache(users.model, userJane.ids, []string{"id", "name", "email"}), ShouldBeTrue)
			})
			Convey("Calling values already in cache should not call the DB", func() {
				userJane = userJane.Load()
				id, dbCalled := userJane.get("id", true)
				So(dbCalled, ShouldBeFalse)
				So(id, ShouldEqual, userJane.ids[0])
				name, dbCalled := userJane.get("name", true)
				So(dbCalled, ShouldBeFalse)
				So(name, ShouldEqual, "Jane A. Smith")
				email, dbCalled := userJane.get("email", true)
				So(dbCalled, ShouldBeFalse)
				So(email, ShouldEqual, "jane.smith@example.com")
			})
		})
	})
}
