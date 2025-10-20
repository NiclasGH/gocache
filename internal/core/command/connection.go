package command

import (
	"gocache/internal/core/resp"
	"gocache/internal/persistence"
	"strings"
)

// / Returns PONG if no argument is provided, otherwise return a copy of the argument as a bulk.
// / PING {name}?
// / Example:
// / Req: PING
// / Res: PONG
func pingStrategy(request resp.Value, _ persistence.Database) resp.Value {
	args := request.GetArgs()

	if len(args) == 0 {
		return resp.Value{Typ: "string", Str: "PONG"}
	}

	return resp.Value{Typ: resp.STRING.Typ, Str: args[0].Bulk}
}

// / Gives information about commands and about available commands
// / COMMAND -> All available commands and their specs (command structure, acl categories, tips, key specification and subcommands). For simplicity reason, I will implement only the first seven categories
// / COMMAND {command} -> Same as Command but filtered to the command
// / COMMAND DOCS -> Docs about the commands. may include: summary, since redis version, functional group, complexity, doc_flags, arguments. We only use summary, group and complexity
func commandMetadataStrategy(request resp.Value, _ persistence.Database) resp.Value {
	args := request.GetArgs()

	commandFilter := ""
	if len(args) >= 1 {
		commandFilter = strings.ToUpper(args[0].Bulk)
	}

	var result []resp.Value

	metadata := commandList()
	if commandFilter == "DOCS" {
		if len(args) >= 2 {
			commandFilter = strings.ToUpper(args[1].Bulk)
		} else {
			commandFilter = ""
		}

		result = filterCommands(metadata, commandFilter, (*commandMetadata).docs)
	} else {
		result = filterCommands(metadata, commandFilter, (*commandMetadata).specs)
	}

	return resp.Value{Typ: resp.ARRAY.Typ, Array: result}
}

func filterCommands(items []commandMetadata, filter string, accessor func(*commandMetadata) []resp.Value) []resp.Value {
	result := make([]resp.Value, 0, len(items))
	for i := range items {
		if filter == "" || items[i].name == filter {
			result = append(result, accessor(&items[i])...)
		}
	}
	return result
}

func commandList() []commandMetadata {
	docs := make([]commandMetadata, 0, len(commandMetadatas))
	for _, v := range commandMetadatas {
		docs = append(docs, v)
	}
	return docs
}
