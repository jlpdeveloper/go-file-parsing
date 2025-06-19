# GO File Parsing Experiments

This application is an experiment in using Go's high concurrency to parse file data and store in a Redis-compatible cache (Valkey). It demonstrates efficient techniques for processing large CSV files using Go's concurrency features.

> [!IMPORTANT]
> This application has no real world use, it is meant to be an experiment and possibly a model for how 
> Go can be used to efficiently read large files.

## Project Overview

This project demonstrates:
- Concurrent CSV file parsing using goroutines
- Efficient validation of data rows against business rules
- Caching valid data in a Redis-compatible database (Valkey)
- Memory-efficient processing using object pools
- Error handling and reporting

## Project Structure

```
go-file-parsing/
├── cache/              # Cache abstraction and implementation
│   ├── cache.go        # Cache interface definition
│   ├── parser_cache.go # Valkey implementation of cache
│   └── cache_test.go   # Tests for cache functionality
├── config/             # Configuration handling
│   └── config.go       # Parser configuration
├── loan_info/          # Domain-specific validation logic
│   ├── loan_info.go    # Main validation rules
│   ├── *_validations.go # Specific validation implementations
│   └── *_test.go       # Tests for validations
├── utils/              # Utility functions
│   └── string_utils.go # String manipulation utilities
├── validator/          # Generic validation framework
│   ├── row_validator.go # Row validation logic
│   ├── column_utils.go  # Column processing utilities
│   └── map_pool.go     # Memory-efficient map pool
├── main.go             # Application entry point
├── config.json         # Parser configuration
├── dev.compose.yml     # Docker Compose for development
└── sample.csv          # Sample data file
```

## How It Works

1. The application reads a CSV file line by line
2. For each row, it:
   - Allocates a validator from a pool
   - Validates the row concurrently using multiple validation rules
   - Stores valid data in the cache
   - Collects and reports validation errors
3. After processing, it cleans up any invalid data from the cache

The application uses Go's concurrency primitives (goroutines, channels, wait groups, and errgroup) to process rows efficiently.

## Dependencies

