module music-bot

go 1.25

require (
	github.com/bwmarrin/discordgo v0.26.2
	github.com/joho/godotenv v1.5.1
	gopkg.in/hraban/opus.v2 v2.0.0-20230925203106-0188a62cb302
)

require (
	github.com/gorilla/websocket v1.5.3 // indirect
	golang.org/x/crypto v0.43.0 // indirect
	golang.org/x/sys v0.37.0 // indirect
)

replace github.com/bwmarrin/discordgo => github.com/ozraru/discordgo v0.26.2-0.20251101184423-6792228f3271
