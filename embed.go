package filotto

import _ "embed"

var (
	//go:embed web/filotto/app.html
	AppHTML []byte

	//go:embed gen/filotto/channel.js
	ChannelJS []byte
)
