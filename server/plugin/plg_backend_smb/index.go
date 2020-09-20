package plg_backend_smb

import (
	"io"
	"net"
	"os"
	"strings"

	"github.com/hirochachacha/go-smb2"
	. "github.com/mickael-kerjean/filestash/server/common"
)

var SmbCache AppCache

func init() {
	Backend.Register("smb", Smb{})
}

type Smb struct {
	params *SmbParams
}

type SmbParams struct {
	server   string
	shared   string
	username string
	password string
}

func (smb Smb) Init(params map[string]string, app *App) (IBackend, error) {
	backend := Smb{
		params: &SmbParams{
			params["server"],
			params["shared"],
			params["username"],
			params["password"],
		},
	}
	return backend, nil
}

func (smb Smb) LoginForm() Form {
	return Form{
		Elmnts: []FormElement{
			FormElement{
				Name:  "type",
				Type:  "hidden",
				Value: "smb",
			},
			FormElement{
				Name:        "server",
				Type:        "text",
				Placeholder: "Server",
			},
			FormElement{
				Name:        "shared",
				Type:        "text",
				Placeholder: "SharedName",
			},
			FormElement{
				Name:        "username",
				Type:        "text",
				Placeholder: "Username",
			},
			FormElement{
				Name:        "password",
				Type:        "password",
				Placeholder: "Password",
			},
		},
	}
}

func (smb Smb) Ls(path string) ([]os.FileInfo, error) {

	conn, err := net.Dial("tcp", smb.params.server+":445")
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	d := &smb2.Dialer{
		Initiator: &smb2.NTLMInitiator{
			User:     smb.params.username,
			Password: smb.params.password,
		},
	}

	session, err := d.Dial(conn)
	if err != nil {
		return nil, err
	}
	defer session.Logoff()

	fs, err := session.Mount(smb.params.shared)
	if err != nil {
		return nil, err
	}
	defer fs.Umount()

	osPath := path
	if osPath == "/" {
		osPath = ""
	}

	fis, err := fs.ReadDir(encodeUNC(osPath))
	if err != nil {
		return nil, err
	}

	return fis, nil
}

func (smb Smb) Cat(path string) (io.ReadCloser, error) {
	return nil, nil
}

func (smb Smb) Mkdir(path string) error {
	return nil
}

func (smb Smb) Rm(path string) error {

	conn, err := net.Dial("tcp", smb.params.server+":445")
	if err != nil {
		return err
	}
	defer conn.Close()

	d := &smb2.Dialer{
		Initiator: &smb2.NTLMInitiator{
			User:     smb.params.username,
			Password: smb.params.password,
		},
	}

	session, err := d.Dial(conn)
	if err != nil {
		return err
	}
	defer session.Logoff()

	fs, err := session.Mount(smb.params.shared)
	if err != nil {
		return err
	}
	defer fs.Umount()

	osPath := path
	if osPath == "/" {
		osPath = ""
	}

	err = fs.RemoveAll(encodeUNC(osPath))
	return err
}

func (smb Smb) Mv(from string, to string) error {
	return nil
}

func (smb Smb) Touch(path string) error {
	return nil
}

func (smb Smb) Save(path string, file io.Reader) error {
	return nil
}

func (smb Smb) Stat(path string) (os.FileInfo, error) {
	return nil, nil
}

func (smb Smb) Close() error {
	return nil
}

func encodeUNC(path string) string {
	osPath := strings.Replace(path, "/", "", 1)
	return strings.Replace(osPath, "/", "\\", -1)
}
