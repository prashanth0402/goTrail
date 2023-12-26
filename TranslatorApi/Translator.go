package Translatorapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/rylio/ytdl"
	"github.com/siongui/goigetmedia"
)

// func DownloadYoutubeVideo(link string) error {
// 	// Create a new video info object
// 	videoInfo, err := ytdl.GetVideoInfo(link)
// 	if err != nil {
// 		return err
// 	}

// 	// Get the highest quality video stream
// 	format := videoInfo.Formats.Best(ytdl.FormatAudioEncodingKey)

// 	// Create a new file to save the video
// 	file, err := os.Create(videoInfo.Title + "." + format.Extension)
// 	if err != nil {
// 		return err
// 	}
// 	defer file.Close()

// 	// Download the video
// 	_, err = videoInfo.Download(format, file)
// 	if err != nil {
// 		return err
// 	}

// 	fmt.Println("Downloaded:", videoInfo.Title)
// 	return nil
// }

// func downloadInstagramMedia(link string) error {
// 	media, err := goigetmedia.Get(link)
// 	if err != nil {
// 		return err
// 	}

// 	// Replace with the desired file name
// 	fileName := "downloaded_media"
// 	if media.IsVideo {
// 		fileName += ".mp4"
// 	} else {
// 		fileName += ".jpg"
// 	}

// 	err = goigetmedia.Download(media.URL, fileName)
// 	if err != nil {
// 		return err
// 	}

// 	fmt.Printf("Downloaded: %s\n", fileName)
// 	return nil
// }

// func main() {
// 	// Replace with the link to the Instagram post you want to download
// 	instagramLink := "https://www.instagram.com/p/example/"

// 	err := downloadInstagramMedia(instagramLink)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// }

func downloadMedia(link string) error {
	if strings.Contains(link, "youtube.com") || strings.Contains(link, "youtu.be") {
		// Download YouTube video
		return downloadYouTubeVideo(link)
	} else if strings.Contains(link, "instagram.com/p/") {
		// Download Instagram media
		return downloadInstagramMedia(link)
	} else {
		return fmt.Errorf("Unsupported media link: %s", link)
	}
}

func downloadYouTubeVideo(link string) error {
	videoInfo, err := ytdl.GetVideoInfo(link)
	if err != nil {
		return err
	}

	format := videoInfo.Formats.Best(ytdl.FormatAudioEncodingKey)

	file, err := os.Create(videoInfo.Title + "." + format.Extension)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = videoInfo.Download(format, file)
	if err != nil {
		return err
	}

	fmt.Println("Downloaded YouTube video:", videoInfo.Title)
	return nil
}

func downloadInstagramMedia(link string) error {
	media, err := goigetmedia.Get(link)
	if err != nil {
		return err
	}

	fileName := "downloaded_media"
	if media.IsVideo {
		fileName += ".mp4"
	} else {
		fileName += ".jpg"
	}

	err = goigetmedia.Download(media.URL, fileName)
	if err != nil {
		return err
	}

	fmt.Printf("Downloaded Instagram media: %s\n", fileName)
	return nil
}

// =============git hub link======================
// go get -u github.com/rylio/ytdl
// go get -u github.com/siongui/goigetmedia
// ==============================================

// ====================================Thumbnail process===================
type ThumbnailResponse struct {
	ThumbnailURL string `json:"thumbnailURL"`
}

func fetchThumbnail(link string) (string, error) {
	if strings.Contains(link, "youtube.com") || strings.Contains(link, "youtu.be") {
		// Fetch YouTube thumbnail
		return fetchYouTubeThumbnail(link)
	} else if strings.Contains(link, "instagram.com/p/") {
		// Fetch Instagram thumbnail
		return fetchInstagramThumbnail(link)
	} else {
		return "", fmt.Errorf("Unsupported media link: %s", link)
	}
}

func fetchYouTubeThumbnail(link string) (string, error) {
	videoInfo, err := ytdl.GetVideoInfo(link)
	if err != nil {
		return "", err
	}

	// Get the highest quality thumbnail
	thumbnailURL := videoInfo.Thumbnails[len(videoInfo.Thumbnails)-1].URL

	return thumbnailURL, nil
}

func fetchInstagramThumbnail(link string) (string, error) {
	media, err := goigetmedia.Get(link)
	if err != nil {
		return "", err
	}

	// Instagram provides thumbnails directly in the media object
	thumbnailURL := media.ThumbnailSrc

	return thumbnailURL, nil
}

func thumbnailHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	link := r.FormValue("link")
	if link == "" {
		http.Error(w, "Missing 'link' parameter", http.StatusBadRequest)
		return
	}

	thumbnailURL, err := fetchThumbnail(link)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching thumbnail: %v", err), http.StatusInternalServerError)
		return
	}

	response := ThumbnailResponse{
		ThumbnailURL: thumbnailURL,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// ===============Audio Extracter==================================
// =====================git for audio===================
// "github.com/rylio/ytdl"
// "github.com/iwittkau/m3u8"
// ======================================================

func extractAudio(link, fileName string) error {
	if strings.Contains(link, "youtube.com") || strings.Contains(link, "youtu.be") {
		// Extract audio from YouTube video
		return extractAudioFromYouTube(link, fileName)
	} else {
		// Extract audio from general video link
		return extractAudioFromVideo(link, fileName)
	}
}

func extractAudioFromYouTube(link, fileName string) error {
	videoInfo, err := ytdl.GetVideoInfo(link)
	if err != nil {
		return err
	}

	format := videoInfo.Formats.Best(ytdl.FormatAudioEncodingKey)

	file, err := os.Create(fileName + "." + format.Extension)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = videoInfo.Download(format, file)
	if err != nil {
		return err
	}

	fmt.Println("Extracted audio from YouTube video:", videoInfo.Title)
	return nil
}

func extractAudioFromVideo(link, fileName string) error {
	resp, err := http.Get(link)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	playlist, listType, err := m3u8.DecodeFrom(resp.Body, true)
	if err != nil {
		return err
	}

	if listType == m3u8.MASTER {
		// If it's a master playlist, choose the first variant
		variant := playlist.(*m3u8.MasterPlaylist).Variants[0]
		resp, err := http.Get(variant.URI)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		playlist, _, err = m3u8.DecodeFrom(resp.Body, true)
		if err != nil {
			return err
		}
	}

	// Find the first audio segment and download it
	audioSegment := findAudioSegment(playlist.(*m3u8.MediaPlaylist))
	if audioSegment == nil {
		return fmt.Errorf("No audio segment found in the playlist")
	}

	file, err := os.Create(fileName + ".mp3")
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, audioSegment.Data)
	if err != nil {
		return err
	}

	fmt.Println("Extracted audio from video:", fileName)
	return nil
}

func findAudioSegment(playlist *m3u8.MediaPlaylist) *m3u8.MediaSegment {
	for _, segment := range playlist.Segments {
		if segment != nil && segment.Key != nil && segment.Key.Method == "AES-128" {
			return segment
		}
	}
	return nil
}

// ==========================================Language Finder=================================

//    "github.com/iwat/talisman"

func identifyLanguage(audioFilePath string) (string, error) {
	recognizer, err := talisman.NewRecognizer()
	if err != nil {
		return "", err
	}
	defer recognizer.Close()

	file, err := os.Open(audioFilePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	language, err := recognizer.Identify(file)
	if err != nil {
		return "", err
	}

	return language, nil
}

//============================Audio Translator=========================================

// speech "cloud.google.com/go/speech/apiv1"
// 	"google.golang.org/api/option"
// 	speechpb "google.golang.org/genproto/googleapis/cloud/speech/v1"

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <audio-file>\n", os.Args[0])
		os.Exit(1)
	}

	audioFile := os.Args[1]

	ctx := context.Background()

	client, err := speech.NewClient(ctx, option.WithCredentialsFile("path/to/your/credentials.json"))
	if err != nil {
		fmt.Printf("Failed to create client: %v\n", err)
		return
	}
	defer client.Close()

	data, err := os.ReadFile(audioFile)
	if err != nil {
		fmt.Printf("Failed to read audio file: %v\n", err)
		return
	}

	req := &speechpb.RecognizeRequest{
		Config: &speechpb.RecognitionConfig{
			Encoding:        speechpb.RecognitionConfig_LINEAR16,
			SampleRateHertz: 16000,
			LanguageCode:    "ta-IN",
		},
		Audio: &speechpb.RecognitionAudio{
			AudioSource: &speechpb.RecognitionAudio_Content{Content: data},
		},
	}

	resp, err := client.Recognize(ctx, req)
	if err != nil {
		fmt.Printf("Failed to recognize: %v\n", err)
		return
	}

	for _, result := range resp.Results {
		for _, alt := range result.Alternatives {
			fmt.Printf("Transcript: %s\n", alt.Transcript)
		}
	}
}

// =======================combineAudioAndVideo============================================================

