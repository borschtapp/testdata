# Testdata Package

This directory contains the test data and utilities for the Krip recipe scraper.
It is organized as a **separate repository** to maintain a clean separation between the core library and its test fixtures.

Probably I will regret this later, but at some point I decided to separate the testdata to isolate large files and dependencies from the main package.
If you have better idea how to organize that, let me know.

Few points for:

- **Isolation**: Keeps large HTML and JSON files out of the main code
- **Separate Tools**: Development utilities for managing test data are included here
- **Separate Dependencies**: Allows test-only dependencies without bloating the main package (only 1 as of now)

Few points against:

- **Maintainability**: Requires extra steps to keep test data in sync with main code
- **Discoverability**: Will be confusing if someone will want to contribute a new scraper, do they need to submit two PRs?

## Directory Structure

```text
testdata/
├── testdata_test.go       # Integration tests for all scrapers
├── ingredients.txt        # Collected ingredients corpus
├── schema_paths.txt       # Schema.org property paths found in test data
├── recipes/               # Expected JSON output for each website
│   ├── allrecipes.json
│   ├── bbcgoodfood.json
│   └── ...                # 200+ recipe fixture files
├── websites/              # Raw HTML snapshots from real websites
│   ├── allrecipes.html
│   ├── bbcgoodfood.html
│   └── ...                # 200+ website HTML files
├── schema/                # Schema.org/Microdata test cases
└── tools/                 # Development utilities
    ├── websites/
    │   └── manage.go      # Use it to add or update websites fixtures
    ├── schema/
    │   └── paths.go       # Analyze Schema.org paths in test data
    └── ingredients/
        └── collect.go     # Extract all ingredients from recipes
```

## Usage

The important part here is `replace github.com/borschtapp/krip => ../` in your `go.mod`,
which tells Go to use the local version of the `krip` module for testing.
