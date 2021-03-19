package athservice

import (
	"fmt"
	"os"

	mtunes "github.com/yinyajiang/go-tunes"
	itunes "github.com/yinyajiang/go-tunes/itunesenv"
)

var (
	itunesVersion string
	deviceGrapa   = []byte{
		1, 1, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11,
		0x11, 0x11, 4, 0x40, 0xbc, 0x27, 0x85, 0xe0, 0xdb, 0xf1, 0x66, 0x36, 30, 7, 0x98, 10,
		0x5e, 0xa4, 0x8d, 0xba, 0x95, 0xb3, 0xb8, 0xea, 0x26, 0x5d, 0x62, 0xae, 0xfe, 0xa5, 0x1b, 0xb7,
		0xb1, 0x90, 0xe0, 0xb7, 0x71, 0x26, 0x29, 10, 0xd3, 0x9b, 0xb1, 0x3f, 0xec, 0xc0, 140, 0x25,
		0xa9, 0x56, 0x1c, 0x51, 0x7a, 0xc1, 30, 100, 0x90, 0x5d, 160, 0x29, 230, 0x1b, 0xdf, 0xd0,
		0xba, 0x22, 0xc3, 0x13}
)

/*
<dict>
	<key>SyncHostName</key>
	<string>EASEUS-YAJIANG</string>
	<key>SyncedDataclasses</key>
	<array/>
	<key>Version</key>
	<string>12.6.0.100</string>
	<key>libraryID</key>
	<string>12.6.0.100</string>
</dict>
*/
func requestProto() (plist []byte, err error) {
	p := struct {
		SyncHostName      string   `plist:"SyncHostName"`
		SyncedDataclasses []string `plist:"SyncedDataclasses"`
		Version           string   `plist:"Version"`
		LibraryID         string   `plist:"libraryID"`
	}{}

	hostName, _ := os.Hostname()
	if len(hostName) == 0 {
		hostName = "MMHost"
	}
	p.SyncHostName = hostName
	p.Version = getiTunesVersion()
	p.LibraryID = getiTunesVersion()

	plist, err = mtunes.MashalPlist(p)
	return
}

/*
<dict>
	<key>DataclassAnchors</key>
	<dict>
		<key>[Keybagk key]</key>
		<string>>[Keybagk value]</string>
	</dict>
	<key>Dataclasses</key>
	<array>
		<string>Keybag</string>
		<string>>[Keybagk key]</string>
	</array>
	<key>HostInfo</key>
	<dict>
		<key>Grappa</key>
		<data>
		AQERERERERERERERERERERERBEC8J4Xg2/FmNh4HmApepI26lbO46iZdYq7+
		pRu3sZDgt3EmKQrTm7E/7MCMJalWHFF6wR5kkF2gKeYb39C6IsMT
		</data>
		<key>SyncHostName</key>
		<string>EASEUS-YAJIANG</string>
		<key>SyncedDataclasses</key>
		<array>
			<string>Data</string>
			<array>
				<string>Keybag</string>
				<string>>[Keybagk key]</string>
			</array>
		</array>
		<key>Version</key>
		<string>12.6.0.100</string>
		<key>libraryID</key>
		<string>12.6.0.100</string>
	</dict>
</dict>
*/
func responseAllowedProto(keybag string, anchors []string) (plist []byte, err error) {
	rsp := struct {
		DataclassAnchors map[string]string `plist:"DataclassAnchors"`
		Dataclasses      []string          `plist:"Dataclasses"`
		HostInfo         struct {
			Grappa            []byte        `plist:"Grappa"`
			SyncHostName      string        `plist:"SyncHostName"`
			Version           string        `plist:"Version"`
			LibraryID         string        `plist:"libraryID"`
			SyncedDataclasses []interface{} `plist:"SyncedDataclasses"`
		} `plist:"HostInfo"`
	}{}
	keybagArray := make([]string, 0, 3)
	keybagArray = append(keybagArray, "Keybag")
	keybagArray = append(keybagArray, keybag)

	rsp.DataclassAnchors = make(map[string]string, 3)
	for _, val := range anchors {
		rsp.DataclassAnchors[keybag] = val
	}
	rsp.Dataclasses = keybagArray

	rsp.HostInfo.Grappa = deviceGrapa
	hostName, _ := os.Hostname()
	if len(hostName) == 0 {
		hostName = "MMHost"
	}
	rsp.HostInfo.SyncHostName = hostName
	rsp.HostInfo.Version = getiTunesVersion()
	rsp.HostInfo.LibraryID = getiTunesVersion()

	rsp.HostInfo.SyncedDataclasses = make([]interface{}, 0, 2)
	rsp.HostInfo.SyncedDataclasses = append(rsp.HostInfo.SyncedDataclasses, "Data")
	rsp.HostInfo.SyncedDataclasses = append(rsp.HostInfo.SyncedDataclasses, keybagArray)

	plist, err = mtunes.MashalPlist(rsp)
	return
}

