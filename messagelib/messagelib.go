package messagelib

import (
	"fmt"
	"log"
	"math"
	"strings"

	"github.com/JoeParrinello/brokerbot/persistencelib"
	"github.com/bwmarrin/discordgo"
)

var (
	test          bool   = false
	messagePrefix string = "TEST"
)

// TickerValue passes values of fetched content.
type TickerValue struct {
	Ticker   string
	Value    float32
	Change   float32
	MiscText string
}

// EnterTestModeWithPrefix enables extra log prefixes to identify a test server.
func EnterTestModeWithPrefix(prefix string) {
	test = true
	messagePrefix = prefix
	log.Printf("BrokerBot running in test mode with prefix: %q", prefix)
}

// ExitTestMode disables extra log prefixes to identify a test server.
func ExitTestMode() {
	test = false
}

// SendMessage sends a plaintext message to a Discord channel.
func SendMessage(s *discordgo.Session, channelID string, msg string) *discordgo.Message {
	msg = fmt.Sprintf("%s%s", getMessagePrefix(), msg)
	message, err := s.ChannelMessageSend(channelID, msg)
	if err != nil {
		log.Printf("failed to send message %q to discord: %v", msg, err)
	}
	return message
}

// SendMessageEmbed sends a rich "embed" message to a Discord channel.
func SendMessageEmbed(s *discordgo.Session, channelID string, msg *discordgo.MessageEmbed) *discordgo.Message {
	message, err := s.ChannelMessageSendEmbed(channelID, msg)
	if err != nil {
		log.Printf("failed to send message %+v to discord: %v", msg, err)
	}
	return message
}

// CreateMessageEmbed creates a rich Discord "embed" message
func CreateMessageEmbed(tickerValue *TickerValue) *discordgo.MessageEmbed {
	return createMessageEmbedWithPrefix(tickerValue, getTestServerID())
}

func createMessageEmbedWithPrefix(tickerValue *TickerValue, prefix string) *discordgo.MessageEmbed {
	if tickerValue == nil {
		return nil
	}

	mesg := fmt.Sprintf("Latest Quote: $%.2f", tickerValue.Value)
	if !math.IsNaN(float64(tickerValue.Change)) && tickerValue.Change != 0 {
		mesg = fmt.Sprintf("%s (%.2f%%)", mesg, tickerValue.Change)
	}
	if len(tickerValue.MiscText) > 0 {
		mesg = fmt.Sprintf("%s\n%s", mesg, tickerValue.MiscText)
	}
	return &discordgo.MessageEmbed{
		Title:       tickerValue.Ticker,
		URL:         fmt.Sprintf("https://www.google.com/search?q=%s", tickerValue.Ticker),
		Description: mesg,
		Footer: &discordgo.MessageEmbedFooter{
			Text: prefix,
		},
	}
}

// CreateMultiMessageEmbed will return an embedded message for multiple tickers.
func CreateMultiMessageEmbed(tickers []*TickerValue) *discordgo.MessageEmbed {
	return createMultiMessageEmbedWithPrefix(tickers, getTestServerID())
}

func createMultiMessageEmbedWithPrefix(tickers []*TickerValue, prefix string) *discordgo.MessageEmbed {
	messageFields := make([]*discordgo.MessageEmbedField, len(tickers))
	for i, ticker := range tickers {
		messageFields[i] = createMessageEmbedField(ticker)
	}
	return &discordgo.MessageEmbed{
		Fields: messageFields,
		Footer: &discordgo.MessageEmbedFooter{
			Text: prefix,
		},
	}
}

func createMessageEmbedField(tickerValue *TickerValue) *discordgo.MessageEmbedField {
	if math.IsNaN(float64(tickerValue.Value)) || tickerValue.Value == 0.0 {
		return &discordgo.MessageEmbedField{
			Name:   tickerValue.Ticker,
			Value:  fmt.Sprintf("%s - %s", "No Data", tickerValue.MiscText),
			Inline: false,
		}
	}

	mesg := fmt.Sprintf("$%.2f", tickerValue.Value)
	if !math.IsNaN(float64(tickerValue.Change)) && tickerValue.Change != 0 {
		mesg = fmt.Sprintf("%s (%.2f%%)", mesg, tickerValue.Change)
	}
	if len(tickerValue.MiscText) > 0 {
		mesg = fmt.Sprintf("%s\n%s", mesg, tickerValue.MiscText)
	}
	return &discordgo.MessageEmbedField{
		Name:   tickerValue.Ticker,
		Value:  mesg,
		Inline: false,
	}
}

func getMessagePrefix() string {
	if test {
		return messagePrefix + ": "
	}
	return ""
}

func getTestServerID() string {
	if test {
		return messagePrefix
	}
	return ""
}

// RemoveMentions removes any @ mentions from a message slice.
func RemoveMentions(s []string) (ret []string) {
	for _, v := range s {
		if !strings.HasPrefix(v, "@") {
			ret = append(ret, v)
		}
	}
	return
}

// CanonicalizeMessage upcases each field in a message slice.
func CanonicalizeMessage(s []string) (ret []string) {
	for _, v := range s {
		ret = append(ret, strings.ToUpper(v))
	}
	return
}

// ExpandAliases takes a string that contains an alias of format "?<alias>" and replaces the alias with the valid ticker string.
func ExpandAliases(s []string) (ret []string) {
	for _, v := range s {
		if strings.HasPrefix(v, "?") {
			expanded := persistencelib.ExpandAlias(v)
			if expanded != nil {
				ret = append(ret, expanded...)
			} else {
				ret = append(ret, v)
			}
		} else {
			ret = append(ret, v)
		}
	}
	return
}

// DedupeSlice returns a list of unique tickers from the provided string slice.
func DedupeSlice(s []string) (ret []string) {
	seen := make(map[string]bool)
	for _, v := range s {
		if _, exists := seen[v]; !exists {
			seen[v] = true
			ret = append(ret, v)
		}
	}
	return
}
