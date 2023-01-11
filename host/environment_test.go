package host_test

import (
	"fmt"
	"github.com/idiomatic-go/middleware/resource"
)

func ExampleDevEnv() {
	fmt.Printf("test: IsDevEnv() -> %v\n", resource.IsDevEnv())

	resource.SetEnv("dev")
	fmt.Printf("test: IsDevEnv(dev) -> %v\n", resource.IsDevEnv())

	resource.SetEnv("devrrr")
	fmt.Printf("test: IsDevEnv(devrrr) -> %v\n", resource.IsDevEnv())

	// Output:
	//test: IsDevEnv() -> true
	//test: IsDevEnv(dev) -> true
	//test: IsDevEnv(devrrr) -> false

}

func ExampleDevEnvOverride() {
	resource.SetEnvironmentMatcher(func(int) bool { return false })

	fmt.Printf("test: IsDevEnv() -> %v\n", resource.IsDevEnv())

	resource.SetEnv("dev")
	fmt.Printf("test: IsDevEnv(dev) -> %v\n", resource.IsDevEnv())

	resource.SetEnv("devrrr")
	fmt.Printf("test: IsDevEnv(devrrr) -> %v\n", resource.IsDevEnv())

	// Output:
	//test: IsDevEnv() -> false
	//test: IsDevEnv(dev) -> false
	//test: IsDevEnv(devrrr) -> false

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
