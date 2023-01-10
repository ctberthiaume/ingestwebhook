package serve

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os/exec"
)

type pathParts struct {
	bucket string
	key    string
}

func (p pathParts) String() string {
	return p.bucket + "/" + p.key
}

// Start starts a webserver to process RockBLOCK messages received at /message
func Start(addr string) error {
	http.HandleFunc("/hooks/healthcheck", handleHealthCheck)
	http.HandleFunc("/hooks/minio", handleJSONMinioMessage(func(p pathParts) {
		log.Printf("received notification for %s\n", p)
		dispatchCmd := exec.Command("nomad", "job", "dispatch", "--meta", "bucket="+p.bucket, "--meta", "key="+p.key, "ingest")
		dispatchOut, err := dispatchCmd.Output()
		if err != nil {
			log.Printf("error dispatching ingest job: %v\n", err)
		}
		log.Printf("command output = %s", string(dispatchOut))
	}))
	err := http.ListenAndServe(addr, nil)
	return err
}

func handleHealthCheck(w http.ResponseWriter, req *http.Request) {
	log.Printf("new request from %v for %v\n", req.RemoteAddr, req.URL.Path)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/text")
	w.Write([]byte("healthy"))
}

func handleJSONMinioMessage(cb func(pathParts)) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		log.Printf("new request from %v for %v\n", req.RemoteAddr, req.URL.Path)
		body, err := io.ReadAll(req.Body)
		if err != nil {
			log.Printf("%v\n", err)
			return
		}

		log.Printf("body = %v\n", string(body))

		parts, err := parseMinioJson(body)
		if err != nil {
			log.Printf("could not decode JSON body: %v\n", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)

		cb(parts)
	}
}

func parseMinioJson(b []byte) (pathParts, error) {
	var data map[string]interface{}
	err := json.Unmarshal(b, &data)
	if err != nil {
		return pathParts{}, fmt.Errorf("could not unmarshall json: %s", err)
	}

	recs, ok := data["Records"]
	if !ok {
		return pathParts{}, fmt.Errorf("'records' not present in JSON")
	}
	recsVal, ok := recs.([]interface{})
	if !ok {
		return pathParts{}, fmt.Errorf("problem reading records array")
	}
	if len(recsVal) == 0 {
		return pathParts{}, fmt.Errorf("no records")
	}
	recVal, ok := recsVal[0].(map[string]interface{})
	if !ok {
		return pathParts{}, fmt.Errorf("bad 'Records' item value")
	}

	s3, ok := recVal["s3"]
	if !ok {
		return pathParts{}, fmt.Errorf("'s3' not present in Records item")
	}
	s3Val, ok := s3.(map[string]interface{})
	if !ok {
		return pathParts{}, fmt.Errorf("bad 's3' item")
	}

	bucket, ok := s3Val["bucket"]
	if !ok {
		return pathParts{}, fmt.Errorf("'s3.bucket' not present Records item")
	}
	bucketVal, ok := bucket.(map[string]interface{})
	if !ok {
		return pathParts{}, fmt.Errorf("bad 's3.bucket' value")
	}
	bucketName, ok := bucketVal["name"]
	if !ok {
		return pathParts{}, fmt.Errorf("'s3.bucket.name' not present Records item")
	}
	bucketNameVal, ok := bucketName.(string)
	if !ok {
		return pathParts{}, fmt.Errorf("bad 's3.bucket.name' value")
	}

	object, ok := s3Val["object"]
	if !ok {
		return pathParts{}, fmt.Errorf("'s3.object' not present Records item")
	}
	objectVal, ok := object.(map[string]interface{})
	if !ok {
		return pathParts{}, fmt.Errorf("bad 's3.object' value")
	}
	objectKey, ok := objectVal["key"]
	if !ok {
		return pathParts{}, fmt.Errorf("'s3.object.key' not present Records item")
	}
	objectKeyVal, ok := objectKey.(string)
	if !ok {
		return pathParts{}, fmt.Errorf("bad 's3.object.key' value")
	}

	return pathParts{bucketNameVal, objectKeyVal}, nil
}
