package yang
import (
	"testing"
	"log"
)

func LoadSampleModule(t *testing.T) (*Module) {
	f := &FileResolver{}
	m, err:= LoadModule(f, "test_data/romancing-the-stone.yang")
	if err != nil {
		t.Error(err.Error())
	}
	return m
}

func TestLoadLoader(t *testing.T) {
	//yyDebug = 4
	LoadSampleModule(t)
}

func TestFindByPath(t *testing.T) {
	m := LoadSampleModule(t)
	groupings := m.GetGroupings().(* DefContainer)
	found := FindByPath(groupings, "team")
	log.Println("found", found)
	AssertIterator(t, m.DataDefs())
	if teams := FindByPathWithoutResolvingProxies(m, "game/teams"); teams != nil {
		if def := FindByPathWithoutResolvingProxies(teams.(DefList), "team"); def != nil {
			if team, isContainer := def.(*Container); isContainer {
				AssertFindGrouping(t, team)
			} else {
				t.Error("Found team but it's not a container type")
			}
		} else {
			t.Error("FindByPathWithoutResolvingProxies Could not find ../teams/team")
		}
		AssertProxies(t, teams.(DefList))
	} else {
		t.Error("Could not find game/teams")
	}
}

func AssertIterator(t *testing.T, defs DefList) {
	i := NewDefListIterator(defs, false)
	game := i.NextDef()
	if game == nil {
		t.Error("first and only child:game not found in module defs")
	} else if game.GetIdent() != "game" {
		t.Error("expected 'game' child but got ", game.GetIdent())
	} else {
		t.Log("Iterator passed")
	}
}

func AssertFindGrouping(t *testing.T, team *Container) {
	if uses := team.GetFirstDef(); uses != nil {
		if grouping := uses.(*Uses).FindGrouping("team"); grouping != nil {
			t.Log("Found team grouping")
		} else {
			t.Error("Could not find 'team' grouping in 'team' container")
		}
	} else {
		t.Error("Could not find uses child in team")
	}
}

func AssertProxies(t *testing.T, teams DefList) {
	if def := FindByPath(teams, "team"); def != nil {
		i := NewDefListIterator(def.(DefList), true)
		t.Log("first team child", i.NextDef().GetIdent())
		i = NewDefListIterator(def.(DefList), true)
		c := FindByIdent(i, "color")
		t.Log("color", c)

		if color := FindByPath(def.(DefList), "color"); color != nil {
			t.Log("Found color from grouping")
		} else {
			t.Error("team grouping didn't resolve")
		}
	} else {
		t.Error("FindByPath could not find ../teams/team")
	}
}