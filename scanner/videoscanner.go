package scanner

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/schollz/progressbar/v3"
	"gocv.io/x/gocv"
)

type Command struct {
	Function string
	IntArg   int
}

type Segment struct {
	OutputFn string
	Series   string
	Episode  int
	Start    int
	End      int
}

type VideoScanner struct {
	FileName string
	FPS      float64
	FrameNum int
	NFrames  float64
	Segments []*Segment

	series  string
	season  int
	episode int
	video   *gocv.VideoCapture
}

func NewScanner(fn string, series string, season, episode int) *VideoScanner {
	scan := VideoScanner{
		FileName: fn,
		Segments: make([]*Segment, 0, 10),
		series:   series,
		season:   season,
		episode:  episode,
	}
	var err error
	scan.video, err = gocv.VideoCaptureFile(fn)
	if err != nil {
		log.Fatal().Err(err).Str("FileName", fn).Msg("error opening file")
	}

	scan.NFrames = scan.video.Get(gocv.VideoCaptureFrameCount)
	scan.FPS = scan.video.Get(gocv.VideoCaptureFPS)
	log.Info().Str("FileName", fn).Float64("TotalFrames", scan.NFrames).Float64("FPS", scan.FPS).Msg("opened video file")

	return &scan
}

func (scan *VideoScanner) NextBlankFrame() int {
	img := gocv.NewMat()
	for scan.video.Read(&img) {
		isBlankFrame := false
		components := gocv.Split(img)
		nonzero := gocv.CountNonZero(components[0])

		if nonzero < 5000 {
			isBlankFrame = true
		}

		for _, c := range components {
			c.Close()
		}

		if isBlankFrame {
			return int(scan.video.Get(gocv.VideoCapturePosFrames)) - 1
		}
	}
	return int(scan.video.Get(gocv.VideoCapturePosFrames)) - 1
}

func (scan *VideoScanner) Encode() {
	fmt.Println("Commands:")
	fmt.Println("=========")
	for idx, segment := range scan.Segments {
		segment.OutputFn = fmt.Sprintf("%s S%dE%d.mp4", scan.series, scan.season, scan.episode+idx)
		cmd := fmt.Sprintf(`HandBrakeCLI --preset jflix -i %s -o %s --start-at frames:%d --stop-at frames:%d --preset-import-file jflix.json`, scan.FileName, segment.OutputFn, segment.Start, segment.Start)
		fmt.Println(cmd)
	}
}

func (scan *VideoScanner) FindSegments() {
	log.Info().Msg("Beginning search for segments")

	window := gocv.NewWindow("Video preview")

	img := gocv.NewMat()
	currFrame := 0
	bar := progressbar.Default(int64(scan.NFrames))

	fmt.Println("  a = forward 1  -  s = backward 1  -  d = forward 20 min  -  f = backward 20 min")
	fmt.Println("  j = forward 100 frames  -  k = back 100 frames")
	fmt.Println("  <space> = mark segment  -  e = exit")

	var currSegment *Segment

	exit := false
	for !exit && scan.video.Read(&img) {
		window.IMShow(img)

		cont := true
		skipFrame := false
		for cont {
			key := window.WaitKey(0)

			switch key {
			case 'e':
				exit = true
				cont = false
			case 'z':
				currFrame = scan.NextBlankFrame()
				cont = false
			case 'a':
				cont = false
			case 's':
				currFrame -= 1
				scan.video.Set(gocv.VideoCapturePosFrames, float64(currFrame))
				skipFrame = true
				cont = false
			case 'd':
				currFrame += int(20 * 60 * scan.FPS)
				scan.video.Set(gocv.VideoCapturePosFrames, float64(currFrame))
				skipFrame = true
				cont = false
			case 'f':
				currFrame -= int(20 * 60 * scan.FPS)
				scan.video.Set(gocv.VideoCapturePosFrames, float64(currFrame))
				skipFrame = true
				cont = false
			case 'j':
				currFrame += 100
				scan.video.Set(gocv.VideoCapturePosFrames, float64(currFrame))
				skipFrame = true
				cont = false
			case 'k':
				currFrame -= 100
				scan.video.Set(gocv.VideoCapturePosFrames, float64(currFrame))
				skipFrame = true
				cont = false
			case ' ':
				if currSegment == nil {
					currSegment = &Segment{
						Start: currFrame,
						End:   -1,
					}
					log.Info().Msg("start new segment")
				} else {
					currSegment.End = currFrame
					scan.Segments = append(scan.Segments, currSegment)
					log.Info().Msg("closed segment")
				}
				cont = false
			default:
				fmt.Println("Unknown key")
			}
		}

		if !skipFrame {
			currFrame++
		}
		bar.Set(currFrame)
	}

	fmt.Printf("Segments: %d\n", len(scan.Segments))
	for _, segment := range scan.Segments {
		fmt.Printf("\tSegment: %d to %d\n", segment.Start, segment.End)
	}
}
