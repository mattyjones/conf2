package yang

import (
	"log"
	"schema"
	"testing"
)

func LoadSampleModule(t *testing.T) *schema.Module {
	//	f := &schema.FileStreamSource{Root:"../testdata"}
	m, err := LoadModuleFromByteArray([]byte(TestDataRomancingTheStone), nil)
	if err != nil {
		t.Fatal(err.Error())
	}
	return m
}

func TestLoadLoader(t *testing.T) {
	//yyDebug = 4
	LoadSampleModule(t)
}

func TestFindByPath(t *testing.T) {
	m := LoadSampleModule(t)
	groupings := m.GetGroupings().(*schema.MetaContainer)
	found := schema.FindByPath(groupings, "team")
	log.Println("found", found)
	AssertIterator(t, m.DataDefs())
	if teams := schema.FindByPathWithoutResolvingProxies(m, "game/teams"); teams != nil {
		if def := schema.FindByPathWithoutResolvingProxies(teams.(schema.MetaList), "team"); def != nil {
			if team, isContainer := def.(*schema.Container); isContainer {
				AssertFindGrouping(t, team)
			} else {
				t.Error("Found team but it's not a container type")
			}
		} else {
			t.Error("FindByPathWithoutResolvingProxies Could not find ../teams/team")
		}
		AssertProxies(t, teams.(schema.MetaList))
	} else {
		t.Error("Could not find game/teams")
	}
}

func AssertIterator(t *testing.T, defs schema.MetaList) {
	i := schema.NewMetaListIterator(defs, false)
	game := i.NextMeta()
	if game == nil {
		t.Error("first and only child:game not found in module defs")
	} else if game.GetIdent() != "game" {
		t.Error("expected 'game' child but got ", game.GetIdent())
	} else {
		t.Log("Iterator passed")
	}
}

func AssertFindGrouping(t *testing.T, team *schema.Container) {
	if uses := team.GetFirstMeta(); uses != nil {
		if grouping := uses.(*schema.Uses).FindGrouping("team"); grouping != nil {
			t.Log("Found team grouping")
		} else {
			t.Error("Could not find 'team' grouping in 'team' container")
		}
	} else {
		t.Error("Could not find uses child in team")
	}
}

func AssertProxies(t *testing.T, teams schema.MetaList) {
	if def := schema.FindByPath(teams, "team"); def != nil {
		i := schema.NewMetaListIterator(def.(schema.MetaList), true)
		t.Log("first team child", i.NextMeta().GetIdent())
		i = schema.NewMetaListIterator(def.(schema.MetaList), true)
		c := schema.FindByIdent(i, "color")
		t.Log("color", c)

		if members := schema.FindByPath(def.(schema.MetaList), "members"); members != nil {
			t.Log("Found members from grouping")
		} else {
			t.Error("team grouping didn't resolve")
		}
	} else {
		t.Error("FindByPath could not find ../teams/team")
	}
}
