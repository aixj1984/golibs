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
	"errors"
	"testing"

	"github.com/aixj1984/golibs/gorm-plus/gplus"
	"gorm.io/gorm"
)

// TestTransactionCommitMultipleOps runs several gplus calls on the same *gorm.DB transaction;
// nil error commits and all rows must be visible.
func TestTransactionCommitMultipleOps(t *testing.T) {
	deleteOldData()

	err := gormDb.Transaction(func(tx *gorm.DB) error {
		u1 := &User{Username: "tx_ok_u1", Password: "p1", Age: 1, Score: 1, Dept: "d1"}
		u2 := &User{Username: "tx_ok_u2", Password: "p2", Age: 2, Score: 2, Dept: "d2"}
		if r := gplus.Insert(u1, gplus.Db(tx)); r.Error != nil {
			return r.Error
		}
		if r := gplus.Insert(u2, gplus.Db(tx)); r.Error != nil {
			return r.Error
		}
		return nil
	})
	if err != nil {
		t.Fatalf("transaction should commit: %v", err)
	}

	got1, db1 := gplus.SelectById[User](mustUserIDByUsername(t, "tx_ok_u1"))
	if db1.Error != nil {
		t.Fatalf("select u1: %v", db1.Error)
	}
	if got1.Username != "tx_ok_u1" {
		t.Fatalf("u1 username: got %q", got1.Username)
	}
	got2, db2 := gplus.SelectById[User](mustUserIDByUsername(t, "tx_ok_u2"))
	if db2.Error != nil {
		t.Fatalf("select u2: %v", db2.Error)
	}
	if got2.Username != "tx_ok_u2" {
		t.Fatalf("u2 username: got %q", got2.Username)
	}
}

// TestTransactionRollbackOnError forces an error after a successful insert inside the same tx;
// the insert must be rolled back.
func TestTransactionRollbackOnError(t *testing.T) {
	deleteOldData()

	errForced := errors.New("force rollback")
	err := gormDb.Transaction(func(tx *gorm.DB) error {
		u := &User{Username: "tx_rb_user", Password: "p", Age: 3, Score: 3, Dept: "d"}
		if r := gplus.Insert(u, gplus.Db(tx)); r.Error != nil {
			return r.Error
		}
		return errForced
	})
	if !errors.Is(err, errForced) {
		t.Fatalf("expected wrapped/forced error, got %v", err)
	}

	q, u := gplus.NewQuery[User]()
	q.Eq(&u.Username, "tx_rb_user")
	n, db := gplus.SelectCount[User](q)
	if db.Error != nil {
		t.Fatalf("count: %v", db.Error)
	}
	if n != 0 {
		t.Fatalf("rollback expected 0 rows for tx_rb_user, got %d", n)
	}
}

// TestTransactionBeginManualCommit uses gplus.Begin then Commit; all ops use gplus.Db(tx).
func TestTransactionBeginManualCommit(t *testing.T) {
	deleteOldData()

	tx := gplus.Begin()
	if tx.Error != nil {
		t.Fatalf("begin: %v", tx.Error)
	}
	u := &User{Username: "tx_manual", Password: "p", Age: 4, Score: 4, Dept: "d"}
	if r := gplus.Insert(u, gplus.Db(tx)); r.Error != nil {
		_ = tx.Rollback()
		t.Fatalf("insert: %v", r.Error)
	}
	if r := tx.Commit(); r.Error != nil {
		t.Fatalf("commit: %v", r.Error)
	}

	got, db := gplus.SelectById[User](mustUserIDByUsername(t, "tx_manual"))
	if db.Error != nil {
		t.Fatalf("select: %v", db.Error)
	}
	if got.Username != "tx_manual" {
		t.Fatalf("username: got %q", got.Username)
	}
}

func mustUserIDByUsername(t *testing.T, username string) int64 {
	t.Helper()
	q, u := gplus.NewQuery[User]()
	q.Eq(&u.Username, username)
	row, db := gplus.SelectOne[User](q)
	if db.Error != nil {
		t.Fatalf("lookup %q: %v", username, db.Error)
	}
	return row.ID
}
