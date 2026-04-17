package vault

import (
	"bytes"
	"fmt"
	"os"
	"text/template"
)

// RenderTemplate renders a Go template file using secrets from Vault
// and writes the output to destPath.
func RenderTemplate(templatePath, destPath string, secrets map[string]interface{}) error {
	tmplBytes, err := os.ReadFile(templatePath)
	if err != nil {
		return fmt.Errorf("reading template %q: %w", templatePath, err)
	}

	tmpl, err := template.New("vault").Option("missingkey=error").Parse(string(tmplBytes))
	if err != nil {
		return fmt.Errorf("parsing template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, secrets); err != nil {
		return fmt.Errorf("executing template: %w", err)
	}

	if err := os.WriteFile(destPath, buf.Bytes(), 0600); err != nil {
		return fmt.Errorf("writing rendered file %q: %w", destPath, err)
	}

	return nil
}

// ListTemplateVars returns the set of variables referenced in a template file.
func ListTemplateVars(templatePath string) ([]string, error) {
	tmplBytes, err := os.ReadFile(templatePath)
	if err != nil {
		return nil, fmt.Errorf("reading template %q: %w", templatePath, err)
	}

	tmpl, err := template.New("vault").Parse(string(tmplBytes))
	if err != nil {
		return nil, fmt.Errorf("parsing template: %w", err)
	}

	var vars []string
	seen := map[string]bool{}
	for _, node := range tmpl.Tree.Root.Nodes {
		if a, ok := node.(*template.Template); ok {
			_ = a
		}
		s := fmt.Sprintf("%s", node)
		if !seen[s] {
			seen[s] = true
			vars = append(vars, s)
		}
	}
	return vars, nil
}
