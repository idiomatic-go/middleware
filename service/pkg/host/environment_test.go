package host

import (
	"fmt"
	"os"
)

func Example_Match() {
	fmt.Printf("test: matchEnvironment() -> [dev:%v] [review:%v] [test:%v] [stage:%v] [prod:%v] \n",
		isDevEnv(), isReviewEnv(), isTestEnv(), isStageEnv(), isProdEnv())

	os.Setenv(runtimeEnvKey, "test")
	fmt.Printf("test: matchEnvironment() -> [dev:%v] [review:%v] [test:%v] [stage:%v] [prod:%v] \n",
		isDevEnv(), isReviewEnv(), isTestEnv(), isStageEnv(), isProdEnv())

	os.Setenv(runtimeEnvKey, "prod")
	fmt.Printf("test: matchEnvironment() -> [dev:%v] [review:%v] [test:%v] [stage:%v] [prod:%v] \n",
		isDevEnv(), isReviewEnv(), isTestEnv(), isStageEnv(), isProdEnv())

	os.Setenv(runtimeEnvKey, "invalid")
	fmt.Printf("test: matchEnvironment() -> [dev:%v] [review:%v] [test:%v] [stage:%v] [prod:%v] \n",
		isDevEnv(), isReviewEnv(), isTestEnv(), isStageEnv(), isProdEnv())

	//Output:
	//test: matchEnvironment() -> [dev:true] [review:false] [test:false] [stage:false] [prod:false]
	//test: matchEnvironment() -> [dev:false] [review:false] [test:true] [stage:false] [prod:false]
	//test: matchEnvironment() -> [dev:false] [review:false] [test:false] [stage:false] [prod:true]
	//test: matchEnvironment() -> [dev:false] [review:false] [test:false] [stage:false] [prod:false]
	
}
