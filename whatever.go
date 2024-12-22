package main

import (
	"bytes"
	"fmt"
	"io"

	"encoding/json"

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

// Manager represents a company manager with specific fields
type Manager struct {
	FullName       string `json:"full_name,omitempty"`
	Position       string `json:"position,omitempty"`
	Age            int32  `json:"age,omitempty"`
	YearsInCompany int32  `json:"years_in_company,omitempty"`
}

// EncodeManager encodes a Manager object to JSON and returns an io.Reader
func EncodeManager(manager *Manager) (io.Reader, error) {
	// Check if manager is nil
	if manager == nil {
		return nil, nil
	}

	// Marshal the manager object to JSON
	jsonData, err := json.Marshal(manager)
	if err != nil {
		return nil, err
	}

	// Create and return a reader containing the JSON data
	return bytes.NewReader(jsonData), nil
}

func asdf() {
	// Example usage
	manager := &Manager{
		FullName:       "Chris Smith",
		Position:       "CISO",
		Age:            32,
		YearsInCompany: 5,
	}

	reader, err := EncodeManager(manager)
	if err != nil {
		panic(err)
	}

	// Read and print the JSON for demonstration
	data, _ := io.ReadAll(reader)
	println(string(data))
}

func camelCaseTo_snack_case(str string) string {
	var result []byte

	for i, b := range []byte(str) {

		upperCased := strings.ToUpper(string(b))[0]
		if b != upperCased {
			result = append(result, b)
		} else {
			if i == 0 {
				result = append(result, strings.ToLower(string(b))[0])
			} else {
				result = append(result, '_', strings.ToLower(string(b))[0])
			}
		}
	}

	return string(result)
}

// Exercise 3: Implement a function that finds the longest increasing subsequence
// Example: [10, 9, 2, 5, 3, 7, 101, 18] -> [2, 5, 7, 101]
// func longestIncreasingSubsequence(nums []int) []int {
// 	var arr [][]int
// 	i := 0
// 	arr = append(arr, []int{nums[i]})
// 	i++
// 	for i < len(nums) {
// 		nextNum := nums[i]
// 		isInserted := false
// 		for j, a := range arr {
// 			lastNum := a[len(a)-1]
// 			if lastNum < nextNum {
// 				arr[j] = append(a, nextNum)
// 				isInserted = true
// 			}
// 		}
// 		if !isInserted {
// 			arr = append(arr, []int{nextNum})
// 		}
// 		i++
// 	}

// 	var maxArr []int
// 	maxLen := 0
// 	for _, a := range arr {
// 		fmt.Println(a)
// 		if maxLen < len(a) {
// 			maxLen = len(a)
// 			maxArr = a
// 		}
// 	}
// 	return maxArr
// }

func longestIncreasingSubsequence(nums []int) []int {
	if len(nums) == 0 {
		return []int{}
	}

	// Initialize DP arrays
	dp := make([]int, len(nums))   // dp[i] will hold the length of the LIS ending at index i
	prev := make([]int, len(nums)) // prev[i] will hold the index of the previous element in the LIS
	maxLength := 1                 // Length of the maximum LIS
	maxIndex := 0                  // Index of the last element in the LIS

	for i := range nums {
		dp[i] = 1    // Every element is a subsequence of length 1
		prev[i] = -1 // No previous element by default
		for j := 0; j < i; j++ {
			if nums[j] < nums[i] && dp[j]+1 > dp[i] {
				dp[i] = dp[j] + 1
				prev[i] = j
			}
		}
		if dp[i] > maxLength {
			maxLength = dp[i]
			maxIndex = i
		}
	}

	// Reconstruct the LIS
	lis := make([]int, 0, maxLength)
	for maxIndex != -1 {
		lis = append([]int{nums[maxIndex]}, lis...)
		maxIndex = prev[maxIndex]
	}

	return lis
}

// Exercise 4: Write a function that rotates a slice by k positions
// Example: ([1,2,3,4,5], k=2) -> [4,5,1,2,3]
func rotateSlice(slice []int, k int) []int {
	arr := make([]int, len(slice))
	for i := 0; i < len(slice); i++ {
		arr[(i+k)%len(slice)] = slice[i]
	}
	return arr
}