/*
p1:
<dict>
	<key>Keybag</key>
	<integer>1</integer>
	<key>Media_test</key>
	<integer>1</integer>
</dict>

p2:
<dict>
	<key>Media_test</key>
	<string>0_syncNumTest</string>
</dict>
*/

func responseReadyProto(keybag, syncNum string) (pl1 []byte, pl2 []byte, err error) {
	p1 := map[string]int{
		"Keybag": 1,
		keybag:   1,
	}

	p2 := map[string]string{
		keybag: syncNum,
	}

	pl1, err = mtunes.MashalPlist(p1)
	if err != nil {
		return
	}

	pl2, err = mtunes.MashalPlist(p2)
	if err != nil {
		return
	}
	return
}

func unmarshalReadyProto(plmsg []byte, keybag string) (syncNum string, grappa []byte, err error) {
	stMsg := struct {
		Params struct {
			DataclassAnchors map[string]string `plist:"DataclassAnchors"`
			DeviceInfo       struct {
				Grappa []byte `plist:"Grappa"`
			} `plist:"DeviceInfo"`
		} `plist:"Params"`
	}{}

	_, err = mtunes.UnmashalPlist(plmsg, &stMsg)
	if err != nil {
		return
	}
	syncNum, ok := stMsg.Params.DataclassAnchors[keybag]
	if !ok {
		syncNum = stMsg.Params.DataclassAnchors["Media"]
	}
	if len(syncNum) == 0 {
		err = fmt.Errorf("Not get sync num")
		return
	}
	grappa = stMsg.Params.DeviceInfo.Grappa
	return
}

func unmarshalManifestProto(plmsg []byte, keybag string) (idArray []string, err error) {
	type AssetItem struct {
		AssetID string `plist:"AssetID"`
	}
	type AssetArray []AssetItem
	stMsg := struct {
		Params struct {
			AssetManifest map[string]AssetArray `plist:"AssetManifest"`
		} `plist:"Params"`
	}{}

	_, err = mtunes.UnmashalPlist(plmsg, &stMsg)
	if err != nil {
		return
	}
	if len(stMsg.Params.AssetManifest) == 0 {
		err = fmt.Errorf("Mainifest event is empty")
		return
	}

	assetArray, ok := stMsg.Params.AssetManifest[keybag]
	if !ok {
		assetArray = stMsg.Params.AssetManifest["Media"]
	}
	if len(assetArray) == 0 {
		err = fmt.Errorf("Mainifest event assetArray is empty")
		return
	}
	idArray = make([]string, 0, 2)
	for _, item := range assetArray {
		idArray = append(idArray, item.AssetID)
	}
	return
}

func getiTunesVersion() string {
	if len(itunesVersion) == 0 {
		itunesVersion = itunes.New().Version
		if len(itunesVersion) == 0 {
			itunesVersion = "12.10.2"
		}
	}
	return itunesVersion
}
