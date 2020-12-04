package main

import "fmt"
import "time"

type Game struct {
	playerLoc               int
	playerInventory         map[string]Item
	playerInventoryCapacity int
	itemsByLoc              map[int][]Item
	npcsByLoc               map[int][]NPC
	playerHealth            int
}

func (g *Game) Spawn() {
	g.AddItem(0, Item{
		name:        "roll of burned bandages",
		description: "They appear to be wrapped around several bottles--like a makeshift pack--probably tonics and salves meant for healing. Probably.",
	})
	g.AddItem(0, Item{
		name:        "blackened sword",
		description: "It appears to be burned--and dull. It'll be heavy to carry, but also heavy to receive a blow from--perfect for protecting oneself on hostile roads.",
	})
	g.AddItem(0, Item{
		name:        "indecipherable book",
		description: "The words are pronounceable--maybe. It's some language you aren't familiar with. Seems impracticle to take with you, but maybe it's valuable. Maybe.",
	})
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
	if npcsAtLoc, ok := g.npcsByLoc[g.playerLoc]; ok {
		for i, _ := range npcsAtLoc {
			if npcsAtLoc[i].IsFound == false {
				go func(i int) {
					cc <- npcsAtLoc[i].Find()
				}(i)
			}
			if npcsAtLoc[i].IsAttacking == true {
				go func(i int) {
					cc <- npcsAtLoc[i].Attack(g)
				}(i)
			}
		}
	}
}

// Game.****Item Methods
const CAPACITY = "CAPACITY"
const NOTFOUND = "NOTFOUND"

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

// Game.****NPC Methods
func (g *Game) AddNPC(loc int, npc NPC) {
	g.npcsByLoc[loc] = append(g.npcsByLoc[loc], npc)
}

func (g *Game) GetNPC(loc int, name string) (npc NPC, ok bool) {
	for _, npc := range g.npcsByLoc[loc] {
		if npc.name == name {
			return npc, true
		}
	}
	return NPC{}, false
}

type NPC struct {
	name         string
	description  string
	IsFound      bool
	IsAttacking  bool
	weaponDamage int
}

func (npc *NPC) Find() Command {
	npc.IsFound = true
	go func() {
		time.Sleep(6 * time.Second)
		npc.IsAttacking = true
	}()
	return Command{text: "/find " + npc.name, origin: GAME}
}

func (npc *NPC) Attack(g *Game) Command {
	dmg := npc.weaponDamage
	g.playerHealth = g.playerHealth - dmg
	if g.playerHealth > 0 {
		return Command{text: fmt.Sprint("/attackByNPC ", dmg, " ", npc.name), origin: GAME}
	} else {
		npc.IsAttacking = false
		return Command{text: fmt.Sprint("/deathByNPC ", dmg, " ", npc.name), origin: GAME}
	}
}
