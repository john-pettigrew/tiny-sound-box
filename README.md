# Tiny Sound Box 📢
A lightweight API for playing local WAV audio files.

## Overview
The Tiny Sound Box API allows you to play local WAV files programmatically using HTTP requests. It supports looping audio files with an optional delay time.

## Getting Started
### Using Docker
To build the docker image, run:
```sh
docker build -t tiny-sound-box .
```
Then you can start the container using:
```sh
docker run -p 8500:8500 -v path/to/sounds:/sounds:ro --device /dev/snd tiny-sound-box
```
This will start Tiny Sound Box on port "8500".

**By default, the docker container will have a "--timeout-seconds" value of "600". Meaning, a sound loop will be stopped early after the sound's current loop and delay are complete if it has been playing and delaying for more than 600 seconds. If needed, this value can be changed. See "CLI Arguments" below.**

### Requirements
- aplay

### CLI Arguments
The following options are available when running the application:
- `--addr=<address>`: Specify on which address/port to listen for incoming requests (default ":8500").
- `--sounds-dir=<directory>`: Directory containing WAV audio files.
- `--num-workers=<number>`: Number of workers to start for playing audio files. **If all workers are busy, new requests to play files will be blocked until a worker is available. A worker is not done playing a sound until all loops and delays for that sound are completed.** So this should be at least the number of audio files that you want to play at the same time. (default 3)
- `--timeout-seconds=<seconds>`: Max number of seconds a worker should continue looping a sound file. Note: A playing sound file should not be interrupted and is allowed to complete the current loop and delay time before it is stopped even if more than "timeout-seconds" has elapsed.

## Routes
### GET /health
Check if the server is running.

### GET /play
Start playing a new sound.

**Note**: If all workers are busy, this request will block until a worker is available. See the "--num-workers" section under "CLI Arguments" above.

#### Query Parameters
- `sound`: (required) The name of the sound to play. This is the sound file name without the ".wav" extension. The sound name should not contain any special characters.
- `loop`: Number of times to loop playing the sound. (default 1)
- `delay`: Number of seconds to wait after playing a sound before the next loop. (default 0)

#### Example
This would play "hello.wav" 3 times with a 1 second delay between loops.
```
/play?sound=hello&loop=3&delay=1
```

### GET /stop-all
Stop each looping sound after the current loop and delay are completed.