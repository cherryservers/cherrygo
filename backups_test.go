package cherrygo

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBackupStorage_Get(t *testing.T) {
	setup()
	defer teardown()

	expected := BackupStorage{
		ID:            123,
		Status:        "deployed",
		State:         "active",
		PrivateIP:     "10.10.10.10",
		PublicIP:      "5.199.199.199",
		SizeGigabytes: 100,
		UsedGigabytes: 1,
	}

	mux.HandleFunc("/v1/backup-storages/123", func(writer http.ResponseWriter, request *http.Request) {
		testMethod(t, request, http.MethodGet)
		_, err := fmt.Fprint(writer, `{
			"id": 123,
			"status": "deployed",
			"state": "active",
			"private_ip": "10.10.10.10",
			"public_ip": "5.199.199.199",
			"size_gigabytes": 100,
			"used_gigabytes": 1
		}`)
		require.NoError(t, err)
	})

	backup, _, err := testClient.Backups.Get(t.Context(), 123, nil)
	if err != nil {
		t.Errorf("Backups.Get returned %+v", err)
	}

	if !reflect.DeepEqual(backup, expected) {
		t.Errorf("Backups.Get returned %+v, expected %+v", backup, expected)
	}
}

func TestBackupStorage_ListBackups(t *testing.T) {
	setup()
	defer teardown()

	expected := []BackupStorage{
		{ID: 123},
		{ID: 321},
	}

	mux.HandleFunc("/v1/projects/"+strconv.Itoa(projectID)+"/backup-storages", func(writer http.ResponseWriter, request *http.Request) {
		testMethod(t, request, http.MethodGet)
		_, err := fmt.Fprint(writer, `[
			{"id": 123},
			{"id": 321}
		]`)
		require.NoError(t, err)
	})

	backups, _, err := testClient.Backups.ListBackups(t.Context(), projectID, nil)
	if err != nil {
		t.Errorf("Backups.ListBackups returned %+v", err)
	}

	if !reflect.DeepEqual(backups, expected) {
		t.Errorf("Backups.ListBackups returned %+v, expected %+v", backups, expected)
	}
}

func TestBackupStorage_ListPlans(t *testing.T) {
	setup()
	defer teardown()

	expected := []BackupStoragePlan{
		{
			ID:            123,
			Name:          "Backup 100",
			Slug:          "backup_100",
			SizeGigabytes: 100,
		},
		{ID: 321},
	}

	mux.HandleFunc("/v1/backup-storage-plans", func(writer http.ResponseWriter, request *http.Request) {
		testMethod(t, request, http.MethodGet)
		_, err := fmt.Fprint(writer, `[
			{
				"id": 123,
				"name": "Backup 100",
				"slug": "backup_100",
				"size_gigabytes": 100
			},
			{"id": 321}
		]`)

		require.NoError(t, err)
	})

	backupPlans, _, err := testClient.Backups.ListPlans(t.Context(), nil)
	if err != nil {
		t.Errorf("Backups.ListPlans returned %+v", err)
	}

	if !reflect.DeepEqual(backupPlans, expected) {
		t.Errorf("Backups.ListPlans returned %+v, expected %+v", backupPlans, expected)
	}
}

func TestBackupStorage_Create(t *testing.T) {
	setup()
	defer teardown()

	serverID := 312
	requestBody := map[string]any{
		"slug":    "backup_100",
		"region":  "eu_nord_1",
		"ssh_key": "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC6ec8eT...",
	}

	mux.HandleFunc("/v1/servers/"+strconv.Itoa(serverID)+"/backup-storages", func(writer http.ResponseWriter, request *http.Request) {
		testMethod(t, request, http.MethodPost)

		var v map[string]interface{}
		err := json.NewDecoder(request.Body).Decode(&v)
		if err != nil {
			t.Fatalf("decode json: %v", err)
		}

		if !reflect.DeepEqual(v, requestBody) {
			t.Errorf("Request body\n sent %#v\n expected %#v", v, requestBody)
		}

		_, err = fmt.Fprint(writer, `{"id": 123}`)
		require.NoError(t, err)
	})

	createBackup := CreateBackup{
		BackupPlanSlug: "backup_100",
		RegionSlug:     "eu_nord_1",
		SSHKey:         "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC6ec8eT...",
	}

	_, _, err := testClient.Backups.Create(t.Context(), serverID, &createBackup)
	if err != nil {
		t.Errorf("Backup.Create returned %+v", err)
	}
}

