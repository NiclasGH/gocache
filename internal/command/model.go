package command

import "gocache/internal/resp"

type Handler = func([]resp.Value) resp.Value

type newCoolCommand struct {
	name string
	spec commandSpec
	doc  commandDoc
}

type commandSpec struct {
	argCount      int
	flags         []string
	firstKey      int
	lastKey       int
	steps         int
	aclCategories []string
}

type commandDoc struct {
	summary    string
	since      string
	group      string
	complexity string
}

func (c newCoolCommand) specs() []resp.Value {
	flags := make([]resp.Value, len(c.spec.flags))
	for i, v := range c.spec.flags {
		flags[i] = resp.Value{Typ: resp.BULK.Typ, Bulk: v}
	}

	aclCategories := make([]resp.Value, len(c.spec.aclCategories))
	for i, v := range c.spec.aclCategories {
		aclCategories[i] = resp.Value{Typ: resp.BULK.Typ, Bulk: v}
	}

	return []resp.Value{
		{
			Typ: resp.ARRAY.Typ,
			Array: []resp.Value{
				{Typ: resp.BULK.Typ, Bulk: c.name},
				{Typ: resp.INTEGER.Typ, Num: c.spec.argCount},
				{
					Typ:   resp.ARRAY.Typ,
					Array: flags,
				},
				{Typ: resp.INTEGER.Typ, Num: c.spec.firstKey},
				{Typ: resp.INTEGER.Typ, Num: c.spec.lastKey},
				{Typ: resp.INTEGER.Typ, Num: c.spec.steps},
				{
					Typ:   resp.ARRAY.Typ,
					Array: aclCategories,
				},
			},
		},
	}
}

func (c newCoolCommand) docs() []resp.Value {
	docs := []resp.Value{
		{
			Typ:  resp.BULK.Typ,
			Bulk: "summary",
		},
		{
			Typ:  resp.BULK.Typ,
			Bulk: c.doc.summary,
		},
		{
			Typ:  resp.BULK.Typ,
			Bulk: "since",
		},
		{
			Typ:  resp.BULK.Typ,
			Bulk: c.doc.since,
		},

		{
			Typ:  resp.BULK.Typ,
			Bulk: "group",
		},
		{
			Typ:  resp.BULK.Typ,
			Bulk: c.doc.group,
		},

		{
			Typ:  resp.BULK.Typ,
			Bulk: "complexity",
		},
		{
			Typ:  resp.BULK.Typ,
			Bulk: c.doc.complexity,
		},
	}

	return []resp.Value{
		{
			Typ:  resp.BULK.Typ,
			Bulk: c.name,
		},
		{
			Typ:   resp.ARRAY.Typ,
			Array: docs,
		},
	}
}
