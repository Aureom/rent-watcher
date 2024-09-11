package discord

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"rent-watcher/internal/models"
	"rent-watcher/internal/notifier"
)

type Discord struct {
	session *discordgo.Session
	channel string
}

func New(token, channel string) (notifier.Notifier, error) {
	session, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, fmt.Errorf("failed to create Discord session: %w", err)
	}

	if err := session.Open(); err != nil {
		return nil, fmt.Errorf("failed to open Discord session: %w", err)
	}

	return &Discord{
		session: session,
		channel: channel,
	}, nil
}

func (d *Discord) NotifyNewProperty(p *models.Property) error {
	embed := &discordgo.MessageEmbed{
		Title:       "🏠 New Property Alert!",
		Description: fmt.Sprintf("[%s - %s, %s](https://www.arantesimoveis.com/detalhes/%s)", p.Logradouro, p.Bairro, p.Cidade, p.ID),
		URL:         fmt.Sprintf("https://www.arantesimoveis.com/detalhes/%s", p.ID),
		Color:       0x00bfff,
		Fields:      createEmbedFields(p),
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Property ID: " + strings.TrimPrefix(p.ID, "/detalhes/"),
		},
		Timestamp: time.Now().Format(time.RFC3339),
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
		return fmt.Errorf("error sending Discord embed: %w", err)
	}

	return nil
}

func createEmbedFields(p *models.Property) []*discordgo.MessageEmbedField {
	fields := []*discordgo.MessageEmbedField{
		{Name: "💰 Price", Value: formatCurrency(p.Price), Inline: true},
		{Name: "🏢 Condo Fee", Value: formatCurrency(p.Condominio), Inline: true},
		{Name: "💵 Total Price", Value: formatCurrency(p.TotalPrice), Inline: true},
		{Name: "📍 Distance", Value: formatDistance(p.DistanceMeters), Inline: true},
		{Name: "🏘️ Type", Value: p.TipoImovel, Inline: true},
		{Name: "📏 Area", Value: fmt.Sprintf("%s m²", p.Metragem), Inline: true},
		{Name: "🛏️ Bedrooms", Value: p.Quartos, Inline: true},
		{Name: "🚿 Bathrooms", Value: p.Banheiros, Inline: true},
		{Name: "🚗 Parking", Value: p.Garagens, Inline: true},
	}

	if p.Suites != "" && p.Suites != "0" {
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   "🛁 Suites",
			Value:  p.Suites,
			Inline: true,
		})
	}

	return fields
}

func formatCurrency(value string) string {
	return fmt.Sprintf("R$ %s", strings.TrimSpace(strings.TrimPrefix(value, "R$")))
}

func formatDistance(meters int) string {
	if meters < 1000 {
		return fmt.Sprintf("%d m", meters)
	}
	return fmt.Sprintf("%.1f km", float64(meters)/1000)
}

func (d *Discord) Close() error {
	return d.session.Close()
}

func isValidURL(url string) bool {
	return strings.HasPrefix(url, "http")
}
