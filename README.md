# WIP

# go-change-delta

go-change-delta is a tool to detect the blast radius of changes in a Go-based repository. It helps you identify which packages are affected by changes in your codebase, enabling more focused and efficient testing.

## Features

- Branch Comparison: Compare changes against a specified base branch.
- Project Path: Specify the project path for subprojects.
- Exclusions: Exclude specific paths from the analysis.
- Recursion Levels: Control the depth of recursion to find affected packages.
- Test Dependencies: Optionally exclude test dependencies in the analysis.
- Embedded files are included in the comparisons

## Installation
To install go-change-delta, you can use go install:

``` bash
go install github.com/tateexon/go-change-delta@latest
```

## Usage
Run go-change-delta with the following flags:

```bash
./go-change-delta -b <branch> -p <project-path> -e <excludes> -l <levels> -t <include-test-deps>
```

### Flags

    -b (required): The base git branch to compare current changes with.
    -p (optional): The path to the project. Default is the current directory.
    -e (optional): Comma-separated list of paths to exclude.
    -l (optional): The number of levels of recursion to search for affected packages. Default is 2. Use 0 for unlimited recursion.
    -t (optional): Include test dependencies. Default is true.

### Example

With go test:

```bash
go test $(go-change-delta -b=origin/main -l=0)
```

This command checks all the go packages between your current branch and origin/main that are affected and runs their tests.

```bash
go-change-delta -b main -p ./my-subproject -e "vendor,third_party" -l 3 -t false
```

This command compares changes against the main branch, analyzes the project located at ./my-subproject, excludes the vendor and third_party directories, searches up to 3 levels of package dependencies, and excludes test dependencies from the analysis.

## Github Action Example

A github action has been provided at tateexon/go-change-delta/.github/actions/go-change-delta that will grab the packages delta and return it in the outputs. Example usage below:

```yaml
jobs:
  test:
    name: Run Go Tests
    runs-on: ubuntu-latest
    steps:
      - name: Checkout the Repo
        uses: actions/checkout@692973e3d937129bcbf40652eb9f2f61becf3332 # v4.1.7
        with:
          fetch-depth: 0 # needed for go-change-delta to work correctly
      - name: Install golang and other environment stuff needed for your tests
        ...
      - name: Get Affected Packages
        uses: tateexon/go-change-delta/.github/actions/go-change-delta
        id: delta
      - name: Run Test
        if: steps.delta.outputs.packages != ''
        run: go test ${{ steps.delta.outputs.packages }}
```

## Contributing
Contributions are welcome! Feel free to open an issue or submit a pull request.

# TODO:

- Setup and add CI
- Add a github action to make using this easier in other CIs
