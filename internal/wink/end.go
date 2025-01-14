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
  final_person_global_name := ""

  if guild.Wink.TotalParticipants-1 == guild.Wink.ConfirmedCount {
    log.Println("map? : ", guild.Wink.ConfirmedUsers)
    for k, i :=	range guild.Wink.ConfirmedUsers {
      if i == false {
        final_person_global_name = k
      }
    }
    
    log.Printf("Ending game!!! final_person_global_name= [%s]", final_person_global_name)
    send_noti_of_final_person(s, guild, final_person_global_name)

    return true
  }

  return false
}

func send_noti_of_final_person(s *dgo.Session, guild *data.Guild, f string) {
  // TODO: f to global variable (yes, failures get global haha)
  if f == guild.Wink.KingName {
    announce_results(s, guild, f)
  } else {
    announce_last_person(s, guild, f)
    go func() {
      // start game end timer (default: 15 sec?)
      time.Sleep(5 * time.Second) // TODO: change to 15
      announce_results(s, guild, f)
    }()
  }
}

// announce to everyone the last person
func announce_last_person (s *dgo.Session, guild *data.Guild, f string) {
  players := guild.Wink.SelectedUsersID

  embed := dgo.MessageEmbed{
    Title:        "마지막 사람이.....",
    Description:  "\""+f+"\" 입니다!!!\n" +
                  "15초 뒤에 게임이 종료되니 얼른 투표해!",
    Color:        0xFC2803,
  }
  data := dgo.MessageSend{
    Embeds: []*dgo.MessageEmbed{ 
      &embed,
    },
  }

  // ignore array index
  for _, i := range players {
    general.SendComplexDM(s, i, &data)
  }
}

// announce to everyone game results
func announce_results(s *dgo.Session, guild *data.Guild, f string) {
  players := guild.Wink.SelectedUsersID
  right := make([]string, guild.Wink.TotalParticipants)
  wrong := make([]string, guild.Wink.TotalParticipants)
  text := ""

  if f == guild.Wink.KingName {
    // king lose
    text += "왕이 투표 안해서 왕이 졌습니다!\n왕은 ["+f+"]님이였습니다!" 
  } else {
    // check final person's target
    target := guild.Wink.UserSelectionFinal[f]
    if len(target) == 0 || target != guild.Wink.KingName {
      // f loses
      text += "["+f+"]님이 왕을 못맞췄습니다!\n왕은 ["+guild.Wink.KingName+"]님이였습니다!"
      wrong = append(wrong, f)
    } else {
      // king lose
      text += "["+f+"]님이 왕을 맞췄습니다!\n왕은 ["+guild.Wink.KingName+"]님이였습니다!"
    }
  }
  
  for u, v := range guild.Wink.UserSelectionFinal {
    if v == guild.Wink.KingName {
      right = append(right, u)
    } else {
      wrong = append(wrong, u)
    }
  }

  text += "\n\n**왕을 정확히 맞춘 사람:**\n"
  for _, v := range right {
    if len(v) != 0 {
      text += " -> " + v + "\n"
    }
  }

  text += "\n**왕을 틀리게 맞춘 사람:**\n"
  for _, v := range wrong {
    if len(v) != 0 {
      if len(guild.Wink.UserSelectionFinal[v]) != 0 {
        text += " -> " + v + " [찍은 사람: \""+guild.Wink.UserSelectionFinal[v]+"\"]\n"
      } else {
        text += " -> " + v + " [투표안한바보]\n"
      }
    }
  }

  embed := dgo.MessageEmbed{
    Title:        "게임 종료!!!",
    Description:  text,
    Color:        0xFFFFFF,
  }

  data := dgo.MessageSend{
    Embeds: []*dgo.MessageEmbed{ 
      &embed,
    },
  }

  // ignore array index
  for _, i := range players {
    general.SendComplexDM(s, i, &data)
  }
  
  // TODO: add a announcement to guild chat as well (with reset buttons as well)
    
}

