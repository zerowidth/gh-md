# gh-md: markdown link and issue reference generation

This [GitHub CLI](https://cli.github.com) extension generates markdown links from input URLs or references.

This isn't particularly useful by itself, but combine it with tools like [Alfred](https://alfred.app) snippet triggers, [Raycast](https://www.raycast.com) scripts, or [Obsidian](https://obsidian.md) templates, you can automate these links for insertion into your markdown documents.

## Installation

`gh extension install zerowidth/gh-md`

## Usage

`gh md --help` for full details.

### `gh md link`

Generates a markdown link from an input URL or issue/pr/discussion reference.

Basic example:

```
$ gh md link https://github.com/cli/cli/pull/123
[cli/cli#123: Tweak flags language](https://github.com/cli/cli/pull/123)
```

Skip title lookup:

```
$ gh md link --simple https://github.com/cli/cli/pull/123
[cli/cli#123](https://github.com/cli/cli/pull/123)
```

Create a link from an issue reference:

```
$ gh md link cli/cli#123
[cli/cli#123: Tweak flags language](https://github.com/cli/cli/pull/123)
```

### `gh md ref`

Generates an issue/pr/discussion reference from an input URL or reference.

```
$ gh md ref https://github.com/cli/cli/pull/123
cli/cli#123
```

### `gh md title`

Looks up the title of the given URL or reference:

```
$ gh md title cli/cli#123
Tweak flags language
```

The title can be sanitized for use as a path, stripping `:`, `/`, and a few other characters.

### `gh md url`

Generates the URL from the given issue reference or markdown link.

```
$ gh md url cli/cli#123
https://github.com/cli/cli/pull/123
```

## Using `gh-md` in [Obsidian](https://obsidian.md)

`gh-md` was built in part to be used in user-defined functions in the [Templater](https://github.com/SilentVoid13/Templater) plugin.

### Markdown links

- define a `markdownLink` user function as `/path/to/gh md link -n "$input"`
- copy a GitHub URL to your clipboard
- use it in a template: `<% tp.user.markdownLink({input: (await tp.system.clipboard())}) %>`

### Sanitized titles

The sanitized issue title can be used for filenames. In a templater template for creating a document for an issue from the clipboard:

- Define a `markdownTitle` user function as `/path/to/gh md title -n --sanitize "$input"`.
- Use the title to rename a new file from that template:

  ```
  <%- await tp.file.rename("Issue - " + tp.date.now("YYYY-MM") + " - " + (await tp.user.markdownTitle({input: (await tp.system.clipboard())}))) -%>
  ```

### Auth issues

If you're using a recent version of the GitHub CLI that uses the system keychain to store credentials, these helper functions may not work. If so, you can prepend the `PATH` to the user functions:

- `PATH=$PATH:/opt/homebrew/bin /opt/homebrew/bin/gh md link -n "$input"`
- `PATH=$PATH:/opt/homebrew/bin /opt/homebrew/bin/gh md title -n --sanitize "$input"`
