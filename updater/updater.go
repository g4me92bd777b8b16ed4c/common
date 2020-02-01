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
var Update1 = Stage1

func Rebuild() error {
	b, err := exec.Command("make").CombinedOutput()
	if err != nil {
		return fmt.Errorf("build error: %v\n%s", err, string(b))
	}
	log.Println("REBUILD:", string(b))
	return nil
}

func Stage1() error {
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

var Update2 = Stage2

func Stage2() {
	// TODO: i dont know if this works lol
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
