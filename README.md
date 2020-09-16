# Simple Whatsapp Login Session

A simple command-line feature to get the WhatsApp session file

this project is to support my main project related to the WhatsApp bot

# How to use 

rune the file, for example `./WhatsappLogin -p 68123123`

prepare your phone to scan the QR Code,

after scan the file `68123123.gob` is available

# Available command

`-o` Output dir, default dir is the current file is located
`-p` Phone Number, default is 6288

Example command : `./WhatsappLogin -o ~/sessions -p 63112311`

# Docker use

Run this single line docker command to easily get session using using docker

```
docker run --rm \
-v $(pwd):/go/src/github.com/k1m0ch1/WhatsappLogin/sessions \
k1m0ch1/whatsapplogin -p 628123123
```