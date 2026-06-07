// Package fzfpreview implements the fzf-image-preview command.
package fzfpreview

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/dmccaffery/dotfiles/internal/cmd/cmdutil"
	"github.com/dmccaffery/dotfiles/internal/envx"
	"github.com/dmccaffery/dotfiles/internal/execx"
)

// NewCmd builds the fzf-image-preview command.
func NewCmd(deps *cmdutil.Deps) *cobra.Command {
	return &cobra.Command{
		Use:   "fzf-image-preview <file>",
		Short: "Preview handler for fzf — chafa for images, bat for text",
		Long: "An fzf --preview handler: directories list, images render via chafa, other binaries\n" +
			"report their type, and text renders via bat (first 100 lines).",
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			out := cmd.OutOrStdout()
			file := cmdutil.Arg(args, 0)
			if file == "" {
				fmt.Fprint(out, "usage: fzf-image-preview <FILE_PATH>\n\n")
				return cmdutil.ErrSilent
			}

			ctx := cmd.Context()
			r := deps.Runner

			if cmdutil.DirExists(file) {
				_ = r.RunIO(ctx, cmdutil.Streams(cmd), "ls", "-la", "--color", file)
				return nil
			}
			if !cmdutil.FileExists(file) {
				fmt.Fprint(out, "file does not exist!\n\n")
				return cmdutil.ErrSilent
			}

			mime := ""
			if res, err := r.Run(ctx, "file", "--mime", file); err == nil {
				mime = res.Stdout
			}

			switch {
			case strings.Contains(mime, "binary") && strings.Contains(mime, "image/"):
				pass := "auto"
				if deps.Env.Get("TMUX") != "" {
					pass = "tmux"
				}
				_ = r.RunIO(ctx, cmdutil.Streams(cmd), "chafa",
					"--passthrough="+pass, "--size="+previewColumns(ctx, r, deps.Env), file)
			case strings.Contains(mime, "binary"):
				fmt.Fprintf(out, "%s is a binary file\n", file)
			default:
				batHead(ctx, r, out, file)
			}
			return nil
		},
	}
}

// previewColumns mirrors `${FZF_PREVIEW_COLUMNS:-$(($(tput cols) / 2))}`.
func previewColumns(ctx context.Context, r execx.Runner, env envx.Env) string {
	if c := env.Get("FZF_PREVIEW_COLUMNS"); c != "" {
		return c
	}
	if res, err := r.Run(ctx, "tput", "cols"); err == nil {
		if n := cmdutil.Atoi(res.Stdout); n > 0 {
			return strconv.Itoa(n / 2)
		}
	}
	return ""
}

// batHead renders a text file with bat and prints the first 100 lines, mirroring
// `bat --style=numbers --color=always FILE 2>/dev/null | head -100`.
func batHead(ctx context.Context, r execx.Runner, out io.Writer, file string) {
	res, _ := r.Run(ctx, "bat", "--style=numbers", "--color=always", file)
	sc := bufio.NewScanner(strings.NewReader(res.Stdout))
	sc.Buffer(make([]byte, 0, 64*1024), 4*1024*1024)
	for i := 0; i < 100 && sc.Scan(); i++ {
		fmt.Fprintln(out, sc.Text())
	}
}
