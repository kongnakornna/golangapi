# Go Code Review Template - AI Skill
 
**Go Code Review** is a comprehensive quality assurance process for Go code that covers 8 critical areas:

1. **Error Handling** - Ensuring all errors are properly handled and propagated
2. **SQL Injection** - Using parameterized queries to prevent security vulnerabilities
3. **Panic Check** - Preventing and gracefully handling runtime panics
4. **Unclosed Transactions/Resources** - Properly closing all resources after use
5. **Memory Leak** - Identifying and preventing memory leaks
6. **Context Propagation** - Correctly propagating context through the call chain
7. **Debug Code Cleanup** - Removing debug code before production deployment
8. **Unit Test Presence** - Ensuring adequate test coverage for new/changed code

**Benefits:**
- Improved system stability and reliability
- Enhanced security posture
- Reduced production incidents
- Better maintainability and code quality
- Faster debugging and issue resolution

**Best Practices:**
- Always review code before merging to dev branch
- Address critical issues before deployment
- Maintain comprehensive test coverage
- Document lessons learned from incidents

---

### ภาษาไทย

**Go Code Review** เป็นกระบวนการตรวจสอบคุณภาพโค้ด Go ที่ครอบคลุม 8 หัวข้อสำคัญ:

1. **Error Handling** - ตรวจสอบการจัดการ error ที่ถูกต้องและครบถ้วน
2. **SQL Injection** - ใช้ parameterized query เพื่อป้องกันช่องโหว่ด้านความปลอดภัย
3. **Panic Check** - ป้องกันและจัดการ runtime panic อย่างเหมาะสม
4. **Unclosed Transactions/Resources** - ปิด resources ทั้งหมดหลังใช้งานเสร็จ
5. **Memory Leak** - ระบุและป้องกันการรั่วไหลของหน่วยความจำ
6. **Context Propagation** - ส่งต่อ Context อย่างถูกต้องใน call chain
7. **Debug Code Cleanup** - ลบโค้ด debug ก่อน deploy ไป production
8. **Unit Test Presence** - ตรวจสอบความครอบคลุมของ test สำหรับโค้ดที่เพิ่ม/แก้ไข

**ประโยชน์ที่ได้รับ:**
- ระบบมีความเสถียรและน่าเชื่อถือมากขึ้น
- เพิ่มความปลอดภัยของระบบ
- ลดปัญหาที่เกิดขึ้นใน production
- โค้ดดูแลรักษาง่ายและมีคุณภาพสูงขึ้น
- แก้ไขปัญหาได้รวดเร็วขึ้น

**แนวปฏิบัติที่ดี:**
- ตรวจสอบโค้ดก่อน merge ไป dev branch ทุกครั้ง
- แก้ไขปัญหาที่สำคัญก่อน deploy
- รักษาความครอบคลุมของ test
- บันทึกบทเรียนที่ได้จากเหตุการณ์ต่างๆ

---

## 🔧 Pseudo Code for AI Agent

