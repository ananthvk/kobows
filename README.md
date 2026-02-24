# kobows (Kobo Web Server)

A web application that can be run on the Kobo e-reader so that `epub` files can be sent to the device over Wi-Fi, it also converts the `epub` file to a `kepub` file using the [kepubify](https://github.com/pgaskin/kepubify) library for improved perfomrance.


This application can be useful when you need to transfer an ebook file in a hurry, and you do not have an usb cable / pc nearby

---

⚠️ **WARNING: Run on LAN only!**

This server has **no authentication** and should only be run on a trusted local network. Do not expose it to the internet.

**Limitations:**
- Maximum file size: **50MB** (can be modified in source code)
- Only **EPUB** files are supported for now

## How to build
```
CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7 go build -o kobows-armv7
```

Then copy, `kobows-armv7`, `templates/` and `static/` folders to `/mnt/onboard/.adds/kobows` (When the reader is connected to the computer, in the root folder, inside the `.adds` directory, make a directory called `kobows`)

Then put the following lines in the NickelMenu config file

```
menu_item :main    :KoboWS Kill              :cmd_output         :500:quiet  :/usr/bin/pkill -f kobows-armv7
  chain_success                              :dbg_toast          :Killed KoboWS
menu_item :main    :KoboWS Start             :cmd_spawn          :sh -c 'cd /mnt/onboard/.adds/kobows && ./kobows-armv7'
  chain_failure                              :dbg_toast          :KoboWS start failed due to error
```

## To access the server go to

```
http://<kobo lan ip>:8000/

For example,

http://192.168.1.100:8000/
```
