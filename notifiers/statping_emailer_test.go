package notifiers

import (
	"github.com/statping-ng/statping-ng/database"
	"github.com/statping-ng/statping-ng/types/core"
	"github.com/statping-ng/statping-ng/types/failures"
	"github.com/statping-ng/statping-ng/types/notifications"
	"github.com/statping-ng/statping-ng/types/null"
	"github.com/statping-ng/statping-ng/types/services"
	"github.com/statping-ng/statping-ng/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

var (
	testEmail string
)

func TestStatpingEmailerNotifier(t *testing.T) {
	err := utils.InitLogs()
	require.Nil(t, err)

	t.Parallel()
	db, err := database.OpenTester()
	require.Nil(t, err)
	db.AutoMigrate(&notifications.Notification{})
	notifications.SetDB(db)
	core.Example()

	testEmail = utils.Params.GetString("TEST_EMAIL")
	statpingMailer.Host = null.NewNullString(testEmail)
	statpingMailer.Enabled = null.NewNullBool(true)

	if testEmail == "" {
		t.Log("statping email notifier testing skipped, missing TEST_EMAIL environment variable")
		t.SkipNow()
	}

	t.Run("Load statping emailer", func(t *testing.T) {
		statpingMailer.Host = null.NewNullString(testEmail)
		statpingMailer.Delay = time.Duration(100 * time.Millisecond)
		statpingMailer.Limits = 3
		Add(statpingMailer)
		assert.Equal(t, "Hunter Long", statpingMailer.Author)
		assert.Equal(t, testEmail, statpingMailer.Host.String)
	})

	t.Run("statping emailer Within Limits", func(t *testing.T) {
		ok := statpingMailer.CanSend()
		assert.True(t, ok)
	})

	t.Run("statping emailer OnSave", func(t *testing.T) {
		_, err := statpingMailer.OnSave()
		assert.Nil(t, err)
	})

	t.Run("statping emailer OnFailure", func(t *testing.T) {
		_, err := statpingMailer.OnFailure(services.Example(false), failures.Example())
		assert.Nil(t, err)
	})

	t.Run("statping emailer OnSuccess", func(t *testing.T) {
		_, err := statpingMailer.OnSuccess(services.Example(true))
		assert.Nil(t, err)
	})

}
