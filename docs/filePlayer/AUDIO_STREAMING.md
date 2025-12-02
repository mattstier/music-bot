## Audio Streaming
### Using FFMPEG to get PCM data
```go
func (player *FilePlayer) startFFMPEG() (*exec.Cmd, context.CancelFunc) {
	fmt.Println(audioPath)
	ctx, cancel := context.WithCancel(context.Background())
	cmd := exec.CommandContext(ctx, "ffmpeg",
		"-ss", fmt.Sprintf("%.3f", player.timestamp.Seconds()),
		"-i", mediaDir+"/"+player.CurrentSong(),
		"-af", "aresample=resampler=soxr:osf=s16:dither_method=shibata",
		"-loglevel", "quiet",
		"-ar", strconv.Itoa(sampleRate),
		"-ac", strconv.Itoa(channels),
		"-f", "s16le",
		"pipe:1",
	)
	return cmd, cancel
}
```
What it does: 
- Wraps FFMPEG command with a cancellable context so audio playback can be cleanly stopped elsewhere.
- `-ss` - Starts playing audio from the player's timestamp using .
- `-i` - Loads the current song from mediaDir.
- `-loglevel quiet` - Hides ffmpeg output, remove it for verbose output
- `-af, aresample=resampler=soxr:osf=s16:dither_method=shibata` - Uses the soxr high-quality resampler with Shibata dithering, and converts the output to signed 16-bit samples.
- `-ac` - number of channels (we use 2, for stereo, which Discord supports)
- `-f s16le` - Denotes signed 16-bit little-endian **_PCM output_**
> [!NOTE]
> This PCM output has to be encoded into **_OPUS frames_** when sending it to a voice channel

### What is called 
```go
func (player *FilePlayer) streamAudio(vc *discordgo.VoiceConnection) {
	//creating pipe with a ffmpeg command
	cmd, cancel := player.startFFMPEG()
	stdout, _ := cmd.StdoutPipe()
	cmd.Stderr = os.Stderr
	cmd.Start()
	defer cmd.Wait()
	//local done channel
	sessionDone := make(chan struct{})

	// buffered channels
	pcmChannel := make(chan []byte, 50)
	opusChannel := make(chan []byte, 50)
}
```
What it does: 
- Gets FFMPEG's PCM output using its _stdout pipe_ 
- Defines a (_sessionDone_) channel for when the music naturally finishes (without interruption) from the player's done channel
- Creates buffered channels for the **raw PCM output** and the **encoded OPUS frames** (see [link](#encoding-pcm-to-opus))

```go
//handles stopping, song finishing etc
	go func() {
		select {
		//triggers if done has been signaled
		case <-player.done:
			stdout.Close() // always close the pipe first, to not break it
			cancel()
			fmt.Println("Player stopped")
		case <-sessionDone:
			//reset timestamp, so next song plays from beginning
			player.timestamp = 0
			fmt.Println("Song finished")
		}
		player.isPlaying = false
		player.done <- struct{}{}

	}()
```
What it does:
- Gets FFMPEG's PCM output using its _stdout pipe_
- Defines a (_sessionDone_) channel for when the music naturally finishes (without interruption) from the player's done channel
- Creates buffered channels for the **raw PCM output** and the **encoded OPUS frames** (see [link](#encoding-pcm-to-opus))

### 'Conversion pipeline'
```go
	//sending/streaming pcm into the pcm channel
	go bufferPCM(stdout, pcmChannel)
	//encoding and sending pcms into opus frames to the opus channel
	go encodePCM(pcmChannel, opusChannel)
	// 200 ms latency cushion
	time.Sleep(200 * time.Millisecond)
	//sending opus frames to the VoiceConnection
	go sendOpus(opusChannel, vc, player, sessionDone) //waiting for done to be true
}
```
What it does:
- Starts **buffering the PCMs** coming from the stdout pipe 
- And in the meantime **concurrently encoding them** into OPUS
- As well as constantly sending them to the voice channel after 200 millisecond 'latency cushion' to give the buffer some time to fill up

### Encoding PCM to OPUS
```go
// encodes pcm to Opus format
func encodePCM(pcmChannel chan []byte, opusChannel chan []byte) {
	enc, _ := opus.NewEncoder(sampleRate, channels, opus.AppVoIP)
	data := make([]byte, frameSize)
	for frame := range pcmChannel {
		samples := bytesToInt16(frame)
		e, err := enc.Encode(samples, data)
		if err != nil {
			break
		}
		pkt := make([]byte, e)
		copy(pkt, data[:e])
		opusChannel <- pkt
	}
	close(opusChannel)
}
```
What it does:
- Takes the buffered channels from above (one for PCM one for the encoded OPUS) and
- Encodes the PCM channels content to the OPUS channel on the flight by
  - Using the [OPUS encoding library by hraban](https://github.com/hraban/opus) by calling `enc.Encode()` and writing to the buffer called '_data_'
  - Converting the input slice ([]byte) to a **16-bit int slice**, as that's what the encoder takes as an input
  - Copying data into a packet to be sent to the opus channel

### Sending OPUS to the Voice Channel
```go
func sendOpus(opusChannel chan []byte, vc *discordgo.VoiceConnection, player *FilePlayer, sessionDone chan struct{}) {
	ticker := time.NewTicker(sendRate)
	defer ticker.Stop()
	for pkt := range opusChannel {
		<-ticker.C
		select {
		case vc.OpusSend <- pkt:
			//send, and increase timestamp if successful
			player.timestamp += sendRate
		case <-player.done:
			return
		default:
			//drop packet if it cannot be sent
		}

	}
	// if no more packets to send, signal done
	close(sessionDone)
}
```
What it does:
- Uses a ticker with a **20ms tick rate** (the expected sending rate for Discord's voice channels)
- To send each OPUS packets to the voice channel's channel with `vc.OpusSend`
- Or stops if the player.done channel triggers it
- Every time it sends a packet the players timestamp is incremented (roughly lining up with the actual timestamp in the song)