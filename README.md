# jflix-split

open-cv tool to split recordings of TV Shows into mulitple files - one per episode.

# Building

* Install opencv 4.5.5+, follow instructions for your OS here: [https://gocv.io/getting-started/]()
* Install Handbrake and Handbrake CLI

```bash
go get
go build
```

# Running

create a new handbrake presets profile and name it `jflix`, export to json and place in the same directory as `jflix-split`

```bash
./jflix-split <file_name> -s 1 -e 1 --show "TV Show name"
```

## Navigation commands

* a = forward 1 frame
* s = backward 1 frame
* d = forward 20 min
* f = backward 20 min 
* j = forward 100 frames
* k = back 100 frames
* <space> = mark segment begin/end
* e = exit
