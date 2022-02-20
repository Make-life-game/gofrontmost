# gofrontmost

Get MacOS frontmost Application by golang.

## how2use

```shell
go install github.com/Make-life-game/gofrontmost
```

Will output json-format process infomation (with active window title).
```shell
{"CreateTime":1645181816691,"Name":"goland","Pid":11262,"Ppid":1,"Title":"gofrontmost – README.md"}
```

## development

```shell
git clone https://github.com/Make-life-game/gofrontmost
cd gofrontmost

go build .
./gofrontmost
```