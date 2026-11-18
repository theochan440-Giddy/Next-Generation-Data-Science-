# Gen AI Prompt Log: China CO2 Emissions Analysis in Go

**Student:** THET HTUT ZAW  |  **Course task:** China CO2 emissions 1970-2022 with Go and Gen AI
**Local AI tool:** Ollama, model `gemma4:e2b`  |  **Other AI tool:** Claude (Anthropic), used to plan the pipeline, draft the Go code and write the prompts
**Date:** 18 November 2026 (Singapore time)

> Ollama was not trained. I used it as a reviewer for my Go code. It cannot see my files, so every prompt has the real column names or numbers typed in. I checked each answer against my code and data before using it.

---

## 1. Summary

**23 prompts were sent to Ollama, in 7 topics.** The final numbers (coefficients, errors, CSV files) come from my Go code, not from Ollama.

| What the prompt led to | Prompts | Count |
|---|---|---|
| **Changed my Go code or tests** | 1b, 2a, 2d | 3 |
| **Confirmed expected values or a decision** | 4a, 4b, 4c, 4d, 4e, 5a, 5a2 | 7 |
| **Used for slide wording or explanation** after checking | 1a, 6a, 6b, 6b2, 6c, 7a | 6 |
| **Not used** (vague, wrong, or no change needed) | 2b, 2c, 3a, 3b, 3c, 5b, 7b | 7 |
| **Total** | | **23** |

**Right and wrong, by topic**

| Topic | Right (checked and used) | Wrong or weak (how I found out) |
|---|---|---|
| 1 Dataset | Explained columns; its "NaN" check found a real bug (Go accepts "NaN" and "Inf") | Said CO2e means all greenhouse gases, but my file is CO2 only (I read the indicator name) |
| 2 Cleaning code | Found 2 real crashes; two specific questions both right | **Said "No problems found"** when there were 2 (I wrote tests) |
| 3 Regression maths | Method and hand example correct | **2 of 3 suggested tests wrong** (I ran them); rated a risk "High" that is small |
| 4 Metrics | MAE, RMSE, MAPE, R2 all correct | First long prompt gave no output |
| 5 Split | "Split by time" correct | Suggested a shorter test set; I kept my own |
| 6 Results | Negative R2 explained correctly | **Repeated my wrong premise** (said over-predicts; the line under-predicts) |
| 7 Slide wording | Period growth bullet correct | Changed "1.7% a year" into "1.7% in total" |

**What I decided myself:** train 1970-2012 and test 2013-2022 (not its split); a year-range check in `Clean`; one simple regression model; no forecast; my own unit tests instead of Ollama's.

**Final code checks:** 9 unit tests pass; the line's coefficients match Python numpy to 4 decimal places (intercept -459.9903, slope 1873.8353 per decade).

---

## 2. How I broke the tasks down for Ollama

`gemma4:e2b` is a small model running on my laptop. It gave poor or empty answers to long prompts, so I changed how I asked:

| Problem I met | What I changed | Result |
|---|---|---|
| A 27B model crashed on my graphics card | Moved to the smaller `gemma4:e2b` | It runs |
| The model showed only its "thinking" and no final answer | Added "Answer immediately. No thinking." and a final command line; sent a short follow-up when it stopped (4e, 5a2, 6b2) | Final answers appeared |
| One long prompt with four metrics gave no output (Prompt 4) | One metric per prompt, with numbered steps (4a to 4e) | All four correct |
| One long prompt with a results table and three questions (Prompt 6) | Typed in only the numbers needed; three short prompts (6a, 6b, 6c) | Answers on topic |
| The open question "any problems?" gave a wrong "none" (2c) | Specific questions with exact inputs (2d) | Both answers right |
| A big function in one message | Sent the code in small pieces, one task each (2a, 2c) | Fewer misses |
| It cannot see my files | Typed the real column names and numbers into every prompt | No invented numbers |

**Prompt structure I used:** ROLE (who the model is) → CONTEXT (real column names or numbers) → TASK (one job) → RULES (length limit, "use only these numbers") → FORMAT (for example "numbered list only") → a final command ("List the crashes now"). I started a new chat before each new topic so earlier answers could not confuse the model.

**How I handled every answer:** prompt → read the answer → run a test or compare with the CSV → keep, fix or reject → write it down here.

