package main

import "fmt"
import "math/rand"

type Item struct {
	name        string
	description string
	IsFound     bool
}

type Game struct {
	playerHealth            int
	playerLoc               int
	playerInventoryCapacity int
	playerInventory         map[string]Item
	playerWeaponDamage      int
	playerWeaponDelay       int
	playerWeaponDelayTimer  int
	playerIsAttacking       bool
	playerTarget            *NPC
	itemsByLoc              map[int][]Item
	npcsByLoc               map[int][]NPC
	messagesByLoc           map[int][]Message
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
	g.AddMessage(0, createMessage("Use '/inspect [item]' to look closer."))
	g.AddMessage(0, createMessage("Use '/take [item]' to take one--and at most one. Time is short."))

	g.AddNPC(1, NPC{
		name:         "invading warrior",
		description:  "He looks dangerous and will attack any second.",
		aggroDelay:   3,
		weaponDamage: 10,
		weaponDelay:  3,
		health:       30,
	})
	g.AddMessage(1, createMessage("Use '/consider [name]' to look closer."))
	g.AddMessage(1, createMessage("Use '/attack [name]' to attack."))

	g.AddMessage(2, createMessage("\nYou've reached the gatehouse! Congratulations! Press Q to quit."))
}

func (g *Game) Tick(cc chan<- Command, mc chan<- Message) {
	if itemsAtLoc, ok := g.itemsByLoc[g.playerLoc]; ok {
		for i, _ := range itemsAtLoc {
			if itemsAtLoc[i].IsFound == false {
				itemsAtLoc[i].IsFound = true
				mc <- createMessage("You find a [" + itemsAtLoc[i].name + "].")
			}
		}
	}
	if npcsAtLoc, ok := g.npcsByLoc[g.playerLoc]; ok {
		for i, _ := range npcsAtLoc {
			if npcsAtLoc[i].IsFound == false {
				npcsAtLoc[i].IsFound = true
				mc <- createMessage("Just ahead, you see a [" + npcsAtLoc[i].name + "].")
			} else if npcsAtLoc[i].aggroDelay > 0 {
				npcsAtLoc[i].aggroDelay--
			} else {
				npcsAtLoc[i].IsAttacking = true
			}
			if npcsAtLoc[i].IsAttacking == true {
				if npcsAtLoc[i].weaponDelayTimer == 0 {
					npcsAtLoc[i].Attack(g, mc)
					npcsAtLoc[i].weaponDelayTimer = npcsAtLoc[i].weaponDelay
				} else {
					npcsAtLoc[i].weaponDelayTimer--
				}
			}
		}
	}
	if messagesAtLoc, ok := g.messagesByLoc[g.playerLoc]; ok {
		for i, _ := range messagesAtLoc {
			mc <- messagesAtLoc[i]
		}
		g.messagesByLoc[g.playerLoc] = nil
	}
	if g.playerIsAttacking == true {
		if g.playerWeaponDelayTimer == 0 {
			g.PlayerAttack(g.playerTarget, mc)
			g.playerWeaponDelayTimer = g.playerWeaponDelay
		} else {
			g.playerWeaponDelayTimer--
		}
	}
}

func (g *Game) AddMessage(loc int, message Message) {
	g.messagesByLoc[loc] = append(g.messagesByLoc[loc], message)
}

