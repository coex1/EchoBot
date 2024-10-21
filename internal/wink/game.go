package wink

import (
  // system packages
	"log"

  // internal packages
  "github.com/coex1/EchoBot/internal/general"

  // external packages
  dgo "github.com/bwmarrin/discordgo"
)

// 사용자 목록에서 왕 선택
func selectKing(players []string) (kingID string){
	kingID = players[general.Random(0, len(players)-1)]
  log.Printf("Selected king! [%s]", kingID)
  return
}

// 역할 공지!
func sendRoleNotice(s *dgo.Session, players []string, kingID string) {
  var message string

  // ignore index
  for _, i := range players {
    if i == kingID {
      message = "당신은 왕 입니다!\n" +
                "주변 사람들에게 조심스럼게 윙크를 보내세요!\n" +
                "(맘에 안드는 사람에게는 윙크 보내실 필요 없습니다)"
    } else {
      message = "당신은 왕이 아닙니다!"
    }

    general.SendDM(s, i, message)
  }
}
