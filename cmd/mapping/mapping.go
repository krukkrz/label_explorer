package mapping

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/krukkrz/label_loader/releasescounts"
)

type Details struct {
	Id     int
	Artist string
	Genres []string
	Styles []string
}

func DetailsToReleaseCount(labelId int, details chan Details) ([]releasescounts.ReleaseCount, error) {
	resultsMap := make(map[string]int)
	for detail := range details {
		for _, style := range detail.Styles {
			key := fmt.Sprintf("%d#%s#%s", labelId, detail.Artist, style)
			resultsMap[key]++
		}
		for _, genre := range detail.Genres {
			key := fmt.Sprintf("%d#%s#%s", labelId, detail.Artist, genre)
			resultsMap[key]++
		}
	}

	results, err := mapToReleaseCounts(resultsMap)
	if err != nil {
		return nil, fmt.Errorf("failed to map to release counts: %v", err)
	}
	return results, nil
}

func mapToReleaseCounts(m map[string]int) ([]releasescounts.ReleaseCount, error) {
	var results []releasescounts.ReleaseCount
	for k, v := range m {
		item, err := parseKeyValue(k, v)
		if err != nil {
			return nil, fmt.Errorf("could not map to release count object: %v", err)
		}
		results = append(results, item)
	}
	return results, nil
}

func parseKeyValue(key string, value int) (releasescounts.ReleaseCount, error) {
	elements := strings.Split(key, "#")
	if len(elements) != 3 {
		return releasescounts.ReleaseCount{}, fmt.Errorf("key does not contain appropriate number of elements")
	}
	labelId, err := strconv.Atoi(elements[0])
	if err != nil {
		return releasescounts.ReleaseCount{}, fmt.Errorf("error while parsing label id: %v", err)
	}
	return releasescounts.ReleaseCount{
		LabelId:      labelId,
		ArtistName:   elements[1],
		StyleGenre:   elements[2],
		ReleaseCount: value,
	}, nil
}
