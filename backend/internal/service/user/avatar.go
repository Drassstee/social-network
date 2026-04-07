package user

import (
	"io"
	"net/http"
	"os"
)

func DownloadFile(filepath, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}

	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

// --------------------------------------------------------------------|

func (us *UserService) GetAvatar() {

}

// --------------------------------------------------------------------|

func (us *UserService) SetAvatar() {

}
