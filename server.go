package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"crypto/sha256"
	"encoding/base64"
	"golang.org/x/crypto/pbkdf2"
)

func main() {
	// HTTP 핸들러 등록
	http.HandleFunc("/", handleRoot)
	http.HandleFunc("/process", handleProcess)

	// 포트 설정 (환경변수 또는 기본값)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// 서버 시작
	fmt.Printf("서버가 http://localhost:%s 에서 시작되었습니다.\n", port)
	fmt.Printf("다음 주소들을 방문하세요:\n")
	fmt.Printf("- http://localhost:%s (입력 폼)\n", port)
	fmt.Printf("- http://localhost:%s/process?value1=10&value2=20 (값 처리)\n", port)

	log.Fatal(http.ListenAndServe(":"+port, nil))
}

// 메인 페이지 - HTML 폼 표시
func handleRoot(w http.ResponseWriter, r *http.Request) {
	html := `
	<!DOCTYPE html>
	<html>
	<head>
		<meta charset="UTF-8">
		<title>Get API Key for CATSA</title>
		<style>
			body {
				font-family: Arial, sans-serif;
				display: flex;
				justify-content: center;
				align-items: center;
				height: 100vh;
				margin: 0;
				background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
			}
			.container {
				background: white;
				padding: 40px;
				border-radius: 10px;
				box-shadow: 0 10px 25px rgba(0,0,0,0.2);
				width: 90%;
				max-width: 400px;
			}
			h1 {
				text-align: center;
				color: #333;
				margin-bottom: 30px;
			}
			.form-group {
				margin-bottom: 20px;
			}
			label {
				display: block;
				margin-bottom: 8px;
				color: #555;
				font-weight: bold;
			}
			input[type="text"], input[type="number"] {
				width: 100%;
				padding: 10px;
				border: 1px solid #ddd;
				border-radius: 5px;
				font-size: 14px;
				box-sizing: border-box;
			}
			input[type="text"]:focus, input[type="number"]:focus {
				outline: none;
				border-color: #667eea;
				box-shadow: 0 0 5px rgba(102,126,234,0.5);
			}
			button {
				width: 100%;
				padding: 12px;
				background: #667eea;
				color: white;
				border: none;
				border-radius: 5px;
				font-size: 16px;
				font-weight: bold;
				cursor: pointer;
				transition: background 0.3s;
			}
			button:hover {
				background: #764ba2;
			}
		</style>
	</head>
	<body>
		<div class="container">
			<h1>Get API Key for CATSA</h1>
			<form action="/process" method="GET">
				<div class="form-group">
					<label for="value1">PASSPHRASE:</label>
					<input type="text" id="value1" name="value1" placeholder="Enter passphrase" required>
				</div>
				<div class="form-group">
					<label for="value2">Serial Number:</label>
					<input type="text" id="value2" name="value2" placeholder="Enter Device Serial Number" required>
				</div>
				<button type="submit">Sumbmit</button>
			</form>
		</div>
	</body>
	</html>
	`
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, html)
}

var CATSASiteKey = []byte{
	0x67, 0x1F, 0x7F, 0x11, 0x5B, 0x3F, 0xB2, 0xB2,
	0x9D, 0xEB, 0xF2, 0xCA, 0x77, 0xAA, 0xE0, 0xAE,
	0x03, 0x0E, 0xCB, 0xF3, 0xBB, 0x1F, 0xCB, 0xB9,
	0x68, 0xEE, 0xC5, 0x30, 0x3A, 0xA2, 0x20, 0xAC,
}

func MakeKey(passphrase string, serialNumber string) (string, string) {

   siteKey:= CATSASiteKey

	apiKey := pbkdf2.Key(
		siteKey,
		[]byte(serialNumber),
		4096,
		32,
		sha256.New,
	)
	
	apiKeyBase64 := base64.StdEncoding.EncodeToString(apiKey)
	return apiKeyBase64, fmt.Sprintf("%x", siteKey)
}

