# Prompt Engineering: โครงสร้างและการเขียน
##  ภาษาอังกฤษ

### Basic Prompt Structure
```
Role + Task + Context + Specifications + Output Format
```

### Key Components
1. **Role Definition**
   ```
   "You are an expert in digital marketing..."
   "As a data scientist with 10 years of experience..."
   ```

2. **Clear Task**
   ```
   "Write a product description for..."
   "Analyze the following dataset and identify trends..."
   ```

3. **Context Provision**
   ```
   "For a startup targeting Gen Z consumers..."
   "In an academic research context..."
   ```

4. **Detailed Specifications**
   ```
   "Use simple language suitable for beginners"
   "Include 5 key takeaways"
   "Limit to 500 words"
   ```

5. **Output Format**
   ```
   "Format as a JSON object"
   "Create a markdown table"
   "Structure as an executive summary"
   ```

### Example English Prompt
```
"As a financial analyst, create an investment risk assessment for renewable energy stocks. 
Consider market volatility, regulatory changes, and technological disruption. 
Present in a structured report with: 1) Executive summary, 2) Risk categories, 
3) Mitigation strategies, 4) Recommendations. Use professional tone and include data points where relevant."
```

### Prompt Writing Techniques
1. **Zero-shot Prompting**
   ```
   "Translate this paragraph to French."
   ```

2. **Few-shot Prompting** (providing examples)
   ```
   "Example 1: [input] -> [output]
    Example 2: [input] -> [output]
    Now process this: [new input]"
   ```

3. **Chain-of-Thought Prompting**
   ```
   "Explain your reasoning step by step."
   "Let's think through this problem systematically."
   ```

### Best Practices
- Be specific and unambiguous
- Use delimiters for complex inputs
- Specify the desired length and depth
- Iterate and refine based on results
- Break complex tasks into subtasks

### Advanced Techniques
- **Temperature setting** (for creativity vs. consistency)
- **System prompts** for setting behavior parameters
- **Template prompts** for reproducible results
- **Meta-prompts** for generating better prompts

