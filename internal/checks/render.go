package checks

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
)

func WriteText(w io.Writer, report Report) error {
	if _, err := fmt.Fprintf(w, "Omarchy Forge check (schema %s, %s)\nTarget: %s\n\n", report.SchemaVersion, report.Compatibility, report.Target); err != nil {
		return err
	}
	if len(report.Findings) == 0 {
		_, err := fmt.Fprintln(w, "PASS  No findings.")
		return err
	}
	for _, f := range report.Findings {
		location := f.Path
		if f.Line > 0 {
			location = fmt.Sprintf("%s:%d", location, f.Line)
		}
		if _, err := fmt.Fprintf(w, "%s  %s  %s  [%s]\n", f.Severity, f.RuleID, location, f.Source); err != nil {
			return err
		}
		if _, err := fmt.Fprintf(w, "  %s\n  Fix: %s\n", f.Message, f.Remediation); err != nil {
			return err
		}
	}
	_, err := fmt.Fprintf(w, "\nSummary: %d error(s), %d warning(s), %d note(s)\n", report.Summary.Errors, report.Summary.Warnings, report.Summary.Notes)
	return err
}

func WriteJSON(w io.Writer, report Report) error {
	encoder := json.NewEncoder(w)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "  ")
	return encoder.Encode(report)
}

type sarifLog struct {
	Version string     `json:"version"`
	Schema  string     `json:"$schema"`
	Runs    []sarifRun `json:"runs"`
}
type sarifRun struct {
	Tool    sarifTool     `json:"tool"`
	Results []sarifResult `json:"results"`
}
type sarifTool struct {
	Driver sarifDriver `json:"driver"`
}
type sarifDriver struct {
	Name           string      `json:"name"`
	InformationURI string      `json:"informationUri"`
	Rules          []sarifRule `json:"rules"`
}
type sarifRule struct {
	ID               string       `json:"id"`
	ShortDescription sarifMessage `json:"shortDescription"`
}
type sarifResult struct {
	RuleID    string          `json:"ruleId"`
	Level     string          `json:"level"`
	Message   sarifMessage    `json:"message"`
	Locations []sarifLocation `json:"locations,omitempty"`
}
type sarifMessage struct {
	Text string `json:"text"`
}
type sarifLocation struct {
	PhysicalLocation sarifPhysical `json:"physicalLocation"`
}
type sarifPhysical struct {
	ArtifactLocation sarifArtifact `json:"artifactLocation"`
	Region           *sarifRegion  `json:"region,omitempty"`
}
type sarifArtifact struct {
	URI string `json:"uri"`
}
type sarifRegion struct {
	StartLine int `json:"startLine"`
}

func WriteSARIF(w io.Writer, report Report) error {
	ruleMessages := map[string]string{}
	results := make([]sarifResult, 0, len(report.Findings))
	for _, f := range report.Findings {
		ruleMessages[f.RuleID] = f.Message
		result := sarifResult{RuleID: f.RuleID, Level: sarifLevel(f.Severity), Message: sarifMessage{Text: f.Message + " Fix: " + f.Remediation}}
		if f.Path != "" {
			physical := sarifPhysical{ArtifactLocation: sarifArtifact{URI: f.Path}}
			if f.Line > 0 {
				physical.Region = &sarifRegion{StartLine: f.Line}
			}
			result.Locations = []sarifLocation{{PhysicalLocation: physical}}
		}
		results = append(results, result)
	}
	ids := make([]string, 0, len(ruleMessages))
	for id := range ruleMessages {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	rules := make([]sarifRule, 0, len(ids))
	for _, id := range ids {
		rules = append(rules, sarifRule{ID: id, ShortDescription: sarifMessage{Text: ruleMessages[id]}})
	}
	log := sarifLog{Version: "2.1.0", Schema: "https://json.schemastore.org/sarif-2.1.0.json", Runs: []sarifRun{{Tool: sarifTool{Driver: sarifDriver{Name: "Omarchy Forge", InformationURI: "https://github.com/omarchy-forge/forge", Rules: rules}}, Results: results}}}
	encoder := json.NewEncoder(w)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "  ")
	return encoder.Encode(log)
}

func sarifLevel(severity Severity) string {
	if severity == Error {
		return "error"
	}
	if severity == Warning {
		return "warning"
	}
	return "note"
}
