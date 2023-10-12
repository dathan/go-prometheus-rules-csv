package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

type PrometheusRules struct {
	KubePrometheusStack struct {
		AdditionalPrometheusRulesMap map[string]Groups `yaml:"additionalPrometheusRulesMap"`
	} `yaml:"kube-prometheus-stack"`
}

type Groups struct {
	Groups []Group `yaml:"groups"`
}

type Group struct {
	Name  string `yaml:"name"`
	Rules []Rule `yaml:"rules"`
}

type Rule struct {
	Alert      string            `yaml:"alert"`
	Expr       string            `yaml:"expr"`
	For        string            `yaml:"for"`
	Labels     map[string]string `yaml:"labels"`
	Annotations map[string]string `yaml:"annotations"`
}

func main() {
	file, err := os.ReadFile("prod.yaml")
	if err != nil {
		panic(err)
	}

	var rules PrometheusRules
	if err := yaml.Unmarshal(file, &rules); err != nil {
		panic(err)
	}

	csvFile, err := os.Create("output.csv")
	if err != nil {
		panic(err)
	}
	defer csvFile.Close()

	writer := csv.NewWriter(csvFile)
	defer writer.Flush()

	writer.Write([]string{"Name", "Alert", "Expr", "For", "Labels", "Annotations"})

	for ruleName, group := range rules.KubePrometheusStack.AdditionalPrometheusRulesMap {
		for _, ruleGroup := range group.Groups {
			for _, rule := range ruleGroup.Rules {
				labels := "None"
				if len(rule.Labels) > 0 {
					var labelParts []string
					for k, v := range rule.Labels {
						labelParts = append(labelParts, fmt.Sprintf("%s:%s", k, v))
					}
					labels = strings.Join(labelParts, ", ")
				}

				annotations := "None"
				if len(rule.Annotations) > 0 {
					var annotationParts []string
					for k, v := range rule.Annotations {
						annotationParts = append(annotationParts, fmt.Sprintf("%s:%s", k, v))
					}
					annotations = strings.Join(annotationParts, ", ")
				}

				writer.Write([]string{
					ruleName,
					rule.Alert,
					rule.Expr,
					rule.For,
					labels,
					annotations,
				})
			}
		}
	}
}

