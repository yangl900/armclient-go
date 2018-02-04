package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func prettyJSON(buffer []byte) string {
	var prettyJSON string
	if len(buffer) > 0 {
		var jsonBuffer bytes.Buffer
		error := json.Indent(&jsonBuffer, buffer, "", "  ")
		if error != nil {
			return string(buffer)
		}
		prettyJSON = jsonBuffer.String()
	} else {
		prettyJSON = ""
	}

	return prettyJSON
}

func responseDetail(response *http.Response) string {
	var buffer bytes.Buffer
	fmt.Fprintf(&buffer, "%s\t\t%s\n", response.Request.Method, response.Request.URL.String())
	fmt.Fprintf(&buffer, "Status Code:\t%d", response.StatusCode)

	return buffer.String()
}
