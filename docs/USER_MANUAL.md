# User Manual of commands
The following manual explains the different commands of this application, as well as **how** and **when** to use them.
## Slash Commands
These are the commands that are globally available, which can be used anywhere, regardless of the current channel.
### General
- `/play` - adds a song to the playlist / queue, or immediately starts auto-playing them if there are no songs in the queue 
  - **options**
    - `query` - this is where your **search term** / song name goes to
    - `platform` - you can manually specify the platform from the preset available platforms, otherwise it defaults to searching in the uploaded files
      - _Note_: This feature is currently not available, and will most likely be removed in the future 
    - `loop` - you can specify how many times the song should be played, defaults to play only once
- `/list` - display all the songs in the playlist / queue, in the order of they are going to be played
- `/pause` - **toggle** stop / resume the current song, depending on whether its currently playing or not
- `/skip` - skips current song, starts playing the next in the queue (if there is a next one)
- `/quit` - remove bot from the current voice channel. _Note_: it does not stop the bot itself

### File player specific
- `/upload` - drag and drop a media file (supported by FFMPEG; see [here](https://ffmpeg.org/ffmpeg-formats.html))
- `/uploaded` - displays all uploaded / hosted songs, available to stream within the current guild 

## Button Actions
Most message feedbacks sent by the bot have buttons, which mostly coincide with the slash commands, but there are some that are only accessible with the use of buttons.
These include:
- ` [Cancel]` - appears when a song was added to the queue, used to undo the addition of the song found to the queue / playlist
- ` [Show queue]` - functionally the same as writing `/list`
- ` [Show all] | [Collapse list]` - if there are many songs in the queue / playlist currently, this toggles showing all or just the **upcoming 5 songs**. The collapsed version is synonymous with the `/list` command mentioned above.