**Earlier attempts (not logged):** before `gemma4:e2b` I tried a larger model (`qwen3.8:27b`). It crashed on my graphics card and often returned only "thinking" text, so I discarded those runs and restarted with the small model. The NaN/Inf bug was first noticed in an earlier run of Prompt 1 on a larger model. I did not keep that answer, so the evidence for it is my test `TestClean`.

---

## 3. When I used Ollama

All prompts below were run on **18 November 2026**, in the order shown, with `gemma4:e2b`. Ollama does not time-stamp its chats, so the times are when I saved each answer into this project (Singapore time).

| Topic | Prompts | Answer saved |
|---|---|---|
| 1 Dataset | 1a, 1b | 17:48 |
| 2 Cleaning code | 2a, 2b, 2c, 2d | 17:53 |
| 3 Regression maths | 3a, 3b, 3c | 17:55 |
| 4 Metrics | 4a, 4b, 4c, 4d, 4e | 18:00 |
| 5 Split | 5a, 5a2, 5b | 18:07 |
| 6 Results | 6a, 6b, 6b2, 6c | 18:14 |
| 7 Slide wording | 7a, 7b | 18:15 |

---

## 4. The prompts, Ollama's answers and my checks

Each entry shows the exact prompt, Ollama's answer as pasted, and what I did with it.


### Topic 1: Understanding the dataset  (answers saved 17:48)

#### Prompt 1a: Explain the dataset columns
**Result:** Used (with 2 gaps)

**Prompt sent:**
```text
ROLE: You are a data analyst who explains things simply to a student.

IMPORTANT: No file is attached. Use ONLY the column names and values I type below. Do NOT invent extra columns or numbers.

CONTEXT: I am analysing a CSV in Go. It has 53 rows, one per year from 1970 to 2022.
Columns: Country Name, Country Code, Indicator Name, Indicator Code, Year, CO2_Mt_CO2e
First two rows:
China,CHN,Carbon dioxide (CO2) emissions (total) excluding LULUCF (Mt CO2e),EN.GHG.CO2.MT.CE.AR5,1970,911.5013
China,CHN,Carbon dioxide (CO2) emissions (total) excluding LULUCF (Mt CO2e),EN.GHG.CO2.MT.CE.AR5,1971,914.9404

TASK: Explain what each column means, and what "Mt CO2e" and "LULUCF" mean.

RULES:
- Simple English, under 120 words.
- One bullet per column.

FORMAT: Bullet points only. Do not repeat my question.

Explain the columns now.
```
**Ollama answer:**
```text
* **Country Name:** This tells you which country the data belongs to.
* **Country Code:** This is the official code used to identify the country (like CHN for China).
* **Indicator Name:** This describes exactly what the measurement is (e.g., total CO2 emissions).
* **Indicator Code:** This is a specific technical label for that exact type of data.
* **Year:** This shows the specific year the emissions data was recorded.
* **CO2\_Mt\_CO2e:** This is the actual amount of carbon dioxide emissions, measured in Megatonnes of CO2 equivalents.

**Definitions:**

* **Mt CO2e:** Stands for Megatonnes of Carbon Dioxide equivalents. This is a way to measure all greenhouse gases as if they were CO2.
* **LULUCF:** Stands for Land Use, Land Use Change, and Forestry. This accounts for emissions and removals that happen due to changes in land use (like forests).
```
**My check:** Right: explained columns and LULUCF. Gaps I corrected: it said CO2e means all greenhouse gases, but my file is CO2 only; it did not say the file excludes LULUCF (the indicator name states this).

#### Prompt 1b: List data-quality checks
**Result:** Changed code

