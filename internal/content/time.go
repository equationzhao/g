package content

import (
	"runtime"
	"strings"
	"time"

	constval "github.com/Equationzhao/g/internal/global"
	"github.com/Equationzhao/g/internal/item"
	"github.com/Equationzhao/g/internal/osbased"
	"github.com/Equationzhao/g/internal/render"
	strftime "github.com/itchyny/timefmt-go"
)

type RelativeTimeEnabler struct {
	Mode string
}

func NewRelativeTimeEnabler() *RelativeTimeEnabler {
	return &RelativeTimeEnabler{}
}

const RelativeTime = constval.NameOfRelativeTime

func (r *RelativeTimeEnabler) Enable(renderer *render.Renderer) ContentOption {
	return func(info *item.FileInfo) (string, string) {
		var t time.Time
		timeType := ""
		switch r.Mode {
		case "mod":
			t = osbased.ModTime(info)
			timeType = timeModified
		case "create":
			t = osbased.CreateTime(info)
			timeType = timeCreated
		case "access":
			t = osbased.AccessTime(info)
			timeType = timeAccessed
		case "birth":
			timeType = timeBirth
			// if darwin, check birth time
			if runtime.GOOS == "darwin" {
				t = osbased.BirthTime(info)
			} else {
				t = osbased.CreateTime(info)
			}
		default:
			t = osbased.ModTime(info)
			timeType = timeModified
		}
		return renderer.RTime(time.Now(), t), RelativeTime + " " + timeType
	}
}

const (
	timeName     = constval.NameOfTime
	timeModified = constval.NameOfTimeModified
	timeCreated  = constval.NameOfTimeCreated
	timeAccessed = constval.NameOfTimeAccessed
	timeBirth    = constval.NameOfTimeBirth
)

// EnableTime enables time
// accepts ['mod', 'modified', 'create', 'access', 'birth']
func EnableTime(format, mode string, renderer *render.Renderer) ContentOption {
	return func(info *item.FileInfo) (string, string) {
		// get mod time/ create time/ access time
		var t time.Time
		timeType := ""
		switch mode {
		case "mod", "modified":
			t = osbased.ModTime(info)
			timeType = timeModified
		case "create", "cr":
			t = osbased.CreateTime(info)
			timeType = timeCreated
		case "access", "ac":
			t = osbased.AccessTime(info)
			timeType = timeAccessed
		case "birth":
			timeType = timeBirth
			// if darwin, check birth time
			if runtime.GOOS == "darwin" {
				t = osbased.BirthTime(info)
			} else {
				t = osbased.CreateTime(info)
			}
		default:
			t = osbased.ModTime(info)
			timeType = timeModified
		}

		var timeString string
		if after, ok := strings.CutPrefix(format, "+"); ok {
			timeString = strftime.Format(t, after)
		} else {
			timeString = t.Format(format)
		}
		return renderer.Time(timeString), timeName + " " + timeType
	}
}
