package climax

import (
	"bytes"
	"strings"
	"text/template"
)

const globalHelpTemplate string = `{{.Brief}}

Usage:

	{{.Name}} {{if .Commands}}command [arguments]{{end}}

{{if .Commands}}The commands are:
{{if .UngroupedCount}}{{range .Commands}}
	{{if not .Group}}{{.Name | printf "%-11s"}} {{.Brief}}{{end}}{{end}}
{{end}}{{range .Groups}}{{if .Commands}}
{{.Name}}
	{{range .Commands}}
	{{.Name | printf "%-11s"}} {{.Brief}}{{end}}
	{{end}}{{end}}
Use "{{.Name}} help [command]" for more information about a command.{{end}}
{{if .Topics}}
Additional help topics:
{{range .Topics}}
	{{.Name | printf "%-11s"}} {{.Brief}}{{end}}

Use "{{.Name}} help [topic]" for more information about a topic.
{{end}}`

const commandHelpTemplate string = `Usage: {{commandUsage .Command}}

{{.Help}}
{{if .Flags}}
Available options:
{{range .Flags}}
	{{flagUsage . false}}
		{{.Help | tabout}}{{end}}{{end}}
{{if .Examples}}
Examples:
{{$app := .App}}{{$cmd := .Name}}{{range .Examples}}
	$ {{$app}} {{$cmd}} {{.Usecase}}
		{{.Description | tabout}}
{{end}}{{end}}`

func templated(canvas string, data interface{}) string {
	t := template.New("")
	t.Funcs(template.FuncMap{
		"tabout":       alignMultilineHelp,
		"commandUsage": commandUsage,
		"flagUsage":    flagUsage,
	})
	template.Must(t.Parse(canvas))

	var b bytes.Buffer

	err := t.Execute(&b, data)
	if err != nil {
		panic(err)
	}

	output := b.String()

	// TODO: Fix this nasty workaround for templating ASAP!
	output = strings.Replace(output, "\n\n\n", "", -1)

	return output
}

func alignMultilineHelp(text string) string {
	return strings.Replace(text, "\n", "\n\t\t", -1)
}

func commandUsage(command Command) string {
	if command.Usage != "" {
		return command.Name + " " + command.Usage
	}

	usage := command.Name
	for _, flag := range command.Flags {
		usage += " [" + flagUsage(flag, true) + "]"
	}

	return usage
}

func flagUsage(flag Flag, tiny bool) string {
	if tiny {
		if flag.Short != "" {
			return "-" + flag.Short
		}

		return "--" + flag.Name
	}

	var short string
	if flag.Short != "" {
		short = "-" + flag.Short + ", "
	}

	if flag.Usage != "" {
		return short + flag.Usage
	}

	usage := "--" + flag.Name
	if flag.Variable {
		usage += "=\"\""
	}

	return short + usage
}

func (a *Application) globalHelp() string {
	return templated(globalHelpTemplate, struct {
		Application
		UngroupedCount int
	}{
		*a,
		a.ungroupedCmdsCount,
	})
}

func (a *Application) commandHelp(command *Command) string {
	return templated(commandHelpTemplate, struct {
		Command
		App string
	}{
		*command,
		a.Name,
	})
}
