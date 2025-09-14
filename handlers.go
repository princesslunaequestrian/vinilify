package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	tg "github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
	"github.com/princesslunaequestrian/vinilify/types"
	"github.com/princesslunaequestrian/vinilify/utils"
	"github.com/princesslunaequestrian/vinilify/utils/converters"
)

const (
	//start command
	MessageSendStart      = "Send \"/start\" command to begin"
	MessageInstruction    = "Upload audio and cover image"
	MessageUnknownCommand = "Unknown command"

	//upload notificators
	MessageUploadedImage   = "Image uploaded"
	MessageUplodadedAudio  = "Audio uploaded"
	MessageReadyToGenerate = "Now send \"/generate\" to generate the vinyl"

	//error notificators
	MessageNoAudio = "You have not uploaded audio file"
	MessageNoImage = "You have not uploaded image file"

	//block notificators
	MessageGenerating = "Your video is being processed right now, wait for it to complete"
	MessageCooldown   = "You are in cooldown, wait..."

	//generator notificators
	MessageAudioDownloadFailed    = "Could not download audio"
	MessageImageDownloadFailed    = "Could not download image"
	MessageDownloadComplete       = "Files downloaded, generating video..."
	MessageDownloadStarted        = "Downloading files..."
	MessagePreparingForGeneration = "Preparing for generation..."
)

func handleStart(ctx *th.Context, update tg.Update) error {

	userID := update.Message.From.ID
	chatID := update.Message.Chat.ChatID()
	createdErr := error(nil)

	//check if user is in the map, if not - add
	_, ok := users[userID]
	if !ok {
		users[userID] = &types.User{
			Id:       userID,
			State:    0,
			Cooldown: time.Now(),
			AudioURL: "",
			ImageURL: "",
		}
	}

	// Checking if "users" directory exists
	usersExists, _ := utils.DirExists("./users")
	if !usersExists {
		// if not -- create
		createdErr = os.Mkdir("./users", 0755)
		if createdErr != nil {
			log.Fatal(createdErr)
		}
	}

	// Checking if user exists
	userDirExists, _ := utils.DirExists(fmt.Sprintf("./users/%d", userID))
	if !userDirExists {
		createdErr = os.Mkdir(fmt.Sprintf("./users/%d", userID), 0755)
		if createdErr != nil {
			log.Fatal(createdErr)
		}
	}

	sendMessage(ctx, types.Message{ChatID: chatID, Content: MessageInstruction})
	return createdErr
}

func handleUpload(ctx *th.Context, update tg.Update) error {
	userID := update.Message.From.ID
	chatID := update.Message.Chat.ChatID()

	user, ok := users[userID]
	if !ok {
		//REMOVE
		println(user)
		sendMessage(ctx, types.Message{ChatID: chatID, Content: MessageSendStart})
		return nil
	}

	if update.Message.Audio == nil && update.Message.Photo == nil {
		sendMessage(ctx, types.Message{ChatID: chatID, Content: MessageUnknownCommand})
		return nil
	}

	if update.Message.Photo != nil {
		photosCount := len(update.Message.Photo)

		file, err := ctx.Bot().GetFile(ctx, &tg.GetFileParams{FileID: update.Message.Photo[photosCount-1].FileID})
		if err != nil {
			log.Panic("can't retreive image file url")
			return err
		}
		user.ImageURL = ctx.Bot().FileDownloadURL(file.FilePath)
		sendMessage(ctx, types.Message{ChatID: chatID, Content: MessageUploadedImage})
	}

	if update.Message.Audio != nil {
		file, err := ctx.Bot().GetFile(ctx, &tg.GetFileParams{FileID: update.Message.Audio.FileID})
		if err != nil {
			log.Panic("can't retreive audio file url")
			return err
		}
		user.AudioURL = ctx.Bot().FileDownloadURL(file.FilePath)
		sendMessage(ctx, types.Message{ChatID: chatID, Content: MessageUplodadedAudio})
	}

	if user.HasAudioURL() && user.HasImageURL() {
		sendMessage(ctx, types.Message{ChatID: chatID, Content: MessageReadyToGenerate})
	}

	fmt.Printf("\n\n\n%+v\n\n\n", user)
	return nil

}

func handleGenerateVideo(ctx *th.Context, update tg.Update) error {
	userID := update.Message.From.ID
	chatID := update.Message.Chat.ChatID()
	user, ok := users[userID]
	if !ok {
		sendMessage(ctx, types.Message{ChatID: chatID, Content: MessageSendStart})
		return nil
	}

	//0. check if user is already generating video
	if user.Generating {
		sendMessage(ctx, types.Message{ChatID: chatID, Content: MessageGenerating})
		return nil
	}

	//0.1 check if user is in cooldown
	if time.Now().Compare(user.Cooldown) <= 0 {
		sendMessage(ctx, types.Message{ChatID: chatID, Content: MessageCooldown})
		return nil
	}

	//1. check if user has both audio and video file links
	if !user.HasAudioURL() {
		sendMessage(ctx, types.Message{ChatID: chatID, Content: MessageNoAudio})
		return nil
	}

	if !user.HasImageURL() {
		sendMessage(ctx, types.Message{ChatID: chatID, Content: MessageNoImage})
		return nil
	}

	//2. Go generate video note
	user.Generating = true

	go GenerateVideo(ctx, update, user)
	return nil
}

