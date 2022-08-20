package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"time"
)

func main() {
	fmt.Println("Server Listening")
	err := http.ListenAndServe(":8000", http.HandlerFunc(handler))
	if err != nil {
		fmt.Println("Failed to ListenAndServe : ", err)
	}
}

func handler(rw http.ResponseWriter, req *http.Request) {
	// fmt.Println("Method : ", req.Method)
	// fmt.Println("URL : ", req.URL)
	// fmt.Println("Header : ", req.Header)

	body, _ := io.ReadAll(req.Body)
	defer req.Body.Close()
	fmt.Println("Body : ", string(body))

	resp := new(TargetRequest)
	_ = json.Unmarshal(body, resp)

	switch resp.Target {
	case "Book":
		resp := new(BookRequest)
		_ = json.Unmarshal(body, resp)

		err, output := RunScript("alice", "create-book",
			resp.BookId, resp.Title, resp.Synopsis, resp.CreatedAt)
		if err != nil {
			fmt.Println("error: ", err)
		}
		rw.Write([]byte(output))
	default:
		rw.Write([]byte("Wrong Target"))
	}
}

type TargetRequest struct {
	Target string `json:"target"`
}

type BookRequest struct {
	BookId    string `json:"bookId"`
	Title     string `json:"title"`
	Synopsis  string `json:"synopsis"`
	CreatedAt string `json:"createdAt"`
	Account   string `json:"accountName"`
}

func RunScript(from string, message string, arguments ...string) (error, string) {
	args := makeArgs(arguments)
	script := []byte("storyblockd tx storyblock " + message + args + " --from " + from + " -y")
	fileName := "./run_" + StringWithCharset(10) + ".sh"
	if err := ioutil.WriteFile(fileName, script, 0644); err != nil {
		_ = os.Remove(fileName)
		return err, ""
	}
	cmd, err := exec.Command("/bin/sh", fileName).Output()
	if err != nil {
		_ = os.Remove(fileName)
		return err, ""
	}
	output := string(cmd)
	_ = os.Remove(fileName)
	return nil, output
}

func makeArgs(args []string) string {
	result := " "
	for _, str := range args {
		result += str + " "
	}
	return result
}

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

func StringWithCharset(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz" +
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}
