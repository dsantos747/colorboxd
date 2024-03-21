package colorboxd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"reflect"
	"slices"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

func HTTPWriteList(w http.ResponseWriter, r *http.Request) {
	var err error

	// Read env variables - local development only - comment out for production
	err = godotenv.Load()
	if err != nil {
		fmt.Printf("Could not load environment variables from .env file: %v\n", err)
		return
	}

	// Set necessary headers for CORS
	w.Header().Set("Access-Control-Allow-Origin", os.Getenv("BASE_URL"))
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	var responseData WriteListRequest
	err = json.NewDecoder(r.Body).Decode(&responseData)
	if err != nil {
		http.Error(w, "Missing or empty 'accessToken' query parameter", http.StatusBadRequest)
		return
	}

	listUpdateRequest, err := prepareListUpdateRequest(responseData.List, responseData.Offset, responseData.SortMethod, responseData.Reverse)
	if err != nil {
		http.Error(w, "Error preparing list update request body", http.StatusInternalServerError)
		return
	}

	message, err := writeListSorting(responseData.AccessToken, responseData.List.ID, *listUpdateRequest)
	if err != nil {
		http.Error(w, "Error updating user list", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(message)
}

// Sort the list as per the specified method, then return a ListUpdateRequest, as required by Letterboxd endpoint
func prepareListUpdateRequest(list ListWithEntries, offset int, sortMethod string, reverse bool) (*ListUpdateRequest, error) {
	generateSortFunction := func(sortMethod string) (func(Entry, Entry) int, error) {
		// Check if sortMethod is valid
		if _, ok := reflect.TypeOf(SortVals{}).FieldByName(sortMethod); !ok {
			return nil, errors.New("error: provided sort method not recognized")
		}

		// Generate the sort function from the type
		sortFunction := func(a, b Entry) int {
			A := reflect.ValueOf(a.SortVals).FieldByName(sortMethod).Int()
			B := reflect.ValueOf(b.SortVals).FieldByName(sortMethod).Int()
			return int(A - B)
		}

		return sortFunction, nil
	}

	sortFunction, err := generateSortFunction(sortMethod)
	if err != nil {
		return nil, err
	}
	slices.SortFunc(list.Entries, sortFunction)

	n := len(list.Entries)
	currentPositions := make(map[string]int)
	var finishSlice []FilmTargetPosition

	for i, entry := range list.Entries {
		initPos, err := strconv.Atoi(entry.EntryID)
		if err != nil {
			fmt.Println("Error parsing entryId to int")
			return nil, err
		}
		endPos := ((i + n) - offset) % n
		if reverse {
			endPos = (n - endPos) % n
		}

		currentPositions[entry.FilmID] = initPos
		finishSlice = append(finishSlice, FilmTargetPosition{entry.FilmID, endPos})
	}

	slices.SortFunc(finishSlice, func(a, b FilmTargetPosition) int {
		return a.Position - b.Position
	})

	updateEntries := createListUpdateEntries(currentPositions, finishSlice)
	request := ListUpdateRequest{Version: list.Version, Entries: updateEntries}

	return &request, nil
}

// Create a set of instructions that, applied in turn, result in the correctly-sorted list.
func createListUpdateEntries(currentPositions map[string]int, finishPositions []FilmTargetPosition) []listUpdateEntry {
	var updateEntries []listUpdateEntry
	for _, film := range finishPositions {
		currPos := currentPositions[film.FilmId]
		if film.Position == currPos {
			continue
		}

		updateEntries = append(updateEntries, listUpdateEntry{Action: "UPDATE", Position: currPos, NewPosition: film.Position})
		currentPositions[film.FilmId] = film.Position

		for f, cP := range currentPositions {
			if f != film.FilmId {
				if film.Position <= cP && cP < currPos {
					currentPositions[f]++
				}
			}
		}
	}
	return updateEntries
}

// Send request to Letterboxd endpoint to update list.
func writeListSorting(token, id string, listUpdateRequest ListUpdateRequest) (*[]string, error) {

	// Prepare endpoint and body for PATCH request
	method := "PATCH"
	endpoint := fmt.Sprintf("%s/list/%s", os.Getenv("LBOXD_BASEURL"), id)
	headers := map[string]string{"Authorization": fmt.Sprintf("Bearer %s", token), "Content-Type": "application/json", "X-HTTP-Method-Override": "PATCH"}
	body, err := json.Marshal(listUpdateRequest)
	if err != nil {
		return nil, err
	}

	response, err := MakeHTTPRequest(method, endpoint, bytes.NewReader(body), headers)
	if err != nil {
		return nil, fmt.Errorf("error making HTTP request: %v", err)
	}
	defer response.Body.Close()

	var responseData ListUpdateResponse
	if err = json.NewDecoder(response.Body).Decode(&responseData); err != nil {
		fmt.Println(err)
		return nil, err
	}

	var message []string
	if len(responseData.Messages) != 0 {
		for _, m := range responseData.Messages {
			message = append(message, fmt.Sprintf("%s: %s - %s", m.Type, m.Code, m.Title))
		}
		errorStr := "The letterboxd API responded with the following errors: " + strings.Join(message, "; ")
		return &message, errors.New(errorStr)
	}

	message = []string{"List updated successfully"}

	return &message, nil
}
