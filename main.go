package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/textproto"
	"os"
	"strings"
)

var quoteEscaper = strings.NewReplacer("\\", "\\\\", `"`, "\\\"")

func escapeQuotes(s string) string {
	return quoteEscaper.Replace(s)
}

type File struct {
	Path    string
	Type    string
	Content []byte
}

type Configuration struct {
	UserDataFiles [][]string `json:"cloud_init_parts"`
}

func buildUserData(files []File) (string, error) {
	// Let's build the launch config userdata
	w := &bytes.Buffer{}
	mimeWriter := multipart.NewWriter(w)
	// Craft a header for our mime type
	fmt.Fprintf(w, "Content-Type: multipart/mixed; boundary=\"%s\"\r\n\r\n", mimeWriter.Boundary())
	for _, file := range files {
		fileParts := strings.Split(file.Path, "/")
		h := make(textproto.MIMEHeader)
		h.Set("Content-Type", fmt.Sprintf("%s; charset=\"us-ascii\"", file.Type))
		h.Set("MIME-Version", "1.0")
		h.Set("Content-Transfer-Encoding", "7bit")
		h.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, escapeQuotes(fileParts[len(fileParts)-1])))

		part, err := mimeWriter.CreatePart(h)
		if err != nil {
			return "", fmt.Errorf("Unable to create mime type: %v", err)
		}
		r := bytes.NewReader(file.Content)
		_, err = io.Copy(part, r)
		if err != nil {
			return "", fmt.Errorf("Unable to copy the data part: %v", err)
		}
	}
	_ = mimeWriter.Close()

	return w.String(), nil
}

func main() {

	var (
		configFile string
		encode     bool
	)

	flag.StringVar(&configFile, "config", "<file>", "Config file containing paths and type of the userdata")
	flag.BoolVar(&encode, "encode", false, "Base64 encode the userdata")

	flag.Parse()

	//Do some validation
	if configFile == "<file>" {
		fmt.Println("No config supplied")
		os.Exit(1)
	}

	// Parse configuration
	config, err := ioutil.ReadFile(configFile)
	if err != nil {
		fmt.Printf("Unable to open config file %s: %v\n", configFile, err)
		os.Exit(1)
	}

	configuration := &Configuration{}
	json.Unmarshal(config, &configuration)

	files := make([]File, 0)

	// Construct the files map

	for _, file := range configuration.UserDataFiles {
		b, err := ioutil.ReadFile(file[0])
		if err != nil {
			fmt.Printf("Unable to open user data file: %s\n", file[0])
			os.Exit(1)
		}
		f := File{
			Content: b,
			Path:    file[0],
			Type:    file[1],
		}
		files = append(files, f)
	}

	userdata, err := buildUserData(files)
	if err != nil {
		fmt.Printf("Error building userdata: %v", err)
		os.Exit(1)
	}

	if encode {
		userdata = base64.StdEncoding.EncodeToString([]byte(userdata))
	}

	fmt.Println(userdata)

	return

}
