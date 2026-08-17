package service

import "testing"

func TestUpdateRecipeCanClearSteps(t *testing.T) {
	recipes, _, _, _, _ := testServices()
	created, err := recipes.Create(sampleRecipe("Soup"))
	if err != nil {
		t.Fatal(err)
	}
	replacement := sampleRecipe("Soup without instructions")
	replacement.Steps = []string{}
	updated, err := recipes.Update(created.ID, replacement)
	if err != nil {
		t.Fatal(err)
	}
	if len(updated.Steps) != 0 {
		t.Fatalf("steps were not cleared: %v", updated.Steps)
	}
}
