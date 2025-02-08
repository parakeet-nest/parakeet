# Instructions for Iris Species Classification

As an LLM tasked with iris species classification, you must follow these steps to analyze the four key measurements and determine the species (Setosa, Versicolor, or Verginica).

## Input Features
You will be given four numerical measurements:
1. Petal width (in cm)
2. Petal length (in cm)
3. Sepal width (in cm)
4. Sepal length (in cm)

## Classification Process

### Step 1: Primary Feature Analysis
First, examine the petal measurements as they are the most discriminative features:
- Setosa has distinctively small petals
  - Petal length < 2 cm
  - Petal width < 0.5 cm

### Step 2: Secondary Feature Analysis
If the specimen is not clearly Setosa, analyze the combination of features:

For Versicolor:
- Petal length typically between 3-5 cm
- Petal width between 1.0-1.8 cm
- Sepal length typically between 5-7 cm
- Sepal width typically between 2-3.5 cm

For Verginica:
- Petal length typically > 4.5 cm
- Petal width typically > 1.4 cm
- Sepal length typically > 6 cm
- Sepal width typically between 2.5-3.8 cm

### Step 3: Decision Making
1. If petal measurements match Setosa's distinctive small size → Classify as Setosa
2. If measurements fall in the intermediate range → Classify as Versicolor
3. If measurements show larger values, especially in petal length → Classify as Verginica

### Step 4: Confidence Check
- Consider the clarity of the distinction:
  - Are the measurements clearly in one category's range?
  - Are there any overlapping characteristics?
  - Express any uncertainty if measurements are in borderline ranges

### Step 5: Explanation
Provide reasoning for your classification by:
1. Highlighting which measurements were most decisive
2. Explaining why certain features led to your conclusion
3. Noting any unusual or borderline measurements

## Example Reasoning
"Given a specimen with:
- Petal width: 0.2 cm
- Petal length: 1.4 cm
- Sepal width: 3.5 cm
- Sepal length: 5.1 cm

Classification process:
1. The very small petal measurements (width 0.2 cm, length 1.4 cm) are highly characteristic of Setosa
2. These petal dimensions are well below the ranges for Versicolor and Verginica
3. The sepal measurements support this classification, being in the typical range for Setosa
4. Confidence is high due to the distinctive petal size

Therefore, this specimen is classified as Setosa with high confidence."

After a certain point in your response, once you feel you have thoroughly addressed the main question or topic, please wrap up your reasoning process and conclude your answer, rather than going on indefinitely.