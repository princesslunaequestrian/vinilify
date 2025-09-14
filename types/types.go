package types

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	tg "github.com/mymmrac/telego"
	"github.com/princesslunaequestrian/vinilify/utils"
)

// Type for representing current state of user
type State int

// State constatnts
const (
	Welcome State = iota
	HasImage
	HasAudio
	HasBoth
	Generating
)

// user record structure
type User struct {
	Id         int64
	State      State
	Generating bool
	Cooldown   time.Time
	ImageURL   string
	AudioURL   string
}

type Message struct {
	ChatID  tg.ChatID
	Content string
}

var (
	// Error showing that there's nothing to display
	ErrorNothingToDisplay = errors.New("nothing to display")
	// Error showing that state value is unknown (not one of constants defined)
	ErrorUnknownState = errors.New("unknown state received")
	// No audio URL found for user
	ErrorNoAudioURL = errors.New("no audio URL")
	// No video URL found for user
	ErrorNoImageURL = errors.New("no video URL")
)

// Returns true if user's audio file specified
func (u User) HasAudioURL() bool {
	return u.AudioURL != ""
}

// Returns true if user's image file specified
func (u User) HasImageURL() bool {
	return u.ImageURL != ""
}

// Builds path to user's local audio file even without file existing
func (u User) GetAudioPath() string {
	return filepath.Join("users", fmt.Sprintf("%d", u.Id), "audio.mp3")
}

// Builds path to user's local image file even without file existing
func (u User) GetImagePath() string {
	return filepath.Join("users", fmt.Sprintf("%d", u.Id), "image.jpg")
}

// Gets path to user's local audio file creating it's local copy if one does not exists.
// Returns
func (u User) GetAudio() (string, error) {
	if !u.HasAudioURL() {
		return "", ErrorNoAudioURL
	}
	_, err := os.Stat(u.GetAudioPath())
	if os.IsNotExist(err) {
		return u.GetAudioPath(), utils.DownloadAttachment(u.GetAudioPath(), u.AudioURL)
	}
	return u.GetAudioPath(), nil
}

func (u User) GetImage() (string, error) {
	if !u.HasImageURL() {
		return "", ErrorNoImageURL
	}
	_, err := os.Stat(u.GetImagePath())
	if os.IsNotExist(err) {
		return u.GetImagePath(), utils.DownloadAttachment(u.GetImagePath(), u.ImageURL)
	}
	return "", err
}
