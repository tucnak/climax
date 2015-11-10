package climax

import (
	"fmt"
	"os"
	"strings"
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

	Commands   []Command
	Topics     []Topic
	Categories map[string][]Command
}

func newApplication(name string) *Application {
	return &Application{
		Name:       name,
		Categories: make(map[string][]Command),
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

func (a Application) isNameAvailable(name string) bool {
	hypo, jypo := a.commandByName(name), a.topicByName(name)
	if hypo != nil || jypo != nil {
		return false
	}

	return true
}

// AddCommand does literally what its name says.
func (a *Application) AddCommand(command Command) {
	a.Commands = append(a.Commands, command)

	category := strings.ToUpper(command.Category)
	if _, ok := a.Categories[category]; !ok {
		a.Categories[category] = []Command{command}
	} else {
		a.Categories[category] = append(a.Categories[category], command)
	}
}

// AddTopic does literally what its name says.
func (a *Application) AddTopic(topic Topic) {
	a.Topics = append(a.Topics, topic)
}

// Run executes a CLI.
//
// Take a note, Run panics if len(os.Args) < 1
func (a Application) Run() int {
	if len(os.Args) < 1 {
		panic("shell-provided arguments are not present")
	}
	arguments := os.Args[1:]
	// $ program
	//           ^ no args
	if len(arguments) == 0 {
		println(a.globalHelp())
		return 0
	}

	yankeeGoHome := func(errMsg string) {
		printerr(fmt.Errorf("%s: %s", a.Name, errMsg))
		os.Exit(1)
	}

	subcommandName := arguments[0]
	subcommand := a.commandByName(subcommandName)

	if subcommandName == "help" {
		// $ program help
		//           ^ one argument
		if len(arguments) <= 1 {
			println(a.globalHelp())
			return 0
		}

		command := a.commandByName(arguments[1])
		if command != nil {
			println(a.commandHelp(command))
			return 0
		}

		topic := a.topicByName(arguments[1])
		if topic != nil {
			println(topic.Text)
			return 0
		}

		yankeeGoHome("no such command or help topic")
	}

	if subcommandName == "version" {
		if subcommand != nil {
			return subcommand.Run(Context{})
		}

		printf("%s version %s\n", a.Name, a.Version)
		return 0
	}

	if subcommand != nil {
		context, err := parseContext(subcommand.Flags, arguments[1:])
		if err != nil {
			yankeeGoHome(err.Error())
		}

		return subcommand.Run(*context)
	}

	yankeeGoHome("unknown subcommand \"" + subcommandName + "\"\n")
	return 1
}
