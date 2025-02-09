//go:build !solution

package hotelbusiness

import "sort"

type Guest struct {
	CheckInDate  int
	CheckOutDate int
}

type Load struct {
	StartDate  int
	GuestCount int
}

func ComputeLoad(guests []Guest) []Load {
	starts := make([]int, len(guests))
	ends := make([]int, len(guests))
	loads := []Load{}

	for index, guest := range guests {
		starts[index] = guest.CheckInDate
		ends[index] = guest.CheckOutDate
	}

	sort.Ints(starts)
	sort.Ints(ends)
	currentLoad := 0

	startPtr, endPtr := 0, 0
	for startPtr < len(starts) || endPtr < len(ends) {
		if startPtr != len(starts) && ends[endPtr] > starts[startPtr] {
			currentLoad++
			for startPtr < len(starts)-1 && starts[startPtr] == starts[startPtr+1] {
				startPtr++
				currentLoad++
			}
			loads = append(loads, Load{StartDate: starts[startPtr], GuestCount: currentLoad})

			startPtr++
			continue
		}

		if (startPtr != len(starts) && ends[endPtr] < starts[startPtr]) || startPtr == len(starts) {
			currentLoad--
			for endPtr < len(ends)-1 && ends[endPtr] == ends[endPtr+1] {
				endPtr++
				currentLoad--
			}

			loads = append(loads, Load{StartDate: ends[endPtr], GuestCount: currentLoad})
			endPtr++
			continue
		}

		startPtr++
		endPtr++

	}

	return loads
}
