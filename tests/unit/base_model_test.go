package unit

import (
	"os"
	"testing"
	"webgos/internal/models"
	"webgos/internal/xlog"

	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	xlog.InitLogger()
	code := m.Run()
	xlog.Xlogger.Close()
	os.Exit(code)
}

func TestUserModel(t *testing.T) {
	t.Run("TestUserModelStructure", func(t *testing.T) {
		user := &models.User{}
		assert.NotNil(t, user)
		assert.Equal(t, 0, user.ID)
	})
}

func TestTT(t *testing.T) {
	t.Run("TestTT", func(t *testing.T) {
		xlog.Debug("0: %v", ss(0))
		xlog.Debug("1: %v", ss(1))
		xlog.Debug("2: %v", ss(2))
		xlog.Debug("3: %v", ss(3))
		xlog.Debug("4: %v", ss(4))
	})
}

func ss(num int) float64 {
	const rate = 0.1
	const drate = 0.07
	const invest = 20

	totalInvest := float64((num + 1) * invest)

	switch num {
	case 0:
		currentGain := float64(invest)
		return currentGain*(1-drate) - totalInvest
	case 1:
		gain := invest + invest*rate + invest
		return gain*(1-drate) - totalInvest
	case 2:
		prevGain := invest + invest*rate + invest
		currentGain := prevGain + prevGain*rate + invest
		return currentGain*(1-drate) - totalInvest
	case 3:
		prevGain1 := invest + invest*rate + invest
		prevGain2 := prevGain1 + prevGain1*rate + invest
		currentGain := prevGain2 + prevGain2*rate + invest
		return currentGain*(1-drate) - totalInvest
	case 4:
		prevGain1 := invest + invest*rate + invest
		prevGain2 := prevGain1 + prevGain1*rate + invest
		prevGain3 := prevGain2 + prevGain2*rate + invest
		currentGain := prevGain3 + prevGain3*rate + invest
		return currentGain*(1-drate) - totalInvest
	default:
		return 0
	}
}
