package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
)

var playerList map[string]string //id username

type GameMap struct {
	Name  string   `json:"name"`
	Point []string `json:"point"`
	URL   string   `json:"url"`
}

var maps = map[string]string{
	"Bind":  "https://static.wikia.nocookie.net/valorant/images/2/23/Loading_Screen_Bind.png",
	"Haven": "https://static.wikia.nocookie.net/valorant/images/7/70/Loading_Screen_Haven.png",
	"Split": "https://static.wikia.nocookie.net/valorant/images/d/d6/Loading_Screen_Split.png",
	"Icebox": "https://static.wikia.nocookie.net/valorant/images/3/34/Loading_Icebox.png",
	"Ascent": "https://static.wikia.nocookie.net/valorant/images/e/e7/Loading_Screen_Ascent.png",
	"Breeze": "https://static.wikia.nocookie.net/valorant/images/1/1e/Valorant_Loading_Breeze.png",
}
var MapPoints = map[string][]string{
	"Bind":  {"A", "B"},
	"Haven": {"A", "B", "C"},
	"Split": {"A", "B"},
	"Icebox" : {"A", "B"},
	"Ascent": {"A", "B"},
	"Breeze": {"A", "B"},
}

var (
	version string
)

// Variables used for command line parameters
var (
	Token string
)

func init() {

	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.Parse()
}

