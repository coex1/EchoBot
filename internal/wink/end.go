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

func checkEndCondition(s *dgo.Session, guild *data.Guild) (bool) {
  final_person_global_name := ""

  if guild.Wink.TotalParticipants-1 == guild.Wink.ConfirmedCount {
    log.Println("map? : ", guild.Wink.ConfirmedUsers)
    for k, i :=	range guild.Wink.ConfirmedUsers {
      log.Printf("---------> iter test = k=%s i=%d", k, i)
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
      time.Sleep(15)
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

  embed_king_win := dgo.MessageEmbed{
    Title:        "게임이 종료되었습니다!!!",
    Description:  "왕은 \""+guild.Wink.KingName+"\"이였습니다!\n" +
                  "마지막 시민은 \""+f+"\"이였습니다\n" +
                  "\""+f+"\"님이 지셨습니다!",
    Color:        0xFC2803,
  }
  embed_last_win := dgo.MessageEmbed{
    Title:        "게임이 종료되었습니다!!!",
    Description:  "왕은 \""+guild.Wink.KingName+"\"이였습니다!\n" +
                  "마지막 시민은 \""+f+"\"이였습니다\n" +
                  "\""+guild.Wink.KingName+"\"님이 지셨습니다!",
    Color:        0xFC2803,
  }

  data := dgo.MessageSend{
    Embeds: []*dgo.MessageEmbed{ 
      &embed_king_win,
    },
  }
  data_2 := dgo.MessageSend{
    Embeds: []*dgo.MessageEmbed{ 
      &embed_last_win,
    },
  }

  // ignore array index
  for _, i := range players {
    general.SendComplexDM(s, i, &data)
    general.SendComplexDM(s, i, &data_2)
  }
  
  // TODO: add a announcement to guild chat as well (with reset buttons as well)
    
}

