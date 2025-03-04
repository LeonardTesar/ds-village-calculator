package TwStatsScraper

import "fmt"

func ScrapeVillageHistory(world string, village int) []int {
	fmt.Printf("Scrap history for village %d in world %s", village, world)
	return []int{62, 64, 7}
}
