package roundtrip_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/oapi-codegen/oapi-codegen-exp/codegen/internal/test/parameters/roundtrip/client"
	stdhttpparams "github.com/oapi-codegen/oapi-codegen-exp/codegen/internal/test/parameters/roundtrip/stdhttp"
)

func TestStdHttpParameterRoundTrip(t *testing.T) {
	var s stdhttpparams.Server
	handler := stdhttpparams.Handler(&s)
	testImpl(t, handler)
}

// testImpl runs the full parameter roundtrip test suite against any http.Handler.
// The generated client serializes Go values into an HTTP request, the server
// deserializes them and echoes them back as JSON, and we compare the response
// body against the original values.
func testImpl(t *testing.T, handler http.Handler) {
	t.Helper()

	server := "http://example.com"

	expectedObject := client.Object{
		FirstName: "Alex",
		Role:      "admin",
	}

	expectedComplexObject := client.ComplexObject{
		Object:  expectedObject,
		ID:      12345,
		IsAdmin: true,
	}

	expectedArray := []int32{3, 4, 5}

	var expectedPrimitive int32 = 5

	// doRoundTrip sends a request to the handler, asserts 200, and decodes the JSON response.
	doRoundTrip := func(t *testing.T, req *http.Request, target interface{}) {
		t.Helper()
		req.RequestURI = req.URL.RequestURI()
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
		if !assert.Equal(t, http.StatusOK, rec.Code, "server returned %d; body: %s", rec.Code, rec.Body.String()) {
			return
		}
		if target != nil {
			require.NoError(t, json.NewDecoder(rec.Body).Decode(target), "failed to decode response body")
		}
	}

	// =========================================================================
	// Path Parameters
	// =========================================================================
	t.Run("path", func(t *testing.T) {
		t.Run("simple", func(t *testing.T) {
			t.Run("primitive", func(t *testing.T) {
				req, err := client.NewGetSimplePrimitiveRequest(server, expectedPrimitive)
				require.NoError(t, err)
				var got int32
				doRoundTrip(t, req, &got)
				assert.Equal(t, expectedPrimitive, got)
			})

			t.Run("primitive explode", func(t *testing.T) {
				req, err := client.NewGetSimpleExplodePrimitiveRequest(server, expectedPrimitive)
				require.NoError(t, err)
				var got int32
				doRoundTrip(t, req, &got)
				assert.Equal(t, expectedPrimitive, got)
			})

			t.Run("array noExplode", func(t *testing.T) {
				req, err := client.NewGetSimpleNoExplodeArrayRequest(server, expectedArray)
				require.NoError(t, err)
				var got []int32
				doRoundTrip(t, req, &got)
				assert.Equal(t, expectedArray, got)
			})

			t.Run("array explode", func(t *testing.T) {
				req, err := client.NewGetSimpleExplodeArrayRequest(server, expectedArray)
				require.NoError(t, err)
				var got []int32
				doRoundTrip(t, req, &got)
				assert.Equal(t, expectedArray, got)
			})

			t.Run("object noExplode", func(t *testing.T) {
				req, err := client.NewGetSimpleNoExplodeObjectRequest(server, expectedObject)
				require.NoError(t, err)
				var got client.Object
				doRoundTrip(t, req, &got)
				assert.Equal(t, expectedObject, got)
			})

			t.Run("object explode", func(t *testing.T) {
				req, err := client.NewGetSimpleExplodeObjectRequest(server, expectedObject)
				require.NoError(t, err)
				var got client.Object
				doRoundTrip(t, req, &got)
				assert.Equal(t, expectedObject, got)
			})
		})

		t.Run("label", func(t *testing.T) {
			t.Run("primitive", func(t *testing.T) {
				req, err := client.NewGetLabelPrimitiveRequest(server, expectedPrimitive)
				require.NoError(t, err)
				var got int32
				doRoundTrip(t, req, &got)
				assert.Equal(t, expectedPrimitive, got)
			})
			t.Run("primitive explode", func(t *testing.T) {
				req, err := client.NewGetLabelExplodePrimitiveRequest(server, expectedPrimitive)
				require.NoError(t, err)
				var got int32
				doRoundTrip(t, req, &got)
				assert.Equal(t, expectedPrimitive, got)
			})
			t.Run("array noExplode", func(t *testing.T) {
				req, err := client.NewGetLabelNoExplodeArrayRequest(server, expectedArray)
				require.NoError(t, err)
				var got []int32
				doRoundTrip(t, req, &got)
				assert.Equal(t, expectedArray, got)
			})
			t.Run("array explode", func(t *testing.T) {
				req, err := client.NewGetLabelExplodeArrayRequest(server, expectedArray)
				require.NoError(t, err)
				var got []int32
				doRoundTrip(t, req, &got)
				assert.Equal(t, expectedArray, got)
			})
			t.Run("object noExplode", func(t *testing.T) {
				req, err := client.NewGetLabelNoExplodeObjectRequest(server, expectedObject)
				require.NoError(t, err)
				var got client.Object
				doRoundTrip(t, req, &got)
				assert.Equal(t, expectedObject, got)
			})
			t.Run("object explode", func(t *testing.T) {
				req, err := client.NewGetLabelExplodeObjectRequest(server, expectedObject)
				require.NoError(t, err)
				var got client.Object
				doRoundTrip(t, req, &got)
				assert.Equal(t, expectedObject, got)
			})
		})

		t.Run("matrix", func(t *testing.T) {
			t.Run("primitive", func(t *testing.T) {
				req, err := client.NewGetMatrixPrimitiveRequest(server, expectedPrimitive)
				require.NoError(t, err)
				var got int32
				doRoundTrip(t, req, &got)
				assert.Equal(t, expectedPrimitive, got)
			})
			t.Run("primitive explode", func(t *testing.T) {
				req, err := client.NewGetMatrixExplodePrimitiveRequest(server, expectedPrimitive)
				require.NoError(t, err)
				var got int32
				doRoundTrip(t, req, &got)
				assert.Equal(t, expectedPrimitive, got)
			})
			t.Run("array noExplode", func(t *testing.T) {
				req, err := client.NewGetMatrixNoExplodeArrayRequest(server, expectedArray)
				require.NoError(t, err)
				var got []int32
				doRoundTrip(t, req, &got)
				assert.Equal(t, expectedArray, got)
			})
			t.Run("array explode", func(t *testing.T) {
				req, err := client.NewGetMatrixExplodeArrayRequest(server, expectedArray)
				require.NoError(t, err)
				var got []int32
				doRoundTrip(t, req, &got)
				assert.Equal(t, expectedArray, got)
			})
			t.Run("object noExplode", func(t *testing.T) {
				req, err := client.NewGetMatrixNoExplodeObjectRequest(server, expectedObject)
				require.NoError(t, err)
				var got client.Object
				doRoundTrip(t, req, &got)
				assert.Equal(t, expectedObject, got)
			})
			t.Run("object explode", func(t *testing.T) {
				req, err := client.NewGetMatrixExplodeObjectRequest(server, expectedObject)
				require.NoError(t, err)
				var got client.Object
				doRoundTrip(t, req, &got)
				assert.Equal(t, expectedObject, got)
			})
		})

		t.Run("content-based", func(t *testing.T) {
			t.Run("json complex object", func(t *testing.T) {
				// V3 client generates string param for content-based path params.
				// Serialize the object to JSON string for the client.
				jsonBytes, err := json.Marshal(expectedComplexObject)
				require.NoError(t, err)
				req, err := client.NewGetContentObjectRequest(server, string(jsonBytes))
				require.NoError(t, err)
				var got string
				doRoundTrip(t, req, &got)
				// The server echoes the string param, so we compare JSON strings.
				assert.JSONEq(t, string(jsonBytes), got)
			})

			t.Run("passthrough string", func(t *testing.T) {
				req, err := client.NewGetPassThroughRequest(server, "hello world")
				require.NoError(t, err)
				var got string
				doRoundTrip(t, req, &got)
				assert.Equal(t, "hello world", got)
			})
		})
	})

	// =========================================================================
	// Query Parameters
	// =========================================================================
	t.Run("query", func(t *testing.T) {
		t.Run("form", func(t *testing.T) {
			expectedArray2 := []int32{6, 7, 8}
			var expectedPrimitive2 int32 = 100
			expectedPrimitiveString := "123;456"

			t.Run("exploded array only", func(t *testing.T) {
				params := client.GetQueryFormParams{Ea: &expectedArray}
				req, err := client.NewGetQueryFormRequest(server, &params)
				require.NoError(t, err)
				var got client.GetQueryFormParams
				doRoundTrip(t, req, &got)
				require.NotNil(t, got.Ea)
				assert.Equal(t, expectedArray, *got.Ea)
			})

			t.Run("unexploded array only", func(t *testing.T) {
				params := client.GetQueryFormParams{A: &expectedArray2}
				req, err := client.NewGetQueryFormRequest(server, &params)
				require.NoError(t, err)
				var got client.GetQueryFormParams
				doRoundTrip(t, req, &got)
				require.NotNil(t, got.A)
				assert.Equal(t, expectedArray2, *got.A)
			})

			t.Run("exploded primitive only", func(t *testing.T) {
				params := client.GetQueryFormParams{Ep: &expectedPrimitive}
				req, err := client.NewGetQueryFormRequest(server, &params)
				require.NoError(t, err)
				var got client.GetQueryFormParams
				doRoundTrip(t, req, &got)
				require.NotNil(t, got.Ep)
				assert.Equal(t, expectedPrimitive, *got.Ep)
			})

			t.Run("unexploded primitive only", func(t *testing.T) {
				params := client.GetQueryFormParams{P: &expectedPrimitive2}
				req, err := client.NewGetQueryFormRequest(server, &params)
				require.NoError(t, err)
				var got client.GetQueryFormParams
				doRoundTrip(t, req, &got)
				require.NotNil(t, got.P)
				assert.Equal(t, expectedPrimitive2, *got.P)
			})

			t.Run("primitive string", func(t *testing.T) {
				params := client.GetQueryFormParams{Ps: &expectedPrimitiveString}
				req, err := client.NewGetQueryFormRequest(server, &params)
				require.NoError(t, err)
				var got client.GetQueryFormParams
				doRoundTrip(t, req, &got)
				require.NotNil(t, got.Ps)
				assert.Equal(t, expectedPrimitiveString, *got.Ps)
			})

			t.Run("exploded object only", func(t *testing.T) {
				params := client.GetQueryFormParams{Eo: &expectedObject}
				req, err := client.NewGetQueryFormRequest(server, &params)
				require.NoError(t, err)
				var got client.GetQueryFormParams
				doRoundTrip(t, req, &got)
				require.NotNil(t, got.Eo)
				assert.Equal(t, expectedObject, *got.Eo)
			})

			t.Run("unexploded object only", func(t *testing.T) {
				params := client.GetQueryFormParams{O: &expectedObject}
				req, err := client.NewGetQueryFormRequest(server, &params)
				require.NoError(t, err)
				var got client.GetQueryFormParams
				doRoundTrip(t, req, &got)
				require.NotNil(t, got.O)
				assert.Equal(t, expectedObject, *got.O)
			})
		})

		t.Run("deepObject", func(t *testing.T) {
			params := client.GetDeepObjectParams{DeepObj: expectedComplexObject}
			req, err := client.NewGetDeepObjectRequest(server, &params)
			require.NoError(t, err)
			var got client.GetDeepObjectParams
			doRoundTrip(t, req, &got)
			// DeepObj is typed as any, compare via JSON.
			gotJSON, _ := json.Marshal(got.DeepObj)
			expectedJSON, _ := json.Marshal(expectedComplexObject)
			assert.JSONEq(t, string(expectedJSON), string(gotJSON))
		})
	})

	// =========================================================================
	// Header Parameters
	// =========================================================================
	t.Run("header", func(t *testing.T) {
		expectedArray2 := []int32{6, 7, 8}
		var expectedPrimitive2 int32 = 100

		t.Run("primitive only", func(t *testing.T) {
			params := client.GetHeaderParams{XPrimitive: &expectedPrimitive}
			req, err := client.NewGetHeaderRequest(server, &params)
			require.NoError(t, err)
			var got client.GetHeaderParams
			doRoundTrip(t, req, &got)
			require.NotNil(t, got.XPrimitive)
			assert.Equal(t, expectedPrimitive, *got.XPrimitive)
		})

		t.Run("primitive exploded only", func(t *testing.T) {
			params := client.GetHeaderParams{XPrimitiveExploded: &expectedPrimitive2}
			req, err := client.NewGetHeaderRequest(server, &params)
			require.NoError(t, err)
			var got client.GetHeaderParams
			doRoundTrip(t, req, &got)
			require.NotNil(t, got.XPrimitiveExploded)
			assert.Equal(t, expectedPrimitive2, *got.XPrimitiveExploded)
		})

		t.Run("array only", func(t *testing.T) {
			params := client.GetHeaderParams{XArray: &expectedArray}
			req, err := client.NewGetHeaderRequest(server, &params)
			require.NoError(t, err)
			var got client.GetHeaderParams
			doRoundTrip(t, req, &got)
			require.NotNil(t, got.XArray)
			assert.Equal(t, expectedArray, *got.XArray)
		})

		t.Run("array exploded only", func(t *testing.T) {
			params := client.GetHeaderParams{XArrayExploded: &expectedArray2}
			req, err := client.NewGetHeaderRequest(server, &params)
			require.NoError(t, err)
			var got client.GetHeaderParams
			doRoundTrip(t, req, &got)
			require.NotNil(t, got.XArrayExploded)
			assert.Equal(t, expectedArray2, *got.XArrayExploded)
		})

		t.Run("object only", func(t *testing.T) {
			params := client.GetHeaderParams{XObject: &expectedObject}
			req, err := client.NewGetHeaderRequest(server, &params)
			require.NoError(t, err)
			var got client.GetHeaderParams
			doRoundTrip(t, req, &got)
			require.NotNil(t, got.XObject)
			assert.Equal(t, expectedObject, *got.XObject)
		})
	})

	// =========================================================================
	// Cookie Parameters
	// =========================================================================
	t.Run("cookie", func(t *testing.T) {
		expectedArray2 := []int32{6, 7, 8}
		var expectedPrimitive2 int32 = 100

		t.Run("primitive only", func(t *testing.T) {
			params := client.GetCookieParams{P: &expectedPrimitive}
			req, err := client.NewGetCookieRequest(server, &params)
			require.NoError(t, err)
			var got client.GetCookieParams
			doRoundTrip(t, req, &got)
			require.NotNil(t, got.P)
			assert.Equal(t, expectedPrimitive, *got.P)
		})

		t.Run("primitive exploded only", func(t *testing.T) {
			params := client.GetCookieParams{Ep: &expectedPrimitive2}
			req, err := client.NewGetCookieRequest(server, &params)
			require.NoError(t, err)
			var got client.GetCookieParams
			doRoundTrip(t, req, &got)
			require.NotNil(t, got.Ep)
			assert.Equal(t, expectedPrimitive2, *got.Ep)
		})

		t.Run("array only", func(t *testing.T) {
			params := client.GetCookieParams{A: &expectedArray}
			req, err := client.NewGetCookieRequest(server, &params)
			require.NoError(t, err)
			var got client.GetCookieParams
			doRoundTrip(t, req, &got)
			require.NotNil(t, got.A)
			assert.Equal(t, expectedArray, *got.A)
		})

		t.Run("array exploded only", func(t *testing.T) {
			params := client.GetCookieParams{Ea: &expectedArray2}
			req, err := client.NewGetCookieRequest(server, &params)
			require.NoError(t, err)
			var got client.GetCookieParams
			doRoundTrip(t, req, &got)
			require.NotNil(t, got.Ea)
			assert.Equal(t, expectedArray2, *got.Ea)
		})

		t.Run("object only", func(t *testing.T) {
			params := client.GetCookieParams{O: &expectedObject}
			req, err := client.NewGetCookieRequest(server, &params)
			require.NoError(t, err)
			var got client.GetCookieParams
			doRoundTrip(t, req, &got)
			require.NotNil(t, got.O)
			assert.Equal(t, expectedObject, *got.O)
		})
	})
}
