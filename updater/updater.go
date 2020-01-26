package updater

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"runtime"
	"syscall"

	"github.com/aerth/tgun"
)

var Version string

var DefaultEndpoint string

func Update1() error {

	arch := runtime.GOARCH
	osys := runtime.GOOS
	var suffixOptional string
	if osys == "windows" {
		suffixOptional = ".exe"
	}
	endpoint := fmt.Sprintf("%s/%s/%s/latest%s?have=%s", DefaultEndpoint, osys, arch, suffixOptional, Version)

	resp, err := new(tgun.Client).Get(endpoint)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("not 200, got %v", resp.StatusCode)
	}

	exeme, err := os.Executable()
	if err != nil {
		return err
	}
	if err := os.Remove(exeme); err != nil {
		return err
	}

	f, err := os.Create(exeme)
	if err != nil {
		return err
	}
	f.Close()
	f, err = os.OpenFile(exeme, os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		return err
	}

	if _, err := io.Copy(f, resp.Body); err != nil {
		return err
	}
	resp.Body.Close()

	f.Close()
	if err := os.Chmod(exeme, 0755); err != nil {
		return err
	}

	return nil

}
func Update2() {

	if runtime.GOOS == "windows" {
		exeme, err := os.Executable()
		if err != nil {
			log.Fatalln(err)
		}
		if err != nil {
			log.Fatalln(err)
		}
		cmd := exec.Command(exeme, os.Args[0:]...)
		if err := cmd.Start(); err != nil {
			log.Fatalln(err)
		}
		os.Exit(0)

		return
	}
	exeme, err := os.Executable()
	if err != nil {
		log.Fatalln(err)
	}
	if err := syscall.Exec(exeme, os.Args[0:], os.Environ()); err != nil {
		log.Fatalln(err)
	}
}
