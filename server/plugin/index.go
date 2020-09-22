package plugin

import (
	. "github.com/mickael-kerjean/filestash/server/common"
	_ "github.com/mickael-kerjean/filestash/server/plugin/plg_backend_ldap"
	_ "github.com/mickael-kerjean/filestash/server/plugin/plg_backend_smb"
	_ "github.com/mickael-kerjean/filestash/server/plugin/plg_editor_onlyoffice"
	_ "github.com/mickael-kerjean/filestash/server/plugin/plg_handler_console"
	_ "github.com/mickael-kerjean/filestash/server/plugin/plg_handler_syncthing"
	_ "github.com/mickael-kerjean/filestash/server/plugin/plg_security_scanner"
	_ "github.com/mickael-kerjean/filestash/server/plugin/plg_security_svg"
	_ "github.com/mickael-kerjean/filestash/server/plugin/plg_starter_tor"
	_ "github.com/mickael-kerjean/filestash/server/plugin/plg_starter_tunnel"
	_ "github.com/mickael-kerjean/filestash/server/plugin/plg_video_transcoder"
)

func init() {
	Log.Debug("Plugin loader")
}
