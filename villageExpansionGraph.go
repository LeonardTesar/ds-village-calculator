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

const StartingScore = 26

type VillageExpansionGraphNode struct {
	buildings map[string]BuildingInfo
	score     int
	next      []string
}

func PossibleVillageExpansions(buildings map[string]BuildingInfo, villageIncreases []int) []VillageExpansionGraphNode {
	completeBuildingInfo(buildings)
	rootNode := VillageExpansionGraphNode{
		buildings: buildings,
		score:     StartingScore,
		next:      []string{},
	}
	villageNodes := []VillageExpansionGraphNode{rootNode}
	for _, increase := range villageIncreases {
		var tempNodes []VillageExpansionGraphNode
		for _, node := range villageNodes {
			tempNodes = append(tempNodes, possibleExpandedVillages(node, increase)...)
		}
		villageNodes = tempNodes
	}
	return villageNodes
}

func generateKey(village map[string]BuildingInfo) string {
	var nodeKey []string
	for key, value := range village {
		nodeKey = append(nodeKey, fmt.Sprintf("%s:%d", key, value.currentLevel))
	}
	slices.Sort(nodeKey)
	return strings.Join(nodeKey, "-")
}

var graph = make(map[string]bool)

func possibleExpandedVillages(rootNode VillageExpansionGraphNode, scoreIncrease int) []VillageExpansionGraphNode {
	var queue []VillageExpansionGraphNode
	var results []VillageExpansionGraphNode
	desiredScore := rootNode.score + scoreIncrease
	queue = append(queue, rootNode)
	iteration := 0
	for len(queue) != 0 {
		if iteration%1024 == 0 {
			fmt.Println("Iteration:", iteration)
		}

		// Pop queue element
		villageNode := queue[0]
		queue = queue[1:]

		if villageNode.score == desiredScore {
			results = append(results, villageNode)
		} else {
			queue = append(queue, villageNode.generatePossibleExpansionNodes(desiredScore)...)
		}

		villageKey := generateKey(villageNode.buildings)
		delete(graph, villageKey)

		iteration++
	}
	fmt.Printf("Took %d iterations\n", iteration)
	return results
}

func (n *VillageExpansionGraphNode) generatePossibleExpansionNodes(desiredScore int) []VillageExpansionGraphNode {
	var results []VillageExpansionGraphNode

	for building, buildingInfo := range n.buildings {
		if !n.isBuildingExpandable(buildingInfo, desiredScore-n.score) {
			continue
		}
		buildingInfo.currentLevel += 1
		expandedVillage := n.expandVillage(building, buildingInfo)
		expandedVillageKey := generateKey(expandedVillage)

		if _, ok := graph[expandedVillageKey]; ok {
			continue
		}

		child := VillageExpansionGraphNode{
			buildings: expandedVillage,
			score:     n.score + buildingInfo.points[expandedVillage[building].currentLevel-1],
		}
		graph[expandedVillageKey] = true
		results = append(results, child)
	}

	return results
}

func (n *VillageExpansionGraphNode) isBuildingExpandable(building BuildingInfo, maxIncrease int) bool {
	if building.currentLevel >= building.maxLevel {
		return false
	}

	// Since current level starts at 1, slice at 0 this already looks at the value for the current level + 1
	if building.points[building.currentLevel] > maxIncrease {
		return false
	}

	// If it is already built, the buildings has to fulfill the requirements.
	if building.currentLevel > 0 {
		return true
	}

	// check building requirements
	expandable := true
	for requiredBuilding, requiredLevel := range building.restrictions {
		if buildingInfo, ok := n.buildings[requiredBuilding]; ok {
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
	for building, info := range n.buildings {
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