func GenerateVideo(ctx *th.Context, update tg.Update, user *types.User) error {

	defer func() {
		user.Generating = false
		user.Cooldown = time.Now().Add(time.Second * 100)
		user.AudioURL = ""
		user.ImageURL = ""
	}()

	chatID := update.Message.Chat.ChatID()

	assetsPath := utils.GetAssets()
	userPath := utils.GetUserPath(user.Id)

	//1. Download audio and image to the folder
	audioPath := user.GetAudioPath()
	imagePath := user.GetImagePath()
	sendMessage(ctx, types.Message{ChatID: chatID, Content: MessageDownloadStarted})

	err := utils.DownloadAttachment(audioPath, user.AudioURL)
	if err != nil {
		sendMessage(ctx, types.Message{ChatID: chatID, Content: MessageAudioDownloadFailed})

		return err
	}
	err = utils.DownloadAttachment(imagePath, user.ImageURL)
	if err != nil {
		sendMessage(ctx, types.Message{ChatID: chatID, Content: MessageImageDownloadFailed})

		return err
	}

	sendMessage(ctx, types.Message{ChatID: chatID, Content: MessageDownloadComplete})

	sendMessage(ctx, types.Message{ChatID: chatID, Content: MessagePreparingForGeneration})

	//2. Mix audio with effect

	effect := assetsPath + "/sounds/vinyl.mp3"
	music := userPath + "/audio.mp3"
	mix := userPath + "/mix.mp3" //mixed audio is stored in users/.../mix.mp3

	err = converters.Mix(effect, music, mix)
	if err != nil {
		sendMessage(ctx, types.Message{ChatID: chatID, Content: "Error mixing audio " + err.Error()})
		return err
	}

	//2.1 Add vinyl effects
	//2.1.1 Convert mp3 to wav
	mixWav := userPath + "/mix.wav"
	converters.Convert(mix, mixWav)

	//2.1.2 Add vinyl effects
	mixWavWithEffects := userPath + "/mix_with_effects.wav"
	converters.AddVinylEffects(mixWav, mixWavWithEffects)

	//2.1.3 Convert back to mp3
	mix = userPath + "/mix.mp3"
	converters.Convert(mixWavWithEffects, mix)

	//3. Generate images

	image := filepath.Join(utils.GetRoot(), "users", fmt.Sprintf("%d", user.Id), "image.jpg")
	imageOut := filepath.Join(utils.GetRoot(), "users", fmt.Sprintf("%d", user.Id))
	err = converters.AssembleImages(image, imageOut) //video frames are stored in users/.../01...32.png
	if err != nil {
		sendMessage(ctx, types.Message{ChatID: chatID, Content: "Error generating images " + err.Error()})
		return err
	}

	//4. Generate 1 second video
	patternPath := userPath + "/%02d.png"
	secondVideoPath := userPath + "/secondvideo.mp4"
	err = converters.SecondVideo(patternPath, secondVideoPath) //video is stored in users/.../secondvideo.mp4
	if err != nil {
		sendMessage(ctx, types.Message{ChatID: chatID, Content: "Error generating a second-long video " + err.Error()})
		return err
	}

	//5. Generate minute long video
	minuteVideoPath := userPath + "/minutevideo.mp4"
	err = converters.LoopVideo(secondVideoPath, minuteVideoPath)
	if err != nil {
		sendMessage(ctx, types.Message{ChatID: chatID, Content: "Error generating a minute-long video " + err.Error()})
		return err
	}

	//6. Mix audio and video together
	videoPath := userPath + "/output.mp4"
	err = converters.AddAudio(mix, minuteVideoPath, videoPath)
	if err != nil {
		sendMessage(ctx, types.Message{ChatID: chatID, Content: "Error generating the final video " + err.Error()})
		return err
	}

	//Check if the video file was actually created
	if _, err := os.Stat(videoPath); os.IsNotExist(err) {
		sendMessage(ctx, types.Message{ChatID: chatID, Content: "Failed to obtain the final video " + err.Error()})
		return err
	}

	//7. Send the video note to the user
	sendMessage(ctx, types.Message{ChatID: chatID, Content: "Video has been generated, sending..."})

	ctx.Bot().SendVideoNote(
		ctx,
		tu.VideoNote(
			update.Message.Chat.ChatID(),
			tu.File(mustOpen(videoPath)),
		),
	)

	return nil
}

func sendMessage(ctx *th.Context, msg types.Message) {
	ctx.Bot().SendMessage(ctx, tu.Message(msg.ChatID, msg.Content))
}

func mustOpen(filename string) *os.File {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	return file
}
