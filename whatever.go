package main

import (
	"bytes"
	"fmt"
	"io"

	"strings"

	"time"
)

func MergeReaders(r1 io.Reader, r2 io.Reader) (io.Reader, error) {
	var result bytes.Buffer

	// Buffers to hold single byte read from readers
	buf1 := make([]byte, 1)
	buf2 := make([]byte, 1)

	for {
		// Read one byte from r1
		n1, err1 := r1.Read(buf1)
		// Read one byte from r2
		n2, err2 := r2.Read(buf2)

		// If one of the readers reaches EOF or returns an error, stop reading
		if err1 != nil && err1 != io.EOF || err2 != nil && err2 != io.EOF {
			return nil, fmt.Errorf("error reading: %v, %v", err1, err2)
		}
		if n1 == 0 || n2 == 0 {
			break
		}

		// Append bytes from r1 and r2 alternately
		result.Write(buf1)
		result.Write(buf2)
	}

	return bytes.NewReader(result.Bytes()), nil
}

// Create a composition structure for operating systems. First,
// define a base class OS which consists of two fields, Name and IsFree.
// Two classes, LinuxOS and WindowsOS inherit from this base class.
// WindowsOS should contain a field EndOfSupport which is a timestamp (time.Time).
// LinuxOS should contain two boolean fields, YumBased and AptBased.

// You should define a method which will accept a distro name and detect
// if it is a Windows or a Linux system. No Windows versions
// are free, and all Linux versions are free. For Windows,
// it should print the end of mainstream support date and for
// Linux it should print if it is yum and if it is apt based.

type OS struct {
	Name   string
	IsFree bool
}

type LinuxOS struct {
	OS
	YumBased bool
	AptBased bool
}

func (l LinuxOS) String() string {
	if l.YumBased {
		return "yam"
	}
	return "apt"
}

type WindowsOS struct {
	OS
	EndOfSupport time.Time
}

func (w WindowsOS) String() string {
	return w.EndOfSupport.Local().String()
}
func DetectOS(name string) {
	windowsSupportDates := map[string]string{
		"Windows XP":    "2009-04-14",
		"Windows Vista": "2012-04-10",
		"Windows 7":     "2015-01-13",
		"Windows 8":     "2016-01-12",
		"Windows 8.1":   "2018-01-09",
	}

	linuxDistros := map[string]LinuxOS{
		"CentOS":   {OS: OS{Name: "Linux", IsFree: true}, YumBased: true, AptBased: false},
		"Debian":   {OS: OS{Name: "Linux", IsFree: true}, YumBased: false, AptBased: true},
		"Fedora":   {OS: OS{Name: "Linux", IsFree: true}, YumBased: true, AptBased: false},
		"Mint":     {OS: OS{Name: "Linux", IsFree: true}, YumBased: false, AptBased: true},
		"Raspbian": {OS: OS{Name: "Linux", IsFree: true}, YumBased: false, AptBased: true},
		"Ubuntu":   {OS: OS{Name: "Linux", IsFree: true}, YumBased: false, AptBased: true},
	}

	if date, ok := windowsSupportDates[name]; ok {
		endDate, _ := time.Parse("2006-01-02", date)
		fmt.Println("Windows")
		fmt.Println("false")
		fmt.Println(endDate)

	} else {

		for key, distro := range linuxDistros {
			upperKey := strings.ToUpper(key)
			upperName := strings.ToUpper(name)

			if strings.Contains(upperName, upperKey) {

				fmt.Println("Linux")
				fmt.Println("ture")
				fmt.Println(distro.YumBased)
				fmt.Println(distro.AptBased)
				break
			}

		}
	}
}

func merge(nums1 []int, m int, nums2 []int, n int) {
	i := 0
	j := 0
	k := 0

	for i < m && j < n {
		if nums1[i] < nums2[j] {
			nums1[k] = nums1[i]
			i++
		} else {
			nums1[k] = nums2[j]
			j++
		}
		k++
	}
	for i < m {
		nums1[k] = nums1[i]
		i++
		k++
	}

	for j < n {
		nums1[k] = nums2[j]
		j++
		k++
	}
}
