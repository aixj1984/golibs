/*
 * Licensed to the AcmeStack under one or more contributor license
 * agreements. See the NOTICE file distributed with this work for
 * additional information regarding copyright ownership.
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package tests

import (
	"testing"

	"github.com/aixj1984/golibs/gorm-plus/gplus"
)

func TestRouteInsert1Name(t *testing.T) {
	expectSql := "INSERT INTO `Users2` (`username`,`password`,`address`,`age`,`phone`,`score`,`dept`) VALUES ('afumu','123456','',18,'',12,'研发部门') RETURNING `id`"
	user := &User{Username: "afumu", Password: "123456", Age: 18, Score: 12, Dept: "研发部门"}
	u := gplus.GetModel[User]()
	sessionDb := checkInsertSql(t, expectSql)
	gplus.Insert(&user, gplus.Table("Users2"), gplus.Db(sessionDb), gplus.Omit(&u.CreatedAt, &u.UpdatedAt))
}

func TestRouteSelectByIdName(t *testing.T) {
	expectSql := "SELECT * FROM `Users2` WHERE id = 1  LIMIT 1"
	sessionDb := checkSelectSql(t, expectSql)
	gplus.SelectById[User](1, gplus.Db(sessionDb), gplus.Table("Users2"))
}

func TestRouteSelectByIdSelect(t *testing.T) {
	expectSql := "SELECT `username`,`age` FROM `Users2` WHERE id = 1  LIMIT 1"
	sessionDb := checkSelectSql(t, expectSql)
	u := gplus.GetModel[User]()
	gplus.SelectById[User](1, gplus.Db(sessionDb), gplus.Table("Users2"), gplus.Select(&u.Username, &u.Age))
}
