package ringtone

import (
	"strconv"

	mtunes "github.com/yinyajiang/go-tunes"
	"github.com/yinyajiang/go-tunes/service/fileservice"
)

/*
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>Ringtones</key>
	<dict>
		<key>SCJZ.m4r</key>
		<dict>
			<key>Name</key>
			<string>家乡 - 赵雷</string>
			<key>GUID</key>
			<string>DF372091E523C114</string>
			<key>Total Time</key>
			<integer>304559</integer>
			<key>PID</key>
			<integer>6157754780804351850</integer>
			<key>Protected Content</key>
			<false/>
		</dict>
		<key>BCMJ.m4r</key>
		<dict>
			<key>Name</key>
			<string>寂寞的人伤心的歌 - 龙梅子,杨海彪</string>
			<key>GUID</key>
			<string>C62E1D551EDCC253</string>
			<key>Total Time</key>
			<integer>210780</integer>
			<key>PID</key>
			<integer>6430244051389038143</integer>
			<key>Protected Content</key>
			<false/>
		</dict>
	</dict>
</dict>
</plist>
*/

//Parse ...
func Parse(db []byte, purchased bool, fs fileservice.Service) (retAll map[uint64]*TrackInfo, nameSet, fileNameSet map[string]struct{}) {
	type plItemInfo map[string]interface{}
	stRoot := struct {
		Ringtones map[string]plItemInfo `plist:"Ringtones"`
	}{}

	mtunes.UnmashalPlist(db, &stRoot)

	_str := func(infoMap *plItemInfo, key string) (s string) {
		v, ok := (*infoMap)[key]
		if !ok {
			return
		}
		s, _ = v.(string)
		return
	}

	_num := func(infoMap *plItemInfo, key string) (i uint64) {
		v, ok := (*infoMap)[key]
		if !ok {
			return
		}
		i, ok = v.(uint64)
		if !ok {
			ii, ok := v.(int64)
			i = uint64(ii)
			if !ok {
				s := _str(infoMap, key)
				if len(s) > 0 {
					i, _ = strconv.ParseUint(s, 0, 64)
				}
			}
		}
		return
	}

	_bool := func(infoMap *plItemInfo, key string) (b bool) {
		v, ok := (*infoMap)[key]
		if !ok {
			return
		}
		b, _ = v.(bool)
		return
	}

	retAll = make(map[uint64]*TrackInfo, 0)
	nameSet = make(map[string]struct{}, 0)
	fileNameSet = make(map[string]struct{}, 0)

	for key, infoMap := range stRoot.Ringtones {

		info := &TrackInfo{
			FileName:  key,
			Name:      _str(&infoMap, "Name"),
			GUID:      _str(&infoMap, "GUID"),
			TotalTime: _num(&infoMap, "Total Time"),
			PID:       _num(&infoMap, "PID"),
			Protected: _bool(&infoMap, "Protected Content"),
		}
		if purchased {
			info.Path = "/Purchases/" + key
		} else {
			info.Path = "/iTunes_Control/Ringtones/" + key
		}
		info.Purchased = purchased

		finfo := fs.GetFileInfo(info.Path)
		if finfo == nil {
			info.isInvalid = true
			info.Path = ""
		} else {
			info.Size = finfo.Size
		}
		nameSet[info.Name] = struct{}{}
		fileNameSet[info.FileName] = struct{}{}
		retAll[info.PID] = info
	}
	return
}
