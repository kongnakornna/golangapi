<img src="https://r2cdn.perplexity.ai/pplx-full-logo-primary-dark%402x.png" style="height:64px;margin-right:32px"/>

# Branch  Enchant Code

คำว่า **"Enchant Code"** น่าจะเป็นการพิมพ์ผิดจากคำว่า **"Enhancement Code"** (การปรับปรุงโค้ดให้ดีขึ้น) หรืออาจหมายถึงการเขียนโค้ดให้ **"Clean" (Clean Code)** และมีคุณภาพสูง

ในบริบทของ Git Branching สำหรับทีมใหญ่ มักจะมีการแยกประเภท Branch สำหรับ "งานปรับปรุง" ออกจาก "ฟีเจอร์ใหม่" เพื่อให้ง่ายต่อการตรวจสอบ (Code Review) และการจัดลำดับความสำคัญ (Prioritization) ดังนี้ครับ:

***

### **1. Enhancement Branch (`enhancement/xxx`)**

ใช้สำหรับงานที่ **"ไม่ใช่ฟีเจอร์ใหม่ แต่ทำให้ระบบดีขึ้น"** เช่น การปรับจูน Performance, ปรับปรุง UI/UX เล็กน้อย หรือเพิ่ม Logging

* **ความแตกต่างจาก Feature:** Feature คือสิ่งที่ User "เห็นและใช้งานได้ใหม่" ส่วน Enhancement มักเป็นการปรับปรุงของเดิมให้ดีกว่าเดิม
* **Base Branch:** แตกออกจาก `develop`
* **Pattern:** `enhancement/<TICKET-ID>-<description>`
* **ตัวอย่าง:**
    * `enhancement/CART-105-improve-loading-speed` (ทำให้โหลดเร็วขึ้น)
    * `enhancement/UI-202-adjust-button-shadow` (ปรับเงาปุ่มให้สวยขึ้น)

***

### **2. Refactor Branch (`refactor/xxx`)**

ใช้สำหรับงาน **"รื้อโครงสร้างโค้ด (Refactoring)"** โดยที่ **"ผลลัพธ์การทำงานต้องเหมือนเดิม"** (User ไม่เห็นความเปลี่ยนแปลง แต่โค้ดอ่านง่ายขึ้น บำรุงรักษาง่ายขึ้น)

* **ความสำคัญ:** ทีมใหญ่แยก Branch นี้ออกมาเพื่อบอก Reviewer ว่า *"ไม่ต้อง Test ฟังก์ชันนะ แค่ดู Logic ว่าเขียนดีขึ้นไหม"*
* **Base Branch:** แตกออกจาก `develop`
* **Pattern:** `refactor/<TICKET-ID>-<scope>`
* **ตัวอย่าง:**
    * `refactor/USER-300-clean-auth-service` (จัดระเบียบโค้ดใน Auth Service)
    * `refactor/CORE-404-remove-unused-imports` (ลบโค้ดที่ไม่ได้ใช้ออก)

***

### **3. Chore Branch (`chore/xxx`)**

ใช้สำหรับงาน **"งานบ้าน/งานจุกจิก"** ที่ไม่กระทบ Code หลัก เช่น อัปเกรด Library, แก้ไฟล์ Config, หรือเขียน Document

* **Base Branch:** แตกออกจาก `develop`
* **Pattern:** `chore/<TICKET-ID>-<task>`
* **ตัวอย่าง:**
    * `chore/DEVOPS-500-update-nestjs-v10` (อัปเกรด Version Framework)
    * `chore/DOC-101-update-readme` (แก้ไฟล์คู่มือ)

***

### **สรุปตารางเปรียบเทียบ "Enchant" (Improvement) Branches**

| ประเภท Branch | ความหมาย | ผลกระทบต่อ User | ต้องเขียน Test เพิ่มไหม? |
| :-- | :-- | :-- | :-- |
| `feature/xxx` | ของใหม่ | **เห็น** และใช้งานได้ | ✅ ต้องมี Unit/E2E Test |
| `enhancement/xxx` | ของเดิมที่ดีขึ้น | **เห็น** (ทำงานดีขึ้น/เร็วขึ้น) | ✅ อาจต้องแก้ Test เดิม |
| `refactor/xxx` | จัดระเบียบโค้ด | **ไม่เห็น** (ทำงานเหมือนเดิมเป๊ะ) | ❌ ไม่ควรแก้ Test (ถ้า Logic ไม่เปลี่ยน) |
| `chore/xxx` | งานจุกจิก/Config | **ไม่เห็น** | ❌ ไม่ต้องเขียน Test |

### **คำแนะนำเพิ่มเติมสำหรับทีมใหญ่**

ถ้าทีมของคุณต้องการเน้นเรื่อง "Enchant Code" (ทำให้โค้ดดูขลัง/เทพขึ้น 🧙‍♂️) แนะนำให้เพิ่มขั้นตอน **Automated Code Quality Check** ใน Pipeline ก่อน Merge:

1. **Linting:** บังคับใช้กฎการเขียนโค้ด (ESLint/Prettier)
2. **SonarQube:** สแกนหา "Code Smell" หรือจุดที่เขียนไม่ดี
3. **Commit Message Lint:** บังคับให้ใส่ชื่อ Branch ประเภทนี้ใน Commit Message เช่น `refactor: clean up user service` เพื่อให้ Log อ่านง่าย
<span style="display:none">[^1][^10][^11][^12][^13][^14][^15][^16][^17][^18][^19][^2][^20][^3][^4][^5][^6][^7][^8][^9]</span>

<div align="center">⁂</div>

[^1]: https://git-scm.com/book/en/v2/Git-Branching-Branches-in-a-Nutshell

[^2]: https://git-scm.com/book/en/v2/Git-Branching-Basic-Branching-and-Merging

[^3]: https://git-scm.com/book/en/v2/Git-Branching-Branch-Management

[^4]: https://git-scm.com/book/en/v2/Git-Branching-Remote-Branches

[^5]: https://www.atlassian.com/git/tutorials/comparing-workflows/feature-branch-workflow

[^6]: https://graphite.com/guides/git-branch-naming-conventions

[^7]: https://docs.aws.amazon.com/prescriptive-guidance/latest/choosing-git-branch-approach/branches-in-a-gitflow-strategy.html

[^8]: https://dev.to/karmpatel/git-branching-strategies-a-comprehensive-guide-24kh

[^9]: https://stackoverflow.blog/2021/04/05/a-look-under-the-hood-how-branches-work-in-git/

[^10]: https://www.geeksforgeeks.org/git/introduction-to-git-branch/

[^11]: https://stackoverflow.com/questions/273695/what-are-some-examples-of-commonly-used-practices-for-naming-git-branches

[^12]: https://www.facebook.com/groups/ThaiPGAssociateSociety/posts/2939005346310715/

[^13]: https://www.w3schools.com/git/git_branch.asp

[^14]: https://gist.github.com/digitaljhelms/4287848

[^15]: https://www.toptal.com/git/enhanced-git-flow-explained

[^16]: https://nvie.com/posts/a-successful-git-branching-model/

[^17]: https://www.abtasty.com/blog/git-branching-strategies/

[^18]: https://www.atlassian.com/git/tutorials/comparing-workflows/gitflow-workflow

[^19]: https://git-scm.com/book/en/v2/Git-Branching-Branching-Workflows

[^20]: https://www.datacamp.com/tutorial/git-branching-strategy-guide

