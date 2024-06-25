package main

import (
	// go-mc

	"encoding/json"
	"io/fs"
	"mcquery/frontend"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Tnze/go-mc/bot"
	"github.com/Tnze/go-mc/chat"
	"github.com/google/uuid"
	"github.com/pelletier/go-toml/v2"
)

type status struct {
	Description chat.Message
	Players     struct {
		Max    int
		Online int
		Sample []struct {
			ID   uuid.UUID
			Name string
		}
	}
	Version struct {
		Name     string
		Protocol int
	}
	Delay time.Duration
}

type Config struct {
	Web struct {
		Address string
	}
	Minecraft struct {
		Address string
	}
}

func (c *Config) Load(path string) {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	err = toml.NewDecoder(file).Decode(c)
	if err != nil {
		panic(err)
	}
}

func main() {
	var config Config
	config.Load("config.toml")

	mux := http.NewServeMux()

	mux.HandleFunc("GET /data", func(w http.ResponseWriter, r *http.Request) {
		res, _, err := bot.PingAndListTimeout(config.Minecraft.Address, time.Second*5)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		var status status
		err = json.Unmarshal(res, &status)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var playerList []string
		for _, player := range status.Players.Sample {
			playerList = append(playerList, player.Name)
		}

		data := map[string]any{
			"address":     config.Minecraft.Address,
			"version":     status.Version.Name,
			"players":     status.Players.Online,
			"maxPlayers":  status.Players.Max,
			"playerList":  playerList,
			"description": strings.TrimSpace(status.Description.String()),
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	// subfs
	fs, err := fs.Sub(frontend.Dist, "dist")
	if err != nil {
		panic(err)
	}
	mux.Handle("/", http.FileServer(http.FS(fs)))

	err = http.ListenAndServe(config.Web.Address, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")

		if r.Method == http.MethodOptions {
			w.Header().Set("Access-Control-Allow-Methods", "GET")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			w.WriteHeader(http.StatusNoContent)
			return
		}

		mux.ServeHTTP(w, r)
	}))
	if err != nil {
		panic(err)
	}
}
