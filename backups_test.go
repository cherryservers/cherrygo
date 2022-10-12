package cherrygo

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"testing"
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
		fmt.Fprint(writer, `{
			"id": 123,
			"status": "deployed",
			"state": "active",
			"private_ip": "10.10.10.10",
			"public_ip": "5.199.199.199",
			"size_gigabytes": 100,
			"used_gigabytes": 1
		}`)
	})

	backup, _, err := client.Backups.Get(123, nil)
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
		fmt.Fprint(writer, `[
			{"id": 123},
			{"id": 321}
		]`)
	})

	backups, _, err := client.Backups.ListBackups(projectID, nil)

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

	serverID := 312
	expected := []BackupStoragePlan{
		{
			ID:            123,
			Name:          "Backup 100",
			Slug:          "backup_100",
			SizeGigabytes: 100,
		},
		{ID: 321},
	}

	mux.HandleFunc("/v1/servers/"+strconv.Itoa(serverID)+"/backup-storage-plans", func(writer http.ResponseWriter, request *http.Request) {
		testMethod(t, request, http.MethodGet)
		fmt.Fprint(writer, `[
			{
				"id": 123,
				"name": "Backup 100",
				"slug": "backup_100",
				"size_gigabytes": 100
			},
			{"id": 321}
		]`)
	})

	backupPlans, _, err := client.Backups.ListPlans(serverID, nil)

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
	requestBody := map[string]interface{}{
		"server_id": float64(serverID),
		"slug":      "backup_100",
		"region":    "eu_nord_1",
		"ssh_key":   "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC6ec8eT...",
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

		fmt.Fprint(writer, `{"id": 123}`)
	})

	createBackup := CreateBackup{
		ServerID:       serverID,
		BackupPlanSlug: "backup_100",
		RegionSlug:     "eu_nord_1",
		SSHKey:         "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC6ec8eT...",
	}

	_, _, err := client.Backups.Create(&createBackup)
	if err != nil {
		t.Errorf("Backup.Create returned %+v", err)
	}
}

func TestBackupStorage_Update(t *testing.T) {
	setup()
	defer teardown()

	requestBody := map[string]interface{}{
		"id":       float64(123),
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

		fmt.Fprint(writer, `{"id": 123}`)
	})

	updateBackupStorage := UpdateBackupStorage{
		BackupStorageID: 123,
		BackupPlanSlug:  "backup_500",
		Password:        "abc123",
	}

	_, _, err := client.Backups.Update(&updateBackupStorage)
	if err != nil {
		t.Errorf("Backups.Update returned %+v", err)
	}
}

func TestBackupStorage_UpdateService(t *testing.T) {
	setup()
	defer teardown()

	serviceName := "FTP"
	requestBody := map[string]interface{}{
		"id":        float64(123),
		"name":      serviceName,
		"enabled":   true,
		"whitelist": []interface{}{"1.1.1.1", "2.2.2.2"},
	}

	mux.HandleFunc("/v1/backup-storages/123/services/"+serviceName, func(writer http.ResponseWriter, request *http.Request) {
		testMethod(t, request, http.MethodPatch)

		var v map[string]interface{}
		err := json.NewDecoder(request.Body).Decode(&v)
		if err != nil {
			t.Fatalf("decode json: %v", err)
		}

		if !reflect.DeepEqual(v, requestBody) {
			t.Errorf("Request body\n sent %#v\n expected %#v", v, requestBody)
		}

		fmt.Fprint(writer, `[{
			"name": "FTP",
			"username": "username",
			"password": "password",
			"host": "host",
			"ssh_key": "ssh_key",
			"enabled": true
			}]`)
	})

	updateBackupService := UpdateBackupService{
		BackupStorageID:   123,
		BackupServiceName: serviceName,
		Enabled:           true,
		Whitelist:         []string{"1.1.1.1", "2.2.2.2"},
	}

	_, _, err := client.Backups.UpdateBackupService(&updateBackupService)
	if err != nil {
		t.Errorf("Backups.UpdateBackupService returned %+v", err)
	}
}

func TestBackupStorage_Delete(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v1/backup-storages/123", func(writer http.ResponseWriter, request *http.Request) {
		testMethod(t, request, http.MethodDelete)
		writer.WriteHeader(http.StatusNoContent)
		fmt.Fprint(writer)
	})

	_, err := client.Backups.Delete(123)
	if err != nil {
		t.Errorf("Backups.Delete returned %+v", err)
	}
}
