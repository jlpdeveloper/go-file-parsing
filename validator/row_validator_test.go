package validator

import (
	"context"
	"fmt"
	"go-file-parsing/config"
	"sync"
	"testing"
	"time"
)

func TestValidate_StopsOnFirstError(t *testing.T) {
	successValidator := func(_ *RowValidatorContext, _ *[]string) error {
		return nil
	}
	errorValidator := func(_ *RowValidatorContext, _ *[]string) error {
		return fmt.Errorf("bad column")
	}
	v := CsvRowValidator{
		config:      &config.ParserConfig{},
		cacheClient: &MockCache{},
		// Config as needed
		colValidators: []ColValidator{
			successValidator,
			errorValidator,   // should cause Validate to stop and return err
			successValidator, // should NOT be called thanks to early exit
		},
	}
	err := v.Validate("irrelevant,row,data")
	if err == nil || err.Error() != "bad column" {
		t.Errorf("expected error from errorValidator; got %v", err)
	}
}

func TestValidate_SuccessfulValidation(t *testing.T) {
	successValidator := func(_ *RowValidatorContext, _ *[]string) error {
		return nil
	}
	v := CsvRowValidator{
		config:      &config.ParserConfig{},
		cacheClient: &MockCache{},
		colValidators: []ColValidator{
			successValidator,
			successValidator,
			successValidator,
		},
	}
	err := v.Validate("valid,row,data")
	if err != nil {
		t.Errorf("expected no error, got: %v", err)
	}
}

func TestValidate_ConcurrentValidation(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(3)

	// Create validators that track if they were called
	validatorCalled := make([]bool, 3)

	validators := make([]ColValidator, 3)
	for i := 0; i < 3; i++ {
		idx := i // Capture loop variable
		validators[i] = func(_ *RowValidatorContext, _ *[]string) error {
			defer wg.Done()
			validatorCalled[idx] = true
			// Add a small delay to ensure concurrency
			time.Sleep(10 * time.Millisecond)
			return nil
		}
	}

	v := CsvRowValidator{
		config:        &config.ParserConfig{},
		cacheClient:   &MockCache{},
		colValidators: validators,
	}

	err := v.Validate("test,row,data")

	// Wait for all validators to complete
	wg.Wait()

	if err != nil {
		t.Errorf("expected no error, got: %v", err)
	}

	// Verify all validators were called
	for i, called := range validatorCalled {
		if !called {
			t.Errorf("validator %d was not called", i)
		}
	}
}

func TestValidate_EmptyRow(t *testing.T) {
	called := false
	validator := func(_ *RowValidatorContext, cols *[]string) error {
		called = true
		if len(*cols) != 1 {
			t.Errorf("expected 1 empty column, got %d", len(*cols))
		}
		if (*cols)[0] != "" {
			t.Errorf("expected empty string, got %s", (*cols)[0])
		}
		return nil
	}

	v := CsvRowValidator{
		config:        &config.ParserConfig{Delimiter: ","},
		cacheClient:   &MockCache{},
		colValidators: []ColValidator{validator},
	}

	err := v.Validate("")

	if err != nil {
		t.Errorf("expected no error, got: %v", err)
	}

	if !called {
		t.Error("validator was not called")
	}
}

