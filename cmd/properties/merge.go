package main

// MergeBundles merges two bundles together, takes care of parent-child relations,
// returns list of bundle ids.
func MergeBundles(curb, newb []Bundle) (merged []int) {
	var (
		childs = map[int]int{}
		set    = map[int]struct{}{}
	)

	// step 1: iterate new ones, collecting id's
	// and parent-child connections (if any)
	for i := 0; i < len(newb); i++ {
		var (
			b        = &newb[i]
			bid, pid = b.ID, b.ParentID
		)

		set[bid] = struct{}{}

		if pid != 0 {
			childs[pid] = bid
		}
	}

	// step 2: iterate current state
	for i := 0; i < len(curb); i++ {
		var (
			b        = &curb[i]
			bid, pid = b.ID, b.ParentID
			ok       bool
		)

		if _, ok = childs[bid]; ok {
			// if our child already in set - downgrade mode, skip self
			continue
		}

		if pid != 0 {
			if _, ok = set[pid]; ok {
				// if our parent already in set - upgrade mode, skip self
				continue
			}
		}

		set[bid] = struct{}{}
	}

	for k, _ := range set {
		merged = append(merged, k)
	}

	return merged
}

// DropBundles removes elements of `cutb` from `curb`, takes care of parent-child relations,
// returns list of bundle ids.
func DropBundles(curb, cutb []Bundle) (merged []int) {
	var set = map[int]struct{}{}

	// step 1: iterate cut set, collecting id's
	for i := 0; i < len(cutb); i++ {
		var (
			b   = &cutb[i]
			bid = b.ID
		)

		set[bid] = struct{}{}
	}

	// step 2: iterate current state
	for i := 0; i < len(curb); i++ {
		var (
			b   = &curb[i]
			bid = b.ID
		)

		// replace bundle id with parent upon deletion
		if _, ok := set[bid]; ok {
			bid = b.ParentID
		}

		merged = append(merged, bid)
	}

	return merged
}