```python
# Go Code Review AI Agent - Pseudo Code

class GoCodeReviewAgent:
    def __init__(self):
        self.review_categories = [
            "Error Handling",
            "SQL Injection",
            "Panic Check",
            "Unclosed Transactions/Resources",
            "Memory Leak",
            "Context Propagation",
            "Debug Code Cleanup",
            "Unit Test Presence"
        ]
        self.critical_findings = []
        self.warning_findings = []
    
    def parse_arguments(self, args):
        """
        Parse user arguments to determine review scope
        
        Args:
            args: Command line arguments
                  - path/to/file.go
                  - path/to/file.go::FunctionName
                  - None (review all changed files)
        Returns:
            dict: {
                'file': str, 
                'function': str or None,
                'mode': 'file' | 'function' | 'branch'
            }
        """
        if not args:
            return {'mode': 'branch'}
        
        if '::' in args:
            file_path, func_name = args.split('::')
            return {'mode': 'function', 'file': file_path, 'function': func_name}
        
        return {'mode': 'file', 'file': args}
    
    def get_changed_files(self):
        """
        Get all changed .go files in current branch
        
        Returns:
            list: List of file paths
        """
        # git rev-parse --abbrev-ref HEAD
        # git diff dev...HEAD --name-only
        # Filter for .go files excluding _test.go
        pass
    
    def get_file_diff(self, file_path, base_branch='dev'):
        """
        Get diff for a specific file
        
        Args:
            file_path: Path to file
            base_branch: Base branch name
        Returns:
            str: Diff content
        """
        # git diff --staged -- {file_path}
        # If empty: git diff {base_branch}...HEAD -- {file_path}
        pass
    
    def extract_function_code(self, file_content, function_name):
        """
        Extract only the code for a specific function
        
        Args:
            file_content: Full file content
            function_name: Name of function to extract
        Returns:
            str: Function code block
        """
        # Find func {function_name} in file
        # Extract function body
        pass
    
    def review_error_handling(self, code):
        """Check for proper error handling"""
        findings = []
        # Check for ignored errors: _ = err
        # Check for nil error handling
        # Check for missing error context
        return findings
    
    def review_sql_injection(self, code):
        """Check for SQL injection vulnerabilities"""
        findings = []
        # Check for fmt.Sprintf in SQL queries
        # Check for string concatenation in SQL
        # Check for raw SQL without parameters
        return findings
    
    def review_panic_check(self, code):
        """Check for panic risks"""
        findings = []
        # Check for unprotected panic()
        # Check for nil pointer dereference
        # Check for type assertion without ok check
        return findings
    
    def review_unclosed_resources(self, code):
        """Check for unclosed resources"""
        findings = []
        # Check for db.Begin() without defer Rollback()
        # Check for http.Get() without defer Body.Close()
        # Check for os.Open() without defer Close()
        # Check for db.Query() without rows.Close()
        return findings
    
    def review_memory_leak(self, code):
        """Check for memory leak risks"""
        findings = []
        # Check for goroutines without lifecycle control
        # Check for tickers without Stop()
        # Check for timers without Stop()
        # Check for channels without drain/close
        return findings
    
    def review_context_propagation(self, code):
        """Check for proper context propagation"""
        findings = []
        # Check for context.Background() usage
        # Check for context.TODO() usage
        # Check for missing WithContext(ctx) in GORM
        # Check for context not passed to downstream calls
        return findings
    
    def review_debug_code(self, code):
        """Check for debug code in production"""
        findings = []
        # Check for fmt.Println, fmt.Printf
        # Check for unnecessary log.Println
        # Check for commented-out code
        # Check for TODO comments
        # Check for hardcoded test values
        return findings
    
    def review_unit_tests(self, file_path, changes):
        """Check for unit test coverage"""
        findings = []
        # Check if _test.go exists
        # Check if test covers new/changed functions
        # Check for happy path, error path, edge cases
        return findings
    
    def generate_report(self):
        """Generate final review report"""
        report = "## สรุปผลการ Code Review\n\n"
        
        for category in self.review_categories:
            report += f"### {category}\n"
            # Add findings for this category
            report += "\n"
        
        # Add summary statistics
        report += "### สรุปภาพรวม\n"
        report += f"- 🔴 Critical: {len(self.critical_findings)}\n"
        report += f"- 🟡 Warning: {len(self.warning_findings)}\n"
        
        return report
    
    def run(self, args):
        """Main execution flow"""
        # Parse arguments
        review_scope = self.parse_arguments(args)
        
        # Get code to review
        if review_scope['mode'] == 'branch':
            files = self.get_changed_files()
            for file_path in files:
                diff = self.get_file_diff(file_path)
                self.review_code(diff, file_path)
        elif review_scope['mode'] == 'function':
            file_path = review_scope['file']
            func_name = review_scope['function']
            diff = self.get_file_diff(file_path)
            func_code = self.extract_function_code(diff, func_name)
            self.review_code(func_code, file_path, func_name)
        else:  # file mode
            file_path = review_scope['file']
            diff = self.get_file_diff(file_path)
            self.review_code(diff, file_path)
        
        # Generate report
        return self.generate_report()
```

---

## 📝 Quick Reference Card

| Priority | Category | Check | Fix |
|----------|----------|-------|-----|
| 🔴 | Error Handling | `_ = err` | `if err != nil { return fmt.Errorf("...: %w", err) }` |
| 🔴 | SQL Injection | `fmt.Sprintf` in SQL | Use parameterized queries |
| 🔴 | Panic Check | `panic()` in business logic | Use error handling instead |
| 🔴 | Resources | Missing `defer` for Close | Add `defer resource.Close()` |
| 🟡 | Memory Leak | Goroutine without context | Use context for lifecycle control |
| 🟡 | Context | `context.Background()` | Use context from caller |
| 🟡 | Debug | `fmt.Println()` | Remove or use proper logging |
| 🟡 | Tests | Missing `_test.go` | Write unit tests |

---

## 🎯 Success Metrics

| Metric | Target | Measurement |
|--------|--------|-------------|
| Critical Issues Found | > 0 before merge | Number of critical issues reported |
| Warnings Found | > 0 before merge | Number of warnings reported |
| Code Quality Improvement | ↑ 20% | Reduce production incidents |
| Review Time | < 5 min per file | Time to complete review |
| Test Coverage | ↑ 80% | Percentage of code covered by tests |

---
 