package main

import (
	"net/http"
	"os"
	"strings"
)

type operatorAssetHandler struct{}

func (operatorAssetHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		rw.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if !strings.HasPrefix(req.URL.Path, operatorAssetURLBase) {
		rw.WriteHeader(http.StatusNotFound)
		return
	}

	cacheDir, err := ensureOperatorCacheDir("")
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	assetPath := operatorAssetFilePath(cacheDir, req.URL.Path)
	file, err := os.Open(assetPath)
	if err != nil {
		if os.IsNotExist(err) {
			rw.WriteHeader(http.StatusNotFound)
			return
		}
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	http.ServeContent(rw, req, stat.Name(), stat.ModTime(), file)
}
