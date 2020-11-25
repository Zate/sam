module github.com/zate/sam

replace github.com/zate/sam => ../sam

replace github.com/zate/sam/pkg/sbase => ../sam/pkg/sbase

go 1.15

require (
	github.com/briandowns/spinner v1.11.1
	github.com/joho/godotenv v1.3.0
	github.com/sirupsen/logrus v1.7.0
)
