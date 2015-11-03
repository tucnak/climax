// Package climax is a handy alternative CLI for Go applications. It looks
// pretty much exactly like the output of the default `go` command and
// incorporates some fancy features from it. For instance, Climax does
// support so-called topics (some sort of Wiki entries for CLI).
// You can also define some annotated use cases of some command that
// would get displayed in the help section of corresponding command.
//
// After running Climax-based application, you'd likely see something
// like this:
//
//		Camus is a modern content writing suite.
//
//		Usage:
//
//			camus command [arguments]
//
//		The commands are:
//
//			init        starts a new project
//			new         creates flavored book parts
//
//		Use "camus help [command]" for more information about a command.
//
//		Additional help topics:
//
//			writing     markdown language cheatsheet
//			metadata    intro to yaml-based metadata
//			realtime    effective real-time writing
//
//		Use "camus help [topic]" for more information about a topic.
//
// Here is an example of trivial CLI application that does nothing, but
// provides a single string split-like functionality:
//
//		demo := climax.New("demo")
//		demo.Brief = "Demo is a funky demonstation of Climax capabilities."
//		demo.Version = "stable"
//
//		joinCmd := climax.Command{
//			Name:  "join",
//			Brief: "merges the strings given",
//			Usage: `[-s=] "a few" distinct strings`,
//			Help:  `Lorem ipsum dolor sit amet amet sit todor...`,
//
//			Flags: []climax.Flag{
//				{
//					Name:     "separator",
//					Short:    "s",
//					Usage:    `--separator="."`,
//					Help:     `Put some separating string between all the strings given.`,
//					Variable: true,
//				},
//			},
//
//			Examples: []climax.Example{
//				{
//					Usecase:     `-s . "google" "com"`,
//					Description: `Results in "google.com"`,
//				},
//			},
//
//			Handle: func(ctx climax.Context) int {
//				var separator string
//				if sep, ok := ctx.Get("separator"); ok {
//					separator = sep
//				}
//
//				fmt.Println(strings.Join(ctx.Args, separator))
//
//				return 0
//			},
//		}
//
//		demo.AddCommand(joinCmd)
//		demo.Run()
//
// Have fun!
package climax

// New constructs a new CLI application with a given name.
// In case of an empty name it will panic.
func New(name string) *Application {
	if name == "" {
		panic("can't construct an app without a name")
	}

	return &Application{Name: name}
}
