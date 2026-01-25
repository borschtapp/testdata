package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"path/filepath"

	"github.com/borschtapp/krip"
	"github.com/borschtapp/krip/model"
	"github.com/borschtapp/krip/scraper"
	"github.com/borschtapp/krip/utils"

	"github.com/borschtapp/testdata"
)

// use as: go run tools/websites/manage.go add <url>
func main() {
	flag.Usage = func() {
		_, _ = fmt.Fprintf(os.Stderr, "Usage: %s <command> [arguments]\n", os.Args[0])
		_, _ = fmt.Fprint(os.Stderr, "\nCommands:\n")
		_, _ = fmt.Fprint(os.Stderr, "  add <url>         Adds a new recipe. Fails if the alias already exists.\n")
		_, _ = fmt.Fprint(os.Stderr, "  update <url|alias> Updates an existing recipe by its URL or alias.\n")
		_, _ = fmt.Fprint(os.Stderr, "  update --all      Updates all existing recipes in the test catalog.\n")
		_, _ = fmt.Fprint(os.Stderr, "  clean             Removes recipes without websites and .new files.\n")
	}

	if len(os.Args) < 2 {
		flag.Usage()
		os.Exit(1)
	}

	cmd := os.Args[1]
	switch cmd {
	case "add":
		addCmd := flag.NewFlagSet("add", flag.ExitOnError)
		addCmd.Parse(os.Args[2:])
		if addCmd.NArg() != 1 {
			_, _ = fmt.Fprintln(os.Stderr, "Usage: manage add <url>")
			os.Exit(1)
		}
		url := addCmd.Arg(0)
		if err := saveTestdata(url, false); err != nil {
			log.Fatal(err)
		}
		fmt.Println("Done!")

	case "update":
		updateCmd := flag.NewFlagSet("update", flag.ExitOnError)
		allFlag := updateCmd.Bool("all", false, "Update all existing websites")
		updateCmd.Parse(os.Args[2:])

		if *allFlag {
			updateAll()
			fmt.Println("Done, all websites updated!")
			return
		}

		if updateCmd.NArg() != 1 {
			_, _ = fmt.Fprintln(os.Stderr, "Usage: manage update <url|alias> OR manage update --all")
			os.Exit(1)
		}

		arg := updateCmd.Arg(0)
		if strings.HasPrefix(arg, "http") {
			if err := saveTestdata(arg, true); err != nil {
				log.Fatal(err)
			}
		} else {
			// It's an alias
			input, err := scraper.FileInput(testdata.WebsitesDir+arg+testdata.HtmlExt, model.InputOptions{SkipSchema: true})
			if err != nil {
				log.Fatal("Unable to read old testdata: " + err.Error())
			}
			fmt.Println("Updating " + input.Url)
			if err := saveTestdata(input.Url, true); err != nil {
				log.Fatal(err)
			}
		}
		fmt.Println("Done!")

	case "clean":
		cleanAll()
		fmt.Println("Done, testdata cleaned!")

	default:
		flag.Usage()
		os.Exit(1)
	}
}

func updateAll() {
	testdata.WalkTestdataWebsites(func(name string, path string) {
		input, err := scraper.FileInput(path, model.InputOptions{SkipSchema: true})
		if err != nil {
			log.Printf("Unable to read old testdata (%s): %v\n", path, err)
			return
		}

		fmt.Println("Updating " + input.Url)
		if err := saveTestdata(input.Url, true); err != nil {
			log.Printf("Failed to update %s: %v\n", input.Url, err)
		}
	})
}

func saveTestdata(url string, force bool) error {
	alias := utils.HostAlias(url)
	websiteFileName := testdata.WebsitesDir + alias + testdata.HtmlExt
	recipeFileName := testdata.RecipesDir + alias + testdata.JsonExt

	if !force {
		if _, err := os.Stat(websiteFileName); err == nil {
			return fmt.Errorf("testdata already exists for the alias: %s. Use `update` to overwrite", alias)
		}
	}

	input, err := scraper.UrlInput(url)
	if err != nil {
		return fmt.Errorf("unable to fetch content: %w", err)
	}

	content := []byte(utils.TrimZeroWidthSpaces(input.Text)) // remove zero-width spaces, http.DetectContentType() doesn't like them
	if err = os.WriteFile(websiteFileName, content, 0644); err != nil {
		return fmt.Errorf("unable to create file: %w", err)
	}

	recipe, err := krip.Scrape(input)
	if err != nil {
		return fmt.Errorf("unable to scrape recipe: %w", err)
	}

	if err = os.WriteFile(recipeFileName, []byte(recipe.String()), 0644); err != nil {
		return fmt.Errorf("unable to create recipe file: %w", err)
	}

	return nil
}

func cleanAll() {
	testdata.WalkTestdataRecipes(func(name string, recipe model.Recipe) {
		alias := strings.TrimSuffix(name, testdata.JsonExt)
		websitePath := testdata.WebsitesDir + alias + testdata.HtmlExt
		if _, err := os.Stat(websitePath); os.IsNotExist(err) {
			recipePath := testdata.RecipesDir + name
			fmt.Println("Removing recipe without website: " + name)
			if err := os.Remove(recipePath); err != nil {
				log.Printf("Failed to remove recipe %s: %v\n", recipePath, err)
			}
		}
	})

	// Remove .new files
	_ = filepath.Walk(testdata.RecipesDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), testdata.NewExt) {
			fmt.Println("Removing .new file: " + info.Name())
			if err := os.Remove(path); err != nil {
				log.Printf("Failed to remove .new file %s: %v\n", path, err)
			}
		}
		return nil
	})
}