**Prompt sent:**
```text
TASK: Using the same CSV, list the data-quality checks my Go program should run before modelling.

RULES:
- At most 8 checks.
- One line each: the check, then what problem it catches.
- Only checks that make sense for 53 yearly rows.

FORMAT: Numbered list only. Do not repeat my question.

List the checks now.
```
**Ollama answer:**
```text
1. Check for missing values (NaN or nulls) in the `CO2_Mt_CO2e` column to identify incomplete records.
2. Validate that the `Year` column contains valid integer values within the expected range (1970 to 2022).
3. Verify that the `CO2_Mt_CO2e` values are correctly formatted as numbers and are not corrupted by text strings.
4. Ensure temporal sequence by checking that the `Year` column is strictly increasing to prevent time-series errors.
5. Identify and flag duplicate rows based on the combination of `Country Name`, `Country Code`, and `Year`.
6. Check the plausibility of emissions data by flagging any negative CO2 values, which are physically impossible.
7. Verify that essential identification columns (`Country Name` and `Country Code`) are present and non-empty for every row.
8. Perform a range check on `CO2_Mt_CO2e` to flag values that are extremely high or unusually low compared to historical norms.
```
**My check:** Right: sensible checks. Its "NaN" check led me to test Go: `strconv.ParseFloat` accepts "NaN" and "Inf", so `Clean` would let them through. Fixed and tested. Vague check 8 ("unusual values"): I used my own range, 0 to 50,000 Mt.


### Topic 2: Reviewing the cleaning code  (answers saved 17:53)

#### Prompt 2a: Find crashes in the row-reading loop
**Result:** Changed code

**Prompt sent:**
```text
ROLE: You are a Go code reviewer.

CONTEXT: This Go loop reads rows of a CSV with columns Year and CO2_Mt_CO2e. The row is a []string. columnIndex returns -1 if a column name is not found in the header.

CODE:
yearCol := columnIndex(header, "Year")
valueCol := columnIndex(header, "CO2_Mt_CO2e")
for _, row := range rows {
	yearText := strings.TrimSpace(row[yearCol])
	valueText := strings.TrimSpace(row[valueCol])
	if yearText == "" || valueText == "" {
		continue
	}
}

TASK: Name any way this code can crash (panic) at runtime.

RULES:
- At most 3 problems.
- One sentence each, naming the exact variable.
- Do not suggest fixes yet.

FORMAT: Numbered list only.

List the crashes now.
```
**Ollama answer:**
```text
1. The line `yearText := strings.TrimSpace(row[yearCol])` will panic if `yearCol` is -1, as it attempts to access an index out of bounds in the `row` slice.
2. The line `valueText := strings.TrimSpace(row[valueCol])` will panic if `valueCol` is -1, as it attempts to access an index out of bounds in the `row` slice.
3. If the returned index from `columnIndex` is greater than or equal to the length of `row`, both lines will panic due to an index out of range error.
```
**My check:** Right: missing column (index -1) and short row both crash. I proved each with `TestCleanBrokenInput` (`index out of range`). Fixed in `Clean`.

#### Prompt 2b: Suggest fixes for those crashes
**Result:** Not used

**Prompt sent:**
```text
TASK: Now give the smallest fix for each of the 3 crashes you listed.

RULES:
- One sentence per fix.
- Show Go code only if needed, max 4 lines per fix.
- The code must compile and use only the standard library.
- Fixes 1 and 2 may share the same code. If so, say so.

FORMAT: Same numbering as your list (1, 2, 3).

Give the fixes now.
```
**Ollama answer:**
```text
1. Before attempting to access the row elements, check if `yearCol` and `valueCol` are valid indices (not -1).
2. You can check both indices before the loop starts to ensure you have valid column references.
3. Inside the loop, before accessing `row`, ensure that both `yearCol` and `valueCol` are non-negative and less than the length of `row`.
```
**My check:** Correct idea but vague, no code, and no instruction for what to do on failure. I wrote the fixes myself: report and stop for a missing column; count and skip for a short row.

#### Prompt 2c: Review the second half of `Clean`
**Result:** Not used (wrong)

**Prompt sent:**
```text
TASK: Review this second part of the same Go function. Name any logic problems in it.

CONTEXT: Each row has a year and a CO2 value. The loop skips a row if the year was already seen, or if the value is negative. After the loop, the records are sorted by year and any missing years are listed.

CODE:
if value < 0 {
	report.NegativeValues++
	continue
}
if seen[year] {
	report.DuplicateYears++
	continue
}
seen[year] = true
records = append(records, Record{Year: year, CO2: value})

sort.Slice(records, func(i, j int) bool { return records[i].Year < records[j].Year })
for i := 1; i < len(records); i++ {
	for y := records[i-1].Year + 1; y < records[i].Year; y++ {
		report.MissingYears = append(report.MissingYears, y)
	}
}

RULES:
- At most 3 problems.
- One sentence each, naming the exact line.
- Do not suggest fixes yet.
- If you find no problem, say "No problems found".

FORMAT: Numbered list only.

List the problems now.
```
**Ollama answer:**
```text
No problems found.
```
**My check:** **Wrong: "No problems found."** The code had 2 problems (see 2d). An open question gave the small model an easy way out.

