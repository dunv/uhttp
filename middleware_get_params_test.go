package uhttp

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/dunv/ulog"
)

func testRequirementFail(requirement R, actual map[string]string, unexpectedKey string, t *testing.T) {
	u := NewUHTTP()
	validatedMap := R{}
	err := u.validateParams(requirement, actual, validatedMap, true)
	if err == nil {
		t.Error(fmt.Errorf("validation mistakenly succeeded"))
	}
	if _, ok := validatedMap[unexpectedKey]; ok {
		t.Error(fmt.Errorf("param in validated map"))
	}
}

func testRequirementSuccess(requirement R, actual map[string]string, expectedKey string, expectedValue interface{}, t *testing.T) {
	u := NewUHTTP()
	validatedMap := R{}
	err := u.validateParams(requirement, actual, validatedMap, true)
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
		R{"test": ENUM("test1", "test2")},
		map[string]string{"test": "test2"},
		"test",
		"test2",
		t,
	)
}

func TestEnumRequirementFail(t *testing.T) {
	testRequirementFail(
		R{"test": ENUM("test1", "test2")},
		map[string]string{"test": "test3"},
		"test",
		t,
	)
}

func TestStringRequirementSuccess(t *testing.T) {
	testRequirementSuccess(
		R{"test": STRING},
		map[string]string{"test": "test2"},
		"test",
		"test2",
		t,
	)
}

func TestBoolRequirementSuccess1(t *testing.T) {
	testRequirementSuccess(
		R{"test": BOOL},
		map[string]string{"test": "true"},
		"test",
		true,
		t,
	)
}

func TestBoolRequirementSuccess2(t *testing.T) {
	testRequirementSuccess(
		R{"test": BOOL},
		map[string]string{"test": "false"},
		"test",
		false,
		t,
	)
}

func TestBoolRequirementFail(t *testing.T) {
	testRequirementFail(
		R{"test": BOOL},
		map[string]string{"test": "ture"},
		"test",
		t,
	)
}

func TestIntRequirementSuccess(t *testing.T) {
	testRequirementSuccess(
		R{"test": INT},
		map[string]string{"test": "2"},
		"test",
		int(2),
		t,
	)
}

func TestIntRequirementFail(t *testing.T) {
	testRequirementFail(
		R{"test": INT},
		map[string]string{"test": "ture"},
		"test",
		t,
	)
}

func TestInt32RequirementSuccess(t *testing.T) {
	testRequirementSuccess(
		R{"test": INT32},
		map[string]string{"test": "2"},
		"test",
		int32(2),
		t,
	)
}

func TestInt32RequirementFail(t *testing.T) {
	testRequirementFail(
		R{"test": INT32},
		map[string]string{"test": "ture"},
		"test",
		t,
	)
}

func TestInt64RequirementSuccess(t *testing.T) {
	testRequirementSuccess(
		R{"test": INT64},
		map[string]string{"test": "2"},
		"test",
		int64(2),
		t,
	)
}

func TestInt64RequirementFail(t *testing.T) {
	testRequirementFail(
		R{"test": INT64},
		map[string]string{"test": "ture"},
		"test",
		t,
	)
}

func TestFloat32RequirementSuccess(t *testing.T) {
	testRequirementSuccess(
		R{"test": FLOAT32},
		map[string]string{"test": "2.2"},
		"test",
		float32(2.2),
		t,
	)
}

func TestFloat32RequirementFail(t *testing.T) {
	testRequirementFail(
		R{"test": FLOAT32},
		map[string]string{"test": "ture"},
		"test",
		t,
	)
}

func TestFloat64RequirementSuccess(t *testing.T) {
	testRequirementSuccess(
		R{"test": FLOAT64},
		map[string]string{"test": "2.2"},
		"test",
		float64(2.2),
		t,
	)
}

func TestFloat64RequirementFail(t *testing.T) {
	testRequirementFail(
		R{"test": FLOAT64},
		map[string]string{"test": "ture"},
		"test",
		t,
	)
}

func TestShortDateRequirementSuccess(t *testing.T) {
	testRequirementSuccess(
		R{"test": SHORT_DATE},
		map[string]string{"test": "2019-08-09"},
		"test",
		time.Date(2019, 8, 9, 0, 0, 0, 0, time.UTC),
		t,
	)
}

func TestShortDateRequirementFail(t *testing.T) {
	testRequirementFail(
		R{"test": SHORT_DATE},
		map[string]string{"test": "2019-13-30"},
		"test",
		t,
	)
}

func TestRFC3339DateRequirementSuccess(t *testing.T) {
	testRequirementSuccess(
		R{"test": RFC3339_DATE},
		map[string]string{"test": "2002-10-02T10:00:00-05:00"},
		"test",
		time.Date(2002, 10, 2, 10, 0, 0, 0, time.FixedZone("UTC-5", -5*60*60)),
		t,
	)
}

func TestRFC3339DateRequirementFail(t *testing.T) {
	testRequirementFail(
		R{"test": RFC3339_DATE},
		map[string]string{"test": "2002-10-02T30:00:00-05:00"},
		"test",
		t,
	)
}

func TestDurationRequirementSuccess(t *testing.T) {
	testRequirementSuccess(
		R{"test": DURATION},
		map[string]string{"test": "10h"},
		"test",
		time.Hour*10,
		t,
	)
}

func TestDurationRequirementFail(t *testing.T) {
	testRequirementFail(
		R{"test": DURATION},
		map[string]string{"test": "5k"},
		"test",
		t,
	)
}

func TestRequirementsInHandler(t *testing.T) {
	ulog.SetWriter(io.Discard, nil)
	u := NewUHTTP()
	handler := NewHandler(
		WithRequiredGet(R{
			"string":      STRING,
			"bool":        BOOL,
			"int":         INT,
			"int32":       INT32,
			"int64":       INT64,
			"float32":     FLOAT32,
			"float64":     FLOAT64,
			"shortDate":   SHORT_DATE,
			"rfc3339Date": RFC3339_DATE,
			"duration":    DURATION,
		}),
		WithGet(func(r *http.Request, ret *int) interface{} {
			return map[string]interface{}{
				"string":      GetAsString("string", r),
				"bool":        GetAsBool("bool", r),
				"int":         GetAsInt("int", r),
				"int32":       GetAsInt32("int32", r),
				"int64":       GetAsInt64("int64", r),
				"float32":     GetAsFloat32("float32", r),
				"float64":     GetAsFloat64("float64", r),
				"shortDate":   GetAsTime("shortDate", r),
				"rfc3339Date": GetAsTime("rfc3339Date", r),
				"duration":    GetAsDuration("duration", r).String(),
			}
		}),
	)
	u.Handle("/test", handler)
	RequireHTTPBodyJSONEq(t, u.ServeMux().ServeHTTP, http.MethodGet, "/test", url.Values{
		"string":      []string{"myString"},
		"bool":        []string{"true"},
		"int":         []string{"42"},
		"int32":       []string{"42"},
		"int64":       []string{"42"},
		"float32":     []string{"42.42"},
		"float64":     []string{"42.42"},
		"shortDate":   []string{"2021-10-15"},
		"rfc3339Date": []string{"2021-10-15T08:30:00Z"},
		"duration":    []string{"5m"},
	}, `{
		"string": "myString",
		"bool": true,
		"int": 42,
		"int32": 42,
		"int64": 42,
		"float32": 42.42,
		"float64": 42.42,
		"shortDate": "2021-10-15T00:00:00Z",
		"rfc3339Date": "2021-10-15T08:30:00Z",
		"duration": "5m0s"
	}`)
}
