module github.com/Speakerkfm/iso

go 1.15

require (
	github.com/Speakerkfm/iso/pkg/proto v0.0.0-20211230104626-9305ab470077
	github.com/emicklei/proto v1.9.1
	github.com/spf13/cobra v1.3.0
	gopkg.in/yaml.v2 v2.4.0
)

replace github.com/Speakerkfm/iso/pkg/proto v0.0.0-20211230104626-9305ab470077 => github.com/Speakerkfm/iso/example/iso v0.0.0-20211230105029-dab7a57ee3b4
