package demoinfocs

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/markus-wa/demoinfocs-golang/common"
)

func TestNewGameState(t *testing.T) {
	gs := newGameState()

	assert.NotNil(t, gs.playersByEntityID)
	assert.NotNil(t, gs.playersByUserID)
	assert.NotNil(t, gs.grenadeProjectiles)
	assert.NotNil(t, gs.infernos)
	assert.NotNil(t, gs.entities)
	assert.Equal(t, common.TeamTerrorists, gs.tState.Team())
	assert.Equal(t, common.TeamCounterTerrorists, gs.ctState.Team())
}

func TestGameState_Participants(t *testing.T) {
	gs := newGameState()
	ptcp := gs.Participants()
	byEntity := ptcp.ByEntityID()
	byUserID := ptcp.ByUserID()

	// Should update ptcp as well since it uses the same map
	gs.playersByEntityID[0] = common.NewPlayer()
	gs.playersByUserID[0] = common.NewPlayer()

	assert.Equal(t, gs.playersByEntityID, ptcp.ByEntityID())
	assert.Equal(t, gs.playersByUserID, ptcp.ByUserID())

	// But should not update byEntity or byUserID since they're copies
	assert.NotEqual(t, byEntity, ptcp.ByEntityID())
	assert.NotEqual(t, byUserID, ptcp.ByUserID())
}

func TestParticipants_All(t *testing.T) {
	gs := newGameState()
	pl := common.NewPlayer()
	gs.playersByUserID[0] = pl

	allPlayers := gs.Participants().All()

	assert.Len(t, allPlayers, 1)
	assert.Equal(t, pl, allPlayers[0])
}

func TestParticipants_Playing(t *testing.T) {
	gs := newGameState()

	terrorist := common.NewPlayer()
	terrorist.Team = common.TeamTerrorists
	gs.playersByUserID[0] = terrorist
	ct := common.NewPlayer()
	ct.Team = common.TeamCounterTerrorists
	gs.playersByUserID[1] = ct
	unassigned := common.NewPlayer()
	unassigned.Team = common.TeamUnassigned
	gs.playersByUserID[2] = unassigned
	spectator := common.NewPlayer()
	spectator.Team = common.TeamSpectators
	gs.playersByUserID[3] = spectator
	def := common.NewPlayer()
	gs.playersByUserID[4] = def

	playing := gs.Participants().Playing()

	assert.Len(t, playing, 2)
	assert.Equal(t, terrorist, playing[0])
	assert.Equal(t, ct, playing[1])
}

func TestParticipants_TeamMembers(t *testing.T) {
	gs := newGameState()

	terrorist := common.NewPlayer()
	terrorist.Team = common.TeamTerrorists
	gs.playersByUserID[0] = terrorist
	ct := common.NewPlayer()
	ct.Team = common.TeamCounterTerrorists
	gs.playersByUserID[1] = ct
	unassigned := common.NewPlayer()
	unassigned.Team = common.TeamUnassigned
	gs.playersByUserID[2] = unassigned
	spectator := common.NewPlayer()
	spectator.Team = common.TeamSpectators
	gs.playersByUserID[3] = spectator
	def := common.NewPlayer()
	gs.playersByUserID[4] = def

	playing := gs.Participants().TeamMembers(common.TeamCounterTerrorists)

	assert.Len(t, playing, 1)
	assert.Equal(t, ct, playing[0])
}

func TestParticipants_FindByHandle(t *testing.T) {
	gs := newGameState()

	pl := common.NewPlayer()
	pl.Team = common.TeamTerrorists
	gs.playersByEntityID[3000&entityHandleIndexMask] = pl

	found := gs.Participants().FindByHandle(3000)

	assert.Equal(t, pl, found)
}

func TestParticipants_FindByHandle_InvalidEntityHandle(t *testing.T) {
	gs := newGameState()

	pl := common.NewPlayer()
	pl.Team = common.TeamTerrorists
	gs.playersByEntityID[invalidEntityHandle&entityHandleIndexMask] = pl

	found := gs.Participants().FindByHandle(invalidEntityHandle)

	assert.Nil(t, found)
}