#### Prompt 2d: Two specific questions about `Clean`
**Result:** Changed code

**Prompt sent:**
```text
TASK: Answer these two questions about the same code. Think through each one step by step.

QUESTION 1: A file has two rows with the same year 2001. The first row has CO2 = 10. The second row has CO2 = -5.
Walk through the code for the SECOND row. Which counter increases: NegativeValues or DuplicateYears?

QUESTION 2: The cleaned records are 2022 and then 2050 (a typing mistake for 2025).
Walk through the gap loop. How many years are added to report.MissingYears?

RULES:
- Keep each answer under 40 words.
- Give the final number or counter name at the end of each answer.

FORMAT: Q1 and Q2, one short paragraph each.

Answer both questions now.
```
**Ollama answer:**
```text
**QUESTION 1:** For the second row with a CO2 value of -5, the condition `if value < 0` will be true. Therefore, the counter `NegativeValues` will increase.

**QUESTION 2:** The gap loop will count the years from 2023 up to 2049, because the condition is strictly less than the end year (2050). This results in 27 years being added to `report.MissingYears`.
```
**My check:** Right on both, confirmed by test: a repeated year with a negative value is counted as negative (`dup=0 neg=1`); a year 2050 adds 27 fake missing years. Fixed: duplicates are checked first; years outside 1900-2100 are rejected.


### Topic 3: Checking the regression maths  (answers saved 17:55)

#### Prompt 3a: Is least squares correct for this model?
**Result:** Not used

**Prompt sent:**
```text
ROLE: You are a statistics lecturer.

IMPORTANT: No file is attached. Do NOT describe or invent a dataset.

CONTEXT: I wrote ordinary least squares (OLS) linear regression in Go from scratch. My method: build the matrices X^T X and X^T y, then solve (X^T X) b = X^T y with Gaussian elimination with partial pivoting. X has a column of 1s for the intercept. I use it on 53 yearly values with at most 3 coefficients (intercept, t, t squared).

TASK: Answer these two questions.
Q1: Is this method correct for OLS? Answer yes or no, then explain in 2 sentences.
Q2: Give 3 situations where this method can fail or be inaccurate. For each one, say if it is a risk for my 53 rows and 3 coefficients.

RULES:
- Under 150 words in total.
- Simple English. Explain any maths term in one line.

FORMAT: Q1 and Q2 as headings. Do not repeat my question.

Answer Q1 and Q2 now.
```
**Ollama answer:**
```text
**Q1: Is this method correct for OLS?**
Yes. Solving the Normal Equations, $(X^T X)b = X^T y$, is a fundamental and mathematically correct approach to finding the Ordinary Least Squares (OLS) estimates. This method finds the coefficient vector $b$ that minimizes the sum of squared residuals (the differences between the predicted and actual values).

**Q2: Give 3 situations where this method can fail or be inaccurate.**
1. **Multicollinearity:** If your predictor variables (t and $t^2$) are highly correlated, the matrix $X^T X$ becomes ill-conditioned, making the solution highly sensitive and inaccurate. (Risk: High, especially with polynomial terms).
2. **Singular Matrix:** If $X^T X$ is singular (non-invertible), the system cannot be solved, and the regression fails completely. (Risk: Moderate, possible if columns are perfectly dependent).
3. **Numerical Instability:** Gaussian elimination with partial pivoting can still suffer from numerical instability, leading to inaccurate results, particularly if the matrix is poorly conditioned. (Risk: Low to Moderate, but present in custom implementations).
```
**My check:** Right: the method is correct (my Go and numpy coefficients match). Overstated: it rated multicollinearity "High"; the real condition number is only about 1,221.

#### Prompt 3b: Check the matrix-building code
**Result:** Not used

