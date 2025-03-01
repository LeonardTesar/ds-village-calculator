package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

type VillageExpansionGraphNode struct {
	village  map[string]BuildingInfo
	score    int
	previous []*VillageExpansionGraphNode
	next     []*VillageExpansionGraphNode
}

func NewVillageExpansionGraph(buildings map[string]BuildingInfo) *VillageExpansionGraphNode {
	completeBuildingInfo(buildings)
	rootNode := &VillageExpansionGraphNode{
		village:  buildings,
		score:    26,
		previous: nil,
		next:     []*VillageExpansionGraphNode{},
	}
	generateGraph(rootNode, 8)
	return rootNode
}

func generateKey(village map[string]BuildingInfo) string {
	var nodeKey []string
	for key, value := range village {
		nodeKey = append(nodeKey, fmt.Sprintf("%s:%d", key, value.currentLevel))
	}
	slices.Sort(nodeKey)
	return strings.Join(nodeKey, "-")
}

var graph = make(map[string]*VillageExpansionGraphNode)

func generateGraph(rootNode *VillageExpansionGraphNode, scoreIncrease int) {
	var queue []*VillageExpansionGraphNode
	desiredScore := rootNode.score + scoreIncrease
	queue = append(queue, rootNode)
	iteration := 0
	for len(queue) != 0 {
		if iteration%1024 == 0 {
			fmt.Println("Iteration:", iteration)
		}
		villageNode := queue[0]
		for building, buildingInfo := range villageNode.village {
			if !villageNode.isBuildingExpandable(buildingInfo, desiredScore-villageNode.score) {
				continue
			}
			buildingInfo.currentLevel += 1
			expandedVillage := villageNode.expandVillage(building, buildingInfo)
			key := generateKey(expandedVillage)

			if village, ok := graph[key]; ok {
				villageNode.next = append(villageNode.next, village)
				village.previous = append(village.previous, villageNode)
				continue
			}

			child := &VillageExpansionGraphNode{
				village:  expandedVillage,
				score:    villageNode.score + buildingInfo.points[expandedVillage[building].currentLevel-1],
				previous: []*VillageExpansionGraphNode{villageNode},
				next:     []*VillageExpansionGraphNode{},
			}
			villageNode.next = append(villageNode.next, child)
			graph[key] = child
			queue = append(queue, child)
		}
		queue = queue[1:]
		iteration++
	}
	fmt.Printf("Took %d iterations\n", iteration)
}

func (n *VillageExpansionGraphNode) isBuildingExpandable(building BuildingInfo, maxIncrease int) bool {
	if building.currentLevel >= building.maxLevel {
		return false
	}
	if building.points[building.currentLevel] > maxIncrease {
		return false
	}
	expandable := true
	for requiredBuilding, requiredLevel := range building.restrictions {
		if buildingInfo, ok := n.village[requiredBuilding]; ok {
			if buildingInfo.currentLevel < requiredLevel {
				expandable = false
			}
			continue
		}
		log.Fatalf("Required building %s not present in config", requiredBuilding)
	}
	return expandable
}

func (n *VillageExpansionGraphNode) expandVillage(expandedBuilding string, expandedInfo BuildingInfo) map[string]BuildingInfo {
	expandedVillage := make(map[string]BuildingInfo)
	for building, info := range n.village {
		if building == expandedBuilding {
			expandedVillage[expandedBuilding] = expandedInfo
			continue
		}
		expandedVillage[building] = info
	}
	return expandedVillage
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

func completeBuildingInfo(buildings map[string]BuildingInfo) {
	buildingLevelPoints := readBuildingLevelPoints()
	buildingRequirements := readBuildingRequirements()

	for buildingName, buildingInfo := range buildings {
		buildingInfo.points = buildingLevelPoints[buildingName]
		buildingInfo.restrictions = buildingRequirements[buildingName]
		buildings[buildingName] = buildingInfo
	}
}