func main() {

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}
	rand.Seed(time.Now().UnixNano())
	playerList = make(map[string]string)
	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the autenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}
	/*if set, err := s.UserSettings(); err != nil {
		fmt.Print(err)
	} else {
		fmt.Print(set.EnableTtsCommand)
	}*/

	if strings.HasPrefix(m.Content, "+g") {
		command := strings.Split(m.Content, " ")
		response := &discordgo.MessageEmbed{
			Color: 0xffcc00,
		}

		respAuthor := &discordgo.MessageEmbedAuthor{
			Name:    m.Author.Username,
			IconURL: m.Author.AvatarURL(""),
		}
		response.Author = respAuthor

		footer := &discordgo.MessageEmbedFooter{
			Text: "Karéjone Training Bot " + version + " - Created by davzo",
		}
		response.Footer = footer
		if len(command) > 1 {
			switch command[1] {
			case "add":
				if len(command) > 2 {
					userMention := m.Message.Mentions
					response.Title = "Modification de la liste d'attente"
					for _, u := range userMention {

						if _, exist := playerList[u.ID]; exist {
							response.Description += u.Username + " existe déjà dans la liste d'attente\n"
						} else {
							response.Description += u.Username + " à été ajouté à la liste d'attente\n"
							playerList[u.ID] = u.Username
						}
					}
				} else {
					if _, exist := playerList[m.Author.ID]; exist {
						response.Title = "Le joueur existe déjà dans la liste"
					} else {
						response.Title = "Le joueur a été ajouté à la liste"
						playerList[m.Author.ID] = m.Author.Username
					}
				}

			case "del":
				if len(command) > 2 {
					userMention := m.Message.Mentions
					response.Title = "Modification de la liste d'attente"
					for _, u := range userMention {
						if _, exist := playerList[u.ID]; exist {
							response.Description += u.Username + " a été supprimé de la liste d'attente\n"
							delete(playerList, u.ID)
						} else {
							response.Description += u.Username + " n'existe pas dans la liste d'attente\n"
							playerList[u.ID] = u.Username
						}
					}
				} else {
					if _, exist := playerList[m.Author.ID]; exist {
						response.Title = "Le joueur a été supprimé de la liste"
						delete(playerList, m.Author.ID)
					} else {
						response.Title = "Le joueur n'existe pas dans la liste"

					}
				}

			case "list":
				response.Title = "Liste des joueurs en attente"
				if len(playerList) != 0 {
					response.Description = stringListPlayer(playerList)
				} else {
					response.Description = "Aucun joueur n'est dans la liste d'attente"

				}
			case "flush":
				response.Title = "La liste des joueurs en attente a été purgée"
				for k := range playerList {
					delete(playerList, k)
				}

			case "generate":
				response.Title = "Génération de la partie d'entrainement"
				if len(playerList) >= 2 {
					mapKeyRand := shuffleList(maps)[0]
					imageMap := &discordgo.MessageEmbedImage{
						URL: maps[mapKeyRand],
					}
					response.Image = imageMap
					fmt.Println(mapKeyRand)
					fmt.Println(MapPoints[mapKeyRand])
					pointRandPicked := pickAnElement(MapPoints[mapKeyRand])
					response.Description = "**" + mapKeyRand + " - " + pointRandPicked + "**" + "\n"
					idShuffled := shuffleList(playerList)
					var attPlayers []string
					var defPlayers []string
					var specPlayers []string
					fmt.Print(len(idShuffled))
					if len(idShuffled)%2 == 1 {
						specPlayers = append(specPlayers, playerList[idShuffled[len(idShuffled)-1]])
						idShuffled[len(idShuffled)-1] = ""
						idShuffled = idShuffled[:len(idShuffled)-1]
					}
					fmt.Print(len(idShuffled))
					for i, s := range idShuffled {

						if i%2 == 0 {
							attPlayers = append(attPlayers, playerList[s])
						} else {
							defPlayers = append(defPlayers, playerList[s])
						}
					}

					field := &discordgo.MessageEmbedField{
						Name:  "Attaque",
						Value: strings.Join(attPlayers, ", "),
					}
					response.Fields = append(response.Fields, field)
					field = &discordgo.MessageEmbedField{
						Name:  "Défense",
						Value: strings.Join(defPlayers, ", "),
					}

					response.Fields = append(response.Fields, field)
					if len(specPlayers) > 0 {
						field = &discordgo.MessageEmbedField{
							Name:  "Spectateur",
							Value: strings.Join(specPlayers, ", "),
						}
						response.Fields = append(response.Fields, field)
					}
				} else {
					response.Description = "Pas assez de joueur"
				}

			default:
				response.Title = "Karéjone Training Bot Generator :woozy_face:"
				response.Description =
					command[1] + ` est commande non reconnue"
						Syntaxe : +g <command>"
						**Help**
						:arrow_right:` + "`add [@mention]`" + `: Ajoute des joueurs dans la liste d'attente, s'il n'y a pas de mention, cela ajoute le joueur qui a tapé la commande
						:arrow_right:` + "`del [@mention]`" + `: Supprime des joueurs de la liste d'attente, s'il n'y a pas de mention, cela supprime le joueur qui a tapé la commande
						:arrow_right:` + "`list`" + `: liste les joueurs en attente
						:arrow_right:` + "`flush`" + `: purge la liste d'attente
						:arrow_right:` + "`generate`" + `: génère une configuration de partie d'entrainement`
			}

		} else {
			response.Title = "Karéjone Training Bot Generator :woozy_face:"
			response.Description =
				`Syntaxe : ` + "\n`+g <command>`" + `
				**Help**
				:arrow_right:` + "`add [@mention]`" + `: Ajoute des joueurs dans la liste d'attente, s'il n'y a pas de mention, cela ajoute le joueur qui a tapé la commande
				:arrow_right:` + "`del [@mention]`" + `: Supprime des joueurs de la liste d'attente, s'il n'y a pas de mention, cela supprime le joueur qui a tapé la commande
				:arrow_right:` + "`list`" + `: liste les joueurs en attente
				:arrow_right:` + "`flush`" + `: purge la liste d'attente
				:arrow_right:` + "`generate`" + `: génère une configuration de partie d'entrainement`
		}
		s.ChannelMessageSendEmbed(m.ChannelID, response)
		if err := s.ChannelMessageDelete(m.ChannelID, m.Message.ID); err != nil {
			fmt.Print(err.Error())
		}
	}
}

func stringListPlayer(m map[string]string) string {
	values := make([]string, 0, len(m))
	for k := range m {
		values = append(values, m[k])
	}
	return strings.Join(values, "\n")
}

func pickAnElement(a []string) string {
	randomIndex := rand.Intn(len(a))
	pick := a[randomIndex]
	return pick
}

func shuffleList(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
		fmt.Println(keys)
	}
	rand.Shuffle(len(keys), func(i, j int) {
		keys[i], keys[j] = keys[j], keys[i]
	})
	return keys
}