**Prompt sent:**
```text
TASK: Check this Go code, which builds X^T X and X^T y. Is it correct?

CODE:
XtX := make([][]float64, k)
for i := range XtX {
	XtX[i] = make([]float64, k)
}
Xty := make([]float64, k)
for r := range X {
	for i := 0; i < k; i++ {
		Xty[i] += X[r][i] * y[r]
		for j := 0; j < k; j++ {
			XtX[i][j] += X[r][i] * X[r][j]
		}
	}
}

CONTEXT: X has one row per year and k columns. y has one value per year. k is the number of coefficients.

RULES:
- First work out by hand: X = [[1,0],[1,1],[1,2]] and y = [1,3,5]. What are XtX and Xty? Show the numbers.
- Then say if the code gives the same numbers.
- Finish with one line: "Verdict: correct" or "Verdict: mistake in line ...".
- Under 150 words.

Do the check now.
```
**Ollama answer:**
```text
**Manual Calculation:**
Given X = [[1,0],[1,1],[1,2]] and y = [1,3,5] (k=2):

* X^T X = [[3, 3], [3, 5]]
* X^T y = [9, 13]

**Code Check:**
The code iterates through each row r and then calculates the required sums. The inner loop correctly computes the matrix product X^T X and the vector product X^T y.

Verdict: correct
```
**My check:** Right: X^T X = [[3,3],[3,5]], X^T y = [9,13]. But it only stated the result and did not trace the loops, so the reasoning was thin. (This matrix code was later removed.)

#### Prompt 3c: Suggest unit tests
**Result:** Not used (2 of 3 wrong)

**Prompt sent:**
```text
TASK: Suggest 3 unit tests for my Go function FitOLS(X, y), which returns the regression coefficients, or an error.

RULES:
- Each test uses tiny data I can check by hand.
- Give the input X and y, and the exact expected coefficients.
- Test 3 must be a case where FitOLS should return an error.
- Describe each test in 3 lines or less. No full Go code.

FORMAT: Numbered list: name, input, expected result.

Suggest the 3 tests now.
```
**Ollama answer:**
```text
1. **Test 1: Basic Linear Regression.**
   X: [[1, 1], [1, 2]]
   y: [1, 3]
   Expected: Intercept = 1.0, Slope = 2.0

2. **Test 2: Zero Slope Fit.**
   X: [[1, 0], [1, 0]]
   y: [5, 5]
   Expected: Intercept = 5.0, Slope = 0.0

3. **Test 3: Test for Singularity/Error.**
   X: [[1, 1], [1, 1]] (Perfect multicollinearity)
   y: [1, 2]
   Expected: Error (indicating a singular matrix).
```
**My check:** **Test 1 wrong:** the line through (1,1) and (2,3) has intercept -1, not 1. **Test 2 wrong:** an all-zero column is singular, so the answer is an error, not [5, 0]. Test 3 (identical columns give an error) was right. Found only by running them. I wrote my own tests.


### Topic 4: Testing the evaluation metrics  (answers saved 18:00)

#### Prompt 4a: MAE on a small example
**Result:** Confirmed values

**Prompt sent:**
```text
Actual values: 10, 20, 30.
Predicted values: 12, 18, 30.

Step 1: Write the error for each pair (actual minus predicted).
Step 2: Write the absolute error for each pair.
Step 3: Calculate MAE = the average of the absolute errors.

Give the three steps as a short list with the numbers.
```
**Ollama answer:**
```text
1. Error 1: 10 - 12 = -2
2. Error 2: 20 - 18 = 2
3. Error 3: 30 - 30 = 0

4. Absolute Errors: |-2| = 2, |2| = 2, |0| = 0.

5. MAE: (2 + 2 + 0) / 3 = 4 / 3 ≈ 1.333
```
**My check:** Correct: errors -2, 2, 0; MAE = 1.333.

#### Prompt 4b: RMSE on the same example
**Result:** Confirmed values

**Prompt sent:**
```text
Actual: 10, 20, 30. Predicted: 12, 18, 30. The errors are -2, 2 and 0.

Step 1: Square each error. Write the three numbers.
Step 2: Add them up and divide by 3. This is the MSE.
Step 3: Take the square root of the MSE. Round to 3 decimal places. This is the RMSE.

Write the three steps with the numbers.
```
**Ollama answer:**
```text
1. Squared Errors: (-2)^2 = 4, (2)^2 = 4, (0)^2 = 0.
2. MSE: (4 + 4 + 0) / 3 = 8 / 3 ≈ 2.667.
3. RMSE: sqrt(2.667) ≈ 1.633.
```
**My check:** Correct: MSE 2.667, RMSE 1.633. A longer version asking for RMSE and MAPE together stopped early.

