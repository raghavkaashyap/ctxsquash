package cli

import (
	"fmt"
	"io"
	"os"

	"github.com/raghavkaashyap/ctxsquash/internal/squash"
	"github.com/spf13/cobra"
)

type config struct {
	output      string
	include     string
	exclude     string
	treeOnly    bool
	stdout      bool
	maxFileSize int64
}

func NewRootCommand() *cobra.Command {
	return newRootCommand(os.Stdout, os.Stderr)
}

func newRootCommand(stdout io.Writer, stderr ...io.Writer) *cobra.Command {
	var cfg config
	warningOutput := io.Discard
	if len(stderr) > 0 && stderr[0] != nil {
		warningOutput = stderr[0]
	}

	cmd := &cobra.Command{
		Use:   "ctxsquash [path]",
		Short: "Convert a repository or folder into Markdown context",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			options := squash.Options{
				Root:        args[0],
				Output:      cfg.output,
				Include:     squash.SplitCSV(cfg.include),
				Exclude:     squash.SplitCSV(cfg.exclude),
				TreeOnly:    cfg.treeOnly,
				MaxFileSize: cfg.maxFileSize,
			}

			result, warnings, err := squash.RenderWithWarnings(options)
			if err != nil {
				return err
			}
			writeSecretWarnings(warningOutput, warnings)

			if cfg.stdout || cfg.output == "" {
				_, err = fmt.Fprint(stdout, result)
				return err
			}

			return os.WriteFile(cfg.output, []byte(result), 0644)
		},
	}

	cmd.Flags().StringVarP(&cfg.output, "output", "o", "", "write Markdown to a file")
	cmd.Flags().StringVar(&cfg.include, "include", "", "comma-separated file extensions to include")
	cmd.Flags().StringVar(&cfg.exclude, "exclude", "", "comma-separated directories to exclude")
	cmd.Flags().BoolVar(&cfg.treeOnly, "tree-only", false, "only render the project tree")
	cmd.Flags().BoolVar(&cfg.stdout, "stdout", false, "print Markdown to stdout")
	cmd.Flags().Int64Var(&cfg.maxFileSize, "max-file-size", squash.DefaultMaxFileSize, "maximum file size in bytes to include")

	return cmd
}

func writeSecretWarnings(writer io.Writer, warnings []squash.SecretWarning) {
	for _, warning := range warnings {
		fmt.Fprintf(writer, "warning: possible %s in %s:%d; value not shown\n", warning.Kind, warning.Path, warning.Line)
	}
}
