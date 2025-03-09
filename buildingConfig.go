package main

import (
	"bytes"
	"fmt"
	"github.com/antchfx/xmlquery"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strconv"
)

type BuildingInfo struct {
	minLevel     int
	maxLevel     int
	startLevel   int
	restrictions map[string]int
	points       []int
}

var authority = "die-staemme.de"
var path = "interface.php"

func getBuildings(world string) map[string]BuildingInfo {
	filePath := filepath.Join("resources", fmt.Sprintf("building_info_%s.xml", world))
	buildingInfoBytes, err := os.ReadFile(filePath)

	if err != nil {
		buildingInfoBytes = fetchBuildingConfig(world)
		err = os.WriteFile(filePath, buildingInfoBytes, 0644)
		if err != nil {
			log.Fatalf("error writing to file: %s\n", err)
		}
	}

	return readFromBuildingConfig(buildingInfoBytes)
}

func fetchBuildingConfig(world string) []byte {
	queryParams := "func=get_building_info"
	requestURL := fmt.Sprintf("https://%s.%s/%s?%s", world, authority, path, queryParams)
	res, err := http.Get(requestURL)
	if err != nil {
		log.Fatalf("error making http request: %s\n", err)
	}

	defer res.Body.Close()

	buildingInfoBytes, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("error reading from body: %s\n", err)
	}
	if res.StatusCode != 200 {
		log.Fatalf("[Status %d] %s", res.StatusCode, string(buildingInfoBytes))
	}
	return buildingInfoBytes
}

func readFromBuildingConfig(buildingConfig []byte) map[string]BuildingInfo {
	doc, err := xmlquery.Parse(bytes.NewReader(buildingConfig))
	if err != nil {
		log.Fatalf("error parsing building config bytes: %s\n", err)
	}
	config := xmlquery.Find(doc, "//config/*")
	buildings := make(map[string]BuildingInfo)
	for _, building := range config {
		name := building.Data
		maxLevel, _ := strconv.Atoi(building.SelectElement("max_level").InnerText())
		minLevel, _ := strconv.Atoi(building.SelectElement("min_level").InnerText())
		currentLevel := 0
		if slices.Contains([]string{"hide", "main", "storage", "place", "farm"}, name) {
			currentLevel = 1
		}

		buildings[name] = BuildingInfo{
			minLevel:     minLevel,
			maxLevel:     maxLevel,
			startLevel:   currentLevel,
			restrictions: map[string]int{},
			points:       []int{},
		}
	}
	return buildings
}
