package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"slices"
)

const StartingScore = 26

type VillageExpansionGraphNode struct {
	buildings map[string]int
	score     int
}

var BuildingConfigs map[string]BuildingInfo

func PossibleVillageExpansions(buildingInfo map[string]BuildingInfo, villageIncreases []int) []VillageExpansionGraphNode {
	completeBuildingInfo(buildingInfo)
	BuildingConfigs = buildingInfo
	startVillage := map[string]int{}
	for buildingName, buildingDetails := range buildingInfo {
		startVillage[buildingName] = buildingDetails.startLevel
	}
	rootNode := VillageExpansionGraphNode{
		buildings: startVillage,
		score:     StartingScore,
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

func generateKey(village map[string]int) string {
	var nodeKey []string
	for key := range village {
		nodeKey = append(nodeKey, key)
	}
	slices.Sort(nodeKey)
	var key []byte
	for _, value := range nodeKey {
		key = append(key, byte(village[value]))
	}
	return string(key)
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

	for buildingName, buildingLevel := range n.buildings {
		if !n.isBuildingExpandable(buildingName, buildingLevel, desiredScore-n.score) {
			continue
		}
		expandedVillage := n.expandVillage(buildingName)
		expandedVillageKey := generateKey(expandedVillage)

		if _, ok := graph[expandedVillageKey]; ok {
			continue
		}

		buildingInfo := BuildingConfigs[buildingName]
		child := VillageExpansionGraphNode{
			buildings: expandedVillage,
			// buildingLevel can be used without -1, since it represents the building level from before expansion
			score: n.score + buildingInfo.points[buildingLevel],
		}
		graph[expandedVillageKey] = true
		results = append(results, child)
	}

	return results
}

func (n *VillageExpansionGraphNode) isBuildingExpandable(buildingName string, buildingLevel int, maxIncrease int) bool {
	buildingInfo := BuildingConfigs[buildingName]
	if buildingLevel >= buildingInfo.maxLevel {
		return false
	}

	// Since current level starts at 1, slice at 0 this already looks at the value for the current level + 1
	if buildingInfo.points[buildingLevel] > maxIncrease {
		return false
	}

	// If it is already built, the buildings has to fulfill the requirements.
	if buildingLevel > 0 {
		return true
	}

	// check building requirements
	expandable := true
	for requiredBuilding, requiredLevel := range buildingInfo.restrictions {
		if buildingLevel, ok := n.buildings[requiredBuilding]; ok {
			if buildingLevel < requiredLevel {
				expandable = false
			}
			continue
		}
		log.Fatalf("Required building %s not present in config", requiredBuilding)
	}
	return expandable
}

func (n *VillageExpansionGraphNode) expandVillage(expandedBuilding string) map[string]int {
	expandedVillage := make(map[string]int)
	for building, level := range n.buildings {
		if building == expandedBuilding {
			expandedVillage[expandedBuilding] = level + 1
			continue
		}
		expandedVillage[building] = level
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
