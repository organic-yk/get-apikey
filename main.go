package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"encoding/json"
	"crypto/sha256"
	"encoding/base64"
	"golang.org/x/crypto/pbkdf2"
)

// HTML 템플릿
var tmpl = template.Must(template.New("form").Parse(`
<!DOCTYPE html>
<html>
<head>
    <title>Input Form</title>
</head>
<body>
    <h2>Enter Information</h2>
    <form method="POST" action="/submit">
        <label>PASSPHRASE:</label><br>
        <input type="password" name="passphrase"><br><br>

        <label>Serial Number:</label><br>
        <input type="text" name="serial"><br><br>

        <input type="submit" value="Submit">
    </form>
</body>
</html>
`))

// 결과 템플릿
var resultTmpl = template.Must(template.New("result").Parse(`
<!DOCTYPE html>
<html>
<head>
    <title>Result</title>
</head>
<body>
    <h2>Result</h2>
    <p><b>APIKey:</b> {{.Pass}}</p>
	<p><b>SITEKey:</b> {{.Serial}}</p>
</body>
</html>
`))

// 메인 페이지
func formHandler(w http.ResponseWriter, r *http.Request) {
	tmpl.Execute(w, nil)
}

// 제출 처리
func submitHandler(w http.ResponseWriter, r *http.Request) {
	// 🔥 무조건 먼저 CORS 헤더 설정
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")

	// 🔥 preflight 요청 처리 (핵심)
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	// POST만 허용
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// form 파싱
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	pass := r.FormValue("passphrase")
	serial := r.FormValue("serial")

	apiKey, siteKey := MakeKey(pass, serial)

	response := map[string]string{
		"apiKey":  apiKey,
		"siteKey": siteKey,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
/*
func MakeKey(passphrase string, serialNumber string) ([]byte, []byte) {
//	key.mutex.Lock()
//	defer key.mutex.Unlock()

	cPassPhrase := C.CString(passphrase)
	defer C.free(unsafe.Pointer(cPassPhrase))

	siteKey := make([]byte, 32)
	C.kdf(cPassPhrase, (*C.uchar)(unsafe.Pointer(&siteKey[0])))

	apiKey := pbkdf2.Key(siteKey, []byte(serialNumber), 4096, 32, sha256.New)

	return siteKey, apiKey
}
*/

var CATSASiteKey = []byte{
	0x67, 0x1F, 0x7F, 0x11, 0x5B, 0x3F, 0xB2, 0xB2,
	0x9D, 0xEB, 0xF2, 0xCA, 0x77, 0xAA, 0xE0, 0xAE,
	0x03, 0x0E, 0xCB, 0xF3, 0xBB, 0x1F, 0xCB, 0xB9,
	0x68, 0xEE, 0xC5, 0x30, 0x3A, 0xA2, 0x20, 0xAC,
}

func MakeKey(passphrase string, serialNumber string) (string, string) {

//	salt := sha256.Sum256([]byte(serialNumber))
/*
	siteKey := pbkdf2.Key(
		[]byte(passphrase),
		salt[:],
		4096,
		32,
		sha256.New,
	)
*/
   siteKey:= CATSASiteKey

	apiKey := pbkdf2.Key(
		siteKey,
		[]byte(serialNumber),
		4096,
		32,
		sha256.New,
	)
	
	apiKeyBase64 := base64.StdEncoding.EncodeToString(apiKey)
/*	
	apiKey := pbkdf2.Key(
		siteKey,
		[]byte(serialNumber),
		4096,
		32,
		sha256.New,
	)
*/
	//return siteKey, apiKey
	return apiKeyBase64, fmt.Sprintf("%x", siteKey)
}

func enableCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
}

func main() {
	http.HandleFunc("/submit", submitHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Println("Server started at :" + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
