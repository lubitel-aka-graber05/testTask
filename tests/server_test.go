package tests

import (
	"bufio"
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os/exec"
	"strings"
	"testing"
	"time"

	"testTask/internal/config"
	"testTask/internal/db/postgres/auth"
	"testTask/internal/db/postgres/notes"
	addnote "testTask/internal/handlers/addNote"
	"testTask/internal/handlers/adduser"
	outputnote "testTask/internal/handlers/outputNote"
	"testTask/internal/logger"

	"github.com/stretchr/testify/assert"
)

const (
	testUserAuthName = "testuser"
	testUserAuthPass = "secret"
	testUserNoteName = "ashot3000"
)

func TestServer(t *testing.T) {
	_ = exec.Command("./start_test_db.sh").Run()
	defer exec.Command("./stop_test_db.sh").Run()

	time.Sleep(time.Second * 2)

	log := logger.SetupLogger()

	filePath := "config.yaml"

	cfg, err := config.CreateConfig(log, filePath)
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
	router.HandleFunc("/adduser", adduser.AddUserHandler(log, authDB))
	router.HandleFunc("/addnote", addnote.AddNoteHandler(log, authDB, notesDB))
	router.HandleFunc("/outbyname", outputnote.OutputByNameHandler(log, authDB, notesDB))

	server := httptest.NewServer(router)

	client := server.Client()

	addUserBody := strings.NewReader(fmt.Sprintf(`{"username":"%s","pass":"%s"}`, testUserAuthName, testUserAuthPass))

	reqAddUser, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/adduser", server.URL), addUserBody)
	if err != nil {
		t.Fatal(err)
	}

	resAddUser, err := client.Do(reqAddUser)
	if err != nil {
		t.Fatal(err)
	}
	// defer resAddUser.Body.Close()

	assert.Equal(t, resAddUser.StatusCode, http.StatusOK, "status code not correct")
	r := bufio.NewReader(resAddUser.Body)
	buf := &bytes.Buffer{}
	_, err = r.WriteTo(buf)
	assert.Nil(t, err)
	assert.Equal(t, buf.String(), fmt.Sprintf(`{"error":"","status":%d}`+"\n", resAddUser.StatusCode), resAddUser.StatusCode)
}
