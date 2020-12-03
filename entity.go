package main

type Game struct {
	playerLoc               int
	playerInventory         map[string]Item
	playerInventoryCapacity int
	itemsByLoc              map[int][]Item
}

func (g *Game) CheckTriggers(cc chan<- Command) {
	if itemsAtLoc, ok := g.itemsByLoc[g.playerLoc]; ok {
		for i, _ := range itemsAtLoc {
			if itemsAtLoc[i].IsFound == false {
				go func(i int) {
					cc <- itemsAtLoc[i].Find()
				}(i)
			}
		}
	}
}

func (g *Game) AddItem(loc int, item Item) {
	g.itemsByLoc[loc] = append(g.itemsByLoc[loc], item)
}

func (g *Game) GetItem(loc int, name string) (item Item, ok bool) {
	for _, item := range g.itemsByLoc[loc] {
		if item.name == name {
			return item, true
		}
	}
	return Item{}, false
}

func (g *Game) GetInventoryItem(name string) (item Item, ok bool) {
	item, ok = g.playerInventory[name]
	return item, ok
}

func (g *Game) deleteItem(loc int, name string) {
	var newItems []Item
	for _, item := range g.itemsByLoc[loc] {
		if item.name != name {
			newItems = append(newItems, item)
		}
	}
	g.itemsByLoc[loc] = newItems
}

const CAPACITY = "CAPACITY"
const NOTFOUND = "NOTFOUND"

func (g *Game) TakeItem(name string) (ok bool, err string) {
	if item, _ok := g.GetItem(g.playerLoc, name); _ok {
		if g.playerInventoryCapacity < len(g.playerInventory)+1 {
			ok, err = false, CAPACITY
		} else {
			g.playerInventory[name] = item
			ok = true
			g.deleteItem(g.playerLoc, name)
		}
	} else {
		ok, err = false, NOTFOUND
	}
	return
}

type Item struct {
	name        string
	description string
	IsFound     bool
}

func (item *Item) Find() Command {
	item.IsFound = true
	return Command{text: "/find " + item.name, origin: GAME}
}
