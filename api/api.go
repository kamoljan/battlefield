package api

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"html"
	"image"
	"image/jpeg"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"github.com/nfnt/resize"

	"github.com/kamoljan/battlefield/conf"
	"github.com/kamoljan/battlefield/json"
)

func genPath(file string) string {
	path := fmt.Sprintf(conf.IkuraStore+"%s/%s/%s", file[5:7], file[7:9], file)
	log.Println(path)
	return path
}

func genFile(eid string, color string, width, height int) string {
	file := fmt.Sprintf("%04x_%s_%s_%d_%d", conf.IkuraId, eid, color, width, height)
	log.Println(file)
	return file
}

/*
 *{
 *	status: "ok"
 * 	result: { newborn: "0001_040db0bc2fc49ab41fd81294c7d195c7d1de358b_ACA0AC_100_160" }
 *}
 */
func Put(w http.ResponseWriter, r *http.Request) {
	if r.Method != "PUT" {
		w.Write(json.Message("ERROR", "Not supported Method"))
		return
	}

	log.Println(r)

	reader, err := r.MultipartReader()
	if err != nil {
		w.Write(json.Message("ERROR", "Client should support multipart/form-data"))
		return
	}

	buf := bytes.NewBufferString("")
	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}
		if part.FileName() == "" { // if empy skip this iteration
			continue
		}
		_, err = io.Copy(buf, part)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	defer r.Body.Close()
	img, _, err := image.Decode(buf)
	if err != nil {
		w.Write(json.Message("ERROR", "Unable to decode your image"))
		return
	}
	t0 := time.Now()

	imgBaby := resize.Resize(conf.BabyWidth, 0, img, resize.NearestNeighbor)
	imgInfant := resize.Resize(conf.InfantWidth, 0, imgBaby, resize.NearestNeighbor)
	imgNewborn := resize.Resize(conf.NewbornWidth, 0, imgInfant, resize.NearestNeighbor)
	imgSperm := resize.Resize(conf.Sperm, conf.Sperm, imgNewborn, resize.NearestNeighbor)

	red, green, blue, _ := imgSperm.At(0, 0).RGBA()
	color := fmt.Sprintf("%X%X%X", red>>8, green>>8, blue>>8) // removing 1 byte 9A16->9A

	fileOrig, err := imgToFile(img, color)
	if err != nil {
		w.Write(json.Message("ERROR", "Unable to save your image"))
	}
	fileBaby, err := imgToFile(imgBaby, color)
	fileInfant, err := imgToFile(imgInfant, color)
	fileNewborn, err := imgToFile(imgNewborn, color)

	result := json.Result{
		Newborn: fileNewborn,
	}
	log.Printf("fileOrig=%s,fileBaby=%s,fileInfant=%s,fileNewborn=%s", fileOrig, fileBaby, fileInfant, fileNewborn)
	if err != nil {
		w.Write(json.Message("ERROR", "Unable to save your image meta into db"))
	} else {
		w.Write(json.Message("OK", &result))
	}

	t1 := time.Now()
	log.Printf("The call took %v to run.\n", t1.Sub(t0))
}

func genHash(img image.Image) (string, error) {
	h := sha1.New()
	err := jpeg.Encode(h, img, nil)
	return fmt.Sprintf("%x", h.Sum(nil)), err // generate hash
}

func imgToFile(img image.Image, color string) (string, error) {
	hash, err := genHash(img)
	if err != nil {
		log.Println("Unable to save a file ", err)
		return "", err
	}
	file := genFile(hash, color, img.Bounds().Size().X, img.Bounds().Size().Y)
	path := genPath(file)
	out, err := os.Create(path)
	if err != nil {
		log.Println("Unable to create a file", err)
		return "", err
	}
	defer out.Close()
	err = jpeg.Encode(out, img, nil) // write image to file
	if err != nil {
		log.Println("Unable to save your image to file")
		return "", err
	}
	return file, err
}

func parsePath(eid string) string {
	return fmt.Sprintf(conf.IkuraStore+"%s/%s/%s", eid[5:7], eid[7:9], eid)
}

//http://localhost:9090/egg/0001_8787bec619ff019fd17fe02599a384d580bf6779_9BA4AA_400_300?type=baby
func Get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", conf.Mime)
	w.Header().Set("Cache-Control", fmt.Sprintf("max-age=%d", conf.CacheMaxAge))
	eid := html.EscapeString(r.URL.Path[5:]) //cutting "/egg/"
	log.Println("GET: eid = " + eid)
	path := parsePath(eid)
	log.Println("GET: path = " + path)
	http.ServeFile(w, r, path)
}
