package components

import "github.com/bwmarrin/discordgo"

const PURPLE = 0xA21DB9
const RED = 0xE02700
const GREEN = 0x0FE000

var CancelButton = discordgo.Button{
	Label:    "Cancel",
	Style:    discordgo.DangerButton,
	CustomID: "button_cancel",
}

var SkipButton = discordgo.Button{
	Label:    "Skip",
	Style:    discordgo.PrimaryButton,
	CustomID: "button_skip",
}

var PauseButton = discordgo.Button{
	Label:    "Pause",
	Style:    discordgo.PremiumButton,
	CustomID: "button_pause",
}

var ResumeButton = discordgo.Button{
	Label:    "Resume",
	Style:    discordgo.PremiumButton,
	CustomID: "button_pause",
}
