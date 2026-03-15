package testdata

import (
	"os"
	"regexp"
	"sort"
	"strings"
	"testing"

	"github.com/borschtapp/krip"
	"github.com/borschtapp/krip/scraper"
	"github.com/borschtapp/krip/utils"
	"github.com/stretchr/testify/assert"
)

func TestFilenames(t *testing.T) {
	WalkTestdataWebsites(func(name string, path string) {
		t.Run(name, func(t *testing.T) {
			input, err := scraper.FileInput(path, krip.ScrapeOptions{SkipMicrodata: true})
			assert.NoError(t, err)
			assert.NotEmpty(t, input.Url)

			// some exceptions, where there is no canonical in html
			if name == "sunbasket.html" {
				return
			}

			expected := utils.HostAlias(input.Url)
			assert.NotRegexp(t, regexp.MustCompile(`^file://.+`), input.Url)
			assert.Equal(t, expected+HtmlExt, name, "Incorrect filename for "+input.Url)
		})
	})
}

func TestWebsiteFiles(t *testing.T) {
	MockRequests(t)
	t.Parallel()

	WalkTestdataWebsites(func(name string, path string) {
		t.Run(name, func(t *testing.T) {
			recipe, err := krip.ScrapeFile(path, krip.ScrapeOptions{})
			assert.NoError(t, err)

			AssertJson(t, recipe, RecipesDir+strings.TrimSuffix(name, HtmlExt))
		})
	})
}

func TestWebsitesOnline(t *testing.T) {
	t.Skip("Skip online tests")
	t.Parallel()

	WalkTestdataWebsites(func(name string, path string) {
		t.Run(name, func(t *testing.T) {
			input, err := scraper.FileInput(path, krip.ScrapeOptions{SkipMicrodata: true})
			assert.NoError(t, err)
			assert.NotEmpty(t, input.Url)

			recipe, err := krip.ScrapeUrl(input.Url, krip.ScrapeOptions{})
			assert.NoError(t, err)

			AssertJson(t, recipe, RecipesDir+strings.TrimSuffix(name, HtmlExt))
		})
	})
}

func TestWebsitesProperties(t *testing.T) {
	t.Parallel()
	t.Skip("TODO: doesn't work for all")

	var domains []string
	WalkTestdataRecipes(func(name string, recipe krip.Recipe) {
		t.Run(name, func(t *testing.T) {
			assert.True(t, recipe.IsValid())

			if !t.Failed() {
				domains = append(domains, strings.ReplaceAll(strings.ReplaceAll(recipe.Publisher.Url, "http://", "https://"), "//www.", "//"))
			}
		})
	})

	assert.NoError(t, updateReadme(domains))
}

func updateReadme(domains []string) error {
	readmeContent, err := os.ReadFile(PackageDir + "../README.md")
	if err != nil {
		return err
	}

	newReadme := ""
	for _, line := range strings.Split(string(readmeContent), "\n") {
		newReadme += line + "\n"

		if line == "[//]: # (This list is generated automatically, do not edit manually)" {
			break
		}
	}

	newReadme += "\n"
	sort.Strings(domains)
	for _, domain := range domains {
		newReadme += "* <" + domain + ">\n"
	}

	if string(readmeContent) != newReadme {
		return os.WriteFile(PackageDir+"../README.md", []byte(newReadme), 0644)
	}
	return nil
}
