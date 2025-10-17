package command

import "gocache/internal/resp"

type intoValue interface {
	values() []resp.Value
	getCommand() string
}

type commandSpec struct {
	command       string
	argCount      int
	flags         []string
	firstKey      int
	lastKey       int
	steps         int
	aclCategories []string
}

func (spec commandSpec) getCommand() string {
	return spec.command
}

func (spec commandSpec) values() []resp.Value {
	flags := make([]resp.Value, len(spec.flags))
	for i, v := range spec.flags {
		flags[i] = resp.Value{Typ: resp.BULK.Typ, Bulk: v}
	}

	aclCategories := make([]resp.Value, len(spec.aclCategories))
	for i, v := range spec.aclCategories {
		aclCategories[i] = resp.Value{Typ: resp.BULK.Typ, Bulk: v}
	}

	return []resp.Value{
		{
			Typ: resp.ARRAY.Typ,
			Array: []resp.Value{
				{Typ: resp.BULK.Typ, Bulk: spec.command},
				{Typ: resp.INTEGER.Typ, Num: spec.argCount},
				{
					Typ:   resp.ARRAY.Typ,
					Array: flags,
				},
				{Typ: resp.INTEGER.Typ, Num: spec.firstKey},
				{Typ: resp.INTEGER.Typ, Num: spec.lastKey},
				{Typ: resp.INTEGER.Typ, Num: spec.steps},
				{
					Typ:   resp.ARRAY.Typ,
					Array: aclCategories,
				},
			},
		},
	}
}

type commandDoc struct {
	command    string
	summary    string
	since      string
	group      string
	complexity string
}

func (doc commandDoc) getCommand() string {
	return doc.command
}

func (doc commandDoc) values() []resp.Value {
	docs := []resp.Value{
		{
			Typ:  resp.BULK.Typ,
			Bulk: "summary",
		},
		{
			Typ:  resp.BULK.Typ,
			Bulk: doc.summary,
		},
		{
			Typ:  resp.BULK.Typ,
			Bulk: "since",
		},
		{
			Typ:  resp.BULK.Typ,
			Bulk: doc.since,
		},

		{
			Typ:  resp.BULK.Typ,
			Bulk: "group",
		},
		{
			Typ:  resp.BULK.Typ,
			Bulk: doc.group,
		},

		{
			Typ:  resp.BULK.Typ,
			Bulk: "complexity",
		},
		{
			Typ:  resp.BULK.Typ,
			Bulk: doc.complexity,
		},
	}

	return []resp.Value{
		{
			Typ:  resp.BULK.Typ,
			Bulk: doc.command,
		},
		{
			Typ:   resp.ARRAY.Typ,
			Array: docs,
		},
	}
}
