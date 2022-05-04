package converter_test

import (
	"scraper/converter"
	"scraper/dto"
	"scraper/storage/model"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFromDTO(t *testing.T) {

	location, _ := time.LoadLocation("Europe/Chisinau")

	testCases := []struct {
		name     string
		incoming dto.Film
		expected model.Film
	}{
		{
			name: "title with spaces and newlines",
			incoming: dto.Film{
				Title: `    The Contractor   


   2D     `,
				Lang: "(RU)",
				Date: "(10.03.2022)",
			},
			expected: model.Film{
				ID:        0,
				Title:     "The Contractor",
				Lang:      "RU",
				Dimension: "2D",
				StartDate: time.Date(2022, time.Month(3), 10, 0, 0, 0, 0, location),
			},
		},
		{
			name: "lang data with spaces and newline",
			incoming: dto.Film{
				Title: `The Contractor 2D`,
				Lang: ` 

( RU ) 
  `,
				Date: "(10.03.2022)",
			},
			expected: model.Film{
				ID:        0,
				Title:     "The Contractor",
				Lang:      "RU",
				Dimension: "2D",
				StartDate: time.Date(2022, time.Month(3), 10, 0, 0, 0, 0, location),
			},
		},
		{
			name: "missing dimension",
			incoming: dto.Film{
				Title: `The Contractor`,
				Lang:  `(RU)`,
				Date:  "(10.03.2022)",
			},
			expected: model.Film{
				ID:        0,
				Title:     "The Contractor",
				Lang:      "RU",
				Dimension: "2D",
				StartDate: time.Date(2022, time.Month(3), 10, 0, 0, 0, 0, location),
			},
		},
		{
			name: "3D dimension ",
			incoming: dto.Film{
				Title: `The Contractor 3D`,
				Lang:  `(RU)`,
				Date:  "(10.03.2022)",
			},
			expected: model.Film{
				ID:        0,
				Title:     "The Contractor",
				Lang:      "RU",
				Dimension: "3D",
				StartDate: time.Date(2022, time.Month(3), 10, 0, 0, 0, 0, location),
			},
		},
		{
			name: "3D dimension ",
			incoming: dto.Film{
				Title: `The Contractor 3D`,
				Lang:  `(KO)`,
				Date:  "(10.03.2022)",
			},
			expected: model.Film{
				ID:        0,
				Title:     "The Contractor",
				Lang:      "KO",
				Dimension: "3D",
				StartDate: time.Date(2022, time.Month(3), 10, 0, 0, 0, 0, location),
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			output, _ := converter.FromDTO(testCase.incoming)
			assert.Equal(t, testCase.expected, output)
		})
	}

}
