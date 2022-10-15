package main

import (
	"fmt"
	"os"

	"github.com/cli/go-gh"
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
		fmt.Println(out)
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
	Run: func(cmd *cobra.Command, args []string) {
		out, err := ref(args[0])
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		fmt.Println(out)
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
	Run: func(cmd *cobra.Command, args []string) {
		sanitize, _ := cmd.Flags().GetBool("sanitize")
		out, err := title(args[0], sanitize)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		fmt.Println(out)
	},
}

func init() {
	linkCmd.Flags().BoolP("simple", "s", false, "Disable title lookup")
	titleCmd.Flags().BoolP("sanitize", "s", false, "Sanitize output for use as a file path")
	rootCmd.AddCommand(linkCmd, refCmd, titleCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func link(input string, simple bool) (string, error) {
	client, err := gh.RESTClient(nil)
	if err != nil {
		return input, err
	}
	response := struct{ Login string }{}
	err = client.Get("user", &response)
	if err != nil {
		return input, err
	}
	return fmt.Sprintf("running as %s", response.Login), nil
}

func ref(input string) (string, error) {
	if reference, matched := match(input); matched {
		return reference.Reference(), nil
	}
	return input, nil
}

func title(input string, sanitize bool) (string, error) {
	return input, nil
}
