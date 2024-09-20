package commands

import (
	"log"
	"math/rand"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var roastList = []string{
	"Vai levar nas nalgas {}, seu paneleiro do caralho.",
	"Já vi montes de merda mais bonitos que o {}.",
	"{} a tua certidão de nascimento é um pedido de desculpas da Durex.",
	"{} só não és meu filho porque a tua mãe não tinha troco de 20€.",
	"{} calado és um poeta, oh filho.",
	"Bom {}, dá a vez ao teu cu de falar, é que dessa tua boca só sai merda fds.",
	"A ouvir as merdas que o {} diz, preferia estar em casa a ver as ervas do quintal crescerem.",
	"Com a carinha de atrasadinho mental do {}, passava a ser Make-Two-Wishes.",
	"{}, vai maseh mamar na quinta pata do cavalo.",
	"Tu não és estúpido {}, tens é um derrame cerebral cada vez que abres a boca.",
	"A chucha do {} devia de ser feita de amianto.",
	// Novas frases
	"Oh {}, devias receber uma medalha por teres sobrevivido com esse QI até agora.",
	"Quando o {} nasceu, a parteira pediu reembolso.",
	"Se o {} fosse mais burro, tinha de ser regado para não secar.",
	"{}, tens uma cara que só uma mãe podia amar… e mesmo ela já tem dúvidas.",
	"Sabes {}, se a ignorância fosse música, fazias uma sinfonia.",
	"Tu és a prova viva de que até o lixo se recicla, {}.",
	"É impressionante {}, tens tanto cérebro como um balde vazio.",
	"Se a tua vida fosse um filme {}, era daqueles que até a pipoca adormecia.",
	"O {} é tipo o Titanic, parecia promissor, mas é só uma tragédia ambulante.",
	"{}, já vi comas alcoólicos a tomarem decisões mais conscientes que tu.",
	"Quando olham para ti, {} as pessoas finalmente entendem porque a contraceção existe.",
	"A tua família toda tem a foto no frigorífico, {}, mas só para lembrar que a porta tem de ficar fechada.",
	"{}, já foste desejado? Ah espera, nem a tua mãe conseguiu fingir isso.",
	"Quando Deus fez o {}, era claramente segunda-feira e ele ainda estava de ressaca.",
	"{}, a única coisa pior do que olharem para ti, é perceberem que ainda estás vivo.",
	"Os médicos deviam ter-te deixado na incubadora, {}, e depois apagado a luz.",
}

var RoastCommand = &discordgo.ApplicationCommand{
	Name:        "roast",
	Description: "Sends a roast message to the target",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionUser,
			Name:        "user",
			Description: "The user to roast",
			Required:    true, // The parameter is mandatory
		},
	},
}

func RoastHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Extract the user parameter from the interaction
	userOption := i.ApplicationCommandData().Options[0].UserValue(s)

	// Select random roast message and replace placeolder by user mention tag
	randomRoast := strings.Replace(roastList[rand.Intn(len(roastList))], "{}", userOption.Mention(), 1)

	// Answer with 'World!" to complete user message
	log.Printf("Roasting %v with '%v' by %v \n", userOption, randomRoast, i.Member.User.Username)
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: randomRoast,
		},
	})

	// Handle error
	if err != nil {
		log.Fatal(err)
	}
}