func TestValidate_DifferentDelimiter(t *testing.T) {
	validator := func(_ *RowValidatorContext, cols *[]string) error {
		if len(*cols) != 3 {
			return fmt.Errorf("expected 3 columns, got %d", len(*cols))
		}
		return nil
	}

	testCases := []struct {
		name      string
		delimiter string
		row       string
		wantErr   bool
	}{
		{
			name:      "comma delimiter",
			delimiter: ",",
			row:       "a,b,c",
			wantErr:   false,
		},
		{
			name:      "semicolon delimiter",
			delimiter: ";",
			row:       "a;b;c",
			wantErr:   false,
		},
		{
			name:      "tab delimiter",
			delimiter: "\t",
			row:       "a\tb\tc",
			wantErr:   false,
		},
		{
			name:      "pipe delimiter",
			delimiter: "|",
			row:       "a|b|c",
			wantErr:   false,
		},
		{
			name:      "wrong delimiter",
			delimiter: ",",
			row:       "a;b;c",
			wantErr:   true, // Will have 1 column instead of 3
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			v := CsvRowValidator{
				config:        &config.ParserConfig{Delimiter: tc.delimiter},
				cacheClient:   &MockCache{},
				colValidators: []ColValidator{validator},
			}

			err := v.Validate(tc.row)

			if (err != nil) != tc.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}

func TestValidate_ColumnCountValidation(t *testing.T) {
	testCases := []struct {
		name            string
		row             string
		expectedColumns int
		wantErr         bool
	}{
		{
			name:            "correct column count",
			row:             "a,b,c",
			expectedColumns: 3,
			wantErr:         false,
		},
		{
			name:            "too few columns",
			row:             "a,b",
			expectedColumns: 3,
			wantErr:         true,
		},
		{
			name:            "too many columns",
			row:             "a,b,c,d",
			expectedColumns: 3,
			wantErr:         true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			v := CsvRowValidator{
				config: &config.ParserConfig{
					Delimiter:       ",",
					ExpectedColumns: tc.expectedColumns,
				},
				cacheClient:   &MockCache{},
				colValidators: []ColValidator{isValidSize},
			}

			err := v.Validate(tc.row)

			if (err != nil) != tc.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}

// MockCacheWithTracking extends MockCache to track calls
type MockCacheWithTracking struct {
	MockCache
	GetCalled      bool
	SetCalled      bool
	SetFieldCalled bool
	GetKey         string
	SetKey         string
	SetValue       string
}

func (m *MockCacheWithTracking) Get(ctx context.Context, key string) (string, error) {
	m.GetCalled = true
	m.GetKey = key
	return "cached-value", nil
}

func (m *MockCacheWithTracking) Set(ctx context.Context, key, value string) error {
	m.SetCalled = true
	m.SetKey = key
	m.SetValue = value
	return nil
}

func (m *MockCacheWithTracking) SetField(ctx context.Context, key, field, value string) error {
	m.SetFieldCalled = true
	return nil
}

func TestValidate_CacheInteraction(t *testing.T) {
	mockCache := &MockCacheWithTracking{}

	// Create a validator that uses the cache
	cacheValidator := func(ctx *RowValidatorContext, cols *[]string) error {
		// Get something from cache
		_, err := ctx.Cache.Get(context.Background(), "test-key")
		if err != nil {
			return err
		}

		// Set something in cache
		if len(*cols) > 0 {
			return ctx.Cache.Set(context.Background(), "result-key", (*cols)[0])
		}

		return nil
	}

	v := CsvRowValidator{
		config:        &config.ParserConfig{Delimiter: ","},
		cacheClient:   mockCache,
		colValidators: []ColValidator{cacheValidator},
	}

	err := v.Validate("value1,value2,value3")

	if err != nil {
		t.Errorf("expected no error, got: %v", err)
	}

	// Verify cache interactions
	if !mockCache.GetCalled {
		t.Error("Cache Get was not called")
	}
	if mockCache.GetKey != "test-key" {
		t.Errorf("Cache Get called with wrong key, got: %s, want: %s", mockCache.GetKey, "test-key")
	}

	if !mockCache.SetCalled {
		t.Error("Cache Set was not called")
	}
	if mockCache.SetKey != "result-key" {
		t.Errorf("Cache Set called with wrong key, got: %s, want: %s", mockCache.SetKey, "result-key")
	}
	if mockCache.SetValue != "value1" {
		t.Errorf("Cache Set called with wrong value, got: %s, want: %s", mockCache.SetValue, "value1")
	}
}

func TestValidate_MultipleValidators(t *testing.T) {
	validationResults := make([]bool, 3)

	validators := []ColValidator{
		// Check first column is not empty
		func(_ *RowValidatorContext, cols *[]string) error {
			if len(*cols) == 0 || (*cols)[0] == "" {
				return fmt.Errorf("first column is empty")
			}
			validationResults[0] = true
			return nil
		},
		// Check second column is a number
		func(_ *RowValidatorContext, cols *[]string) error {
			if len(*cols) < 2 || !isNumeric((*cols)[1]) {
				return fmt.Errorf("second column is not a number")
			}
			validationResults[1] = true
			return nil
		},
		// Check third column is one of allowed values
		func(_ *RowValidatorContext, cols *[]string) error {
			if len(*cols) < 3 {
				return fmt.Errorf("missing third column")
			}

			allowedValues := []string{"A", "B", "C"}
			value := (*cols)[2]

			for _, allowed := range allowedValues {
				if value == allowed {
					validationResults[2] = true
					return nil
				}
			}

			return fmt.Errorf("third column value '%s' not in allowed list: %v", value, allowedValues)
		},
	}

	testCases := []struct {
		name    string
		row     string
		wantErr bool
	}{
		{
			name:    "all valid",
			row:     "name,123,A",
			wantErr: false,
		},
		{
			name:    "empty first column",
			row:     ",123,A",
			wantErr: true,
		},
		{
			name:    "non-numeric second column",
			row:     "name,abc,A",
			wantErr: true,
		},
		{
			name:    "invalid third column",
			row:     "name,123,X",
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Reset validation results
			for i := range validationResults {
				validationResults[i] = false
			}

			v := CsvRowValidator{
				config:        &config.ParserConfig{Delimiter: ","},
				cacheClient:   &MockCache{},
				colValidators: validators,
			}

			err := v.Validate(tc.row)

			if (err != nil) != tc.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tc.wantErr)
			}

			if !tc.wantErr {
				// If no error expected, all validators should have passed
				for i, result := range validationResults {
					if !result {
						t.Errorf("validator %d did not pass", i)
					}
				}
			}
		})
	}
}

// Helper function to check if a string is numeric
func isNumeric(s string) bool {
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return len(s) > 0
}
