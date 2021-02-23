package ringtone

import (
	"time"

	mtunes "github.com/yinyajiang/go-tunes/mTunes"
)

/*
/iTunes_Control/Sync/Media/
*/
func operaitionsProto(syncNum int) []byte {
	proto := []byte(`
	<?xml version="1.0" encoding="UTF-8"?>
	<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
	<plist version="1.0">
	<dict>
		<key>timestamp</key>
		<date>2001-01-01T00:00:00Z</date>
		<key>operations</key>
		<array>
			<dict>
				<key>operation</key>
				<string>update_db_info</string>
				<key>db_info</key>
				<dict>
					<key>audio_language</key>
					<integer>0</integer>
					<key>subtitle_language</key>
					<integer>0</integer>
					<key>primary_container_pid</key>
					<integer>0</integer>
				</dict>
				<key>pid</key>
				<integer>0</integer>
			</dict>
		</array>
		<key>revision</key>
		<integer>0</integer>
	</dict>
	</plist>
	`)
	stProto := struct {
		Revision   int                      `plist:"revision"`
		Timestamp  time.Time                `plist:"timestamp"`
		Operations []map[string]interface{} `plist:"operations"`
	}{}
	mtunes.UnmashalPlist(proto, &stProto)
	stProto.Timestamp = time.Now()
	stProto.Revision = syncNum
	proto, _ = mtunes.MashalPlist(stProto)
	return proto
}

/*
/iTunes_Control/Ringtones/Sync/

<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>revision</key>
	<integer>0</integer>
	<key>timestamp</key>
	<date>2021-02-08T15:10:56Z</date>
	<key>operations</key>
	<array>
		<dict>
			<key>pid</key>
			<integer>0</integer>
			<key>operation</key>
			<string>update_db_info</string>
			<key>db_info</key>
			<dict>
				<key>primary_container_pid</key>
				<integer>0</integer>
				<key>audio_language</key>
				<integer>0</integer>
				<key>subtitle_language</key>
				<integer>0</integer>
			</dict>
		</dict>
		<dict>
			<key>operation</key>
			<string>insert_track</string>
			<key>item</key>
			<dict>
				<key>title</key>
				<string>独家记忆(1)</string>
				<key>is_ringtone</key>
				<true/>
				<key>sort_name</key>
				<string>独家记忆(1)</string>
				<key>total_time_ms</key>
				<integer>307131</integer>
			</dict>
			<key>pid</key>
			<integer>8583372321666431737</integer>
			<key>ringtone_info</key>
			<dict>
				<key>guid</key>
				<string>771E43A65FD8FAF9</string>
			</dict>
		</dict>
		<dict>
			<key>pid</key>
			<integer>8375388374480384279</integer>
			<key>operation</key>
			<string>delete_track</string>
			<key>device_path</key>
			<string>/iTunes_Control/Ringtones/KFUT.m4r</string>
			<key>path</key>
			<string>/iTunes_Control/Ringtones/KFUT.m4r</string>
		</dict>
	</array>
</dict>
</plist>
*/
func mashalUpdateProto(syncNum int, inserts, dels []*TrackInfo) (plist []byte, err error) {
	oppl := operaitionsProto(syncNum)

	stProto := struct {
		Revision   int                      `plist:"revision"`
		Timestamp  time.Time                `plist:"timestamp"`
		Operations []map[string]interface{} `plist:"operations"`
	}{}
	_, err = mtunes.UnmashalPlist(oppl, &stProto)
	if err != nil {
		return
	}

	for _, delInfo := range dels {
		del := make(map[string]interface{}, 4)
		del["pid"] = delInfo.PID
		del["operation"] = "delete_track"
		del["device_path"] = delInfo.Path
		del["path"] = delInfo.Path
		stProto.Operations = append(stProto.Operations, del)
	}

	for _, insInfo := range inserts {
		ins := make(map[string]interface{}, 4)
		ins["pid"] = insInfo.PID
		ins["operation"] = "insert_track"

		item := make(map[string]interface{}, 4)
		item["title"] = insInfo.Name
		item["is_ringtone"] = true
		item["sort_name"] = insInfo.Name
		item["total_time_ms"] = insInfo.TotalTime
		ins["item"] = item

		ringtoneInfo := make(map[string]interface{}, 1)
		ringtoneInfo["guid"] = insInfo.GUID
		ins["ringtone_info"] = ringtoneInfo

		stProto.Operations = append(stProto.Operations, ins)
	}
	plist, err = mtunes.MashalPlist(stProto)
	return
}
