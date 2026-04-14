package env

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

// WriteOptions configures how the .env file is written.
type WriteOptions struct {
	OutputPath string
	Overwrite  bool
	Append     bool
}

// WriteEnvFile writes the given secrets map to a .env file.
// Keys are sorted alphabetically for deterministic output.
func WriteEnvFile(secrets map[string]string, opts WriteOptions) error {
	if opts.OutputPath == "" {
		return fmt.Errorf("output path must not be empty")
	}

	if !opts.Overwrite && !opts.Append {
		if _, err := os.Stat(opts.OutputPath); err == nil {
			return fmt.Errorf("file %q already exists; use --overwrite or --append", opts.OutputPath)
		}
	}

	flag := os.O_CREATE | os.O_WRONLY
	if opts.Append {
		flag |= os.O_APPEND
	} else {
		flag |= os.O_TRUNC
	}

	f, err := os.OpenFile(opts.OutputPath, flag, 0600)
	if err != nil {
		return fmt.Errorf("opening env file: %w", err)
	}
	defer f.Close()

	keys := make([]string, 0, len(secrets))
	for k := range secrets {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var sb strings.Builder
	for _, k := range keys {
		sb.WriteString(fmt.Sprintf("%s=%s\n", k, quoteValue(secrets[k])))
	}

	_, err = f.WriteString(sb.String())
	if err != nil {
		return fmt.Errorf("writing env file: %w", err)
	}

	return nil
}

// quoteValue wraps values containing spaces or special characters in double quotes.
func quoteValue(v string) string {
	if strings.ContainsAny(v, " \t\n#") {
		v = strings.ReplaceAll(v, `"`, `\"`)
		return `"` + v + `"`
	}
	return v
}
