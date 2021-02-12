package mod

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/jinzhu/gorm"
	"github.com/matryer/is"
	"github.com/rs/zerolog/log"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"

	"github.com/domino14/liwords/pkg/entity"
	"github.com/domino14/liwords/pkg/stores/user"
	pkguser "github.com/domino14/liwords/pkg/user"
	ms "github.com/domino14/liwords/rpc/api/proto/mod_service"
	macondoconfig "github.com/domino14/macondo/config"
)

var TestDBHost = os.Getenv("TEST_DB_HOST")
var TestingDBConnStr = "host=" + TestDBHost + " port=5432 user=postgres password=pass sslmode=disable"

var DefaultConfig = macondoconfig.Config{
	LexiconPath:               os.Getenv("LEXICON_PATH"),
	LetterDistributionPath:    os.Getenv("LETTER_DISTRIBUTION_PATH"),
	DefaultLexicon:            "CSW19",
	DefaultLetterDistribution: "English",
}

func recreateDB() {
	// Create a database.
	db, err := gorm.Open("postgres", TestingDBConnStr+" dbname=postgres")
	if err != nil {
		log.Fatal().Err(err).Msg("error")
	}
	defer db.Close()
	db = db.Exec("DROP DATABASE IF EXISTS liwords_test")
	if db.Error != nil {
		log.Fatal().Err(db.Error).Msg("error")
	}
	db = db.Exec("CREATE DATABASE liwords_test")
	if db.Error != nil {
		log.Fatal().Err(db.Error).Msg("error")
	}

	ustore := userStore(TestingDBConnStr + " dbname=liwords_test")

	for _, u := range []*entity.User{
		{Username: "Spammer", Email: "spammer@woogles.io", UUID: "Spammer"},
		{Username: "Sandbagger", Email: "sandbagger@gmail.com", UUID: "Sandbagger"},
		{Username: "Cheater", Email: "cheater@woogles.io", UUID: "Cheater"},
	} {
		err = ustore.New(context.Background(), u)
		if err != nil {
			log.Fatal().Err(err).Msg("error")
		}
	}
	ustore.(*user.DBStore).Disconnect()
}

func userStore(dbURL string) pkguser.Store {
	ustore, err := user.NewDBStore(TestingDBConnStr + " dbname=liwords_test")
	if err != nil {
		log.Fatal().Err(err).Msg("error")
	}
	return ustore
}