#### Prompt 4c: MAPE on the same example
**Result:** Confirmed values

**Prompt sent:**
```text
Actual: 10, 20, 30. Predicted: 12, 18, 30. The absolute errors are 2, 2 and 0.

Step 1: For each pair, divide the absolute error by the actual value, then multiply by 100. Write the three numbers.
Step 2: Take the average of the three numbers. This is the MAPE, in percent.

Write the two steps with the numbers.
```
**Ollama answer:**
```text
1. Percentage Errors: (2/10 x 100) = 20.0%, (2/20 x 100) = 10.0%, (0/30 x 100) = 0.0%.
2. MAPE: (20.0 + 10.0 + 0.0) / 3 = 30.0 / 3 = 10.0%.
```
**My check:** Correct: 20%, 10%, 0%, MAPE = 10%.

#### Prompt 4d: R2, steps 1-3
**Result:** Confirmed values

**Prompt sent:**
```text
Actual: 10, 20, 30. Predicted: 12, 18, 30. The squared errors are 4, 4 and 0.

Step 1: Calculate the mean (average) of the actual values.
Step 2: For each actual value, subtract the mean and square the result. Write the three numbers.
Step 3: Add the three numbers from Step 2. Call it SS_total.
Step 4: Add the three squared errors. Call it SS_error.
Step 5: Calculate R2 = 1 - SS_error / SS_total. Round to 3 decimal places.

Write the five steps with the numbers.
```
**Ollama answer:**
```text
1. Mean of Actual Values: (10 + 20 + 30) / 3 = 20.0.
2. Squared Deviations: (10 - 20)^2 = (-10)^2 = 100, (20 - 20)^2 = 0^2 = 0, (30 - 20)^2 = 10^2 = 100.
3. Sum of Squared Deviations (SS_total): 100 + 0 + 100 = 200.
```
**My check:** Correct: mean 20, SS_total 200. The reply stopped before steps 4 and 5.

#### Prompt 4e: R2, steps 4-5 (follow-up)
**Result:** Confirmed values

**Prompt sent:**
```text
Actual: 10, 20, 30. Predicted: 12, 18, 30. SS_total = 200.

Step 4: The errors are 10-12, 20-18 and 30-30. Square each one and add them up. Call the result SS_error.
Step 5: Calculate R2 = 1 - SS_error / SS_total. Round to 3 decimal places.

Write the two steps with the numbers.
```
**Ollama answer:**
```text
1. Sum of Squared Errors (SS_error): (-2)^2 + (2)^2 + (0)^2 = 4 + 4 + 0 = 8.
2. R^2: 1 - 8 / 200 = 1 - 0.04 = 0.960.
```
**My check:** Correct: SS_error 8, R2 = 0.960. These are the expected values in `TestEvaluate`; I had worked them out by hand first.


### Topic 5: Choosing the train/test split  (answers saved 18:07)

#### Prompt 5a: Random or time-based split?
**Result:** Confirmed decision

**Prompt sent:**
```text
Answer immediately. Do not show any analysis or thinking.

Yearly data from 1970 to 2022. I want to predict future values.
Should I split train and test data randomly, or by time?

Reply with exactly two lines:
Line 1: Random or By time
Line 2: One short reason
```
**Ollama answer:**
```text
(The model showed only its thinking and stopped before the final lines, so no final answer was given. I sent 5a2.)
```
**My check:** No final answer this time (thinking only). I sent the follow-up 5a2.

#### Prompt 5a2: Follow-up: final answer only
**Result:** Confirmed decision

**Prompt sent:**
```text
Now write the final answer only. Two lines, no thinking:
Line 1: Random or By time
Line 2: One short reason
```
**Ollama answer:**
```text
By time
This preserves the temporal order, preventing future data from leaking into the training set.
```
**My check:** Right: split by time, because a random split lets the model learn from future years.

