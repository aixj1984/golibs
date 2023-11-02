package gorm

import "gorm.io/gorm"

// DBConditions DB常用的查询条件封装
type DBConditions struct {
	And       map[string]interface{}
	Or        map[string]interface{}
	Not       map[string]interface{}
	Limit     int
	Offset    int
	Order     interface{}
	Select    interface{}
	Group     string
	Having    interface{}
	NeedCount bool
	Count     int64
	Distinct  interface{} // grom v1 暂不支持
}

// Fill 填充查询条件
func (d *DBConditions) Fill(db *gorm.DB) *gorm.DB {
	if d.Select != nil {
		db = db.Select(d.Select)
	}

	for cond, val := range d.And {
		db = db.Where(cond, val)
	}
	for cond, val := range d.Not {
		db = db.Not(cond, val)
	}
	for cond, val := range d.Or {
		db = db.Or(cond, val)
	}

	if d.NeedCount {
		db = db.Count(&d.Count)
	}
	if d.Order != nil {
		db = db.Order(d.Order)
	}
	if d.Limit != 0 {
		db = db.Limit(d.Limit)
	}
	if d.Offset != 0 {
		db = db.Offset(d.Offset)
	}
	if d.Group != "" {
		db = db.Group(d.Group)
	}
	if d.Having != nil {
		db = db.Having(d.Having)
	}
	if d.Distinct != nil {
		db = db.Distinct(d.Distinct)
	}

	return db
}

/* demo
cond := &base.DBConditions{
	And: map[string][]interface{}{
		"id IN (?)": {95,96,97},
	},
	Not: map[string][]interface{}{
		"id": {96},
	},
	Limit: 1,
	Offset: 1,
	Order: "id DESC",
}
*/
