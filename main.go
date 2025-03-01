package main

import (
	"ds-village-calculator/TwStatsScraper"
	"fmt"
	"os"
	"strconv"
)

func main() {
	arguments := os.Args
	if len(arguments) != 3 {
		fmt.Println("Please provide the world and village id as arguments")
		fmt.Println("Scheme: TwStatScraper.exe world villageId")
		fmt.Println("Example: TwStatScraper.exe de238 4165")
		return
	}

	world := arguments[1]
	villageId, err := strconv.Atoi(arguments[2])
	if err != nil {
		fmt.Println("Invalid village id:", err.Error())
	}
	buildings := getBuildings(world)
	villageExpansionTree := NewVillageExpansionTree(buildings...)
	villageHistory := TwStatsScraper.ScrapeVillageHistory(world, villageId)

	possibleBuildingCompositions := findPossibleBuildingCompositions(villageExpansionTree, villageHistory)
	for _, village := range possibleBuildingCompositions {
		for building, level := range village {
			fmt.Printf("%s: %d", building, level)
		}
	}
}

func findPossibleBuildingCompositions(tree VillageExpansionTree, history []TwStatsScraper.VillageHistoryEntry) []map[string]int {
	//TODO: implement
	panic("Not yet implemented")
}
