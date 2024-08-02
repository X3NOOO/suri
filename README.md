# suri

AI assistant 

## API docs

Suri provides an (hopefully) easy-to-use API that allows for an easy implementation of different clients.

### Endpoints

- GET /status
- POST /login
- POST /query
- - Custom headers: `X-Mute-Server-Audio`, `X-No-Response`, Supported `Accept:` types: `text/plain`, `application/json` (default), `application/wav`
- GET,POST /knowledge
- GET,DELETE,PUT /knowledge/{name}
- ??? /settings


### Todo

- Authentication

### Credits

Special thanks to:
- The [Piper Project](https://github.com/rhasspy/piper) for an amazing local TTS software
- [Lingoose](https://github.com/henomis/lingoose), we use a lot of the abstraction layer they provide for interacting with the LLMs
- [ffmpeg](https://ffmpeg.org), of course
