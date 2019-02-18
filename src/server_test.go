package main

import (
	"fmt"
	"github.com/nlopes/slack"
	"net/http"
	"testing"
)

func TestCommandResponseHello(t *testing.T) {
	s := slack.SlashCommand{
		Token:      "Gq0rL1ZBphpmFvqYcGGwLMQF",
		TeamID:     "T5BTQ05J4",
		TeamDomain: "ube-computerclub",
		Command:    "/hello",
		Text:       "",
	}

	code, params := commandResponse(s)

	if code != http.StatusOK {
		t.Fatal("Failed test: /hello expect StatusOK")
	}
	if params.Text != "Hello" {
		t.Fatal("Failed test: /hello expect StatusOK")
	}

	s.Token = "THIS_IS_DUMMY_TOKEN"

	code, params = commandResponse(s)
	fmt.Println(code)

	if code != http.StatusInternalServerError {
		t.Fatal("Failed test: /hello expect StatusInternalServerError")
	}
}