func combineAudioAndVideo(w http.ResponseWriter, r *http.Request) {
	// Parse multipart form
	err := r.ParseMultipartForm(10 << 20) // 10 MB limit
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	// Get the files
	audioFile, _, err := r.FormFile("audio")
	if err != nil {
		http.Error(w, "Error retrieving audio file", http.StatusBadRequest)
		return
	}
	defer audioFile.Close()

	videoFile, _, err := r.FormFile("video")
	if err != nil {
		http.Error(w, "Error retrieving video file", http.StatusBadRequest)
		return
	}
	defer videoFile.Close()

	// Save audio file
	audioPath := saveFile("audio", audioFile)
	defer os.Remove(audioPath)

	// Save video file
	videoPath := saveFile("video", videoFile)
	defer os.Remove(videoPath)

	// Combine audio and video using FFmpeg
	outputPath := filepath.Join("output", fmt.Sprintf("output_%d.mp4", time.Now().Unix()))
	cmd := exec.Command("ffmpeg", "-i", videoPath, "-i", audioPath, "-c:v", "copy", "-c:a", "aac", outputPath)

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	err = cmd.Run()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error combining audio and video: %s", out.String()), http.StatusInternalServerError)
		return
	}

	// Open the combined video file
	outputFile, err := os.Open(outputPath)
	if err != nil {
		http.Error(w, "Error opening combined video file", http.StatusInternalServerError)
		return
	}
	defer outputFile.Close()

	// Set content type and send the combined video file to the front end
	w.Header().Set("Content-Type", "video/mp4")
	io.Copy(w, outputFile)
}

func saveFile(fileType string, file multipart.File) string {
	// Create a unique filename for the saved file
	filename := fmt.Sprintf("%s_%d", fileType, time.Now().Unix())
	filePath := filepath.Join("uploads", filename)

	// Create the file
	out, err := os.Create(filePath)
	if err != nil {
		panic(err)
	}
	defer out.Close()

	// Copy the file data to the new file
	_, err = io.Copy(out, file)
	if err != nil {
		panic(err)
	}

	return filePath
}

// =================================convert audio to string

// go get github.com/go-audio/audio
// go get github.com/go-audio/wav

func ExtractAudioToString(filePath string) (string, error) {
	// Open the WAV file
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Create a WAV decoder
	decoder := wav.NewDecoder(file)

	// Read audio data
	buf := &audio.IntBuffer{Data: make([]int, 1024)}

	var audioData strings.Builder

	for {
		if err := decoder.Read(buf); err != nil {
			break
		}

		// Process the audio data (for example, convert it to a string)
		for _, val := range buf.Data {
			audioData.WriteString(fmt.Sprintf("%d ", val))
		}
	}

	return audioData.String(), nil
}

// =====================translate string =======================================
// "github.com/bregydoc/gtranslate"
// go get github.com/bregydoc/gtranslate

func stringTranslator() {
	originalString := "Hello, world!"

	// Translate to Tamil
	translatedTamil, err := gtranslate.Translate(originalString, "ta", "")
	if err != nil {
		fmt.Println("Error translating to Tamil:", err)
		return
	}

	// Translate to Hindi
	translatedHindi, err := gtranslate.Translate(originalString, "hi", "")
	if err != nil {
		fmt.Println("Error translating to Hindi:", err)
		return
	}

	// Translate to Malayalam
	translatedMalayalam, err := gtranslate.Translate(originalString, "ml", "")
	if err != nil {
		fmt.Println("Error translating to Malayalam:", err)
		return
	}

	// Translate to Telugu
	translatedTelugu, err := gtranslate.Translate(originalString, "te", "")
	if err != nil {
		fmt.Println("Error translating to Telugu:", err)
		return
	}

	fmt.Println("Original:", originalString)
	fmt.Println("Translated to Tamil:", translatedTamil)
	fmt.Println("Translated to Hindi:", translatedHindi)
	fmt.Println("Translated to Malayalam:", translatedMalayalam)
	fmt.Println("Translated to Telugu:", translatedTelugu)
}

// =======================================text to speech=================================

