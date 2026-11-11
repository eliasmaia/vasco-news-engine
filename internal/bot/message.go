package bot

// NewsMessage represents the news DTO (Data Transfer Object)
type NewsMessage struct {
	Title 	string
	Link  	string
	Source 	string
}

// Format isolates the responsabilidade of drawing the message template
func (n NewsMessage) Format() string {
	return  "💢 *NOVIDADE NO GIGANTE* 💢\n\n" +
		n.Title + "\n\n" +
		"📍 Fonte: " + n.Source + "\n" +
		"🔗 [Leia mais](" + n.Link + ")"
}