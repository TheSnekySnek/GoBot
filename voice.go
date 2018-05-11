package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"sync"

	"github.com/bwmarrin/discordgo"
	"layeh.com/gopus"
)

const (
	channels  int = 2               // 1 for mono, 2 for stereo
	frameRate int = 48000           // audio sampling rate
	frameSize int = 960             // uint16 size of each audio frame
	maxBytes  int = (frameSize * 2) // max size of opus data
	nightcore int = 38400
	daycore   int = 57600
)

var (
	speakers    map[uint32]*gopus.Decoder
	opusEncoder *gopus.Encoder
	mu          sync.Mutex
	hasStopped  bool
)

// OnError gets called by dgvoice when an error is encountered.
// By default logs to STDERR
var OnError = func(str string, err error) {
	prefix := "Voice Engine: " + str

	if err != nil {
		os.Stderr.WriteString(prefix + ": " + err.Error())
	} else {
		os.Stderr.WriteString(prefix)
	}
}

// SendPCM will receive on the provied channel encode
// received PCM data into Opus then send that to Discordgo
func SendPCM(v *discordgo.VoiceConnection, pcm <-chan []int16) {
	if pcm == nil {
		return
	}

	var err error

	opusEncoder, err = gopus.NewEncoder(frameRate, channels, gopus.Audio)
	opusEncoder.SetBitrate(96000)

	if err != nil {
		OnError("NewEncoder Error", err)
		return
	}

	for {

		// read pcm from chan, exit if channel is closed.
		recv, ok := <-pcm
		if !ok {
			OnError("PCM Channel closed", nil)
			return
		}

		// try encoding pcm frame with Opus

		opus, err := opusEncoder.Encode(recv, frameSize, maxBytes)
		if err != nil {
			OnError("Encoding Error", err)
			return
		}

		if v.Ready == false || v.OpusSend == nil {
			// OnError(fmt.Sprintf("Discordgo not ready for opus packets. %+v : %+v", v.Ready, v.OpusSend), nil)
			// Sending errors here might not be suited
			return
		}
		// send encoded opus data to the sendOpus channel
		v.OpusSend <- opus
	}
}

// PlayAudioFile will play the given filename to the already connected
// Discord voice server/channel.  voice websocket and udp socket
// must already be setup before this will work.
func PlayAudioFile(v *discordgo.VoiceConnection, filename string, mod int, volume float64, stop <-chan bool) {
	hasStopped = false
	freq := frameRate
	if mod == 1 {
		freq = nightcore
	} else if mod == 2 {
		freq = daycore
	}
	//vo := strconv.FormatFloat(volume/100, 'f', 2, 64)
	// Create a shell command "object" to run.
	run := exec.Command("ffmpeg", "-i", filename, "-f", "s16le", "-ar", strconv.Itoa(freq), "-ac", strconv.Itoa(channels), "-analyzeduration", "0", "pipe:1")
	ffmpegout, err := run.StdoutPipe()
	if err != nil {
		OnError("StdoutPipe Error", err)
		return
	}

	ffmpegbuf := bufio.NewReaderSize(ffmpegout, 16344)

	// Starts the ffmpeg command
	err = run.Start()
	if err != nil {
		OnError("RunStart Error", err)
		return
	}

	go func() {
		<-stop
		fmt.Println("KILL")
		err = run.Process.Kill()
		if err != nil {
			fmt.Println(err.Error())
		}
	}()

	// Send "speaking" packet over the voice websocket
	err = v.Speaking(true)
	fmt.Println("START")
	if err != nil {
		OnError("Couldn't set speaking", err)
	}

	// Send not "speaking" packet over the websocket when we finish
	defer func() {
		hasStopped = true
		err := v.Speaking(false)
		fmt.Println("STOP")
		if err != nil {
			OnError("Couldn't stop speaking", err)
		}
	}()

	send := make(chan []int16, 2)
	defer close(send)

	close := make(chan bool)
	go func() {
		SendPCM(v, send)

		close <- true
	}()

	for {
		// read data from ffmpeg stdout
		if !hasStopped {
			audiobuf := make([]int16, frameSize*channels)
			err = binary.Read(ffmpegbuf, binary.LittleEndian, &audiobuf)
			if err == io.EOF || err == io.ErrUnexpectedEOF {
				return
			}
			if err != nil {
				OnError("error reading from ffmpeg stdout", err)
				return
			}
			// Send received PCM to the sendPCM channel
			select {
			case send <- audiobuf:
			case <-close:
				return
			}
		}
	}
}
