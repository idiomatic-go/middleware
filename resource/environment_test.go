package resource_test

import (
	"fmt"
	"github.com/idiomatic-go/middleware/resource"
)

func ExampleDevEnv() {
	fmt.Println(resource.IsDevEnv())
	resource.SetEnv("dev")
	fmt.Println(resource.IsDevEnv())
	resource.SetEnv("devrrr")
	fmt.Println(resource.IsDevEnv())

	// Output:
	// true
	// true
	// false
}

func ExampleDevEnvOverride() {
	resource.OverrideIsDevEnv(func() bool { return false })
	fmt.Println(resource.IsDevEnv())
	resource.SetEnv("dev")
	fmt.Println(resource.IsDevEnv())
	resource.SetEnv("devrrr")
	fmt.Println(resource.IsDevEnv())

	// Output:
	// false
	// false
	// false
}

/*
func ExampleProdEnv() {
	fmt.Println(resource.IsProdEnv())
	os.Setenv(resource.RuntimeEnvKey, "prod")
	fmt.Println(resource.IsProdEnv())
	os.Setenv(resource.RuntimeEnvKey, "production")
	fmt.Println(resource.IsProdEnv())

	// Output:
	// false
	// true
	// false
}

func ExampleReviewEnv() {
	fmt.Println(resource.IsReviewEnv())
	os.Setenv(resource.RuntimeEnvKey, "review")
	fmt.Println(resource.IsReviewEnv())
	os.Setenv(resource.RuntimeEnvKey, "revvrrr")
	fmt.Println(resource.IsReviewEnv())

	// Output:
	// false
	// true
	// false
}

func ExampleStageEnv() {
	fmt.Println(resource.IsStageEnv())
	os.Setenv(resource.RuntimeEnvKey, "stage")
	fmt.Println(resource.IsStageEnv())
	os.Setenv(resource.RuntimeEnvKey, "")
	fmt.Println(resource.IsStageEnv())

	// Output:
	// false
	// true
	// false
}

func ExampleTestEnv() {
	fmt.Println(resource.IsTestEnv())
	os.Setenv(resource.RuntimeEnvKey, "test")
	fmt.Println(resource.IsTestEnv())
	os.Setenv(resource.RuntimeEnvKey, "atvrrr")
	fmt.Println(resource.IsTestEnv())

	// Output:
	// false
	// true
	// false
}


*/
