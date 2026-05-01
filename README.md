# 体脂率计算器

### 编译方法
```bash
git clone https://github.com/zyoung11/BodyFatCalculation.git

go mod tidy

go install fyne.io/fyne/v2/cmd/fyne@latest

fyne package -os android --app-id com.example.bodyfat --name "体脂率计算器"
```
