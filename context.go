package climax

import (
	"bytes"
	"fmt"
	"strings"
)

// Context is a set of arguments and options of command call.
type Context struct {
	// Real arguments, excluding any option flags and their values.
	//
	// Example:
	//
	//     $ app command --force -s="magic" 42 fairy
	//
	//     []string{"42", "fairy"}
	//
	Args []string

	NonVariable map[string]bool
	Variable    map[string]string

	app *Application
}

// Log prints the message to stderrr (each argument takes a distinct line).
func (c *Context) Log(data ...interface{}) {
	c.app.Log(data...)
}

// Is returns true if a flag with corresponding name is defined.
func (c *Context) Is(flagName string) bool {
	if _, ok := c.NonVariable[flagName]; ok {
		return true
	}

	if _, ok := c.Variable[flagName]; ok {
		return true
	}

	return false
}

// Get returns a value of corresponding variable flag.
// Second (bool) parameter says whether it's really defined or not.
func (c *Context) Get(variableFlagName string) (string, bool) {
	value, ok := c.Variable[variableFlagName]
	return value, ok
}

func looksLikeFlag(flag string) bool {
	if strings.HasPrefix(flag, "-") || strings.HasPrefix(flag, "--") {
		return true
	}

	return false
}

func parseFlagSignature(flag string) (string, string) {
	flag = strings.TrimLeft(flag, "-")

	equalPos := strings.Index(flag, "=")
	if equalPos < 0 {
		return flag, ""
	}

	return flag[:equalPos], flag[equalPos+1:]
}

func flagByName(flags *[]Flag, name string) *Flag {
	for i, flag := range *flags {
		if flag.Name == name || flag.Short == name {
			return &(*flags)[i]
		}
	}

	return nil
}

func newContext(app *Application) *Context {
	ctx := Context{}

	ctx.Args = []string{}
	ctx.NonVariable = make(map[string]bool)
	ctx.Variable = make(map[string]string)

	ctx.app = app

	return &ctx
}

func (a *Application) parseContext(flags []Flag, argv []string) (*Context, error) {
	ctx := newContext(a)

	for i := 0; i < len(argv); i++ {
		argument := argv[i]

		if !looksLikeFlag(argument) {
			ctx.Args = append(ctx.Args, argument)
			continue
		}

		name, value := parseFlagSignature(argument)
		flag := flagByName(&flags, name)

		if flag == nil {
			return nil, fmt.Errorf(`option -%s does not exist`, name)
		}

		if flag.Variable {
			if value == "" {
				if strings.HasSuffix(argument, "=") {
					ctx.Variable[flag.Name] = ""
					continue
				}

				if i+1 >= len(argv) {
					return nil, fmt.Errorf(`option -%s is invalid`, name)
				}

				if looksLikeFlag(argv[i+1]) {
					return nil, fmt.Errorf(`option -%s is invalid`, name)
				}

				ctx.Variable[flag.Name] = argv[i+1]
				i++
				continue
			}

			ctx.Variable[flag.Name] = value

		} else {
			if value != "" {
				return nil, fmt.Errorf(`-%s is not variable option`, name)
			}

			ctx.NonVariable[flag.Name] = true
		}
	}

	return ctx, nil
}

func (c Context) String() string {
	var b bytes.Buffer

	fmt.Fprintf(&b, "Context {\n")
	fmt.Fprintf(&b, "\tArgs: %q\n", c.Args)
	fmt.Fprintf(&b, "\tFlags:\n")

	for flag := range c.NonVariable {
		fmt.Fprintf(&b, "\t\t%s\n", flag)
	}

	for flag, value := range c.Variable {
		fmt.Fprintf(&b, "\t\t%s=%s\n", flag, value)
	}

	fmt.Fprintf(&b, "}")

	return b.String()
}
