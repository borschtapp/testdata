package main

import (
	"log"
	"os"
	"strings"

	"github.com/borschtapp/krip/model"

	"github.com/borschtapp/testdata"
)

func main() {
	var ingredients []string

	testdata.WalkTestdataRecipes(func(name string, recipe model.Recipe) {
		if recipe.Ingredients != nil {
			ingredients = append(ingredients, recipe.Ingredients...)
		}
	})

	log.Printf("Done, collected %d ingredients.\n", len(ingredients))
	_ = os.WriteFile(testdata.PackageDir+"/ingredients.txt", []byte(strings.Join(ingredients, "\n")), 0644)
}