// 값 처리 핸들러
func handleProcess(w http.ResponseWriter, r *http.Request) {
	// URL 쿼리 파라미터에서 값 받기
	pass := r.URL.Query().Get("value1")
	serial := r.URL.Query().Get("value2")

	// 값 검증
	if pass == "" || serial == "" {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, `
		<!DOCTYPE html>
		<html>
		<head>
			<meta charset="UTF-8">
			<title>Error</title>
			<style>
				body { font-family: Arial; text-align: center; margin-top: 50px; }
				.error { color: red; }
			</style>
		</head>
		<body>
			<h2 class="error">Error: Enter the input value!</h2>
			<a href="/">다시 시도</a>
		</body>
		</html>
		`)
		return
	}

	apiKey, siteKey := MakeKey(pass, serial)

/*
	response := map[string]string{
		"apiKey":  apiKey,
		"siteKey": siteKey,
	}
*/	
	// 숫자 변환 시도 (선택사항)
//	num1, err1 := strconv.ParseFloat(value1, 64)
//	num2, err2 := strconv.ParseFloat(value2, 64)

	// 응답 생성
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	html := `
	<!DOCTYPE html>
	<html>
	<head>
		<meta charset="UTF-8">
		<title>결과</title>
		<style>
			body {
				font-family: Arial, sans-serif;
				display: flex;
				justify-content: center;
				align-items: center;
				min-height: 100vh;
				margin: 0;
				background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
			}
			.container {
				background: white;
				padding: 40px;
				border-radius: 10px;
				box-shadow: 0 10px 25px rgba(0,0,0,0.2);
				width: 90%;
				max-width: 600px;
			}
			h1 {
				text-align: center;
				color: #333;
			}
			.result {
				background: #f5f5f5;
				padding: 20px;
				border-radius: 5px;
				margin-bottom: 20px;
			}
			.result-item {
				margin-bottom: 15px;
			}
			.label {
				font-weight: bold;
				color: #555;
			}
			.value {
				color: #333;
				font-size: 13px;
				margin-top: 5px;
				word-break: break-all;
				word-wrap: break-word;
				overflow-wrap: break-word;
				background: #f9f9f9;
				padding: 10px;
				border: 1px solid #ddd;
				border-radius: 3px;
				max-height: 120px;
				overflow-y: auto;
				font-family: 'Courier New', monospace;
				line-height: 1.5;
				white-space: pre-wrap;
			}
			.calculation {
				background: #e8f4f8;
				padding: 15px;
				border-radius: 5px;
				margin-bottom: 20px;
			}
			.back-btn {
				display: inline-block;
				width: 100%;
				padding: 12px;
				background: #667eea;
				color: white;
				text-align: center;
				border-radius: 5px;
				text-decoration: none;
				font-weight: bold;
				cursor: pointer;
				border: none;
				font-size: 16px;
				transition: background 0.3s;
			}
			.back-btn:hover {
				background: #764ba2;
			}
		</style>
	</head>
	<body>
		<div class="container">
			<h1>✓ Complete</h1>
			<div class="result">
				<div class="result-item">
					<div class="label">API Key:</div>
					<div class="value">` + apiKey + `</div>
				</div>
				<div class="result-item">
					<div class="label">Site Key:</div>
					<div class="value">` + siteKey + `</div>
				</div>
			</div>
	`

	// 숫자 계산이 가능한 경우 추가 정보 표시
	/*
	if err1 == nil && err2 == nil {
		sum := num1 + num2
		diff := num1 - num2
		product := num1 * num2
		var quotient string
		if num2 != 0 {
			quotient = fmt.Sprintf("%.2f", num1/num2)
		} else {
			quotient = "불가능 (0으로 나눔)"
		}

		html += `
			<div class="calculation">
				<div class="label">수치 계산 결과:</div>
				<div style="margin-top: 10px; line-height: 1.8;">
					<div>합계: <strong>` + fmt.Sprintf("%.2f", sum) + `</strong></div>
					<div>차이: <strong>` + fmt.Sprintf("%.2f", diff) + `</strong></div>
					<div>곱하기: <strong>` + fmt.Sprintf("%.2f", product) + `</strong></div>
					<div>나누기: <strong>` + quotient + `</strong></div>
				</div>
			</div>
		`
	}
	*/

	html += `
			<a href="/" class="back-btn">Go to Home</a>
		</div>
	</body>
	</html>
	`

	fmt.Fprint(w, html)

	// 콘솔에도 출력
	log.Printf("값1: %s, 값2: %s\n", apiKey, siteKey)
}
