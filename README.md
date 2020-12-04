### RPG in Go/Golang + TermDash UI

This is a text adventure game modeled after the chat interface used in RPGs/MMORPGs (e.g. EverQuest). User actions are entirely based on "chat commands" e.g. "/inspect sword" or "/attack thing". The game tracks spatial points and the associated triggers at each point. These can be items, NPCs, or messages that the user interacts with as they advance through the spatial points in the game. The NPC and player character "/attack thing" actions make use of the game's ticker (e.g. it has a sense of time).

For now, this is just a fun project to practice some Go concepts. Eventually I could add a Map to the UI, add more interesting free-movement, etc. For now it's WIP, not finished, but runnable and playable.
