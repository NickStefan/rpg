package main

import "strconv"

type LocatedEntity struct {
	name        string
	description string
	xLoc        int
	yLoc        int
	command     string
}

func (e LocatedEntity) Loc() string {
	return "x" + strconv.Itoa(e.xLoc) + ", y" + strconv.Itoa(e.yLoc)
}

func (e LocatedEntity) Name() string {
	return e.name
}

type Character struct {
	health    int
	maxHealth int
	inventory []Item
	LocatedEntity
}

type Item struct {
	LocatedEntity
}

type Zone struct {
	minX     int
	maxX     int
	minY     int
	maxY     int
	npcLocs  map[string]map[string]Character
	itemLocs map[string]map[string]Item
}

func (z *Zone) addNpc(c Character) {
	z.npcLocs[c.Loc()][c.Name()] = c
}

func (z *Zone) addItem(i Item) {
	z.itemLocs[i.Loc()][i.Name()] = i
}
