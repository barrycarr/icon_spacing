package main

import (
	"bufio"
	"flag"
	"fmt"
	"golang.org/x/sys/windows/registry"
	"log"
	"os"
	"strconv"
	"strings"
)

type distanceOption int
type options struct {
	distance distanceOption
	update   bool
}

const (
	WIDE distanceOption = iota + 1
	MEDIUM
	NARROW
)

var (
	distancesMap = map[string]distanceOption{
		"w": WIDE,
		"m": MEDIUM,
		"n": NARROW,
	}
)

func isWindowsEleven() (bool, error) {
	const Windows11BuildNo = 22000

	k, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows NT\CurrentVersion`, registry.QUERY_VALUE)
	if err != nil {
		return false, err
	}
	defer func(k registry.Key) {
		if err := k.Close(); err != nil {
			log.Fatal(err)
		}
	}(k)

	s, _, err := k.GetStringValue("CurrentBuildNumber")
	if err != nil {
		return false, err
	}
	buildNo, err := strconv.Atoi(s)
	if err != nil {
		return false, err
	}
	return buildNo >= Windows11BuildNo, nil
}

func getOptions() options {
	distanceOption := flag.String("distance", "WIDE", "Distance between desktop icons. Valid values are: wide, medium or narrow. Alternatively: w, m or n")
	update := flag.Bool("update", false, "Automatically update the distance value")
	flag.Parse()
	result, _ := distancesMap[strings.ToLower(*distanceOption)[0:1]]

	return options{distance: result, update: *update}
}

func setKeys(k registry.Key, value string) error {
	const (
		iconSpacing         = "IconSpacing"
		iconVerticalSpacing = "IconVerticalSpacing"
	)

	fmt.Println("Updating registry...")
	if err := k.SetStringValue(iconSpacing, value); err != nil {
		return err
	}

	if err := k.SetStringValue(iconVerticalSpacing, value); err != nil {
		return err
	}

	return nil
}

func setIconDistance(distance distanceOption) error {
	const (
		wide   = "-2056"
		narrow = "-1128"
		medium = "-1592"
	)

	k, err := registry.OpenKey(registry.CURRENT_USER, `Control Panel\Desktop\WindowMetrics\`, registry.SET_VALUE)
	if err != nil {
		return err
	}

	switch distance {
	case WIDE:
		if err := setKeys(k, wide); err != nil {
			return err
		}
	case MEDIUM:
		if err := setKeys(k, medium); err != nil {
			return err
		}
	case NARROW:
		if err := setKeys(k, narrow); err != nil {
			return err
		}
	}

	defer func(k registry.Key) {
		if err := k.Close(); err != nil {
			log.Fatal(err)
		}
	}(k)

	return nil
}

func userAffirmed(alreadyAgreed bool) (bool, error) {
	const (
		y = 121
		Y = 89
	)

	if alreadyAgreed {
		return alreadyAgreed, nil
	}

	fmt.Println("Proceeding will make changes to your registry. ")
	fmt.Println("Are you sure? (enter `y` or `Y` to confirm, any other key to cancel)")
	fmt.Print("> ")

	consoleReader := bufio.NewReaderSize(os.Stdin, 1)
	input, err := consoleReader.ReadByte()
	if err != nil {
		return false, err
	}
	return input == y || input == Y, nil
}

func main() {
	fmt.Println("Windows 11 desktop icon spacing utility")
	isWin11, err := isWindowsEleven()
	if err != nil {
		fmt.Println("Error getting Windows version number: ", err)
		os.Exit(-12)
	}
	if !isWin11 {
		fmt.Println("This utility is for Windows 11 or greater.")
		os.Exit(-11)
	}

	opts := getOptions()
	if opts.distance == 0 {
		fmt.Println("Invalid distance option received.")
		os.Exit(-13)
	}

	if yn, err := userAffirmed(opts.update); yn && err == nil {
		e := setIconDistance(opts.distance)
		if e != nil {
			fmt.Println("Couldn't update registry", e)
			os.Exit(-14)
		}
	} else if !yn {
		fmt.Println("User declined update. No changes have been made to your system.")
		os.Exit(1)
	} else {
		fmt.Println("Couldn't get response from user", err)
		os.Exit(-15)
	}

	fmt.Println("Done!")
	fmt.Println("You may need to reboot for changes to take effect.")
	os.Exit(0)
}
