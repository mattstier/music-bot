package components

import "github.com/bwmarrin/discordgo"

const PURPLE = 0xA21DB9
const RED = 0xE02700
const GREEN = 0x0FE000
const BLUE = 0x003FBE

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
	Style:    discordgo.SecondaryButton,
	CustomID: "button_pause",
}

var ResumeButton = discordgo.Button{
	Label:    "Resume",
	Style:    discordgo.SuccessButton,
	CustomID: "button_pause",
}

var ForwardJumpButton = discordgo.Button{
	Label:    "+10s",
	Style:    discordgo.SecondaryButton,
	CustomID: "button_jump_10s_forward",
}

var BackwardJumpButton = discordgo.Button{
	Label:    "-10s",
	Style:    discordgo.SecondaryButton,
	CustomID: "button_jump_10s_backward",
}

var ListQueueButton = discordgo.Button{
	Label:    "Show queue",
	Style:    discordgo.SecondaryButton,
	CustomID: "button_list_queue",
}

var ExpandListButton = discordgo.Button{
	Label:    "Show all",
	Style:    discordgo.SecondaryButton,
	CustomID: "expand_list",
}

var CollapseListButton = discordgo.Button{
	Label:    "Collapse list",
	Style:    discordgo.SecondaryButton,
	CustomID: "collapse_list",
}
