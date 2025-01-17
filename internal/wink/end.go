package wink

import (
  // system packages
	"log"
  "time"

  // internal packages
  "github.com/coex1/EchoBot/internal/general"
  "github.com/coex1/EchoBot/internal/data"

  // external packages
  dgo "github.com/bwmarrin/discordgo"
)

// check if the game should continue or end
func end_checkEndCondition(s *dgo.Session, g *data.Guild) (bool) {
  if g.Wink.TotalParticipants-1 == g.Wink.ConfirmedCount {
    for k, i :=	range g.Wink.ConfirmedUsers {
      if i == false {
        g.Wink.FinalPlayerID = k
        break
      }
    }

    log.Printf("User [%s] is the last person to vote!", g.Wink.NameList[g.Wink.FinalPlayerID])
    g.Wink.State = data.LAST_PLAYER

    if g.Wink.FinalPlayerID == g.Wink.KingID {
      end_broadcastResults(s, g)
    } else {
      end_broadcastFinalPlayer(s, g)

      go func() {
        // start game end timer (default: 15 sec? should be enough)
        time.Sleep(5 * time.Second) // TODO: change to 15
        end_broadcastResults(s, g)
      }()
    }

    return true
  }

  return false
}

// announce to everyone the last person
func end_broadcastFinalPlayer(s *dgo.Session, g *data.Guild) {
  players := g.Wink.SelectedUsersID

  data := dgo.MessageSend{
    Embeds: []*dgo.MessageEmbed{ 
      {
        Title:        "마지막 사람은 [ " + g.Wink.NameList[g.Wink.FinalPlayerID] + " ] 닙입니다!",
        Description:  "15초 뒤에 게임이 종료됩니다!\n" +
                      "그 전에  투표해 주세요!",
        Color:        0xFC2803,
      },
    },
  }

  // ignore array index
  for _, i := range players {
    general.SendComplexDM(s, i, &data)
  }
}

// announce to everyone game results
func end_broadcastResults(s *dgo.Session, g *data.Guild) {
  var kingName string = g.Wink.NameList[g.Wink.KingID]
  var finalName string = g.Wink.NameList[g.Wink.FinalPlayerID]
  var players []string = g.Wink.SelectedUsersID

  var text_r string = ""
  var text_w string = ""
  var text string = ""

  if g.Wink.KingID == g.Wink.FinalPlayerID {
    // when the king didn't vote
    text += "왕이 투표 안해서 왕이 졌습니다!\n" +
            "왕은 [ " + finalName + " ]님이였습니다!" 
  } else {
    // check final person's target
    t := g.Wink.UserSelectionFinal[g.Wink.FinalPlayerID]
    if len(t) == 0 || t != g.Wink.KingID {
      // if final player didn't select a target
      text += "[ " + finalName + " ]님이 왕을 못맞췄습니다!\n" +
              "왕은 [ " + kingName + " ]님이였습니다!"
    } else {
      // if the final player selected the king
      text += "[ " + finalName + " ]님이 왕을 맞췄습니다!\n" +
              "왕은 [ " + kingName + " ]님이였습니다!"
    }
  }
  
  for u, v := range g.Wink.UserSelectionFinal {
    if u == g.Wink.KingID {
      continue
    } else if v == g.Wink.KingID {
      text_r += " -> " + g.Wink.NameList[u] + "\n"
    } else {
      if len(g.Wink.UserSelectionFinal[u]) == 0 {
        text_w += " -> " + g.Wink.NameList[u] + " [ 투표 안한 바보 :P ]\n"
      } else {
        text_w += " -> " + g.Wink.NameList[u] + " [ \"" + g.Wink.NameList[g.Wink.UserSelectionFinal[u]] + "\"님을 찍었습니다 ^.^ ]\n"
      }
    }
  }

  text += "\n\n**왕을 맞게 찍은 사람:**\n"
  text += text_r
  text += "\n**왕을 틀리게 찍은 사람:**\n"
  text += text_w

  message := dgo.MessageSend{
    Embeds: []*dgo.MessageEmbed{ 
      {
        Title:        "게임 종료!!!",
        Description:  text,
        Color:        0x9534EB,
      },
    },
  }

  // ignore array index
  for _, i := range players {
    general.SendComplexDM(s, i, &message)
  }

  g.Wink.State = data.ENDED
  
  end_announceToGuild()
}

// TODO: add a announcement to guild chat as well (with reset buttons as well)
func end_announceToGuild() {

}

// for when ending the game forcibly
func End_Game(g *data.Guild) {
  resetGame(g)
}