func TestMod(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	cstr := TestingDBConnStr + " dbname=liwords_test"
	recreateDB()
	us := userStore(cstr)

	var muteDuration int32 = 2

	muteAction := &ms.ModAction{UserId: "Spammer", Type: ms.ModActionType_MUTE, Duration: muteDuration}
	// Negative value for duration should not matter for transient actions
	resetAction := &ms.ModAction{UserId: "Sandbagger", Type: ms.ModActionType_RESET_STATS_AND_RATINGS, Duration: -10}
	suspendAction := &ms.ModAction{UserId: "Cheater", Type: ms.ModActionType_SUSPEND_ACCOUNT, Duration: 100}

	// Remove an action that does not exist
	err := RemoveActions(ctx, us, []*ms.ModAction{muteAction})
	errString := fmt.Sprintf("user does not have current action %s", muteAction.Type.String())
	is.True(err.Error() == errString)

	// Apply Actions
	err = ApplyActions(ctx, us, []*ms.ModAction{muteAction, resetAction, suspendAction})
	is.NoErr(err)

	is.True(ActionExists(ctx, us, "Spammer", muteAction.Type).Error() == "this user is not permitted to perform this action")
	is.NoErr(ActionExists(ctx, us, "Sandbagger", resetAction.Type))
	is.True(ActionExists(ctx, us, "Cheater", suspendAction.Type).Error() == "this user is not permitted to perform this action")

	// Check Actions
	expectedSpammerActions, err := GetActions(ctx, us, "Spammer")
	is.NoErr(err)
	is.NoErr(equalActionMaps(expectedSpammerActions, makeActionMap([]*ms.ModAction{muteAction})))
	is.True(expectedSpammerActions[muteAction.Type.String()].EndTime != nil)
	is.True(expectedSpammerActions[muteAction.Type.String()].StartTime != nil)

	expectedSpammerHistory, err := GetActionHistory(ctx, us, "Spammer")
	is.NoErr(err)
	is.NoErr(equalActionHistories(expectedSpammerHistory, []*ms.ModAction{}))

	expectedSandbaggerActions, err := GetActions(ctx, us, "Sandbagger")
	is.NoErr(err)
	is.NoErr(equalActionMaps(expectedSandbaggerActions, makeActionMap([]*ms.ModAction{})))

	expectedSandbaggerHistory, err := GetActionHistory(ctx, us, "Sandbagger")
	is.NoErr(err)
	is.NoErr(equalActionHistories(expectedSandbaggerHistory, []*ms.ModAction{resetAction}))
	is.True(expectedSandbaggerHistory[0] != nil)
	is.True(expectedSandbaggerHistory[0].EndTime != nil)
	is.True(expectedSandbaggerHistory[0].StartTime != nil)
	is.True(expectedSandbaggerHistory[0].Expired)
	is.NoErr(equalTimes(expectedSandbaggerHistory[0].EndTime, expectedSandbaggerHistory[0].StartTime))
	is.NoErr(equalTimes(expectedSandbaggerHistory[0].EndTime, expectedSandbaggerHistory[0].RemovedTime))
	is.NoErr(equalTimes(expectedSandbaggerHistory[0].StartTime, expectedSandbaggerHistory[0].EndTime))

	expectedCheaterActions, err := GetActions(ctx, us, "Cheater")
	is.NoErr(err)
	is.NoErr(equalActionMaps(expectedCheaterActions, makeActionMap([]*ms.ModAction{suspendAction})))
	is.True(expectedCheaterActions[suspendAction.Type.String()].EndTime != nil)
	is.True(expectedCheaterActions[suspendAction.Type.String()].StartTime != nil)

	expectedCheaterHistory, err := GetActionHistory(ctx, us, "Cheater")
	is.NoErr(err)
	is.NoErr(equalActionHistories(expectedCheaterHistory, []*ms.ModAction{}))

	longerSuspendAction := &ms.ModAction{UserId: "Cheater", Type: ms.ModActionType_SUSPEND_ACCOUNT, Duration: 200}

	// Overwrite some actions
	err = ApplyActions(ctx, us, []*ms.ModAction{longerSuspendAction})
	is.NoErr(err)

	expectedCheaterActions, err = GetActions(ctx, us, "Cheater")
	is.NoErr(err)
	is.NoErr(equalActionMaps(expectedCheaterActions, makeActionMap([]*ms.ModAction{longerSuspendAction})))
	is.True(expectedCheaterActions[suspendAction.Type.String()].EndTime != nil)
	is.True(expectedCheaterActions[suspendAction.Type.String()].StartTime != nil)
	is.True(expectedCheaterActions[suspendAction.Type.String()].Duration == 200)

	expectedCheaterHistory, err = GetActionHistory(ctx, us, "Cheater")
	is.NoErr(err)
	is.NoErr(equalActionHistories(expectedCheaterHistory, []*ms.ModAction{suspendAction}))
	is.True(!expectedCheaterHistory[0].Expired)
	is.NoErr(equalTimes(expectedCheaterHistory[0].EndTime, expectedCheaterHistory[0].StartTime))
	is.NoErr(equalTimes(expectedCheaterHistory[0].EndTime, expectedCheaterHistory[0].RemovedTime))

	// Recheck Spammer actions
	is.True(ActionExists(ctx, us, "Spammer", muteAction.Type).Error() == "this user is not permitted to perform this action")

	expectedSpammerActions, err = GetActions(ctx, us, "Spammer")
	is.NoErr(err)
	is.NoErr(equalActionMaps(expectedSpammerActions, makeActionMap([]*ms.ModAction{muteAction})))
	is.True(expectedSpammerActions[muteAction.Type.String()].EndTime != nil)
	is.True(expectedSpammerActions[muteAction.Type.String()].StartTime != nil)

	expectedSpammerHistory, err = GetActionHistory(ctx, us, "Spammer")
	is.NoErr(err)
	is.NoErr(equalActionHistories(expectedSpammerHistory, []*ms.ModAction{}))

	// Wait
	time.Sleep(time.Duration(muteDuration+1) * time.Second)

	// Recheck Spammer actions
	is.NoErr(ActionExists(ctx, us, "Spammer", muteAction.Type))
	expectedSpammerActions, err = GetActions(ctx, us, "Spammer")
	is.NoErr(err)
	is.NoErr(equalActionMaps(expectedSpammerActions, makeActionMap([]*ms.ModAction{})))

	expectedSpammerHistory, err = GetActionHistory(ctx, us, "Spammer")
	is.NoErr(err)
	is.NoErr(equalActionHistories(expectedSpammerHistory, []*ms.ModAction{muteAction}))
	is.True(expectedSpammerHistory[0].EndTime != nil)
	is.True(expectedSpammerHistory[0].StartTime != nil)
	is.True(expectedSpammerHistory[0].Expired)
	is.NoErr(equalTimes(expectedSpammerHistory[0].EndTime, expectedSpammerHistory[0].StartTime))
	is.NoErr(equalTimes(expectedSpammerHistory[0].EndTime, expectedSpammerHistory[0].RemovedTime))
	// Test negative durations
	invalidSuspendAction := &ms.ModAction{UserId: "Cheater", Type: ms.ModActionType_SUSPEND_ACCOUNT, Duration: -100}

	err = ApplyActions(ctx, us, []*ms.ModAction{invalidSuspendAction})
	is.True(err.Error() == "nontransient moderator action has a negative duration: -100")

	// Apply a permanent action

	permanentSuspendAction := &ms.ModAction{UserId: "Sandbagger", Type: ms.ModActionType_SUSPEND_ACCOUNT, Duration: 0}

	err = ApplyActions(ctx, us, []*ms.ModAction{permanentSuspendAction})
	is.NoErr(err)

	is.True(ActionExists(ctx, us, "Sandbagger", permanentSuspendAction.Type).Error() == "this user is not permitted to perform this action")

	expectedSandbaggerActions, err = GetActions(ctx, us, "Sandbagger")
	is.NoErr(err)
	is.NoErr(equalActionMaps(expectedSandbaggerActions, makeActionMap([]*ms.ModAction{permanentSuspendAction})))
	is.True(expectedSandbaggerActions[permanentSuspendAction.Type.String()].EndTime == nil)
	is.True(expectedSandbaggerActions[permanentSuspendAction.Type.String()].StartTime != nil)

	expectedSandbaggerHistory, err = GetActionHistory(ctx, us, "Sandbagger")
	is.NoErr(err)
	is.NoErr(equalActionHistories(expectedSandbaggerHistory, []*ms.ModAction{resetAction}))

	// Remove an action
	err = RemoveActions(ctx, us, []*ms.ModAction{permanentSuspendAction})
	is.NoErr(err)

	is.NoErr(ActionExists(ctx, us, "Sandbagger", permanentSuspendAction.Type))

	expectedSandbaggerActions, err = GetActions(ctx, us, "Sandbagger")
	is.NoErr(err)
	is.NoErr(equalActionMaps(expectedSandbaggerActions, makeActionMap([]*ms.ModAction{})))

	expectedSandbaggerHistory, err = GetActionHistory(ctx, us, "Sandbagger")
	is.NoErr(err)
	is.NoErr(equalActionHistories(expectedSandbaggerHistory, []*ms.ModAction{resetAction, permanentSuspendAction}))
	is.True(expectedSandbaggerHistory[0].Expired)
	is.NoErr(equalTimes(expectedSandbaggerHistory[0].EndTime, expectedSandbaggerHistory[0].RemovedTime))
	is.NoErr(equalTimes(expectedSandbaggerHistory[0].StartTime, expectedSandbaggerHistory[0].EndTime))
}

