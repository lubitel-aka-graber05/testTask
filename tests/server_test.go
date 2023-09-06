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
	"github.com/stretchr/testify/require"
)

const (
	testUserAuthName    = "testuser"
	testUserAuthPass    = "secret"
	testUserNoteName    = "ashot3000"
	testUserNoteBody    = "Здраствуйте, миня завут Ашот"
	testUserNoteCorrect = "Здравствуйте, меня зовут Ашот"
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
		require.Nil(t, err)
	}

	resAddUser, err := client.Do(reqAddUser)
	if err != nil {
		require.Nil(t, err, "request with auth info return error")
	}
	defer resAddUser.Body.Close()

	assert.Equal(t, resAddUser.StatusCode, http.StatusOK, "status code not correct")
	r := bufio.NewReader(resAddUser.Body)
	buf := &bytes.Buffer{}
	_, err = r.WriteTo(buf)
	assert.Nil(t, err)
	assert.Equal(t, buf.String(), fmt.Sprintf(`{"statuscode":%d}`+"\n", http.StatusOK))

	buf.Reset()

	addNoteBody := strings.NewReader(fmt.Sprintf(`{"username":"%s","body":"%s"}`, testUserNoteName, testUserNoteBody))

	reqAddNote, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/addnote", server.URL), addNoteBody)
	if err != nil {
		require.Nil(t, err)
	}
	reqAddNote.SetBasicAuth(testUserAuthName, testUserAuthPass)

	resAddNote, err := client.Do(reqAddNote)
	if err != nil {
		require.Nil(t, err)
	}
	defer resAddNote.Body.Close()

	require.Equal(t, resAddNote.StatusCode, http.StatusOK)

	_, err = buf.ReadFrom(resAddNote.Body)
	require.Nil(t, err)
	require.Equal(t, buf.String(), fmt.Sprintf(`{"statuscode":%d}`+"\n", http.StatusOK))

	buf.Reset()

	outputBody := strings.NewReader(fmt.Sprintf(`{"username":"%s"}`, testUserNoteName))

	reqOutNote, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/outbyname", server.URL), outputBody)
	require.Nil(t, err)
	reqOutNote.SetBasicAuth(testUserAuthName, testUserAuthPass)

	resOutNode, err := client.Do(reqOutNote)
	require.Nil(t, err)
	defer resOutNode.Body.Close()

	require.Equal(t, resOutNode.StatusCode, http.StatusOK)
	_, err = buf.ReadFrom(resOutNode.Body)
	require.Nil(t, err)

	require.Equal(t, buf.String(), fmt.Sprintf(`{"username":"%s","body":["%s"],"statuscode":%d}`+"\n", testUserNoteName, testUserNoteCorrect, http.StatusOK))
}
