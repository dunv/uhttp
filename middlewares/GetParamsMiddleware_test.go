package middlewares

import (
	"fmt"
	"testing"
	"time"

	"github.com/dunv/uhttp/params"
)

func testRequirementFail(requirement params.R, actual map[string]string, unexpectedKey string, t *testing.T) {
	validatedMap := params.R{}
	err := params.ValidateParams(requirement, actual, validatedMap, true)
	if err == nil {
		t.Error(fmt.Errorf("validation mistakenly succeeded"))
	}
	if _, ok := validatedMap[unexpectedKey]; ok {
		t.Error(fmt.Errorf("param in validated map"))
	}
}

func testRequirementSuccess(requirement params.R, actual map[string]string, expectedKey string, expectedValue interface{}, t *testing.T) {
	validatedMap := params.R{}
	err := params.ValidateParams(requirement, actual, validatedMap, true)
	if err != nil {
		t.Error(fmt.Errorf("validation mistakenly failed"))
	}

	if _, ok := validatedMap[expectedKey]; !ok {
		t.Error(fmt.Errorf("param not in validated map"))
	}

	if expectedTime, ok := expectedValue.(time.Time); ok {
		if actualTime, ok := validatedMap[expectedKey].(time.Time); ok {
			if !expectedTime.Equal(actualTime) {
				t.Error(fmt.Errorf("incorrect paramValue in validated map actual: %v expected: %v", validatedMap[expectedKey], expectedValue))
			}
		}
	} else if validatedMap[expectedKey] != expectedValue {
		t.Error(fmt.Errorf("incorrect paramValue in validated map actual: %v expected: %v", validatedMap[expectedKey], expectedValue))
	}
}

func TestEnumRequirementSuccess(t *testing.T) {
	testRequirementSuccess(
		params.R{"test": params.ENUM("test1", "test2")},
		map[string]string{"test": "test2"},
		"test",
		"test2",
		t,
	)
}

func TestEnumRequirementFail(t *testing.T) {
	testRequirementFail(
		params.R{"test": params.ENUM("test1", "test2")},
		map[string]string{"test": "test3"},
		"test",
		t,
	)
}

func TestStringRequirementSuccess(t *testing.T) {
	testRequirementSuccess(
		params.R{"test": params.STRING},
		map[string]string{"test": "test2"},
		"test",
		"test2",
		t,
	)
}

func TestBoolRequirementSuccess1(t *testing.T) {
	testRequirementSuccess(
		params.R{"test": params.BOOL},
		map[string]string{"test": "true"},
		"test",
		true,
		t,
	)
}

func TestBoolRequirementSuccess2(t *testing.T) {
	testRequirementSuccess(
		params.R{"test": params.BOOL},
		map[string]string{"test": "false"},
		"test",
		false,
		t,
	)
}

func TestBoolRequirementFail(t *testing.T) {
	testRequirementFail(
		params.R{"test": params.BOOL},
		map[string]string{"test": "ture"},
		"test",
		t,
	)
}

func TestIntRequirementSuccess(t *testing.T) {
	testRequirementSuccess(
		params.R{"test": params.INT},
		map[string]string{"test": "2"},
		"test",
		int(2),
		t,
	)
}

func TestIntRequirementFail(t *testing.T) {
	testRequirementFail(
		params.R{"test": params.INT},
		map[string]string{"test": "ture"},
		"test",
		t,
	)
}

func TestInt32RequirementSuccess(t *testing.T) {
	testRequirementSuccess(
		params.R{"test": params.INT32},
		map[string]string{"test": "2"},
		"test",
		int32(2),
		t,
	)
}

func TestInt32RequirementFail(t *testing.T) {
	testRequirementFail(
		params.R{"test": params.INT32},
		map[string]string{"test": "ture"},
		"test",
		t,
	)
}

func TestInt64RequirementSuccess(t *testing.T) {
	testRequirementSuccess(
		params.R{"test": params.INT64},
		map[string]string{"test": "2"},
		"test",
		int64(2),
		t,
	)
}

func TestInt64RequirementFail(t *testing.T) {
	testRequirementFail(
		params.R{"test": params.INT64},
		map[string]string{"test": "ture"},
		"test",
		t,
	)
}

func TestFloat32RequirementSuccess(t *testing.T) {
	testRequirementSuccess(
		params.R{"test": params.FLOAT32},
		map[string]string{"test": "2.2"},
		"test",
		float32(2.2),
		t,
	)
}

func TestFloat32RequirementFail(t *testing.T) {
	testRequirementFail(
		params.R{"test": params.FLOAT32},
		map[string]string{"test": "ture"},
		"test",
		t,
	)
}

func TestFloat64RequirementSuccess(t *testing.T) {
	testRequirementSuccess(
		params.R{"test": params.FLOAT64},
		map[string]string{"test": "2.2"},
		"test",
		float64(2.2),
		t,
	)
}

func TestFloat64RequirementFail(t *testing.T) {
	testRequirementFail(
		params.R{"test": params.FLOAT64},
		map[string]string{"test": "ture"},
		"test",
		t,
	)
}

func TestShortDateRequirementSuccess(t *testing.T) {
	testRequirementSuccess(
		params.R{"test": params.SHORT_DATE},
		map[string]string{"test": "2019-08-09"},
		"test",
		time.Date(2019, 8, 9, 0, 0, 0, 0, time.UTC),
		t,
	)
}

func TestShortDateRequirementFail(t *testing.T) {
	testRequirementFail(
		params.R{"test": params.SHORT_DATE},
		map[string]string{"test": "2019-13-30"},
		"test",
		t,
	)
}

func TestRFC3339DateRequirementSuccess(t *testing.T) {
	testRequirementSuccess(
		params.R{"test": params.RFC3339_DATE},
		map[string]string{"test": "2002-10-02T10:00:00-05:00"},
		"test",
		time.Date(2002, 10, 2, 10, 0, 0, 0, time.FixedZone("UTC-5", -5*60*60)),
		t,
	)
}

func TestRFC3339DateRequirementFail(t *testing.T) {
	testRequirementFail(
		params.R{"test": params.RFC3339_DATE},
		map[string]string{"test": "2002-10-02T30:00:00-05:00"},
		"test",
		t,
	)
}
