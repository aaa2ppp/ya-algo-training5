package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
)

type Team struct {
	Name       string // not need for this task
	Games      int
	Goals      int
	ScoreOpens int
}

type Player struct {
	Name          string // not need for this task
	Goals         int
	ScoreOpens    int
	Team          *Team
	GoalsOnMinute []int
}

type Game struct {
	Team             [2]*Team
	Goals            [2]int
	curTeam          int
	ScoreOpensMinute int
	ScoreOpensPlayer *Player
}

type Teams map[string]*Team

// Create if not exists. Alway returns team with name.
func (ts Teams) Create(name string) *Team {
	return ts.Get(name)
}

func (ts Teams) Get(name string) *Team {
	it := ts[name]
	if it == nil {
		it = &Team{Name: name}
		ts[name] = it
	}
	return it
}

type Players map[string]*Player

func (ps Players) Create(name string, team *Team) *Player {
	it := ps.Get(name)
	it.Team = team
	return it
}

func (ps Players) Get(name string) *Player {
	it := ps[name]
	if it == nil {
		it = &Player{
			Name:          name,
			GoalsOnMinute: make([]int, 91),
		}
		ps[name] = it
	}
	return it
}

func run(in io.Reader, out io.Writer) error {
	sc := bufio.NewScanner(in)
	bw := bufio.NewWriter(out)
	defer bw.Flush()

	teams := Teams{}
	players := Players{}

	var game Game // current game
	for sc.Scan() {
		b := bytes.TrimSpace(sc.Bytes())

		// XXX это порно!
		switch {
		case b[0] == '"':
			// "<Название 1-й команды>" - "<Название 2-й команды>" <Счет 1-й команды>:<Счет 2-й команды>

			var (
				teamName1 string
				teamName2 string
				n1, n2    int
			)

			line := bytes.NewReader(b)
			_, err := fmt.Fscanf(line, "%q - %q %d:%d", &teamName1, &teamName2, &n1, &n2)
			if err != nil {
				return err
			}

			team1 := teams.Create(teamName1)
			team1.Games++
			team1.Goals += n1

			team2 := teams.Create(teamName2)
			team2.Games++
			team2.Goals += n2

			game = Game{
				Team:             [2]*Team{team1, team2},
				Goals:            [2]int{n1, n2},
				ScoreOpensMinute: -1,
			}

			for game.curTeam < 2 && game.Goals[game.curTeam] == 0 {
				game.curTeam++
			}

		case b[len(b)-1] == '\'':
			// <Автор x-го забитого мяча x-й команды> <Минута, на которой был забит мяч>'

			var (
				playerName string
				minute     int
				err        error
			)

			i := bytes.LastIndexByte(b, ' ')
			playerName = string(bytes.TrimSpace(b[:i]))
			minute, err = strconv.Atoi(string(bytes.TrimSpace(b[i : len(b)-1])))
			if err != nil {
				return err
			}

			player := players.Create(playerName, game.Team[game.curTeam])
			player.Goals++
			player.GoalsOnMinute[minute]++

			if game.ScoreOpensMinute == -1 || minute < game.ScoreOpensMinute {
				game.ScoreOpensMinute = minute
				game.ScoreOpensPlayer = player
			}

			game.Goals[game.curTeam]--
			for game.curTeam < 2 && game.Goals[game.curTeam] == 0 {
				game.curTeam++
			}

			if game.curTeam == 2 {
				game.ScoreOpensPlayer.ScoreOpens++
				game.ScoreOpensPlayer.Team.ScoreOpens++
			}

		case bytes.HasPrefix(b, []byte("Total goals for ")):
			// — количество голов, забитое данной командой за все матчи.

			b := bytes.TrimPrefix(b, []byte("Total goals for "))
			teamName, err := strconv.Unquote(string(b))
			if err != nil {
				return err
			}

			team := teams.Get(teamName)
			fmt.Fprintln(bw, team.Goals)

		case bytes.HasPrefix(b, []byte("Mean goals per game for ")):
			// — среднее количество голов, забиваемое данной командой за один матч.
			// Гарантирутся, что к моменту подачи такого запроса команда уже сыграла хотя бы один матч.

			b := bytes.TrimPrefix(b, []byte("Mean goals per game for "))
			teamName, err := strconv.Unquote(string(b))
			if err != nil {
				return err
			}

			team := teams.Get(teamName)
			fmt.Fprintf(bw, "%g\n", float64(team.Goals)/float64(team.Games))

		case bytes.HasPrefix(b, []byte("Total goals by ")):
			// — количество голов, забитое данным игроком за все матчи.

			b = bytes.TrimPrefix(b, []byte("Total goals by "))
			playerName := string(b)

			player := players.Get(playerName)
			fmt.Fprintln(bw, player.Goals)

		case bytes.HasPrefix(b, []byte("Mean goals per game by ")):
			// — среднее количество голов, забиваемое данным игроком за один матч его команды.
			// Гарантирутся, что к моменту подачи такого запроса игрок уже забил хотя бы один гол.

			b := bytes.TrimPrefix(b, []byte("Mean goals per game by "))
			playerName := string(b)

			player := players.Get(playerName)
			fmt.Fprintf(bw, "%g\n", float64(player.Goals)/float64(player.Team.Games))

		case bytes.HasPrefix(b, []byte("Goals on minute ")):
			// — количество голов, забитых данным игроком ровно на указанной минуте матча.

			b = bytes.TrimPrefix(b, []byte("Goals on minute "))
			pos := bytes.IndexByte(b, ' ')
			minute, err := strconv.Atoi(string(b[:pos]))
			if err != nil {
				return err
			}

			b = bytes.TrimPrefix(b[pos+1:], []byte("by "))
			playerName := string(b)

			player := players.Get(playerName)
			fmt.Fprintln(bw, player.GoalsOnMinute[minute])

		case bytes.HasPrefix(b, []byte("Goals on first ")):
			// — количество голов, забитых данным игроком на минутах с первой по T-ю включительно.

			b = bytes.TrimPrefix(b, []byte("Goals on first "))
			pos := bytes.IndexByte(b, ' ')
			minutes, err := strconv.Atoi(string(b[:pos]))
			if err != nil {
				return err
			}

			b = bytes.TrimPrefix(b[pos+1:], []byte("minutes by "))
			playerName := string(b)

			player := players.Get(playerName)
			goals := 0
			for i := 1; i <= minutes; i++ {
				goals += player.GoalsOnMinute[i]
			}

			fmt.Fprintln(bw, goals)

		case bytes.HasPrefix(b, []byte("Goals on last ")):
			// — количество голов, забитых данным игроком на минутах с (91 - T)-й по 90-ю включительно.

			b = bytes.TrimPrefix(b, []byte("Goals on last "))
			pos := bytes.IndexByte(b, ' ')
			minutes, err := strconv.Atoi(string(b[:pos]))
			if err != nil {
				return err
			}

			b = bytes.TrimPrefix(b[pos+1:], []byte("minutes by "))
			playerName := string(b)

			player := players.Get(playerName)
			goals := 0
			for i := 91 - minutes; i <= 90; i++ {
				goals += player.GoalsOnMinute[i]
			}

			fmt.Fprintln(bw, goals)

		case b[len(b)-1] == '"' && bytes.HasPrefix(b, []byte("Score opens by ")):
			// — сколько раз данная команда открывала счет в матче.

			b = bytes.TrimPrefix(b, []byte("Score opens by "))
			teamName, err := strconv.Unquote(string(b))
			if err != nil {
				return err
			}

			team := teams.Get(teamName)
			fmt.Fprintln(bw, team.ScoreOpens)

		case bytes.HasPrefix(b, []byte("Score opens by ")):
			// — сколько раз данный игрок открывал счет в матче.

			b = bytes.TrimPrefix(b, []byte("Score opens by "))
			playerName := string(b)

			player := players.Get(playerName)
			fmt.Fprintln(bw, player.ScoreOpens)
		}
	}

	return nil
}

var _, debugEnable = os.LookupEnv("DEBUG")

func main() {
	_ = debugEnable
	err := run(os.Stdin, os.Stdout)
	if err != nil {
		log.Fatal(err)
	}
}
