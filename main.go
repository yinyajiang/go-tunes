package main

import (
	"fmt"

	mtunes "github.com/yinyajiang/go-tunes/mTunes"
)

func main() {
	// dev := mtunes.OnceEventLoopForWaitDevice(context.Background(), "651280804511790")
	// if dev == nil {
	// 	fmt.Println("err")
	// 	return
	// }
	// service, _ := fileservice.New(dev)
	// service.RenameAndMove("/HelloWord/ww.txt", "/HelloWord/1/ww.txt")

	data := `<?xml version="1.0" encoding="UTF-8"?>
	<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
	<plist version="1.0">
	<dict>
		<key>Command</key>
		<string>AssetManifest</string>
		<key>Id</key>
		<integer>447</integer>
		<key>Params</key>
		<dict>
			<key>AssetManifest</key>
			<dict>
				<key>Media</key>
				<array>
					<dict>
						<key>AssetID</key>
						<string>8048284188921304141</string>
						<key>AssetType</key>
						<string>Movie</string>
						<key>IsDownload</key>
						<true/>
						<key>Variant</key>
						<dict>
							<key>AssetParts</key>
							<integer>2</integer>
						</dict>
					</dict>
					<dict>
						<key>AssetID</key>
						<string>506887312052933397</string>
						<key>AssetType</key>
						<string>Audiobook</string>
						<key>IsDownload</key>
						<true/>
						<key>Variant</key>
						<dict>
							<key>AssetParts</key>
							<integer>2</integer>
						</dict>
					</dict>
					<dict>
						<key>AssetID</key>
						<string>5940101454307634274</string>
						<key>AssetType</key>
						<string>Music</string>
						<key>IsDownload</key>
						<true/>
						<key>Variant</key>
						<dict>
							<key>AssetParts</key>
							<integer>2</integer>
						</dict>
					</dict>
					<dict>
						<key>AssetID</key>
						<string>5681164399141447303</string>
						<key>AssetType</key>
						<string>Movie</string>
						<key>IsDownload</key>
						<true/>
						<key>Variant</key>
						<dict>
							<key>AssetParts</key>
							<integer>2</integer>
						</dict>
					</dict>
				</array>
				<key>Photo</key>
				<array>
					<dict>
						<key>AssetID</key>
						<string>8048284188921304141</string>
						<key>AssetType</key>
						<string>Movie</string>
						<key>IsDownload</key>
						<true/>
						<key>Variant</key>
						<dict>
							<key>AssetParts</key>
							<integer>2</integer>
						</dict>
					</dict>
					<dict>
						<key>AssetID</key>
						<string>506887312052933397</string>
						<key>AssetType</key>
						<string>Audiobook</string>
						<key>IsDownload</key>
						<true/>
						<key>Variant</key>
						<dict>
							<key>AssetParts</key>
							<integer>2</integer>
						</dict>
					</dict>
					<dict>
						<key>AssetID</key>
						<string>5940101454307634274</string>
						<key>AssetType</key>
						<string>Music</string>
						<key>IsDownload</key>
						<true/>
						<key>Variant</key>
						<dict>
							<key>AssetParts</key>
							<integer>2</integer>
						</dict>
					</dict>
					<dict>
						<key>AssetID</key>
						<string>5681164399141447303</string>
						<key>AssetType</key>
						<string>Movie</string>
						<key>IsDownload</key>
						<true/>
						<key>Variant</key>
						<dict>
							<key>AssetParts</key>
							<integer>2</integer>
						</dict>
					</dict>
				</array>
			</dict>
		</dict>
		<key>Session</key>
		<integer>0</integer>
		<key>Type</key>
		<integer>0</integer>
	</dict>
	</plist>
	`

	type AssetItem struct {
		AssetID   string `plist:"AssetID"`
		AssetType string `plist:"AssetType"`
	}

	type AssetArray []AssetItem

	stMsg := struct {
		Params struct {
			AssetManifest map[string]AssetArray `plist:"AssetManifest"`
		} `plist:"Params"`
	}{}

	mtunes.UnmashalPlist([]byte(data), &stMsg)

	for keybag, arr := range stMsg.Params.AssetManifest {
		fmt.Println(keybag)
		for _, item := range arr {
			fmt.Println(item)
		}
	}
}
