package db

import (
	"os"
	"testing"

	"github.com/ONSdigital/dp-search-scrubber-api/config"
	"github.com/stretchr/testify/assert"
)

func TestGetIndustry(t *testing.T) {
	// Split API tests from Unit tests
	skipUnitTests(t)

	// create a test file with test data and write test data to the test file
	ir := mockIndustryData(t)

	// check if the function returns the expected result
	assert.Len(t, ir, 2, "Expected 2 areas, but got %d", len(ir))

	expected := []struct {
		code string
		name string
	}{
		{
			code: "TestCode1",
			name: "TestName1",
		},
		{
			code: "TestCode2",
			name: "TestName2",
		},
	}

	for i, exp := range expected {
		assert.Equalf(t, exp.code, ir[i].Code, "Unexpected Code at index %d: %s", i, ir[i].Code)
		assert.Equalf(t, exp.name, ir[i].Name, "Unexpected Name at index %d: %s", i, ir[i].Name)
	}
}

func mockIndustryData(t *testing.T) []Industry {
	testFile, err := os.Create("test.csv")
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	defer testFile.Close()
	defer os.Remove("test.csv")

	_, err = testFile.WriteString("SIC Code,Description\nTestCode1,TestName1\nTestCode2,TestName2\n")
	if err != nil {
		t.Fatalf("Failed to write test data: %v", err)
	}

	cfg := config.Config{
		IndustryDataFile: "test.csv",
	}

	ir, err := getIndustry(&cfg)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	return ir
}
func skipUnitTests(t *testing.T) {
	if os.Getenv("UNIT") != "" {
		t.Skip("Skipping Unit tests in CI environment")
	}
}
