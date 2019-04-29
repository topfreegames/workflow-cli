package cmd

import (
	"fmt"
	"github.com/deis/pkg/prettyprint"

	"k8s.io/api/core/v1"
	"sort"
	"strings"

	"github.com/deis/controller-sdk-go/api"
	"github.com/deis/controller-sdk-go/config"
)

// TolerationList lists an app tolerations
func (d *DeisCmd) TolerationList(appID string, format string) error {
	settings, appID, err := load(d.ConfigFile, appID)
	if err != nil {
		return err
	}

	config, err := config.List(settings.Client, appID)
	if d.checkAPICompatibility(settings.Client, err) != nil {
		return err
	}
	var configOutput strings.Builder

	appTypes := make([]string, 0, len(config.Tolerations))
	for k := range config.Tolerations {
		appTypes = append(appTypes, k)
	}

	sort.Strings(appTypes)

	for _, appType := range appTypes {
		tolerations := config.Tolerations[appType]
		var identifiers []string
		for identifier := range tolerations {
			identifiers = append(identifiers, identifier)
		}
		sort.Strings(identifiers)
		switch format {
		case "oneline":
            configOutput.WriteString(fmt.Sprintf("%s:", appType))
			for _, identifier := range identifiers {
				toleration := tolerations[identifier]
				configOutput.WriteString(fmt.Sprintf(" %s|", identifier))
				configOutput.WriteString(fmt.Sprintf("Key=%s", toleration.Key))
				configOutput.WriteString(fmt.Sprintf(",Operator=%s", toleration.Operator))
				configOutput.WriteString(fmt.Sprintf(",Value=%s", toleration.Value))
				configOutput.WriteString(fmt.Sprintf(",Effect=%s", toleration.Effect))
				configOutput.WriteString(fmt.Sprintf(",TolerationSeconds=%d", *toleration.TolerationSeconds))
			}
			configOutput.WriteString(fmt.Sprintf(";\n"))
		case "diff":
			configOutput.WriteString(fmt.Sprintf("%s:\n", appType))
			for _, identifier := range identifiers {
				toleration := tolerations[identifier]
				configOutput.WriteString(fmt.Sprintf("---- %s\n", identifier))
				configOutput.WriteString(fmt.Sprintf("    Key=%s\n", toleration.Key))
				configOutput.WriteString(fmt.Sprintf("    Operator=%s\n", toleration.Operator))
				configOutput.WriteString(fmt.Sprintf("    Value=%s\n", toleration.Value))
				configOutput.WriteString(fmt.Sprintf("    Effect=%s\n", toleration.Effect))
				configOutput.WriteString(fmt.Sprintf("    TolerationSeconds=%d\n", *toleration.TolerationSeconds))
			}
		default:
			configOutput.WriteString(fmt.Sprintf("=== %s Tolerations\n", appType))
			for _, identifier := range identifiers {
				toleration := tolerations[identifier]
				configOutput.WriteString(fmt.Sprintf("---- %s\n", identifier))
				var output = make(map[string]string, 5)
				if toleration.Key != "" {
					output["Key"] = toleration.Key
				}
				output["Operator"] = fmt.Sprintf("%s", toleration.Operator)
				if toleration.Value != "" {
					output["Value"] = toleration.Value
				}
				if toleration.Effect != "" {
					output["Effect"] = fmt.Sprintf("%s", toleration.Effect)
				}
				output["Toleration Seconds"] = fmt.Sprintf("%d", *toleration.TolerationSeconds)

				configOutput.WriteString(fmt.Sprintf(prettyprint.PrettyTabs(output, 6)))
			}
		}
	}

	_, err = d.Print(configOutput.String())
	return err
}

// TolerationSet sets an app's tolerations.
func (d *DeisCmd) TolerationSet(appID string, appType string, identifier string, toleration v1.Toleration) error {
	settings, appID, err := load(d.ConfigFile, appID)

	if err != nil {
		return err
	}

	d.Print("Creating Tolerations... ")

	quit := progress(d.WOut)

	appTolerations := make(map[string]map[string]*v1.Toleration)
	appTolerations[appType] = map[string]*v1.Toleration{identifier: &toleration}
	configObj := api.Config{Tolerations: appTolerations}
	configObj, err = config.Set(settings.Client, appID, configObj)

	quit <- true
	<-quit
	if d.checkAPICompatibility(settings.Client, err) != nil {
		return err
	}

	d.Print("done\n\n")

	return d.TolerationList(appID, "")
}

// TolerationUnset removes an toleration from an app.
func (d *DeisCmd) TolerationUnset(appID string, appType string, tolerationIdentifiers []string) error {
	settings, appID, err := load(d.ConfigFile, appID)

	if err != nil {
		return err
	}

	d.Print("Removing Tolerations... ")

	quit := progress(d.WOut)

	valuesMap := make(map[string]*v1.Toleration)
	for _, identifier := range tolerationIdentifiers {
		valuesMap[identifier] = nil
	}

	appTolerations := make(map[string]map[string]*v1.Toleration)
	appTolerations[appType] = valuesMap
	configObj := api.Config{Tolerations: appTolerations}
	configObj, err = config.Set(settings.Client, appID, configObj)

	quit <- true
	<-quit
	if d.checkAPICompatibility(settings.Client, err) != nil {
		return err
	}

	d.Print("done\n\n")

	return d.TolerationList(appID, "")
}

