package main_test

// func getRecipeFromFile() (*brewfather.Recipe, error) {
// 	jsonFile, err := os.Open("./model/recipe.json")
// 	if err != nil {
// 		return nil, errors.Wrap(err, "Could not open file")
// 	}

// 	defer jsonFile.Close()

// 	var recipe *brewfather.Recipe

// 	if err = json.NewDecoder(jsonFile).Decode(&recipe); err != nil {
// 		return nil, errors.Wrap(err, "Could not decode Recipe")
// 	}

// 	return recipe, nil
// }
