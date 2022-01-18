

// Package election implements dpos's election method.
package election

import (
	"math/rand"
	"sort"

	"github.com/gcchains/chain/commons/log"
	"github.com/gcchains/chain/consensus/dpos/rpt"
	"github.com/ethereum/go-ethereum/common"
)


func Elect(rpts rpt.RptList, seed int64, totalSeats int, lowRptCount int, lowRptSeats int) []common.Address {
	if lowRptCount > rpts.Len() || lowRptSeats > totalSeats || totalSeats > rpts.Len() {
		return []common.Address{}
	}

	sort.Sort(rpts)

	highRptCount := len(rpts) - lowRptCount
	highRptSeats := totalSeats - lowRptSeats

	log.Debug("elect parameters", "seed", seed, "total seats", totalSeats, "lowRptCount", lowRptCount, "lowRptSeats", lowRptSeats, "highRptCount", highRptCount, "highRptSeats", highRptSeats)

	lowRpts := rpts[:lowRptCount]
	highRpts := rpts[lowRptCount:]

	randSource := rand.NewSource(seed)
	myRand := rand.New(randSource)

	log.Debug("random select by rpt parameters", "lowRptSeats", lowRptSeats, "highRptSeats", totalSeats-lowRptSeats)

	var lowElected, highElected []common.Address

	if lowRptCount > lowRptSeats {
		lowElected = randomSelectByRpt(lowRpts, myRand, lowRptSeats)
	} else {
		lowElected = lowRpts.Addrs()
	}

	if highRptCount > highRptSeats {
		highElected = randomSelectByRpt(highRpts, myRand, highRptSeats)
	} else {
		highElected = highRpts.Addrs()
	}

	return append(lowElected, highElected...)
}


func randomSelectByRpt(rpts rpt.RptList, myRand *rand.Rand, seats int) (result []common.Address) {
	
	sort.Sort(rpts)

	log.Debug("seats", "seats", seats)

	sums, sum := sumOfFirstN(rpts)
	selected := make(map[int]struct{})

	log.Debug("sums and sum", "sums", sums, "sum", sum)

	for seats > 0 {
		log.Debug("seats in for loop", "seats", seats)

		randI := myRand.Int63n(sum)
		resultIdx := findHit(randI, sums)

		log.Debug("randI", "randI", randI, "result idx", resultIdx, "sum", sum, "sums", sums, "result", result)

		// if already selected, continue
		if _, already := selected[resultIdx]; already {
			continue
		}

		// not selected yet, append it!
		selected[resultIdx] = struct{}{}
		result = append(result, rpts[resultIdx].Address)

		seats--

	}
	return result
}

func findHit(hit int64, hitSums []int64) int {
	for idx, x := range hitSums {
		log.Debug("find hit", "hit", hit, "idx", idx, "x", x, "hit sums", hitSums)
		if hit <= x {
			return idx
		}
	}
	return len(hitSums) - 1
}

func sumOfFirstN(rpts rpt.RptList) (sums []int64, sum int64) {
	sum = 0
	for _, rpt := range rpts {
		sum += rpt.Rpt
		sums = append(sums, sum)
	}
	return
}
