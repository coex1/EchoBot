# 설명
Echo(메아리) + Bot(봇) 프로젝트는 Go언어를 공부하자는 의미에서 시작을 했습니다.  
이 프로젝트는, 사람들끼리 모여서 할 수 있을 만한 게임에서 중재 매체의 역할 혹은 사회자의 필요성을 없애기 위한 목표가 있습니다.  
  
Discord Bot이며, [discordgo](https://github.com/bwmarrin/discordgo) 패키지를 사용해서 개발했습니다.  

# 봇 실행 전!
Go가 설치되어 있어야 합니다.  
[윈도우에서 Go 설치하기](https://artist-developer.tistory.com/4) 

# 봇 실행 절차
패키지를 다운 받은 후, example/config.json 파일을 프로젝트 폴더로 옮겨 주세요.  
복사한 파일에 봇의 TOKEN ID, 서버 ID를 입력해 주세요.
- [봇 TOKEN ID 발행](https://www.writebots.com/discord-bot-token/)
- [서버 ID 찾기](https://www.alphr.com/discord-find-server-id/)
이후 다음 명령어를 입력해서 봇을 실행시켜 주세요.

```
build.bat
bot.exe
```

# 기능 설명
현재까지 개발한 게임 지원 기능입니다.  
| 명령어    | 설명 |  옵션  |
| -------- | ------- |
| wink  | 5+ 이상이 있는 서버 채팅 방에서 명령어를 입력하시면 됩니다. |
| mafia | 5+ 이상이 있는 서버 채팅 방에서 명령어를 입력하시면 됩니다. | [Mafia][Police][Doctor] 각 역할의 인원 수 입니다. 나머지 인원은 모두 시민입니다. |


# 윙크 게임 기능 설명
[윙크 게임 설명](https://anyoutplay.tistory.com/57) 
서버 채팅창에서 명령어를 입력하면, 목록 창이 표기 됩니다.  
게임에 참여할 사람들을 목록에서 선택합니다.  
선택 한 후, **게임 시작** 버튼을 클릭 하시면 개인 문자로 역할이 공지됩니다.  
시민들은 왕이 누군지 목록에서 선택한 후 제출 버튼을 클릭하면 제출할 수 있게 됩니다.  
한 사람만 남았을 때, 15초 뒤에 게임이 종료됩니다.  
그때, 투표 안한 사람과 잘못 투표한 사람이 패배하게 됩니다.  

# 마피아 게임 기능 설명
[역할 설명]
- 마피아: 밤에 1명을 지목합니다. 의견이 엇갈려 여러명이 지목됐을 경우 랜덤으로 선택 됩니다.
          지목 당한 시민은 낮이 되면 죽습니다. (마피아 팀)
- 시민: 아무 능력이 없습니다. 토론을 통해 마피아를 찾아내야 합니다.
        밤에는 랜덤으로 주어지는 문장을 입력해야 합니다. (시민 팀)
- 경찰: 밤에 1명을 지목해 그 사람이 마피아인지 아닌지 알 수 있습니다. (시민 팀)
- 의사: 밤에 1명을 지목해 치료할 수 있습니다. 마피아가 지목한 사람과 동일할 경우, 그 사람은 죽지
        않습니다. 자신을 지목할 수 있습니다. (시민 팀)
  
[마피아 게임 설명]
서버 채팅창에 명령어를 입력하여 게임에 참여할 사람들을 목록에서 선택합니다.
**게임 시작** 버튼을 클릭 하시면 개인 문자로 역할이 공지되며 게임이 시작됩니다.
마피아 게임은 낮과 밤으로 진행됩니다. 
낮에는 회의 시간이 10분 주어집니다.
회의를 통해 죽일 사람을 목록에서 선택 후 투표(or 기권) 버튼을 클릭하여 제출 할 수 있습니다.
10분이 지나기 전에 모든 사람이 투표를 완료했을 경우 회의 시간이 종료됩니다.
만약 10분이 지날 떄 까지 투표를 하지 못했을 경우 자동으로 기권한 것으로 처리됩니다.
과반수 이상의 표를 받았을 경우 처형됩니다. 처형된 사람은 게임 종료 시 까지 어떠한 대화도 할 수 없고 정체를 공개하서도 안됩니다. 동점이거나 과반수를 넘지 않았을 경우 아무도 처형되지 않습니다.
회의 시간이 끝나면 밤으로 넘어갑니다.
밤에는 능력을 사용 할 수 있습니다. (반드시 사용해야됨)
모든 사람이 능력 사용을 완료했을 경우 낮으로 넘어갑니다.
투표로 처형되거나 밤에 마피아에 의해 죽을 때 마다 마피아와 시민 팀의 수를 비교하여
마피아의 수가 시민의 수와 같아지면 마피아 팀의 승리! 모든 마피아가 처형되면 시민 팀 승리!

