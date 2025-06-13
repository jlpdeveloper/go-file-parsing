# GO File Parsing Experiments
This application is an experiment in using Go's high concurrency to parse file data and store in a redis cache. 

> [!IMPORTANT]
> This application has no real world use, it is meant to be an experiment and possibly a model for how 
> Go can be used to efficiently read large files.

## Data
The data I'm using for this experiment is from [kaggle](https://www.kaggle.com/datasets/wordsforthewise/lending-club?resource=download)
The data has the following columns:
```python

['id',
 'member_id',
 'loan_amnt',
 'funded_amnt',
 'funded_amnt_inv',
 'term',
 'int_rate',
 'installment',
 'grade',
 'sub_grade',
 'emp_title',
 'emp_length',
 'home_ownership',
 'annual_inc',
 'verification_status',
 'issue_d',
 'loan_status',
 'pymnt_plan',
 'url',
 'desc',
 'purpose',
 'title',
 'zip_code',
 'addr_state',
 'dti',
 'delinq_2yrs',
 'earliest_cr_line',
 'fico_range_low',
 'fico_range_high',
 'inq_last_6mths',
 'mths_since_last_delinq',
 'mths_since_last_record',
 'open_acc',
 'pub_rec',
 'revol_bal',
 'revol_util',
 'total_acc',
 'initial_list_status',
 'out_prncp',
 'out_prncp_inv',
 'total_pymnt',
 'total_pymnt_inv',
 'total_rec_prncp',
 'total_rec_int',
 'total_rec_late_fee',
 'recoveries',
 'collection_recovery_fee',
 'last_pymnt_d',
 'last_pymnt_amnt',
 'next_pymnt_d',
 'last_credit_pull_d',
 'last_fico_range_high',
 'last_fico_range_low',
 'collections_12_mths_ex_med',
 'mths_since_last_major_derog',
 'policy_code',
 'application_type',
 'annual_inc_joint',
 'dti_joint',
 'verification_status_joint',
 'acc_now_delinq',
 'tot_coll_amt',
 'tot_cur_bal',
 'open_acc_6m',
 'open_act_il',
 'open_il_12m',
 'open_il_24m',
 'mths_since_rcnt_il',
 'total_bal_il',
 'il_util',
 'open_rv_12m',
 'open_rv_24m',
 'max_bal_bc',
 'all_util',
 'total_rev_hi_lim',
 'inq_fi',
 'total_cu_tl',
 'inq_last_12m',
 'acc_open_past_24mths',
 'avg_cur_bal',
 'bc_open_to_buy',
 'bc_util',
 'chargeoff_within_12_mths',
 'delinq_amnt',
 'mo_sin_old_il_acct',
 'mo_sin_old_rev_tl_op',
 'mo_sin_rcnt_rev_tl_op',
 'mo_sin_rcnt_tl',
 'mort_acc',
 'mths_since_recent_bc',
 'mths_since_recent_bc_dlq',
 'mths_since_recent_inq',
 'mths_since_recent_revol_delinq',
 'num_accts_ever_120_pd',
 'num_actv_bc_tl',
 'num_actv_rev_tl',
 'num_bc_sats',
 'num_bc_tl',
 'num_il_tl',
 'num_op_rev_tl',
 'num_rev_accts',
 'num_rev_tl_bal_gt_0',
 'num_sats',
 'num_tl_120dpd_2m',
 'num_tl_30dpd',
 'num_tl_90g_dpd_24m',
 'num_tl_op_past_12m',
 'pct_tl_nvr_dlq',
 'percent_bc_gt_75',
 'pub_rec_bankruptcies',
 'tax_liens',
 'tot_hi_cred_lim',
 'total_bal_ex_mort',
 'total_bc_limit',
 'total_il_high_credit_limit',
 'revol_bal_joint',
 'sec_app_fico_range_low',
 'sec_app_fico_range_high',
 'sec_app_earliest_cr_line',
 'sec_app_inq_last_6mths',
 'sec_app_mort_acc',
 'sec_app_open_acc',
 'sec_app_revol_util',
 'sec_app_open_act_il',
 'sec_app_num_rev_accts',
 'sec_app_chargeoff_within_12_mths',
 'sec_app_collections_12_mths_ex_med',
 'sec_app_mths_since_last_major_derog',
 'hardship_flag',
 'hardship_type',
 'hardship_reason',
 'hardship_status',
 'deferral_term',
 'hardship_amount',
 'hardship_start_date',
 'hardship_end_date',
 'payment_plan_start_date',
 'hardship_length',
 'hardship_dpd',
 'hardship_loan_status',
 'orig_projected_additional_accrued_interest',
 'hardship_payoff_balance_amount',
 'hardship_last_payment_amount',
 'disbursement_method',
 'debt_settlement_flag',
 'debt_settlement_flag_date',
 'settlement_status',
 'settlement_date',
 'settlement_amount',
 'settlement_percentage',
 'settlement_term']
```
A sample file has been generated using ChatGPT. The rules each row satisfy are listed [here](sample.csv.md)

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

10. **Stable Employment**
    - `emp_length` in [5 years, 6 years, 7 years, 8 years, 9 years, 10+ years].

11. **No Public Record or Bankruptcies**
    - `pub_rec` == 0 and `pub_rec_bankruptcies` == 0 and `tax_liens` == 0.

12. **Verified with Income**
    - `verification_status` in [Source Verified, Verified] and `annual_inc` > 30,000.


## Additional Data to Store
- `avg_cur_bal`
- `application_type`
  - `annual_inc_joint` if type: `Joint App`
- `tot_coll_amt`
- `acc_now_delinq`