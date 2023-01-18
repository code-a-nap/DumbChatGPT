package main

import (
	"os"
	"regexp"

	"github.com/google/uuid"
)

var botName = "DumbChatGPT"

var secretKey = []byte("You_Cant_Guess_This")

var answers = []string{
	"I'm not sure what you're asking, but I'm pretty sure it's above my pay grade.",
	"I think I need a translator for this question.",
	"I'm as confused as a chameleon in a bag of skittles.",
	"I must have missed the class on understanding this type of question.",
	"I think my brain just shut off in self defense.",
	"I'm trying to understand, but it's like trying to fit a square peg in a round hole.",
	"I'm not sure what you're asking, but I'm pretty sure the answer is 42.",
	"I think the question is trying to trap me like a fly in a spiderweb.",
	"I'm not sure if you're asking a question or speaking in tongues.",
	"I'm trying to understand, but it's like trying to solve a Rubik's cube blindfolded.",
	"I'm not sure what you're asking, but I'm pretty sure it's not in my job description.",
	"I'm trying to understand, but it's like trying to find a needle in a haystack.",
	"I'm not sure what you're asking, but I'm pretty sure it's impossible.",
	"I'm trying to understand, but it's like trying to tie a shoelace with oven mitts.",
	"I'm not sure what you're asking, but I'm pretty sure it's a trap.",
	"I'm trying to understand, but it's like trying to navigate a maze with no walls.",
	"I'm not sure what you're asking, but I'm pretty sure it's a trick question.",
	"I'm trying to understand, but it's like trying to unscramble a jumbled word.",
	"I'm not sure what you're asking, but I'm pretty sure it's a riddle.",
	"I'm trying to understand, but it's like trying to make sense of a dream.",
	"I'm not sure what you're asking, but I'm pretty sure it's a conspiracy.",
	"I'm trying to understand, but it's like trying to decode a secret message.",
	"I'm not sure what you're asking, but I'm pretty sure it's a puzzle.",
	"I'm trying to understand, but it's like trying to climb a mountain with no ropes.",
	"I'm not sure what you're asking, but I'm pretty sure it's a labyrinth.",
	"I'm trying to understand, but it's like trying to swim in a sea of molasses.",
	"I'm not sure what you're asking, but I'm pretty sure it's a conundrum.",
	"I'm trying to understand, but it's like trying to fly a kite in a tornado.",
	"I'm not sure what you're asking, but I'm pretty sure it's a mystery.",
	"I'm trying to understand, but it's like trying to read a book written in hieroglyphics.",
	"I'm not sure what you're asking, but I'm pretty sure it's a brain-teaser.",
	"I'm trying to understand, but it's like trying to balance on a tightrope.",
	"I'm trying to understand, but it's like trying to fit a round peg in a square hole.",
	"I'm not sure what you're asking, but I'm pretty sure it's a red herring.",
	"I'm trying to understand, but it's like trying to put together a jigsaw puzzle with missing pieces.",
	"I'm not sure what you're asking, but I'm pretty sure it's a wild goose chase.",
	"I'm trying to understand, but it's like trying to navigate a minefield.",
	"I'm not sure what you're asking, but I'm pretty sure it's a Gordian knot.",
	"I'm trying to understand, but it's like trying to thread a needle in a dark room.",
	"I'm not sure what you're asking, but I'm pretty sure it's a paradox.",
	"I'm trying to understand, but it's like trying to herd cats.",
	"I'm not sure what you're asking, but I'm pretty sure it's an enigma.",
	"I'm trying to understand, but it's like trying to catch a greased pig.",
	"I'm not sure what you're asking, but I'm pretty sure it's a double-edged sword.",
	"I'm trying to understand, but it's like trying to untangle a ball of yarn.",
	"I'm not sure what you're asking, but I'm pretty sure it's a can of worms.",
	"I'm trying to understand, but it's like trying to find a needle in a haystack of needles.",
	"I'm not sure what you're asking, but I'm pretty sure it's a Pandora's box.",
	"I'm trying to understand, but it's like trying to fit a square peg in a round hole, with gloves on.",
	"I'm not sure what you're asking, but I'm pretty sure it's a catch-22.",
	"I'm trying to understand, but it's like trying to drink from a fire hose.",
	"I'm not sure what you're asking, but I'm pretty sure it's a Chinese puzzle.",
	"I'm trying to understand, but it's like trying to solve a Rubik's cube while blindfolded and wearing gloves.",
	"I'm not sure what you're asking, but I'm pretty sure it's a Sisyphean task.",
	"I'm trying to understand, but it's like trying to find a needle in a haystack of needles, in the dark.",
	"I'm not sure what you're asking, but I'm pretty sure it's a game of cat and mouse.",
	"I'm trying to understand, but it's like trying to build a sandcastle in a hurricane.",
	"I'm not sure what you're asking, but I'm pretty sure it's a Herculean task.",
	"I'm trying to understand, but it's like trying to navigate a labyrinth while drunk.",
	"I'm not sure what you're asking, but I'm pretty sure it's a needle in a haystack of needles, in a dark room, with gloves on.",
	"I'm trying to understand, but it's like trying to put together a puzzle with missing pieces, with one hand tied behind my back.",
	"I'm not sure what you're asking, but I'm pretty sure it's a game of whack-a-mole.",
	"Why was the cyber criminal's computer always cold? It left its Windows open.",
	"How do you know when a hacker is around? Your computer starts to feel violated.",
	"Why do hackers always wear hoodies? So they can keep a low IP address.",
	"How do you stop a cyber attack? With a good firewall and a strong password.",
	"Why do hackers always use Linux? Because Windows is just too easy to hack.",
	"How do you know when a website is secure? When it has a padlock and an 'https' in the URL.",
	"Why was the hacker afraid of the dark web? Because it was full of scary malware.",
	"How do you protect yourself from phishing? By not falling for the bait.",
	"Why do hackers always use strong encryption? Because they have something to hide.",
	"How do you know when your computer is infected with malware? When it starts to act strange.",
	"Why was the hacker's keyboard always sticky? Because they used too much 'Ctrl C' and 'Ctrl V'.",
	"How do you know when a cyber attack is coming? When your computer starts to feel paranoid.",
	"Why did the cyber criminal get arrested? Because they left their laptop open at the library.",
	"How do you stop an online scam? By not falling for the clickbait.",
	"Why was the hacker's computer so slow? Because they had too many viruses.",
	"How do you protect your online identity? By not giving away personal information.",
	"Why did the hacker's computer crash? Because they opened too many windows.",
	"How do you know if a website is fake? By checking the URL and looking for spelling mistakes.",
	"Why was the hacker's computer so hot? Because they left the encryption running overnight.",
	"How do you protect your online privacy? By using a VPN and strong passwords.",
}

func doesFileExist(fileName string) bool {
	_, error := os.Stat(fileName)

	if os.IsNotExist(error) {
		return false
	} else {
		return true
	}
}

func IsValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return (err == nil) && regexp.MustCompile(`^[a-f-0-9]*$`).MatchString(u)
}
