package parser

import (
	"fmt"
	"github.com/deis/workflow-cli/cmd"
	"github.com/docopt/docopt-go"
	"k8s.io/api/core/v1"
	"strconv"
)

// Toleration routes toleration commands to their specific function.
func Toleration(argv []string, cmdr cmd.Commander) error {
	usage := `
Valid commands for toleration:

toleration:list        list toleration for an app
toleration:set         set toleration for an app
toleration:unset       unset toleration for an app

Use 'deis help [command]' to learn more.
`

	switch argv[0] {
	case "toleration:list":
		return tolerationList(argv, cmdr)
	case "toleration:set":
		return tolerationSet(argv, cmdr)
	case "toleration:unset":
		return tolerationUnset(argv, cmdr)
	default:
		if printHelp(argv, usage) {
			return nil
		}

		if argv[0] == "toleration" {
			argv[0] = "toleration:list"
			return tolerationList(argv, cmdr)
		}

		PrintUsage(cmdr)
		return nil
	}
}

func tolerationList(argv []string, cmdr cmd.Commander) error {
	usage := `
Lists toleration for the pods of an application.

Usage: deis toleration:list [options]

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

	return cmdr.TolerationList(app, format)
}

func tolerationSet(argv []string, cmdr cmd.Commander) error {
	usage := `
Sets toleration for the pods of an application. For more details on how to use tolerations, refer to the official Kubernetes documentation https://kubernetes.io/docs/concepts/configuration/taint-and-toleration/

Usage: deis toleration:set --type=<app_type> <name> [--key=<key>] [--value=<value>] [--operator=<operator>] [--effect=<effect>] [--toleration-seconds=<seconds>] [options]

Arguments:
  <app_type>
    the process type as defined in your Procfile, such as 'web' or 'worker'.
    Note that Dockerfile apps have a default 'cmd' process type.
  <name>
	A unique identifier for this toleration.
  <key>
    the key to check in this toleration.
  <value>
    the value of the key to check. If not set, will match any value.
  <operator>
    the operator to be used to compare the key's value. Default: Equal.
  <effect>
	the effect that should happen on a toleration match. Default: NoSchedule.
  <seconds>
	the time (in seconds) to wait before evicting a running pod from a newly-tained node.

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

	var toleration v1.Toleration
	identifier := safeGetValue(args, "<name>")
	if identifier == "" {
        return fmt.Errorf("expected identifier not to be null")
	}
	key := safeGetValue(args, "--key")
	if key != "" {
		toleration.Key = key
	}
	value := safeGetValue(args, "--value")
	if value != "" {
		toleration.Value = value
	}
	operator := safeGetValue(args, "--operator")
	if operator != "" {
        toleration.Operator = v1.TolerationOperator(operator)
	}
	effect := safeGetValue(args, "--effect")
	if effect != "" {
        toleration.Effect = v1.TaintEffect(effect)
	}
	seconds := safeGetValue(args, "--toleration-seconds")
	if seconds != "" {
		secondsInt, err := strconv.ParseInt(seconds, 10, 0)
		if err != nil {
			return err
		}
		toleration.TolerationSeconds = &secondsInt
	}

	return cmdr.TolerationSet(app, appType, identifier, toleration)
}

func tolerationUnset(argv []string, cmdr cmd.Commander) error {
	usage := `
Unsets a toleration for the pods of an application.

Usage: deis toleration:unset --type=<app_type> <key>... [options]

Arguments:
  <app_type>
    the process type as defined in your Procfile, such as 'web' or 'worker'.
    Note that Dockerfile apps have a default 'cmd' process type.
  <key>
    the toleration key to remove from the application's pod.

Options:
  -a --app=<app>
    the uniquely identifiable name for the toleration.
  -t --type=<app_type>
	the process type to be affected by these tolerations.
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

	return cmdr.TolerationUnset(app, appType, args["<key>"].([]string))
}
