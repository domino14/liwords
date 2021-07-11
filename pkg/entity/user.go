package entity

import (
	"errors"
	"math"
	"strconv"
	"sync"
	"time"

	"github.com/domino14/liwords/pkg/glicko"
	ms "github.com/domino14/liwords/rpc/api/proto/mod_service"
	"github.com/domino14/liwords/rpc/api/proto/realtime"
	"github.com/rs/zerolog/log"
)

const (
	// SessionExpiration - Expire a session after this much time.
	SessionExpiration = time.Hour * 24 * 30
)

type Actions struct {
	Current map[string]*ms.ModAction
	History []*ms.ModAction
}

// User - the db-specific details are in the store package.
type User struct {
	sync.RWMutex

	Anonymous bool
	// ID is the database ID. Since this increases monotonically, we should
	// not expose it to the user
	ID uint
	// UUID is the "user-exposed" ID, in any APIs.
	UUID     string
	Username string
	Password string
	Email    string
	Profile  *Profile
	// CurrentChannel tracks presence; where is the user currently?
	CurrentChannel string
	IsBot          bool
	IsDirector     bool
	IsMod          bool
	IsAdmin        bool

	Actions   *Actions
	Notoriety int
}

type UserPermission int

const (
	PermDirector UserPermission = iota
	PermMod
	PermAdmin
	PermBot
)

// Session - The db specific-details are in the store package.
type Session struct {
	ID       string
	Username string
	UserUUID string
	Expiry   time.Time
}

// Profile is a user profile. It might not be defined for anonymous users.
type Profile struct {
	FirstName string
	LastName  string
	// BirthDate uses ISO format YYYY-MM-DD
	BirthDate   string
	CountryCode string
	Title       string
	About       string
	Ratings     Ratings
	Stats       ProfileStats
	AvatarUrl   string
}

// If the RD is <= this number, the rating is "known"
const RatingDeviationConfidence = float64(glicko.MinimumRatingDeviation + 30)

// RelevantRating returns the rating from a Ratings object given a rating key.
func RelevantRating(ratings Ratings, ratingKey VariantKey) string {

	unknownRating := "?"

	if ratings.Data == nil {
		// This is not an unrated user. Use default rating.
		return strconv.Itoa(glicko.InitialRating) + unknownRating
	}
	ratdict, ok := ratings.Data[ratingKey]
	if ok {
		if ratdict.RatingDeviation <= RatingDeviationConfidence {
			unknownRating = ""
		}
		return strconv.Itoa(int(math.Round(ratdict.Rating))) + unknownRating
	}
	// User has no rating in this particular variant.
	return strconv.Itoa(glicko.InitialRating) + unknownRating
}

// GetRelevantRating gets a displayable rating for this user, based on the passed-in
// rating key (encoding variant, time control, etc)
func (u *User) GetRelevantRating(ratingKey VariantKey) string {
	if u.Profile == nil {
		return "UnratedAnon"
	}
	return RelevantRating(u.Profile.Ratings, ratingKey)
}

// GetRating gets a full Glicko-225 rating for this user, based on the
// passed-in rating key.
func (u *User) GetRating(ratingKey VariantKey) (*SingleRating, error) {
	if u.Profile == nil {
		return nil, errors.New("anonymous user has no rating")
	}
	defaultRating := &SingleRating{
		Rating:          float64(glicko.InitialRating),
		RatingDeviation: float64(glicko.InitialRatingDeviation),
		Volatility:      glicko.InitialVolatility,
	}
	if u.Profile.Ratings.Data == nil {
		return defaultRating, nil
	}
	ratdict, ok := u.Profile.Ratings.Data[ratingKey]
	if !ok {
		// Ratings dictionary exists, but user has no rating for this variant.
		return defaultRating, nil
	}
	return &ratdict, nil
}

// RealName returns a user's real name, or an empty string if anonymous.
func (u *User) RealName() string {
	if u.Profile != nil {
		if u.Profile.FirstName != "" {
			if u.Profile.LastName != "" {
				return u.Profile.FirstName + " " + u.Profile.LastName
			} else {
				return u.Profile.FirstName
			}
		} else {
			return u.Profile.LastName
		}
	}
	return ""
}

// RealNameIfNotYouth returns a user's real name, only if they are older than
// 13. If a birth date has not been provided, do not show it.
func (u *User) RealNameIfNotYouth() string {
	if u.Profile == nil {
		return ""
	}
	if u.IsChild() == realtime.ChildStatus_NOT_CHILD {
		return u.RealName()
	}
	return ""
}

func (u *User) AvatarUrl() string {
	if u.IsBot && u.Profile.AvatarUrl == "" {
		return "https://woogles-prod-assets.s3.amazonaws.com/macondog.png"
	} else {
		return u.Profile.AvatarUrl
	}
}

// TournamentID returns the "player ID" of a user. UUID:username is probably not
// a good design, but let's at least narrow it down to this function.
func (u *User) TournamentID() string {
	return u.UUID + ":" + u.Username
}

func InferChildStatus(dob string, now time.Time) realtime.ChildStatus {
	// The birth date must be in the form YYYY-MM-DD
	birthDateTime, err := time.Parse(time.RFC3339Nano, dob+"T00:00:00.000Z")
	if err != nil {
		// This means the birth date was either not defined or malformed
		// Either way, the child status should be unknown
		return realtime.ChildStatus_UNKNOWN
	} else {
		timeOfNotChild := birthDateTime.AddDate(13, 0, 0)
		if now.After(timeOfNotChild) {
			return realtime.ChildStatus_NOT_CHILD
		} else {
			return realtime.ChildStatus_CHILD
		}
	}
}

func IsAdult(dob string, now time.Time) bool {
	return InferChildStatus(dob, now) == realtime.ChildStatus_NOT_CHILD
}

func (u *User) IsChild() realtime.ChildStatus {
	if u.Profile == nil {
		log.Error().Str("uuid", u.UUID).Msg("unexpected-nil-profile")
		return realtime.ChildStatus_UNKNOWN
	}
	return InferChildStatus(u.Profile.BirthDate, time.Now())
}
