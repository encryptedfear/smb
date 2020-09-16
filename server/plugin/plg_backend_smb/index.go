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
	SmbCache = NewAppCache()

	SmbCache.OnEvict(func(key string, value interface{}) {
		s := value.(*Smb)
		s.Close()
	})
}

type Smb struct {
	SmbClient *smb2.Share
}

func (smb Smb) Init(params map[string]string, app *App) (IBackend, error) {
	p := struct {
		server   string
		shared   string
		username string
		password string
	}{
		params["server"],
		params["shared"],
		params["username"],
		params["password"],
	}

	if c := SmbCache.Get(params); c != nil {
		d := c.(*Smb)
		return d, nil
	}

	// conn, err := net.Dial("tcp", p.server + ":445")
	conn, err := net.Dial("tcp", "192.168.0.14:445")
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	d := &smb2.Dialer{
		Initiator: &smb2.NTLMInitiator{
			User:     "sephgallo",
			Password: "Crimsonblue_@!23",
		},
	}

	s, err := d.Dial(conn)
	if err != nil {
		panic(err)
	}
	defer s.Logoff()

	fs, err := s.Mount("Shared")
	if err != nil {
		panic(err)
	}
	defer fs.Umount()

	smb.SmbClient = fs
	SmbCache.Set(params, &smb)
	return &smb, nil

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

	dir, err := smb.SmbClient.Open("")
	if err != nil {
		return nil, err
	}

	fis, err := dir.Readdir(0)
	if err != nil {
		return nil, err
	}

	return fis, nil
}

func (smb Smb) Cat(path string) (io.ReadCloser, error) {
	remoteFile, err := smb.SmbClient.Open(path)
	if err != nil {
		return nil, err
	}

	return remoteFile, nil
}

func (smb Smb) Mkdir(path string) error {
	err := smb.SmbClient.Mkdir(path, os.FileMode(0777))
	return err
}

func (smb Smb) Rm(path string) error {
	err := smb.SmbClient.RemoveAll(path)
	return err
}

func (smb Smb) Mv(from string, to string) error {
	err := smb.SmbClient.Rename(from, to)
	return err
}

func (smb Smb) Touch(path string) error {
	file, err := smb.SmbClient.Create(path)
	if err != nil {
		return err
	}
	_, err = file.ReadFrom(strings.NewReader(""))
	return err
}

func (smb Smb) Save(path string, file io.Reader) error {
	remoteFile, err := smb.SmbClient.Create(path)
	if err != nil {
		return err
	}
	_, err = io.Copy(remoteFile, file)
	remoteFile.Close()
	return err
}

func (smb Smb) Stat(path string) (os.FileInfo, error) {
	f, err := smb.SmbClient.Stat(path)
	return f, err
}

func (smb Smb) Close() error {
	err0 := smb.SmbClient.Umount()
	return err0
}
