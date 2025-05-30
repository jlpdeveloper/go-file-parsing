# Sample CSV Data

The `sample.csv` rule has 10 records and follow these rules:

- Row 1: Passes all rules
- Row 2: Violates Rule 1 (Empty `item_id`)
- Row 3: Violates Rule 2 (`price` is not a number)
- Row 4: Violates Rule 3 (`quantity` is negative)
- Row 5: Violates Rule 4 (`category` is invalid)
- Row 6: Violates Rule 5 (`in_stock` is not 'true' or 'false')
- Row 7: Violates Rule 6 (`release_date` is not a valid date format)
- Row 8: Violates Rule 7 (`rating` is not a float between 0 and 5)
- Row 9: Violates Rule 8 (`description` is over 100 characters)
- Row 10: Passes all rules