func textToSpeech(text string, languageCode string, outputPath string) error {
	ctx := context.Background()

	client, err := tts.NewClient(ctx)
	if err != nil {
		return err
	}
	defer client.Close()

	req := &ttspb.SynthesizeSpeechRequest{
		Input: &ttspb.SynthesisInput{
			InputSource: &ttspb.SynthesisInput_Text{Text: text},
		},
		Voice: &ttspb.VoiceSelectionParams{
			LanguageCode: languageCode,
			Name:         "en-US-Standard-C",
		},
		AudioConfig: &ttspb.AudioConfig{
			AudioEncoding: ttspb.AudioEncoding_LINEAR16,
		},
	}

	resp, err := client.SynthesizeSpeech(ctx, req)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(outputPath, resp.AudioContent, 0644)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	// Your text to convert into audio
	textToConvert := "Hello, world!"

	// Specify the language code (e.g., "en-US" for English)
	languageCode := "en-US"

	// Output file path for the generated audio
	outputPath := "output.wav"

	// Convert text to speech
	err := textToSpeech(textToConvert, languageCode, outputPath)
	if err != nil {
		log.Fatal(err)
		return
	}

	fmt.Printf("Text converted to audio. Audio saved to %s\n", outputPath)
}
// =============================================converting pitches ========================================================

package main

import (
	"fmt"
	"log"
	"os"

	"github.com/go-audio/audio"
	"github.com/go-audio/wav"
)

func main() {
	// Replace with the paths to your input audio files
	inputPath1 := "input1.wav"
	inputPath2 := "input2.wav"

	// Replace with the desired paths for the output audio files
	outputPath1 := "output1_pitch_matched.wav"
	outputPath2 := "output2_pitch_matched.wav"

	// Extract pitch from the first audio file
	pitch, err := extractPitch(inputPath1)
	if err != nil {
		log.Fatal(err)
	}

	// Change the pitch of the second audio file to match the extracted pitch
	err = changePitch(inputPath2, outputPath2, pitch)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Pitch of %s extracted as %.2f Hz\n", inputPath1, pitch)
	fmt.Printf("Pitch of %s matched to %s\n", inputPath2, outputPath2)
}

// extractPitch extracts the pitch (fundamental frequency) of an entire audio file.
func extractPitch(inputPath string) (float64, error) {
	// Open the input file
	inputFile, err := os.Open(inputPath)
	if err != nil {
		return 0, err
	}
	defer inputFile.Close()

	// Create a WAV decoder for the input file
	decoder := wav.NewDecoder(inputFile)

	// Create a buffer for audio processing
	bufferSize := decoder.Length()
	buf := &audio.IntBuffer{Data: make([]int, bufferSize)}

	// Read the entire audio file into the buffer
	if err := decoder.Read(buf); err != nil {
		return 0, err
	}

	// Perform pitch estimation on the entire audio file
	pitch := calculatePitch(buf.Data, decoder.SampleRate())

	return pitch, nil
}

// calculatePitch estimates the pitch (fundamental frequency) of an audio buffer.
func calculatePitch(buffer []int, sampleRate int) float64 {
	// Replace this with a more sophisticated pitch estimation algorithm if needed
	// This is a simple example using the autocorrelation method
	// You may need a dedicated pitch estimation library for better results

	// Autocorrelation
	var acf []float64
	for lag := 0; lag < len(buffer); lag++ {
		sum := 0.0
		for i := 0; i < len(buffer)-lag; i++ {
			sum += float64(buffer[i]) * float64(buffer[i+lag])
		}
		acf = append(acf, sum)
	}

	// Find the index of the maximum autocorrelation value
	maxIndex := 0
	maxValue := acf[0]
	for i, value := range acf {
		if value > maxValue {
			maxIndex = i
			maxValue = value
		}
	}

	// Calculate pitch in Hertz
	pitch := float64(sampleRate) / float64(maxIndex)
	return pitch
}

// changePitch changes the pitch of an audio file and saves the result.
func changePitch(inputPath, outputPath string, pitchFactor float64) error {
	// Open the input file
	inputFile, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer inputFile.Close()

	// Create a WAV decoder for the input file
	decoder := wav.NewDecoder(inputFile)

	// Create a WAV encoder for the output file
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	encoder := wav.NewEncoder(outputFile, decoder.SampleRate(), decoder.BitDepth(), decoder.NumChannels(), 1)

	// Create a buffer for audio processing
	bufferSize := 1024
	buf := &audio.FloatBuffer{Data: make([]float64, bufferSize)}

	// Process audio samples
	for {
		if err := decoder.ReadFloat(buf); err != nil {
			break
		}

		// Modify the pitch by changing the speed of playback
		buf.ResampleTo(buf, pitchFactor)

		// Write the modified samples to the output file
		if err := encoder.WriteFloat(buf); err != nil {
			return err
		}
	}

	// Close the encoder to finish writing the output file
	if err := encoder.Close(); err != nil {
		return err
	}

	return nil
}
