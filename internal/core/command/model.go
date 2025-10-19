package command

import (
	"gocache/internal/core/resp"
	"gocache/internal/persistence"
)

type CommandStrategy = func(resp.Value, persistence.Database) resp.Value

type commandMetadata struct {
	name        string
	subCommands []commandMetadata
	spec        commandSpec
	doc         commandDoc
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

func (c commandMetadata) specs() []resp.Value {
	flags := make([]resp.Value, len(c.spec.flags))
	for i, v := range c.spec.flags {
		flags[i] = resp.Value{Typ: resp.BULK.Typ, Bulk: v}
	}

	aclCategories := make([]resp.Value, len(c.spec.aclCategories))
	for i, v := range c.spec.aclCategories {
		aclCategories[i] = resp.Value{Typ: resp.BULK.Typ, Bulk: v}
	}

	commandSpecs := []resp.Value{
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

	for _, v := range c.subCommands {
		commandSpecs = append(commandSpecs, v.specs()...)
	}

	return commandSpecs
}

func (c commandMetadata) docs() []resp.Value {
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

	commandDocs := []resp.Value{
		{
			Typ:  resp.BULK.Typ,
			Bulk: c.name,
		},
		{
			Typ:   resp.ARRAY.Typ,
			Array: docs,
		},
	}

	for _, v := range c.subCommands {
		commandDocs = append(commandDocs, v.docs()...)
	}

	return commandDocs
}
