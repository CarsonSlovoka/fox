package flag

import (
	"flag"
	"fmt"
)

type Command struct {
	*flag.FlagSet
	Fields   map[string][]CmdField // "myType": {{"varName", "value", "usage"}, ...}
	Init     func()
	Parse    func(args []string, reset bool) error
	MainFunc func(args []string) error
}

type CmdField struct {
	Name         string
	DefaultValue interface{}
	Usage        string
}

func defaultParse(cmd *Command, args []string, reset bool) error {
	if reset {
		for _, cmdFields := range cmd.Fields {
			for _, curCmdField := range cmdFields {
				field := cmd.FlagSet.Lookup(curCmdField.Name)
				if field != nil {
					if err := field.Value.Set(fmt.Sprintf("%v", curCmdField.DefaultValue)); err != nil {
						return err
					}
				}
			}
		}
	}
	if err := cmd.FlagSet.Parse(args); err != nil {
		return err
	}
	return nil
}

func NewCommand(flagSet *flag.FlagSet, fields map[string][]CmdField) *Command {
	cmd := &Command{FlagSet: flagSet, Fields: fields}
	cmd.Parse = func(args []string, reset bool) error {
		return defaultParse(cmd, args, reset)
	}

	// defaultInit
	cmd.Init = func() {
		for typeName, cmdFields := range cmd.Fields {
			switch typeName {
			case "string":
				for _, curCmdField := range cmdFields {
					cmd.FlagSet.String(curCmdField.Name, fmt.Sprintf("%v", curCmdField.DefaultValue), curCmdField.Usage)
				}
			case "int":
				for _, curCmdField := range cmdFields {
					cmd.FlagSet.Int(curCmdField.Name, curCmdField.DefaultValue.(int), curCmdField.Usage)
				}
			case "bool":
				for _, curCmdField := range cmdFields {
					cmd.FlagSet.Bool(curCmdField.Name, curCmdField.DefaultValue.(bool), curCmdField.Usage)
				}
			}
		}
	}
	cmd.Init()
	return cmd
}
