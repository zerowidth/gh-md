package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/cli/go-gh"
	"github.com/cli/go-gh/pkg/api"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Short: "GitHub markdown link tools",
}

var linkCmd = &cobra.Command{
	Use:   "link <url>",
	Short: "Convert an input URL into a markdown link",
	Long: `Convert an input URL into a markdown link.

The input can be one of:

   * a GitHub repository URL
   * a GitHub issue URL
   * a GitHub pull request URL
   * a GitHub discussion URL
   * a GitHub issue reference (e.g. "cli/cli#123")
   * a markdown link containing one of the above URL types

The output is a markdown link to the input URL, with the link text being the
issue/PR/discussion reference as well as its title fetched from the GitHub API.

Include the --simple flag to disable title lookup.

If the input is unrecognized it will be returned as-is.

Examples:

   $ gh md link https://github.com/cli/cli/pull/123
   [cli/cli#123: Tweak flags language](https://github.com/cli/cli/pull/123)

   $ gh md link --simple https://github.com/cli/cli/pull/123
   [cli/cli#123](https://github.com/cli/cli/pull/123)

   $ gh md link cli/cli#123
   [cli/cli#123: Tweak flags language](https://github.com/cli/cli/pull/123)`,

	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		simple, _ := cmd.Flags().GetBool("simple")
		out, err := link(args[0], simple)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		printOptionalNewline(cmd, out)
	},
}

var refCmd = &cobra.Command{
	Use:     "ref <url>",
	Aliases: []string{"reference"},
	Short:   "Convert an input URL into an issue reference",
	Long: `Convert an input URL into an issue reference.

The input can be one of:

   * A GitHub issue URL
   * A GitHub pull request URL
   * A GitHub discussion URL
   * An issue reference ("cli/cli#123")
   * A markdown link containing one of the above

Example:

   $ gh md ref https://github.com/cli/cli/pull/123
   cli/cli#123`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		out, err := ref(args[0])
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		printOptionalNewline(cmd, out)
	},
}

var titleCmd = &cobra.Command{
	Use:   "title <url>",
	Short: "Fetch the title of a GitHub issue, pull request, or discussion",
	Long: `Fetch the title of a GitHub issue, pull request, or discussion.

The input can be one of:

   * A GitHub issue URL
   * A GitHub pull request URL
   * A GitHub discussion URL
   * An issue reference ("cli/cli#123")
   * A markdown link containing one of the above

Example:

   $ gh md title cli/cli#123
   Tweak flags language`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		sanitize, _ := cmd.Flags().GetBool("sanitize")
		out, err := title(args[0], sanitize)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		printOptionalNewline(cmd, out)
	},
}

var urlCmd = &cobra.Command{
	Use:   "url <ref or markdown link>",
	Short: "Convert an issue reference or markdown link into a bare GitHub URL",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		out, err := url(args[0])
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		printOptionalNewline(cmd, out)
	},
}

var client api.GQLClient

func init() {
	for _, cmd := range []*cobra.Command{linkCmd, refCmd, titleCmd, urlCmd} {
		cmd.Flags().BoolP("no-newline", "n", false, "Do not print trailing newline")
	}
	linkCmd.Flags().Bool("simple", false, "Disable title lookup")
	titleCmd.Flags().Bool("sanitize", false, "Sanitize output for use as a file path")
	rootCmd.AddCommand(linkCmd, refCmd, titleCmd, urlCmd)
}

func main() {
	var err error
	opts := &api.ClientOptions{
		EnableCache: true,
		Timeout:     10 * time.Second,
	}
	client, err = gh.GQLClient(opts)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func link(input string, simple bool) (string, error) {
	reference, matched := match(input)
	if !matched {
		return input, nil
	}

	ref := reference.Reference()
	if simple {
		return fmt.Sprintf("[%s](%s)", ref, reference.URL()), nil
	} else {
		title, err := reference.Title(client)
		if err != nil {
			return input, err
		}

		if strings.Count(title, "::") > 1 {
			title = strings.ReplaceAll(title, "::", "-")
		}

		return fmt.Sprintf("[%s: %s](%s)", ref, title, reference.URL()), nil
	}
}

func ref(input string) (string, error) {
	reference, matched := match(input)
	if !matched {
		return input, nil
	}

	return reference.Reference(), nil
}

func title(input string, sanitize bool) (string, error) {
	reference, matched := match(input)
	if !matched {
		fmt.Fprintf(os.Stderr, "didn't match %s", input)
		return input, nil
	}

	title, err := reference.Title(client)
	if err != nil {
		return input, err
	}

	// sanitize for paths, remove : / ? chars
	if sanitize {
		title = regexp.MustCompile(`:+|/`).ReplaceAllString(title, " - ")
		title = strings.ReplaceAll(title, "/", "-")
		title = strings.ReplaceAll(title, "?", "")
	}
	title = regexp.MustCompile(`\s+`).ReplaceAllString(title, " ")
	title = strings.ReplaceAll(title, "[", "(")
	title = strings.ReplaceAll(title, "]", ")")
	if strings.Count(title, "::") > 1 {
		title = strings.ReplaceAll(title, "::", "-")
	}

	return title, nil
}

func url(input string) (string, error) {
	reference, matched := match(input)
	if !matched {
		return input, nil
	}

	return reference.URL(), nil
}

func printOptionalNewline(cmd *cobra.Command, output string) {
	noNewline, _ := cmd.Flags().GetBool("no-newline")
	if noNewline {
		fmt.Print(output)
	} else {
		fmt.Println(output)
	}
}
