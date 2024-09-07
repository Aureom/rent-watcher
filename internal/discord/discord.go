package discord

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"rent-watcher/internal/models"
	"rent-watcher/internal/notifier"
	"strings"
	"time"
)

type Discord struct {
	session *discordgo.Session
	channel string
}

func New(token, channel string) (notifier.Notifier, error) {
	session, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, err
	}

	if err := session.Open(); err != nil {
		return nil, err
	}

	return &Discord{
		session: session,
		channel: channel,
	}, nil
}

func (d *Discord) NotifyNewProperty(p *models.Property) error {
	embed := &discordgo.MessageEmbed{
		Title: "New Property Listed! 🏠",
		URL:   fmt.Sprintf("https://www.arantesimoveis.com%s", p.ID),
		Color: 0x00ff00, // Green color
		Fields: []*discordgo.MessageEmbedField{
			{Name: "📍 Address", Value: fmt.Sprintf("%s, %s, %s", p.Logradouro, p.Bairro, p.Cidade), Inline: false},
			{Name: "💰 Price", Value: fmt.Sprintf("R$ %s", p.Price), Inline: true},
			{Name: "🏘️ Type", Value: p.TipoImovel, Inline: true},
			{Name: "🛏️ Bedrooms", Value: p.Quartos, Inline: true},
			{Name: "🚿 Bathrooms", Value: p.Banheiros, Inline: true},
			{Name: "🚗 Parking Spaces", Value: p.Garagens, Inline: true},
			{Name: "📏 Area", Value: fmt.Sprintf("%s m²", p.Metragem), Inline: true},
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Property ID: " + strings.TrimPrefix(p.ID, "/detalhes/"),
		},
		Timestamp: time.Now().Format(time.RFC3339),
	}

	if p.Suites != "" && p.Suites != "0" {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   "🛁 Suites",
			Value:  p.Suites,
			Inline: true,
		})
	}

	if isValidURL(p.FirstPhoto) {
		embed.Image = &discordgo.MessageEmbedImage{
			URL: p.FirstPhoto,
		}
	} else {
		log.Printf("Invalid or missing photo URL for property ID %s", p.ID)
	}

	_, err := d.session.ChannelMessageSendEmbed(d.channel, embed)

	if err != nil {
		log.Printf("Error sending Discord embed: %v", err)
	}

	return err
}

func (d *Discord) Close() error {
	return d.session.Close()
}

func isValidURL(url string) bool {
	return strings.HasPrefix(url, "http")
}
