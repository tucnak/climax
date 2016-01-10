package climax

import (
	"fmt"
	"io"
	"os"
)

var (
	outputDevice io.Writer = os.Stdout
	errorDevice  io.Writer = os.Stderr
)

// Application is a main CLI instance.
//
// By default, Climax provides its own implementation of version
// command, but it will use "version" command instead if you
// provide one.
type Application struct {
	Name    string // `go`
	Brief   string // `Go is a tool for managing Go source code.`
	Version string // `1.5`

	Commands []Command
	Topics   []Topic
	Groups   []Group

	// Default is a default handler. It gets executed if there are
	// no command line arguments (except the program name), when
	// otherwise, by default, the help entry is being shown.
	Default CmdHandler

	ungroupedCmdsCount int
}

// Group connects a list of commands with a descriptive string.
//
// The Name is used in the help output to group related commands together.
type Group struct {
	Name     string
	Commands []*Command
}

func (a *Application) println(stuff ...interface{}) {
	fmt.Fprintln(outputDevice, stuff...)
}

func (a *Application) printf(format string, stuff ...interface{}) {
	fmt.Fprintf(outputDevice, format, stuff...)
}

func (a *Application) printerr(err ...interface{}) {
	for _, each := range err {
		fmt.Fprintln(errorDevice, a.Name+":", each)
	}
}

func (a *Application) commandByName(name string) *Command {
	for i, command := range a.Commands {
		if command.Name == name {
			return &a.Commands[i]
		}
	}

	return nil
}

func (a *Application) topicByName(name string) *Topic {
	for i, topic := range a.Topics {
		if topic.Name == name {
			return &a.Topics[i]
		}
	}

	return nil
}

func (a *Application) groupByName(name string) *Group {
	for i, group := range a.Groups {
		if group.Name == name {
			return &a.Groups[i]
		}
	}

	return nil
}

func (a *Application) isNameAvailable(name string) bool {
	hypo, jypo := a.commandByName(name), a.topicByName(name)
	if hypo != nil || jypo != nil {
		return false
	}

	return true
}

// AddGroup adds a new empty, named group.
//
// Pass the returned group name to Command's Group member
// to make the command part of the group.
func (a *Application) AddGroup(name string) string {
	a.Groups = append(a.Groups, Group{Name: name})
	return name
}

// AddCommand does literally what its name says.
func (a *Application) AddCommand(command Command) {
	a.Commands = append(a.Commands, command)

	newCmd := &a.Commands[len(a.Commands)-1]
	if newCmd.Group != "" {
		group := a.groupByName(newCmd.Group)
		if group == nil {
			panic("group doesn't exist")
		}

		group.Commands = append(group.Commands, newCmd)
	} else {
		a.ungroupedCmdsCount++
	}
}

// AddTopic does literally what its name says.
func (a *Application) AddTopic(topic Topic) {
	a.Topics = append(a.Topics, topic)
}

// Run executes a CLI.
//
// Take a note, Run panics if len(os.Args) < 1
func (a *Application) Run() int {
	if len(os.Args) < 1 {
		panic("shell-provided arguments are not present")
	}
	arguments := os.Args[1:]
	// $ program
	//           ^ no args
	if len(arguments) == 0 {
		if a.Default == nil {
			a.println(a.globalHelp())
			return 0
		}

		return a.Default(*newContext(a))
	}

	yankeeGoHome := func(errMsg string) {
		a.printerr(errMsg)
		os.Exit(1)
	}

	subcommandName := arguments[0]
	subcommand := a.commandByName(subcommandName)

	if subcommandName == "help" {
		// $ program help
		//           ^ one argument
		if len(arguments) <= 1 {
			a.println(a.globalHelp())
			return 0
		}

		command := a.commandByName(arguments[1])
		if command != nil {
			a.println(a.commandHelp(command))
			return 0
		}

		topic := a.topicByName(arguments[1])
		if topic != nil {
			a.println(topic.Text)
			return 0
		}

		yankeeGoHome("no such command or help topic")
	}

	if subcommandName == "version" {
		if subcommand != nil {
			return subcommand.Run(Context{})
		}

		a.printf("%s version %s\n", a.Name, a.Version)
		return 0
	}

	if subcommand != nil {
		context, err := a.parseContext(subcommand.Flags, arguments[1:])
		if err != nil {
			yankeeGoHome(err.Error())
		}

		return subcommand.Run(*context)
	}

	yankeeGoHome("unknown subcommand \"" + subcommandName + "\"\n")
	return 1
}

// Log prints the message to stderrr (each argument takes a distinct line).
func (a *Application) Log(lines ...interface{}) {
	a.printerr(lines...)
}
