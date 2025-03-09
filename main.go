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
	villageIncreases := TwStatsScraper.ScrapeVillageHistory(world, villageId)
	possibleVillageNodes := PossibleVillageExpansions(buildings, villageIncreases)

	if len(possibleVillageNodes) < 50 {
		for id, village := range possibleVillageNodes {
			fmt.Printf("%i:", id)
			for buildingName, buildingLevel := range village.buildings {
				fmt.Printf("%s: %d\n", buildingName, buildingLevel)
			}
		}
	}

	fmt.Println(len(possibleVillageNodes))
}
