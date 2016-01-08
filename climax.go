// Package climax is a handy alternative CLI for Go applications. It looks
// pretty much exactly like the output of the default `go` command and
// incorporates some fancy features from it. For instance, Climax does
// support so-called topics (some sort of Wiki entries for CLI).
// You can also define some annotated use cases of some command that
// would get displayed in the help section of corresponding command.
//
// Climax-based applications produce this sort of output:
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
