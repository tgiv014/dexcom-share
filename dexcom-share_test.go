package dexcomshare

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/dnaeon/go-vcr.v3/cassette"
	"gopkg.in/dnaeon/go-vcr.v3/recorder"
)

func beforeSaveHook(i *cassette.Interaction) error {
	// Remove account ID, name, and password from the request body
	tmp := map[string]interface{}{}

	err := json.Unmarshal([]byte(i.Request.Body), &tmp)
	if err == nil {
		if _, ok := tmp["accountId"]; ok {
			tmp["accountId"] = "REDACTED"
		}

		if _, ok := tmp["accountName"]; ok {
			tmp["accountName"] = "REDACTED"
		}

		if _, ok := tmp["password"]; ok {
			tmp["password"] = "REDACTED"
		}

		if _, ok := tmp["sessionId"]; ok {
			tmp["sessionId"] = "REDACTED"
		}

		b, err := json.Marshal(tmp)
		if err != nil {
			return err
		}

		i.Request.Body = string(b)
	}

	if i.Request.URL == "https://share2.dexcom.com/ShareWebServices/Services/General/AuthenticatePublisherAccount" {
		i.Response.Body = `"accountID"`
	}

	if i.Request.URL == "https://share2.dexcom.com/ShareWebServices/Services/General/LoginPublisherAccountByName" {
		i.Response.Body = `"sessionID"`
	}

	return nil
}

func Test_Client(t *testing.T) {
	r, err := recorder.New("testdata/Test_NewClient")
	assert.NoError(t, err)
	defer r.Stop()

	r.AddHook(beforeSaveHook, recorder.BeforeSaveHook)

	client := r.GetDefaultClient()

	c, err := NewClient("username", "password", WithClient(client))
	assert.NoError(t, err)
	assert.NotNil(t, c)

	entries, err := c.ReadGlucose(1440, 1)
	assert.NoError(t, err)
	assert.Len(t, entries, 1)

	entries, err = c.ReadGlucose(1440, 100)
	assert.NoError(t, err)
	assert.Len(t, entries, 100)
}
