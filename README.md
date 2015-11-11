# Climax
>Climax is an alternative CLI that looks like Go command

[![GoDoc](https://godoc.org/github.com/tucnak/climax?status.svg)](https://godoc.org/github.com/tucnak/climax)
[![Travis](https://travis-ci.org/tucnak/climax.svg?branch=master)](https://travis-ci.org/tucnak/climax)

**Climax** is a handy alternative CLI (command-line interface) for Go apps.
It looks pretty much exactly like the output of the default `go` command and
incorporates some fancy features from it. For instance, Climax does support
so-called topics (some sort of Wiki entries for CLI). You can define some
annotated use cases of some command that would get displayed in the
help section of corresponding command also.

##### Why creating another CLI?
I didn't like existing solutions (e.g. codegangsta/cli | spf13/cobra) either for
bloated codebase (I dislike the huge complex libraries) or poor output
style / API. This project is just an another view on the subject, it has
slightly different API than, let's say, Cobra; I find it much more convenient.
<hr>

A sample application output, Climax produces:
```
Camus is a modern content writing suite.

Usage:

	camus command [arguments]

The commands are:

	init        starts a new project
	new         creates flavored book parts

Use "camus help [command]" for more information about a command.

Additional help topics:

	writing     markdown language cheatsheet
	metadata    intro to yaml-based metadata
	realtime    effective real-time writing

Use "camus help [topic]" for more information about a topic.
```

Here is an example of a trivial CLI application that does nothing,
but provides a single string split-like functionality:
```
demo := climax.New("demo")
demo.Brief = "Demo is a funky demonstation of Climax capabilities."
demo.Version = "stable"

joinCmd := climax.Command{
	Name:  "join",
	Brief: "merges the strings given",
	Usage: `[-s=] "a few" distinct strings`,
	Help:  `Lorem ipsum dolor sit amet amet sit todor...`,

	Flags: []climax.Flag{
		{
			Name:     "separator",
			Short:    "s",
			Usage:    `--separator="."`,
			Help:     `Put some separating string between all the strings given.`,
			Variable: true,
		},
	},

	Examples: []climax.Example{
		{
			Usecase:     `-s . "google" "com"`,
			Description: `Results in "google.com"`,
		},
	},

	Handle: func(ctx climax.Context) int {
		var separator string
		if sep, ok := ctx.Get("separator"); ok {
			separator = sep
		}

		fmt.Println(strings.Join(ctx.Args, separator))

		return 0
	},
}

demo.AddCommand(joinCmd)
demo.Run()
```

Have fun!
