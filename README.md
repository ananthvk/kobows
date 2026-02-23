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