func equalActionHistories(ah1 []*ms.ModAction, ah2 []*ms.ModAction) error {
	if len(ah1) != len(ah2) {
		return errors.New("history lengths are not the same")
	}
	for i := 0; i < len(ah1); i++ {
		a1 := ah1[i]
		a2 := ah2[i]
		if !equalActions(a1, a2) {
			return fmt.Errorf("actions are not equal:\n  a1.UserId: %s a1.Type: %s, a1.Duration: %d\n"+
				"  a1.UserId: %s a1.Type: %s, a1.Duration: %d\n", a1.UserId, a1.Type, a1.Duration,
				a2.UserId, a2.Type, a2.Duration)
		}
	}
	return nil
}

func equalActionMaps(am1 map[string]*ms.ModAction, am2 map[string]*ms.ModAction) error {
	for key, _ := range ms.ModActionType_value {
		a1 := am1[key]
		a2 := am2[key]
		if a1 == nil && a2 == nil {
			continue
		}
		if a1 == nil || a2 == nil {
			return fmt.Errorf("exactly one actions is nil: %s", key)
		}
		if !equalActions(a1, a2) {
			return fmt.Errorf("actions are not equal:\n  a1.UserId: %s a1.Type: %s, a1.Duration: %d\n"+
				"  a1.UserId: %s a1.Type: %s, a1.Duration: %d\n", a1.UserId, a1.Type, a1.Duration,
				a2.UserId, a2.Type, a2.Duration)
		}
	}
	return nil
}

func equalActions(a1 *ms.ModAction, a2 *ms.ModAction) bool {
	return a1.UserId == a2.UserId &&
		a1.Type == a2.Type &&
		a1.Duration == a2.Duration
}

func equalTimes(t1 *timestamppb.Timestamp, t2 *timestamppb.Timestamp) error {
	gt1, err := ptypes.Timestamp(t1)
	if err != nil {
		return err
	}
	gt2, err := ptypes.Timestamp(t1)
	if err != nil {
		return err
	}
	if !gt1.Equal(gt2) {
		return fmt.Errorf("times are not equal:\n%v\n%v", gt1, gt2)
	}
	return nil
}

func makeActionMap(actions []*ms.ModAction) map[string]*ms.ModAction {
	actionMap := make(map[string]*ms.ModAction)
	for _, action := range actions {
		actionMap[action.Type.String()] = action
	}
	return actionMap
}
