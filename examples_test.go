package climax_test

import (
	"fmt"
	"github.com/tucnak/climax"
	"strings"
)

func Example_application() {
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
}

// Handler accepts a climax.Context object and returns an exitcode integer.
func Example_handler(ctx climax.Context) int {
	if len(ctx.Args) < 2 {
		ctx.Log("not enough arguments")

		// with os.Exit(1)
		return 1
	}

	if name, ok := ctx.Get("name"); ok {
		// argument `name` parsed
		fmt.Println(name)

	} else {
		ctx.Log("name not specified")

		return 1
	}

	return 0
}
