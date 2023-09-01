package yandexspeller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
)

type RequestSpeller struct {
	Lang    string   `json:"lang,omitempty"`
	Format  string   `json:"format,omitempty"`
	Text    []string `json:"text"`
	Options int      `json:"options,omitempty"`
}

type ResponseSpeller struct {
	Error         map[string]string `json:"error,omitempty"`
	Word          string            `json:"word,omitempty"`
	CorrectedWord []string          `json:"s,omitempty"`
}

const serviceAddress = "https://speller.yandex.net/services/spellservice.json/checkTexts?text="

func YandexSpeller(log *slog.Logger, inputNotes map[string][]string, userName *string) (outputWithCorrectNotes map[string][]string) {
	reqSpeller := RequestSpeller{}
	reqSpeller.Options = 14
	client := http.Client{}
	resSpeller := [][]ResponseSpeller{}

	var sb strings.Builder

	outputWithCorrectNotes = make(map[string][]string)

	notes := inputNotes[*userName]

	for _, note := range notes {
		reqSpeller.Text = strings.Split(note, " ")
		request := fmt.Sprint(serviceAddress + strings.Join(reqSpeller.Text, "&text="))
		input, err := json.Marshal(reqSpeller)
		if err != nil {
			log.Error("YandexSpeller", "json.Unmarshal", err)

			return
		}

		yandexRes, err := client.Post(request, "application/json", bytes.NewReader(input))
		if err != nil {
			log.Error("YandexSpeller", "fn:client.Post", err)

			return
		}

		if err = json.NewDecoder(yandexRes.Body).Decode(&resSpeller); err != nil {
			log.Error("YandexSpeller", "fn:json.Decode", err)

			return
		}

		for i, res := range resSpeller {
			if len(res) == 0 {
				sb.WriteString(reqSpeller.Text[i])
				sb.WriteString(" ")

				continue
			}
			sb.WriteString(res[0].CorrectedWord[0])
			switch {
			case strings.Contains(reqSpeller.Text[i], ","):
				sb.WriteString(",")
			case strings.Contains(reqSpeller.Text[i], "."):
				sb.WriteString(".")
			case strings.Contains(reqSpeller.Text[i], ":"):
				sb.WriteString(":")
			case strings.Contains(reqSpeller.Text[i], "!"):
				sb.WriteString("!")
			}
			sb.WriteString(" ")
		}

		s, _ := strings.CutSuffix(sb.String(), " ")
		sb.Reset()

		if outputWithCorrectNotes[*userName] == nil {
			outputWithCorrectNotes[*userName] = make([]string, 0, len(notes))
		}
		outputWithCorrectNotes[*userName] = append(outputWithCorrectNotes[*userName], s)

	}

	return outputWithCorrectNotes
}
