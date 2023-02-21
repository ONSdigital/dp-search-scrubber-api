package mock

import (
	"os"
	"testing"
)

type DbMockStruct struct {
	testAreaFile     *os.File
	testIndustryFile *os.File
}

func CreateFiles(t *testing.T) DbMockStruct {
	af, err := os.Create("area.csv")
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	ti, err := os.Create("industry.csv")
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	_, err = af.WriteString("Output Area Code,Local Authority Code,Local Authority Name,Region/Country Code,Region/Country Name\nTest Output Area Code1,Test LAC1,Test LAN1,Test RC1,Test RN 1\nTest Output Area Code2,Test LAC2,Test LAN2,Test RC2,Test RN 2\n")
	if err != nil {
		t.Fatalf("Failed to write test data: %v", err)
	}

	_, err = ti.WriteString("SIC Code,Description\nTestCode1,TestName1\nTestCode2,TestName2\n")
	if err != nil {
		t.Fatalf("Failed to write test data: %v", err)
	}

	return DbMockStruct{
		testAreaFile:     af,
		testIndustryFile: ti,
	}
}

func (m *DbMockStruct) CloseFiles() {
	m.testAreaFile.Close()
	os.Remove("area.csv")

	m.testIndustryFile.Close()
	os.Remove("industry.csv")
}
