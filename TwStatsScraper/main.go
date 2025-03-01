package TwStatsScraper

import "fmt"

func ScrapeVillageHistory(world string, village int) []VillageHistoryEntry {
	fmt.Printf("Scrap history for village %d in world %s", village, world)
	return []VillageHistoryEntry{}
}
