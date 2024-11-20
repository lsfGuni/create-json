package main

import (
    "crypto/sha256"
    "encoding/hex"
    "encoding/json"
    "fmt"
    "net/http"
    "strconv"
    "text/template"
)

// Hash 구조체는 SEQ와 그에 해당하는 Hash 값을 저장합니다.
type Hash struct {
    Seq      int64  `json:"seq"`
    HashCode string `json:"hashcode"`
}

// 해시 생성 함수
func generateHash(seq int64) string {
    hash := sha256.New()
    hash.Write([]byte(fmt.Sprintf("seq%d", seq)))
    return hex.EncodeToString(hash.Sum(nil))
}

// 결과 JSON 생성 함수
func generateHashResults(startSeq, endSeq int64) []Hash {
    var results []Hash
    for seq := startSeq; seq <= endSeq; seq++ {
        hashCode := generateHash(seq)
        results = append(results, Hash{Seq: seq, HashCode: hashCode})
    }
    return results
}

// 메인 핸들러
func mainPage(w http.ResponseWriter, r *http.Request) {
    tmpl, err := template.ParseFiles("templates/index.html")
    if err != nil {
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }
    tmpl.Execute(w, nil)
}

// 해시 결과 생성 핸들러
func generateHashHandler(w http.ResponseWriter, r *http.Request) {
    // 쿼리 파라미터에서 startSeq와 endSeq 값을 가져옵니다.
    startSeq, err := strconv.ParseInt(r.URL.Query().Get("startSeq"), 10, 64)
    if err != nil {
        http.Error(w, "Invalid startSeq", http.StatusBadRequest)
        return
    }
    endSeq, err := strconv.ParseInt(r.URL.Query().Get("endSeq"), 10, 64)
    if err != nil {
        http.Error(w, "Invalid endSeq", http.StatusBadRequest)
        return
    }

    // SEQ 값에 대한 해시 계산
    results := generateHashResults(startSeq, endSeq)

    // JSON 형식으로 결과를 출력합니다.
    w.Header().Set("Content-Type", "application/json")
    if err := json.NewEncoder(w).Encode(results); err != nil {
        http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
        return
    }
}

func main() {
    http.HandleFunc("/", mainPage)
    http.HandleFunc("/generate", generateHashHandler)
    // 정적 파일(스타일 CSS 등)을 서빙합니다.
    http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
    // 서버 시작
    fmt.Println("서버가 http://localhost:8180에서 실행 중입니다.")
    http.ListenAndServe(":8180", nil)
}
