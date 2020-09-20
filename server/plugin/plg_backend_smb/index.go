package plg_backend_smb

import (
	"fmt"
	"io"
	"net"
	"os"

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
	SmbClient *smb2.Client
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

	fmt.Printf("This is stupid")

	conn, err := net.Dial("tcp", p.server+":445")
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	fmt.Printf("This is stupid 2")

	d := &smb2.Dialer{
		Initiator: &smb2.NTLMInitiator{
			User:     p.username,
			Password: p.password,
		},
	}

	fmt.Printf("This is stupid 3")

	client, err := d.Dial(conn)
	if err != nil {
		return nil, err
	}
	defer client.Logoff()

	fmt.Printf("This is stupid 4")

	smb.SmbClient = client

	SmbCache.Set(params, &smb)
	fmt.Printf("This is stupid 7")
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
	fmt.Printf("This is stupid 8")
	rfs, err := smb.SmbClient.Mount("Shared")
	if err != nil {
		return nil, err
	}

	dir, err := rfs.Open("")
	fmt.Printf("This is stupid 9")
	if err != nil {
		return nil, err
	}
	fmt.Printf("This is stupid 10")
	fis, err := dir.Readdir(0)
	if err != nil {
		return nil, err
	}

	return fis, nil
}

func (smb Smb) Cat(path string) (io.ReadCloser, error) {
	rfs, err := smb.SmbClient.Mount("Shared")
	if err != nil {
		return nil, err
	}
	remoteFile, err := rfs.Open(path)
	if err != nil {
		return nil, err
	}

	return remoteFile, nil
}

func (smb Smb) Mkdir(path string) error {
	return nil
}

func (smb Smb) Rm(path string) error {
	return nil
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
	return nil
}

func (smb Smb) Close() error {
	return nil
}