- Go 1.24 or later
- [valkey-io/valkey-go](https://github.com/valkey-io/valkey-go) v1.0.60 - Redis-compatible client library
- [golang.org/x/sync](https://pkg.go.dev/golang.org/x/sync) v0.15.0 - Additional synchronization primitives
- [Valkey](https://valkey.io/) - Redis-compatible database (via Docker)

## Setup and Installation

### Prerequisites

- Go 1.24 or later (for local development)
- Docker and Docker Compose (for running with Docker)

### Option 1: Running with Docker

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/go-file-parsing.git
   cd go-file-parsing
   ```

2. Build and run the application with Docker Compose:
   ```bash
   docker-compose up -d
   ```

   This will start both the Valkey container and the application container.

3. To use the sample file instead of the large dataset, modify the `docker-compose.yml` file:
   ```yaml
   app:
     # ... other settings ...
     command: ["sample.csv"]
   ```

4. To view logs:
   ```bash
   docker-compose logs -f app
   ```

### Option 2: Running Locally

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/go-file-parsing.git
   cd go-file-parsing
   ```

2. Start the Valkey container:
   ```bash
   docker-compose -f dev.compose.yml up -d
   ```

3. Set the environment variable for Valkey:
   ```bash
   # For Windows PowerShell
   $env:VALKEY_URLS = "localhost:6379"

   # For Linux/macOS
   export VALKEY_URLS="localhost:6379"
   ```

4. Run the application:
   ```bash
   go run main.go
   ```

## Usage

The application is configured via the `config.json` file:

```json
{
  "HasHeader": true,
  "Delimiter": ",",
  "ExpectedColumns": 156
}
```

- `HasHeader`: Set to true if the CSV file has a header row
- `Delimiter`: The character used to separate columns
- `ExpectedColumns`: The expected number of columns in each row

To use your own CSV file, modify the `parseFile` function call in `main.go`:

```
// In main.go
parseFile("your-file.csv", cacheClient)
```

## Performance Considerations

- The application uses a pool of validators to limit memory usage
- Each row is processed concurrently, with validation rules applied in parallel
- The map pool pattern is used to reduce garbage collection pressure
- Buffer sizes are configurable to balance memory usage and performance

## Reuse considerations
If you wish to reuse this project, here are some considerations to help you adapt it to your needs:

### Extracting Validations as a Separate Package

The validation framework in this project is designed to be modular and reusable. You could extract the validation logic into its own package or even a separate library:

1. **Core Validation Framework**: The `validator` package contains the core validation framework, including:
   - `RowValidator` interface and `CsvRowValidator` implementation
   - `ColValidator` function type for individual column validations
   - Map pooling for efficient memory usage

2. **Domain-Specific Validations**: The `loan_info` package contains domain-specific validations that could be moved to a separate package:
   - Validation functions like `hasValidLoanAmount`, `hasValidInterestRate`, etc.
   - Error definitions specific to loan data validation
   - The validator registration in `loan_info.go`

### Adding Different File-Based Validators

To add support for different file types or data domains:

1. **Create a New Domain Package**: Similar to the `loan_info` package, create a new package for your domain:
   ```
   your-domain/
   ├── domain.go             # Register validators and create validator pool
   ├── errors.go             # Define domain-specific errors
   ├── validations.go        # Implement domain-specific validation functions
   └── row_shape_validations.go # Basic structure validations
   ```

2. **Implement Validation Functions**: Create functions that follow the `ColValidator` signature:
   ```
   // Example validation function
   func yourValidationFunction(ctx *validator.RowValidatorContext, cols []string) (map[string]string, error) {
       // Validation logic here

       // If validation passes, optionally return data to cache
       result := ctx.GetMap()
       result["key"] = "value"

       return result, nil // or return nil, yourError
   }
   ```

3. **Register Validators**: Create a slice of validators in your domain package:
   ```
   // Example validator registration
   var validators = []validator.ColValidator{
       isValidSize,
       yourValidationFunction1,
       yourValidationFunction2,
       // ...
   }
   ```

4. **Create Validator Pool**: Implement a function to create a pool of validators:
   ```
   // Example validator pool creation
   func NewRowValidatorPool(conf *config.ParserConfig, cache cache.DistributedCache, poolSize int) chan validator.CsvRowValidator {
       pool := make(chan validator.CsvRowValidator, poolSize)
       for i := 0; i < poolSize; i++ {
           pool <- validator.New(conf, cache, validators)
       }
       return pool
   }
   ```

5. **Update Main Application**: Modify `main.go` to use your new validator pool:
   ```
   // Example usage in main.go
   pool := your_domain.NewRowValidatorPool(&conf, cacheClient, 100)
   ```

### Best Practices for Extension

1. **Keep Validators Focused**: Each validator function should validate one specific aspect of the data.

2. **Use Error Constants**: Define error constants in an `errors.go` file for consistent error messages.

3. **Reuse Maps**: Always use the map pool (`ctx.GetMap()`) to get maps for returning data.

4. **Concurrent Safety**: Ensure your validators are safe for concurrent use.

5. **Testing**: Write comprehensive tests for each validator function.

6. **Configuration**: Use the configuration system to make your validators configurable.

7. **Documentation**: Document your validators and their expected behavior.

By following these guidelines, you can extend this project to handle different types of file parsing and validation while maintaining its performance and memory efficiency.

## Data
The data I'm using for this experiment is from [kaggle](https://www.kaggle.com/datasets/wordsforthewise/lending-club?resource=download)
The data has the following columns:

| Field Name                   | Description/Value    |
|------------------------------|----------------------|
| id                           |                      |
| member_id                    |                      |
| loan_amnt                    |                      |
| funded_amnt                  |                      |
| funded_amnt_inv              |                      |
| term                         |                      |
| int_rate                     |                      |
| installment                  |                      |
| grade                        |                      |
| sub_grade                    |                      |
| emp_title                    |                      |
| emp_length                   |                      |
| home_ownership               |                      |
| annual_inc                   |                      |
| verification_status          |                      |
| issue_d                      |                      |
| loan_status                  |                      |
| pymnt_plan                   |                      |
| url                          |                      |
| desc                         |                      |
| purpose                      |                      |
| title                        |                      |
| zip_code                     |                      |
| addr_state                   |                      |
| dti                          |                      |
| delinq_2yrs                  |                      |
| earliest_cr_line             |                      |
| fico_range_low               |                      |
| fico_range_high              |                      |
| inq_last_6mths               |                      |
| mths_since_last_delinq       |                      |
| mths_since_last_record       |                      |
| open_acc                     |                      |
| pub_rec                      |                      |
| revol_bal                    |                      |
| revol_util                   |                      |
| total_acc                    |                      |
| initial_list_status          |                      |
| out_prncp                    |                      |
| out_prncp_inv                |                      |
| total_pymnt                  |                      |
| total_pymnt_inv              |                      |
| total_rec_prncp              |                      |
| total_rec_int                |                      |
| total_rec_late_fee           |                      |
| recoveries                   |                      |
| collection_recovery_fee      |                      |
| last_pymnt_d                 |                      |
| last_pymnt_amnt              |                      |
| next_pymnt_d                 |                      |
| last_credit_pull_d           |                      |
| last_fico_range_high         |                      |
| last_fico_range_low          |                      |
| collections_12_mths_ex_med   |                      |
| mths_since_last_major_derog  |                      |
| policy_code                  |                      |
| application_type             |                      |
| annual_inc_joint             |                      |
| dti_joint                    |                      |
| verification_status_joint    |                      |
| acc_now_delinq               |                      |
| tot_coll_amt                 |                      |
| tot_cur_bal                  |                      |
| open_acc_6m                  |                      |
| open_act_il                  |                      |
| open_il_12m                  |                      |
| open_il_24m                  |                      |
| mths_since_rcnt_il           |                      |
| total_bal_il                 |                      |
| il_util                      |                      |
| open_rv_12m                  |                      |
| open_rv_24m                  |                      |
| max_bal_bc                   |                      |
| all_util                     |                      |
| total_rev_hi_lim             |                      |
| inq_fi                       |                      |
| total_cu_tl                  |                      |
| inq_last_12m                 |                      |
| acc_open_past_24mths         |                      |
| avg_cur_bal                  |                      |
| bc_open_to_buy               |                      |
| bc_util                      |                      |
| chargeoff_within_12_mths     |                      |
| delinq_amnt                  |                      |
| mo_sin_old_il_acct           |                      |
| mo_sin_old_rev_tl_op         |                      |
| mo_sin_rcnt_rev_tl_op        |                      |
| mo_sin_rcnt_tl               |                      |
| mort_acc                     |                      |
| mths_since_recent_bc         |                      |
| mths_since_recent_bc_dlq     |                      |
| mths_since_recent_inq        |                      |
| mths_since_recent_revol_delinq|                     |
| num_accts_ever_120_pd        |                      |
| num_actv_bc_tl               |                      |
| num_actv_rev_tl              |                      |
| num_bc_sats                  |                      |
| num_bc_tl                    |                      |
| num_il_tl                    |                      |
| num_op_rev_tl                |                      |
| num_rev_accts                |                      |
| num_rev_tl_bal_gt_0          |                      |
| num_sats                     |                      |
| num_tl_120dpd_2m             |                      |
| num_tl_30dpd                 |                      |
| num_tl_90g_dpd_24m           |                      |
| num_tl_op_past_12m           |                      |
| pct_tl_nvr_dlq               |                      |
| percent_bc_gt_75             |                      |
| pub_rec_bankruptcies         |                      |
| tax_liens                    |                      |
| tot_hi_cred_lim              |                      |
| total_bal_ex_mort            |                      |
| total_bc_limit               |                      |
| total_il_high_credit_limit   |                      |
| revol_bal_joint              |                      |
| sec_app_fico_range_low       |                      |
| sec_app_fico_range_high      |                      |
| sec_app_earliest_cr_line     |                      |
| sec_app_inq_last_6mths       |                      |
| sec_app_mort_acc             |                      |
| sec_app_open_acc             |                      |
| sec_app_revol_util           |                      |
| sec_app_open_act_il          |                      |
| sec_app_num_rev_accts        |                      |
| sec_app_chargeoff_within_12_mths|                   |
| sec_app_collections_12_mths_ex_med|                |
| sec_app_mths_since_last_major_derog|               |
| hardship_flag                |                      |
| hardship_type                |                      |
| hardship_reason              |                      |
| hardship_status              |                      |
| deferral_term                |                      |
| hardship_amount              |                      |
| hardship_start_date          |                      |
| hardship_end_date            |                      |
| payment_plan_start_date      |                      |
| hardship_length              |                      |
| hardship_dpd                 |                      |
| hardship_loan_status         |                      |
| orig_projected_additional_accrued_interest|         |
| hardship_payoff_balance_amount|                     |
| hardship_last_payment_amount |                      |
| disbursement_method          |                      |
| debt_settlement_flag         |                      |
| debt_settlement_flag_date    |                      |
| settlement_status            |                      |
| settlement_date              |                      |
| settlement_amount            |                      |
| settlement_percentage        |                      |
| settlement_term              |                      |


A sample file has been generated using ChatGPT.

## Rules
The below rules were generated as ways to determine if data that is being parsed in is "good" or "bad." They were generated
using ChatGPT to analyze the columns and give recommendations on rules to add complexity to the parsing


### ✅ Data Validation Rules

1. **Valid Loan and Funding** 
    - `loan_amnt` > 0 and `funded_amnt` == `funded_amnt_inv`.

2. **Reasonable Interest Rate** 
    - `int_rate` between 5% and 35%.

3. **Valid Grade/Subgrade** 
    - `grade` in [A–G], `sub_grade` matches pattern like `B3`.

4. **Valid Term** 
   - `term` is between 12 and 72 months

5. **Has Employment Info**
    - Non-empty `emp_title` and `emp_length` is not null.

6. **Low DTI and Home Ownership**
    - `dti` < 20, `home_ownership` in [MORTGAGE, OWN], and `annual_inc` > 40,000.

7. **Established Credit History**
    - `earliest_cr_line` not null and is > 10 years ago.

8. **Healthy FICO Score**
    - `fico_range_low` >= 660 and `fico_range_high` <= 850.

9. **Has Sufficient Accounts**
   - `total_acc` >= 5 and `open_acc` >= 2.

10. **Valid Income**
    - `annual_inc` > 30,000.

## Additional Data to Store
- `avg_cur_bal`
- `application_type`
  - `annual_inc_joint` if type: `Joint App`
- `tot_coll_amt`
- `acc_now_delinq`

## Contributing

This project is an experiment and demonstration. Feel free to fork it and adapt it to your needs. If you have suggestions for improvements, please open an issue or submit a pull request.

## License

This project is licensed under the terms found in the LICENSE file in the root of this repository.
