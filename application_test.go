package climax

import (
	"os"
	"testing"
)

func TestRun_Bare(t *testing.T) {
	mustPanic(t, "no shell args", func() {
		os.Args = []string{}
		defer setArguments()

		a := Application{}
		a.Run()
	})
}

func TestRun_Version(t *testing.T) {
	a := New("application")
	a.Version = "5.0"
	setArguments("version")
	defer setArguments()
	defer output.Reset()

	if exitcode := a.Run(); exitcode != 0 {
		t.Errorf("finished with code %d, expected 0", exitcode)
	}

	expected := "application version 5.0\n"
	if output.String() != expected {
		t.Errorf("actual output is different to expected:\n")
		t.Logf("- expected: %q", expected)
		t.Logf("- recieved: %q", output.String())
	}
}

const expectedAppHelp string = `application is a thing

Usage:

	application command [arguments]

The commands are:

	open        opens smth
	close       closes smth
	attach      does some attaching
	detach      detaches your mind
	server      starts a web server

Use "application help [command]" for more information about a command.

Additional help topics:

	writing     how to write
	reading     how to read
	listening   how to listen to people
	speaking    how to talk

Use "application help [topic]" for more information about a topic.

`

const expectedCommandHelp string = `Usage: server -http=4747

Lorem ipsum dolor sit amet amet sit.

Available options:

	-p, -http=4747
		Specify the port of server's HTTP interface.

Examples:

	$ application server -http 4747
		Start a server on http port 4747.

`

func TestRun_Help(t *testing.T) {
	a := New("application")
	a.Brief = "application is a thing"

	a.AddCommand(Command{Name: "open", Brief: "opens smth"})
	a.AddCommand(Command{Name: "close", Brief: "closes smth"})
	a.AddCommand(Command{Name: "attach", Brief: "does some attaching"})
	a.AddCommand(Command{Name: "detach", Brief: "detaches your mind"})
	a.AddCommand(Command{
		Name:  "server",
		Brief: "starts a web server",
		Usage: "-http=4747",
		Help:  "Lorem ipsum dolor sit amet amet sit.",
		Flags: []Flag{
			{
				Name:  "http",
				Short: "p",
				Usage: "-http=4747",
				Help:  "Specify the port of server's HTTP interface.",
			},
		},
		Examples: []Example{
			{
				Usecase:     "-http 4747",
				Description: "Start a server on http port 4747.",
			},
		},
	})

	a.AddTopic(Topic{Name: "writing", Brief: "how to write"})
	a.AddTopic(Topic{Name: "reading", Brief: "how to read"})
	a.AddTopic(Topic{Name: "listening", Brief: "how to listen to people"})
	a.AddTopic(Topic{Name: "speaking", Brief: "how to talk"})

	if exitcode := a.Run(); exitcode != 0 {
		t.Errorf("finished with code %d, expected 0", exitcode)
	}

	if output.String() != expectedAppHelp {
		t.Errorf("global help output is different to expected:\n")
		t.Logf("- expected:\n%s", expectedAppHelp)
		t.Logf("- recieved:\n%s", output.String())
	}

	output.Reset()
	setArguments("help", "server")

	if exitcode := a.Run(); exitcode != 0 {
		t.Errorf("finished with code %d, expected 0", exitcode)
	}

	if output.String() != expectedCommandHelp {
		t.Errorf("command help output is different to expected:\n")
		t.Logf("- expected:\n%s", expectedCommandHelp)
		t.Logf("- recieved:\n%s", output.String())
	}

	setArguments()
	output.Reset()
}

func TestAddCommand(t *testing.T) {
	a := Application{}
	a.AddCommand(Command{})
	if len(a.Commands) != 1 {
		t.Error("broken")
	}
}

func TestAddTopic(t *testing.T) {
	a := Application{}
	a.AddTopic(Topic{})
	if len(a.Topics) != 1 {
		t.Error("broken")
	}
}

func TestAddFlag(t *testing.T) {
	var c Command
	c.AddFlag(Flag{})
	if len(c.Flags) != 1 {
		t.Error("broken")
	}
}

func TestAddExample(t *testing.T) {
	var c Command
	c.AddExample(Example{})
	if len(c.Examples) == 0 {
		t.Error("failed to add example")
	}
}