#### Prompt 5b: Which years for train and test?
**Result:** Not used

**Prompt sent:**
```text
Answer immediately. No thinking.

I have yearly data from 1970 to 2022 (53 years). I split by time.
Which years should be the training set, and which years the test set?

Reply with exactly two lines:
Train: first year to last year
Test: first year to last year
```
**Ollama answer:**
```text
Train: 1970 to 2015
Test: 2016 to 2022
```
**My check:** Suggested train 1970-2015, test 2016-2022 (7 test years). I kept train to 2012, test 2013-2022 (10 years, all in the slow-growth period). I ran the pipeline with its split too: same model ranking, so the choice does not change the conclusion.


### Topic 6: Interpreting the results  (answers saved 18:14)

#### Prompt 6a: Why does the naive guess beat the line?
**Result:** Used (partly)

**Prompt sent:**
```text
Answer immediately. No thinking.

I tested two models that predict China's CO2 emissions on the years 2013 to 2022.
Model A guesses "same as last year". Its average error (MAPE) is 2.2%.
Model B is a straight-line trend. Its average error (MAPE) is 26.3%.

Which model is better? Why can the simple Model A beat the trend Model B?
Answer in 2 short sentences.
```
**Ollama answer:**
```text
Model A is better because it captures a consistent pattern (autocorrelation), while Model B assumes a simple, continuous linear trend that may not fully describe the complex, fluctuating emissions data.
```
**My check:** Verdict right (2.24% vs 26.31% in `model_metrics.csv`). Reason partly misleading: "complex, fluctuating data". The real cause is a change in trend around 2013.

#### Prompt 6b: Why did the trend model fail?
**Result:** Used (premise wrong)

**Prompt sent:**
```text
Answer immediately. No thinking.

China's CO2 emissions grew by 8.7% a year from 2001 to 2013, then by only 1.7% a year from 2014 to 2022.
I trained a straight-line trend model on the years 1970 to 2012, then tested it on 2013 to 2022.
The model's predictions kept rising fast, but the real values rose slowly.

Explain in 3 short bullet points why the trend model made large errors on the test years.
```
**Ollama answer:**
```text
* **Non-linear Growth:** The emissions growth rate changed significantly (8.7% to 1.7%), which a single straight line cannot capture accurately.
* **Trend Mismatch:** The model extrapolated the high growth rate observed earlier, leading to overestimation in the later, slower growth period.
```
**My check:** The growth change (8.7% to 1.7% a year) matches `period_summary.csv`. **But it repeated a wrong premise I wrote in the prompt:** it said the model over-predicts. The straight line **under-predicts** by about 3,000 Mt (2022: predicted 9,284, real 12,519). Found later in the predictions file.

#### Prompt 6b2: Third bullet of 6b (follow-up)
**Result:** Used

**Prompt sent:**
```text
Write only the third bullet point, no thinking. Explain in one sentence why a trend model cannot predict a change in direction.
```
**Ollama answer:**
```text
A simple trend model assumes the rate of change remains constant, making it unable to predict a fundamental shift or deceleration in the underlying growth dynamic.
```
**My check:** The reply to 6b stopped at bullet 2, so I asked for the third separately. Correct: a trend model assumes a constant rate of change.

#### Prompt 6c: What does a negative R2 mean?
**Result:** Used

**Prompt sent:**
```text
Answer immediately. No thinking.

A prediction model has an R squared (R2) of -3.3 on the test data.

Explain in 2 short sentences, in simple English:
1. What does a negative R2 mean?
2. Is the model better or worse than just guessing the average of the test values?
```
**Ollama answer:**
```text
A negative R-squared means your model performs worse than simply guessing the average of the actual values. The model is actively making poorer predictions than a baseline guess.
```
**My check:** Correct: worse than guessing the average of the actual values. Matches my result (R2 = -18.42).


### Topic 7: Wording the presentation  (answers saved 18:15)

#### Prompt 7a: Two slide bullets on growth by period
**Result:** Used (bullet 1 only)