func TestBackupStorage_Update(t *testing.T) {
	setup()
	defer teardown()

	requestBody := map[string]any{
		"slug":     "backup_500",
		"password": "abc123",
	}

	mux.HandleFunc("/v1/backup-storages/123", func(writer http.ResponseWriter, request *http.Request) {
		testMethod(t, request, http.MethodPut)

		var v map[string]interface{}
		err := json.NewDecoder(request.Body).Decode(&v)
		if err != nil {
			t.Fatalf("decode json: %v", err)
		}

		if !reflect.DeepEqual(v, requestBody) {
			t.Errorf("Request body\n sent %#v\n expected %#v", v, requestBody)
		}

		_, err = fmt.Fprint(writer, `{"id": 123}`)
		require.NoError(t, err)
	})

	updateBackupStorage := UpdateBackupStorage{
		BackupPlanSlug: "backup_500",
		Password:       "abc123",
	}

	_, _, err := testClient.Backups.Update(t.Context(), 123, &updateBackupStorage)
	if err != nil {
		t.Errorf("Backups.Update returned %+v", err)
	}
}

func TestBackupStorage_UpdateMethod(t *testing.T) {
	setup()
	defer teardown()

	methodName := "FTP"
	requestBody := map[string]any{
		"enabled":   true,
		"whitelist": []any{"1.1.1.1", "2.2.2.2"},
	}

	mux.HandleFunc("/v1/backup-storages/123/methods/"+methodName, func(writer http.ResponseWriter, request *http.Request) {
		testMethod(t, request, http.MethodPatch)

		var v map[string]interface{}
		err := json.NewDecoder(request.Body).Decode(&v)
		if err != nil {
			t.Fatalf("decode json: %v", err)
		}

		if !reflect.DeepEqual(v, requestBody) {
			t.Errorf("Request body\n sent %#v\n expected %#v", v, requestBody)
		}

		_, err = fmt.Fprint(writer, `[{
			"name": "FTP",
			"username": "username",
			"password": "password",
			"host": "host",
			"ssh_key": "ssh_key",
			"enabled": true
			}]`)

		require.NoError(t, err)
	})

	var (
		enabled   = true
		whitelist = []string{"1.1.1.1", "2.2.2.2"}
	)
	updateBackupMethod := UpdateBackupMethod{
		Enabled:   &enabled,
		Whitelist: &whitelist,
	}

	_, _, err := testClient.Backups.UpdateBackupMethod(t.Context(), 123, methodName, &updateBackupMethod)
	if err != nil {
		t.Errorf("Backups.UpdateBackupMethod returned %+v", err)
	}
}

func TestBackupStorage_UpdateMethodBodyFieldsOmitted(t *testing.T) {
	setup()
	defer teardown()

	bod := UpdateBackupMethod{}

	mux.HandleFunc("PATCH /v1/backup-storages/123/methods/FTP", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPatch)

		got, err := io.ReadAll(r.Body)
		require.NoError(t, err)
		assert.JSONEq(t, "{}", string(got))

		_, err = fmt.Fprint(w, `[{
			"name": "FTP"
			}]`)
		require.NoError(t, err)
	})

	_, _, err := testClient.Backups.UpdateBackupMethod(t.Context(), 123, "FTP", &bod)
	require.NoError(t, err)
}

func TestBackupStorage_Delete(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v1/backup-storages/123", func(writer http.ResponseWriter, request *http.Request) {
		testMethod(t, request, http.MethodDelete)
		writer.WriteHeader(http.StatusNoContent)
		_, err := fmt.Fprint(writer)
		require.NoError(t, err)
	})

	_, err := testClient.Backups.Delete(t.Context(), 123)
	if err != nil {
		t.Errorf("Backups.Delete returned %+v", err)
	}
}
