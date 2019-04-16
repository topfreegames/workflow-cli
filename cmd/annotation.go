package cmd

import (
	"bytes"
	"fmt"
	"github.com/deis/pkg/prettyprint"
	"regexp"
	"sort"

	"github.com/deis/controller-sdk-go/api"
	"github.com/deis/controller-sdk-go/config"
)

// AnnotationList lists an app annotations
func (d *DeisCmd) AnnotationList(appID string, format string) error {
	settings, appID, err := load(d.ConfigFile, appID)

	if err != nil {
		return err
	}

	config, err := config.List(settings.Client, appID)
	if d.checkAPICompatibility(settings.Client, err) != nil {
		return err
	}
	var configOutput = new(bytes.Buffer)

	appTypes := make([]string, 0, len(config.Annotations))
	for k := range config.Annotations {
		appTypes = append(appTypes, k)
	}

	sort.Strings(appTypes)

	for _, appType := range appTypes {
		annotations := config.Annotations[appType]
		keys := sortKeys(annotations)
		switch format {
		case "oneline":
			fmt.Fprintf(configOutput, "%s:", appType)
			for _, key := range keys {
				value := annotations[key]
				fmt.Fprintf(configOutput, " %s=%s", key, value)
			}
			fmt.Fprintf(configOutput, "\n")
		case "diff":
			fmt.Fprintf(configOutput, "%s:\n", appType)
			for _, key := range keys {
				value := annotations[key]
				fmt.Fprintf(configOutput, "    %s=%s\n", key, value)
			}
		default:
			fmt.Fprintf(configOutput, "=== %s Annotations\n", appType)
			prettyPrintAnnotations := make(map[string]string)
            for _, key := range keys {
				value := annotations[key]
				prettyPrintAnnotations[key] = value.(string)
			}
			fmt.Fprint(configOutput, prettyprint.PrettyTabs(prettyPrintAnnotations, 6))
		}
	}

	d.Print(configOutput)
	return nil
}

// AnnotationSet sets an app's annotations.
func (d *DeisCmd) AnnotationSet(appID string, appType string, annotationCommandLine []string) error {
	settings, appID, err := load(d.ConfigFile, appID)

	if err != nil {
		return err
	}

	annotations, err := parseAnnotations(annotationCommandLine)
	if err != nil {
		return err
	}

	d.Print("Creating Annotations... ")

	quit := progress(d.WOut)

	appAnnotations := make(map[string]api.Annotation)
	appAnnotations[appType] = annotations
	configObj := api.Config{Annotations: appAnnotations}
	configObj, err = config.Set(settings.Client, appID, configObj)

	quit <- true
	<-quit
	if d.checkAPICompatibility(settings.Client, err) != nil {
		return err
	}

	d.Print("done\n\n")

	return d.AnnotationList(appID, "")
}

// AnnotationUnset removes an annotation from an app.
func (d *DeisCmd) AnnotationUnset(appID string, appType string, annotations []string) error {
	s, appID, err := load(d.ConfigFile, appID)

	if err != nil {
		return err
	}

	d.Print("Removing Annotations... ")

	quit := progress(d.WOut)

	configObj := api.Config{}

	valuesMap := make(map[string]interface{})

	for _, configVar := range annotations {
		valuesMap[configVar] = nil
	}

	annotationMap := make(map[string]api.Annotation)
	annotationMap[appType] = valuesMap
	configObj.Annotations = annotationMap

	_, err = config.Set(s.Client, appID, configObj)
	quit <- true
	<-quit
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	d.Print("done\n\n")

	return d.AnnotationList(appID, "")
}

func parseAnnotations(annotations []string) (api.Annotation, error) {
	annotationsMap := make(api.Annotation)

	regex := regexp.MustCompile(`^([A-Za-z_]+.*)=([\s\S]+)$`)
	for _, annotation := range annotations {
		if len(annotation) > 0 && annotation[0] == '#' {
			continue
		}
		if regex.MatchString(annotation) {
			captures := regex.FindStringSubmatch(annotation)
			annotationsMap[captures[1]] = captures[2]
		} else {
			return nil, fmt.Errorf("'%s' does not match the pattern 'key=var', ex: MODE=test\n", annotation)
		}
	}

	return annotationsMap, nil
}
