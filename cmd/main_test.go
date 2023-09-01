package main

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"testTask/internal/config"
	"testTask/internal/db/postgres/auth"
	"testTask/internal/db/postgres/notes"
	addnote "testTask/internal/handlers/addNote"
	"testTask/internal/handlers/adduser"
	outputnote "testTask/internal/handlers/outputNote"
	"testTask/internal/logger"

	_ "github.com/lib/pq"
)

func TestServer(t *testing.T) {
	log := logger.SetupLogger()

	cfg, err := config.CreateConfig(log)
	if err != nil {
		t.Fatal(err)
	}

	authDB, err := auth.ConnectUserDB(log, cfg.AddressDB)
	if err != nil {
		t.Fatal(err)
	}

	notesDB, err := notes.ConnectNoteDB(log, cfg.AddressDB)
	if err != nil {
		t.Fatal(err)
	}

	router := http.NewServeMux()
	router.HandleFunc("/droptest", func(w http.ResponseWriter, r *http.Request) {
		db, err := sql.Open("postgres", cfg.AddressDB)
		if err != nil {
			t.Fatal(err)
		}
		defer db.Close()

		_, err = db.Exec("DROP TABLE note")
		if err != nil {
			t.Fatal(err)
		}
	})

	router.HandleFunc("/adduser", adduser.AddUserHandler(log, authDB))
	router.HandleFunc("/addnote", addnote.AddNoteHandler(log, authDB, notesDB))
	router.HandleFunc("/outbyname", outputnote.OutputByNameHandler(log, authDB, notesDB))

	server := httptest.NewServer(router)
	defer server.Close()

	dropTable, err := http.NewRequest(http.MethodGet, server.URL+"/droptest", nil)
	if err != nil {
		t.Fatal(err)
	}

	bodyAddNote := strings.NewReader(`{"username":"testuser","body":"Здравствуте, миня завут Ашот"}`)

	addNote, err := http.NewRequest(http.MethodPost, server.URL+"/addnote", bodyAddNote)
	if err != nil {
		t.Fatal(err)
	}
	addNote.SetBasicAuth("testuser", "secret")

	client := server.Client()

	resAddNote, err := client.Do(addNote)
	if err != nil {
		t.Fatal(err)
	}
	defer resAddNote.Body.Close()

	bodyAddNote2 := strings.NewReader(`{"username":"testuser","body":"Прывет, я хачу сспать"}`)

	addNote2, err := http.NewRequest(http.MethodPost, server.URL+"/addnote", bodyAddNote2)
	if err != nil {
		t.Fatal(err)
	}
	addNote2.SetBasicAuth("testuser", "secret")

	resAddNote2, err := client.Do(addNote2)
	if err != nil {
		t.Fatal(err)
	}
	defer resAddNote2.Body.Close()

	bodyOutByName := strings.NewReader(`{"username":"testuser"}`)

	outByName, err := http.NewRequest(http.MethodPost, server.URL+"/outbyname", bodyOutByName)
	if err != nil {
		t.Fatal(err)
	}
	outByName.SetBasicAuth("testuser", "secret")

	resOutByName, err := client.Do(outByName)
	if err != nil {
		t.Fatal(err)
	}
	defer resOutByName.Body.Close()

	outAll, err := http.NewRequest(http.MethodPost, server.URL+"/outall", bodyOutByName)
	if err != nil {
		t.Fatal(err)
	}
	outAll.SetBasicAuth("testuser", "secret")

	resOutAll, err := client.Do(outAll)
	if err != nil {
		t.Fatal(err)
	}
	defer resOutAll.Body.Close()

	resDropTable, err := client.Do(dropTable)
	if err != nil {
		t.Fatal(err)
	}
	resDropTable.Body.Close()
}