**Prompt sent:**
```text
Answer immediately. No thinking.

China's CO2 emissions, in million tonnes (Mt):
1970-1977: grew 6.2% a year, from 912 to 1,385 Mt
1978-2000: grew 4.3% a year, from 1,385 to 3,648 Mt
2001-2013: grew 8.7% a year, from 3,648 to 10,760 Mt
2014-2022: grew 1.7% a year, from 10,760 to 12,519 Mt

Write exactly 2 bullet points for a presentation slide, for a non-technical audience.
Each bullet under 20 words. Use only the numbers above.
```
**Ollama answer:**
```text
* Emissions grew rapidly from 1,385 Mt in 1977 to 10,760 Mt in 2013.
* Growth slowed significantly after 2013, increasing by only 1.7% between 2014 and 2022.
```
**My check:** Bullet 1 correct (1,385 Mt in 1977 to 10,760 Mt in 2013). Bullet 2 changed "1.7% a year" into "1.7% between 2014 and 2022"; the real total rise is about 16.3%. I wrote the bullet myself.

#### Prompt 7b: One slide bullet on the forecast
**Result:** Not used

**Prompt sent:**
```text
Answer immediately. No thinking.

Forecast of China's CO2 emissions in 2025, in million tonnes (Mt), under three scenarios:
- No growth: 12,519 Mt
- Recent growth (1.7% a year): 13,166 Mt
- Fast growth like the 2000s: 15,615 Mt

Write exactly 1 bullet point for a presentation slide, for a non-technical audience.
Under 25 words. Use only the numbers above.
```
**Ollama answer:**
```text
2025 predictions range from 12,519 Mt (No growth) to 15,615 Mt (Fast growth).
```
**My check:** Numbers matched the forecast file, but I later removed the forecast: a line fitted on all years predicts 11,457 Mt for 2023, below the real 2022 value of 12,519 Mt.


---

## 5. Change after the prompts: one simple model

After the prompts I simplified the project to **one regression model**, a straight line (slope and intercept from the least-squares formula). I removed the matrix method, the quadratic, exponential and lag models and the forecast. The "same as last year" guess stays as a comparison.

| Entry | Still describes the final code? |
|---|---|
| Topics 1, 2, 4, 5 | Yes. Cleaning, metrics and the split are unchanged. |
| Topic 3 | Partly. The matrix code was removed; what still applies is that Go and numpy coefficients are identical and `TestFitLineExact` recovers a known line. |
| Topic 6 | Only for the two remaining models. **Correct explanation:** the line adds a fixed +187 Mt a year, but real emissions curved upward in the 2000s, so it ends about 3,000 Mt too low. R2 is very negative mainly because the test years vary little (10,750 to 12,552 Mt). |
| Topic 7 | 7a yes; 7b no (forecast removed). Kept as the original record. |

## 6. How I checked the code

| Check | Result |
|---|---|
| `go test -v`, 9 tests (`docs/test_output.txt`) | All pass |
| Line fitted with Python numpy (`validation/crosscheck_numpy.py`, output in `docs/numpy_crosscheck_output.txt`) | Identical to 4 decimal places |
| Real data and model (`TestRealData`, `TestRealModel`) | 53 rows, 1970-2022, no gaps; test errors 2.24% (naive) and 26.31% (line) |
| Full run (`docs/pipeline_output.txt`) | Matches the CSV files in `data/processed/` |

## 7. Responsible-use checklist

- [ ] I recorded every AI tool I used, including Claude and Ollama.
- [ ] I did not paste any personal or confidential data into an AI tool (the dataset is public).
- [ ] I tested AI-suggested code before using it and kept the test results.
- [ ] I checked every number in my slides against the CSV outputs.
- [ ] I can explain every line of the Go code in my own words.
- [ ] Where an AI answer was wrong, I noted it and how I found out.

## 8. Reflection

I used Ollama (`gemma4:e2b`) as a reviewer for my Go code, not as the author. It was useful for step-by-step maths and it found two real crashes in my cleaning code, but it also said "No problems found" about code that had two problems, gave two wrong unit tests, and repeated a wrong premise from my own prompt. Because of this I did not trust any answer until I ran a test or checked it against the CSV files. A small model worked best with short, specific prompts that had one task each and the real numbers typed in. In the end I made my own decisions: a time-based split with test years 2013-2022, a year-range check in the cleaning code, one simple regression model, and no forecast. The main limit of the project is that the data holds only CO2 by year, so the model describes the trend but cannot explain it or link it to heat or health.
