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
func checkEndCondition(s *dgo.Session, guild *data.Guild) (bool) {
  var kingName string = guild.Wink.NameList[guild.Wink.KingID]

  if guild.Wink.TotalParticipants-1 == guild.Wink.ConfirmedCount {
    for k, i :=	range guild.Wink.ConfirmedUsers {
      if i == false {
        guild.Wink.FinalName = guild.Wink.NameList[k]
        break
      }
    }
    log.Printf("User [%s] is the last person to vote!", guild.Wink.FinalName)

    if guild.Wink.FinalName == kingName {
      broadcastResults(s, guild)
    } else {
      broadcastFinalPlayer(s, guild)

      go func() {
        // start game end timer (default: 15 sec? should be enough)
        time.Sleep(5 * time.Second) // TODO: change to 15
        broadcastResults(s, guild)
      }()
    }

    return true
  }

  return false
}

// announce to everyone the last person
func broadcastFinalPlayer(s *dgo.Session, guild *data.Guild) {
  players := guild.Wink.SelectedUsersID

  data := dgo.MessageSend{
    Embeds: []*dgo.MessageEmbed{ 
      {
        Title:        "마지막 사람은 [ " + guild.Wink.FinalName + " ] 닙입니다!",
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
func broadcastResults(s *dgo.Session, guild *data.Guild) {
  var kingName string = guild.Wink.NameList[guild.Wink.KingID]
  var players []string = guild.Wink.SelectedUsersID

  var text_r string = ""
  var text_w string = ""
  var text string = ""

  if guild.Wink.FinalName == kingName {
    // when the king didn't vote
    text += "왕이 투표 안해서 왕이 졌습니다!\n" +
            "왕은 [ " + guild.Wink.FinalName + " ]님이였습니다!" 
  } else {
    // check final person's target
    t := guild.Wink.UserSelectionFinal[guild.Wink.FinalName]
    if len(t) == 0 || t != guild.Wink.KingID {
      // if final player didn't select a target
      text += "[ " + guild.Wink.FinalName + " ]님이 왕을 못맞췄습니다!\n" +
              "왕은 [ " + kingName + " ]님이였습니다!"
    } else {
      // if the final player selected the king
      text += "[ " + guild.Wink.FinalName + " ]님이 왕을 맞췄습니다!\n" +
              "왕은 [ " + kingName + " ]님이였습니다!"
    }
  }
  
  for u, v := range guild.Wink.UserSelectionFinal {
    log.Printf("OKAY STUPID: u[%s] v[%s] kingName[%s]\n", u, v, kingName)
    if v == kingName {
      text_r += " -> " + u + "\n"
    } else {
      if len(guild.Wink.UserSelectionFinal[u]) != 0 {
        text_w += " -> " + guild.Wink.NameList[u] + " [찍은 사람: \""+guild.Wink.UserSelectionFinal[u]+"\"]\n"
      } else {
        text_w += " -> " + guild.Wink.NameList[u] + " [투표 안한 바보 :P]\n"
      }
    }
  }

  text += "\n\n**왕을 정확히 맞춘 사람:**\n"
  text += text_r
  text += "\n**왕을 틀리게 맞춘 사람:**\n"
  text += text_w

  data := dgo.MessageSend{
    Embeds: []*dgo.MessageEmbed{ 
      {
        Title:        "게임 종료!!!",
        Description:  text,
        Color:        0xFFFFFF,
      },
    },
  }

  // ignore array index
  for _, i := range players {
    general.SendComplexDM(s, i, &data)
  }
  
  // TODO: add a announcement to guild chat as well (with reset buttons as well)
}
