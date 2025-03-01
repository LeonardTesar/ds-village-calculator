package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

type VillageExpansionTree struct {
}

func NewVillageExpansionTree(building ...Building) VillageExpansionTree {
	buildingLevelPoints := readBuildingLevelPoints()
	buildingRequirements := readBuildingRequirements()
	fmt.Println(buildingLevelPoints)
	fmt.Println(buildingRequirements)
	//TODO: implement
	panic("Not yet implemented")
}

func readBuildingRequirements() map[string]map[string]int {
	filePath := filepath.Join("resources", "building_requirements.json")
	buildingRequirementsBytes, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("error reading building requirements from file: %s\n", err)
	}

	var buildingRequirements map[string]map[string]int
	err = json.Unmarshal(buildingRequirementsBytes, &buildingRequirements)
	if err != nil {
		log.Fatalf("error unmarshalling building requirements: %s\n", err)
	}

	return buildingRequirements
}

func readBuildingLevelPoints() map[string][]int {
	filePath := filepath.Join("resources", "building_points.json")
	buildingLevelPointsBytes, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("error reading building level points from file: %s\n", err)
	}

	var buildingLevelPoints map[string][]int
	err = json.Unmarshal(buildingLevelPointsBytes, &buildingLevelPoints)
	if err != nil {
		log.Fatalf("error unmarshalling building level points: %s\n", err)
	}

	return buildingLevelPoints
}