func (g *Game) InspectItem(name string, mc chan<- Message) {
	item, ok := g.GetItem(g.playerLoc, name)
	if ok == false {
		item, ok = g.GetInventoryItem(name)
	}
	if ok {
		mc <- createMessage("You inspect the " + item.name + ". " + item.description)
	} else {
		mc <- createMessage("There isn't anything by that name here. Maybe it's gone?")
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

func (g *Game) TakeItem(name string, mc chan<- Message) {
	if item, ok := g.GetItem(g.playerLoc, name); ok {
		if ok, _ := g.AddItemToInventory(item); ok {
			g.deleteItem(g.playerLoc, name)
			mc <- createMessage("You take the " + item.name + ".")
		} else {
			mc <- createMessage("You can't carry anything more.'")
		}
	} else {
		mc <- createMessage("There isn't anything by that name here. Maybe it's gone?")
	}
}

func (g *Game) AddItemToInventory(item Item) (ok bool, err string) {
	if g.playerInventoryCapacity < len(g.playerInventory)+1 {
		ok, err = false, "capacity"
	} else {
		g.playerInventory[item.name] = item
		ok = true
	}
	return
}

func (g *Game) ListInventory(mc chan<- Message) {
	mc <- createMessage(fmt.Sprint("Inventory Items(", len(g.playerInventory), "):"))
	for _, item := range g.playerInventory {
		mc <- createMessage("- " + item.name + " --- " + item.description)
	}
}

func (g *Game) AddNPC(loc int, npc NPC) {
	g.npcsByLoc[loc] = append(g.npcsByLoc[loc], npc)
}

func (g *Game) deleteNPC(loc int, name string) {
	var newNPCs []NPC
	for _, npc := range g.npcsByLoc[loc] {
		if npc.name != name {
			newNPCs = append(newNPCs, npc)
		}
	}
	g.npcsByLoc[loc] = newNPCs
}

func (g *Game) ConsiderNPC(name string, mc chan<- Message) {
	if npc, ok := g.GetNPC(g.playerLoc, name); ok {
		mc <- createMessage("You consider the " + name + ". " + npc.description)
	} else {
		mc <- createMessage("There isn't anything by that name here. Maybe it's gone?")
	}
}

func (g *Game) GetNPC(loc int, name string) (npc NPC, ok bool) {
	for _, npc := range g.npcsByLoc[loc] {
		if npc.name == name {
			return npc, true
		}
	}
	return NPC{}, false
}

func (g *Game) ToggleAttack(name string, mc chan<- Message) {
	if npc, ok := g.GetNPC(g.playerLoc, name); ok {
		g.playerTarget = &npc
		g.playerIsAttacking = true
		// TODO use the inventory to change weapon damage!
		// should reflect the item they chose earlier!
		g.playerWeaponDamage = 10
		g.playerWeaponDelay = 5
		g.playerWeaponDelayTimer = 5
	} else {
		mc <- createMessage("There isn't anything by that name here. Maybe it's gone?")
	}
}

func (g *Game) PlayerAttack(npc *NPC, mc chan<- Message) {
	dmg := rand.Intn(g.playerWeaponDamage)
	npc.health = npc.health - dmg
	if npc.health > 0 {
		mc <- createMessage(fmt.Sprint("You hit a ", npc.name, " for ", dmg, " damage."))
	} else {
		npc.IsAttacking = false
		g.playerIsAttacking = false
		g.playerTarget = nil
		g.deleteNPC(g.playerLoc, npc.name)
		g.playerLoc++
		mc <- createMessage(fmt.Sprint("You hit a ", npc.name, " for ", dmg, " damage AND DESTROY IT."))
	}
}

type NPC struct {
	name             string
	description      string
	health           int
	IsFound          bool
	aggroDelay       int
	IsAttacking      bool
	weaponDamage     int
	weaponDelay      int
	weaponDelayTimer int
}

func (npc *NPC) Attack(g *Game, mc chan<- Message) {
	dmg := rand.Intn(npc.weaponDamage)
	g.playerHealth = g.playerHealth - dmg
	if g.playerHealth > 0 {
		mc <- createMessage(fmt.Sprint("A ", npc.name, " hits your for ", dmg, " damage."))
	} else {
		npc.IsAttacking = false
		g.playerIsAttacking = false
		mc <- createMessage(fmt.Sprint("A ", npc.name, " hits your for ", dmg, " damage AND DESTROYS YOU."))
		mc <- createMessage("Game over. Press Q to quit.")
	}
}
