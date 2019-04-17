package parser

import (
	"github.com/deis/workflow-cli/cmd"
	docopt "github.com/docopt/docopt-go"
)

// Annnotations routes annotation commands to their specific function.
func Annotation(argv []string, cmdr cmd.Commander) error {
	usage := `
Valid commands for annotation:

annotation:list        list annotations for an app
annotation:set         set annotations for an app
annotation:unset       unset annotations for an app

Use 'deis help [command]' to learn more.
`

	switch argv[0] {
	case "annotation:list":
		return annotationList(argv, cmdr)
	case "annotation:set":
		return annotationSet(argv, cmdr)
	case "annotation:unset":
		return annotationUnset(argv, cmdr)
	default:
		if printHelp(argv, usage) {
			return nil
		}

		if argv[0] == "annotation" {
			argv[0] = "annotation:list"
			return annotationList(argv, cmdr)
		}

		PrintUsage(cmdr)
		return nil
	}
}

func annotationList(argv []string, cmdr cmd.Commander) error {
	usage := `
Lists annotations for the pods of an application.

Usage: deis annotation:list [options]

Options:
  --oneline
    print output on one line.
  -a --app=<app>
    the uniquely identifiable name of the application.
`

	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}
	app := safeGetValue(args, "--app")
	oneline := args["--oneline"].(bool)

	format := ""
	if oneline {
		format = "oneline"
	}

	return cmdr.AnnotationList(app, format)
}

func annotationSet(argv []string, cmdr cmd.Commander) error {
	usage := `
Sets annotations for the pods of an application.

Usage: deis annotation:set --type=<app_type> <var>=<value> [<var>=<value>...] [options]

Arguments:
  <app_type>
    the process type as defined in your Procfile, such as 'web' or 'worker'.
    Note that Dockerfile apps have a default 'cmd' process type.
  <var>
    the uniquely identifiable name for the annotation.
  <value>
    the value of said annotation.

Options:
  -a --app=<app>
    the uniquely identifiable name for the application.
  -t --type=<app_type>
	the process type to be affected by these annotations.
`

	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}

	app := safeGetValue(args, "--app")

	appType := safeGetValue(args, "--type")
	if appType == "" {
		appType = "cmd"
	}

	return cmdr.AnnotationSet(app, appType, args["<var>=<value>"].([]string))
}

func annotationUnset(argv []string, cmdr cmd.Commander) error {
	usage := `
Unsets an annotation for the pods of an application.

Usage: deis annotation:unset --type=<app_type> <key>... [options]

Arguments:
  <app_type>
    the process type as defined in your Procfile, such as 'web' or 'worker'.
    Note that Dockerfile apps have a default 'cmd' process type.
  <key>
    the annotation to remove from the application's pod.

Options:
  -a --app=<app>
    the uniquely identifiable name for the application.
  -t --type=<app_type>
	the process type to be affected by these annotations.
`

	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}
	app := safeGetValue(args, "--app")

	appType := safeGetValue(args, "--type")
	if appType == "" {
		appType = "cmd"
	}

	return cmdr.AnnotationUnset(app, appType, args["<key>"].([]string))
}
